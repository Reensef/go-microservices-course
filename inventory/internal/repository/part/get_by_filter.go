package part

import (
	"context"

	"github.com/Reensef/go-microservices-course/inventory/internal/model"
	converter "github.com/Reensef/go-microservices-course/inventory/internal/repository/converter"
	filters "github.com/Reensef/go-microservices-course/inventory/internal/repository/filters"
)

func (r *repository) GetByFilter(
	ctx context.Context,
	filter *model.PartsFilter,
) []*model.Part {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Если нет фильтра, отдаем все
	if filter == nil {
		modelParts := make([]*model.Part, 0, len(r.parts))
		for _, part := range r.parts {
			modelParts = append(modelParts, converter.ToModelPart(part))
		}
		return modelParts
	}

	// Если есть фильтр по UUID, отбираем по ключу
	modelParts := make([]*model.Part, 0)
	if len(filter.Uuids) > 0 {
		for _, uuid := range filter.Uuids {
			if part, exists := r.parts[uuid]; exists {
				if filters.MatchPartFilters(part, filter) {
					modelParts = append(modelParts, converter.ToModelPart(part))
				}
			}
		}
	} else {
		for _, part := range r.parts {
			if filters.MatchPartFilters(part, filter) {
				modelParts = append(modelParts, converter.ToModelPart(part))
			}
		}
	}

	return modelParts
}
