package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type UserInfo struct {
	UserId  string `json:"userId"`
	UserPwd string `json:"userPwd"`
}

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Static("/static", "assets")

	e.File("/", "assets/index.html")
	e.File("/list", "assets/list.html")
	e.File("/write", "assets/write.html")
	e.File("/memo", "assets/memo.html")
	e.File("/create", "assets/create.html")

	e.POST("/", checkUser)

	//e.GET("/list", serveList)

	e.Logger.Fatal(e.Start(":8000"))

}

func checkUser(c echo.Context) error {
	fmt.Println("checkUser")
	sampUsr := UserInfo{
		UserId:  "kim",
		UserPwd: "1234",
	}

	var u UserInfo
	if err := c.Bind(&u); err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "wrong data"})
	}
	fmt.Println(u)
	if sampUsr == u {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "성공",
			"data":    u,
		})
	} else {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "different user"})
	}
}

func serveList(c echo.Context) error {
	return c.File("assets/list.html")
}
