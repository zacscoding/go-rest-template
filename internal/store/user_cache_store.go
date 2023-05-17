package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/zacscoding/go-rest-template/internal/config"
	"github.com/zacscoding/go-rest-template/internal/metrics"
	"github.com/zacscoding/go-rest-template/internal/model"
	"github.com/zacscoding/go-rest-template/pkg/cache"
)

var _ UserStore = (*userCacheStore)(nil)

const (
	cacheKeyUserByEmail = "user-by-email"
)

type userCacheStore struct {
	cacher   cache.Cacher
	mp       metrics.Provider
	delegate UserStore
}

func newUserCacheStore(_ *config.Config,
	cacher cache.Cacher,
	mp metrics.Provider,
	delegate UserStore,
) (UserStore, error) {
	if cacher == nil {
		return nil, errors.New("require cacher")
	}
	return &userCacheStore{
		cacher:   cacher,
		mp:       mp,
		delegate: delegate,
	}, nil
}

func (uc *userCacheStore) Save(ctx context.Context, u *model.User) error {
	return uc.delegate.Save(ctx, u)
}

func (uc *userCacheStore) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var (
		item     model.User
		key      = uc.userByEmailKey(email)
		cacheHit = true
	)
	err := uc.cacher.Fetch(ctx, key, &item, func() (interface{}, error) {
		cacheHit = false
		return uc.delegate.FindByEmail(ctx, email)
	})
	if err != nil {
		return nil, err
	}
	uc.mp.RecordCache(cacheKeyUserByEmail, cacheHit)
	return &item, nil
}

func (uc *userCacheStore) userByEmailKey(email string) string {
	return fmt.Sprintf("%s.%s", cacheKeyUserByEmail, email)
}
