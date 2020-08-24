package repository

import "github.com/jackc/pgx"

type Repository struct {
	pool             *pgx.ConnPool
	TransactionsRepo *TransactionsRepo
	BalanceRepo      *BalanceRepo
}

var repo Repository

const dbConnections = 20

func Init(config pgx.ConnConfig) error {
	var err error
	repo.pool, err = pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig:     config,
		MaxConnections: dbConnections,
	})
	if err != nil {
		return err
	}
	err = repo.createTables()
	if err != nil {
		return err
	}
	repo.TransactionsRepo = &TransactionsRepo{}
	repo.BalanceRepo = &BalanceRepo{}
	return nil
}

// relation style tables creation

func (repo *Repository) createTables() error {
	_, err := repo.pool.Exec(`
CREATE TABLE IF NOT EXISTS balance (
    id SERIAL NOT NULL PRIMARY KEY,
    user_id int NOT NULL UNIQUE,
    balance numeric(20, 2)  DEFAULT 0 CONSTRAINT non_negative_balance CHECK (balance >=0)
);

CREATE INDEX IF NOT EXISTS balance_user_id ON balance (user_id );

CREATE TABLE IF NOT EXISTS transactions  (
    id SERIAL NOT NULL PRIMARY KEY,
    user_id int NOT NULL REFERENCES balance(user_id) ON DELETE SET NULL,
    user_from_id int DEFAULT 0,
    operation int CONSTRAINT op_types CHECK (operation >=1 AND operation <= 3),
    sum numeric(20, 2) NOT NULL CONSTRAINT positive_sum CHECK (sum > 0),
    balance numeric(20, 2) CONSTRAINT non_negative_balance CHECK (balance >= 0),
    balance_from numeric(20, 2) CONSTRAINT non_negative_balance_from CHECK (balance >= 0),
    created TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS transactions_user_id ON transactions (user_id );

CREATE OR REPLACE FUNCTION update_balance() RETURNS TRIGGER
LANGUAGE  plpgsql
AS $add_transaction$
BEGIN
   UPDATE balance SET balance = NEW.balance WHERE user_id = NEW.user_id;
   IF NEW.user_from_id != 0 THEN
    BEGIN
        UPDATE balance SET  balance = NEW.balance_from WHERE user_id = NEW.user_from_id;
    END;
    END IF;
   RETURN NEW;
END
$add_transaction$;

DROP TRIGGER IF EXISTS UpdateBalance on transactions;

CREATE TRIGGER  UpdateBalance
    AFTER INSERT on transactions
    FOR EACH ROW
    EXECUTE PROCEDURE update_balance();
`)
	if err != nil {
		return err
	}
	return nil
}

func getPool() *pgx.ConnPool {
	return repo.pool
}

func GetBalanceRepo() BalanceRepoI {
	return repo.BalanceRepo
}

func GetTransactionsRepo() TransactionsRepoI {
	return repo.TransactionsRepo
}
