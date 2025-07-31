package db

import (
	"database/sql"
	"time"
)

var (
	dbURL string
	DB    *sql.DB
)

const defaultURL = "postgres://postgres:postgres@127.0.0.1:5432/swag_dev?sslmode=disable"

func init() {
	Open(defaultURL)
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

type Customer struct {
	Name  string
	Email string
}

type Address struct {
	Street1 string
	Street2 string
	City    string
	State   string
	Zip     string
	Country string

	// In case the format above fails
	Raw string
}

type Payment struct {
	Source     string
	CustomerID string
	ChargeID   string
}

type Order struct {
	ID         int
	CampaignID int
	Customer   Customer
	Address    Address
	Payment    Payment
}

func CreateOrder(order *Order) error {
	statement := `
insert into orders (
                    campaign_id,
                    cus_name, cus_email,
                    adr_street1, adr_street2, adr_city, adr_state, adr_zip, adr_country,
                    adr_raw,
                    pay_source, pay_customer_id, pay_charge_id
)
values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
returning id`

	if err := DB.QueryRow(statement,
		order.CampaignID,
		order.Customer.Name,
		order.Customer.Email,
		order.Address.Street1,
		order.Address.Street2,
		order.Address.City,
		order.Address.State,
		order.Address.Zip,
		order.Address.Country,
		order.Address.Raw,
		order.Payment.Source,
		order.Payment.CustomerID,
		order.Payment.ChargeID,
	).Scan(&order.ID); err != nil {
		return err
	}

	return nil
}

func GetOrderViaPayCus(payCustomerID string) (*Order, error) {
	statement := `
	select * from orders
    where pay_customer_id = $1`

	row := DB.QueryRow(statement, payCustomerID)

	var order Order
	if err := row.Scan(
		&order.CampaignID,
		&order.Customer.Name,
		&order.Customer.Email,
		&order.Address.Street1,
		&order.Address.Street2,
		&order.Address.City,
		&order.Address.State,
		&order.Address.Zip,
		&order.Address.Country,
		&order.Address.Raw,
		&order.Payment.Source,
		&order.Payment.CustomerID,
		&order.Payment.ChargeID); err != nil {
		return nil, err
	}

	return &order, nil
}
