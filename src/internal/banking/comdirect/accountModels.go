package comdirect

// ValueUnit represents a value with its unit (e.g., currency).
type ValueUnit struct {
	Value string `json:"value"`
	Unit  string `json:"unit"`
}

type KeyTextType struct {
	Key  string `json:"key"`
	Text string `json:"text"`
}

// Account represents the account information.
type Account struct {
	AccountId        string      `json:"accountId"`
	AccountDisplayId string      `json:"accountDisplayId"`
	Currency         string      `json:"currency"`
	ClientId         string      `json:"clientId"`
	AccountType      KeyTextType `json:"accountType"`
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

type AccountTransaction struct {
	Reference       string      `json:"reference"`
	BookingStatus   string      `json:"bookingStatus"`
	BookingDate     string      `json:"bookingDate"`
	ValutaDate      string      `json:"valueDate"`
	Amount          ValueUnit   `json:"amount"`
	RemittanceInfo  string      `json:"remittanceInfo"`
	TransactionType KeyTextType `json:"transactionType"`
}

type AccountTransactionsResponse struct {
	Paging Paging `json:"paging"`
	// Aggregated, not implemented yet
	Values []AccountTransaction `json:"values"`
}

type Depot struct {
	DepotId                    string `json:"depotId"`
	DepotDisplayId             string `json:"depotDisplayId"`
	ClientId                   string `json:"clientId"`
	DepotType                  string `json:"depotType"`
	DefaultSettlementAccountId string `json:"defaultSettlementAccountId"`
	TargetMarket               string `json:"targetMarket"`
}

type DepotsResponse struct {
	Paging Paging  `json:"paging"`
	Values []Depot `json:"values"`
}

type DepotAggregated struct {
	Depot                 Depot     `json:"depot"`
	PrevDayValue          ValueUnit `json:"prevDayValue"`
	CurrentValue          ValueUnit `json:"currentValue"`
	PurchaseValue         ValueUnit `json:"purchaseValue"`
	ProfitLossPurchaseAbs ValueUnit `json:"profitLossPurchaseAbs"`
	ProfitLossPurchaseRel string    `json:"profitLossPurchaseRel"`
	ProfitLossPrevDayAbs  ValueUnit `json:"profitLossPrevDayAbs"`
	ProfitLossPrevDayRel  string    `json:"profitLossPrevDayRel"`
}

type Venue struct {
	VenueId string `json:"venueId"`
	Name    string `json:"name"`
	Country string `json:"country"`
	Type    string `json:"type"`
}

type CurrentPrice struct {
	Price         ValueUnit `json:"price"`
	PriceDateTime FlexTime  `json:"priceDateTime"`
	Venue         Venue     `json:"venue"`
}

type DepotPosition struct {
	PositionId   string       `json:"positionId"`
	DepotId      string       `json:"depotId"`
	Wkn          string       `json:"wkn"`
	Quantity     ValueUnit    `json:"quantity"`
	CurrentPrice CurrentPrice `json:"currentPrice"`
}

type DepotPositionsResponse struct {
	Paging     Paging          `json:"paging"`
	Aggregated DepotAggregated `json:"aggregated"`
	Values     []DepotPosition `json:"values"`
}

type StaticData struct {
	InstrumentType     string `json:"instrumentType"`
	SettlementCurrency string `json:"settlementCurrency"`
}

type Instrument struct {
	Wkn        string     `json:"wkn"`
	Isin       string     `json:"isin"`
	Name       string     `json:"name"`
	StaticData StaticData `json:"staticData"`
}

type DepotTransaction struct {
	TransactionId    string     `json:"transactionId"`
	BookingStatus    string     `json:"bookingStatus"`
	BookingDate      string     `json:"bookingDate"`
	BusinessDate     string     `json:"businessDate"`
	Quantity         ValueUnit  `json:"quantity"`
	Instrument       Instrument `json:"instrument"`
	ExecutionPrice   ValueUnit  `json:"executionPrice"`
	TransactionValue ValueUnit  `json:"transactionValue"`
	TransactionType  string     `json:"transactionType"`
}

type DepotTransactionsResponse struct {
	Paging Paging             `json:"paging"`
	Values []DepotTransaction `json:"values"`
}
