package part

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/Reensef/go-microservices-course/inventory/internal/model"
	converter "github.com/Reensef/go-microservices-course/inventory/internal/repository/converter"
	repoModel "github.com/Reensef/go-microservices-course/inventory/internal/repository/model"
)

func (r *repository) GetByID(
	ctx context.Context,
	id string,
) (*model.Part, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, model.ErrPartIdInvalidFormat
	}

	part := &repoModel.Part{}
	err = r.collection.FindOne(ctx, bson.M{"_id": objectId}).Decode(&part)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, model.ErrPartNotFound
		}

		return nil, err
	}

	return converter.ToModelPart(part), nil
}
