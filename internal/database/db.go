package database

import (
	"context"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

const defaultTimeout = 3 * time.Second

type DB struct {
	*sqlx.DB
	UserModel
	UserWorkoutBackupModel
	FeedbackModel
}

func New(dsn string) (*DB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	db, err := sqlx.ConnectContext(ctx, "mysql", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(2 * time.Hour)

	return &DB{
		DB:                     db,
		UserModel:              UserModel{db},
		UserWorkoutBackupModel: UserWorkoutBackupModel{db},
		FeedbackModel:          FeedbackModel{db},
	}, nil
}
