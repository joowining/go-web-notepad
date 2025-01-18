package routes

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"

	"go-web-notepad/middleware"
	"go-web-notepad/models"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

type ChatRequest struct {
	Model    string `json:"model"`
	Messages []struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"messages"`
	Temperature float64 `json:"temperature"`
}

type ChatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func AnalysisRouteGroup(e *echo.Echo, db *sql.DB) {
	routeGroup := e.Group("/analysis")
	routeGroup.Use(middleware.LoginChecker)

	routeGroup.GET("/list", func(c echo.Context) error {
		return analysisListHandler(c, db)
	})
	routeGroup.POST("/chat", func(c echo.Context) error {
		return gptResponse(c, db)
	})
	routeGroup.GET("/delete/:id", func(c echo.Context) error {
		return deleteAnalysisHandler(c, db)
	})

	routeGroup.GET("/:id", func(c echo.Context) error {
		return analysisHandler(c, db)
	})
}

func analysisListHandler(c echo.Context, db *sql.DB) error {
	session, _ := store.Get(c.Request(), "my-session")
	writerId := session.Values["uniqueId"].(int)
	var analysises []models.Analysis
	var start int

	if "" == c.Param("id") {
		start = 1
	} else {
		start, _ = strconv.Atoi(c.Param("id"))
	}

	data := make(map[string]interface{})
	analysises = models.LoadAnalysisesByInput(writerId, start, 10, db)
	data["Analysises"] = analysises

	// pagination
	totalCount := models.FindUsersAnalysisCount(writerId, db)
	totalPage := int(math.Ceil(float64(totalCount) / 10))
	var pages []int

	pageGroup := int(math.Ceil(float64(start) / 10))
	startPage := ((pageGroup - 1) * 10) + 1
	for i := startPage; i < startPage+10; i++ {
		if i > totalPage {
			break
		}
		pages = append(pages, i)
	}

	data["Pages"] = pages
	data["LeftBtn"] = false
	data["LeftVal"] = 0
	data["RightBtn"] = false
	data["RightVal"] = 0
	data["StartInt"] = start

	if len(pages) > 0 && pages[0] > 10 {
		data["LeftBtn"] = true
		data["LeftVal"] = pages[0] - 1
	}

	if len(pages) > 0 && pages[len(pages)-1] < totalPage {
		data["RightBtn"] = true
		data["RightVal"] = pages[len(pages)-1] + 1
	}

	return c.Render(http.StatusFound, "analysis-list.html", data)
}

func gptResponse(c echo.Context, db *sql.DB) error {
	var message string
	var concept string
	session, _ := store.Get(c.Request(), "my-session")
	writerId := session.Values["uniqueId"].(int)
	conceptOption, _ := strconv.Atoi(c.FormValue("concepts"))
	numberOption, _ := strconv.Atoi(c.FormValue("numbers"))

	var memos []models.Memo
	memos = models.FindUsersMemoAsInput(writerId, 1, numberOption, db)
	// fmt.Print(memos)

	if conceptOption == 1 {
		message = fmt.Sprintf(`다음에 작성된 나의 일기를 바탕으로 내가 어떠한 사람인지 분석해서 알려줘 알려줄 때 그 내용의 제목도 함께 보내줘 다음 데이터를 JSON 형식으로 반환해주세요. 다음은 예시입니다:
		{
 			"Title": "예시 제목",
 	 		"Content": "예시 내용"
		} 반드시 JSON 형식으로만 응답하세요.`)
		concept = "나를 알아가기"
	} else if conceptOption == 2 {
		message = fmt.Sprintf(`다음에 작성된 나의 일기를 바탕으로 내가 다음에 하면 좋을 활동을 하나만 추천해줘 알려줄 때 그 내용의 제목도 함께 보내줘 다음 데이터를 JSON 형식으로 반환해주세요. 다음은 예시입니다:
		{
 			"Title": "예시 제목",
 	 		"Content": "예시 내용"
		} 반드시 JSON 형식으로만 응답하세요.`)
		concept = "활동 제안하기"
	} else if conceptOption == 3 {
		message = fmt.Sprintf(`다음에 작성된 나의 일기를 바탕으로 다음에 하면 좋을 콘텐츠를 하나만 구체적으로 추천해줘 알려줄 때 그 내용의 제목도 함께 보내줘 다음 데이터를 JSON 형식으로 반환해주세요. 다음은 예시입니다:
		{
 			"Title": "예시 제목",
 	 		"Content": "예시 내용"
		} 반드시 JSON 형식으로만 응답하세요.`)
		concept = "콘텐츠 추천하기"
	} else {
		return nil
	}

	for _, memo := range memos {
		message = message + "title : " + memo.Title + "\n"
		message = message + "content: " + memo.Content + "\n"
	}

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Retrieve the API key
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("API_KEY not found in .env file")
	}

	// var userMessage struct {
	// 	Message string `json:"message"`
	// }
	// if err := c.Bind(&userMessage); err != nil {
	// 	return c.String(http.StatusBadRequest, "Invalid request")
	// }

	requestBody := ChatRequest{
		Model: "gpt-4o",
		Messages: []struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		}{
			{"user", message},
		},
		Temperature: 0.0,
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	responseData, _ := io.ReadAll(resp.Body)
	var chatResponse ChatResponse
	json.Unmarshal(responseData, &chatResponse)

	messageVal := chatResponse.Choices[0].Message.Content
	fmt.Println("받은 것: ", messageVal)

	if strings.HasPrefix(messageVal, "```") {
		messageVal = strings.TrimPrefix(messageVal, "```json")
		messageVal = strings.TrimPrefix(messageVal, "```")
		messageVal = strings.TrimSuffix(messageVal, "```")
	}

	var analysis models.Analysis
	err = json.Unmarshal([]byte(messageVal), &analysis)
	if err != nil {
		fmt.Println("err in marshaling : ", err)
		return err
	}

	analysis.Concept = concept
	analysis.User = writerId

	err = models.SaveAnalysis(analysis, db)

	if err != nil {
		fmt.Println("err in saving analysis : ", err)
		return err
	}

	return c.Redirect(http.StatusFound, "/analysis/list")
}

func analysisHandler(c echo.Context, db *sql.DB) error {
	var analysis models.Analysis
	id, _ := strconv.Atoi(c.Param("id"))
	analysis = models.LoadAnalysisbyAnalysisId(id, db)

	data := make(map[string]interface{})

	data["Content"] = analysis.Content
	data["Title"] = analysis.Title
	data["Concept"] = analysis.Concept

	return c.Render(http.StatusOK, "analysis.html", data)
}

func deleteAnalysisHandler(c echo.Context, db *sql.DB) error {
	id, _ := strconv.Atoi(c.Param("id"))

	err := models.DeleteAnalysisByInput(id, db)
	if err != nil {
		return err
	}

	return c.Redirect(http.StatusFound, "/analysis/list")
}
