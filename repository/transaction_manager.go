package repository

import (
	"context"
	"database/sql"
	"log"
)

type transactionManager struct {
	dbConn *sql.DB
}

func NewTransactionManager(dbConn *sql.DB) TransactionManagerI {
	return &transactionManager{
		dbConn: dbConn,
	}
}

func (tm *transactionManager) BeginTransaction(ctx context.Context) (*sql.Tx, error) {
	tx, err := tm.dbConn.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
		return nil, err
	}
	return tx, nil
}
