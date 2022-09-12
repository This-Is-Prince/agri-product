package services

import (
	"context"
	"fmt"

	"github.com/This-Is-Prince/agri-product/db"
	"github.com/This-Is-Prince/agri-product/pb"
	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ListProductService struct {
	pb.UnimplementedListProductServiceServer
	db *db.DB
}

func NewListProductService(db *db.DB) *ListProductService {
	return &ListProductService{db: db}
}

func (u *ListProductService) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
	return ctx, nil
}

func (s *ListProductService) ListProduct(req *pb.ListProductReq, stream pb.ListProductService_ListProductServer) error {
	productsChan, errChan := s.db.Product().Find(bson.M{}, bson.D{}, 0, 0)
	select {
	case products := <-productsChan:
		for _, product := range products {
			stream.Send(&pb.ListProductRes{
				Product: &pb.Product{
					Id:          product.Id(),
					Name:        product.Name,
					Description: product.Description,
					Price:       product.Price,
					Weight:      product.Weight,
				},
			})
		}
	case err := <-errChan:
		return status.Errorf(codes.NotFound, fmt.Sprintf("Could not find the product: %v", err))
	}
	return nil
}
