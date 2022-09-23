package services

import (
	"context"
	"strings"

	"github.com/This-Is-Prince/agri-product/db"
	"github.com/This-Is-Prince/agri-product/models"
	"github.com/This-Is-Prince/agri-product/pb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	if long > 180 || long < -180 {
		return nil, status.Error(codes.InvalidArgument, "Invalid longitude value")
	}
	if lat > 90 || lat < -90 {
		return nil, status.Error(codes.InvalidArgument, "Invalid latitude value")
	}

	shop := <-s.db.Shop().FindNearByShop(long, lat)
	return &pb.SearchNearbyShopRes{
		Shop: &pb.Shop{
			Id:   shop.Id(),
			Name: shop.Name,
		},
	}, nil
}

func (s *SearchService) SearchByProduct(ctx context.Context, req *pb.SearchByProductReq) (*pb.SearchByProductRes, error) {
	var err error
	var shop *models.ShopModel

	var shopIdObject primitive.ObjectID
	if shopIdHex := req.GetShopId(); shopIdHex != "" {
		shopIdObject, err = primitive.ObjectIDFromHex(shopIdHex)

		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "Invalid shop id")
		}

		shop = <-s.db.Shop().FindShopById(shopIdObject)
	}

	var productCatalogues []models.ProductCatalogueModel
	if shop != nil {
		productCatalogues = <-s.db.ProductCatalogue().FindProductCatalogues(bson.M{"_id": shop.ProductCatalogueId}, bson.D{}, 0, 0)
	} else {
		productCatalogues = <-s.db.ProductCatalogue().FindProductCatalogues(bson.M{}, bson.D{}, 0, 0)
	}

	var productIdObject primitive.ObjectID
	if productIdHex := req.GetProductId(); productIdHex != "" {
		productIdObject, err = primitive.ObjectIDFromHex(productIdHex)

		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "Invalid product id")
		}
	}

	nameQuery := strings.ToLower(strings.TrimSpace(req.GetName()))
	product := <-s.db.ProductCatalogue().FindProductFromProductCatalogues(productCatalogues, productIdObject, nameQuery)

	return &pb.SearchByProductRes{
		Product: &pb.Product{
			Id:          product.ID.Hex(),
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			Weight:      product.Weight,
		},
	}, nil
}
