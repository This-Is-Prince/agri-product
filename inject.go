package main

import (
	"github.com/This-Is-Prince/agri-product/db"
	services "github.com/This-Is-Prince/agri-product/services"
)

type Inject struct {
	Db            *db.DB
	SearchService *services.SearchService
}

func NewInject() *Inject {
	inj := &Inject{}
	inj.Db = &db.DB{}

	inj.SearchService = services.NewSearchService(inj.Db)
	return inj
}
