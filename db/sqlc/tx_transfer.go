package db

import "context"

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"_to_ntry"`
}

func (store *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult
	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		//transfer
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		//entry
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    +arg.Amount,
		})
		if err != nil {
			return err
		}

		//TODO: UPDATE ACCOUNT BALANCE
		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, err = AddBalanceIntoAccount(ctx, q, arg.FromAccountID, -arg.Amount)
			if err != nil {
				return err
			}

			result.ToAccount, err = AddBalanceIntoAccount(ctx, q, arg.ToAccountID, arg.Amount)

			if err != nil {
				return err
			}
		} else {
			result.ToAccount, err = AddBalanceIntoAccount(ctx, q, arg.ToAccountID, arg.Amount)

			if err != nil {
				return err
			}
			result.FromAccount, err = AddBalanceIntoAccount(ctx, q, arg.FromAccountID, -arg.Amount)
			if err != nil {
				return err
			}

		}

		return nil
	})
	return result, err
}

func AddBalanceIntoAccount(ctx context.Context, q *Queries, id int64, amount int64) (Account, error) {
	res, err := q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     id,
		Amount: amount,
	})

	if err != nil {
		return Account{}, err
	}
	return res, err
}
