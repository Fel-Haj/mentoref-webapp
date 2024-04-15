package db

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID          primitive.ObjectID `bson:"_id"`
	Email       string             `bson:"email"`
	Password    string             `bson:"password"`
	FirstName   string             `bson:"first_name"`
	Surname     string             `bson:"surname"`
	Phone       string             `bson:"phone,omitempty"`
	CompanyName string             `bson:"company_name,omitempty"`
}
