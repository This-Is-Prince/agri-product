package services

import (
	"context"
	"strings"

	"github.com/This-Is-Prince/agri-product/db"
	"github.com/This-Is-Prince/agri-product/pb"
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
	var err error
	var shopId primitive.ObjectID

	if id := req.GetId(); id != "" {
		shopId, err = primitive.ObjectIDFromHex(id)
		if err != nil {
			return status.Errorf(codes.InvalidArgument, "Invalid shop id")
		}
	}

	long := req.GetLong()
	lat := req.GetLat()
	maxDistance := req.GetMaxDistance()

	if long > 180 || long < -180 {
		return status.Error(codes.InvalidArgument, "Invalid longitude value")
	}
	if lat > 90 || lat < -90 {
		return status.Error(codes.InvalidArgument, "Invalid latitude value")
	}

	shops := <-s.db.Shop().FindShops(shopId, long, lat, maxDistance)

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
	return nil
}
