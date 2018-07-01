package main

import (
	"golang.org/x/net/context"
	"github.com/lysu/go-saga"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/lysu/go-saga/storage/rmdb"
	"fmt"
	"github.com/giantswarm/retry-go"
	"time"
	"errors"
	"os"
)

// This example show how to initialize an Saga execution coordinator(SEC) and add Sub-transaction to it, then start a transfer transaction.
// In transfer transaction we deduce `100` from foo at first, then deposit 100 into `bar`, deduce & deduce wil both success or rollbacked.
func main() {
	fooAmount := 500
	barAmount := 300


	// 1. Define sub-transaction method, anonymous method is NOT required, Just define them as normal way.
	DeduceAccount := func(ctx context.Context, account string, amount int) error {
		// Do deduce amount from account, like: account.money - amount
		fooAmount -= amount
		return nil
	}
	CompensateDeduce := func(ctx context.Context, account string, amount int) error {
		// Compensate deduce amount to account, like: account.money + amount
		fooAmount += amount
		return nil
	}
	DepositAccount := func(ctx context.Context, account string, amount int) error {
		// Do deposit amount to account, like: account.money + amount
		barAmount += amount
		return nil
		//return errors.New("xxx")
	}
	CompensateDeposit := func(ctx context.Context, account string, amount int) error {
		// Compensate deposit amount from account, like: account.money - amount
		return retry.Do(func() error {
				barAmount -= amount
				return nil
				//return errors.New("yyy")
			},
			retry.MaxTries(3),
			retry.Sleep(time.Second * 1),
			retry.Timeout(15 * time.Second),
		)
	}

	// 2. Init SEC as global SINGLETON(this demo not..), and add Sub-transaction definition into SEC.
	saga.StorageConfig.RMDB.DBDialect = "sqlite3"
	saga.StorageConfig.RMDB.DBUrl = "test.db"

	saga.AddSubTxDef("deduce", DeduceAccount, CompensateDeduce).
		AddSubTxDef("deposit", DepositAccount, CompensateDeposit)


	// 3. Start a saga to transfer 100 from foo to bar.
	err := execSagaTx()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
	}

	// 4. done.
	fmt.Printf("fooAmount=%d, barAmount=%d\n", fooAmount, barAmount)
}

func execSagaTx() (err error) {
	defer func() {
		r := recover()
		if r != nil {
			switch r.(type) {
			case error:
				err = r.(error)
			case fmt.Stringer:
				err = errors.New((r.(fmt.Stringer)).String())
			default:
				err = errors.New(fmt.Sprintf("%v", r))
			}
		}
	}()

	from, to := "foo", "bar"
	amount := 100
	ctx := context.Background()
	var sagaID uint64 = 2

	saga.StartSaga(ctx, sagaID).
		ExecSub("deduce", from, amount).
		ExecSub("deposit", to, amount).
		EndSaga()

	return nil
}

