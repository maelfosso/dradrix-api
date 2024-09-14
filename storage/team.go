package storage

import (
	"context"
	"time"

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

type AddMemberIntoOrganizationParams struct {
	OrganizationId primitive.ObjectID
	UserId         primitive.ObjectID
	InvitedAt      time.Time
	ConfirmedAt    *time.Time
}

func (q *Queries) AddMemberIntoOrganization(ctx context.Context, arg AddMemberIntoOrganizationParams) (*models.Member, error) {
	var member models.Member = models.Member{
		Id:             primitive.NewObjectID(),
		OrganizationId: arg.OrganizationId,
		MemberId:       arg.UserId,

		InvitedAt:   arg.InvitedAt,
		ConfirmedAt: arg.ConfirmedAt,
		DeletedAt:   nil,

		Status: "confirmed",
		Role:   "member",
	}

	_, err := q.teamsCollection.InsertOne(ctx, member)
	if err != nil {
		return nil, err
	} else {
		return &member, nil
	}
}
