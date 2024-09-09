package storage

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"stockinos.com/api/models"
)

type GetMembersFromOrganizationParams struct {
	OrganizationId primitive.ObjectID
}

func (q *Queries) GetMembersFromOrganization(ctx context.Context, arg GetMembersFromOrganizationParams) ([]models.Member, error) {
	matchStage := bson.D{
		{
			Key: "$match",
			Value: bson.M{
				"organization_id": arg.OrganizationId,
				"deleted_at":      nil,
			},
		},
	}
	lookupStage := bson.D{
		{
			Key: "$lookup",
			Value: bson.M{
				"from":         "users",
				"localField":   "member_id",
				"foreignField": "_id",
				"as":           "user",
			},
		},
	}
	unwindStage := bson.D{
		{
			Key: "$unwind",
			Value: bson.D{
				{Key: "path", Value: "$user"},
				{Key: "preserveNullAndEmptyArrays", Value: false},
			},
		},
	}

	showLoadedCursor, err := q.teamsCollection.Aggregate(
		ctx,
		mongo.Pipeline{
			matchStage,
			lookupStage,
			unwindStage,
		},
	)
	if err != nil {
		return nil, err
	}
	var members []models.Member
	if err = showLoadedCursor.All(ctx, &members); err != nil {
		return nil, err
	}

	return members, nil
}
