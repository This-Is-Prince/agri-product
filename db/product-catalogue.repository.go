package db

import (
	"github.com/SaiNageswarS/go-api-boot/odm"
	"github.com/This-Is-Prince/agri-product/models"
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
