package handler

import (
	"context"
	pb "github.com/ReStorePUC/protobucket/product"
)

// ProductServer is used to implement Product.
type ProductServer struct {
	pb.UnimplementedProductServer

	controller controller
}

func NewProductServer(c controller) *ProductServer {
	return &ProductServer{
		controller: c,
	}
}

func (s *ProductServer) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.GetProductResponse, error) {
	product, err := s.controller.GetProduct(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	images := []*pb.Image{}
	for _, img := range product.Images {
		images = append(images, &pb.Image{
			Id:        int32(img.ID),
			ImagePath: img.ImagePath,
			ProductId: int32(img.ProductID),
		})
	}

	return &pb.GetProductResponse{
		Id:          int32(product.ID),
		Name:        product.Name,
		Description: product.Description,
		Categories:  product.Categories,
		Size:        product.Size,
		Price:       float32(product.Price),
		Tax:         float32(product.Tax),
		Available:   product.Available,
		StoreId:     int32(product.StoreID),
		Images:      images,
	}, nil
}

func (s *ProductServer) UnavailableProduct(ctx context.Context, req *pb.UnavailableProductRequest) (*pb.UnavailableProductResponse, error) {
	err := s.controller.UnavailablePayment(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.UnavailableProductResponse{}, nil
}
