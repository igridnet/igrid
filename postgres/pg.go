package postgres

import (
	"context"
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/igridnet/users/models"
	"time"
)

type (
	Config struct {
		Host               string
		Port               string
		User               string
		Password           string
		Name               string
		SSLMode            string
		MaxConnectTimeout  time.Duration
		MaxConnectAttempts int
		MaxWaitTime        time.Duration
	}
)

func Initialize(ctx context.Context, config *Config) error {
	db, err := ConnectWithConfig(ctx, config)
	if err != nil {
		return fmt.Errorf("failed to initialize: %w\n", err)
	}
	err = createSchema(db)
	if err != nil {
		return fmt.Errorf("failed to initialize: %w\n", err)
	}
	return nil
}

func (c *Config) DSN() string {
	return fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.Name, c.SSLMode)
}

func ConnectWithConfig(ctx context.Context, config *Config) (*pg.DB, error) {
	//_, _ = fmt.Fprintf(io.Stderr, "using configs: %s\n", config.DSN())
	db := pg.Connect(&pg.Options{
		Addr:            fmt.Sprintf("%s:%s", config.Host, config.Port),
		User:            config.User,
		Password:        config.Password,
		Database:        config.Name,
		DialTimeout:     time.Minute,
		MaxRetries:      5,
		MinRetryBackoff: 3,
		MaxRetryBackoff: 5,
		PoolSize:        100,
		MaxConnAge:      time.Minute,
	})

	if err := db.Ping(ctx); err != nil {
		return nil, err
	}

	return db, nil
}

func createSchema(db *pg.DB) error {
	m := []interface{}{
		(*models.Admin)(nil),
		(*models.Node)(nil),
		(*models.Region)(nil),

	}

	for _, model := range m {
		opts := &orm.CreateTableOptions{
			IfNotExists: true,
		}
		err := db.Model(model).CreateTable(opts)
		if err != nil {
			return err
		}
	}
	return nil
}
