package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type ProductCatalogueModel struct {
	ID         primitive.ObjectID   `bson:"_id,omitempty" json:"_id,omitempty"`
	ProductIds []primitive.ObjectID `bson:"product_ids" json:"product_ids"`
}

func (pc *ProductCatalogueModel) Id() string {
	return pc.ID.Hex()
}
