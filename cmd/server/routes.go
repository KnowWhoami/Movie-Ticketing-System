package server

import (
	"KnowWhoami/movie-ticketing/internal/server"
	"KnowWhoami/movie-ticketing/pkgs/api"
	"KnowWhoami/movie-ticketing/pkgs/models"
)

func (s *Server) setupRoutes() {
	apiHandler := api.NewAPIHandler(s.Db, s.Logger)
	adminOnly := api.RequireRole(models.UserTypeAdmin)
	theatreOwnerOnly := api.RequireRole(models.UserTypeCinemaOwner)
	ownerOrAdmin := api.RequireRole(models.UserTypeAdmin, models.UserTypeCinemaOwner)
	s.Routes = []server.Route{
		{
			Method:  "GET",
			Path:    "/healthz",
			Handler: defaultHandler("healthz"),
		},
		// auth routes (public)
		{
			Method:  "POST",
			Path:    "/register",
			Handler: wrapHandler(apiHandler.Register),
		},
		{
			Method:  "POST",
			Path:    "/login",
			Handler: wrapHandler(apiHandler.Login),
		},
		// user management routes
		{
			Method:  "POST",
			Path:    "/user",
			Handler: adminOnly(wrapHandler(apiHandler.CreateUser)),
		},
		{
			Method:  "GET",
			Path:    "/users",
			Handler: adminOnly(wrapHandler(apiHandler.ListUsers)),
		},
		{
			Method:  "GET",
			Path:    "/user/:id",
			Handler: adminOnly(wrapHandler(apiHandler.GetUserByID)),
		},
		{
			Method:  "PATCH",
			Path:    "/user/:id",
			Handler: adminOnly(wrapHandler(apiHandler.UpdateUser)),
		},
		{
			Method:  "DELETE",
			Path:    "/user/:id",
			Handler: adminOnly(wrapHandler(apiHandler.DeleteUser)),
		},
		// city routes
		{
			Method:  "POST",
			Path:    "/city",
			Handler: adminOnly(wrapHandler(apiHandler.AddCity)),
		},
		{
			Method:  "GET",
			Path:    "/cities",
			Handler: wrapHandler(apiHandler.ListCities),
		},
		// cinema routes
		{
			Method:  "GET",
			Path:    "/cinemas",
			Handler: wrapHandler(apiHandler.ListCinemas),
		},
		{
			Method:  "GET",
			Path:    "/cinemas/:cinema_owner_id",
			Handler: adminOnly(wrapHandler(apiHandler.GetCinemasByOwnerID)),
		},
		{
			Method:  "GET",
			Path:    "/cinema/:id",
			Handler: wrapHandler(apiHandler.GetCinemaByID),
		},
		{
			Method:  "GET",
			Path:    "/my-cinemas",
			Handler: theatreOwnerOnly(wrapHandler(apiHandler.GetMyCinemas)),
		},
		{
			Method:  "POST",
			Path:    "/cinema",
			Handler: adminOnly(wrapHandler(apiHandler.AddCinema)),
		},
		{
			Method:  "PATCH",
			Path:    "/cinema/:id",
			Handler: adminOnly(wrapHandler(apiHandler.UpdateCinema)),
		},
		{
			Method:  "POST",
			Path:    "/screen",
			Handler: theatreOwnerOnly(wrapHandler(apiHandler.AddCinemaScreen)),
		},
		{
			Method:  "GET",
			Path:    "/screens",
			Handler: wrapHandler(apiHandler.ListScreens),
		},
		{
			Method:  "GET",
			Path:    "/screens/:screen_id/seats",
			Handler: wrapHandler(apiHandler.ListSeats),
		},
		{
			Method:  "GET",
			Path:    "/screens/:screen_id/seats/:seat_number",
			Handler: wrapHandler(apiHandler.GetSeatByNumber),
		},
		// movie routes
		{
			Method:  "POST",
			Path:    "/movie",
			Handler: adminOnly(wrapHandler(apiHandler.AddMovie)),
		},
		{
			Method:  "GET",
			Path:    "/movies",
			Handler: wrapHandler(apiHandler.ListMovies),
		},
		{
			Method:  "GET",
			Path:    "/movie/:id",
			Handler: wrapHandler(apiHandler.GetMovieByID),
		},
		{
			Method:  "PATCH",
			Path:    "/movie/:id",
			Handler: adminOnly(wrapHandler(apiHandler.UpdateMovie)),
		},
		// show routes
		{
			Method:  "POST",
			Path:    "/show",
			Handler: theatreOwnerOnly(wrapHandler(apiHandler.AddShow)),
		},
		{
			Method:  "GET",
			Path:    "/shows",
			Handler: wrapHandler(apiHandler.ListShows),
		},
		{
			Method:  "GET",
			Path:    "/show/:id",
			Handler: wrapHandler(apiHandler.GetShowByID),
		},
		{
			Method:  "GET",
			Path:    "/show/:id/seats",
			Handler: wrapHandler(apiHandler.ListShowSeats),
		},
		{
			Method:  "PATCH",
			Path:    "/show/:id/cancel",
			Handler: theatreOwnerOnly(wrapHandler(apiHandler.CancelShow)),
		},
		// booking related routes
		{
			Method:  "GET",
			Path:    "/bookings",
			Handler: ownerOrAdmin(wrapHandler(apiHandler.ListBookings)),
		},
		{
			Method:  "GET",
			Path:    "/my-bookings",
			Handler: api.AuthMiddleware(wrapHandler(apiHandler.ListMyBookings)),
		},
		{
			Method:  "POST",
			Path:    "/booking",
			Handler: api.AuthMiddleware(wrapHandler(apiHandler.BookSeats)),
		},
	}
}
