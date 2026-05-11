package cinema

import (
	"fmt"

	"github.com/go-kit/kit/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"KnowWhoami/movie-ticketing/internal/cache"
	"KnowWhoami/movie-ticketing/pkgs/models"
)

type Service struct {
	db     *gorm.DB
	cache  cache.Cache
	logger log.Logger
}

func NewService(db *gorm.DB, cache cache.Cache, logger log.Logger) *Service {
	return &Service{
		db:     db,
		cache:  cache,
		logger: logger,
	}
}

func (s *Service) AddCinema(inp *AddCinemaInput) (*AddCinemaOutput, error) {
	var existing models.Cinema
	if err := s.db.Where("name = ? AND city_id = ?", inp.CinemaName, inp.CityID).First(&existing).Error; err == nil {
		return nil, fmt.Errorf("cinema %q already exists in this city", inp.CinemaName)
	}

	c := models.Cinema{
		Name:          inp.CinemaName,
		CityID:        inp.CityID,
		CinemaOwnerID: inp.CinemaOwnerID,
	}
	if err := s.db.Create(&c).Error; err != nil {
		_ = s.logger.Log("method", "AddCinema", "err", err)
		return nil, err
	}
	if err := s.db.Preload("City").Preload("CinemaOwner").First(&c, c.ID).Error; err != nil {
		return nil, err
	}
	s.cache.Delete("ListCinemasOutput")
	_ = s.logger.Log("method", "AddCinema", "cinema_id", c.ID)
	return &AddCinemaOutput{Cinema: c}, nil
}

func (s *Service) ListCinemas(inp *ListCinemasInput) (*ListCinemasOutput, error) {
	query := s.db.Preload("CinemaScreens.CinemaSeats").Preload(clause.Associations)

	if inp.Name != "" {
		query = query.Where("name LIKE ?", "%"+inp.Name+"%")
	}
	if inp.CityID != 0 {
		query = query.Where("city_id = ?", inp.CityID)
	}

	var cinemas []models.Cinema
	if err := query.Find(&cinemas).Error; err != nil {
		_ = s.logger.Log("method", "ListCinemas", "err", err)
		return nil, err
	}
	return &ListCinemasOutput{Cinemas: cinemas}, nil
}

func (s *Service) GetCinemaByID(id int) (*GetCinemaByIDOutput, error) {
	var c models.Cinema
	err := s.db.Preload("CinemaScreens.CinemaSeats").Preload(clause.Associations).First(&c, id).Error
	if err != nil {
		_ = s.logger.Log("method", "GetCinemaByID", "err", err)
		return nil, fmt.Errorf("cinema not found")
	}
	return &GetCinemaByIDOutput{Cinema: c}, nil
}

func (s *Service) GetMyCinemas(cinemaOwnerID int) (*ListCinemasOutput, error) {
	var cinemas []models.Cinema
	err := s.db.Preload("CinemaScreens.CinemaSeats").Preload(clause.Associations).
		Where("cinema_owner_id = ?", cinemaOwnerID).Find(&cinemas).Error
	if err != nil {
		_ = s.logger.Log("method", "GetMyCinemas", "err", err)
		return nil, err
	}
	return &ListCinemasOutput{Cinemas: cinemas}, nil
}

func (s *Service) UpdateCinema(inp *UpdateCinemaInput) (*UpdateCinemaOutput, error) {
	var c models.Cinema
	if err := s.db.First(&c, inp.ID).Error; err != nil {
		return nil, fmt.Errorf("cinema not found")
	}

	if inp.CinemaName != "" {
		c.Name = inp.CinemaName
	}
	if inp.CityID != 0 {
		c.CityID = inp.CityID
	}
	if inp.CinemaOwnerID != 0 {
		c.CinemaOwnerID = inp.CinemaOwnerID
	}

	if err := s.db.Save(&c).Error; err != nil {
		_ = s.logger.Log("method", "UpdateCinema", "err", err)
		return nil, err
	}
	if err := s.db.Preload("City").Preload("CinemaOwner").First(&c, c.ID).Error; err != nil {
		return nil, err
	}
	s.cache.Delete("ListCinemasOutput")
	_ = s.logger.Log("method", "UpdateCinema", "cinema_id", c.ID)
	return &UpdateCinemaOutput{Cinema: c}, nil
}

func (s *Service) AddCinemaScreen(inp *AddCinemaScreenInput) (*AddCinemaScreenOutput, error) {
	// verify the cinema belongs to the requesting owner
	var cinema models.Cinema
	if err := s.db.First(&cinema, inp.CinemaID).Error; err != nil {
		return nil, fmt.Errorf("cinema not found")
	}
	if cinema.CinemaOwnerID != inp.CinemaOwnerID {
		return nil, fmt.Errorf("cinema does not belong to your account")
	}

	var existing models.CinemaScreen
	if err := s.db.Where("name = ? AND cinema_id = ?", inp.ScreenName, inp.CinemaID).First(&existing).Error; err == nil {
		return nil, fmt.Errorf("screen %q already exists in this cinema", inp.ScreenName)
	}

	var out AddCinemaScreenOutput

	txErr := s.db.Transaction(func(tx *gorm.DB) error {
		screen := models.CinemaScreen{
			Name:     inp.ScreenName,
			CinemaID: inp.CinemaID,
		}
		if err := tx.Create(&screen).Error; err != nil {
			return err
		}

		// generate seats with globally sequential numbers across all types
		seatNumber := 1
		var seats []models.CinemaSeat
		for seatType, count := range inp.Seats {
			for i := 0; i < count; i++ {
				seats = append(seats, models.CinemaSeat{
					SeatNumber:     seatNumber,
					Type:           seatType,
					CinemaScreenID: screen.ID,
				})
				seatNumber++
			}
		}
		if err := tx.Create(&seats).Error; err != nil {
			return err
		}

		out.CinemaScreen = buildScreenSummary(screen, seats)
		return nil
	})

	if txErr != nil {
		_ = s.logger.Log("method", "AddCinemaScreen", "err", txErr)
		return nil, txErr
	}

	s.cache.Delete("ListCinemasOutput")
	_ = s.logger.Log("method", "AddCinemaScreen", "screen_id", out.CinemaScreen.ID)
	return &out, nil
}

func (s *Service) ListScreens(inp *ListScreensInput) (*ListScreensOutput, error) {
	query := s.db.Preload("CinemaSeats")
	if inp.CinemaID != 0 {
		query = query.Where("cinema_id = ?", inp.CinemaID)
	}
	if inp.Name != "" {
		query = query.Where("LOWER(name) LIKE LOWER(?)", "%"+inp.Name+"%")
	}
	var screens []models.CinemaScreen
	if err := query.Find(&screens).Error; err != nil {
		_ = s.logger.Log("method", "ListScreens", "err", err)
		return nil, err
	}
	summaries := make([]CinemaScreenSummary, len(screens))
	for i, screen := range screens {
		summaries[i] = buildScreenSummary(screen, screen.CinemaSeats)
	}
	return &ListScreensOutput{CinemaScreens: summaries}, nil
}

// buildScreenSummary aggregates individual seats into a per-type count map.
func buildScreenSummary(screen models.CinemaScreen, seats []models.CinemaSeat) CinemaScreenSummary {
	summary := CinemaScreenSummary{
		ID:           screen.ID,
		Name:         screen.Name,
		CinemaID:     screen.CinemaID,
		SeatsSummary: make(map[string]int),
	}
	for _, seat := range seats {
		summary.SeatsSummary[string(seat.Type)]++
	}
	return summary
}

func (s *Service) ListSeats(inp *ListSeatsInput) (*ListSeatsOutput, error) {
	query := s.db.Where("cinema_screen_id = ?", inp.CinemaScreenID)
	if inp.Type != "" {
		query = query.Where("type = ?", inp.Type)
	}
	var seats []models.CinemaSeat
	if err := query.Find(&seats).Error; err != nil {
		_ = s.logger.Log("method", "ListSeats", "err", err)
		return nil, err
	}
	return &ListSeatsOutput{Seats: seats}, nil
}

func (s *Service) GetSeatByNumber(screenID, seatNumber int) (*GetSeatOutput, error) {
	var seat models.CinemaSeat
	err := s.db.Where("cinema_screen_id = ? AND seat_number = ?", screenID, seatNumber).First(&seat).Error
	if err != nil {
		return nil, fmt.Errorf("seat %d not found in screen %d", seatNumber, screenID)
	}
	return &GetSeatOutput{Seat: seat}, nil
}
