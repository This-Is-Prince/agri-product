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

type ListShopService struct {
	pb.UnimplementedListShopServiceServer
	db *db.DB
}

func NewListShopService(db *db.DB) *ListShopService {
	return &ListShopService{db: db}
}

func (u *ListShopService) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
	return ctx, nil
}

func (s *ListShopService) ListShop(req *pb.ListShopReq, stream pb.ListShopService_ListShopServer) error {
	long := req.GetLong()
	lat := req.GetLat()
	maxDistance := req.GetMaxDistance()
	fmt.Println(maxDistance)

	geoFilters := bson.M{
		"$nearSphere": bson.M{
			"$geometry": bson.M{
				"type": "Point", "coordinates": []float64{long, lat},
			},
			"$maxDistance": maxDistance,
		},
	}
	filters := bson.M{
		"location": geoFilters,
	}
	shopsChan, errChan := s.db.Shop().Find(filters, bson.D{}, 0, 0)
	select {
	case shops := <-shopsChan:
		for _, shop := range shops {
			stream.Send(&pb.ListShopRes{
				Shop: &pb.Shop{
					Id:   shop.Id(),
					Name: shop.Name,
				},
			})
		}
	case err := <-errChan:
		return status.Errorf(codes.NotFound, fmt.Sprintf("Could not find the shops: %v", err))
	}
	return nil
}
