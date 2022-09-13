package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type ProductCatalogueModel struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Products []Product          `bson:"products" json:"products"`
}

type Product struct {
	ID          primitive.ObjectID `bson:"id,omitempty" json:"id,omitempty"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	Price       float64            `bson:"price" json:"price"` // Assuming Rupees
	Weight      float64            `bson:"weight" json:"weight"`
}

func (pc *ProductCatalogueModel) Id() string {
	return pc.ID.Hex()
}
