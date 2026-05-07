package cinema

import (
	"os"
	"syscall"
	"testing"

	"gorm.io/gorm"

	"KnowWhoami/movie-ticketing/internal/testhelpers"
)

// testDB is shared across all tests in this package — AutoMigrate runs once.
var testDB *gorm.DB

// TestMain serialises the three service packages that share ticketing_test.
// When `go test ./...` runs them in parallel each process blocks on an
// exclusive OS file lock until the previous package finishes.
func TestMain(m *testing.M) {
	f, _ := os.OpenFile("/tmp/.movie_ticketing_test.lock", os.O_CREATE|os.O_RDWR, 0644)
	defer f.Close()
	syscall.Flock(int(f.Fd()), syscall.LOCK_EX) //nolint:errcheck
	defer syscall.Flock(int(f.Fd()), syscall.LOCK_UN)

	testDB = testhelpers.SetupDB()
	testhelpers.CleanupDB(testDB)
	os.Exit(m.Run())
}
