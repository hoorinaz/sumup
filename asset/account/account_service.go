package account

// impleemnt an interface to create account and transfer money between accounts
// there is no depository to keep all accounts so it is a empty struct

type AccountService interface {
	CreateAccount(owner string) *Account
	Deposit(account *Account, amount float64)
	Withdraw(account *Account, amount float64) error
}

type AccountServiceImp struct {
	account *[]Account
}

func (as *AccountServiceImp) CreateAccount(name string) *Account {
	return &Account{
		Owner: name,
	}
}

func (as *AccountServiceImp) Deposit(account *Account, amount float64) {
	//account.Deposit(amount)
}

func (as *AccountServiceImp) Withdraw(account *Account, amount float64) error {
	//	err := account.Withdraw(amount)
	//	if err != nil {
	//		return account.Withdraw(amount)
	//	}
	//	return err
	return nil
}

func (as *AccountServiceImp) GetAccount(owner string) *Account {
	for _, account := range *as.account {
		if account.Owner == owner {
			return &account
		}
	}

	return nil
}

func NewAccountServiceImp() *AccountServiceImp {
	return &AccountServiceImp{
		account: &[]Account{},
	}
}
