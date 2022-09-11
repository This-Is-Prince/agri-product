package db

type DB struct{}

func (db *DB) Shop() *ShopRepository {
	repo := NewShopRepo()
	return repo
}
