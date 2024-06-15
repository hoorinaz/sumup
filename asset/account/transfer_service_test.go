package account

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestTransfer(t *testing.T) {
	// Create a new mock DB
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock DB: %v", err)
	}
	defer db.Close()

	// Create instances of TransferService using the mock DB
	transferService := NewTransferService(db)

	// Define the test case variables
	fromAccount := Account{ID: 1, Balance: 100.0}
	toAccount := Account{ID: 2, Balance: 50.0}
	amount := 30.0

	// Set up expectations for the first UPDATE query
	mock.ExpectExec("UPDATE accounts SET balance = balance - (.+) WHERE id = (.+)").
		WithArgs(amount, fromAccount.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Set up expectations for the second UPDATE query
	mock.ExpectExec("UPDATE accounts SET balance = balance + (.+) WHERE id = (.+)").
		WithArgs(amount, toAccount.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Set up expectations for the INSERT query
	mock.ExpectExec("INSERT INTO transactions \\(from_account, to_account, amount\\) VALUES \\(\\?, \\?, \\?\\)").
		WithArgs(fromAccount.ID, toAccount.ID, amount).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Perform the transfer
	err = transferService.Transfer(&fromAccount, &toAccount, amount)

	// Assert that there were no errors during the transfer
	assert.NoError(t, err, "Transfer should succeed without errors")

	// Assert that all mock expectations were met
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err, "All expectations should be met")
}
