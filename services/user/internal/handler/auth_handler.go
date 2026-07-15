package handler

import (
	"context"

	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Muhammadpurwanto/ecommerce-baju/services/user/internal/dto"
	"github.com/Muhammadpurwanto/ecommerce-baju/services/user/internal/service"
	pb "github.com/Muhammadpurwanto/ecommerce-baju/services/common/pb/auth"
)

type AuthGrpcHandler struct {
	pb.UnimplementedAuthServiceServer
	authService service.AuthService
	validate    *validator.Validate
}

func NewAuthGrpcHandler(authService service.AuthService) *AuthGrpcHandler {
	return &AuthGrpcHandler{
		authService: authService,
		validate:    validator.New(),
	}
}

func (h *AuthGrpcHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.TokenResponse, error) {
	dtoReq := &dto.RegisterRequest{
		Name:     req.GetName(),
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}

	if err := h.validate.Struct(dtoReq); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validation error: %v", err)
	}

	res, err := h.authService.Register(dtoReq)
	if err != nil {
		if err.Error() == "email already exists" {
			return nil, status.Errorf(codes.AlreadyExists, "failed to register: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to register: %v", err)
	}

	return &pb.TokenResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
		ExpiresIn:    res.ExpiresIn,
	}, nil
}

func (h *AuthGrpcHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.TokenResponse, error) {
	dtoReq := &dto.LoginRequest{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}

	if err := h.validate.Struct(dtoReq); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validation error: %v", err)
	}

	res, err := h.authService.Login(dtoReq)
	if err != nil {
		if err.Error() == "invalid email or password" {
			return nil, status.Errorf(codes.Unauthenticated, "login failed: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "login failed: %v", err)
	}

	return &pb.TokenResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
		ExpiresIn:    res.ExpiresIn,
	}, nil
}

func (h *AuthGrpcHandler) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	err := h.authService.Logout(req.GetAccessToken(), req.GetRefreshToken())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "logout failed: %v", err)
	}

	return &pb.LogoutResponse{
		Success: true,
	}, nil
}

func (h *AuthGrpcHandler) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.TokenResponse, error) {
	dtoReq := &dto.RefreshTokenRequest{
		RefreshToken: req.GetRefreshToken(),
	}

	if err := h.validate.Struct(dtoReq); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validation error: %v", err)
	}

	res, err := h.authService.RefreshToken(dtoReq)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "failed to refresh token: %v", err)
	}

	return &pb.TokenResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
		ExpiresIn:    res.ExpiresIn,
	}, nil
}
