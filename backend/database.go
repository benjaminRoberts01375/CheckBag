package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strconv"

	Coms "github.com/benjaminRoberts01375/Go-Communicate"
	"github.com/benjaminRoberts01375/Web-Tech-Stack/models"
	_ "github.com/lib/pq"
)

type DatabaseSpec interface {
	// Connection Management
	Setup()
	Close()

	// Basic operations
	DatabaseGetSetSpec

	// Atomic operations
	BeginTx(ctx context.Context) (TxSpec, error)
	CommitTx(ctx context.Context) error
	RollbackTx(ctx context.Context) error
}

type DatabaseGetSetSpec interface {
	GetRow(ctx context.Context, query string, args ...any) RowSpec
	GetRows(ctx context.Context, query string, args ...any) (RowsSpec, error)
	SetRows(ctx context.Context, query string, args ...any) error
}

type TxSpec interface {
	Commit() error
	Rollback() error
	DatabaseGetSetSpec
}

type RowSpec interface {
	Scan(dest ...any) error
}

type RowsSpec interface {
	Next() bool
	Scan(dest ...any) error
	Close() error
}

type DBLayer struct { // Implements DatabaseSpec
	DB              *sql.DB
	DBPort          int    `json:"db_port"`
	DBName          string `json:"db_name"`
	DBUser          string `json:"db_user"`
	DBPassword      string `json:"db_password"`
	DBContainerName string `json:"db_container_name"`
}

type DBClient struct {
	raw DBLayer
}

type DBTransaction struct {
	Tx *sql.Tx
}

// Underlying DB operations

func (db *DBLayer) Setup() error {
	// DB Configuration
	DBPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil || DBPort <= 0 {
		panic("Failed to parse DB_PORT: " + err.Error())
	}
	DBContainerName := os.Getenv("DB_CONTAINER_NAME")
	if DBContainerName == "" {
		panic("No database container name specified")
	}

	Coms.ReadExternalConfig("db.json", db)
	db.DBPort = DBPort
	db.DBContainerName = DBContainerName
	if db.DBName == "" {
		panic("No database name specified")
	} else if db.DBUser == "" {
		panic("No database user specified")
	} else if db.DBPassword == "" {
		panic("No database password specified")
	}

	// Connect to DB
	url := fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s?sslmode=disable",
		db.DBUser, db.DBPassword, db.DBContainerName, db.DBPort, db.DBName,
	)
	sqlDB, err := sql.Open("postgres", url)
	if err != nil {
		panic("Could not connect to database: " + err.Error())
	}
	err = sqlDB.Ping()
	if err != nil {
		panic("Could not ping database: " + err.Error())
	}
	db.DB = sqlDB
	return nil
}

func (db *DBLayer) Close() error {
	return db.DB.Close()
}

func (db *DBLayer) GetRow(ctx context.Context, query string, args ...any) RowSpec {
	return db.DB.QueryRowContext(ctx, query, args...)
}

func (db *DBLayer) GetRows(ctx context.Context, query string, args ...any) (RowsSpec, error) {
	return db.DB.QueryContext(ctx, query, args...)
}

func (db *DBLayer) ExecRows(ctx context.Context, query string, args ...any) error {
	_, err := db.DB.ExecContext(ctx, query, args...)
	return err
}

func (db *DBLayer) BeginTx(ctx context.Context) (TxSpec, error) {
	transaction, err := db.DB.BeginTx(ctx, nil)
	return &DBTransaction{Tx: transaction}, err
}

// Underlying DB transaction operations

func (tx *DBTransaction) GetRow(ctx context.Context, query string, args ...any) RowSpec {
	return tx.Tx.QueryRowContext(ctx, query, args...)
}

func (tx *DBTransaction) GetRows(ctx context.Context, query string, args ...any) (RowsSpec, error) {
	return tx.Tx.QueryContext(ctx, query, args...)
}

func (tx *DBTransaction) SetRows(ctx context.Context, query string, args ...any) error {
	_, err := tx.Tx.ExecContext(ctx, query, args...)
	return err
}

func (tx *DBTransaction) Commit() error {
	return tx.Tx.Commit()
}

func (tx *DBTransaction) Rollback() error {
	return tx.Tx.Rollback()
}

// Higher-level DBClient operations
func (db *DBClient) SetNewUser(user models.UserCreate) error {
	statement := `
	INSERT INTO users (email, password, first_name, last_name)
	VALUES ($1, $2, $3, $4)
	`
	_, err := db.raw.DB.Exec(statement, user.Email, string(user.Password), user.FirstName, user.LastName)
	return err
}

func (db *DBClient) SetNewUserConfirmed(email string) error {
	statement := `UPDATE users SET confirmed = TRUE WHERE email = $1;`
	_, err := db.raw.DB.Exec(statement, email)
	return err
}

func (db *DBClient) SetUserHasLoggedIn(userID string) error {
	statement := `UPDATE users SET last_login = NOW() WHERE id = $1;`
	_, err := db.raw.DB.Exec(statement, userID)
	return err
}

func (db *DBClient) GetUserPasswordAndID(email string) ([]byte, string, error) {
	statement := `SELECT password, id FROM users WHERE email = $1 AND confirmed = TRUE;`
	var password []byte
	var userID string
	err := db.raw.DB.QueryRow(statement, email).Scan(&password, &userID)
	return password, userID, err
}

func (db *DBClient) SetUserPassword(userID string, newPassword []byte) error {
	statement := `UPDATE users SET password = $1 WHERE id = $2;`
	_, err := db.raw.DB.Exec(statement, newPassword, userID)
	return err
}

func (db *DBClient) SetUserEmail(originalEmail string, newEmail string) error {
	statement := `UPDATE users SET email = $1 WHERE email = $2;`
	_, err := db.raw.DB.Exec(statement, newEmail, originalEmail)
	return err
}

func (db *DBClient) SetUserForgotPassword(email string, newPassword []byte) error {
	statement := "UPDATE users SET password = $1 WHERE email = $2"
	_, err := db.raw.DB.Exec(statement, newPassword, email)
	return err
}

func (db *DBClient) GetUserData(userID string) (models.User, error) {
	statement := `SELECT first_name, last_name, email FROM users WHERE id = $1;`
	var user models.User
	err := db.raw.DB.QueryRow(statement, userID).Scan(&user.FirstName, &user.LastName, &user.Email)
	return user, err
}
