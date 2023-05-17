package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/zacscoding/go-rest-template/internal/config"
	"github.com/zacscoding/go-rest-template/internal/handler/apierr"
	"github.com/zacscoding/go-rest-template/internal/model"
	"github.com/zacscoding/go-rest-template/internal/store"
	"github.com/zacscoding/go-rest-template/pkg/database"
	"github.com/zacscoding/go-rest-template/pkg/logging"
	"github.com/zacscoding/go-rest-template/pkg/utils/authutil"
)

type UserController struct {
	conf      *config.Config
	userStore store.UserStore
}

func NewUserController(conf *config.Config, userStore store.UserStore) (*UserController, error) {
	return &UserController{
		conf:      conf,
		userStore: userStore,
	}, nil
}

type SignUpReq struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=5"`
}

// HandleSignUp handles "POST /api/v1/signup".
func (c *UserController) HandleSignUp(gctx *gin.Context) (interface{}, error) {
	var (
		ctx = gctx.Request.Context()
		req SignUpReq
	)
	if err := gctx.ShouldBind(&req); err != nil {
		return nil, apierr.ErrInvalidRequest.WithMessage(err.Error())
	}

	password, err := authutil.EncodePassword(req.Password, 0)
	if err != nil {
		logging.FromContext(ctx).Errorw("failed to encode password", "err", err)
		return nil, err
	}

	user := model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: password,
		Roles:    []string{model.RoleUser.String()},
		RolesMap: nil,
	}

	if err := c.userStore.Save(ctx, &user); err != nil {
		if err != database.ErrKeyConflict {
			return nil, err
		}
		return nil, apierr.ErrResourceConflict.WithMessagef("email %s already exists", req.Email)
	}
	return &user, nil
}

// HandleMe handles "GET /api/v1/user/me"
func (c *UserController) HandleMe(gctx *gin.Context) (interface{}, error) {
	currentUser := authutil.CurrentUser(gctx.Request.Context())
	return currentUser, nil
}
