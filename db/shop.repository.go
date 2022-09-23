package db

import (
	"github.com/SaiNageswarS/go-api-boot/logger"
	"github.com/SaiNageswarS/go-api-boot/odm"
	"github.com/This-Is-Prince/agri-product/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

type ShopRepository struct {
	odm.AbstractRepository[models.ShopModel]
}

func NewShopRepo() *ShopRepository {
	repo := odm.AbstractRepository[models.ShopModel]{
		Database:       "agri-product",
		CollectionName: "shop",
	}
	return &ShopRepository{repo}
}

func (s *ShopRepository) FindShops(shopId primitive.ObjectID, long, lat, maxDistance float64) chan []models.ShopModel {
	c := make(chan []models.ShopModel)

	go func() {
		filters := bson.M{}
		if !shopId.IsZero() {
			filters["_id"] = shopId
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
		shopsChan, errChan := s.Find(filters, bson.D{}, 0, 0)
		select {
		case shops := <-shopsChan:
			c <- shops
		case err := <-errChan:
			logger.Error("Error finding the shops", zap.Error(err))
		}
	}()
	return c
}

func (s *ShopRepository) FindNearByShop(long, lat float64) chan *models.ShopModel {
	shopChan := make(chan *models.ShopModel)
	go func() {
		shopChan, errChan := s.FindOne(
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
			shopChan <- shop
		case err := <-errChan:
			logger.Error("Error finding near by shop", zap.Error(err))
		}
	}()
	return shopChan
}

func (s *ShopRepository) FindShopById(shopId primitive.ObjectID) chan *models.ShopModel {
	shopChan := make(chan *models.ShopModel)
	go func() {
		shopChan, errChan := s.FindOne(bson.M{"_id": shopId})
		select {
		case shop := <-shopChan:
			shopChan <- shop
		case err := <-errChan:
			logger.Error("Error finding the shop by id", zap.Error(err))
		}
	}()
	return shopChan
}
