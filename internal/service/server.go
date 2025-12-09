package service

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ServerService struct {
	db *pgxpool.Pool
}

func NewServerService(db *pgxpool.Pool) *ServerService {
	return &ServerService{
		db: db,
	}
}

func (s *ServerService) Ping(ctx context.Context) error {
	if err := s.db.Ping(ctx); err != nil {
		return fmt.Errorf("数据库不可用: %w", err)
	}
	return nil
}
