package schema

import (
	"errors"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var TransactionTypes []string = []string{"CREDIT", "DEBIT"}

func IsValidTransactionType(transactionType string) bool {
	for i := range TransactionTypes {
		if transactionType == TransactionTypes[i] {
			return true
		}
	}

	return false
}

type Transaction struct {
	Id       primitive.ObjectID `json:"id"`
	Type     string             `json:"type"`
	Amount   float32            `json:"amount"`
	Date     time.Time          `json:"time"`
	Category string             `json:"category"`
}

type AddTransaction struct {
	Type     *string    `json:"type"`
	Amount   *float32   `json:"amount"`
	Date     *time.Time `json:"time"`
	Category *string    `json:"category"`
}

func (a *AddTransaction) Validate() (int, error) {
	var err error
	if a.Type == nil {
		err = errors.New("Missing transaction type")
	} else if !IsValidTransactionType(*a.Type) {
		err = errors.New("Invalid transaction type")
	} else if a.Amount == nil {
		err = errors.New("Missing transaction amount")
	} else if *a.Amount < 1.0 {
		err = errors.New("Invalid transaction amount")
	} else if a.Date == nil {
		err = errors.New("Missing transaction date")
	}

	if err != nil {
		return http.StatusBadRequest, err
	}
	return 0, err
}

type TransactionHistory struct {
	Transactions []Transaction `json:"transactions"`
}
