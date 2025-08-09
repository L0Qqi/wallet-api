package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/L0Qqi/wallet-api/internal/model"
	"github.com/L0Qqi/wallet-api/internal/service"
	"github.com/google/uuid"
)

type mockWalletRepo struct{}

func (m *mockWalletRepo) UpdateBalanceTx(ctx context.Context, req model.WalletOperationRequest) error {
	if req.Amount > 1000 {
		return errors.New("not enouth money")
	}
	return nil
}

func (m *mockWalletRepo) GetBalance(ctx context.Context, id uuid.UUID) (int64, error) {
	return 500, nil
}

func TestPostWallet(t *testing.T) {
	mockRepo := &mockWalletRepo{}
	walletService := service.NewWalletService(mockRepo)
	router := SetupRouter(walletService)

	walletID := uuid.New()

	// Тест успешного депозита
	body := model.WalletOperationRequest{
		WalletID:      walletID,
		OperationType: model.Deposit,
		Amount:        500,
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/api/v1/wallet", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.Code)
	}

	// Тест снятия с недостатком средств
	body.Amount = 1500
	body.OperationType = model.Withdraw
	jsonBody, _ = json.Marshal(body)

	req, _ = http.NewRequest("POST", "/api/v1/wallet", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", resp.Code)
	}
}

func TestGetWalletBalance(t *testing.T) {
	mockRepo := &mockWalletRepo{}
	walletService := service.NewWalletService(mockRepo)
	router := SetupRouter(walletService)

	walletID := uuid.New()

	req, _ := http.NewRequest("GET", "/api/v1/wallets/"+walletID.String(), nil)

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.Code)
	}
}
