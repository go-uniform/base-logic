package service

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const (
	CollectionAdministrators = "administrators"
	CollectionAdministratorRoles = "administratorRoles"
)

type AdministratorRole struct {
	// System
	Id primitive.ObjectID `bson:"id"`
	CreatedAt time.Time `bson:"createdAt"`
	ModifiedAt time.Time `bson:"modifiedAt"`
	DeletedAt *time.Time `bson:"deletedAt"`

	// Fields
	Name string `json:"name"`
	AllowTags []string `json:"allowTags"`
	DenyTags []string `json:"denyTags"`
}

type Administrator struct {
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
	Password string `bson:"password"`
	Inverted bool `bson:"inverted"`
	AllowTags []string `bson:"allowTags"`
	DenyTags []string `bson:"denyTags"`

	// Links
	Role *AdministratorRole `bson:"role"`
}