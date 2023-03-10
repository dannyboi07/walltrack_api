package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id           primitive.ObjectID `bson:"_id,omitempty"`
	Email        string             `bson:"email,omitempty"`
	PasswordHash string             `bson:"password,omitempty"`
	Balance      float32            `bson:"balance,omitempty"`
}

type Transaction struct {
	Id     primitive.ObjectID `bson:"_id,omitempty"`
	UserId primitive.ObjectID `bson:"user_id,omitempty"`
	// CREDIT|DEBIT
	Type     string    `bson:"type,omitempty"`
	Amount   float32   `bson:"amount,omitempty"`
	Date     time.Time `bson:"date,omitempty"`
	Category string    `bson:"category,omitempty"`
}

type TransactionHistory struct {
	UserId       primitive.ObjectID `bson:"user_id,omitempty"`
	Transactions []Transaction      `bson:"transactions,omitempty"`
}

// type Profile struct {
// 	Id      primitive.ObjectID `bson:"_id,omitempty"`
// 	UserId  primitive.ObjectID `bson:"user_id,omitempty"`
// 	Balance float32            `bson:"balance"`
// }
