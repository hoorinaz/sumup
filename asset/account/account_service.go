package account

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"sync"
)

// Account represents a bank account
type Account struct {
	ID      int64
	Owner   string
	Balance float64
	mutex   sync.Mutex
}

// AccountService provides account operations
type AccountService struct {
	db *sql.DB
}

// NewAccountService creates a new AccountService
func NewAccountService(db *sql.DB) *AccountService {
	return &AccountService{db: db}
}

// CreateAccount creates a new account
func (as *AccountService) CreateAccount(owner string, balance float64) (*Account, error) {
	tx, err := as.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	res, err := as.db.Exec("INSERT INTO accounts (owner, balance) VALUES (?, ?)", owner, balance)
	if err != nil {
		log.Fatalf("Failed to insert account: %v", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Fatalf("Failed to get last insert ID: %v", err)
	}

	acc := &Account{
		ID:      id,
		Owner:   owner,
		Balance: balance,
	}

	return acc, nil
}

// GetAccount retrieves an account by ID
func (as *AccountService) GetAccount(id int64) (*Account, error) {
	tx, err := as.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	var acc Account
	err = as.db.QueryRow("SELECT id, owner, balance FROM accounts WHERE id = ?", id).Scan(&acc.ID, &acc.Owner, &acc.Balance)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("failed to begin connection: %v", err)
	}

	return &acc, nil
}

// Deposit adds an amount to the account balance
func (as *AccountService) Deposit(account *Account, amount float64) error {
	tx, err := as.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	account.mutex.Lock()
	defer account.mutex.Unlock()

	_, err = as.db.Exec("UPDATE accounts SET balance = balance + ? WHERE id = ?", amount, account.ID)
	if err != nil {
		return err
	}

	account.Balance += amount

	return nil
}

// Withdraw subtracts an amount from the account balance
func (as *AccountService) Withdraw(account *Account, amount float64) error {
	tx, err := as.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	account.mutex.Lock()
	defer account.mutex.Unlock()

	if account.Balance <= amount {
		return errors.New("insufficient funds")
	}
	_, err = as.db.Exec("UPDATE accounts SET balance = balance - ? WHERE id = ?", amount, account.ID)
	if err != nil {
		return errors.New("transaction failed, try again")
	}
	account.Balance -= amount
	return nil
}
