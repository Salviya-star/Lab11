package api

import (
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

// Get возвращает приветствие пользователю
func (srv *Server) GetQuery(e echo.Context) error {
	name := e.QueryParam("name")
	msg, err := srv.uc.SelectNameQuery(name)
	if err != nil {
		return e.String(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusOK, msg)
}

// Post Помещает новый вариант приветствия в БД
func (srv *Server) PostQuery(e echo.Context) error {
	name := e.QueryParam("name")

	err := srv.uc.InsertNameQuery(name)
	if err != nil {
		return e.String(http.StatusInternalServerError, err.Error())
	}

	return e.NoContent(http.StatusOK)
}
