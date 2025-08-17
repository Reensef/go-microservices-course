package part

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/Reensef/go-microservices-course/inventory/internal/model"
	"github.com/Reensef/go-microservices-course/inventory/internal/repository/converter"
	repoModel "github.com/Reensef/go-microservices-course/inventory/internal/repository/model"
)

func (r *repository) GetAll(ctx context.Context) ([]*model.Part, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer func() {
		err := cursor.Close(ctx)
		if err != nil {
			log.Printf("failed to close cursor: %v\n", err)
		}
	}()

	var parts []*repoModel.Part
	err = cursor.All(ctx, &parts)
	if err != nil {
		return nil, err
	}

	return converter.ToModelParts(parts), nil
}
