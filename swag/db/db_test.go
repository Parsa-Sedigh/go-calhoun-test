package db_test

import (
	"fmt"
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

func TestCampaigns(t *testing.T) {
	/* Here, we set up the DB. We could potentially opening a DB conn, drop the db entirely, recreate a new one,
	run migrations(it will run in each testcase). But we don't do it here, it's not needed.

	We wanna run each individual test case with a clean DB.*/
	setup := dbReset
	teardown := dbReset

	//////////////////// APPROACH 1: ////////////////////////

	//t.Run("Create", func(t *testing.T) {
	//	// we can put the testcases here or in the testCreateCampaign func itself.
	//	tests := map[string]*db.Campaign{
	//		"active": &db.Campaign{
	//			StartsAt: time.Now(),
	//			EndsAt:   time.Now().Add(2 * time.Hour),
	//			Price:    1000,
	//		},
	//		"expired": &db.Campaign{
	//			StartsAt: time.Now().Add(-2 * time.Hour),
	//			EndsAt:   time.Now().Add(-1 * time.Hour),
	//			Price:    1000,
	//		},
	//	}
	//
	//	// table-driven test
	//	for name, campaign := range tests {
	//		t.Run(name, func(t *testing.T) {
	//			// NOTE: Instead of calling setup & teardown before & after the test func, we can pass those funcs so that
	//			// if we wanted to run that test with different setup & teardown, we could do that.
	//
	//			setup(t)
	//			testCreateCampaign(t, campaign)
	//			teardown(t)
	//		})
	//	}
	//})

	//////////////////// APPROACH 2: ////////////////////////
	testCreateCampaign(t, setup, teardown)
}

// At this stage, we just want to verify that any changes that we make aren't breaking CreateCampaign func.
// we made this private, so we don't allow it to run without setup & teardown.
func testCreateCampaign(t *testing.T, setup, teardown func(t *testing.T)) {
	tests := map[string]*db.Campaign{
		"active": &db.Campaign{
			StartsAt: time.Now(),
			EndsAt:   time.Now().Add(2 * time.Hour),
			Price:    1000,
		},
		"expired": &db.Campaign{
			StartsAt: time.Now().Add(-2 * time.Hour),
			EndsAt:   time.Now().Add(-1 * time.Hour),
			Price:    1000,
		},
	}

	for name, campaign := range tests {
		t.Run(name, func(t *testing.T) {
			setup(t)
			defer teardown(t)

			nBefore := count(t, "campaigns")

			start := campaign.StartsAt
			end := campaign.EndsAt
			price := campaign.Price

			// instead of using != , we should use .Equal() method because the timezone could change
			created, err := db.CreateCampaign(start, end, price)
			if err != nil {
				t.Fatalf("CreateCampaign() err = %v; want nil", err)
			}

			if created.ID <= 0 {
				t.Errorf("ID = %d; want > 0", created.ID)
			}

			if !created.StartsAt.Equal(start) {
				t.Errorf("StartsAt = %v; want %v", created.StartsAt, start)
			}

			nAfter := count(t, "campaigns")

			// if we didn't create exactly one record, fail
			if diff := nAfter - nBefore; diff != 1 {
				t.Fatalf("campaign count difference = %d; want %d", diff, 1)
			}

			got, err := db.GetCampaign(created.ID)
			if err != nil {
				t.Fatalf("GetCampaign() err = %v; want nil", err)
			}

			if got.ID <= 0 {
				t.Errorf("ID = %d; want > 0", got.ID)
			}

			if !got.StartsAt.Equal(start) {
				t.Errorf("StartsAt = %v; want %v", got.StartsAt, start)
			}
		})
	}
}

func dbReset(t *testing.T) {
	// first delete the orders, since it references other tables
	_, err := db.DB.Exec("delete from orders")
	if err != nil {
		t.Fatalf("dbReset failed: %v", err)
	}

	_, err = db.DB.Exec("delete from campaigns")
	if err != nil {
		t.Fatalf("dbReset failed: %v", err)
	}
}

func count(t *testing.T, table string) int {
	var n int

	// constructing sql queries using Sprintf is a terrible idea. But since this is a testcase and it won't get user input(the
	// tests are literally run by devs), it's ok here.
	if err := db.DB.QueryRow(fmt.Sprintf("select count(*) from campaigns", table)).Scan(&n); err != nil {
		t.Fatalf("Scan() err = %v; want nil", err)
	}

	return n
}
