package middleware

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
)

var store = sessions.NewCookieStore([]byte("secret-key"))

func LoginChecker(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		session, err := store.Get(c.Request(), "my-session")
		if err != nil {
			fmt.Println(" login checker err :", err)
			return c.Redirect(http.StatusFound, "/")
		}

		auth, ok := session.Values["authorized"].(bool)
		if !ok || !auth {
			return c.Redirect(http.StatusFound, "/")
		}

		user := session.Values["userId"]
		c.Set("user", user)

		return next(c)
	}
}

func IdParamChecker(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		if "" == c.Param("id") {
			return next(c)
		} else {
			_, err := strconv.Atoi(c.Param("id"))
			errData := struct {
				Code    int
				Message string
			}{
				Code:    404,
				Message: "Sorry there is no page for you",
			}
			if err != nil {
				return c.Render(http.StatusNotFound, "error.html", errData)
			}

			// if memoId > memos[len(memos)-1].Id {
			// 	return c.Render(http.StatusNotFound, "error.html", errData)
			// }
			return next(c)
		}

	}
}
