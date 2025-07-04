package db

import (
	"database/sql"
	"os"
	"time"
)

var (
	dbURL string
	DB    *sql.DB
)

const DefaultURL = "postgres://postgres:postgres@127.0.0.1:5432/swag_dev?sslmode=disable"

func init() {
	dbURL = os.Getenv("PSQL_URL")
	if dbURL == "" {
		dbURL = DefaultURL
	}

	// if we haven't closed the conn, close it
	if DB != nil {
		DB.Close()
	}

	Open(dbURL)
}

// Open will open a database connection using the provided postgres URL.
// Be sure to close this using db.DB.Close func.
func Open(psqlURL string) error {
	db, err := sql.Open("postgres", psqlURL)
	if err != nil {
		return err
	}

	// be sure to close the DB!
	DB = db

	return nil
}

type Campaign struct {
	ID       int
	StartsAt time.Time
	EndsAt   time.Time
	Price    int
}

func CreateCampaign(start, end time.Time, price int) (*Campaign, error) {
	statement := `
	INSERT INTO campaigns (starts_at, ends_at, price)
	VALUES ($1, $2, $3)
	RETURNING id`

	var id int
	if err := DB.QueryRow(statement, start, end, price).Scan(&id); err != nil {
		return nil, err
	}

	return &Campaign{
		ID:       id,
		StartsAt: start,
		EndsAt:   end,
		Price:    price,
	}, nil
}

func ActiveCampaign() (*Campaign, error) {
	statement := `
	SELECT * FROM campaigns WHERE starts_at <= $1
	AND ends_at >= $1`

	row := DB.QueryRow(statement, time.Now(), time.Now())

	var camp Campaign
	if err := row.Scan(&camp.ID, &camp.StartsAt, &camp.EndsAt); err != nil {
		return nil, err
	}

	return &camp, nil
}

func GetCampaign(id int) (*Campaign, error) {
	statement := `SELECT * FROM campaigns WHERE id = $1`

	var camp Campaign
	if err := DB.QueryRow(statement, id).Scan(&camp); err != nil {
		return nil, err
	}

	return &camp, nil
}
