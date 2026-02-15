package comdirect

// ValueUnit represents a value with its unit (e.g., currency).
type ValueUnit struct {
	Value string `json:"value"`
	Unit  string `json:"unit"`
}

// AccountType represents the type of an account.
type AccountType struct {
	Key  string `json:"key"`
	Text string `json:"text"`
}

// Account represents the account information.
type Account struct {
	AccountId        string      `json:"accountId"`
	AccountDisplayId string      `json:"accountDisplayId"`
	Currency         string      `json:"currency"`
	ClientId         string      `json:"clientId"`
	AccountType      AccountType `json:"accountType"`
	IBAN             string      `json:"iban"`
	BIC              string      `json:"bic"`
	CreditLimit      ValueUnit   `json:"creditLimit"`
}

// AccountBalance represents the full response structure for account balance.
type AccountBalance struct {
	Account                Account   `json:"account"`
	AccountId              string    `json:"accountId"`
	Balance                ValueUnit `json:"balance"`
	BalanceEUR             ValueUnit `json:"balanceEUR"`
	AvailableCashAmount    ValueUnit `json:"availableCashAmount"`
	AvailableCashAmountEUR ValueUnit `json:"availableCashAmountEUR"`
}

type Paging struct {
	Index   int `json:"index"`
	Matches int `json:"matches"`
}

type AccountBalancesResponse struct {
	Paging Paging           `json:"paging"`
	Values []AccountBalance `json:"values"`
}
