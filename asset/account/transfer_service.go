package account

import (
	"database/sql"
	"fmt"
	"log"
)

type TransferService struct {
	db *sql.DB
}

func NewTransferService(db *sql.DB) *TransferService {
	return &TransferService{db: db}
}

func (s *TransferService) Transfer(fromAccount, toAccount *Account, amount float64) (*Account, *Account, error) {
	tx, err := s.db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()

	fromBalance := fromAccount.Balance - amount
	if fromBalance <= 0 {
		return nil, nil, fmt.Errorf("insufficient funds")
	}

	_, err = tx.Exec("UPDATE accounts SET balance = balance - ? WHERE id = ?", amount, fromAccount.ID)
	if err != nil {
		return nil, nil, err
	}

	_, err = tx.Exec("UPDATE accounts SET balance = balance + ? WHERE id = ?", amount, toAccount.ID)
	if err != nil {
		return nil, nil, err
	}

	_, err = tx.Exec("INSERT INTO transactions (from_account, to_account, amount) VALUES (?, ?, ?)", fromAccount.ID, toAccount.ID, amount)
	if err != nil {
		return nil, nil, err
	}

	fromAccount.Balance -= amount
	toAccount.Balance += amount
	if err = tx.Commit(); err != nil {
		return nil, nil, err
	}
	return fromAccount, toAccount, nil
}
