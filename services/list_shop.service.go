package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/This-Is-Prince/agri-product/db"
	"github.com/This-Is-Prince/agri-product/pb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	filters := bson.M{}
	long := req.GetLong()
	lat := req.GetLat()
	maxDistance := req.GetMaxDistance()

	if id := req.GetId(); id != "" {
		oid, err := primitive.ObjectIDFromHex(id)
		// Check for errors
		if err != nil {
			return status.Errorf(codes.InvalidArgument, fmt.Sprintf("Could not convert to ObjectId: %v", err))
		}
		filters["_id"] = oid
	}

	geoFilters := bson.M{
		"$nearSphere": bson.M{
			"$geometry": bson.M{
				"type": "Point", "coordinates": []float64{long, lat},
			},
			"$maxDistance": maxDistance,
		},
	}
	filters["location"] = geoFilters

	shopsChan, errChan := s.db.Shop().Find(filters, bson.D{}, 0, 0)
	select {
	case shops := <-shopsChan:
		for _, shop := range shops {
			if strings.Contains(strings.TrimSpace(strings.ToLower(shop.Name)), strings.TrimSpace(strings.ToLower(req.GetName()))) {
				stream.Send(&pb.ListShopRes{
					Shop: &pb.Shop{
						Id:   shop.Id(),
						Name: shop.Name,
					},
				})
			}
		}
	case err := <-errChan:
		return status.Errorf(codes.NotFound, fmt.Sprintf("Could not find the shops: %v", err))
	}
	return nil
}
