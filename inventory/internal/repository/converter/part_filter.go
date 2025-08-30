package converter

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	model "github.com/Reensef/go-microservices-course/inventory/internal/model"
)

func ToMongoPartFilter(filter *model.PartsFilter) (*bson.M, error) {
	mongoFilter := bson.M{}

	objectIds := make([]primitive.ObjectID, 0, len(filter.IDs))
	for _, id := range filter.IDs {
		objectId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, model.ErrPartIdInvalidFormat
		}
		objectIds = append(objectIds, objectId)
	}

	if len(filter.IDs) > 0 {
		mongoFilter["_id"] = bson.M{"$in": objectIds}
	}
	if len(filter.Names) > 0 {
		mongoFilter["info.name"] = bson.M{"$in": filter.Names}
	}
	if len(filter.Categories) > 0 {
		mongoFilter["info.category"] = bson.M{"$in": filter.Categories}
	}
	if len(filter.ManufacturerCountries) > 0 {
		mongoFilter["info.manufacturer.country"] = bson.M{"$in": filter.ManufacturerCountries}
	}
	if len(filter.Tags) > 0 {
		mongoFilter["info.tags"] = bson.M{"$in": filter.Tags}
	}

	return &mongoFilter, nil
}
