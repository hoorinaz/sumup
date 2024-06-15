package account

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateAccount(t *testing.T) {
	// Create a new mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Create the account service with the mock database
	accountService := NewAccountService(db)

	// Define the expected behavior for the mock database
	mock.ExpectExec("INSERT INTO accounts").
		WithArgs("Hoorie Nazari", 0.0).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Test case: Successful account creation
	owner := "Hoorie Nazari"
	account := accountService.CreateAccount(owner)

	// Assert that account is created correctly
	assert.NotNil(t, account, "Account should not be nil")
	assert.Equal(t, owner, account.Owner, "Owner name should match")
	assert.Equal(t, int64(1), account.ID, "Account ID should match")

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetAccount(t *testing.T) {
	// Create a new mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Create the account service with the mock database
	accountService := NewAccountService(db)

	// Define the expected behavior for the mock database
	rows := sqlmock.NewRows([]string{"id", "owner", "balance"}).
		AddRow(1, "Hoorie Nazari", 100.0)
	mock.ExpectQuery("SELECT id, owner, balance FROM accounts WHERE id = ?").
		WithArgs(1).
		WillReturnRows(rows)

	// Test case: Successful account retrieval
	id := int64(1)
	account := accountService.GetAccount(id)

	// Assert that account is retrieved correctly
	assert.NotNil(t, account, "Account should not be nil")
	assert.Equal(t, id, account.ID, "Account ID should match")
	assert.Equal(t, "Hoorie Nazari", account.Owner, "Owner name should match")
	assert.Equal(t, 100.0, account.Balance, "Balance should match")

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetAccount_NotFound(t *testing.T) {
	// Create a new mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Create the account service with the mock database
	accountService := NewAccountService(db)

	// Define the expected behavior for the mock database
	mock.ExpectQuery("SELECT id, owner, balance FROM accounts WHERE id = ?").
		WithArgs(1).
		WillReturnError(sql.ErrNoRows)

	// Test case: Account not found
	id := int64(1)
	account := accountService.GetAccount(id)

	// Assert that account is not retrieved
	assert.Nil(t, account, "Account should be nil when not found")

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDeposit(t *testing.T) {
	// Create a new mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Create the account service with the mock database
	accountService := NewAccountService(db)

	// Create a test account
	account := &Account{
		ID:      1,
		Owner:   "Navid Gasparof",
		Balance: 100.0,
	}

	// Define the expected behavior for the mock database
	mock.ExpectExec("UPDATE accounts SET balance = balance \\+ \\? WHERE id = \\?").
		WithArgs(50.0, account.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Test case: Successful deposit
	err = accountService.Deposit(account, 50.0)

	// Assert that the deposit was successful
	assert.NoError(t, err, "Deposit should not return an error")
	assert.Equal(t, 150.0, account.Balance, "Account balance should be updated")

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
func TestWithdraw(t *testing.T) {
	// Create a new mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Create the account service with the mock database
	accountService := NewAccountService(db)

	// Test case: Successful withdrawal
	account := &Account{
		ID:      1,
		Owner:   "Samira Batman",
		Balance: 100.0,
	}

	// Define the expected behavior for the mock database
	mock.ExpectExec("UPDATE accounts SET balance = balance - \\? WHERE id = \\?").
		WithArgs(50.0, account.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = accountService.Withdraw(account, 50.0)

	// Assert that the withdrawal was successful
	assert.NoError(t, err, "Withdraw should not return an error")
	assert.Equal(t, 50.0, account.Balance, "Account balance should be updated")

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	// Test case: Insufficient funds
	account.Balance = 30.0
	err = accountService.Withdraw(account, 50.0)

	// Assert that the withdrawal failed due to insufficient funds
	assert.EqualError(t, err, "insufficient funds", "Should return an insufficient funds error")
	assert.Equal(t, 30.0, account.Balance, "Account balance should not change")

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	// Test case: Transaction failure
	account.Balance = 100.0
	mock.ExpectExec("UPDATE accounts SET balance = balance - \\? WHERE id = \\?").
		WithArgs(50.0, account.ID).
		WillReturnError(errors.New("transaction failed"))

	err = accountService.Withdraw(account, 50.0)

	// Assert that the withdrawal failed due to a transaction error
	assert.EqualError(t, err, "transaction failed, try again", "Should return a transaction failed error")
	assert.Equal(t, 100.0, account.Balance, "Account balance should not change")

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestNewAccountService(t *testing.T) {
	// Create a new mock database
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Call NewAccountService with the mock database
	accountService := NewAccountService(db)

	// Assert that the returned AccountServic instance is not nil
	assert.NotNil(t, accountService, "NewAccountService should not return nil")

	// Assert that the db field is correctly set
	assert.Equal(t, db, accountService.db, "NewAccountService should set the db field correctly")
}
