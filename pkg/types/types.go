package types

type Money int64

type Currency string

type Category string

type Status string

const (
	TJS Currency = "TJS"
	RUB Currency = "RUB"
	USD Currency = "USD"
)

const (
	StatusOk         Status = "OK"
	StatusFail       Status = "FAIL"
	StatusInProgress Status = "INPROGRESS"
)

type PAN string

type Card struct {
	ID         int
	PAN        PAN
	Balance    Money
	Currency   Currency
	Color      string
	Name       string
	Active     bool
	MinBalance Money
}

type Payment struct {
	ID       int
	Amount   Money
	Category Category
	Status   Status
}

type PaymentSource struct {
	Type    string // "card"
	Number  string // номер вида '5058 xxxx xxxx 8888
	Balance Money  // баланс в дирамах
}