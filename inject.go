package main

import (
	"github.com/This-Is-Prince/agri-product/db"
	services "github.com/This-Is-Prince/agri-product/services"
)

type Inject struct {
	Db          *db.DB
	ShopService *services.ShopService
}

func NewInject() *Inject {
	inj := &Inject{}
	inj.Db = &db.DB{}

	inj.ShopService = services.NewShopService(inj.Db)
	return inj
}
