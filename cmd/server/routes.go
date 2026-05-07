package server

import (
	"KnowWhoami/movie-ticketing/internal/server"
	"KnowWhoami/movie-ticketing/pkgs/api"
	"KnowWhoami/movie-ticketing/pkgs/models"
)

func (s *Server) setupRoutes() {
	apiHandler := api.NewAPIHandler(s.Db, s.Logger)
	theatreOwnerOnly := api.RequireRole(models.UserTypeTheatreOwner)
	regularOnly := api.RequireRole(models.UserTypeRegular)
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
		// cinema routes
		{
			Method:  "GET",
			Path:    "/cinemas",
			Handler: wrapHandler(apiHandler.ListCinemas),
		},
		{
			Method:  "POST",
			Path:    "/cinema",
			Handler: theatreOwnerOnly(wrapHandler(apiHandler.AddCinema)),
		},
		{
			Method:  "POST",
			Path:    "/screen",
			Handler: theatreOwnerOnly(wrapHandler(apiHandler.AddCinemaScreen)),
		},
		// movie routes
		{
			Method:  "GET",
			Path:    "/show",
			Handler: wrapHandler(apiHandler.GetShow),
		},
		{
			Method:  "POST",
			Path:    "/show",
			Handler: theatreOwnerOnly(wrapHandler(apiHandler.AddShow)),
		},
		{
			Method:  "POST",
			Path:    "/movie",
			Handler: theatreOwnerOnly(wrapHandler(apiHandler.AddMovie)),
		},
		// booking related routes
		{
			Method:  "GET",
			Path:    "/bookings",
			Handler: wrapHandler(apiHandler.ListBookings),
		},
		{
			Method:  "POST",
			Path:    "/book",
			Handler: regularOnly(wrapHandler(apiHandler.BookSeats)),
		},
	}
}
