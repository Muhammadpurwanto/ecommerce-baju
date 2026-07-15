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

type ProductGrpcHandler struct {
	pb.UnimplementedProductServiceServer
	productService service.ProductService
	validate       *validator.Validate
}

func NewProductGrpcHandler(productService service.ProductService) *ProductGrpcHandler {
	return &ProductGrpcHandler{
		productService: productService,
		validate:       validator.New(),
	}
}

func toPbProductResponse(p *model.Product) *pb.ProductResponse {
	return &pb.ProductResponse{
		Id:          uint32(p.ID),
		CategoryId:  uint32(p.CategoryID),
		Name:        p.Name,
		Slug:        p.Slug,
		Description: p.Description,
		Brand:       p.Brand,
		Gender:      p.Gender,
		BasePrice:   p.BasePrice,
		Weight:      p.Weight,
		IsActive:    p.IsActive,
		ImageUrl:    p.ImageURL,
		Stock:       int32(p.Stock),
		CreatedAt:   p.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   p.UpdatedAt.Format(time.RFC3339),
	}
}

func (h *ProductGrpcHandler) GetAllProducts(ctx context.Context, req *pb.GetAllProductsRequest) (*pb.ProductListResponse, error) {
	page := int(req.GetPage())
	if page < 1 {
		page = 1
	}
	perPage := int(req.GetPerPage())
	if perPage < 1 {
		perPage = 20
	}

	products, err := h.productService.GetAll(page, perPage)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get products: %v", err)
	}

	var pbData []*pb.ProductResponse
	for _, p := range products {
		val := p
		pbData = append(pbData, toPbProductResponse(&val))
	}

	return &pb.ProductListResponse{
		Data: pbData,
		Page: int32(page),
		PerPage: int32(perPage),
		Total: 0, 
		TotalPages: 0,
	}, nil
}

func (h *ProductGrpcHandler) GetProductByID(ctx context.Context, req *pb.GetByIDRequest) (*pb.ProductResponse, error) {
	if req.GetId() == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "id is required")
	}

	product, err := h.productService.GetByID(uint(req.GetId()))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "product not found: %v", err)
	}
	return toPbProductResponse(product), nil
}

func (h *ProductGrpcHandler) CreateProduct(ctx context.Context, req *pb.ProductRequest) (*pb.ProductResponse, error) {
	dtoReq := &dto.ProductRequest{
		CategoryID:  uint(req.GetCategoryId()),
		Name:        req.GetName(),
		Description: req.GetDescription(),
		Brand:       req.GetBrand(),
		Gender:      req.GetGender(),
		BasePrice:   req.GetBasePrice(),
		Weight:      req.GetWeight(),
		Stock:       int(req.GetStock()),
		ImageURL:    req.GetImageUrl(),
		IsActive:    req.GetIsActive(),
	}

	if err := h.validate.Struct(dtoReq); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validation error: %v", err)
	}

	product, err := h.productService.Create(dtoReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create product: %v", err)
	}
	return toPbProductResponse(product), nil
}

func (h *ProductGrpcHandler) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.ProductResponse, error) {
	if req.GetId() == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "id is required")
	}

	dtoReq := &dto.ProductRequest{
		CategoryID:  uint(req.GetCategoryId()),
		Name:        req.GetName(),
		Description: req.GetDescription(),
		Brand:       req.GetBrand(),
		Gender:      req.GetGender(),
		BasePrice:   req.GetBasePrice(),
		Weight:      req.GetWeight(),
		Stock:       int(req.GetStock()),
		ImageURL:    req.GetImageUrl(),
		IsActive:    req.GetIsActive(),
	}

	if err := h.validate.Struct(dtoReq); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validation error: %v", err)
	}

	product, err := h.productService.Update(uint(req.GetId()), dtoReq)
	if err != nil {
		if err.Error() == "record not found" || err.Error() == "product not found" {
			return nil, status.Errorf(codes.NotFound, "product not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to update product: %v", err)
	}
	return toPbProductResponse(product), nil
}

func (h *ProductGrpcHandler) DeleteProduct(ctx context.Context, req *pb.GetByIDRequest) (*pb.EmptyResponse, error) {
	if req.GetId() == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "id is required")
	}

	if err := h.productService.Delete(uint(req.GetId())); err != nil {
		if err.Error() == "record not found" || err.Error() == "product not found" {
			return nil, status.Errorf(codes.NotFound, "product not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to delete product: %v", err)
	}
	return &pb.EmptyResponse{Success: true}, nil
}

func (h *ProductGrpcHandler) ReserveStock(ctx context.Context, req *pb.ReserveStockRequest) (*pb.StockResponse, error) {
	var dtoItems []dto.StockItem
	for _, item := range req.GetItems() {
		dtoItems = append(dtoItems, dto.StockItem{
			ProductID: uint(item.GetProductId()),
			Quantity:  int(item.GetQuantity()),
		})
	}

	if err := h.productService.ReserveStock(dtoItems); err != nil {
		return &pb.StockResponse{Success: false, Message: err.Error()}, nil
	}
	return &pb.StockResponse{Success: true, Message: "Stock reserved successfully"}, nil
}

func (h *ProductGrpcHandler) ReleaseStock(ctx context.Context, req *pb.ReleaseStockRequest) (*pb.StockResponse, error) {
	var dtoItems []dto.StockItem
	for _, item := range req.GetItems() {
		dtoItems = append(dtoItems, dto.StockItem{
			ProductID: uint(item.GetProductId()),
			Quantity:  int(item.GetQuantity()),
		})
	}

	if err := h.productService.ReleaseStock(dtoItems); err != nil {
		return &pb.StockResponse{Success: false, Message: err.Error()}, nil
	}
	return &pb.StockResponse{Success: true, Message: "Stock released successfully"}, nil
}
