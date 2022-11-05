package main

import "go.mongodb.org/mongo-driver/bson/primitive"

/* ==============================================
Copyright (c) Eensymachines
Developed by 		: kneerunjun@gmail.com
Developed on 		: OCT'22
Datashape of the models when with accounting
============================================== */

type Account interface {
	GetEmail() string
	GetPhone() string
	GetTitle() string
}

// UserAccount : account data that signifies the user information
// this has a direct representation in the database
type UserAccount struct {
	// ID : maps directly to the _id on the database
	// when inserting new its empty or with zero value, hence omitempty
	// when inserting new the driver generates a new id
	ID    primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Email string             `bson:"email, unique" json:"email"`
	Phone string             `bson:"phone, unique" json:"phone"`
	Title string             `bson:"title" json:"title"`
}

func (ua *UserAccount) GetEmail() string {
	return ua.Email
}
func (ua *UserAccount) GetPhone() string {
	return ua.Phone
}
func (ua *UserAccount) GetTitle() string {
	return ua.Title
}
