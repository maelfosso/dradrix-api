package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Address struct {
	Street     string `bson:"street" json:"street,omitempty"`
	City       string `bson:"city" json:"city,omitempty"`
	Region     string `bson:"region" json:"region,omitempty"`
	PostalCode string `bson:"postal_code" json:"postal_code,omitempty"`
	Country    string `bson:"country" json:"country,omitempty"`
}

type Organization struct {
	Id      primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Name    string             `bson:"name" json:"name,omitempty"`
	Bio     string             `bson:"bio" json:"bio,omitempty"`
	Email   string             `bson:"email" json:"email,omitempty"`
	Address Address            `bson:"address" json:"address,omitempty"`

	CreatedAt time.Time  `bson:"created_at" json:"created_at,omitempty"`
	UpdatedAt time.Time  `bson:"updated_at" json:"updated_at,omitempty"`
	DeletedAt *time.Time `bson:"deleted_at" json:"deleted_at,omitempty"`

	CreatedBy primitive.ObjectID `bson:"created_by" json:"created_by,omitempty"`

	OwnedBy         primitive.ObjectID `bson:"owned_by" json:"owned_by"`
	InvitationToken string             `bson:"invitation_token" json:"invitation_token"`
}

type Member struct {
	Id             primitive.ObjectID `bson:"_id" json:"id"`
	OrganizationId primitive.ObjectID `bson:"organization_id" json:"organization_id"`

	MemberId primitive.ObjectID `bson:"member_id"`
	User     User               `bson:"user,omitempty" json:"user"`

	InvitedAt   time.Time  `bson:"invited_at" json:"invited_at"`
	ConfirmedAt *time.Time `bson:"confirmed_at" json:"confirmed_at"`
	DeletedAt   *time.Time `bson:"deleted_at" json:"deleted_at"`

	Status string `bson:"status" json:"status"`
	Role   string `bson:"role" json:"role"`
}
