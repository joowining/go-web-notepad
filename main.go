package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"

	"go-web-notepad/routes"

	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type TemplateRenderer struct {
	templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

// sessions
var store = sessions.NewCookieStore([]byte("secret-key"))

// DB
var db *sql.DB

func main() {

	cfg := mysql.Config{
		User:      os.Getenv("DB_USER"),
		Passwd:    os.Getenv("DB_PASSWORD"),
		Net:       "tcp",
		Addr:      "localhost:3306",
		DBName:    "notepad",
		ParseTime: true,
	}

	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		fmt.Println("52", err)
	}
	pingErr := db.Ping()
	if pingErr != nil {
		fmt.Println("56", pingErr)
	}
	fmt.Println("DB Connected!")

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Static("/static", "assets")

	t := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("templates/**/*.html")),
	}

	e.Renderer = t

	e.GET("/", rootHandler)
	// e.GET("/", redirectHandler)
	e.POST("/", func(c echo.Context) error {
		return loginHandler(c, db)
	})
	e.GET("/logout", logoutHandler)
	e.GET("/any", anyHandler)

	e.HTTPErrorHandler = func(err error, c echo.Context) {
		code := http.StatusInternalServerError
		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
		}

		data := struct {
			Code    int
			Message string
		}{
			Code:    code,
			Message: "Sorry there is no page for you",
		}
		// 에러 코드에 따라 템플릿을 렌더링
		switch code {
		case http.StatusNotFound:
			err := c.Render(http.StatusNotFound, "error.html", data)
			if err != nil {
				c.Logger().Error(err)
			}
		case http.StatusInternalServerError:
			c.Render(http.StatusInternalServerError, "error.html", data)
		default:
			c.Render(code, "error.html", nil)
		}
	}

	routes.MemoRouteGroup(e, db)
	routes.CreateRouteGroup(e, db)
	routes.AnalysisRouteGroup(e, db)
	routes.UserRouterGroup(e)

	e.Logger.Fatal(e.Start(":8000"))

}

func rootHandler(c echo.Context) error {
	session, _ := store.Get(c.Request(), "my-session")

	// 사용자 목록을 콘텍스트에 저장

	if _, ok := session.Values["authorized"].(bool); ok {
		return c.Redirect(http.StatusFound, "/memo/list")
	} else {
		return c.Render(http.StatusOK, "index-tailwind.html", nil)
	}
}

func loginHandler(c echo.Context, db *sql.DB) error {
	session, _ := store.Get(c.Request(), "my-session")

	userId := c.FormValue("user-id")
	password := c.FormValue("user-pwd")

	fmt.Println("login handler userID and Pwd", userId, password)
	var id int
	err := db.QueryRow("SELECT id FROM members WHERE userid = ? AND password = ?", userId, password).Scan(&id)
	fmt.Println("id : ", id)
	fmt.Println("out err :", err)
	if err == nil {
		fmt.Println("in if block")
		session.Values["authorized"] = true
		session.Values["userId"] = userId
		session.Values["uniqueId"] = id
		session.Save(c.Request(), c.Response())
		fmt.Println("auth in login: ", session.Values["authorized"])
		return c.Redirect(http.StatusFound, "/memo/list")
	} else {
		fmt.Println(err)
		return c.Redirect(http.StatusOK, "/")
	}
}

func logoutHandler(c echo.Context) error {
	session, _ := store.Get(c.Request(), "my-session")

	session.Options.MaxAge = -1                        // 세션 만료 설정
	session.Values = make(map[interface{}]interface{}) // 세션 값 초기화

	// 세션 저장
	session.Save(c.Request(), c.Response())

	return c.Redirect(http.StatusFound, "/")
}

func redirectHandler(c echo.Context) error {
	session, _ := store.Get(c.Request(), "my-session")

	if _, ok := session.Values["authorized"].(bool); ok {
		return c.File("assets/index.html")
	} else {
		return c.Redirect(http.StatusFound, "/memo/list")
	}
}

func anyHandler(c echo.Context) error {
	return c.Render(http.StatusOK, "any.html", nil)
}
