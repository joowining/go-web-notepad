package models

import (
	"database/sql"
	"fmt"
	"strings"
)

type Memo struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	Writer      int    `json:"writer"`
	CreatedTime string `json:"time"`
}

// 한 사용자의 모든 메모 찾기
func LoadUserALLMemos(writerId int, db *sql.DB) []Memo {
	var memos []Memo

	rows, err := db.Query("SELECT id, title, content, createdTime FROM memo WHERE writer = ?", writerId)
	defer rows.Close()

	if err != nil {
		fmt.Println("database ", err)
		return nil
	}

	for rows.Next() {
		var id int
		var title string
		var content string
		var createdTime string

		err := rows.Scan(&id, &title, &content, &createdTime)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		memo := Memo{
			Id:          id,
			Title:       title,
			Content:     content,
			Writer:      writerId,
			CreatedTime: createdTime,
		}
		memos = append(memos, memo)
	}

	return memos
}

// 특정 사용자가 가진 전체 메모의 갯수 가져오기
func FindUserMemosCount(writerId int, db *sql.DB) int {
	var count int
	err := db.QueryRow("SELECT COUNT(*) AS count FROM memo WHERE writer = ? ;", writerId).Scan(&count)
	if err != nil {
		fmt.Println(err)
		return -1
	}
	return count
}

// 특정 갯수의 데이터만 특정 시작점으로부터 가져오기
func FindUsersMemoAsInput(writerId int, start int, limit int, db *sql.DB) []Memo {
	var memos []Memo
	fmt.Println("value : ", (start-1)*limit)
	rows, err := db.Query("SELECT * FROM memo WHERE writer = ? ORDER BY id DESC LIMIT ? OFFSET ?", writerId, limit, (start-1)*limit)
	defer rows.Close()

	if err != nil {
		fmt.Println("database ", err)
		return nil
	}

	for rows.Next() {
		var id int
		var title string
		var content string
		var writer int
		var createdTime string

		err := rows.Scan(&id, &title, &content, &writer, &createdTime)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		memo := Memo{
			Id:          id,
			Title:       title,
			Content:     content,
			Writer:      writer,
			CreatedTime: createdTime,
		}
		memos = append(memos, memo)
	}

	return memos
}

// 사용자가 사용하는 아이디 찾기
func FindUserId(uniqueId int, db *sql.DB) string {
	var userId string
	err := db.QueryRow("SELECT userid FROM members WHERE id = ?", uniqueId).Scan(&userId)
	if err != nil {
		return ""
	}

	return userId
}

// 사용자의 식별 아이디 찾기
func FindUserUniqueId(userId string, db *sql.DB) int {
	var uniqueId int
	err := db.QueryRow("SELECT id FROM members WHERE userid = ?", userId).Scan(&uniqueId)
	if err != nil {
		return -1
	}
	return uniqueId
}

// 사용자의 특정 메모 내용 찾기
func FindUsersMemo(memoId int, db *sql.DB) Memo {
	var content Memo
	err := db.QueryRow("SELECT id, title, content FROM memo WHERE id = ?", memoId).Scan(&content.Id, &content.Title, &content.Content)
	if err != nil {
		if err == sql.ErrNoRows {
			return Memo{}
		}
		return Memo{}
	}
	return content
}

// 사용자가 입력한 검색어와 일치 , 포함하는 메모 내용 찾기
func FindMemoBySearch(word string, db *sql.DB) []Memo {
	var memos []Memo

	rows, err := db.Query("SELECT * FROM memo WHERE title REGEXP ?", word)
	defer rows.Close()
	if err != nil {
		fmt.Println("DB error : ", err)
		return nil
	}

	for rows.Next() {
		var id int
		var title string
		var content string
		var writer int
		var createdTime string

		err := rows.Scan(&id, &title, &content, &writer, &createdTime)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		memo := Memo{
			Id:          id,
			Title:       title,
			Content:     content,
			Writer:      writer,
			CreatedTime: createdTime,
		}
		memos = append(memos, memo)
	}

	return memos
}

// 메모 저장하기
func SaveMemo(data Memo, db *sql.DB) (int, error) {
	result, err := db.Exec("INSERT INTO memo ( title, content, writer) VALUES( ?, ?, ?)", data.Title, data.Content, data.Writer)
	lastId, _ := result.LastInsertId()
	if err != nil {
		fmt.Println(err)
		return int(lastId), err
	}
	return int(lastId), nil
}

// memoId로 메모 삭제하기
func DeleteMemo(memoId int, db *sql.DB) error {
	_, err := db.Exec("DELETE FROM memo WHERE id = ?", memoId)
	if err != nil {
		return err
	}
	return nil
}

func DeleteAllMemo(writerId int, db *sql.DB) error {
	_, err := db.Exec("DELETE FROM memo WHERE writer = ?", writerId)
	if err != nil {
		return err
	}
	return nil
}

// writer로 메모 하나만 삭제하기
func DeleteWriterMemoOnce(writer int, memoId int, db *sql.DB) error {
	_, err := db.Exec("DELETE FROM memo WHERE writer = ? AND id = ? ", writer, memoId)
	if err != nil {
		return err
	}
	return nil
}

func DeleteMemoDynamicaly(db *sql.DB, ids []int) error {
	if len(ids) == 0 {
		return nil // 삭제할 항목이 없는 경우
	}

	// IN 구문을 위한 플레이스홀더 생성
	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = fmt.Sprintf("?") // $1, $2, ...
		args[i] = id
	}

	// DELETE 쿼리 생성
	query := fmt.Sprintf("DELETE FROM memo WHERE id IN (%s)",
		strings.Join(placeholders, ","))

	// 쿼리 실행
	_, err := db.Exec(query, args...)
	return err
}

// 메모 수정하기
func EditMemo(data Memo, db *sql.DB) error {
	_, err := db.Exec("UPDATE memo SET title = ?, content = ? WHERE id = ?", data.Title, data.Content, data.Id)
	if err != nil {
		return err
	}
	return nil
}
