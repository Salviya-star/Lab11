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

// GetHello возвращает случайное приветствие пользователю
func (srv *Server) GetHello(e echo.Context) error {
	msg, err := srv.uc.FetchHelloMessage()
	if err != nil {
		return e.String(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusOK, msg)
}

// PostHello Помещает новый вариант приветствия в БД
func (srv *Server) PostHello(e echo.Context) error {
	input := struct {
		Msg *string `json:"msg"`
	}{}

	err := e.Bind(&input)
	if err != nil {
		return e.String(http.StatusInternalServerError, err.Error())
	}

	if input.Msg == nil {
		return e.String(http.StatusBadRequest, "msg is empty")
	}

	if len([]rune(*input.Msg)) > srv.maxSize {
		return e.String(http.StatusBadRequest, "hello message too large")
	}

	err = srv.uc.SetHelloMessage(*input.Msg)
	if err != nil {
		if errors.Is(err, vars.ErrAlreadyExist) {
			return e.String(http.StatusConflict, err.Error())
		}
		return e.String(http.StatusInternalServerError, err.Error())
	}

	return e.String(http.StatusCreated, "OK")
}