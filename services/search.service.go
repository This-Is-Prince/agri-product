package services

import (
	"context"
	"fmt"
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

func (s *SearchService) SearchByProduct(ctx context.Context, req *pb.SearchByProductReq) (*pb.SearchByProductRes, error) {
	nameQuery := strings.ToLower(strings.TrimSpace(req.GetName()))

	if shopIdHex := req.GetShopId(); shopIdHex != "" {
		shopIdObject, err := primitive.ObjectIDFromHex(shopIdHex)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Could not convert shop hex id to ObjectId: %v", err))
		}
		shopChan, errChan := s.db.Shop().FindOne(bson.M{"_id": shopIdObject})

		select {
		case shop := <-shopChan:
			productCatalogueId := shop.ProductCatalogueId
			productCataloguesChan, errChan := s.db.ProductCatalogue().Find(bson.M{"_id": productCatalogueId}, bson.D{}, 0, 0)

			select {
			case productCatalogues := <-productCataloguesChan:
				productChan, errChan := findProduct(productCatalogues, req.GetProductId(), nameQuery)
				select {
				case product := <-productChan:
					return &pb.SearchByProductRes{
						Product: product,
					}, nil
				case err := <-errChan:
					return nil, err
				}
			case err := <-errChan:
				return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Could not find the product catalogue with id %s: %v", productCatalogueId.Hex(), err))
			}
		case err := <-errChan:
			return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Could not find the shop with id %s: %v", shopIdHex, err))
		}
	}
	productCataloguesChan, errChan := s.db.ProductCatalogue().Find(bson.M{}, bson.D{}, 0, 0)
	select {
	case productCatalogues := <-productCataloguesChan:
		productChan, errChan := findProduct(productCatalogues, req.GetProductId(), nameQuery)
		select {
		case product := <-productChan:
			return &pb.SearchByProductRes{
				Product: product,
			}, nil
		case err := <-errChan:
			return nil, err
		}
	case err := <-errChan:
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Could not find the product with Id %s: %v", req.GetProductId(), err))
	}
}

func findProduct(productCatalogues []models.ProductCatalogueModel, productId, nameQuery string) (chan *pb.Product, chan error) {
	product := make(chan *pb.Product)
	err := make(chan error)

	go func() {
		var result *pb.Product
		results := make([]*pb.Product, 0)
		for _, productCatalogue := range productCatalogues {
			for _, product := range productCatalogue.Products {
				tmp := pb.Product{
					Id:          product.ID.Hex(),
					Name:        product.Name,
					Description: product.Description,
					Price:       product.Price,
					Weight:      product.Weight,
				}

				if product.ID.Hex() == productId {
					result = &tmp
				}

				productName := strings.ToLower(strings.TrimSpace(product.Name))
				if strings.Contains(productName, nameQuery) {
					results = append(results, &tmp)
				}
			}
		}
		if result != nil {
			product <- result
		} else if len(results) > 0 {
			product <- results[0]
		} else {
			err <- status.Error(codes.NotFound, fmt.Sprintf("Could not find the product with name: %s", nameQuery))
		}
	}()

	return product, err
}
