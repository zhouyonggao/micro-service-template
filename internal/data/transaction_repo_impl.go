package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"microServiceTemplate/internal/biz/repository"
	"microServiceTemplate/internal/data/ent"
)

type TransactionRepoImpl struct {
	repository.TransactionRepo
	db   *ent.Client
	hLog *log.Helper
}

// ContextTxKey 用来承载事务的上下文
type ContextTxKey struct{}

// NewTransaction .
func NewTransaction(d *Data, hLog *log.Helper) repository.TransactionRepo {
	return &TransactionRepoImpl{db: d.Db, hLog: hLog}
}

func (t *TransactionRepoImpl) Exec(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := t.db.Tx(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if v := recover(); v != nil {
			if err = tx.Rollback(); err != nil {
				t.hLog.Error("rollback failure from panic recover", err)
			}
			panic(v)
		}
	}()

	txCtx := context.WithValue(ctx, ContextTxKey{}, tx)
	err = fn(txCtx)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			err = errors.Wrap(err, rollbackErr.Error())
			t.hLog.Error("rollback failure from panic recover", err)
		}
		return err
	}

	err = tx.Commit()
	if err != nil {
		t.hLog.Error("commit failure", err)
		return err
	}

	return err
}
