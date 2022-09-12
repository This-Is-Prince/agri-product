package services

import (
	"context"
	"fmt"

	"github.com/This-Is-Prince/agri-product/db"
	"github.com/This-Is-Prince/agri-product/pb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	filters := bson.M{}

	if id := req.GetId(); id != "" {
		oid, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return status.Errorf(codes.InvalidArgument, fmt.Sprintf("Could not convert to ObjectId: %v", err))
		}
		filters["_id"] = oid
	}

	price_gte := req.GetPriceGte()
	price_lte := req.GetPriceLte()
	price := bson.M{}
	if price_gte != 0 && price_lte != 0 {
		price["price"] = bson.M{"$gte": price_gte, "$lte": price_lte}
	} else if price_gte != 0 {
		price["price"] = bson.M{"$gte": price_gte}
	} else if price_lte != 0 {
		price["price"] = bson.M{"$lte": price_lte}
	}

	weight_gte := req.GetWeightGte()
	weight_lte := req.GetWeightLte()
	weight := bson.M{}
	if weight_gte != 0 && weight_lte != 0 {
		weight["weight"] = bson.M{"$gte": weight_gte, "$lte": weight_lte}
	} else if weight_gte != 0 {
		weight["weight"] = bson.M{"$gte": weight_gte}
	} else if weight_lte != 0 {
		weight["weight"] = bson.M{"$lte": weight_lte}
	}
	filters["$and"] = []interface{}{price, weight}

	if name := req.GetName(); name != "" {
		filters["$text"] = bson.M{"$search": name}
	}

	productsChan, errChan := s.db.Product().Find(filters, bson.D{}, 0, 0)
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
