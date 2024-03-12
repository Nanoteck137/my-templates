package server

import (
	"github.com/MadAppGang/httplog/echolog"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nanoteck137/{{ .ProjectName }}/database"
	"github.com/nanoteck137/{{ .ProjectName }}/handlers"
	"github.com/nanoteck137/{{ .ProjectName }}/types"
)

func New(db *database.Database) *echo.Echo {
	e := echo.New()

	e.HTTPErrorHandler = func(err error, c echo.Context) {
		switch err := err.(type) {
		case *types.ApiError:
			c.JSON(err.Code, types.ApiResponse{
				Status: types.StatusError,
				Error:  err,
			})
		default:
			c.JSON(500, types.ApiResponse{
				Status: types.StatusError,
				Error: &types.ApiError{
					Code:    500,
					Message: err.Error(),
				},
			})
		}

	}

	e.Use(echolog.LoggerWithName("{{ .ProjectName }}"))
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	h := handlers.New(db)
	apiGroup := e.Group("/api/v1")

	h.InstallTodoHandlers(apiGroup)

	return e
}
