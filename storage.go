package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage interface {
	GetAccountByID(int) (*Account, error)
	GetAccounts() ([]*Account, error)
	CreateAccount(*Account) error
	UpdateAccount(*Account) error
	DeleteAcccount(int) error
}

type PostgresStore struct {
	DB *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	connStr := "user=postgres dbname=postgres password=gobank sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{
		DB: db,
	}, nil
}

func (s *PostgresStore) Init() error {
	return s.CreateAccountTable()
}

func (s *PostgresStore) CreateAccountTable() error {
	query := `CREATE TABLE IF NOT EXISTS Accounts (
		id SERIAL PRIMARY KEY,
		first_name VARCHAR(50),
		last_name VARCHAR(50),
		number SERIAL,
		balance SERIAL,
		created_at TIMESTAMP
	)`

	_, err := s.DB.Exec(query)

	return err
}

func (s *PostgresStore) GetAccountByID(id int) (*Account, error) {
	return nil, nil
}

func (s *PostgresStore) GetAccounts() ([]*Account, error) {
	rows, err := s.DB.Query(`SELECT * FROM Accounts`)
	if err != nil {
		return nil, err
	}

	accounts := []*Account{}
	for rows.Next() {
		account := new(Account)

		err := rows.Scan(
			&account.ID,
			&account.FirstName,
			&account.LastName,
			&account.Number,
			&account.Balance,
			&account.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		accounts = append(accounts, account)
	}

	return accounts, nil
}

func (s *PostgresStore) CreateAccount(account *Account) error {
	query := `INSERT INTO Accounts 
	(first_name, last_name, number, balance, created_at) 
	VALUES ($1, $2, $3, $4, $5)`

	resp, err := s.DB.Query(query, account.FirstName, account.LastName, account.Number, account.Balance, account.CreatedAt)

	if err != nil {
		return err
	}

	fmt.Printf("%v+\n", resp)

	return nil
}

func (s *PostgresStore) UpdateAccount(account *Account) error {
	return nil
}

func (s *PostgresStore) DeleteAcccount(id int) error {
	return nil
}
