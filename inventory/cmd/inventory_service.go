package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	inventory "github.com/Reensef/go-microservices-course/shared/pkg/proto/inventory/v1"
)

// Реализует gRPC сервис для хранения информации о деталях
type inventoryService struct {
	inventory.UnimplementedInventoryServiceServer

	mu    sync.RWMutex
	parts map[string]*inventory.Part
}

func NewInventoryService() *inventoryService {
	service := &inventoryService{
		parts: make(map[string]*inventory.Part),
	}

	// Генерируем данные о деталях
	for range 10 {
		part := generateRandomPart()
		service.parts[part.Uuid] = part
		log.Println("Generated part:", part.Uuid)
	}

	return service
}

func generateRandomPart() *inventory.Part {
	// Генерация случайных метаданных
	metadata := map[string]*inventory.Value{
		"color":       {Kind: &inventory.Value_StringValue{StringValue: gofakeit.Color()}},
		"is_fragile":  {Kind: &inventory.Value_BoolValue{BoolValue: gofakeit.Bool()}},
		"temperature": {Kind: &inventory.Value_DoubleValue{DoubleValue: gofakeit.Float64Range(-50, 100)}},
	}

	// Генерация случайной категории
	categories := []inventory.Category{
		inventory.Category_CATEGORY_UNSPECIFIED,
		inventory.Category_CATEGORY_ENGINE,
		inventory.Category_CATEGORY_FUEL,
		inventory.Category_CATEGORY_PORTHOLE,
		inventory.Category_CATEGORY_WING,
	}
	category := categories[gofakeit.Number(0, len(categories)-1)]

	// Генерация случайных тегов
	tags := gofakeit.RandomString([]string{"metal", "plastic", "electronics", "rare", "expensive"})

	// Создание детали
	part := &inventory.Part{
		Uuid:          gofakeit.UUID(),
		Name:          gofakeit.ProductName(),
		Description:   gofakeit.Sentence(10),
		Price:         float64(gofakeit.Price(100, 10000)),
		StockQuantity: int64(gofakeit.Number(0, 1000)),
		Category:      category,
		Dimensions: &inventory.Dimensions{
			Length: gofakeit.Float64Range(1, 100),
			Width:  gofakeit.Float64Range(1, 100),
			Height: gofakeit.Float64Range(1, 100),
			Weight: gofakeit.Float64Range(0.1, 500),
		},
		Manufacturer: &inventory.Manufacturer{
			Name:    gofakeit.Company(),
			Country: gofakeit.Country(),
			Website: gofakeit.URL(),
		},
		Tags:      []string{tags},
		Metadata:  metadata,
		CreatedAt: timestamppb.New(gofakeit.DateRange(time.Now().AddDate(-1, 0, 0), time.Now())),
		UpdatedAt: timestamppb.New(time.Now()),
	}

	return part
}

// Обработчик запроса на получение детали по UUID
func (s *inventoryService) GetPart(
	ctx context.Context,
	req *inventory.GetPartRequest,
) (*inventory.GetPartResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	uuid := req.GetUuid()

	part, ok := s.parts[uuid]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "part with UUID %s not found", req.GetUuid())
	}

	return &inventory.GetPartResponse{
		Part: part,
	}, nil
}

// Обработчик запроса на получения списка деталей с фильтрацией
func (s *inventoryService) ListParts(
	ctx context.Context,
	req *inventory.ListPartsRequest,
) (*inventory.ListPartsResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	parts := make([]*inventory.Part, 0, len(s.parts))

	filter := req.GetFilter()

	// Если нет фильтра, отдаем все
	if filter == nil {
		for _, part := range s.parts {
			parts = append(parts, part)
		}
		return &inventory.ListPartsResponse{
			Parts: parts,
		}, nil
	}

	// Если есть фильтр по UUID, отбираем по ключу
	uuids := filter.GetUuids()
	if len(uuids) > 0 {
		for _, uuid := range uuids {
			if part, exists := s.parts[uuid]; exists {
				if matchPartFilters(part, filter) {
					parts = append(parts, part)
				}
			}
		}
	} else {
		for _, part := range s.parts {
			if matchPartFilters(part, filter) {
				parts = append(parts, part)
			}
		}
	}

	return &inventory.ListPartsResponse{
		Parts: parts,
	}, nil
}
