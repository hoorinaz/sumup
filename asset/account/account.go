package account

import (
	"database/sql"
	"errors"
	"log"
)

// implemet a struct to have concret account
type Account struct {
	ID      int64
	Owner   string
	Balance float64
}

type AccountServic struct {
	db *sql.DB
}

func (as *AccountServic) CreateAccount(owner string) *Account {
	res, err := as.db.Exec("INSERT INTO accounts (owner, balance) VALUES (?, ?)", owner, 0.0)
	if err != nil {
		log.Fatalf("failed to insert account: %v", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Fatalf("failed to insert account: %v", err)
	}

	return &Account{
		ID:    id,
		Owner: owner,
	}
}

func (as *AccountServic) GetAccount(id int64) *Account {
	var acc Account
	err := as.db.QueryRow("SELECT id, owner, balance FROM accounts WHERE id = ?", id).Scan(&acc.ID, &acc.Owner, &acc.Balance)
	if err != nil {
		log.Println(err)
		return nil
	}

	return &acc
}

func (as *AccountServic) Deposit(account *Account, amount float64) {
	_, err := as.db.Exec("UPDATE accounts SET balance = balance + ? WHERE id = ?", amount, account.ID)
	if err != nil {
		log.Fatal(err)
	}
	account.Balance += amount
}

func (as *AccountServic) Withdraw(account *Account, amount float64) error {
	if account.Balance <= amount {
		return errors.New("insufficient funds")
	}
	_, err := as.db.Exec("UPDATE accounts SET balance = balance - ? WHERE id = ?", amount, account.ID)
	if err != nil {
		return errors.New("transaction failed, try again")
	}
	account.Balance -= amount
	return nil
}

func NewAccountService(db *sql.DB) *AccountServic {
	return &AccountServic{
		db: db,
	}
}
