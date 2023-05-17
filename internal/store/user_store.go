package store

import (
	"context"

	"github.com/zacscoding/go-rest-template/internal/config"
	"github.com/zacscoding/go-rest-template/internal/metrics"
	"github.com/zacscoding/go-rest-template/internal/model"
	"github.com/zacscoding/go-rest-template/pkg/cache"
	"github.com/zacscoding/go-rest-template/pkg/database"
	"github.com/zacscoding/go-rest-template/pkg/logging"
	"gorm.io/gorm"
)

var _ UserStore = (*userStore)(nil)

//go:generate mockery --name UserStore --filename user_store.go
type UserStore interface {
	// Save saves a given u user.
	Save(ctx context.Context, u *model.User) error

	// FindByEmail returns an user with given email if exists, otherwise database.ErrRecordNotFound.
	FindByEmail(ctx context.Context, email string) (*model.User, error)
}

func NewUserStore(conf *config.Config, db *gorm.DB, cacher cache.Cacher, mp metrics.Provider) (UserStore, error) {
	if cacher == nil {
		return &userStore{db: db, mp: mp}, nil
	}
	return newUserCacheStore(conf, cacher, mp, &userStore{db: db, mp: mp})
}

type userStore struct {
	db *gorm.DB
	mp metrics.Provider
}

func (s *userStore) Save(ctx context.Context, u *model.User) error {
	if err := database.FromContext(ctx, s.db).
		WithContext(ctx).
		Save(u).Error; err != nil {
		logging.FromContext(ctx).Errorw("failed to save an user", "err", err)
		return database.WrapError(err)
	}
	return nil
}

func (s *userStore) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var (
		db     = database.FromContext(ctx, s.db).WithContext(ctx)
		result model.User
	)
	if err := db.Where("email = ?", email).First(&result).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			logging.FromContext(ctx).Errorw("failed to find an user by email", "email", email, "err", err)
		}
		return nil, database.WrapError(err)
	}
	return &result, nil
}
