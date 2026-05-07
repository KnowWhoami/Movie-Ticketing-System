package seed

import (
	"fmt"
	"os"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/spf13/cobra"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"KnowWhoami/movie-ticketing/pkgs/auth"
	"KnowWhoami/movie-ticketing/pkgs/models"
)

type CMD struct {
	RootCmd *cobra.Command
	Logger  log.Logger
}

var logger log.Logger

var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Seed the database with sample data",
	Run:   runSeed,
}

func Init(cmd *CMD) {
	logger = cmd.Logger
	cmd.RootCmd.AddCommand(seedCmd)
}

func dbConnect() (*gorm.DB, error) {
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "db"
	}
	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "3306"
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "ticketing"
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), dbHost, dbPort, dbName)
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}

func runSeed(_ *cobra.Command, _ []string) {
	db, err := dbConnect()
	if err != nil {
		_ = logger.Log("seed", "db_connect", "err", err)
		os.Exit(1)
	}

	steps := []struct {
		name string
		fn   func(*gorm.DB) error
	}{
		{"users", seedUsers},
		{"cities", seedCities},
		{"cinemas", seedCinemas},
		{"movies", seedMovies},
		{"shows", seedShows},
	}

	for _, step := range steps {
		if err := step.fn(db); err != nil {
			_ = logger.Log("seed", step.name, "err", err)
			os.Exit(1)
		}
		_ = logger.Log("seed", step.name, "status", "ok")
	}
	_ = logger.Log("seed", "done")
}

func seedUsers(db *gorm.DB) error {
	users := []models.User{
		{
			Name:     "Joy",
			Email:    "joy@example.com",
			Password: auth.HashPassword("password123"),
			UserType: models.UserTypeTheatreOwner,
		},
		{
			Name:     "Alice",
			Email:    "alice@example.com",
			Password: auth.HashPassword("password123"),
			UserType: models.UserTypeRegular,
		},
	}
	for i := range users {
		if err := db.Where(models.User{Email: users[i].Email}).FirstOrCreate(&users[i]).Error; err != nil {
			return err
		}
	}
	return nil
}

func seedCities(db *gorm.DB) error {
	cities := []models.City{
		{Name: "Bengaluru", ZipCode: "560068"},
		{Name: "Mumbai", ZipCode: "400001"},
		{Name: "Delhi", ZipCode: "110001"},
	}
	for i := range cities {
		if err := db.Where(models.City{Name: cities[i].Name}).FirstOrCreate(&cities[i]).Error; err != nil {
			return err
		}
	}
	return nil
}

type seatGroup struct {
	count    int
	seatType models.SeatType
}

func makeSeats(groups ...seatGroup) []models.CinemaSeat {
	var seats []models.CinemaSeat
	seatNum := 1
	for _, g := range groups {
		for i := 0; i < g.count; i++ {
			seats = append(seats, models.CinemaSeat{
				SeatNumber: seatNum,
				Type:       g.seatType,
			})
			seatNum++
		}
	}
	return seats
}

func seedCinemas(db *gorm.DB) error {
	var bengaluru, mumbai models.City
	if err := db.Where("name = ?", "Bengaluru").First(&bengaluru).Error; err != nil {
		return fmt.Errorf("city Bengaluru not found: %w", err)
	}
	if err := db.Where("name = ?", "Mumbai").First(&mumbai).Error; err != nil {
		return fmt.Errorf("city Mumbai not found: %w", err)
	}

	cinemas := []struct {
		name   string
		cityID int
	}{
		{"PVR Cinemas", bengaluru.ID},
		{"INOX Multiplex", bengaluru.ID},
		{"Cinepolis", mumbai.ID},
	}

	cinemaIDs := make(map[string]int)
	for _, c := range cinemas {
		var row models.Cinema
		if err := db.Where("name = ? AND city_id = ?", c.name, c.cityID).FirstOrCreate(&row, models.Cinema{
			Name:   c.name,
			CityID: c.cityID,
		}).Error; err != nil {
			return err
		}
		cinemaIDs[c.name] = row.ID
	}

	screens := []struct {
		cinemaName string
		screenName string
		seats      []models.CinemaSeat
	}{
		{"PVR Cinemas", "Screen 1", makeSeats(
			seatGroup{10, models.Recliner},
			seatGroup{5, models.Premium},
			seatGroup{5, models.Balcony},
		)},
		{"PVR Cinemas", "Screen 2", makeSeats(
			seatGroup{8, models.Premium},
			seatGroup{8, models.FrontRow},
		)},
		{"INOX Multiplex", "Screen 1", makeSeats(
			seatGroup{12, models.Recliner},
			seatGroup{8, models.Premium},
		)},
		{"Cinepolis", "Screen 1", makeSeats(
			seatGroup{10, models.Recliner},
			seatGroup{10, models.Balcony},
		)},
	}

	for _, s := range screens {
		cinemaID := cinemaIDs[s.cinemaName]
		var existing models.CinemaScreen
		if err := db.Where("cinema_id = ? AND name = ?", cinemaID, s.screenName).First(&existing).Error; err == nil {
			continue
		}

		screen := models.CinemaScreen{
			Name:     s.screenName,
			CinemaID: cinemaID,
		}
		if err := db.Create(&screen).Error; err != nil {
			return err
		}

		for i := range s.seats {
			s.seats[i].CinemaScreenID = screen.ID
		}
		if err := db.Create(&s.seats).Error; err != nil {
			return err
		}
	}

	return nil
}

func seedMovies(db *gorm.DB) error {
	movies := []models.Movie{
		{
			Name:        "Inception",
			Description: "A thief who steals corporate secrets through dream-sharing technology.",
			Duration:    148 * time.Minute,
		},
		{
			Name:        "The Dark Knight",
			Description: "Batman faces the Joker, a criminal mastermind who plunges Gotham into anarchy.",
			Duration:    152 * time.Minute,
		},
		{
			Name:        "Interstellar",
			Description: "A team of explorers travel through a wormhole in space to ensure humanity's survival.",
			Duration:    169 * time.Minute,
		},
		{
			Name:        "Avengers: Endgame",
			Description: "The Avengers assemble once more to reverse Thanos's actions and restore order to the universe.",
			Duration:    181 * time.Minute,
		},
	}
	for i := range movies {
		if err := db.Where("name = ?", movies[i].Name).FirstOrCreate(&movies[i]).Error; err != nil {
			return err
		}
	}
	return nil
}

func seedShows(db *gorm.DB) error {
	movieMap := make(map[string]models.Movie)
	for _, name := range []string{"Inception", "The Dark Knight", "Interstellar", "Avengers: Endgame"} {
		var m models.Movie
		if err := db.Where("name = ?", name).First(&m).Error; err != nil {
			return fmt.Errorf("movie %q not found: %w", name, err)
		}
		movieMap[name] = m
	}

	screenMap := make(map[string]models.CinemaScreen)
	for _, key := range []struct{ cinema, screen string }{
		{"PVR Cinemas", "Screen 1"},
		{"PVR Cinemas", "Screen 2"},
		{"INOX Multiplex", "Screen 1"},
		{"Cinepolis", "Screen 1"},
	} {
		var cinema models.Cinema
		if err := db.Where("name = ?", key.cinema).First(&cinema).Error; err != nil {
			return fmt.Errorf("cinema %q not found: %w", key.cinema, err)
		}
		var screen models.CinemaScreen
		if err := db.Where("cinema_id = ? AND name = ?", cinema.ID, key.screen).First(&screen).Error; err != nil {
			return fmt.Errorf("screen %q in %q not found: %w", key.screen, key.cinema, err)
		}
		screenMap[key.cinema+"/"+key.screen] = screen
	}

	base := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	shows := []struct {
		screenKey string
		movieName string
		startHour int
	}{
		{"PVR Cinemas/Screen 1", "Inception", 10},
		{"PVR Cinemas/Screen 1", "The Dark Knight", 14},
		{"PVR Cinemas/Screen 1", "Interstellar", 18},
		{"PVR Cinemas/Screen 2", "Avengers: Endgame", 11},
		{"PVR Cinemas/Screen 2", "Inception", 16},
		{"INOX Multiplex/Screen 1", "The Dark Knight", 10},
		{"INOX Multiplex/Screen 1", "Avengers: Endgame", 15},
		{"Cinepolis/Screen 1", "Interstellar", 12},
		{"Cinepolis/Screen 1", "Inception", 17},
	}

	for _, s := range shows {
		screen := screenMap[s.screenKey]
		movie := movieMap[s.movieName]
		start := base.Add(time.Duration(s.startHour) * time.Hour)
		end := start.Add(movie.Duration)

		var existing models.MovieShow
		if err := db.Where("cinema_screen_id = ? AND start_time = ?", screen.ID, start).First(&existing).Error; err == nil {
			continue
		}

		show := models.MovieShow{
			CinemaScreenID: screen.ID,
			MovieID:        movie.ID,
			StartTime:      start,
			EndTime:        end,
		}
		if err := db.Create(&show).Error; err != nil {
			return fmt.Errorf("creating show %q on %q: %w", s.movieName, s.screenKey, err)
		}
	}

	return nil
}
