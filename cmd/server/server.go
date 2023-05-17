package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/spf13/cobra"
	"github.com/zacscoding/go-rest-template/internal/config"
	"github.com/zacscoding/go-rest-template/internal/controller"
	"github.com/zacscoding/go-rest-template/internal/metrics"
	"github.com/zacscoding/go-rest-template/internal/server"
	"github.com/zacscoding/go-rest-template/internal/store"
	"github.com/zacscoding/go-rest-template/pkg/cache"
	"github.com/zacscoding/go-rest-template/pkg/database"
	"github.com/zacscoding/go-rest-template/pkg/logging"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func runApplication(*cobra.Command, []string) {
	conf, err := config.Load(configPath, nil)
	if err != nil {
		log.Fatal(err)
	}

	// setup logging
	logging.SetConfig(&logging.Config{
		Encoding:          conf.Logging.Encoding,
		Level:             zapcore.Level(conf.Logging.Level),
		Development:       false,
		EncoderConfig:     logging.NewEncoderConfig(),
		DisableStacktrace: conf.Logging.DisableStacktrace,
	})
	defer logging.DefaultLogger().Sync()

	b, _ := json.MarshalIndent(conf, "", "    ")
	logging.DefaultLogger().Infof("Starting applicaltion server. configs: %s", string(b))

	// setup global components at here
	if tr, ok := http.DefaultTransport.(*http.Transport); ok {
		tr.MaxIdleConnsPerHost = tr.MaxIdleConns
	}

	runApplicationReal(conf)
}

func runApplicationReal(conf *config.Config) {
	fx.New(
		fx.Supply(conf),
		fx.Supply(&conf.DB),
		fx.Supply(&conf.Cache),
		fx.Supply(logging.DefaultLogger().Desugar()),
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log.Named("fx")}
		}),
		fx.StartTimeout(20*time.Second),
		fx.StopTimeout(conf.Server.GracefulShutdown+time.Second),
		fx.Provide(
			// setup metrics provider
			metrics.NewProvider,

			// setup database and stores
			database.Open,
			cache.NewCacher,
			store.NewUserStore,

			// setup controllers
			controller.NewAuthController,
			controller.NewUserController,

			server.NewServer,
		),
		fx.Invoke(
			func(srv *server.Server) error {
				return srv.RouteAPI()
			}),
	).Run()
}
