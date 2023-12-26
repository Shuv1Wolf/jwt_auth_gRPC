package auth

import (
	"context"
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

type serverAPI struct {
	ssov1.UnimplementedAuthServer
}

func Register(gRPC *grpc.Server) {
	ssov1.RegisterAuthServer(gRPC, &serverAPI{})
}

const (
	emptyValue = 0
)

func (s *serverAPI) Login(ctx context.Context, req *ssov1.LoginRequest) (*ssov1.LoginResponse, error) {
	err := validateLogin(req)
	if err != nil {
		return nil, err
	}
	return &ssov1.LoginResponse{
		AccessToken: "token",
	}, nil
}

func (s *serverAPI) Register(ctx context.Context, req *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error) {
	err := validateRegister(req)
	if err != nil {
		return nil, err
	}
	return &ssov1.RegisterResponse{
		UserId: 3,
	}, nil
}

func (s *serverAPI) IsAdmin(ctx context.Context, req *ssov1.IsAdminRequest) (*ssov1.IsAdminResponse, error) {
	err := validateIsAdmin(req)
	if err != nil {
		return nil, err
	}
	return &ssov1.IsAdminResponse{
		IsAdmin: true,
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
