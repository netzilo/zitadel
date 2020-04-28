package eventsourcing

import (
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/cache"
	"github.com/caos/zitadel/internal/cache/config"
	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/usergrant/repository/eventsourcing/model"
)

type UserGrantCache struct {
	userGrantCache cache.Cache
}

func StartCache(conf *config.CacheConfig) (*UserGrantCache, error) {
	userGrantCache, err := conf.Config.NewCache()
	logging.Log("EVENT-vDneN").OnError(err).Panic("unable to create user cache")

	return &UserGrantCache{userGrantCache: userGrantCache}, nil
}

func (c *UserGrantCache) getUserGrant(ID string) *model.UserGrant {
	user := &model.UserGrant{ObjectRoot: models.ObjectRoot{AggregateID: ID}}
	if err := c.userGrantCache.Get(ID, user); err != nil {
		logging.Log("EVENT-4eTZh").WithError(err).Debug("error in getting cache")
	}
	return user
}

func (c *UserGrantCache) cacheUserGrant(grant *model.UserGrant) {
	err := c.userGrantCache.Set(grant.AggregateID, grant)
	if err != nil {
		logging.Log("EVENT-ThnBb").WithError(err).Debug("error in setting project cache")
	}
}
