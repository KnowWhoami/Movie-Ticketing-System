package booking

import (
	"testing"
	"time"

	"github.com/go-kit/kit/log"
	"gorm.io/gorm"

	"KnowWhoami/movie-ticketing/internal/cache"
	"KnowWhoami/movie-ticketing/internal/testhelpers"
	"KnowWhoami/movie-ticketing/pkgs/auth"
	"KnowWhoami/movie-ticketing/pkgs/models"
)

func setupTest(t *testing.T) (*Service, *gorm.DB) {
	t.Helper()
	testhelpers.CleanupDB(testDB)
	if sqlDB, err := testDB.DB(); err == nil {
		sqlDB.SetMaxIdleConns(0)
		sqlDB.SetMaxIdleConns(10)
	}
	t.Cleanup(func() { testhelpers.CleanupDB(testDB) })
	svc := NewService(testDB, cache.NewCache(cache.InMemoryCache), log.NewNopLogger())
	return svc, testDB
}

type seedResult struct {
	user   models.User
	show   models.MovieShow
	screen models.CinemaScreen
}

// seedFullChain creates city → cinema → screen + 2 recliner seats → movie → show → user.
// The AfterCreate hook on MovieShow generates the BookingSeats automatically.
func seedFullChain(t *testing.T, db *gorm.DB) seedResult {
	t.Helper()

	city := models.City{Name: "Delhi", ZipCode: "110001"}
	if err := db.Create(&city).Error; err != nil {
		t.Fatalf("seed city: %v", err)
	}

	cinema := models.Cinema{Name: "PVR", CityID: city.ID}
	if err := db.Create(&cinema).Error; err != nil {
		t.Fatalf("seed cinema: %v", err)
	}

	screen := models.CinemaScreen{Name: "Screen 1", CinemaID: cinema.ID}
	if err := db.Create(&screen).Error; err != nil {
		t.Fatalf("seed screen: %v", err)
	}

	seats := []models.CinemaSeat{
		{SeatNumber: 1, Type: models.Recliner, CinemaScreenID: screen.ID},
		{SeatNumber: 2, Type: models.Recliner, CinemaScreenID: screen.ID},
	}
	if err := db.Create(&seats).Error; err != nil {
		t.Fatalf("seed seats: %v", err)
	}

	mv := models.Movie{Name: "Inception", Description: "A sci-fi thriller", Duration: 2 * time.Hour}
	if err := db.Create(&mv).Error; err != nil {
		t.Fatalf("seed movie: %v", err)
	}

	now := time.Now().Truncate(time.Second)
	show := models.MovieShow{
		MovieID:        mv.ID,
		CinemaScreenID: screen.ID,
		StartTime:      now.Add(1 * time.Hour),
		EndTime:        now.Add(3 * time.Hour),
	}
	if err := db.Create(&show).Error; err != nil {
		t.Fatalf("seed show: %v", err)
	}

	user := models.User{
		Name:     "Alice",
		Email:    "alice@test.com",
		Password: auth.HashPassword("password"),
		UserType: models.UserTypeRegular,
	}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("seed user: %v", err)
	}

	return seedResult{user: user, show: show, screen: screen}
}

func TestService_BookSeats_Confirmed(t *testing.T) {
	svc, db := setupTest(t)
	seed := seedFullChain(t, db)

	out, err := svc.BookSeats(&BookSeatsInput{
		ShowID:      seed.show.ID,
		SeatNumbers: []int{1, 2},
		UserID:      seed.user.ID,
		SeatType:    models.Recliner,
	})
	if err != nil {
		t.Fatalf("BookSeats() unexpected error: %v", err)
	}
	if out.Booking.Status != models.BookingConfirmed {
		t.Errorf("status = %s, want %s", out.Booking.Status, models.BookingConfirmed)
	}
	if out.Booking.SeatCount != 2 {
		t.Errorf("SeatCount = %d, want 2", out.Booking.SeatCount)
	}
}

func TestService_BookSeats_SeatAlreadyBooked(t *testing.T) {
	svc, db := setupTest(t)
	seed := seedFullChain(t, db)

	inp := &BookSeatsInput{
		ShowID:      seed.show.ID,
		SeatNumbers: []int{1, 2},
		UserID:      seed.user.ID,
		SeatType:    models.Recliner,
	}

	// First booking succeeds.
	if _, err := svc.BookSeats(inp); err != nil {
		t.Fatalf("first BookSeats() unexpected error: %v", err)
	}

	// Second booking for the same seats: transaction commits but returns a FAILED booking.
	out, err := svc.BookSeats(inp)
	if err == nil {
		t.Fatal("expected error on second booking of already-booked seats")
	}
	if out == nil {
		t.Fatal("expected non-nil output even on booking failure")
	}
	if out.Booking.Status != models.BookingFailed {
		t.Errorf("status = %s, want %s", out.Booking.Status, models.BookingFailed)
	}
}

func TestService_ListBookings_ServedFromCacheOnSecondCall(t *testing.T) {
	svc, db := setupTest(t)
	seed := seedFullChain(t, db)

	// Create a confirmed booking so the cache has something non-trivial.
	if _, err := svc.BookSeats(&BookSeatsInput{
		ShowID:      seed.show.ID,
		SeatNumbers: []int{1},
		UserID:      seed.user.ID,
		SeatType:    models.Recliner,
	}); err != nil {
		t.Fatalf("BookSeats() unexpected error: %v", err)
	}

	// First ListBookings — hits DB and populates cache.
	first, err := svc.ListBookings()
	if err != nil {
		t.Fatalf("first ListBookings() error: %v", err)
	}
	if len(first.Bookings) == 0 {
		t.Fatal("expected at least one booking after BookSeats")
	}

	// Bypass the service and wipe the bookings table directly.
	// DELETE (not TRUNCATE) avoids a schema-version change that would break
	// any transaction still holding a connection from before this point.
	// booking_seats references bookings, so delete the child rows first.
	db.Exec("DELETE FROM booking_seats")
	db.Exec("DELETE FROM bookings")

	// Second call must return the cached result, not the now-empty DB.
	second, err := svc.ListBookings()
	if err != nil {
		t.Fatalf("second ListBookings() error: %v", err)
	}
	if len(second.Bookings) != len(first.Bookings) {
		t.Errorf("cache miss: want %d bookings from cache, got %d", len(first.Bookings), len(second.Bookings))
	}
}
