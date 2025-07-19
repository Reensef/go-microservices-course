package main

import (
	"context"
	"log"
	"slices"
	"sync"
	"time"

	inventory "github.com/Reensef/go-microservices-course/shared/pkg/proto/inventory/v1"
	"github.com/brianvoe/gofakeit/v7"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
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
	s.mu.Lock()
	defer s.mu.Unlock()

	filter := req.GetFilter()
	if filter == nil {
		result := make([]*inventory.Part, 0)
		for _, part := range s.parts {
			result = append(result, part)
		}
		return &inventory.ListPartsResponse{
			Parts: result,
		}, nil
	}

	result := make([]*inventory.Part, 0)

	for _, part := range s.filterPartsByUuids(s.parts, filter.GetUuids()) {
		if !s.filterPartByNames(part, filter.GetNames()) {
			continue
		} else if !s.filterPartByCategories(part, filter.GetCategories()) {
			continue
		} else if !s.filterPartByManufacturerCountries(part, filter.GetManufacturerCountries()) {
			continue
		} else if !s.filterPartByTags(part, filter.GetTags()) {
			continue
		} else {
			result = append(result, part)
		}
	}

	return &inventory.ListPartsResponse{
		Parts: result,
	}, nil
}

func (s *inventoryService) filterPartsByUuids(
	parts map[string]*inventory.Part,
	uuids []string,
) []*inventory.Part {
	result := make([]*inventory.Part, 0)

	if len(uuids) == 0 {
		for _, p := range parts {
			if p != nil {
				result = append(result, p)
			}
		}
		return result
	}

	for _, uuid := range uuids {
		if parts[uuid] != nil {
			result = append(result, parts[uuid])
		}
	}

	return result
}

func (s *inventoryService) filterPartByNames(
	part *inventory.Part,
	names []string,
) bool {
	if part == nil {
		return false
	}

	if len(names) == 0 {
		return true
	}

	return slices.Contains(names, part.Name)
}

func (s *inventoryService) filterPartByCategories(
	part *inventory.Part,
	categories []inventory.Category,
) bool {
	if part == nil {
		return false
	}

	if len(categories) == 0 {
		return true
	}

	return slices.Contains(categories, part.Category)
}

func (s *inventoryService) filterPartByManufacturerCountries(
	part *inventory.Part,
	countries []string,
) bool {
	if part == nil {
		return false
	}

	if len(countries) == 0 {
		return true
	}

	return slices.Contains(countries, part.Manufacturer.Country)
}

func (s *inventoryService) filterPartByTags(
	part *inventory.Part,
	tags []string,
) bool {
	if part == nil {
		return false
	}

	if len(tags) == 0 {
		return true
	}

	for _, reqTag := range tags {
		if slices.Contains(part.Tags, reqTag) {
			return true
		}
	}

	return false
}
