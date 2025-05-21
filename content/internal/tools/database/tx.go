package database

import (
	"context"
	"github.com/jackc/pgx/v4"
)

func (db Database) ConfirmOperation(operationTransaction any) error {
	return operationTransaction.(pgx.Tx).Commit(context.Background())
}

func (db Database) CancelOperation(operationTransaction any) error {
	return operationTransaction.(pgx.Tx).Rollback(context.Background())
}
