package db

import (
	"github.com/SaiNageswarS/go-api-boot/odm"
	"github.com/This-Is-Prince/agri-product/models"
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
