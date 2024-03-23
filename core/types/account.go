package types

type Account struct {
	Address [20]byte // The address of the account
	Nonce   uint64   // The nonce of the account
	Balance uint64   // The balance of the account
}

func (ac *Account) CreateAccount(address [20]byte, nonce uint64, balance uint64) Account {
	return Account{address, nonce, balance}
	// propogate in the network
}

// Get account from network and save to db
func (ac *Account) AddAccount(accnt Account) {
	// get account from network
	// save to db
}

// send account over network
func (ac *Account) SendAccount(accnt Account) {
	// propogate in the network
}
