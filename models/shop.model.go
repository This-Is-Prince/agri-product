package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type ShopModel struct {
	ID                 primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	ProductCatalogueId primitive.ObjectID `bson:"product_catalogue_id,omitempty" json:"product_catalogue_id,omitempty"`
	Name               string             `bson:"name" json:"name"`
	Location           GeoJson            `bson:"location" json:"location"`
	ProductInventoryId primitive.ObjectID `bson:"product_inventory_id" json:"product_inventory_id"`
}

type GeoJson struct {
	Type        string    `bson:"type" json:"type"`
	Coordinates []float64 `bson:"coordinates" json:"coordinates"`
}

func (m *ShopModel) Id() string {
	return m.ID.Hex()
}
