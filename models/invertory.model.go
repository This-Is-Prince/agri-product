package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type ProductInventoryModel struct {
	ID          primitive.ObjectID   `bson:"_id,omitempty" json:"_id,omitempty"`
	ProductsIDs []primitive.ObjectID `bson:"productsIds" json:"productsIds"`
}

type ProductsModel struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Quantity  uint32             `bson:"quantity" json:"quantity"`
	ProductID primitive.ObjectID `bson:"productId" json:"productId"`
}

func (m *ProductInventoryModel) Id() string {
	return m.ID.Hex()
}

func (m *ProductsModel) Id() string {
	return m.ID.Hex()
}
