package movie

import (
	"testing"
	"time"

	"github.com/go-kit/kit/log"
	"gorm.io/gorm"

	"KnowWhoami/movie-ticketing/internal/cache"
	"KnowWhoami/movie-ticketing/internal/testhelpers"
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

// seedScreen creates a minimal city → cinema → screen + 2 premium seats and returns the screen.
// Seats are required so that AddMovieShow's AfterCreate hook (GenerateShowSeats) doesn't fail.
func seedScreen(t *testing.T, db *gorm.DB) models.CinemaScreen {
	t.Helper()
	city := models.City{Name: "Mumbai", ZipCode: "400001"}
	if err := db.Create(&city).Error; err != nil {
		t.Fatalf("seed city: %v", err)
	}
	cinema := models.Cinema{Name: "INOX", CityID: city.ID}
	if err := db.Create(&cinema).Error; err != nil {
		t.Fatalf("seed cinema: %v", err)
	}
	screen := models.CinemaScreen{Name: "Screen A", CinemaID: cinema.ID}
	if err := db.Create(&screen).Error; err != nil {
		t.Fatalf("seed screen: %v", err)
	}
	seats := []models.CinemaSeat{
		{SeatNumber: 1, Type: models.Premium, CinemaScreenID: screen.ID},
		{SeatNumber: 2, Type: models.Premium, CinemaScreenID: screen.ID},
	}
	if err := db.Create(&seats).Error; err != nil {
		t.Fatalf("seed seats: %v", err)
	}
	return screen
}

func TestService_AddMovie(t *testing.T) {
	svc, _ := setupTest(t)

	out, err := svc.AddMovie(&AddMovieInput{
		Name:        "Interstellar",
		Description: "Space odyssey",
		Duration:    2*time.Hour + 49*time.Minute,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Movie.ID == 0 {
		t.Error("expected non-zero movie ID")
	}
	if out.Movie.Name != "Interstellar" {
		t.Errorf("got name %q, want %q", out.Movie.Name, "Interstellar")
	}
}

func TestService_AddMovieShow(t *testing.T) {
	svc, db := setupTest(t)
	screen := seedScreen(t, db)

	out, err := svc.AddMovie(&AddMovieInput{
		Name:        "The Dark Knight",
		Description: "Batman",
		Duration:    2*time.Hour + 32*time.Minute,
	})
	if err != nil {
		t.Fatalf("seed movie: %v", err)
	}
	movieID := out.Movie.ID
	base := time.Now().Truncate(time.Second)

	t.Run("Valid", func(t *testing.T) {
		show, err := svc.AddMovieShow(&AddMovieShowInput{
			MovieID:        movieID,
			CinemaScreenID: screen.ID,
			StartTime:      base.Add(1 * time.Hour),
			EndTime:        base.Add(3 * time.Hour),
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if show.Show.ID == 0 {
			t.Error("expected non-zero show ID")
		}
		if show.Show.MovieID != movieID {
			t.Errorf("MovieID = %d, want %d", show.Show.MovieID, movieID)
		}
	})

	t.Run("OverlappingShow", func(t *testing.T) {
		// First show: 5h–7h
		if _, err := svc.AddMovieShow(&AddMovieShowInput{
			MovieID:        movieID,
			CinemaScreenID: screen.ID,
			StartTime:      base.Add(5 * time.Hour),
			EndTime:        base.Add(7 * time.Hour),
		}); err != nil {
			t.Fatalf("seed first show: %v", err)
		}

		// Overlapping show: 6h–8h — BeforeCreate hook must reject this.
		_, err := svc.AddMovieShow(&AddMovieShowInput{
			MovieID:        movieID,
			CinemaScreenID: screen.ID,
			StartTime:      base.Add(6 * time.Hour),
			EndTime:        base.Add(8 * time.Hour),
		})
		if err == nil {
			t.Error("expected error for overlapping show, got nil")
		}
	})

	t.Run("AdjacentShowAllowed", func(t *testing.T) {
		// Show at 9h-11h; one starting exactly at 11h is not an overlap.
		if _, err := svc.AddMovieShow(&AddMovieShowInput{
			MovieID:        movieID,
			CinemaScreenID: screen.ID,
			StartTime:      base.Add(9 * time.Hour),
			EndTime:        base.Add(11 * time.Hour),
		}); err != nil {
			t.Fatalf("seed show: %v", err)
		}
		_, err := svc.AddMovieShow(&AddMovieShowInput{
			MovieID:        movieID,
			CinemaScreenID: screen.ID,
			StartTime:      base.Add(11 * time.Hour),
			EndTime:        base.Add(13 * time.Hour),
		})
		if err != nil {
			t.Fatalf("adjacent show should be allowed: %v", err)
		}
	})

	t.Run("DifferentScreenAllowed", func(t *testing.T) {
		// Same time window on a separate screen must never be blocked.
		screen2 := models.CinemaScreen{Name: "Screen B", CinemaID: screen.CinemaID}
		if err := db.Create(&screen2).Error; err != nil {
			t.Fatalf("seed screen2: %v", err)
		}
		s2seats := []models.CinemaSeat{
			{SeatNumber: 1, Type: models.Premium, CinemaScreenID: screen2.ID},
		}
		if err := db.Create(&s2seats).Error; err != nil {
			t.Fatalf("seed screen2 seats: %v", err)
		}
		_, err := svc.AddMovieShow(&AddMovieShowInput{
			MovieID:        movieID,
			CinemaScreenID: screen2.ID,
			StartTime:      base.Add(1 * time.Hour),
			EndTime:        base.Add(3 * time.Hour),
		})
		if err != nil {
			t.Fatalf("same time on different screen should be allowed: %v", err)
		}
	})
}

func TestService_GetMovieShow(t *testing.T) {
	svc, db := setupTest(t)
	screen := seedScreen(t, db)

	out, err := svc.AddMovie(&AddMovieInput{
		Name:        "Inception",
		Description: "Dreams within dreams",
		Duration:    2*time.Hour + 28*time.Minute,
	})
	if err != nil {
		t.Fatalf("seed movie: %v", err)
	}

	base := time.Now().Truncate(time.Second)
	added, err := svc.AddMovieShow(&AddMovieShowInput{
		MovieID:        out.Movie.ID,
		CinemaScreenID: screen.ID,
		StartTime:      base.Add(1 * time.Hour),
		EndTime:        base.Add(3 * time.Hour),
	})
	if err != nil {
		t.Fatalf("seed show: %v", err)
	}

	t.Run("Valid", func(t *testing.T) {
		got, err := svc.GetMovieShow(&GetMovieShowInput{ShowID: added.Show.ID})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.Show.ID != added.Show.ID {
			t.Errorf("ShowID = %d, want %d", got.Show.ID, added.Show.ID)
		}
		if got.Show.MovieID != out.Movie.ID {
			t.Errorf("MovieID = %d, want %d", got.Show.MovieID, out.Movie.ID)
		}
	})

	t.Run("NonExistentShow", func(t *testing.T) {
		// GORM Find (not First) returns a zero-value struct without error for missing rows.
		got, err := svc.GetMovieShow(&GetMovieShowInput{ShowID: 99999})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.Show.ID != 0 {
			t.Errorf("expected zero show ID for non-existent record, got %d", got.Show.ID)
		}
	})
}
