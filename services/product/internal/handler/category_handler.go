package handler

import (
	"context"
	"time"

	"github.com/go-playground/validator/v10"

	"github.com/Muhammadpurwanto/ecommerce-baju/services/product/internal/dto"
	"github.com/Muhammadpurwanto/ecommerce-baju/services/product/internal/model"
	"github.com/Muhammadpurwanto/ecommerce-baju/services/product/internal/service"
	pb "github.com/Muhammadpurwanto/ecommerce-baju/services/common/pb/product"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CategoryGrpcHandler struct {
	pb.UnimplementedCategoryServiceServer
	categoryService service.CategoryService
	validate        *validator.Validate
}

func NewCategoryGrpcHandler(categoryService service.CategoryService) *CategoryGrpcHandler {
	return &CategoryGrpcHandler{
		categoryService: categoryService,
		validate:        validator.New(),
	}
}

func toPbCategoryResponse(c *model.Category) *pb.CategoryResponse {
	return &pb.CategoryResponse{
		Id:        uint32(c.ID),
		Name:      c.Name,
		Slug:      c.Slug,
		CreatedAt: c.CreatedAt.Format(time.RFC3339),
		UpdatedAt: c.UpdatedAt.Format(time.RFC3339),
	}
}

func (h *CategoryGrpcHandler) GetAllCategories(ctx context.Context, req *pb.Empty) (*pb.CategoryListResponse, error) {
	categories, err := h.categoryService.GetAll()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get categories: %v", err)
	}

	var pbData []*pb.CategoryResponse
	for _, c := range categories {
		val := c
		pbData = append(pbData, toPbCategoryResponse(&val))
	}

	return &pb.CategoryListResponse{Data: pbData}, nil
}

func (h *CategoryGrpcHandler) GetCategoryByID(ctx context.Context, req *pb.GetByIDRequest) (*pb.CategoryResponse, error) {
	if req.GetId() == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "id is required")
	}

	category, err := h.categoryService.GetByID(uint(req.GetId()))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "category not found: %v", err)
	}
	return toPbCategoryResponse(category), nil
}

func (h *CategoryGrpcHandler) CreateCategory(ctx context.Context, req *pb.CategoryRequest) (*pb.CategoryResponse, error) {
	dtoReq := &dto.CategoryRequest{
		Name: req.GetName(),
	}

	if err := h.validate.Struct(dtoReq); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validation error: %v", err)
	}

	category, err := h.categoryService.Create(dtoReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create category: %v", err)
	}
	return toPbCategoryResponse(category), nil
}

func (h *CategoryGrpcHandler) UpdateCategory(ctx context.Context, req *pb.UpdateCategoryRequest) (*pb.CategoryResponse, error) {
	if req.GetId() == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "id is required")
	}

	dtoReq := &dto.CategoryRequest{
		Name: req.GetName(),
	}

	if err := h.validate.Struct(dtoReq); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validation error: %v", err)
	}

	category, err := h.categoryService.Update(uint(req.GetId()), dtoReq)
	if err != nil {
		if err.Error() == "record not found" || err.Error() == "category not found" {
			return nil, status.Errorf(codes.NotFound, "category not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to update category: %v", err)
	}
	return toPbCategoryResponse(category), nil
}

func (h *CategoryGrpcHandler) DeleteCategory(ctx context.Context, req *pb.GetByIDRequest) (*pb.EmptyResponse, error) {
	if req.GetId() == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "id is required")
	}

	if err := h.categoryService.Delete(uint(req.GetId())); err != nil {
		if err.Error() == "record not found" || err.Error() == "category not found" {
			return nil, status.Errorf(codes.NotFound, "category not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to delete category: %v", err)
	}
	return &pb.EmptyResponse{Success: true}, nil
}
