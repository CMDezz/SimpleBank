package db

import "context"

type CreateUserTransferTxParams struct {
	CreateUserParams
	AfterCreate func(user User) error
}

type CreateUserTransferTxResult struct {
	User User
}

func (store *SQLStore) CreateUserTransferTx(ctx context.Context, arg CreateUserTransferTxParams) (CreateUserTransferTxResult, error) {
	var result CreateUserTransferTxResult
	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		result.User, err = q.CreateUser(ctx, arg.CreateUserParams)
		if err != nil {
			return err
		}
		return arg.AfterCreate(result.User)
	})
	return result, err
}
