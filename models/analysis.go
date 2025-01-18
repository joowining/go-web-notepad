package models

import (
	"database/sql"
	"fmt"
	"time"
)

type Analysis struct {
	Id        int       `json:"id"`
	Concept   string    `json:"concept"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	User      int       `json:"user"`
	CreatedAt time.Time `json:"createdat"`
}

// 특정 사용자의 분석결과들을 limit과 offset에 맞게 가져오기
func LoadAnalysisesByInput(userId int, start int, limit int, db *sql.DB) []Analysis {
	var analysises []Analysis
	rows, err := db.Query("SELECT * FROM Analysis WHERE User = ? ORDER BY id DESC LIMIT ? OFFSET ?;", userId, limit, (start-1)*limit)
	defer rows.Close()
	if err != nil {
		fmt.Println("database : ", err)
		return nil
	}

	for rows.Next() {
		var id int
		var concept string
		var title string
		var content string
		var user int
		var cAt time.Time

		err := rows.Scan(&id, &concept, &title, &content, &user, &cAt)
		if err != nil {
			fmt.Println("rows err : ", err)
			return nil
		}

		analysis := Analysis{
			Id:        id,
			Concept:   concept,
			Title:     title,
			Content:   content,
			User:      user,
			CreatedAt: cAt,
		}
		analysises = append(analysises, analysis)
	}
	return analysises
}

// 분석 결과 저장하기
func SaveAnalysis(data Analysis, db *sql.DB) error {
	_, err := db.Exec("INSERT INTO Analysis( Concept, Title, Content, User ) VALUES(?,?,?,?);", data.Concept, data.Title, data.Content, data.User)
	if err != nil {
		fmt.Print(err)
		return err
	}
	return nil
}

// 분석 결과 가져오기
func LoadAnalysisbyAnalysisId(id int, db *sql.DB) Analysis {
	var content Analysis
	err := db.QueryRow("SELECT Id, Concept, Title, Content, User FROM Analysis WHERE id = ? ;", id).Scan(&content.Id, &content.Concept, &content.Title, &content.Content, &content.User)
	if err != nil {
		fmt.Print(err)
		return Analysis{}
	}

	return content
}

// 특정 사용자가 가진 전체 메모의 갯수 가져오기
func FindUsersAnalysisCount(writerId int, db *sql.DB) int {
	var count int
	err := db.QueryRow("SELECT COUNT(*) AS count FROM Analysis WHERE user = ?;", writerId).Scan(&count)
	if err != nil {
		fmt.Println("database err : ", err)
		return -1
	}
	return count
}

// 특정 사용자의 특정 메모 삭제하기
func DeleteAnalysisByInput(id int, db *sql.DB) error {
	_, err := db.Exec("DELETE FROM Analysis WHERE id = ?", id)
	if err != nil {
		return err
	}
	return nil
}
