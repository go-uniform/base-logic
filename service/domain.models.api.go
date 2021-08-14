package service

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type ApiAdministratorsRequest struct {
	// Fields
	FirstName string `bson:"firstName"`
	LastName string `bson:"lastName"`
	Email string `bson:"email"`
	Mobile string `bson:"mobile"`
	Inverted bool `bson:"inverted"`
	AllowTags []string `bson:"allowTags"`
	DenyTags []string `bson:"denyTags"`

	// Links
	Role *string `bson:"role"`
}

type ApiAdministratorsResponse struct {
	// System
	Id primitive.ObjectID `bson:"id"`
	CreatedAt time.Time `bson:"createdAt"`
	ModifiedAt time.Time `bson:"modifiedAt"`
	DeletedAt *time.Time `bson:"deletedAt"`

	// Fields
	FirstName string `bson:"firstName"`
	LastName string `bson:"lastName"`
	Email string `bson:"email"`
	Mobile string `bson:"mobile"`
	Inverted bool `bson:"inverted"`
	AllowTags []string `bson:"allowTags"`
	DenyTags []string `bson:"denyTags"`

	// Links
	Role *AdministratorRole `bson:"role"`
}

type ApiAdministratorsResponseList struct {
	// System
	Id primitive.ObjectID `bson:"id"`
	CreatedAt time.Time `bson:"createdAt"`
	ModifiedAt time.Time `bson:"modifiedAt"`
	DeletedAt *time.Time `bson:"deletedAt"`

	// Fields
	FirstName string `bson:"firstName"`
	LastName string `bson:"lastName"`
	Email string `bson:"email"`
	Mobile string `bson:"mobile"`
	Inverted bool `bson:"inverted"`
	AllowTags []string `bson:"allowTags"`
	DenyTags []string `bson:"denyTags"`

	// Links
	Role *AdministratorRole `bson:"role"`
}