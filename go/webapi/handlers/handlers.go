package handlers

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/nanoteck137/{{ .ProjectName }}/database"
)

type Handlers struct {
	validate   *validator.Validate
	db         *database.Database
}

func New(db *database.Database) *Handlers {
	var validate = validator.New()
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})

	return &Handlers{
		validate:   validate,
		db:         db,
	}
}

func (api *Handlers) validateBody(body any) map[string]string {
	err := api.validate.Struct(body)
	if err != nil {
		type ValidationError struct {
			Field   string `json:"field"`
			Message string `json:"message"`
		}

		validationErrs := make(map[string]string)
		for _, err := range err.(validator.ValidationErrors) {
			field := err.Field()
			validationErrs[field] = fmt.Sprintf("'%v' not satisfying tags '%v'", field, err.Tag())
		}

		return validationErrs
	}

	return nil
}

func ConvertURL(c echo.Context, path string) string {
	host := c.Request().Host

	return fmt.Sprintf("http://%s%s", host, path)
}
