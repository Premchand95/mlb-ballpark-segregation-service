package front_service

import (
	"encoding/json"
	"fmt"
	"time"

	scheduler "github.com/mlb/mlb-ballpark-segregation-service/front_service/gen/scheduler"
)

/*
* isNotNil
* Description: returns false if pointer value is empty
 */
func isNotNil[T any](i *T) bool {
	return i != nil
}

/*
* prettyPrint: returns the struct in readable JSON string
* TODO: handle the marshalIndent error
 */
func prettyPrint[T any](temp T) string {
	//MarshalIndent
	empJSON, err := json.MarshalIndent(temp, "", "  ")
	if err != nil {
		return ""
	}
	return string(empJSON)
}

/*
* getIndex
* Description: returns the index of a game
 */
func indexOfGame(s []*scheduler.Game, x *scheduler.Game) int {
	for i := range s {
		if *s[i].GamePk == *x.GamePk {
			return i
		}
	}
	return -1
}

/*
* GetIndexOfGames
* Description: returns the array of indices
 */
func GetIndexOfGames(s []*scheduler.Game, x []*scheduler.Game) []int {
	var indexArray []int
	for i := range x {
		indexArray = append(indexArray, indexOfGame(s, x[i]))
	}
	return indexArray
}

/*
* deleteGame
* Description: deletes the value from the array
 */
func deleteGame(s []*scheduler.Game, i int) []*scheduler.Game {
	if i == len(s)-1 {
		return s[:i]
	}
	return append(s[:i], s[i+1:]...)
}

/*
* incrementByOne
* Description: increment every value in a array by 1
 */
func incrementByOne(slice []int) []int {
	for i := range slice {
		slice[i]++
	}
	return slice
}

/*
* parseDate
* Description: parse the date string and returns in Time value
 */
func parseDate(date string) (time.Time, error) {
	tm, err := time.Parse(string(DateFormat), date)
	if err != nil {
		return tm, fmt.Errorf("error parsing date format: %s", date)
	}
	return tm, nil
}
