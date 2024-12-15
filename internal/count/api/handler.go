package api

import (
	"Lab11/pkg/vars"
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func (srv *Server) Login(e echo.Context) error {
	username := e.FormValue("username")
	password := e.FormValue("password")

	if username != "admin" || password != "admin" {
		return echo.ErrUnauthorized
	}

	claims := &jwtCustomClaims{
		"admin",
		true,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return err
	}

	return e.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
}

// GetCount возвращает cчетчик
func (srv *Server) GetCount(e echo.Context) error {
	count, err := srv.uc.SelectCount()
	if err != nil {
		return e.String(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusOK, map[string]int{"count": count})
}

// PostСount Помещает новый вариант в БД
func (srv *Server) PostCount(e echo.Context) error {
	input := struct {
		Count int `json:"count"`
	}{}

	err := e.Bind(&input)
	if err != nil {
		return e.String(http.StatusInternalServerError, "Invalid input data: "+err.Error())
	}

	if input.Count <= 0 {
		return e.String(http.StatusBadRequest, "count is negative")
	}

	err = srv.uc.UpdateCount(input.Count)
	if err != nil {
		if errors.Is(err, vars.ErrAlreadyExist) {
			return e.String(http.StatusConflict, err.Error())
		}
		return e.String(http.StatusInternalServerError, err.Error())
	}

	return e.String(http.StatusCreated, "Count updated successfully")
}
