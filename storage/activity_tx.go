package storage

import (
	"context"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"stockinos.com/api/models"
)

type UpdateSetInActivityTxParams struct {
	Activity       models.Activity
	OrganizationId primitive.ObjectID

	FieldsToSet map[string]any

	Field   string
	Details any // models.ActivityFieldType
}

func (store *MongoStorage) UpdateSetInActivityTx(ctx context.Context, arg UpdateSetInActivityTxParams) (*models.Activity, error) {
	result, err := store.withTx(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		fieldSplitten := strings.Split(arg.Field, ".")
		if len(fieldSplitten) > 1 {
			position, _ := strconv.Atoi(fieldSplitten[1])
			field := arg.Activity.Fields[position]

			if fieldSplitten[len(fieldSplitten)-1] == "details" {
				switch arg.Details.(type) {
				case models.ActivityFieldKey:
					details := arg.Details.(models.ActivityFieldKey)

					// Check if there exists a relationship defined on that field in Activity
					// If yes, delete it
					var fieldRelationship *models.ActivityRelationship = nil
					for i := range arg.Activity.Relationships {
						r := arg.Activity.Relationships[i]
						if r.ConcernedFieldId == field.Id {
							fieldRelationship = &r

							break
						}
					}
					if fieldRelationship != nil {
						_, err := store.RemoveRelationshipFromActivity(ctx, RemoveRelationshipFromActivityParams{
							Id:             arg.Activity.Id,
							OrganizationId: arg.OrganizationId,

							Type:             fieldRelationship.Type,
							ActivityId:       fieldRelationship.ActivityId,
							FieldId:          fieldRelationship.FieldId,
							ConcernedFieldId: fieldRelationship.ConcernedFieldId,
						})

						if err != nil {
							return nil, err
						}
					}

					// Add relationship to activity: id
					// belongs-to
					// type: "belongs-to"
					// activityId: activity_id
					// field_id: field_id
					_, err := store.AddRelationshipIntoActivity(ctx, AddRelationshipIntoActivityParams{
						Id:             arg.Activity.Id,
						OrganizationId: arg.OrganizationId,

						Type:             "belongs_to",
						ActivityId:       details.ActivityId,
						FieldId:          details.FieldId,
						ConcernedFieldId: field.Id,
					})
					if err != nil {
						return nil, err
					}

					// Add relationship to activity: activity_id
					// has-many or has-one
					// type: "has-many or has-one"
					// activityId: id
					// field_id:
					_, err = store.AddRelationshipIntoActivity(ctx, AddRelationshipIntoActivityParams{
						Id:             details.ActivityId,
						OrganizationId: arg.OrganizationId,

						Type:             "has_many",
						ActivityId:       arg.Activity.Id,
						FieldId:          field.Id,
						ConcernedFieldId: details.FieldId,
					})
					if err != nil {
						return nil, err
					}

				default:

				}
			}

			if fieldSplitten[len(fieldSplitten)-1] == "type" {
				if field.Type == "key" {
					var fieldRelationship *models.ActivityRelationship = nil
					for i := range arg.Activity.Relationships {
						r := arg.Activity.Relationships[i]
						if r.ConcernedFieldId == field.Id {
							fieldRelationship = &r

							break
						}
					}
					_, err := store.RemoveRelationshipFromActivity(ctx, RemoveRelationshipFromActivityParams{
						Id:             arg.Activity.Id,
						OrganizationId: arg.OrganizationId,

						Type:             fieldRelationship.Type,
						ActivityId:       fieldRelationship.ActivityId,
						FieldId:          fieldRelationship.FieldId,
						ConcernedFieldId: fieldRelationship.ConcernedFieldId,
					})

					if err != nil {
						return nil, err
					}
				}
			}
		}

		return store.UpdateSetInActivity(ctx, UpdateSetInActivityParams{
			Id:             arg.Activity.Id,
			OrganizationId: arg.OrganizationId,

			FieldsToSet: arg.FieldsToSet,
		})
	})

	if err != nil {
		return nil, err
	}

	if updatedActivity, ok := result.(*models.Activity); ok {
		return updatedActivity, err
	} else {
		return nil, err
	}
}

type UpdateRemoveFromActivityTxParams struct {
	Activity       models.Activity
	OrganizationId primitive.ObjectID

	Position uint
	Field    string
	// Value interface{}
}

func (store *MongoStorage) UpdateRemoveFromActivityTx(ctx context.Context, arg UpdateRemoveFromActivityTxParams) (*models.Activity, error) {
	result, err := store.withTx(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		field := arg.Activity.Fields[arg.Position]
		if field.Type == "key" {
			var fieldRelationship *models.ActivityRelationship = nil
			for i := range arg.Activity.Relationships {
				r := arg.Activity.Relationships[i]
				if r.ConcernedFieldId == field.Id {
					fieldRelationship = &r

					break
				}
			}
			_, err := store.RemoveRelationshipFromActivity(ctx, RemoveRelationshipFromActivityParams{
				Id:             arg.Activity.Id,
				OrganizationId: arg.OrganizationId,

				Type:             fieldRelationship.Type,
				ActivityId:       fieldRelationship.ActivityId,
				FieldId:          fieldRelationship.FieldId,
				ConcernedFieldId: fieldRelationship.ConcernedFieldId,
			})

			if err != nil {
				return nil, err
			}
		}

		return store.UpdateRemoveFromActivity(ctx, UpdateRemoveFromActivityParams{
			Id:             arg.Activity.Id,
			OrganizationId: arg.OrganizationId,

			Field:    arg.Field,
			Position: arg.Position,
		})
	})

	if err != nil {
		return nil, err
	}

	if updatedActivity, ok := result.(*models.Activity); ok {
		return updatedActivity, err
	} else {
		return nil, err
	}

}
