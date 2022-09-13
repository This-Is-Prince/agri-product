package services

import (
	"context"
	"fmt"
	"math"
	"strings"

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

	if shopIdHex := req.GetShopId(); shopIdHex != "" {
		shopIdObject, err := primitive.ObjectIDFromHex(shopIdHex)
		if err != nil {
			return status.Errorf(codes.InvalidArgument, fmt.Sprintf("Could not convert shop hex id to ObjectId: %v", err))
		}

		shopChan, errChan := s.db.Shop().FindOne(bson.M{"_id": shopIdObject})

		select {
		case shop := <-shopChan:
			productCatalogueId := shop.ProductCatalogueId
			productCatalogueChan, errChan := s.db.ProductCatalogue().FindOne(bson.M{"_id": productCatalogueId})

			select {
			case productCatalogue := <-productCatalogueChan:
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
				return nil
			case err := <-errChan:
				return status.Errorf(codes.NotFound, fmt.Sprintf("Could not find the shop catalogue with id %s: %v", productCatalogueId.Hex(), err))
			}
		case err := <-errChan:
			return status.Errorf(codes.NotFound, fmt.Sprintf("Could not find the shop with id %s: %v", req.GetShopId(), err))
		}
	}

	productCataloguesChan, errChan := s.db.ProductCatalogue().Find(filters, bson.D{}, 0, 0)
	select {
	case productCatalogues := <-productCataloguesChan:
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
	case err := <-errChan:
		return status.Errorf(codes.NotFound, fmt.Sprintf("Could not find the product: %v", err))
	}
}
