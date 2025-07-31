package db_test

import (
	"github.com/Parsa-Sedigh/go-calhoun-test/db"
	"os"
	"testing"
	"time"
)

const defaultURL = "postgres://postgres:postgres@127.0.0.1:5432/swag_test?sslmode=disable"

func init() {
	testURL := os.Getenv("PSQL_URL")
	if testURL == "" {
		testURL = defaultURL
	}

	// if we haven't closed the conn, close it
	if db.DB != nil {
		db.DB.Close()
	}

	db.Open(testURL)
}

// At this stage, we just want to verify that any changes that we make aren't breaking CreateCampaign func.
func TestCreateCampaign(t *testing.T) {
	var beforeCount int
	if err := db.DB.QueryRow("select count(*) from campaigns").Scan(&beforeCount); err != nil {
		t.Fatalf("Scan() err = %v; want nil", err)
	}

	start := time.Now()
	end := time.Now().Add(1 * time.Hour)
	price := 1000

	// instead of using != , we should use .Equal() method because the timezone could change
	campaign, err := db.CreateCampaign(start, end, price)
	if err != nil {
		t.Fatalf("CreateCampaign() err = %v; want nil", err)
	}

	if campaign.ID <= 0 {
		t.Errorf("ID = %d; want > 0", campaign.ID)
	}

	if !campaign.StartsAt.Equal(start) {
		t.Errorf("StartsAt = %v; want %v", campaign.StartsAt, start)
	}

	var afterCount int
	if err := db.DB.QueryRow("select count(*) from campaigns").Scan(&afterCount); err != nil {
		t.Fatalf("Scan() err = %v; want nil", err)
	}

	// if we didn't create exactly one record, fail
	if afterCount-beforeCount != 1 {
		t.Fatalf("afterCount - beforeCount = %d; want %d", afterCount-beforeCount, 1)
	}

	got, err := db.GetCampaign(campaign.ID)
	if err != nil {
		t.Fatalf("GetCampaign() err = %v; want nil", err)
	}

	if got.ID <= 0 {
		t.Errorf("ID = %d; want > 0", got.ID)
	}

	if !got.StartsAt.Equal(start) {
		t.Errorf("StartsAt = %v; want %v", got.StartsAt, start)
	}
}
