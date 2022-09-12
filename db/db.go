package db

type DB struct{}

func (db *DB) Shop() *ShopRepository {
	repo := NewShopRepo()
	return repo
}

func (db *DB) Product() *ProductRepository {
	repo := NewProductRepo()
	return repo
}
