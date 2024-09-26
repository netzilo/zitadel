package query

import (
	"context"
	"fmt"
	"strings"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/cache"
	"github.com/zitadel/zitadel/internal/cache/gomap"
	"github.com/zitadel/zitadel/internal/cache/noop"
	"github.com/zitadel/zitadel/internal/eventstore"
)

type Caches struct {
	connectors *cacheConnectors
	instance   cache.Cache[instanceIndex, string, *authzInstance]
}

func startCaches(background context.Context, conf *cache.CachesConfig) (_ *Caches, err error) {
	caches := &Caches{
		instance: noop.NewCache[instanceIndex, string, *authzInstance](),
	}
	if conf == nil {
		return caches, nil
	}
	caches.connectors, err = startCacheConnectors(background, conf)
	if err != nil {
		return nil, err
	}
	caches.instance, err = startCache[instanceIndex, string, *authzInstance](background, instanceIndexValues(), "authz_instance", conf.Instance, caches.connectors)
	if err != nil {
		return nil, err
	}
	caches.registerInstanceInvalidation()

	return caches, nil
}

type cacheConnectors struct {
	memory *cache.AutoPruneConfig
	// pool   *pgxpool.Pool
}

func startCacheConnectors(_ context.Context, conf *cache.CachesConfig) (*cacheConnectors, error) {
	connectors := new(cacheConnectors)
	if conf.Connectors.Memory.Enabled {
		connectors.memory = &conf.Connectors.Memory.AutoPrune
	}

	return connectors, nil
}

func startCache[I, K comparable, V cache.Entry[I, K]](background context.Context, indices []I, name string, conf *cache.CacheConfig, connectors *cacheConnectors) (cache.Cache[I, K, V], error) {
	if conf == nil || conf.Connector == "" {
		return noop.NewCache[I, K, V](), nil
	}
	if strings.EqualFold(conf.Connector, "memory") && connectors.memory != nil {
		c := gomap.NewCache[I, K, V](background, indices, *conf)
		connectors.memory.StartAutoPrune(background, c, name)
		return c, nil
	}

	/* TODO
	if strings.EqualFold(conf.Connector, "sql") && connectors.pool != nil {
		return ...
	}
	*/

	return nil, fmt.Errorf("cache connector %q not enabled", conf.Connector)
}

type invalidator[I comparable] interface {
	Invalidate(ctx context.Context, index I, key ...string) error
}

func cacheInvalidationFunc[I comparable](cache invalidator[I], index I, getID func(*eventstore.Aggregate) string) func(context.Context, []*eventstore.Aggregate) {
	return func(ctx context.Context, aggregates []*eventstore.Aggregate) {
		ids := make([]string, len(aggregates))
		for i, aggregate := range aggregates {
			ids[i] = getID(aggregate)
		}
		err := cache.Invalidate(ctx, index, ids...)
		logging.OnError(err).Warn("cache invalidation failed")
	}
}

func getAggregateID(aggregate *eventstore.Aggregate) string {
	return aggregate.ID
}

func getResourceOwner(aggregate *eventstore.Aggregate) string {
	return aggregate.ResourceOwner
}
