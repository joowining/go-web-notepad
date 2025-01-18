package routes

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

var duplicateCheck bool
var completeCreate bool
var disabled bool
var idcheck bool
var userInputId string

func CreateRouteGroup(e *echo.Echo, db *sql.DB) {
	routeGroup := e.Group("/create")

	routeGroup.GET("", createPageHandler)
	routeGroup.POST("", func(c echo.Context) error {
		return createUserHandler(c, db)
	})
	routeGroup.POST("/check", func(c echo.Context) error {
		return checkDuplicationHandler(c, db)
	})
}

func createPageHandler(c echo.Context) error {
	fmt.Println("create page handler check")

	data := map[string]interface{}{
		"DupCheck": duplicateCheck,
		"Disabled": disabled,
		"IdCheck":  idcheck,
		"Id":       userInputId,
	}
	return c.Render(http.StatusOK, "create.html", data)
}

// 동일한 아이디를 가진 사용자가 있는지 확인 SELECT
func checkDuplicationHandler(c echo.Context, db *sql.DB) error {
	userInputId = c.FormValue("id")
	var checkValue bool
	err := db.QueryRow("SELECT 1 FROM members WHERE userid = ? ", userInputId).Scan(&checkValue)
	fmt.Println("user check error :", err)
	if err == nil {
		return c.Redirect(http.StatusFound, "/create")
	}
	if !checkValue {
		duplicateCheck = true
		disabled = true
		idcheck = true
	}

	return c.Redirect(http.StatusFound, "/create")
}

// 사용자 정보를 받아서 사용자를 새롭게 생성하기 INSERT
func createUserHandler(c echo.Context, db *sql.DB) error {
	userId := c.FormValue("id")
	userPwd := c.FormValue("password")
	userName := c.FormValue("name")
	userEmail := c.FormValue("email")

	result, err := db.Exec("INSERT INTO members (userid, password, name, email) VALUES (?, ?, ?, ?)", userId, userPwd, userName, userEmail)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "database error"})
	}

	fmt.Println("create user : ", result)

	session, _ := store.Get(c.Request(), "my-session")
	session.Values["authorized"] = true
	session.Values["userId"] = userId
	session.Save(c.Request(), c.Response())

	return c.Redirect(http.StatusFound, "/memo/list")
}
