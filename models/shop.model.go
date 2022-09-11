package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type ShopModel struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Name         string             `bson:"name" json:"name"`
	Location     GeoJson            `bson:"location" json:"location"`
	InventoryIds primitive.ObjectID `bson:"inventoryIds" json:"inventoryIds"`
}

type GeoJson struct {
	Type        string    `json:"-"`
	Coordinates []float64 `json:"coordinates"`
}

func (m *ShopModel) Id() string {
	return m.ID.Hex()
}
