package entities

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const CollectionAdministratorRoles = "administratorRoles"

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