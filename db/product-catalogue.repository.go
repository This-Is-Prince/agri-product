package db

import (
	"strings"

	"github.com/SaiNageswarS/go-api-boot/logger"
	"github.com/SaiNageswarS/go-api-boot/odm"
	"github.com/This-Is-Prince/agri-product/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

type ProductCatalogueRepository struct {
	odm.AbstractRepository[models.ProductCatalogueModel]
}

func NewProductCatalogueRepo() *ProductCatalogueRepository {
	repo := odm.AbstractRepository[models.ProductCatalogueModel]{
		Database:       "agri-product",
		CollectionName: "product-catalogue",
	}
	return &ProductCatalogueRepository{repo}
}

func (p *ProductCatalogueRepository) FindProductCatalogues(filters primitive.M, sort primitive.D, limit int64, skip int64) chan []models.ProductCatalogueModel {
	c := make(chan []models.ProductCatalogueModel)

	go func() {
		productCataloguesChan, errChan := p.Find(filters, sort, limit, skip)
		select {
		case productCatalogues := <-productCataloguesChan:
			c <- productCatalogues
		case err := <-errChan:
			logger.Error("Error finding the product catalogues", zap.Error(err))
		}
	}()
	return c
}

func (p *ProductCatalogueRepository) FindProductFromProductCatalogues(productCatalogues []models.ProductCatalogueModel, productId primitive.ObjectID, nameQuery string) chan *models.Product {
	c := make(chan *models.Product)

	go func() {
		var result *models.Product
		results := make([]*models.Product, 0)
		for _, productCatalogue := range productCatalogues {
			for _, product := range productCatalogue.Products {

				if product.ID == productId {
					result = &product
				}

				productName := strings.ToLower(strings.TrimSpace(product.Name))
				if strings.Contains(productName, nameQuery) {
					results = append(results, &product)
				}
			}
		}
		if result != nil {
			c <- result
		} else if len(results) > 0 {
			c <- results[0]
		} else {
			c <- nil
		}
	}()
	return c
}
