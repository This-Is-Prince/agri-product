package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type ProductInventory struct {
	ID          primitive.ObjectID   `bson:"_id,omitempty" json:"_id,omitempty"`
	ProductsIDs []primitive.ObjectID `bson:"productsIds" json:"productsIds"`
}

type Products struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Quantity  uint32             `bson:"quantity" json:"quantity"`
	ProductID primitive.ObjectID `bson:"productId" json:"productId"`
}
