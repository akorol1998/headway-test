package db

import (
	"context"
	"database/sql"
	"fmt"
	"payment-api/internal/config"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func Open(ctx context.Context, c *config.Config) (*sql.DB, error) {
	connectionString := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", c.Database.Host, c.Database.User, c.Database.Password, c.Database.DbName, c.Database.Port)
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}
	return db, err
}

// InitialSeed seed database with values
func InitialSeed(conn *sql.DB, stmnts []string, log *zap.SugaredLogger) error {
	_, checkErr := conn.Query("SELECT * FROM providers")
	if checkErr == nil {
		log.Infof("db already exists")
		return nil
	}
	for _, stm := range stmnts {
		log.Infof("executing stmnt: %v", stm)
		if _, err := conn.Exec(stm); err != nil {
			conn.Close()
			return fmt.Errorf("failed to execute query: %v, error: %v", stm, err)
		}
	}
	return nil
}

func Providers(conn *sql.DB, log *zap.SugaredLogger) {
	res, err := conn.Query("SELECT id, name FROM providers")
	if err != nil {
		log.Errorf("failed to fetch providers")
	}
	for res.Next() {
		var ID string
		var Name string
		if err := res.Scan(&ID, &Name); err != nil {
			log.Error("failed to scan columns from providers")
		}
		log.Infow("Product",
			"ID", ID,
			"Name", Name,
		)
	}
}
