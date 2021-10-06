package entities

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const CollectionAdministrators = "administrators"

type Administrator struct {
	// System
	Id primitive.ObjectID `bson:"_id"`
	CreatedAt time.Time `bson:"createdAt"`
	ModifiedAt time.Time `bson:"modifiedAt"`
	DeletedAt *time.Time `bson:"deletedAt"`
	BlockedAt *time.Time `bson:"blockedAt"`

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

	// Account
	Password string `bson:"password"`
	LastLoginAt *time.Time `bson:"lastLoginAt"`
	LoginAttemptCounter int64 `bson:"loginAttemptCounter"`
}