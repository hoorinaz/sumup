package account

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateAccount(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	accountService := &AccountService{db: db}

	owner := "test_owner"
	balance := 1000.0
	accountID := int64(1)

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO accounts").
		WithArgs(owner, balance).
		WillReturnResult(sqlmock.NewResult(accountID, 1))
	mock.ExpectCommit()

	account, err := accountService.CreateAccount(owner, balance)
	assert.NoError(t, err)
	assert.NotNil(t, account)
	assert.Equal(t, accountID, account.ID)
	assert.Equal(t, owner, account.Owner)
	assert.Equal(t, balance, account.Balance)
}

func TestDeposit(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	accountService := &AccountService{db: db}

	account := &Account{
		ID:      1,
		Owner:   "test_owner",
		Balance: 1000.0,
	}

	amount := 500.0

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE accounts SET balance = balance \\+ \\? WHERE id = \\?").
		WithArgs(amount, account.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = accountService.Deposit(account, amount)
	assert.NoError(t, err)
	assert.Equal(t, 1500.0, account.Balance)
}

func TestWithdraw(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	accountService := &AccountService{db: db}

	account := &Account{
		ID:      1,
		Owner:   "test_owner",
		Balance: 1000.0,
	}

	amount := 500.0

	t.Run("Transaction Failure", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec("UPDATE accounts SET balance = balance - \\? WHERE id = \\?").
			WithArgs(amount, account.ID).
			WillReturnError(errors.New("update failed"))
		mock.ExpectRollback()

		err := accountService.Withdraw(account, amount)
		assert.EqualError(t, err, "transaction failed, try again")
		assert.Equal(t, 1000.0, account.Balance) // The balance should not change due to rollback

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("Insufficient Funds", func(t *testing.T) {
		mock.ExpectBegin()
		lowBalanceAccount := &Account{
			ID:      2,
			Owner:   "test_owner2",
			Balance: 100.0,
		}

		err := accountService.Withdraw(lowBalanceAccount, amount)
		assert.EqualError(t, err, "insufficient funds")
		assert.Equal(t, 100.0, lowBalanceAccount.Balance)
	})

	t.Run("Successful Withdrawal", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec("UPDATE accounts SET balance = balance - \\? WHERE id = \\?").
			WithArgs(amount, account.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := accountService.Withdraw(account, amount)
		assert.NoError(t, err)
		assert.Equal(t, 500.0, account.Balance)
	})

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
