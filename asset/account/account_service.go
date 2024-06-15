package account

import (
	"database/sql"
	"errors"
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
func (as *AccountService) CreateAccount(owner string) *Account {
	res, err := as.db.Exec("INSERT INTO accounts (owner, balance) VALUES (?, ?)", owner, 0.0)
	if err != nil {
		log.Fatalf("Failed to insert account: %v", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Fatalf("Failed to get last insert ID: %v", err)
	}

	return &Account{
		ID:    id,
		Owner: owner,
	}
}

// GetAccount retrieves an account by ID
func (as *AccountService) GetAccount(id int64) *Account {
	var acc Account
	err := as.db.QueryRow("SELECT id, owner, balance FROM accounts WHERE id = ?", id).Scan(&acc.ID, &acc.Owner, &acc.Balance)
	if err != nil {
		log.Println(err)
		return nil
	}

	return &acc
}

// Deposit adds an amount to the account balance
func (as *AccountService) Deposit(account *Account, amount float64) error {
	account.mutex.Lock()
	defer account.mutex.Unlock()

	_, err := as.db.Exec("UPDATE accounts SET balance = balance + ? WHERE id = ?", amount, account.ID)
	if err != nil {
		return err
	}
	account.Balance += amount
	return nil
}

// Withdraw subtracts an amount from the account balance
func (as *AccountService) Withdraw(account *Account, amount float64) error {
	account.mutex.Lock()
	defer account.mutex.Unlock()

	if account.Balance < amount {
		return errors.New("insufficient funds")
	}
	_, err := as.db.Exec("UPDATE accounts SET balance = balance - ? WHERE id = ?", amount, account.ID)
	if err != nil {
		return errors.New("transaction failed, try again")
	}
	account.Balance -= amount
	return nil
}
