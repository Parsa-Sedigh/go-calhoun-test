package db_test

import (
	"database/sql"
	"fmt"
	"github.com/Parsa-Sedigh/go-calhoun-test/db"
	"github.com/lib/pq"
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

	for name, want := range tests {
		t.Run(name, func(t *testing.T) {
			setup(t)
			defer teardown(t)

			nBefore := count(t, "campaigns")

			start := want.StartsAt
			end := want.EndsAt
			price := want.Price

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
				t.Fatalf("CreateCampaign increased campaign count by %d; want %d", diff, 1)
			}

			got, err := db.GetCampaign(created.ID)
			if err != nil {
				t.Fatalf("GetCampaign() err = %v; want nil", err)
			}

			if got.ID <= 0 {
				t.Errorf("GetCampaign() = %d; want > 0", got.ID)
			}

			if !got.StartsAt.Equal(start) {
				t.Errorf("StartsAt = %v; want %v", got.StartsAt, start)
			}
		})
	}
}

func TestActiveCampaign(t *testing.T) {
	dbReset(t)

	// Table driven tests
	// NOTE: Here each func is a setup func where each sets up an individual test and return what we expect to get.

	// each testcase returns the campaign and error it gets from a call to ActiveCampaign in the for loop
	tests := map[string]func(t *testing.T) (*db.Campaign, error){
		"just started": func(t *testing.T) (*db.Campaign, error) {
			// this testcase is less than perfect ... . Can we fix it?
			/* This is because when we call ActiveCampaign(), we want the cur time to be equal to the time that the campaign is supposed to start.
			We can't do that in this case because when we pass time.Now() to CreateCampaign(), by the time we get to ActiveCampaign(),
			some milliseconds will be passed, so the time when we call ActiveCampaign() won't be just now.*/
			want, err := db.CreateCampaign(time.Now(), time.Now().Add(1*time.Hour), 900)
			if err != nil {
				t.Fatalf("CreateCampaign() err = %v; want nil", err)
			}

			return want, nil
		},

		// this testcase fails, but we'd like it to pass. How do we fix it?
		/* This fails because of the same issue we had in the prev testcase(a little time passes as code executes). When we get
		to the next lines of code where ActiveCampaign() is called, the campaign is already expired.

		In a sticker selling store, this is not a huge deal because a millisecond for an active campaign being active or not, doesn't matter.
		But in a betting or financial matters.*/

		// name -> setup func
		"nearly ended": func(t *testing.T) (*db.Campaign, error) {
			want, err := db.CreateCampaign(time.Now().Add(-1*time.Hour), time.Now(), 900)
			if err != nil {
				t.Fatalf("CreateCampaign() err = %v; want nil", err)
			}

			return want, nil
		},
		"mid campaign": func(t *testing.T) (*db.Campaign, error) {
			// this testcase is less than perfect ... . Can we fix it?
			want, err := db.CreateCampaign(time.Now().Add(-1*time.Hour), time.Now().Add(-1*time.Hour), 900)
			if err != nil {
				t.Fatalf("CreateCampaign() err = %v; want nil", err)
			}

			return want, nil
		},
		"none": func(t *testing.T) (*db.Campaign, error) {
			return nil, sql.ErrNoRows
		},
		"expired recently": func(t *testing.T) (*db.Campaign, error) {
			_, err := db.CreateCampaign(time.Now().Add(-7*24*time.Hour), time.Now().Add(-1*time.Second), 900)
			if err != nil {
				t.Fatalf("CreateCampaign() err = %v; want nil", err)
			}

			return nil, sql.ErrNoRows
		},
		"future": func(t *testing.T) (*db.Campaign, error) {
			_, err := db.CreateCampaign(time.Now().Add(7*24*time.Hour), time.Now().Add(10*time.Hour), 900)
			if err != nil {
				t.Fatalf("CreateCampaign() err = %v; want nil", err)
			}

			return nil, sql.ErrNoRows
		},
	}

	for name, setup := range tests {
		t.Run(name, func(t *testing.T) {
			want, wantErr := setup(t)
			defer dbReset(t)

			campaign, err := db.ActiveCampaign()
			if err := campaignEq(campaign, want); err != nil {
				t.Errorf("ActiveCampaign() err = %v; want nil", err)
			}

			if err != wantErr {
				t.Fatalf("ActiveCampaign() err = %v; want %v", err, wantErr)
			}
		})
	}
}

func TestGetCampaign(t *testing.T) {
	dbReset(t)

	// Table driven tests
	// NOTE: Here each func is a setup func where each sets up an individual test and return what we expect to get.

	// each testcase returns the id of the campaign that we try to get, campaign we expect to get, err we expect to get
	tests := map[string]func(t *testing.T) (int, *db.Campaign, error){
		"missing": func(t *testing.T) (int, *db.Campaign, error) {
			return 123, nil, sql.ErrNoRows
		},
		"expired recently": func(t *testing.T) (int, *db.Campaign, error) {
			campaign, err := db.CreateCampaign(time.Now().Add(-7*24*time.Hour), time.Now().Add(-1*time.Second), 900)
			if err != nil {
				t.Fatalf("CreateCampaign() err = %v; want nil", err)
			}

			return campaign.ID, campaign, nil
		},
		"future": func(t *testing.T) (int, *db.Campaign, error) {
			campaign, err := db.CreateCampaign(time.Now().Add(7*24*time.Hour), time.Now().Add(10*time.Hour), 900)
			if err != nil {
				t.Fatalf("CreateCampaign() err = %v; want nil", err)
			}

			return campaign.ID, campaign, nil
		},
		"active": func(t *testing.T) (int, *db.Campaign, error) {
			campaign, err := db.CreateCampaign(time.Now().Add(7*24*time.Hour), time.Now().Add(10*time.Hour), 900)
			if err != nil {
				t.Fatalf("CreateCampaign() err = %v; want nil", err)
			}

			return campaign.ID, campaign, nil
		},
		"just started": func(t *testing.T) (int, *db.Campaign, error) {
			campaign, err := db.CreateCampaign(time.Now().Add(-7*24*time.Hour), time.Now().Add(10*time.Hour), 900)
			if err != nil {
				t.Fatalf("CreateCampaign() err = %v; want nil", err)
			}

			return campaign.ID, campaign, nil
		},
	}

	for name, setup := range tests {
		t.Run(name, func(t *testing.T) {
			id, want, wantErr := setup(t)
			defer dbReset(t)

			campaign, err := db.GetCampaign(id)
			if err := campaignEq(campaign, want); err != nil {
				t.Errorf("GetCampaign() err = %v; want nil", err)
			}

			if err != wantErr {
				t.Fatalf("GetCampaign() err = %v; want %v", err, wantErr)
			}
		})
	}
}

func TestCreateOrder_valid(t *testing.T) {
	dbReset(t)

	tests := map[string]func(*testing.T) db.Order{
		"valid": func(t *testing.T) db.Order {
			campaign, err := db.CreateCampaign(time.Now().Add(-1*time.Hour), time.Now().Add(1*time.Hour), 1000)
			if err != nil {
				t.Fatalf("CreateCampaign() err = %v; want nil", err)
			}

			return db.Order{
				CampaignID: campaign.ID,
				Customer:   testCustomer(),
				Address:    testAddress(),
				Payment:    testPayment(),
			}
		},
	}

	for name, setup := range tests {
		t.Run(name, func(t *testing.T) {
			want := setup(t)
			defer dbReset(t)

			nBefore := count(t, "orders")
			created := want

			if err := db.CreateOrder(&created); err != nil {
				t.Fatalf("CreateOrder() err = %v; want nil", err)
			}

			if created.ID <= 0 {
				t.Errorf("CreateOrder() ID = %d; want > 0", created.ID)
			}

			want.ID = created.ID
			if created != want {
				t.Errorf("CreateOrder() = %d; want %d", created, want)
			}

			nAfter := count(t, "orders")
			if diff := nAfter - nBefore; diff != 1 {
				t.Fatalf("CreateOrder() increased order count by %d; want %d", diff, 1)
			}

			got, err := db.GetOrderViaPayCus(want.Payment.CustomerID)
			if err != nil {
				t.Fatalf("GetOrderViaPayCus() err = %v; want nil", err)
			}

			if *got != want {
				t.Errorf("CreateOrder() = %v; want %v", *got, want)
			}
		})
	}
}

const (
	// Ideally our package would return better error here, but it returns sql or driver errors.
	//But since we aren't refactoring right now, I'm just trying to capture whatever the current backend uses as best as I can.
	pqForeignKeyCode = "23503"
)

func TestCreateOrder_invalid(t *testing.T) {
	dbReset(t)

	type checkFn func(error) error

	// verifies that the passed err is not nil and it is actually a pq error.
	checkPqError := func(code pq.ErrorCode) func(error) error {
		return func(err error) error {
			if err == nil {
				return fmt.Errorf("got nil; want *pq.Error")
			}

			if pqErr, ok := err.(*pq.Error); ok {
				if pqErr.Code != code {
					return fmt.Errorf("pq.Error.Code = %s; want %s", pqErr.Code, code)
				}
			}

			return nil
		}
	}

	tests := map[string]func(*testing.T) (db.Order, []checkFn){
		"missing campaign": func(t *testing.T) (db.Order, []checkFn) {
			return db.Order{
				Customer: testCustomer(),
				Address:  testAddress(),
				Payment:  testPayment(),
			}, []checkFn{checkPqError(pqForeignKeyCode)}
		},
	}

	for name, setup := range tests {
		t.Run(name, func(t *testing.T) {
			order, checks := setup(t)
			defer dbReset(t)

			nBefore := count(t, "orders")
			created := order
			createErr := db.CreateOrder(&created)

			for _, check := range checks {
				if err := check(createErr); err != nil {
					t.Errorf("CreateOrder() err = %v; want nil", err)
				}
			}

			nAfter := count(t, "orders")

			if diff := nAfter - nBefore; diff != 0 {
				t.Fatalf("CreateOrder() increased order count by %d; want %d", diff, 0)
			}

			got, err := db.GetOrderViaPayCus(order.Payment.CustomerID)

			// since the order created was invalid, it shouldn't get persisted, so we should get sql.ErrNoRows
			if err != sql.ErrNoRows {
				t.Fatalf("GetOrderViaPayCus() err = %v; want %v", err, sql.ErrNoRows)
			}

			if got != nil {
				t.Errorf("GetOrderViaPayCus() got = %v; want nil", got)
			}
		})
	}
}

func TestGetOrderViaPayCus(t *testing.T) {
	dbReset(t)

	// each testcase returns the id to search the campaign and the err that it wants from the call to GetCampaign()
	tests := map[string]func(*testing.T) (string, *db.Order, error){
		"missing": func(t *testing.T) (string, *db.Order, error) {
			// the order doesn't exist, so we expect getting sql.ErrNoRows
			return "fake_id", nil, sql.ErrNoRows
		},
		"expired recently": func(t *testing.T) (string, *db.Order, error) {
			campaign, err := db.CreateCampaign(time.Now().Add(-7*24*time.Hour), time.Now().Add(-1*time.Second), 900)
			if err != nil {
				t.Fatalf("CreateCampaign() err = %v; want nil", err)
			}

			order := db.Order{
				CampaignID: campaign.ID,
				Customer:   testCustomer(),
				Address:    testAddress(),
				Payment:    testPayment(),
			}

			order.Payment.CustomerID = "cus_123abc"

			if err := db.CreateOrder(&order); err != nil {
				t.Fatalf("CreateOrder() err = %v; want nil", err)
			}

			return order.Payment.ChargeID, &order, nil
		},
		"future campaign": func(t *testing.T) (string, *db.Order, error) {
			campaign, err := db.CreateCampaign(time.Now().Add(7*24*time.Hour), time.Now().Add(10*time.Hour), 900)
			if err != nil {
				t.Fatalf("CreateCampaign() err = %v; want nil", err)
			}

			order := db.Order{
				CampaignID: campaign.ID,
				Customer:   testCustomer(),
				Address:    testAddress(),
				Payment:    testPayment(),
			}

			order.Payment.CustomerID = "cus_888zzz"

			if err := db.CreateOrder(&order); err != nil {
				t.Fatalf("CreateOrder() err = %v; want nil", err)
			}

			return order.Payment.ChargeID, &order, nil
		},
		"active campaign": func(t *testing.T) (string, *db.Order, error) {
			campaign, err := db.CreateCampaign(time.Now().Add(-7*24*time.Hour), time.Now().Add(10*time.Hour), 900)
			if err != nil {
				t.Fatalf("CreateCampaign() err = %v; want nil", err)
			}

			order := db.Order{
				CampaignID: campaign.ID,
				Customer:   testCustomer(),
				Address:    testAddress(),
				Payment:    testPayment(),
			}

			order.Payment.CustomerID = "non_cus_prefix_string"

			if err := db.CreateOrder(&order); err != nil {
				t.Fatalf("CreateOrder() err = %v; want nil", err)
			}

			return order.Payment.ChargeID, &order, nil
		},
	}

	for name, setup := range tests {
		t.Run(name, func(t *testing.T) {
			id, want, wantErr := setup(t)
			defer dbReset(t)

			order, err := db.GetOrderViaPayCus(id)
			if err != wantErr {
				t.Fatalf("GetOrderViaPayCus() err = %v; want %v", err, wantErr)
			}

			if order == nil && want == nil {
				return
			}

			if order == nil || want == nil {
				t.Fatalf("GetOrderViaPayCus() order = %v; want %v", order, want)
			}

			// since these two don't have things like time.Time in them, we can just use the == operator instead of complex comparisons.
			// NOTE: We shouldn't compare the pointers themselves like order != want, because the memory addr of vals they point to is different.
			// Instead, first dereference them and compare the vals they point.
			if *order != *want {
				t.Fatalf("GetOrderViaPayCus() order = %+v; want %+v", *order, *want)
			}
		})
	}
}

func testCustomer() db.Customer {
	return db.Customer{}
}

func testAddress() db.Address {
	return db.Address{}
}

func testPayment() db.Payment {
	return db.Payment{
		Source:     "stripe",
		CustomerID: "cus_123abc",
	}
}

func campaignEq(got *db.Campaign, want *db.Campaign) error {
	// nil == nil
	if got == want {
		return nil
	}

	if got == nil {
		return fmt.Errorf("got nil; want %v", want)
	}

	if want == nil {
		return fmt.Errorf("got %v; want nil", got)
	}

	if got.ID != want.ID {
		return fmt.Errorf("got.ID = %d; want %d", got.ID, want.ID)
	}

	if !got.StartsAt.Equal(want.StartsAt) {
		return fmt.Errorf("got.StartsAt = %v; want %v", got.StartsAt, want.StartsAt)
	}

	if got.EndsAt.Equal(want.EndsAt) {
		return fmt.Errorf("got.EndsAt = %v; want %v", got.EndsAt, want.EndsAt)
	}

	return nil
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
