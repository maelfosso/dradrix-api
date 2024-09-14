package handlers_test

import (
	gofaker "github.com/go-faker/faker/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"stockinos.com/api/models"
)

var authenticatedUser = &models.User{
	Id:          primitive.NewObjectID(),
	FirstName:   gofaker.FirstName(),
	LastName:    gofaker.LastName(),
	PhoneNumber: gofaker.Phonenumber(),
}
