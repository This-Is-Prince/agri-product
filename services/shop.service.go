package services

import "github.com/This-Is-Prince/agri-product/db"

type ShopService struct {
	db *db.DB
}

func NewShopService(db *db.DB) *ShopService {
	return &ShopService{db}
}
