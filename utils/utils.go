package utils

import (
	"go-web-notepad/models"
	"math/rand"
	"strconv"
	"time"
)

func RemoveSliceBaseId(s []models.Memo, id int) []models.Memo {
	for index, memo := range s {
		if id == memo.Id {
			s = append(s[:index], s[index+1:]...)
			return s
		}
	}

	return s
}

func CreateMultiMemo(userId string, uniqueId int, amount int) []models.Memo {
	var memos []models.Memo
	for i := 0; i < amount; i++ {
		memo := models.Memo{
			Title:   "Title " + strconv.Itoa(i),
			Content: "Content " + strconv.Itoa(i),
			Writer:  uniqueId,
		}

		memos = append(memos, memo)
	}

	return memos
}

func GenerateRandomNumber(max int) int {
	seed := time.Now().UnixNano()
	randomGenerator := rand.New(rand.NewSource(seed))
	randomIntInRange := randomGenerator.Intn(max)
	return randomIntInRange
}
