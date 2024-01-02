package app

import (
	grpcapp "jwt_auth_gRPC/sso/internal/app/grpc"
	"jwt_auth_gRPC/sso/internal/services/auth"
	"jwt_auth_gRPC/sso/internal/services/ping"
	"jwt_auth_gRPC/sso/internal/storage/sqlite"
	"log/slog"
	"time"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(log *slog.Logger, grpcPort int, storagePath string, tokenTTL time.Duration) *App {
	storage, err := sqlite.New(storagePath)
	if err != nil {
		panic(err)
	}

	pingService := ping.New(log, storage)
	authService := auth.New(log, storage, storage, storage, tokenTTL)

	grpcApp := grpcapp.New(log, authService, pingService, grpcPort)

	return &App{
		GRPCServer: grpcApp,
	}
}
