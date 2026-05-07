package cinema

import (
	"testing"

	"github.com/go-kit/kit/log"
	"gorm.io/gorm"

	"KnowWhoami/movie-ticketing/internal/cache"
	"KnowWhoami/movie-ticketing/internal/testhelpers"
	"KnowWhoami/movie-ticketing/pkgs/models"
)

// setupTest returns a clean service and the underlying DB for seeding.
// It registers a cleanup that truncates all tables after the test.
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

func seedCity(t *testing.T, db *gorm.DB) models.City {
	t.Helper()
	city := models.City{Name: "Mumbai", ZipCode: "400001"}
	if err := db.Create(&city).Error; err != nil {
		t.Fatalf("seed city: %v", err)
	}
	return city
}

func seedCinema(t *testing.T, db *gorm.DB, cityID int) models.Cinema {
	t.Helper()
	cinema := models.Cinema{Name: "PVR", CityID: cityID}
	if err := db.Create(&cinema).Error; err != nil {
		t.Fatalf("seed cinema: %v", err)
	}
	return cinema
}

func TestService_AddCinema(t *testing.T) {
	svc, db := setupTest(t)
	city := seedCity(t, db)

	t.Run("Valid", func(t *testing.T) {
		out, err := svc.AddCinema(&AddCinemaInput{CinemaName: "INOX", CityID: city.ID})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if out.Cinema.ID == 0 {
			t.Error("expected non-zero cinema ID")
		}
		if out.Cinema.Name != "INOX" {
			t.Errorf("got name %q, want %q", out.Cinema.Name, "INOX")
		}
	})

	t.Run("NonExistentCity", func(t *testing.T) {
		_, err := svc.AddCinema(&AddCinemaInput{CinemaName: "INOX", CityID: 99999})
		if err == nil {
			t.Error("expected error for non-existent city ID")
		}
	})

	t.Run("ZeroCityID", func(t *testing.T) {
		_, err := svc.AddCinema(&AddCinemaInput{CinemaName: "INOX", CityID: 0})
		if err == nil {
			t.Error("expected error for zero city ID")
		}
	})
}

func TestService_AddCinemaScreen(t *testing.T) {
	svc, db := setupTest(t)
	city := seedCity(t, db)
	cinema := seedCinema(t, db, city.ID)

	t.Run("Valid", func(t *testing.T) {
		inp := &AddCinemaScreenInput{
			CinemaID:   cinema.ID,
			ScreenName: "Screen 1",
			Seats: []*SeatInfo{
				{SeatNumber: 1, SeatType: models.Recliner},
				{SeatNumber: 2, SeatType: models.Premium},
			},
		}
		out, err := svc.AddCinemaScreen(inp)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if out.CinemaScreen.ID == 0 {
			t.Error("expected non-zero screen ID")
		}
		if out.CinemaScreen.Name != "Screen 1" {
			t.Errorf("got screen name %q, want %q", out.CinemaScreen.Name, "Screen 1")
		}
		if len(out.CinemaScreen.CinemaSeats) != 2 {
			t.Errorf("expected 2 seats, got %d", len(out.CinemaScreen.CinemaSeats))
		}
	})

	t.Run("NilSeats", func(t *testing.T) {
		// Service does not validate — nil seats creates a screen with zero seats.
		out, err := svc.AddCinemaScreen(&AddCinemaScreenInput{
			CinemaID:   cinema.ID,
			ScreenName: "Empty Screen",
			Seats:      nil,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(out.CinemaScreen.CinemaSeats) != 0 {
			t.Errorf("expected 0 seats for nil input, got %d", len(out.CinemaScreen.CinemaSeats))
		}
	})

	t.Run("NonExistentCinema", func(t *testing.T) {
		_, err := svc.AddCinemaScreen(&AddCinemaScreenInput{
			CinemaID:   99999,
			ScreenName: "Screen X",
			Seats:      []*SeatInfo{{SeatNumber: 1, SeatType: models.Recliner}},
		})
		if err == nil {
			t.Error("expected error for non-existent cinema ID")
		}
	})

	t.Run("ZeroCinemaID", func(t *testing.T) {
		_, err := svc.AddCinemaScreen(&AddCinemaScreenInput{
			CinemaID:   0,
			ScreenName: "Screen X",
			Seats:      []*SeatInfo{{SeatNumber: 1, SeatType: models.Recliner}},
		})
		if err == nil {
			t.Error("expected error for zero cinema ID")
		}
	})

	t.Run("RollsBackScreenOnSeatFailure", func(t *testing.T) {
		var screensBefore int64
		db.Model(&models.CinemaScreen{}).Count(&screensBefore)

		// Duplicate seat entries violate the unique constraint, causing the seat
		// batch insert to fail. The screen must not be committed either.
		_, err := svc.AddCinemaScreen(&AddCinemaScreenInput{
			CinemaID:   cinema.ID,
			ScreenName: "Rollback Screen",
			Seats: []*SeatInfo{
				{SeatNumber: 1, SeatType: models.Recliner},
				{SeatNumber: 1, SeatType: models.Recliner}, // duplicate — triggers constraint violation
			},
		})
		if err == nil {
			t.Fatal("expected error for duplicate seat numbers")
		}

		var screensAfter int64
		db.Model(&models.CinemaScreen{}).Count(&screensAfter)
		if screensAfter != screensBefore {
			t.Errorf("screen was not rolled back: count before=%d after=%d", screensBefore, screensAfter)
		}
	})
}

func TestService_ListCinemas(t *testing.T) {
	t.Run("EmptyDB", func(t *testing.T) {
		svc, _ := setupTest(t)
		out, err := svc.ListCinemas()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(out.Cinemas) != 0 {
			t.Errorf("expected 0 cinemas, got %d", len(out.Cinemas))
		}
	})

	t.Run("WithData", func(t *testing.T) {
		svc, db := setupTest(t)
		city := seedCity(t, db)
		seedCinema(t, db, city.ID)

		out, err := svc.ListCinemas()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(out.Cinemas) != 1 {
			t.Errorf("expected 1 cinema, got %d", len(out.Cinemas))
		}
		if out.Cinemas[0].Name != "PVR" {
			t.Errorf("got cinema name %q, want %q", out.Cinemas[0].Name, "PVR")
		}
	})

	t.Run("ServedFromCacheOnSecondCall", func(t *testing.T) {
		svc, db := setupTest(t)
		city := seedCity(t, db)
		seedCinema(t, db, city.ID)

		// First call hits DB and populates the cache.
		first, err := svc.ListCinemas()
		if err != nil {
			t.Fatalf("unexpected error on first call: %v", err)
		}

		// Delete the cinema directly from the DB, bypassing the service
		// (which would also clear the cache). The cache still holds the result.
		db.Exec("DELETE FROM cinemas")

		// Second call must return the cached result, not the now-empty DB.
		second, err := svc.ListCinemas()
		if err != nil {
			t.Fatalf("unexpected error on second call: %v", err)
		}
		if len(second.Cinemas) != len(first.Cinemas) {
			t.Errorf("cache miss: expected %d cinemas from cache, got %d", len(first.Cinemas), len(second.Cinemas))
		}
	})
}
