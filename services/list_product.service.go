package services

import (
	"context"
	"math"
	"strings"

	"github.com/This-Is-Prince/agri-product/db"
	"github.com/This-Is-Prince/agri-product/models"
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

	var shop *models.ShopModel
	if shopIdHex := req.GetShopId(); shopIdHex != "" {
		shopIdObject, err := primitive.ObjectIDFromHex(shopIdHex)
		if err != nil {
			return status.Errorf(codes.InvalidArgument, "Invalid shop id")
		}

		shop = <-s.db.Shop().FindShopById(shopIdObject)
	}

	var productCatalogues []models.ProductCatalogueModel
	if shop != nil {
		productCatalogues = <-s.db.ProductCatalogue().FindProductCatalogues(bson.M{"_id": shop.ProductCatalogueId}, bson.D{}, 0, 0)
	} else {
		productCatalogues = <-s.db.ProductCatalogue().FindProductCatalogues(bson.M{}, bson.D{}, 0, 0)
	}

	price_gte := req.GetPriceGte()
	price_lte := req.GetPriceLte()
	if price_lte == 0 {
		price_lte = math.MaxFloat64
	}

	weight_gte := req.GetWeightGte()
	weight_lte := req.GetWeightLte()
	if weight_lte == 0 {
		weight_lte = math.MaxFloat64
	}

	nameQuery := strings.ToLower(strings.TrimSpace(req.GetName()))

	for _, productCatalogue := range productCatalogues {
		for _, product := range productCatalogue.Products {
			productName := strings.ToLower(strings.TrimSpace(product.Name))

			if product.Price >= price_gte && product.Price <= price_lte && product.Weight >= weight_gte && product.Weight <= weight_lte && strings.Contains(productName, nameQuery) {
				stream.Send(&pb.ListProductRes{
					Product: &pb.Product{
						Id:          product.ID.Hex(),
						Name:        product.Name,
						Description: product.Description,
						Price:       product.Price,
						Weight:      product.Weight,
					},
				})
			}
		}
	}
	return nil
}
