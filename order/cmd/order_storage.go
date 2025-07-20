package main

import (
	"fmt"
	"sync"

	"github.com/google/uuid"

	orderV1 "github.com/Reensef/go-microservices-course/shared/pkg/openapi/order/v1"
)

// Представляет потокобезопасное хранилище данных о заказах
// TODO: Нужно добавить версионирование
type OrderStorage struct {
	mu     sync.RWMutex
	orders map[uuid.UUID]*orderV1.OrderDto
}

// Создает новое хранилище данных о заказах
func NewOrderStorage() *OrderStorage {
	return &OrderStorage{
		orders: make(map[uuid.UUID]*orderV1.OrderDto),
	}
}

// Возвращает указатель на копию данных о заказе по его UUID
func (s *OrderStorage) GetOrder(uuid uuid.UUID) *orderV1.OrderDto {
	s.mu.RLock()
	defer s.mu.RUnlock()

	order, ok := s.orders[uuid]
	if !ok {
		return nil
	}

	copy := *order
	return &copy
}

// Генерирует свободный UUID и добавляет новый заказ
func (s *OrderStorage) AddOrder(order *orderV1.OrderDto) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if order == nil {
		return
	}

	order.OrderUUID = s.availableUuid()

	copy := *order
	s.orders[copy.OrderUUID] = &copy
}

// Обновляет данные о заказе
func (s *OrderStorage) UpdateOrder(order *orderV1.OrderDto) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if order == nil {
		return fmt.Errorf("order is nil")
	}

	uuid := order.OrderUUID
	order.OrderUUID = uuid

	copy := *order
	s.orders[uuid] = &copy
	return nil
}

// Возвращает свободный UUID для заказа
// WARNING: не потокобезопасна
func (s *OrderStorage) availableUuid() uuid.UUID {
	for {
		uuid := uuid.New()
		if _, ok := s.orders[uuid]; !ok {
			return uuid
		}
	}
}
