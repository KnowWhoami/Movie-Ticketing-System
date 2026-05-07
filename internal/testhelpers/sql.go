package testhelpers

import (
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"KnowWhoami/movie-ticketing/pkgs/models"
)

func SetupDB() *gorm.DB {
	dsn := "user:user@tcp(127.0.0.1:3306)/ticketing_test?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("error connecting to DB: ", err.Error())
		os.Exit(1)
	}
	_ = db.AutoMigrate(&models.Booking{},
		&models.BookingSeat{},
		&models.Movie{},
		&models.MovieShow{},
		&models.User{},
		&models.City{},
		&models.Cinema{},
		&models.CinemaScreen{},
		&models.CinemaSeat{},
	)
	return db
}

// CleanupDB truncates all tables in dependency order so each test starts clean.
func CleanupDB(db *gorm.DB) {
	db.Exec("SET FOREIGN_KEY_CHECKS=0")
	db.Exec("TRUNCATE TABLE booking_seats")
	db.Exec("TRUNCATE TABLE bookings")
	db.Exec("TRUNCATE TABLE movie_shows")
	db.Exec("TRUNCATE TABLE movies")
	db.Exec("TRUNCATE TABLE cinema_seats")
	db.Exec("TRUNCATE TABLE cinema_screens")
	db.Exec("TRUNCATE TABLE cinemas")
	db.Exec("TRUNCATE TABLE cities")
	db.Exec("TRUNCATE TABLE users")
	db.Exec("SET FOREIGN_KEY_CHECKS=1")
}
