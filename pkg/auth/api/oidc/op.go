package oidc

import (
	"context"
	"time"

	"github.com/caos/logging"
	"github.com/caos/oidc/pkg/op"

	http_utils "github.com/caos/zitadel/internal/api/http"
	"github.com/caos/zitadel/internal/auth/repository"
	"github.com/caos/zitadel/internal/id"
)

type OPHandlerConfig struct {
	OPConfig              *op.Config
	StorageConfig         StorageConfig
	UserAgentCookieConfig *http_utils.UserAgentCookieConfig
	Endpoints             *EndpointConfig
}

type StorageConfig struct {
	DefaultLoginURL string
	//TokenLifetime         string
}

type EndpointConfig struct {
	Auth       *Endpoint
	Token      *Endpoint
	Userinfo   *Endpoint
	EndSession *Endpoint
	Keys       *Endpoint
}

type Endpoint struct {
	Path string
	URL  string
}

type OPStorage struct {
	repo repository.Repository
	//config          *op.Config
	defaultLoginURL string
	tokenLifetime   time.Duration
}

func NewProvider(ctx context.Context, config OPHandlerConfig, repo repository.Repository) op.OpenIDProvider {
	cookieHandler, err := http_utils.NewUserAgentHandler(config.UserAgentCookieConfig, id.SonyFlakeGenerator)
	logging.Log("OIDC-sd4fd").OnError(err).Panic("cannot user agent handler")
	provider, err := op.NewDefaultOP(
		ctx,
		config.OPConfig,
		newStorage(config.StorageConfig, repo),
		op.WithHttpInterceptor(
			UserAgentCookieHandler(
				cookieHandler,
				http_utils.CopyHeadersToContext,
			),
		),
		op.WithCustomAuthEndpoint(op.NewEndpointWithURL(config.Endpoints.Auth.Path, config.Endpoints.Auth.URL)),
		op.WithCustomTokenEndpoint(op.NewEndpointWithURL(config.Endpoints.Token.Path, config.Endpoints.Token.URL)),
		op.WithCustomUserinfoEndpoint(op.NewEndpointWithURL(config.Endpoints.Userinfo.Path, config.Endpoints.Userinfo.URL)),
		op.WithCustomEndSessionEndpoint(op.NewEndpointWithURL(config.Endpoints.EndSession.Path, config.Endpoints.EndSession.URL)),
		op.WithCustomKeysEndpoint(op.NewEndpointWithURL(config.Endpoints.Keys.Path, config.Endpoints.Keys.URL)),
	)
	logging.Log("OIDC-asf13").OnError(err).Panic("cannot create provider")
	return provider
}

func newStorage(config StorageConfig, repo repository.Repository) *OPStorage {
	return &OPStorage{
		repo: repo,
		//config:          config.OPConfig,
		defaultLoginURL: config.DefaultLoginURL,
		//op.tokenLifetime, _ = time.ParseDuration(c.TokenLifetime)
	}
}

func (o *OPStorage) Health(ctx context.Context) error {
	return o.repo.Health(ctx)
}
