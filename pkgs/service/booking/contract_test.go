package booking

import (
	"testing"

	"gorm.io/gorm"

	"KnowWhoami/movie-ticketing/internal/testhelpers"
	"KnowWhoami/movie-ticketing/pkgs/models"
)

func TestBookSeatsInput_Validate(t *testing.T) {
	db := testhelpers.SetupDB()

	type fields struct {
		ShowID      int
		SeatNumbers []int
		UserID      int
		SeatType    models.SeatType
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
			name: "ValidInput",
			fields: fields{
				ShowID:      1,
				SeatNumbers: []int{1, 2, 3},
				UserID:      1,
				SeatType:    models.Recliner,
			},
			args: args{db: db},
		},
		{
			name: "ZeroShowID",
			fields: fields{
				ShowID:      0,
				SeatNumbers: []int{1, 2},
				UserID:      1,
				SeatType:    models.Premium,
			},
			args:    args{db: db},
			wantErr: true,
		},
		{
			name: "ZeroUserID",
			fields: fields{
				ShowID:      1,
				SeatNumbers: []int{1},
				UserID:      0,
				SeatType:    models.FrontRow,
			},
			args:    args{db: db},
			wantErr: true,
		},
		{
			name: "NilSeatNumbers",
			fields: fields{
				ShowID:      1,
				SeatNumbers: nil,
				UserID:      1,
				SeatType:    models.Recliner,
			},
			args:    args{db: db},
			wantErr: true,
		},
		{
			name: "EmptySeatType",
			fields: fields{
				ShowID:      1,
				SeatNumbers: []int{1},
				UserID:      1,
				SeatType:    "",
			},
			args:    args{db: db},
			wantErr: true,
		},
		{
			name: "AllZero",
			fields: fields{
				ShowID:      0,
				SeatNumbers: nil,
				UserID:      0,
				SeatType:    "",
			},
			args:    args{db: db},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bs := &BookSeatsInput{
				ShowID:      tt.fields.ShowID,
				SeatNumbers: tt.fields.SeatNumbers,
				UserID:      tt.fields.UserID,
				SeatType:    tt.fields.SeatType,
			}
			if err := bs.Validate(tt.args.db); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
