package cinema

import (
	"testing"

	"gorm.io/gorm"

	"KnowWhoami/movie-ticketing/internal/testhelpers"
	"KnowWhoami/movie-ticketing/pkgs/models"
)

func TestAddCinemaInput_Validate(t *testing.T) {
	db := testhelpers.SetupDB()
	type fields struct {
		CinemaName string
		CityID     int
	}
	type args struct {
		db *gorm.DB
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "MissingCity",
			fields:  fields{CinemaName: "test", CityID: 1},
			args:    args{db: db},
			wantErr: true,
		},
		{
			name:    "EmptyName",
			fields:  fields{CinemaName: "", CityID: 1},
			args:    args{db: db},
			wantErr: true,
		},
		{
			name:    "ZeroCityID",
			fields:  fields{CinemaName: "test", CityID: 0},
			args:    args{db: db},
			wantErr: true,
		},
		{
			name:    "BothEmpty",
			fields:  fields{CinemaName: "", CityID: 0},
			args:    args{db: db},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ac := &AddCinemaInput{
				CinemaName: tt.fields.CinemaName,
				CityID:     tt.fields.CityID,
			}
			if err := ac.Validate(tt.args.db); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAddCinemaScreenInput_Validate(t *testing.T) {
	db := testhelpers.SetupDB()
	type fields struct {
		CinemaID   int
		ScreenName string
		Seats      []*SeatInfo
	}
	type args struct {
		db *gorm.DB
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "MissingCinema",
			fields:  fields{CinemaID: 1, ScreenName: "test", Seats: []*SeatInfo{}},
			args:    args{db: db},
			wantErr: true,
		},
		{
			name:    "ZeroCinemaID",
			fields:  fields{CinemaID: 0, ScreenName: "test", Seats: []*SeatInfo{}},
			args:    args{db: db},
			wantErr: true,
		},
		{
			name:    "EmptyScreenName",
			fields:  fields{CinemaID: 1, ScreenName: "", Seats: []*SeatInfo{}},
			args:    args{db: db},
			wantErr: true,
		},
		{
			name: "NilSeats",
			fields: fields{
				CinemaID:   1,
				ScreenName: "test",
				Seats:      nil,
			},
			args:    args{db: db},
			wantErr: true,
		},
		{
			name: "ValidSeatInfo",
			fields: fields{
				CinemaID:   1,
				ScreenName: "test",
				Seats: []*SeatInfo{
					{SeatNumber: 1, SeatType: models.Recliner},
				},
			},
			args:    args{db: db},
			wantErr: true, // cinema_id=1 doesn't exist in test DB
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			acs := &AddCinemaScreenInput{
				CinemaID:   tt.fields.CinemaID,
				ScreenName: tt.fields.ScreenName,
				Seats:      tt.fields.Seats,
			}
			if err := acs.Validate(tt.args.db); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
