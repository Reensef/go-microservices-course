package part

import (
	"context"
	"fmt"
	"log"

	"github.com/Reensef/go-microservices-course/inventory/internal/model"
	"github.com/Reensef/go-microservices-course/inventory/internal/repository/converter"
	repoModel "github.com/Reensef/go-microservices-course/inventory/internal/repository/model"
)

func (r *repository) GetByFilter(
	ctx context.Context,
	filter *model.PartsFilter,
) ([]*model.Part, error) {
	if filter == nil {
		return r.GetAll(ctx)
	}

	mongoFilter, err := converter.ToMongoPartFilter(filter)
	if err != nil {
		return nil, err
	}

	cursor, err := r.collection.Find(ctx, mongoFilter)
	if err != nil {
		return nil, err
	}
	defer func() {
		cerr := cursor.Close(ctx)
		if cerr != nil {
			log.Printf("failed to close cursor: %v\n", cerr)
		}
	}()

	repoParts := make([]*repoModel.Part, 0)
	err = cursor.All(ctx, &repoParts)
	if err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return converter.ToModelParts(repoParts), nil
}
