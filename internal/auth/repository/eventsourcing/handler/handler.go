package handler

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/auth/repository/eventsourcing/view"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
	handler2 "github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	query2 "github.com/zitadel/zitadel/internal/query"
)

type Config struct {
	Client     *database.DB
	Eventstore *eventstore.Eventstore

	BulkLimit             uint64
	FailureCountUntilSkip uint64
	HandleActiveInstances time.Duration
	TransactionDuration   time.Duration
	Handlers              map[string]*ConfigOverwrites
}

type ConfigOverwrites struct {
	MinimumCycleDuration time.Duration
}

func Register(ctx context.Context, configs Config, view *view.View, queries *query2.Queries) {
	newUser(ctx,
		configs.overwrite("User"),
		view,
		queries,
	).Start(ctx)

	newUserSession(ctx,
		configs.overwrite("UserSession"),
		view,
		queries,
	).Start(ctx)

	newToken(ctx,
		configs.overwrite("Token"),
		view,
	).Start(ctx)

	newRefreshToken(ctx,
		configs.overwrite("RefreshToken"),
		view,
	).Start(ctx)
}

func (config Config) overwrite(viewModel string) handler2.Config {
	c := handler2.Config{
		Client:                config.Client,
		Eventstore:            config.Eventstore,
		BulkLimit:             uint16(config.BulkLimit),
		RequeueEvery:          3 * time.Minute,
		HandleActiveInstances: config.HandleActiveInstances,
		MaxFailureCount:       uint8(config.FailureCountUntilSkip),
		TransactionDuration:   config.TransactionDuration,
	}
	overwrite, ok := config.Handlers[viewModel]
	if !ok {
		return c
	}
	if overwrite.MinimumCycleDuration > 0 {
		c.RequeueEvery = overwrite.MinimumCycleDuration
	}
	return c
}
