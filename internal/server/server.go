package server

import (
	"context"
	"fmt"
	"net/http"
	"sync/atomic"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	_ "github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/zacscoding/go-rest-template/internal/config"
	"github.com/zacscoding/go-rest-template/internal/controller"
	"github.com/zacscoding/go-rest-template/internal/handler"
	"github.com/zacscoding/go-rest-template/internal/handler/middleware"
	"github.com/zacscoding/go-rest-template/internal/metrics"
	"github.com/zacscoding/go-rest-template/pkg/logging"
	"github.com/zacscoding/go-rest-template/pkg/version"
	"go.uber.org/fx"
)

type Server struct {
	apiserver    *http.Server
	metricserver *http.Server

	running      int32
	apiEngine    *gin.Engine
	metricEngine *gin.Engine

	conf           *config.Config
	mp             metrics.Provider
	authController *controller.AuthController
	userController *controller.UserController
}

func NewServer(
	lc fx.Lifecycle,
	conf *config.Config,
	mp metrics.Provider,
	authController *controller.AuthController,
	userController *controller.UserController,
) (*Server, error) {
	gin.SetMode(gin.ReleaseMode)
	srv := Server{
		conf:           conf,
		apiEngine:      gin.New(),
		mp:             mp,
		authController: authController,
		userController: userController,
	}

	// setup gin
	corscfg := cors.DefaultConfig()
	corscfg.AllowBrowserExtensions = conf.Server.Cors.BrowserExt
	corscfg.AllowAllOrigins = true
	if !conf.Server.Cors.AllowAll {
		corscfg.AllowAllOrigins = false
		corscfg.AllowOrigins = conf.Server.Cors.Origin
	}
	srv.apiEngine.Use(
		middleware.LoggingMiddleware("/healthz", "/version", "/metrics"),
		gin.Recovery(),
		cors.New(corscfg),
		middleware.RequestIDMiddleware(),
		middleware.TimeoutMiddleware(conf.Server.WriteTimeout),
		metrics.NewMiddleware(srv.mp, "/version", "/metrics"),
	)
	if conf.Server.Docs.Enabled {
		srv.apiEngine.StaticFile("/docs/docs.html", conf.Server.Docs.Path)
	}

	srv.apiserver = &http.Server{
		Addr:         fmt.Sprintf(":%d", conf.Server.Port),
		Handler:      srv.apiEngine,
		ReadTimeout:  conf.Server.ReadTimeout,
		WriteTimeout: conf.Server.WriteTimeout,
	}
	if conf.Metric.Enabled {
		if conf.Server.Port == conf.Metric.Port {
			srv.metricEngine = srv.apiEngine
		} else {
			srv.metricEngine = gin.New()
			srv.metricEngine.Use(gin.Recovery())
			srv.metricserver = &http.Server{
				Addr:    fmt.Sprintf(":%d", conf.Metric.Port),
				Handler: srv.metricEngine,
			}
		}
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logging.FromContext(ctx).Infof("Start to rest api server :%d", srv.conf.Server.Port)
			return srv.Start()
		},
		OnStop: func(ctx context.Context) error {
			logging.FromContext(ctx).Infof("Stopped rest api server")
			return srv.Stop(ctx)
		},
	})
	return &srv, nil
}

func (srv *Server) Start() error {
	if !atomic.CompareAndSwapInt32(&srv.running, 0, 1) {
		return errors.New("server already started")
	}
	go func() {
		err := srv.apiserver.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			logging.DefaultLogger().Fatalw("failed to close http server", "err", err)
		}
	}()
	if srv.metricserver != nil {
		go func() {
			err := srv.metricserver.ListenAndServe()
			if err != nil && err != http.ErrServerClosed {
				logging.DefaultLogger().Fatalw("failed to close http metric server", "err", err)
			}
		}()
	}
	return nil
}

func (srv *Server) Stop(ctx context.Context) error {
	if !atomic.CompareAndSwapInt32(&srv.running, 1, 0) {
		return errors.New("server already stopped")
	}

	var result error
	if err := srv.apiserver.Shutdown(ctx); err != nil {
		result = multierror.Append(result, fmt.Errorf("shutdown server. err: %v", err))
	}
	if srv.metricserver != nil {
		if err := srv.metricserver.Shutdown(ctx); err != nil {
			result = multierror.Append(result, fmt.Errorf("shutdown metric server. err: %v", err))
		}
	}
	return result
}

func (srv *Server) RouteAPI() error {
	if err := srv.routeAPI(); err != nil {
		return err
	}
	if err := srv.routeMetricAPI(); err != nil {
		return err
	}
	return nil
}

func (srv *Server) routeAPI() error {
	// Route common apis
	srv.apiEngine.GET("version", func(gctx *gin.Context) {
		gctx.JSON(http.StatusOK, version.Get())
	})

	// Route v1
	v1 := srv.apiEngine.Group("/api/v1")

	anonymousGroup := v1.Group("")
	anonymousGroup.POST("login", srv.authController.JWTMiddleware.LoginHandler)
	anonymousGroup.POST("signup", handler.Wrap(srv.userController.HandleSignUp))

	authGroup := v1.Group("")
	authGroup.Use(srv.authController.AuthMiddleware())

	userGroup := authGroup.Group("user")
	userGroup.POST("refresh-token", srv.authController.JWTMiddleware.RefreshHandler)
	userGroup.GET("me", handler.Wrap(srv.userController.HandleMe))
	return nil
}

func (srv *Server) routeMetricAPI() error {
	if !srv.conf.Metric.Enabled {
		return nil
	}
	srv.metricEngine.GET("metrics", gin.WrapH(promhttp.Handler()))
	return nil
}
