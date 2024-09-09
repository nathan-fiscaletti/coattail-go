package database

import (
	"context"
	"errors"
	"os"

	"github.com/nathan-fiscaletti/coattail-go/internal/keys"
	"github.com/nathan-fiscaletti/coattail-go/pkg/coattailmodels"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	ErrDatabaseNotFound = errors.New("database not found in context")
)

type DatabaseConfig struct {
	Path string
}

func ContextWithDatabase(ctx context.Context, cfg DatabaseConfig) (context.Context, error) {
	db, err := newDatabase(cfg.Path)
	if err != nil {
		return nil, err
	}

	return context.WithValue(ctx, keys.DatabaseKey, db), nil
}

func GetDatabase(ctx context.Context) (*Database, error) {
	db, ok := ctx.Value(keys.DatabaseKey).(*Database)
	if !ok {
		return nil, ErrDatabaseNotFound
	}

	return db, nil
}

type Database struct {
	*gorm.DB
}

func newDatabase(path string) (*Database, error) {
	// Check if the database file exists, if not create it
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// create the database file
		_, err := os.Create(path)
		if err != nil {
			return nil, err
		}
	}

	// Create the connection
	ormDb, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db := &Database{ormDb}
	if err = db.migrate(); err != nil {
		return nil, err
	}

	return db, nil
}

func (db *Database) migrate() error {
	err := db.AutoMigrate(&coattailmodels.Subscription{})

	return err
}
