package pg

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/cache"
)

type testIndex int

const (
	testIndexID testIndex = iota
	testIndexName
)

var testIndices = []testIndex{
	testIndexID,
	testIndexName,
}

type testObject struct {
	ID   string
	Name []string
}

func (o *testObject) Keys(index testIndex) []string {
	switch index {
	case testIndexID:
		return []string{o.ID}
	case testIndexName:
		return o.Name
	default:
		return nil
	}
}

func TestNewCache(t *testing.T) {
	tests := []struct {
		name    string
		expect  func(pgxmock.PgxCommonIface)
		wantErr error
	}{
		{
			name: "error",
			expect: func(pci pgxmock.PgxCommonIface) {
				pci.ExpectExec(regexp.QuoteMeta(expectedCreatePartitionQuery)).
					WillReturnError(pgx.ErrTxClosed)
			},
			wantErr: pgx.ErrTxClosed,
		},
		{
			name: "success",
			expect: func(pci pgxmock.PgxCommonIface) {
				pci.ExpectExec(regexp.QuoteMeta(expectedCreatePartitionQuery)).
					WillReturnResult(pgxmock.NewResult("CREATE TABLE", 0))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conf := cache.CacheConfig{
				Log: &logging.Config{
					Level:     "debug",
					AddSource: true,
				},
			}
			pool, err := pgxmock.NewPool()
			require.NoError(t, err)
			tt.expect(pool)

			c, err := NewCache[testIndex, string, *testObject](context.Background(), cacheName, conf, testIndices, pool)
			require.ErrorIs(t, err, tt.wantErr)
			if tt.wantErr == nil {
				assert.NotNil(t, c)
			}

			err = pool.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}

}

func Test_pgCache_Set(t *testing.T) {
	queryExpect := regexp.QuoteMeta(setQuery)
	type args struct {
		entry *testObject
	}
	tests := []struct {
		name    string
		args    args
		expect  func(pgxmock.PgxCommonIface)
		wantErr error
	}{
		{
			name: "error",
			args: args{
				&testObject{
					ID:   "id1",
					Name: []string{"foo", "bar"},
				},
			},
			expect: func(ppi pgxmock.PgxCommonIface) {
				ppi.ExpectExec(queryExpect).
					WithArgs("test",
						[]indexKey[testIndex, string]{
							{IndexID: testIndexID, IndexKey: "id1"},
							{IndexID: testIndexName, IndexKey: "foo"},
							{IndexID: testIndexName, IndexKey: "bar"},
						},
						&testObject{
							ID:   "id1",
							Name: []string{"foo", "bar"},
						}).
					WillReturnError(pgx.ErrTxClosed)
			},
			wantErr: pgx.ErrTxClosed,
		},
		{
			name: "success",
			args: args{
				&testObject{
					ID:   "id1",
					Name: []string{"foo", "bar"},
				},
			},
			expect: func(ppi pgxmock.PgxCommonIface) {
				ppi.ExpectExec(queryExpect).
					WithArgs("test",
						[]indexKey[testIndex, string]{
							{IndexID: testIndexID, IndexKey: "id1"},
							{IndexID: testIndexName, IndexKey: "foo"},
							{IndexID: testIndexName, IndexKey: "bar"},
						},
						&testObject{
							ID:   "id1",
							Name: []string{"foo", "bar"},
						}).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, pool := prepareCache(t, cache.CacheConfig{})
			defer pool.Close()
			tt.expect(pool)

			err := c.(*pgCache[testIndex, string, *testObject]).
				set(context.Background(), tt.args.entry)
			require.ErrorIs(t, err, tt.wantErr)

			err = pool.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}

func Test_pgCache_Get(t *testing.T) {
	queryExpect := regexp.QuoteMeta(getQuery)
	type args struct {
		index testIndex
		key   string
	}
	tests := []struct {
		name   string
		config cache.CacheConfig
		args   args
		expect func(pgxmock.PgxCommonIface)
		want   *testObject
		wantOk bool
	}{
		{
			name: "invalid index",
			config: cache.CacheConfig{
				MaxAge:     time.Minute,
				LastUseAge: time.Second,
			},
			args: args{
				index: 99,
				key:   "id1",
			},
			expect: func(pci pgxmock.PgxCommonIface) {},
			wantOk: false,
		},
		{
			name: "no rows",
			config: cache.CacheConfig{
				MaxAge:     0,
				LastUseAge: 0,
			},
			args: args{
				index: testIndexID,
				key:   "id1",
			},
			expect: func(pci pgxmock.PgxCommonIface) {
				pci.ExpectQuery(queryExpect).
					WithArgs("test", testIndexID, "id1", time.Duration(0), time.Duration(0)).
					WillReturnRows(pgxmock.NewRows([]string{"payload"}))
			},
			wantOk: false,
		},
		{
			name: "error",
			config: cache.CacheConfig{
				MaxAge:     0,
				LastUseAge: 0,
			},
			args: args{
				index: testIndexID,
				key:   "id1",
			},
			expect: func(pci pgxmock.PgxCommonIface) {
				pci.ExpectQuery(queryExpect).
					WithArgs("test", testIndexID, "id1", time.Duration(0), time.Duration(0)).
					WillReturnError(pgx.ErrTxClosed)
			},
			wantOk: false,
		},
		{
			name: "ok",
			config: cache.CacheConfig{
				MaxAge:     time.Minute,
				LastUseAge: time.Second,
			},
			args: args{
				index: testIndexID,
				key:   "id1",
			},
			expect: func(pci pgxmock.PgxCommonIface) {
				pci.ExpectQuery(queryExpect).
					WithArgs("test", testIndexID, "id1", time.Minute, time.Second).
					WillReturnRows(
						pgxmock.NewRows([]string{"payload"}).AddRow(&testObject{
							ID:   "id1",
							Name: []string{"foo", "bar"},
						}),
					)
			},
			want: &testObject{
				ID:   "id1",
				Name: []string{"foo", "bar"},
			},
			wantOk: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, pool := prepareCache(t, tt.config)
			defer pool.Close()
			tt.expect(pool)

			got, ok := c.Get(context.Background(), tt.args.index, tt.args.key)
			assert.Equal(t, tt.wantOk, ok)
			assert.Equal(t, tt.want, got)
			err := pool.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}

func Test_pgCache_Invalidate(t *testing.T) {
	queryExpect := regexp.QuoteMeta(invalidateQuery)
	type args struct {
		index testIndex
		keys  []string
	}
	tests := []struct {
		name    string
		config  cache.CacheConfig
		args    args
		expect  func(pgxmock.PgxCommonIface)
		wantErr error
	}{
		{
			name: "error",
			config: cache.CacheConfig{
				MaxAge:     0,
				LastUseAge: 0,
			},
			args: args{
				index: testIndexID,
				keys:  []string{"id1", "id2"},
			},
			expect: func(pci pgxmock.PgxCommonIface) {
				pci.ExpectExec(queryExpect).
					WithArgs("test", testIndexID, []string{"id1", "id2"}).
					WillReturnError(pgx.ErrTxClosed)
			},
			wantErr: pgx.ErrTxClosed,
		},
		{
			name: "ok",
			config: cache.CacheConfig{
				MaxAge:     time.Minute,
				LastUseAge: time.Second,
			},
			args: args{
				index: testIndexID,
				keys:  []string{"id1", "id2"},
			},
			expect: func(pci pgxmock.PgxCommonIface) {
				pci.ExpectExec(queryExpect).
					WithArgs("test", testIndexID, []string{"id1", "id2"}).
					WillReturnResult(pgxmock.NewResult("DELETE", 1))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, pool := prepareCache(t, tt.config)
			defer pool.Close()
			tt.expect(pool)

			err := c.Invalidate(context.Background(), tt.args.index, tt.args.keys...)
			assert.ErrorIs(t, err, tt.wantErr)

			err = pool.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}

func Test_pgCache_Delete(t *testing.T) {
	queryExpect := regexp.QuoteMeta(deleteQuery)
	type args struct {
		index testIndex
		keys  []string
	}
	tests := []struct {
		name    string
		config  cache.CacheConfig
		args    args
		expect  func(pgxmock.PgxCommonIface)
		wantErr error
	}{
		{
			name: "error",
			config: cache.CacheConfig{
				MaxAge:     0,
				LastUseAge: 0,
			},
			args: args{
				index: testIndexID,
				keys:  []string{"id1", "id2"},
			},
			expect: func(pci pgxmock.PgxCommonIface) {
				pci.ExpectExec(queryExpect).
					WithArgs("test", testIndexID, []string{"id1", "id2"}).
					WillReturnError(pgx.ErrTxClosed)
			},
			wantErr: pgx.ErrTxClosed,
		},
		{
			name: "ok",
			config: cache.CacheConfig{
				MaxAge:     time.Minute,
				LastUseAge: time.Second,
			},
			args: args{
				index: testIndexID,
				keys:  []string{"id1", "id2"},
			},
			expect: func(pci pgxmock.PgxCommonIface) {
				pci.ExpectExec(queryExpect).
					WithArgs("test", testIndexID, []string{"id1", "id2"}).
					WillReturnResult(pgxmock.NewResult("DELETE", 1))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, pool := prepareCache(t, tt.config)
			defer pool.Close()
			tt.expect(pool)

			err := c.Delete(context.Background(), tt.args.index, tt.args.keys...)
			assert.ErrorIs(t, err, tt.wantErr)

			err = pool.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}

func Test_pgCache_Prune(t *testing.T) {
	queryExpect := regexp.QuoteMeta(pruneQuery)
	tests := []struct {
		name    string
		config  cache.CacheConfig
		expect  func(pgxmock.PgxCommonIface)
		wantErr error
	}{
		{
			name: "error",
			config: cache.CacheConfig{
				MaxAge:     0,
				LastUseAge: 0,
			},
			expect: func(pci pgxmock.PgxCommonIface) {
				pci.ExpectExec(queryExpect).
					WithArgs("test", time.Duration(0), time.Duration(0)).
					WillReturnError(pgx.ErrTxClosed)
			},
			wantErr: pgx.ErrTxClosed,
		},
		{
			name: "ok",
			config: cache.CacheConfig{
				MaxAge:     time.Minute,
				LastUseAge: time.Second,
			},
			expect: func(pci pgxmock.PgxCommonIface) {
				pci.ExpectExec(queryExpect).
					WithArgs("test", time.Minute, time.Second).
					WillReturnResult(pgxmock.NewResult("DELETE", 1))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, pool := prepareCache(t, tt.config)
			defer pool.Close()
			tt.expect(pool)

			err := c.Prune(context.Background())
			assert.ErrorIs(t, err, tt.wantErr)

			err = pool.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}

func Test_pgCache_Truncate(t *testing.T) {
	queryExpect := regexp.QuoteMeta(truncateQuery)
	tests := []struct {
		name    string
		config  cache.CacheConfig
		expect  func(pgxmock.PgxCommonIface)
		wantErr error
	}{
		{
			name: "error",
			config: cache.CacheConfig{
				MaxAge:     0,
				LastUseAge: 0,
			},
			expect: func(pci pgxmock.PgxCommonIface) {
				pci.ExpectExec(queryExpect).
					WithArgs("test").
					WillReturnError(pgx.ErrTxClosed)
			},
			wantErr: pgx.ErrTxClosed,
		},
		{
			name: "ok",
			config: cache.CacheConfig{
				MaxAge:     time.Minute,
				LastUseAge: time.Second,
			},
			expect: func(pci pgxmock.PgxCommonIface) {
				pci.ExpectExec(queryExpect).
					WithArgs("test").
					WillReturnResult(pgxmock.NewResult("DELETE", 1))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, pool := prepareCache(t, tt.config)
			defer pool.Close()
			tt.expect(pool)

			err := c.Truncate(context.Background())
			assert.ErrorIs(t, err, tt.wantErr)

			err = pool.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}

const (
	cacheName                    = "test"
	expectedCreatePartitionQuery = `create unlogged table if not exists cache.objects_test
partition of cache.objects
for values in ('test');

create unlogged table if not exists cache.string_keys_test
partition of cache.string_keys
for values in ('test');
`
)

func prepareCache(t *testing.T, conf cache.CacheConfig) (cache.PrunerCache[testIndex, string, *testObject], pgxmock.PgxPoolIface) {
	conf.Log = &logging.Config{
		Level:     "debug",
		AddSource: true,
	}
	pool, err := pgxmock.NewPool()
	require.NoError(t, err)

	pool.ExpectExec(regexp.QuoteMeta(expectedCreatePartitionQuery)).
		WillReturnResult(pgxmock.NewResult("CREATE TABLE", 0))

	c, err := NewCache[testIndex, string, *testObject](context.Background(), cacheName, conf, testIndices, pool)
	require.NoError(t, err)
	return c, pool
}
