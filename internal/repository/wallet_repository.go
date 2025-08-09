package repository

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/L0Qqi/wallet-api/internal/model"
	"github.com/google/uuid"
)

type WalletRepository interface {
	UpdateBalanceTx(ctx context.Context, req model.WalletOperationRequest) error
	GetBalance(ctx context.Context, walletID uuid.UUID) (int64, error)
}

type walletRepository struct {
	db *sql.DB
}

func NewWalletRepository(db *sql.DB) WalletRepository {
	return &walletRepository{db: db}
}

func (r *walletRepository) UpdateBalanceTx(ctx context.Context, req model.WalletOperationRequest) error {
	log.Printf("OperationType from req: %q", req.OperationType)
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var currentBalance int64

	//читаем баланс и блокируем строку, чтобы другие транзакции ждали
	row := tx.QueryRowContext(ctx, `SELECT balance FROM wallets WHERE id = $1 FOR UPDATE`, req.WalletID)

	err = row.Scan(&currentBalance)
	if err == sql.ErrNoRows {
		// если кошелька нет — создаём
		if req.OperationType == model.Withdraw {
			return errors.New("Wallet not found")
		}
		_, err := tx.ExecContext(ctx, `INSERT INTO wallets (id, balance) VALUES ($1, $2)`, req.WalletID, req.Amount)
		if err != nil {
			return err
		}
		return tx.Commit()

	} else if err != nil {
		return err
	}

	if req.OperationType == model.Withdraw && currentBalance < req.Amount {
		return errors.New("not enouth money")
	}

	// считаем новый баланс
	var newBalance int64

	if req.OperationType == model.Deposit {
		newBalance = currentBalance + req.Amount
	} else {
		newBalance = currentBalance - req.Amount
	}

	_, err = tx.ExecContext(ctx, `UPDATE wallets SET balance = $1 WHERE id = $2`, newBalance, req.WalletID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *walletRepository) GetBalance(ctx context.Context, walletID uuid.UUID) (int64, error) {
	var balance int64
	err := r.db.QueryRowContext(ctx, `
		SELECT balance FROM wallets WHERE id = $1
	`, walletID).Scan(&balance)
	if err == sql.ErrNoRows {
		return 0, errors.New("wallet not found")
	}
	return balance, err
}
