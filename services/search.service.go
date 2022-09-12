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

type SearchService struct {
	pb.UnimplementedSearchServiceServer
	db *db.DB
}

func NewSearchService(db *db.DB) *SearchService {
	return &SearchService{db: db}
}

func (u *SearchService) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
	return ctx, nil
}

func (s *SearchService) SearchNearbyShop(ctx context.Context, req *pb.SearchNearbyShopReq) (*pb.SearchNearbyShopRes, error) {
	long := req.GetLong()
	lat := req.GetLat()

	shopChan, errChan := s.db.Shop().FindOne(
		bson.M{
			"location": bson.M{
				"$nearSphere": bson.M{
					"$geometry": bson.M{
						"type": "Point", "coordinates": []float64{long, lat},
					},
				},
			},
		},
	)
	select {
	case shop := <-shopChan:
		return &pb.SearchNearbyShopRes{
			Shop: &pb.Shop{
				Id:   shop.Id(),
				Name: shop.Name,
			},
		}, nil
	case err := <-errChan:
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Could not find the near shop: %v", err))
	}
}
