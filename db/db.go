package db

type DB struct{}

func (db *DB) Shop() *ShopRepository {
	repo := NewShopRepo()
	return repo
}

func (db *DB) ProductCatalogue() *ProductCatalogueRepository {
	repo := NewProductCatalogueRepo()
	return repo
}
