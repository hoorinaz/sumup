package main

import (
	"SumUp_Task/asset/db"
	"flag"
	"fmt"
	"log"

	account "SumUp_Task/asset/account"
)

var (
	action       = flag.String("action", "", "Action to perform: create, deposit, withdraw")
	ownerName    = flag.String("owner", "", "Owner name for creating an account")
	accountOwner = flag.String("account", "", "Account owner for deposit or withdraw")
	amount       = flag.Float64("amount", 0, "Amount for deposit or withdraw")
)

func main() {
	// Debug: Print all command-line arguments
	//fmt.Println("Command-line arguments:", os.Args)
	flag.Parse()
	if *action == "" {
		log.Fatal("No action provided. Please specify -action=create|deposit|withdraw")
	}

	dsn := "mybankuser:mybankpass@tcp(mybank-db:3306)/mybank"
	db := db.InitDB(dsn)
	accountService := account.NewAccountService(db)
	transferService := account.NewTransferService(db)

	accountService := account.NewAccountServiceImp()
	transferService := account.TransferService(db)

	switch *action {
	case "create":
		fmt.Println("create")
		if *ownerName == "" {
			log.Fatal("Owner name is required for creating an account")
		}
		newAccount := accountService.CreateAccount(*ownerName)
		fmt.Printf("Account created: %+v\n", newAccount)
	case "deposit":
		fmt.Println("deposit", "owner= ", *accountOwner)
		if *accountOwner == "" || *amount <= 0 {
			log.Fatal("Valid account owner and amount are required for deposit")
		}
		acc := accountService.GetAccount(*accountOwner)
		if acc == nil {
			log.Fatalf("Account with name %d not found", *accountOwner)
		}
		accountService.Deposit(acc, *amount)
		fmt.Printf("Deposited %f to account %d. New balance: %f\n", *amount, acc.Owner, acc.Balance)
	case "withdraw":
		fmt.Println("withdraw")
		if *accountOwner == "" || *amount <= 0 {
			log.Fatal("Valid account ID and amount are required for withdrawal")
		}
		account := accountService.GetAccount(*accountOwner)
		if account == nil {
			log.Fatalf("Account with ID %d not found", *accountOwner)
		}
		err := accountService.Withdraw(account, *amount)
		if err != nil {
			log.Fatalf("Withdrawal failed: %v", err)
		}
		fmt.Printf("Withdrew %.2f from account %d. New balance: %.2f\n", *amount, account.Owner, account.Balance)
	case "transfer":
		if *fromAccount == 0 || *toAccount == 0 || *amount <= 0 {
			log.Fatal("Valid from account ID, to account ID, and amount are required for transfer")
		}
		fromAcc := accountService.GetAccount(*fromAccount)
		toAcc := accountService.GetAccount(*toAccount)
		if fromAcc == nil || toAcc == nil {
			log.Fatalf("From account ID %d or to account ID %d not found", *fromAccount, *toAccount)
		}
		err := transferService.Transfer(fromAcc, toAcc, *amount)
		if err != nil {
			log.Fatalf("Transfer failed: %v", err)
		}
		fmt.Printf("Transferred %.2f from account %d to account %d. New balances: %.2f, %.2f\n", *amount, fromAcc.ID, toAcc.ID, fromAcc.Balance, toAcc.Balance)

	default:
		log.Fatalf("Unknown action: %s", *action)
	}
	//
	//// create account
	//acc1 := accountService.CreateAccount("Hoorinaz")
	//acc2 := accountService.CreateAccount("asghar")
	//
	//// deposit money
	//acc1.Deposit(230)
	//acc2.Deposit(500)
	//
	//// transfer
	//transferService := account.NewTransferImp()
	//if err := transferService.Transfer(acc1, acc2, 30); err != nil {
	//	fmt.Printf("Transfer failed: %v\n", err)
	//} else {
	//	fmt.Printf("Transfer succeeded!\n")
	//}
	//
	//fmt.Printf("Account 1 balance: %f\n", acc1.Balance)
	//fmt.Printf("Account 2 balance: %f\n", acc2.Balance)

}
