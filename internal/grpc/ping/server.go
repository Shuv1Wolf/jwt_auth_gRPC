package ping

import (
	"context"
	"errors"
	"jwt_auth_gRPC/sso/internal/storage"

	ssov1 "github.com/Shuv1Wolf/jwt_protos/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Ping interface {
	Ping(ctx context.Context,
		appID int64,
	) (bool, error)
	NewApp(ctx context.Context,
		id int64,
		name string,
		secret string,
	) (int64, error)
}

type serverAPI struct {
	ssov1.UnimplementedPingServer
	ping Ping
}

func Register(gRPC *grpc.Server, ping Ping) {
	ssov1.RegisterPingServer(gRPC, &serverAPI{ping: ping})
}

// TODO: добавить валидацию
func (s *serverAPI) Ping(ctx context.Context, req *ssov1.IsPingRequest) (*ssov1.IsPingResponse, error) {
	ping, err := s.ping.Ping(ctx, req.GetAppId())
	if err != nil {
		if errors.Is(err, storage.ErrAppNotFound) {
			return &ssov1.IsPingResponse{
				Client: ping,
			}, nil
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &ssov1.IsPingResponse{
		Client: ping,
	}, nil
}

func (s *serverAPI) NewApp(ctx context.Context, req *ssov1.IsNewAppRequest) (*ssov1.IsNewAppResponse, error) {
	saveApp, err := s.ping.NewApp(ctx, req.GetId(), req.GetName(), req.GetSecret())
	if err != nil {
		if errors.Is(err, storage.ErrAppExists) {
			return &ssov1.IsNewAppResponse{
				Id: saveApp,
			}, nil
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.IsNewAppResponse{
		Id: saveApp,
	}, nil
}
