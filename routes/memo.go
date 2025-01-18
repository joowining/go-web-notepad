package routes

import (
	"database/sql"
	"fmt"
	"math"
	"net/http"
	"strconv"

	"go-web-notepad/middleware"
	"go-web-notepad/models"
	"go-web-notepad/utils"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
)

var store = sessions.NewCookieStore([]byte("secret-key"))
var latestMemoId int

func MemoRouteGroup(e *echo.Echo, db *sql.DB) {
	routeGroup := e.Group("/memo")
	routeGroup.Use(middleware.LoginChecker)
	routeGroup.Use(middleware.IdParamChecker)

	routeGroup.GET("/list", func(c echo.Context) error {
		return listHandler(c, db)
	})
	routeGroup.GET("/list/:id", func(c echo.Context) error {
		return listHandler(c, db)
	})
	routeGroup.GET("/list/delete", func(c echo.Context) error {
		return listDeleteHandler(c, db)
	})
	routeGroup.POST("/list/delete", func(c echo.Context) error {
		return deleteCheckedHandler(c, db)
	})

	routeGroup.GET("/:id", func(c echo.Context) error {
		return memoHandler(c, db)
	})
	routeGroup.GET("/write", writePageHandler)
	routeGroup.POST("/write", func(c echo.Context) error {
		return writeMemoHandler(c, db)
	})

	routeGroup.GET("/edit/:id", func(c echo.Context) error {
		return editPageHandler(c, db)
	})
	routeGroup.POST("/edit/:id", func(c echo.Context) error {
		return editHandler(c, db)
	})
	routeGroup.GET("/delete/:id", func(c echo.Context) error {
		return deleteHandler(c, db)
	})

	routeGroup.GET("/dump/add/:user/:number", func(c echo.Context) error {
		return addDumpData(c, db)
	})

	routeGroup.GET("/dump/delete/:user/:number", func(c echo.Context) error {
		return deleteDumpData(c, db)
	})

	// test area
	routeGroup.GET("/listitem", func(c echo.Context) error {
		return listItemHandler(c, db)
	})
}

// 한 명의 사용자가 작성한 모든 메모 데이터 가져오기 SELECT
func listHandler(c echo.Context, db *sql.DB) error {
	var memos []models.Memo
	data := make(map[string]interface{})
	var userUniqueId int
	var userId string
	var start string
	var searchWord string

	start = c.Param("id")
	searchWord = c.FormValue("search")

	if start == "" {
		start = "1"
	}

	session, _ := store.Get(c.Request(), "my-session")
	userUniqueId = session.Values["uniqueId"].(int)

	userId = models.FindUserId(userUniqueId, db)
	if userId == "" {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "database error"})
	}

	data["User"] = userId
	data["Message"] = "logg in"

	c.Set("writerId", userUniqueId)

	startInt, _ := strconv.Atoi(start)
	// 검색하는 부분 통합하기
	if searchWord == "" {
		memos = models.FindUsersMemoAsInput(userUniqueId, startInt, 10, db)
	} else {
		// 검색어가 존재하면 해당 검색어에 맞게 내용을 찾도록 만듬
		memos = models.FindMemoBySearch(searchWord, db)
	}
	data["Memos"] = memos

	// totalCount 다르게 설정
	totalCount := models.FindUserMemosCount(userUniqueId, db)
	totalPage := int(math.Ceil(float64(totalCount) / 10))
	var pages []int

	pageGroup := int(math.Ceil(float64(startInt) / 10))
	startPage := ((pageGroup - 1) * 10) + 1
	for i := startPage; i < startPage+10; i++ {
		if i > totalPage {
			//continue
			break
		}

		pages = append(pages, i)
	}

	data["LeftBtn"] = false
	data["LeftVal"] = 0
	data["RightBtn"] = false
	data["RightVal"] = 0
	data["StartInt"] = startInt

	if len(pages) > 0 && pages[0] > 10 {
		data["LeftBtn"] = true
		data["LeftVal"] = pages[0] - 1
	}

	if len(pages) > 0 && pages[len(pages)-1] < totalPage {
		data["RightBtn"] = true
		data["RightVal"] = pages[len(pages)-1] + 1
	}

	data["Pages"] = pages
	if c.Get("delete") == true {
		data["DeletePart"] = true
	} else {
		data["DeletePart"] = false
	}

	return c.Render(http.StatusOK, "list.html", data)
}

func listDeleteHandler(c echo.Context, db *sql.DB) error {
	c.Set("delete", true)
	return listHandler(c, db)
}

func deleteCheckedHandler(c echo.Context, db *sql.DB) error {
	formParams, err := c.FormParams()
	if err != nil {
		return err
	}

	checkedOne := formParams["delete-box"]
	var ids []int
	for _, str := range checkedOne {
		// string을 int로 변환
		num, err := strconv.Atoi(str)
		if err != nil {
			return err // 변환 실패시 에러 반환
		}
		ids = append(ids, num)
	}

	err = models.DeleteMemoDynamicaly(db, ids)
	if err != nil {
		return err
	}

	return c.Redirect(http.StatusFound, "/memo/list")
}

// 메모 id값을 찾아서 제공하기 SELECT
func memoHandler(c echo.Context, db *sql.DB) error {
	memoId, _ := strconv.Atoi(c.Param("id"))

	content := models.FindUsersMemo(memoId, db)

	return c.Render(http.StatusOK, "memo.html", content)
}

func writePageHandler(c echo.Context) error {

	return c.Render(http.StatusOK, "write.html", "none")
}

// 메모 저장하기 INSERT
func writeMemoHandler(c echo.Context, db *sql.DB) error {

	// 사용자 정보 가져오기
	session, _ := store.Get(c.Request(), "my-session")
	writerId := session.Values["uniqueId"].(int)

	// 입력한 메모 정도 가져오기
	title := c.FormValue("title")
	content := c.FormValue("content")

	tempMemo := models.Memo{
		Title:   title,
		Content: content,
		Writer:  writerId,
	}
	// 입력한 메모정보를 데이터베이스에 입력
	_, err := models.SaveMemo(tempMemo, db)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "data not found"})
	}

	return c.Redirect(http.StatusFound, "/memo/list")
}

// 메모 삭제하기 DELETE
func deleteHandler(c echo.Context, db *sql.DB) error {
	memoId, _ := strconv.Atoi(c.Param("id"))

	// 실제 데이터베이스에서 지우고자 하는 데이터를 삭제
	err := models.DeleteMemo(memoId, db)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "data not found"})
	}

	return c.Redirect(http.StatusFound, "/memo/list")
}

// 기존의 메모 찾기 SELECT
// func editPageHandler(c echo.Context, db *sql.DB) error {
// 	memoId, _ := strconv.Atoi(c.Param("id"))
// 	var presentMemo models.Memo

// 	row, err := db.Query("SELECT * FROM memo WHERE id = ? ", memoId)
// 	if err != nil {
// 		return err
// 	}
// 	row.Scan(&presentMemo.Id, &presentMemo.Title, &presentMemo.Content)

// 	return c.Render(http.StatusOK, "edit.html", presentMemo)
// }

func editPageHandler(c echo.Context, db *sql.DB) error {
	memoId, _ := strconv.Atoi(c.Param("id"))
	var presentMemo models.Memo

	err := db.QueryRow("SELECT * FROM memo WHERE id = ? ", memoId).Scan(&presentMemo.Id, &presentMemo.Title, &presentMemo.Content, &presentMemo.Writer, &presentMemo.CreatedTime)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return c.Render(http.StatusOK, "edit.html", presentMemo)
}

// 메모 내용 업데이트 하기 UPDATE
func editHandler(c echo.Context, db *sql.DB) error {
	memoId, _ := strconv.Atoi(c.Param("id"))

	newTitle := c.FormValue("title")
	newContent := c.FormValue("content")

	tempMemo := models.Memo{
		Id:      memoId,
		Title:   newTitle,
		Content: newContent,
	}

	err := models.EditMemo(tempMemo, db)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "database error"})
	}

	return c.Redirect(http.StatusFound, "/memo/list")
}

// 덤프 데이터 추가하기
func addDumpData(c echo.Context, db *sql.DB) error {
	var memos []models.Memo

	user := c.Param("user")
	amount, err := strconv.Atoi(c.Param("number"))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "parse error with param data"})
	}

	uniqueId := models.FindUserUniqueId(user, db)
	if uniqueId < 0 {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Cann't find appropriate user"})
	}

	memos = utils.CreateMultiMemo(user, uniqueId, amount)

	for i := 0; i < amount; i++ {
		_, _ = models.SaveMemo(memos[i], db)
	}

	return c.Redirect(http.StatusFound, "/memo/list")
}

// 덤프 데이터 삭제하기
func deleteDumpData(c echo.Context, db *sql.DB) error {
	fmt.Println("dump delete start")
	user := c.Param("user")
	amount, err := strconv.Atoi(c.Param("number"))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "parse error with param data"})
	}

	uniqueId := models.FindUserUniqueId(user, db)
	if amount == 0 {
		models.DeleteAllMemo(uniqueId, db)
		return c.Redirect(http.StatusFound, "/memo/list")
	}

	_, err = db.Exec(`DELETE memo
						FROM memo
						JOIN (
							SELECT id
							FROM memo
							WHERE writer = ?
							ORDER BY id DESC
							LIMIT ?
						) AS to_delete ON memo.id = to_delete.id;
					`, uniqueId, amount)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("dump delete end")
	return c.Redirect(http.StatusFound, "/memo/list")
}

// test area
func listItemHandler(c echo.Context, db *sql.DB) error {
	session, _ := store.Get(c.Request(), "my-session")
	userUniqueId := session.Values["uniqueId"].(int)
	memos := models.FindUsersMemoAsInput(userUniqueId, 1, 10, db)
	fmt.Println("Here is listItemHandler")
	data := make(map[string]interface{})

	data["memos"] = memos

	fmt.Println(data)
	return c.JSON(http.StatusOK, data)
}
