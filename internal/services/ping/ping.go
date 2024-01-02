package ping

import (
	"context"
	"errors"
	"fmt"
	"jwt_auth_gRPC/sso/internal/domain/models"
	"jwt_auth_gRPC/sso/internal/lib/logger/sl"
	"jwt_auth_gRPC/sso/internal/storage"
	"log/slog"
)

type Ping struct {
	log         *slog.Logger
	appProvider AppProvider
}

type AppProvider interface {
	App(ctx context.Context, appID int64) (models.App, error)
}

func New(
	log *slog.Logger,
	appProvider AppProvider,
) *Ping {
	return &Ping{
		log:         log,
		appProvider: appProvider,
	}
}

func (p *Ping) Ping(ctx context.Context, appID int64) (bool, error) {
	const op = "auth.Ping"

	log := p.log.With(slog.String("op", op))
	log.Info("Ping app")

	_, err := p.appProvider.App(ctx, appID)
	if err != nil {
		if errors.Is(err, storage.ErrAppNotFound) {
			log.Warn("app not found", sl.Err(err))

			return false, nil
		}

		return false, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("app found", slog.Bool("client", true))

	return true, nil
}
