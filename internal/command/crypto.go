package command

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
)

type cryptoCode struct {
	value  *crypto.CryptoValue
	plain  string
	expiry time.Duration
}

func newCryptoCodeWithExpiry(ctx context.Context, filter preparation.FilterToQueryReducer, typ domain.SecretGeneratorType, alg crypto.Crypto) (*cryptoCode, error) {
	config, err := secretGeneratorConfig(ctx, filter, typ)
	if err != nil {
		return nil, err
	}
	code := &cryptoCode{
		expiry: config.Expiry,
	}

	switch a := alg.(type) {
	case crypto.HashAlgorithm:
		code.value, code.plain, err = crypto.NewCode(crypto.NewHashGenerator(*config, a))
	case crypto.EncryptionAlgorithm:
		code.value, code.plain, err = crypto.NewCode(crypto.NewEncryptionGenerator(*config, a))
	default:
		return nil, errors.ThrowInternal(nil, "COMMA-RreV6", "Errors.Internal")
	}
	if err != nil {
		return nil, err
	}
	return code, nil
}

func newCryptoCodeWithPlain(ctx context.Context, filter preparation.FilterToQueryReducer, typ domain.SecretGeneratorType, alg crypto.Crypto) (value *crypto.CryptoValue, plain string, err error) {
	config, err := secretGeneratorConfig(ctx, filter, typ)
	if err != nil {
		return nil, "", err
	}

	switch a := alg.(type) {
	case crypto.HashAlgorithm:
		return crypto.NewCode(crypto.NewHashGenerator(*config, a))
	case crypto.EncryptionAlgorithm:
		return crypto.NewCode(crypto.NewEncryptionGenerator(*config, a))
	}

	return nil, "", errors.ThrowInvalidArgument(nil, "V2-NGESt", "Errors.Internal")
}

func secretGeneratorConfig(ctx context.Context, filter preparation.FilterToQueryReducer, typ domain.SecretGeneratorType) (*crypto.GeneratorConfig, error) {
	wm := NewInstanceSecretGeneratorConfigWriteModel(ctx, typ)
	events, err := filter(ctx, wm.Query())
	if err != nil {
		return nil, err
	}
	wm.AppendEvents(events...)
	if err := wm.Reduce(); err != nil {
		return nil, err
	}
	return &crypto.GeneratorConfig{
		Length:              wm.Length,
		Expiry:              wm.Expiry,
		IncludeLowerLetters: wm.IncludeLowerLetters,
		IncludeUpperLetters: wm.IncludeUpperLetters,
		IncludeDigits:       wm.IncludeDigits,
		IncludeSymbols:      wm.IncludeSymbols,
	}, nil
}
