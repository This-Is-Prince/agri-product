package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type ProductModel struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	Price       float64            `bson:"price" json:"price"` // Assuming Rupees
	Weight      float64            `bson:"weight" json:"weight"`
}

func (m *ProductModel) Id() string {
	return m.ID.Hex()
}
