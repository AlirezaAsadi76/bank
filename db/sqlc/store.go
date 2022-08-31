package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db, Queries: New(db)}
}
func (s *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := New(tx)
	err = fn(q)
	if err != nil {
		erb := tx.Rollback()
		if erb != nil {
			return fmt.Errorf("tx error:%v roolback error :%v", err, erb)
		}
		return err
	}
	return tx.Commit()
}

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"fromAccount"`
	ToAccount   Account  `json:"toAccount"`
	FromEntry   Entry    `json:"fromEntry"`
	ToEntry     Entry    `json:"toEntry"`
}

func (s *Store) TransferTx(ctx context.Context, tfp TransferTxParams) (TransferTxResult, error) {

	var Result TransferTxResult

	err := s.execTx(ctx, func(queries *Queries) error {
		var errr error

		Result.Transfer, errr = queries.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: tfp.FromAccountID,
			ToAccountID:   tfp.ToAccountID,
			Amount:        tfp.Amount,
		})
		if errr != nil {
			return errr
		}
		Result.FromEntry, errr = queries.CreateEntry(ctx, CreateEntryParams{
			AccountID: tfp.FromAccountID,
			Amount:    -tfp.Amount,
		})
		if errr != nil {
			return errr
		}
		Result.ToEntry, errr = queries.CreateEntry(ctx, CreateEntryParams{
			AccountID: tfp.ToAccountID,
			Amount:    tfp.Amount,
		})
		if errr != nil {
			return errr
		}
		if tfp.FromAccountID < tfp.ToAccountID {
			Result.FromAccount, errr = queries.UpdateAccountMoney(ctx, tfp.FromAccountID, int64(-1)*tfp.Amount)
			if errr != nil {
				return errr
			}

			Result.ToAccount, errr = queries.UpdateAccountMoney(ctx, tfp.ToAccountID, tfp.Amount)
			if errr != nil {
				return errr
			}
		} else {
			Result.ToAccount, errr = queries.UpdateAccountMoney(ctx, tfp.ToAccountID, tfp.Amount)
			if errr != nil {
				return errr
			}
			Result.FromAccount, errr = queries.UpdateAccountMoney(ctx, tfp.FromAccountID, int64(-1)*tfp.Amount)
			if errr != nil {
				return errr
			}
		}
		return nil
	})
	return Result, err
}
func (q *Queries) UpdateAccountMoney(ctx context.Context, ID int64, Amount int64) (Account, error) {
	var Acc Account
	acc1, err := q.GetAccountForUpdate(ctx, ID)
	if err != nil {
		return Acc, err
	}
	Acc, err = q.UpdateAccount(ctx, UpdateAccountParams{
		ID:      acc1.ID,
		Balance: acc1.Balance + Amount,
	})
	if err != nil {
		return Acc, err
	}
	return Acc, nil
}
