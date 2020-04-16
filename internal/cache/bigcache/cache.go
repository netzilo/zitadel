package bigcache

import (
	"bytes"
	"encoding/gob"
	"time"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/errors"

	a_cache "github.com/allegro/bigcache"
)

type Bigcache struct {
	cache *a_cache.BigCache
}

func NewBigcache(c *Config) (*Bigcache, error) {
	cacheConfig := a_cache.DefaultConfig(c.CacheLifetime)
	cacheConfig.HardMaxCacheSize = c.MaxCacheSizeInMB
	if c.CacheLifetime > 0 {
		cacheConfig.CleanWindow = 1 * time.Minute
	}
	cache, err := a_cache.NewBigCache(cacheConfig)
	if err != nil {
		return nil, err
	}

	return &Bigcache{
		cache: cache,
	}, nil
}

func (c *Bigcache) Set(key string, object interface{}) error {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	if err := enc.Encode(object); err != nil {
		return errors.ThrowInvalidArgument(err, "FASTC-RUyxI", "unable to encode object")
	}
	return c.cache.Set(key, b.Bytes())
}

func (c *Bigcache) Get(key string, ptrToObject interface{}) error {
	value, err := c.cache.Get(key)
	if err == a_cache.ErrEntryNotFound {
		return errors.ThrowNotFound(err, "BIGCA-we32s", "not in cache")
	}
	if err != nil {
		logging.Log("BIGCA-ftofbc").WithError(err).Info("read from cache failed")
		return errors.ThrowInvalidArgument(err, "BIGCA-3idls", "error in reading from cache")
	}

	b := bytes.NewBuffer(value)
	dec := gob.NewDecoder(b)

	return dec.Decode(ptrToObject)
}

func (c *Bigcache) Delete(key string) error {
	return c.cache.Delete(key)
}
