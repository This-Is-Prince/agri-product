package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type ProductInventoryModel struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	ProductStocks []ProductStock     `bson:"product_stocks" json:"product_stocks"`
}

type ProductStock struct {
	Quantity  uint32             `bson:"quantity" json:"quantity"`
	ProductID primitive.ObjectID `bson:"product_id" json:"product_id"`
}

func (m *ProductInventoryModel) Id() string {
	return m.ID.Hex()
}
