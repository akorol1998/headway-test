package repository

import (
	"database/sql"
	"errors"
	"payment-api/internal/models"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Providerrepo struct {
	log  *zap.SugaredLogger
	conn *sql.DB
}

func NewProviderRepo(log *zap.SugaredLogger, conn *sql.DB) *Providerrepo {
	return &Providerrepo{log: log, conn: conn}
}

// FetchByID fetches single providers record by id
func (r *Providerrepo) FetchByID(id string) (*models.Provider, error) {
	_, checkErr := uuid.Parse(id)
	if checkErr != nil {
		return nil, ErrUuidInvalidFormat
	}

	// For the sake of simplicity we fetch all the fields, in real case scenario
	// we would parse `fields` parameter to the method and fetch only those listed there
	stmnt := "SELECT id, name, api_key, secret FROM providers WHERE id = $1"
	row := r.conn.QueryRow(stmnt, id)

	provider := models.Provider{}
	if err := row.Scan(&provider.ID, &provider.Name, &provider.ApiKey, &provider.Secret); err != nil {
		r.log.Errorw("failed to fetch provider by ID",
			"id", id,
			"error", err)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &provider, nil
}
