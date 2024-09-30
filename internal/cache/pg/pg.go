package pg

import (
	"context"
	_ "embed"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/cache"
)

var (
	//go:embed set.sql
	setQuery string
	//go:embed get.sql
	getQuery string
	//go:embed invalidate.sql
	invalidateQuery string
	//go:embed delete.sql
	deleteQuery string
	//go:embed prune.sql
	pruneQuery string
	//go:embed truncate.sql
	truncateQuery string
)

type PGXPool interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type pgCache[I ~int, K ~string, V cache.Entry[I, K]] struct {
	name    string
	config  *cache.CacheConfig
	indices []I
	pool    PGXPool
	logger  *slog.Logger
}

// NewCache returns a cache that stores and retrieves objects using PostgreSQL unlogged tables.
func NewCache[I ~int, K ~string, V cache.Entry[I, K]](name string, config cache.CacheConfig, indices []I, pool PGXPool) cache.PrunerCache[I, K, V] {
	return &pgCache[I, K, V]{
		name:    name,
		config:  &config,
		indices: indices,
		pool:    pool,
		logger:  config.Log.Slog().With("cache_name", name),
	}
}

func (c *pgCache[I, K, V]) Set(ctx context.Context, entry V) {
	//nolint:errcheck
	c.set(ctx, entry)
}

func (c *pgCache[I, K, V]) set(ctx context.Context, entry V) error {
	_, err := c.pool.Exec(ctx, setQuery, c.name, c.indexSetFromEntry(entry), entry)
	if err != nil {
		c.logger.ErrorContext(ctx, "pg cache set", "err", err)
		return err
	}
	return nil
}

func (c *pgCache[I, K, V]) Get(ctx context.Context, index I, key K) (value V, ok bool) {
	err := c.pool.QueryRow(ctx, getQuery, c.name, index, key, c.config.MaxAge, c.config.LastUseAge).Scan(&value)
	if err == nil {
		return value, true
	}
	logger := c.logger.With("err", err, "index", index, "key", key)
	if errors.Is(err, pgx.ErrNoRows) {
		logger.InfoContext(ctx, "pg cache miss")
		return value, false
	}
	logger.ErrorContext(ctx, "pg cache get")
	return value, false
}

func (c *pgCache[I, K, V]) Invalidate(ctx context.Context, index I, keys ...K) (err error) {
	_, err = c.pool.Exec(ctx, invalidateQuery, c.name, index, keys)
	logging.OnError(err).Error("pg cache invalidate")
	return err
}

func (c *pgCache[I, K, V]) Delete(ctx context.Context, index I, keys ...K) (err error) {
	_, err = c.pool.Exec(ctx, deleteQuery, c.name, index, keys)
	logging.OnError(err).Error("pg cache delete")
	return err
}

func (c *pgCache[I, K, V]) Prune(ctx context.Context) (err error) {
	_, err = c.pool.Exec(ctx, pruneQuery, c.name, c.config.MaxAge, c.config.LastUseAge)
	logging.OnError(err).Error("pg cache delete")
	return err
}

func (c *pgCache[I, K, V]) Truncate(ctx context.Context) (err error) {
	_, err = c.pool.Exec(ctx, truncateQuery, c.name)
	logging.OnError(err).Error("pg cache delete")
	return err
}

func (c *pgCache[I, K, V]) Close(context.Context) (err error) { return }

type indexKey[I, K comparable] struct {
	IndexID  I `json:"index_id"`
	IndexKey K `json:"index_key"`
}

func (c *pgCache[I, K, V]) indexSetFromEntry(entry V) []indexKey[I, K] {
	keys := make([]indexKey[I, K], 0, len(c.indices)*3) // naive assumption
	for _, index := range c.indices {
		for _, key := range entry.Keys(index) {
			keys = append(keys, indexKey[I, K]{
				IndexID:  index,
				IndexKey: key,
			})
		}
	}
	return keys
}
