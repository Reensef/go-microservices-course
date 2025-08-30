package part

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Reensef/go-microservices-course/inventory/internal/model"
	converter "github.com/Reensef/go-microservices-course/inventory/internal/repository/converter"
	repoModel "github.com/Reensef/go-microservices-course/inventory/internal/repository/model"
)

func (r *repository) Create(
	ctx context.Context,
	partInfo *model.PartInfo,
) (*model.Part, error) {
	repoPart := &repoModel.Part{
		Info:      *converter.ToRepoModelPartInfo(partInfo),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	result, err := r.collection.InsertOne(
		ctx,
		bson.M{
			"info":       repoPart.Info,
			"created_at": repoPart.CreatedAt,
			"updated_at": repoPart.UpdatedAt,
		},
	)
	if err != nil {
		return nil, err
	}

	id := result.InsertedID.(primitive.ObjectID).Hex()

	return &model.Part{
		ID:        id,
		Info:      *partInfo,
		CreatedAt: repoPart.CreatedAt,
		UpdatedAt: repoPart.UpdatedAt,
	}, nil
}
