package handler

import (
	"context"
	"time"

	"github.com/go-playground/validator/v10"

	"github.com/Muhammadpurwanto/ecommerce-baju/services/user/internal/dto"
	"github.com/Muhammadpurwanto/ecommerce-baju/services/user/internal/service"
	pb "github.com/Muhammadpurwanto/ecommerce-baju/services/common/pb/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserGrpcHandler struct {
	pb.UnimplementedUserServiceServer
	userService    service.UserService
	addressService service.AddressService
	validate       *validator.Validate
}

func NewUserGrpcHandler(userService service.UserService, addressService service.AddressService) *UserGrpcHandler {
	return &UserGrpcHandler{
		userService:    userService,
		addressService: addressService,
		validate:       validator.New(),
	}
}

// Fungsi helper kecil untuk mengubah format
func toPbUserResponse(res *dto.UserResponse) *pb.UserResponse {
	pbRes := &pb.UserResponse{
		Id:        res.ID,
		Email:     res.Email,
		Name:      res.Name,
		Role:      res.Role,
		IsActive:  res.IsActive,
		CreatedAt: res.CreatedAt.Format(time.RFC3339),
		UpdatedAt: res.UpdatedAt.Format(time.RFC3339),
	}
	if res.Phone != nil {
		pbRes.Phone = *res.Phone
	}
	if res.AvatarURL != nil {
		pbRes.AvatarUrl = *res.AvatarURL
	}
	return pbRes
}

func (h *UserGrpcHandler) GetProfile(ctx context.Context, req *pb.GetProfileRequest) (*pb.UserResponse, error) {
	if req.GetUserId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "user_id is required")
	}

	res, err := h.userService.GetProfile(req.GetUserId())
	if err != nil {
		if err.Error() == "user not found" {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to fetch profile: %v", err)
	}
	return toPbUserResponse(res), nil
}

func (h *UserGrpcHandler) UpdateProfile(ctx context.Context, req *pb.UpdateProfileRequest) (*pb.UserResponse, error) {
	if req.GetUserId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "user_id is required")
	}

	dtoReq := &dto.UpdateUserRequest{}
	
	if req.Name != nil {
		dtoReq.Name = req.GetName()
	}
	if req.Phone != nil {
		p := req.GetPhone()
		dtoReq.Phone = &p
	}
	if req.AvatarUrl != nil {
		a := req.GetAvatarUrl()
		dtoReq.AvatarURL = &a
	}

	// Validasi input
	if err := h.validate.Struct(dtoReq); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validation error: %v", err)
	}

	res, err := h.userService.UpdateProfile(req.GetUserId(), dtoReq)
	if err != nil {
		if err.Error() == "user not found" {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to update profile: %v", err)
	}

	return toPbUserResponse(res), nil
}

func (h *UserGrpcHandler) GetAllUsers(ctx context.Context, req *pb.GetAllUsersRequest) (*pb.GetAllUsersResponse, error) {
	users, meta, err := h.userService.GetAllUsers(int(req.GetPage()), int(req.GetPerPage()))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch users: %v", err)
	}

	var pbUsers []*pb.UserResponse
	for _, u := range users {
		// Buat salinan agar pointer aman
		val := u
		pbUsers = append(pbUsers, toPbUserResponse(&val))
	}

	return &pb.GetAllUsersResponse{
		Data: pbUsers,
		Meta: &pb.MetaData{
			Page:       int32(meta.Page),
			PerPage:    int32(meta.PerPage),
			Total:      meta.Total,
			TotalPages: int32(meta.TotalPages),
		},
	}, nil
}

func (h *UserGrpcHandler) GetUserByID(ctx context.Context, req *pb.GetUserByIDRequest) (*pb.UserResponse, error) {
	if req.GetId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "id is required")
	}

	res, err := h.userService.GetUserByID(req.GetId())
	if err != nil {
		if err.Error() == "user not found" {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to fetch user: %v", err)
	}
	return toPbUserResponse(res), nil
}

// ==========================================
// ADDRESS gRPC HANDLERS
// ==========================================

func toPbAddressResponse(res *dto.AddressResponse) *pb.AddressResponse {
	return &pb.AddressResponse{
		Id:         uint64(res.ID),
		UserId:     res.UserID,
		Label:      res.Label,
		Recipient:  res.Recipient,
		Phone:      res.Phone,
		Province:   res.Province,
		City:       res.City,
		District:   res.District,
		PostalCode: res.PostalCode,
		Detail:     res.Detail,
		IsDefault:  res.IsDefault,
		CreatedAt:  res.CreatedAt.Format(time.RFC3339),
	}
}

func (h *UserGrpcHandler) GetAddresses(ctx context.Context, req *pb.GetAddressesRequest) (*pb.GetAddressesResponse, error) {
	if req.GetUserId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "user_id is required")
	}

	addresses, err := h.addressService.GetAddresses(req.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get addresses: %v", err)
	}

	var pbAddresses []*pb.AddressResponse
	for _, a := range addresses {
		val := a
		pbAddresses = append(pbAddresses, toPbAddressResponse(&val))
	}

	return &pb.GetAddressesResponse{Data: pbAddresses}, nil
}

func (h *UserGrpcHandler) GetAddressByID(ctx context.Context, req *pb.GetAddressByIDRequest) (*pb.AddressResponse, error) {
	if req.GetUserId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "user_id is required")
	}

	address, err := h.addressService.GetAddressByID(uint(req.GetId()), req.GetUserId())
	if err != nil {
		if err.Error() == "address not found" {
			return nil, status.Errorf(codes.NotFound, "address not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get address: %v", err)
	}

	return toPbAddressResponse(address), nil
}

func (h *UserGrpcHandler) CreateAddress(ctx context.Context, req *pb.CreateAddressRequest) (*pb.AddressResponse, error) {
	if req.GetUserId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "user_id is required")
	}

	dtoReq := &dto.CreateAddressRequest{
		Label:      req.GetLabel(),
		Recipient:  req.GetRecipient(),
		Phone:      req.GetPhone(),
		Province:   req.GetProvince(),
		City:       req.GetCity(),
		District:   req.GetDistrict(),
		PostalCode: req.GetPostalCode(),
		Detail:     req.GetDetail(),
		IsDefault:  req.GetIsDefault(),
	}

	// Validasi input
	if err := h.validate.Struct(dtoReq); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validation error: %v", err)
	}

	address, err := h.addressService.CreateAddress(req.GetUserId(), dtoReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create address: %v", err)
	}

	return toPbAddressResponse(address), nil
}

func (h *UserGrpcHandler) UpdateAddress(ctx context.Context, req *pb.UpdateAddressRequest) (*pb.AddressResponse, error) {
	if req.GetUserId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "user_id is required")
	}

	dtoReq := &dto.UpdateAddressRequest{
		Label:      req.GetLabel(),
		Recipient:  req.GetRecipient(),
		Phone:      req.GetPhone(),
		Province:   req.GetProvince(),
		City:       req.GetCity(),
		District:   req.GetDistrict(),
		PostalCode: req.GetPostalCode(),
		Detail:     req.GetDetail(),
	}
	if req.IsDefault != nil {
		isDefault := req.GetIsDefault()
		dtoReq.IsDefault = &isDefault
	}

	// Validasi input
	if err := h.validate.Struct(dtoReq); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validation error: %v", err)
	}

	address, err := h.addressService.UpdateAddress(uint(req.GetId()), req.GetUserId(), dtoReq)
	if err != nil {
		if err.Error() == "address not found" {
			return nil, status.Errorf(codes.NotFound, "address not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to update address: %v", err)
	}

	return toPbAddressResponse(address), nil
}

func (h *UserGrpcHandler) DeleteAddress(ctx context.Context, req *pb.DeleteAddressRequest) (*pb.DeleteAddressResponse, error) {
	if req.GetUserId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "user_id is required")
	}

	err := h.addressService.DeleteAddress(uint(req.GetId()), req.GetUserId())
	if err != nil {
		if err.Error() == "address not found" {
			return nil, status.Errorf(codes.NotFound, "address not found")
		}
		return &pb.DeleteAddressResponse{Success: false, Message: err.Error()}, nil
	}

	return &pb.DeleteAddressResponse{Success: true, Message: "Address deleted successfully"}, nil
}

func (h *UserGrpcHandler) SetDefaultAddress(ctx context.Context, req *pb.SetDefaultAddressRequest) (*pb.AddressResponse, error) {
	if req.GetUserId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "user_id is required")
	}

	err := h.addressService.SetDefaultAddress(uint(req.GetId()), req.GetUserId())
	if err != nil {
		if err.Error() == "address not found" {
			return nil, status.Errorf(codes.NotFound, "address not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to set default address: %v", err)
	}

	// Ambil data address setelah diperbarui
	address, err := h.addressService.GetAddressByID(uint(req.GetId()), req.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch updated address: %v", err)
	}

	return toPbAddressResponse(address), nil
}
