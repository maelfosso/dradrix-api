package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"stockinos.com/api/models"
)

type CreateActivityParams struct {
	Name        string
	Description string
	Fields      []models.ActivityField

	OrganizationId primitive.ObjectID
	CreatedBy      primitive.ObjectID
}

func (q *Queries) CreateActivity(ctx context.Context, arg CreateActivityParams) (*models.Activity, error) {
	var activity models.Activity = models.Activity{
		Id:          primitive.NewObjectID(),
		Name:        arg.Name,
		Description: arg.Description,
		Fields:      arg.Fields, // Default to [] empty array instead of null

		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		// DeletedAt:   null.,

		OrganizationId: arg.OrganizationId,
		CreatedBy:      arg.CreatedBy,
	}

	_, err := q.activitiesCollection.InsertOne(ctx, activity)
	if err != nil {
		return nil, err
	} else {
		return &activity, nil
	}
}

type GetActivityParams struct {
	Id             primitive.ObjectID
	OrganizationId primitive.ObjectID
}

func (q *Queries) GetActivity(ctx context.Context, arg GetActivityParams) (*models.Activity, error) {
	var activity models.Activity

	filter := bson.M{
		"_id":             arg.Id,
		"organization_id": arg.OrganizationId,
		"deleted_at":      nil,
	}
	err := q.activitiesCollection.FindOne(ctx, filter).Decode(&activity)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &activity, nil
}

type GetAllActivitiesParams struct {
	OrganizationId primitive.ObjectID
}

func (q *Queries) GetAllActivities(ctx context.Context, arg GetAllActivitiesParams) ([]*models.Activity, error) {
	var activities []*models.Activity

	filter := bson.M{
		"organization_id": arg.OrganizationId,
		"deleted_at":      nil,
	}
	cursor, err := q.activitiesCollection.Find(ctx, filter)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		} else {
			return nil, err
		}
	}
	if err = cursor.All(ctx, &activities); err != nil {
		return nil, err
	}
	return activities, nil
}

type DeleteActivityParams struct {
	Id             primitive.ObjectID
	OrganizationId primitive.ObjectID
}

func (q *Queries) DeleteActivity(ctx context.Context, arg DeleteActivityParams) error {
	filter := bson.M{
		"_id":             arg.Id,
		"organization_id": arg.OrganizationId,
	}
	update := bson.M{
		"$set": bson.M{
			"deleted_at": time.Now(),
		},
	}

	var activity models.Activity
	err := q.activitiesCollection.FindOneAndUpdate(
		ctx,
		filter,
		update,
	).Decode(&activity)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil
		} else {
			return err
		}
	}

	return nil
}

type UpdateSetInActivityParams struct {
	Id             primitive.ObjectID
	OrganizationId primitive.ObjectID

	FieldsToSet map[string]any
}

func (q *Queries) UpdateSetInActivity(ctx context.Context, arg UpdateSetInActivityParams) (*models.Activity, error) {
	filter := bson.M{
		"_id":             arg.Id,
		"organization_id": arg.OrganizationId,
	}

	set := bson.M{}
	for field, value := range arg.FieldsToSet {
		set[field] = value
	}
	update := bson.M{
		"$set": set,
	}

	return CommonUpdateQuery[models.Activity](ctx, *q.activitiesCollection, filter, update)
}

type UpdateAddToActivityParams struct {
	Id             primitive.ObjectID
	OrganizationId primitive.ObjectID

	Position uint
	Field    string
	Value    models.ActivityField // interface{}
}

func (q *Queries) UpdateAddToActivity(ctx context.Context, arg UpdateAddToActivityParams) (*models.Activity, error) {
	filter := bson.M{
		"_id":             arg.Id,
		"organization_id": arg.OrganizationId,
	}
	update := bson.M{
		"$push": bson.M{
			arg.Field: bson.M{
				"$each": bson.A{
					arg.Value,
				},
				"$position": arg.Position,
			},
		},
	}

	return CommonUpdateQuery[models.Activity](ctx, *q.activitiesCollection, filter, update)
}

type UpdateRemoveFromActivityParams struct {
	Id             primitive.ObjectID
	OrganizationId primitive.ObjectID

	Position uint
	Field    string
	// Value interface{}
}

func (q *Queries) UpdateRemoveFromActivity(ctx context.Context, arg UpdateRemoveFromActivityParams) (*models.Activity, error) {
	fieldWithPosition := fmt.Sprintf("%s.%d", arg.Field, arg.Position)
	filter := bson.M{
		"_id":             arg.Id,
		"organization_id": arg.OrganizationId,
	}
	update := bson.M{
		"$unset": bson.M{
			fieldWithPosition: 1,
		},
	}
	_, err := CommonUpdateQuery[models.Activity](ctx, *q.activitiesCollection, filter, update)
	if err != nil {
		return nil, err
	}

	update = bson.M{
		"$pull": bson.M{
			arg.Field: nil,
		},
	}

	return CommonUpdateQuery[models.Activity](ctx, *q.activitiesCollection, filter, update)
}

type AddRelationshipIntoActivityParams struct {
	Id             primitive.ObjectID
	OrganizationId primitive.ObjectID

	Type             string
	ActivityId       primitive.ObjectID
	FieldId          primitive.ObjectID
	ConcernedFieldId primitive.ObjectID
}

func (q *Queries) AddRelationshipIntoActivity(ctx context.Context, arg AddRelationshipIntoActivityParams) (*models.Activity, error) {
	relationship := models.ActivityRelationship{
		Id:               primitive.NewObjectID(),
		Type:             arg.Type,
		ActivityId:       arg.ActivityId,
		FieldId:          arg.FieldId,
		ConcernedFieldId: arg.ConcernedFieldId,
	}

	filter := bson.M{
		"_id":             arg.Id,
		"organization_id": arg.OrganizationId,
	}
	update := bson.M{
		"$push": bson.M{
			"relationships": bson.M{
				"$each": bson.A{
					relationship,
				},
			},
		},
	}

	return CommonUpdateQuery[models.Activity](ctx, *q.activitiesCollection, filter, update)
}

type RemoveRelationshipFromActivityParams struct {
	Id             primitive.ObjectID
	OrganizationId primitive.ObjectID

	Type             string
	ActivityId       primitive.ObjectID
	FieldId          primitive.ObjectID
	ConcernedFieldId primitive.ObjectID
}

func (q *Queries) RemoveRelationshipFromActivity(ctx context.Context, arg RemoveRelationshipFromActivityParams) (*models.Activity, error) {
	filter := bson.M{
		"_id":             arg.Id,
		"organization_id": arg.OrganizationId,
	}
	// update := bson.M{
	// 	"$unset": bson.M{
	// 		// fieldWithPosition: 1,
	// 	},
	// }
	// _, err := CommonUpdateQuery[models.Activity](ctx, *q.activitiesCollection, filter, update)
	// if err != nil {
	// 	return nil, err
	// }

	update := bson.M{
		"$pull": bson.M{
			"relationships": bson.M{
				// "type": arg.Type,
				"activity_id":        arg.ActivityId,
				"field_id":           arg.FieldId,
				"concerned_field_id": arg.ConcernedFieldId,
			},
		},
	}

	_, err := CommonUpdateQuery[models.Activity](ctx, *q.activitiesCollection, filter, update)
	if err != nil {
		return nil, err
	}

	// Other side of the relationship
	filter = bson.M{
		"_id":             arg.ActivityId,
		"organization_id": arg.OrganizationId,
	}

	update = bson.M{
		"$pull": bson.M{
			"relationships": bson.M{
				"activity_id":        arg.Id,
				"field_id":           arg.ConcernedFieldId,
				"concerned_field_id": arg.FieldId,
			},
		},
	}

	return CommonUpdateQuery[models.Activity](ctx, *q.activitiesCollection, filter, update)
}
