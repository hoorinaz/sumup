package account

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestTransfer(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.ExpectCommit()

	transferService := &TransferService{db: db}

	fromAccount := &Account{
		ID:      1,
		Owner:   "from_owner",
		Balance: 1000.0,
	}

	toAccount := &Account{
		ID:      2,
		Owner:   "to_owner",
		Balance: 500.0,
	}

	amount := 200.0

	t.Run("Successful Transfer", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec("UPDATE accounts SET balance = balance - \\? WHERE id = \\?").
			WithArgs(amount, fromAccount.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("UPDATE accounts SET balance = balance \\+ \\? WHERE id = \\?").
			WithArgs(amount, toAccount.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("INSERT INTO transactions \\(from_account, to_account, amount\\) VALUES \\(\\?, \\?, \\?\\)").
			WithArgs(fromAccount.ID, toAccount.ID, amount).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		from, to, err := transferService.Transfer(fromAccount, toAccount, amount)
		assert.NoError(t, err)
		assert.Equal(t, 800.0, from.Balance)
		assert.Equal(t, 700.0, to.Balance)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("Insufficient Funds", func(t *testing.T) {
		lowBalanceAccount := &Account{
			ID:      3,
			Owner:   "low_balance_owner",
			Balance: 100.0,
		}

		from, to, err := transferService.Transfer(lowBalanceAccount, toAccount, amount)
		assert.EqualError(t, err, "insufficient funds")
		assert.Nil(t, from)
		assert.Nil(t, to)
		assert.Equal(t, 100.0, lowBalanceAccount.Balance)
		assert.Equal(t, 700.0, toAccount.Balance) // Ensure toAccount balance remains unchanged
	})

	t.Run("Transaction Failure", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec("UPDATE accounts SET balance = balance - \\? WHERE id = \\?").
			WithArgs(amount, fromAccount.ID).
			WillReturnError(errors.New("update failed"))
		mock.ExpectRollback()

		from, to, err := transferService.Transfer(fromAccount, toAccount, amount)
		assert.EqualError(t, err, "update failed")
		assert.Nil(t, from)
		assert.Nil(t, to)
		assert.Equal(t, 800.0, fromAccount.Balance) // Balance should not change due to rollback
		assert.Equal(t, 700.0, toAccount.Balance)   // Balance should not change due to rollback

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}
