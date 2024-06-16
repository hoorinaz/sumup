package main

import (
	"log"
	"net/http"
	"sumup/asset/account"
	"sumup/asset/db"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var (
	accountService  *account.AccountService
	transferService *account.TransferService
)

func main() {
	db, err := db.InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	accountService = account.NewAccountService(db)
	transferService = account.NewTransferService(db)

	// Setting up routes
	r := gin.Default()
	r.POST("/create", createAccountHandler)
	r.POST("/deposit", depositHandler)
	r.POST("/transfer", transferHandler)
	r.POST("/get", getAccountHandler)

	// Starting the server
	log.Println("Starting server on :8080")
	log.Fatal(r.Run(":8080"))
}

func getAccountHandler(c *gin.Context) {
	var req struct {
		ID int64 `json:"id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	acc, err := accountService.GetAccount(req.ID)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, gin.H{"account": acc})
}

func createAccountHandler(c *gin.Context) {
	var req struct {
		Owner   string  `json:"owner" binding:"required"`
		Balance float64 `json:"balance" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	acc, err := accountService.CreateAccount(req.Owner, req.Balance)
	if err != nil {
		return
	}
	c.JSON(http.StatusOK, acc)
}

func depositHandler(c *gin.Context) {
	var req struct {
		AccountID int64   `json:"account_id" binding:"required"`
		Amount    float64 `json:"amount" binding:"required,gt=0"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	acc, err := accountService.GetAccount(req.AccountID)
	if acc == nil || err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}

	if err := accountService.Deposit(acc, req.Amount); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, acc)
}

func transferHandler(c *gin.Context) {
	var req struct {
		FromAccountID int64   `json:"from_account_id" binding:"required"`
		ToAccountID   int64   `json:"to_account_id" binding:"required"`
		Amount        float64 `json:"amount" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fromAcc, err := accountService.GetAccount(req.FromAccountID)
	if fromAcc == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "fromAccount not found"})
		return
	}

	toAcc, err := accountService.GetAccount(req.ToAccountID)
	if toAcc == nil || err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "toAccount not found"})
		return
	}
	_, _, err = transferService.Transfer(fromAcc, toAcc, req.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transfer successful, from account: " + fromAcc.Owner + " to account: " + toAcc.Owner})
}
