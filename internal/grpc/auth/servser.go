package auth

import (
	"context"
	"errors"
	"jwt_auth_gRPC/sso/internal/services/auth"
	"jwt_auth_gRPC/sso/internal/storage"
	"net/mail"
	"regexp"

	ssov1 "github.com/Shuv1Wolf/jwt_protos/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	Login(ctx context.Context,
		email string,
		password string,
		appID int,
	) (token string, err error)
	RegisterNewUser(ctx context.Context,
		email string,
		password string,
	) (userID int64, err error)
	IsAdmin(ctx context.Context,
		userID int64,
	) (bool, error)
}

type Ping interface {
	Ping(ctx context.Context,
		appID int64,
	) (bool, error)
}

type serverAPI struct {
	ssov1.UnimplementedAuthServer
	auth Auth
	ping Ping
}

func Register(gRPC *grpc.Server, auth Auth, ping Ping) {
	ssov1.RegisterAuthServer(gRPC, &serverAPI{auth: auth, ping: ping})
}

const (
	emptyValue = 0
)

func (s *serverAPI) Login(ctx context.Context, req *ssov1.LoginRequest) (*ssov1.LoginResponse, error) {
	err := validateLogin(req)
	if err != nil {
		return nil, err
	}

	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword(), int(req.GetAppId()))
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid email or password")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.LoginResponse{
		AccessToken: token,
	}, nil
}

func (s *serverAPI) Register(ctx context.Context, req *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error) {
	err := validateRegister(req)
	if err != nil {
		return nil, err
	}

	userID, err := s.auth.RegisterNewUser(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.RegisterResponse{
		UserId: userID,
	}, nil
}

func (s *serverAPI) IsAdmin(ctx context.Context, req *ssov1.IsAdminRequest) (*ssov1.IsAdminResponse, error) {
	err := validateIsAdmin(req)
	if err != nil {
		return nil, err
	}

	isAdmin, err := s.auth.IsAdmin(ctx, req.GetUserId())
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.IsAdminResponse{
		IsAdmin: isAdmin,
	}, nil
}

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

// validation
func validateLogin(req *ssov1.LoginRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email is required")
	}
	_, err := mail.ParseAddress(req.GetEmail())
	if err != nil {
		return status.Error(codes.InvalidArgument, "incorrect email address")
	}

	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}

	if req.GetAppId() == emptyValue {
		return status.Error(codes.InvalidArgument, "app_id is required")
	}
	return nil
}

func validateRegister(req *ssov1.RegisterRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email is required")
	}

	_, err := mail.ParseAddress(req.GetEmail())
	if err != nil {
		return status.Error(codes.InvalidArgument, "incorrect email address")
	}

	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}

	if len(req.GetPassword()) < 8 {
		return status.Error(codes.InvalidArgument, "password must be at least 8 characters long")
	}

	match, _ := regexp.MatchString(`[A-Z]`, req.GetPassword())
	if !match {
		return status.Error(codes.InvalidArgument, "password must contain at least one uppercase letter")
	}

	match, _ = regexp.MatchString(`[a-z]`, req.GetPassword())
	if !match {
		return status.Error(codes.InvalidArgument, "password must contain at least one lowercase letter")
	}

	match, _ = regexp.MatchString(`[0-9]`, req.GetPassword())
	if !match {
		return status.Error(codes.InvalidArgument, "password must contain at least one number")
	}

	return nil
}

func validateIsAdmin(req *ssov1.IsAdminRequest) error {
	if req.GetUserId() == emptyValue {
		return status.Error(codes.InvalidArgument, "user_id is required")
	}
	return nil
}
