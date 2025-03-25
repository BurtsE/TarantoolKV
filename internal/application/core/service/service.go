package service

import (
	"TarantoolKV/internal/application/core/domain"
	"context"
)

type Storage interface {
	Get(ctx context.Context, key string) (domain.Entity, error)
	Create(ctx context.Context, entity domain.Entity) error
	Update(ctx context.Context, entity domain.Entity) error
	Delete(ctx context.Context, key string) error
}

type Service struct {
	storage Storage
}

func NewService(storage Storage) *Service {
	return &Service{storage: storage}
}

func (s *Service) Create(ctx context.Context, entity domain.Entity) error {
	return s.storage.Create(ctx, entity)
}
func (s *Service) Update(ctx context.Context, entity domain.Entity) error {
	return s.storage.Update(ctx, entity)
}
func (s *Service) Delete(ctx context.Context, key string) error {
	return s.storage.Delete(ctx, key)
}
func (s *Service) Get(ctx context.Context, key string) (domain.Entity, error) {
	return s.storage.Get(ctx, key)
}
