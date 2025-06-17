package db

import (
	"database/sql"
	"fmt"
	"time"
)

var (
	DB *sql.DB
)

const (
	host    = "localhost"
	port    = "5432"
	user    = "parsa"
	dbName  = "swag_dev"
	sslMode = "disable"
)

func init() {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s", host, port, user, dbName, sslMode)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	// be sure to close the DB!
	DB = db
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
