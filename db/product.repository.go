package db

import (
	"github.com/SaiNageswarS/go-api-boot/odm"
	"github.com/This-Is-Prince/agri-product/models"
)

type ProductRepository struct {
	odm.AbstractRepository[models.ProductModel]
}

func NewProductRepo() *ProductRepository {
	repo := odm.AbstractRepository[models.ProductModel]{
		Database:       "agri-product",
		CollectionName: "product",
	}
	return &ProductRepository{repo}
}
