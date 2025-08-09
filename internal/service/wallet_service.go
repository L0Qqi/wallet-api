package service

import (
	"context"

	"github.com/L0Qqi/wallet-api/internal/model"
	"github.com/L0Qqi/wallet-api/internal/repository"
	"github.com/google/uuid"
)

type WalletService struct {
	repo repository.WalletRepository
}

func NewWalletService(repo repository.WalletRepository) *WalletService {
	return &WalletService{repo: repo}
}

func (s *WalletService) Operate(ctx context.Context, req model.WalletOperationRequest) error {
	return s.repo.UpdateBalanceTx(ctx, req)
}

func (s *WalletService) GetBalance(ctx context.Context, id uuid.UUID) (int64, error) {
	return s.repo.GetBalance(ctx, id)
}
