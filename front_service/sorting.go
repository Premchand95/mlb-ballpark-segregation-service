package front_service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	design "github.com/mlb/mlb-ballpark-segregation-service/front_service/design"
	scheduler "github.com/mlb/mlb-ballpark-segregation-service/front_service/gen/scheduler"
)

func (s *schedulersrvc) getFavoriteTeamGames(ctx context.Context, teamID uint, games []*scheduler.Game) []*scheduler.Game {
	var favTeamGames []*scheduler.Game
	for _, game := range games {
		teams := s.extractTeams(game)
		for _, team := range teams {
			if isNotNil(team.ID) {
				if *team.ID == teamID {
					// if you find your fav team in this game, no need to check for opponent game
					// so break out
					favTeamGames = append(favTeamGames, game)
					break
				}
			}
		}
	}
	return favTeamGames
}

func (s *schedulersrvc) extractTeams(game *scheduler.Game) []*scheduler.Team {
	teams := make([]*scheduler.Team, 0)
	if isNotNil(game) {
		if isNotNil(game.Teams) {
			//Away
			if isNotNil(game.Teams.Away) {
				if isNotNil(game.Teams.Away.Team) {
					teams = append(teams, game.Teams.Away.Team)
				}
			}
			//Home
			if isNotNil(game.Teams.Home) {
				if isNotNil(game.Teams.Home.Team) {
					teams = append(teams, game.Teams.Home.Team)
				}
			}
		}
	}
	return teams
}

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
* getIndex: returns the index of the value in an array
 */
func indexOfGame(s []*scheduler.Game, x *scheduler.Game) int {
	for i := range s {
		if *s[i].GamePk == *x.GamePk {
			return i
		}
	}
	return -1
}

func GetIndexOfGames(s []*scheduler.Game, x []*scheduler.Game) []int {
	var indexArray []int
	for i := range x {
		indexArray = append(indexArray, indexOfGame(s, x[i]))
	}
	return indexArray
}

func deleteGame(s []*scheduler.Game, i int) []*scheduler.Game {
	if i == len(s)-1 {
		return s[:i]
	}
	return append(s[:i], s[i+1:]...)
}

func rearrangeGameSlice(s []*scheduler.Game, val scheduler.Game, gameIndex int) []*scheduler.Game {
	var temp *scheduler.Game
	s = append(s, temp)
	copy(s[1:], s[:])
	s[0] = &val
	s = deleteGame(s, gameIndex+1)
	return s
}

func recursiveRearrange(s []*scheduler.Game, favGamesIndex []int) []*scheduler.Game {
	if len(favGamesIndex) == 0 {
		return s
	}
	index := len(favGamesIndex) - 1
	gameIndex := favGamesIndex[index]
	s = rearrangeGameSlice(s, *s[gameIndex], gameIndex)
	return recursiveRearrange(s, incrementByOne(favGamesIndex[:index]))
}

func incrementByOne(slice []int) []int {
	for i := range slice {
		slice[i]++
	}
	return slice
}

/*
* CreateCustomIndex: returns the custom array of indices for fav games.
* Description: This function will handle the double header situation & finds the earlier & later games in a given day.
* 			   Also, It will take care of live game constrains & sort the games based future or past events
*		  doubleheader: Y (single admission)
*				second/later game is one with startTImeTBD value true
*         doubleheader: S (split admission)
*               decide earlier & later games based on gameDate
*
*
 */
func CreateCustomIndex(ctx context.Context, favIndices []int, favGames []*scheduler.Game) ([]int, error) {
	var earlier, later *scheduler.Game
	if len(favGames) != 2 {
		/* if we got less than two games or more than two games, just return the order as it is. we are not handling more
		than two games played by one team in this API version. */
		return favIndices, nil
	}
	g1 := favGames[0]
	g2 := favGames[1]
	// if it is single admission both games will gave doubleheader "Y"
	if isSingleAdmission(ctx, g1) || isSingleAdmission(ctx, g2) {
		if getStartTImeTBD(ctx, g1) && !getStartTImeTBD(ctx, g2) {
			later = g1
			earlier = g2
		}
		later = g2
		earlier = g1
	} else {
		g1Time, err := getGameDate(ctx, g1)
		if err != nil {
			return nil, err
		}
		g2Time, err := getGameDate(ctx, g2)
		if err != nil {
			return nil, err
		}
		fmt.Println(g1Time, g2Time)
	}

	fmt.Println(earlier, later)
	fmt.Println(earlier, later)

	return favIndices, nil
}

func getDoubleheader(ctx context.Context, game *scheduler.Game) *string {
	if isNotNil(game) {
		if isNotNil(game.DoubleHeader) {
			return game.DoubleHeader
		}
	}
	return nil
}

func isSingleAdmission(ctx context.Context, game *scheduler.Game) bool {
	if s := getDoubleheader(ctx, game); s != nil {
		return *s == "Y"
	}
	return false
}

func getStartTImeTBD(ctx context.Context, game *scheduler.Game) bool {
	if isNotNil(game) {
		if isNotNil(game.Status) {
			if isNotNil(game.Status.StartTimeTBD) {
				return *game.Status.StartTimeTBD
			}
		}
	}
	return false
}

func getGameDate(ctx context.Context, game *scheduler.Game) (time.Time, error) {
	var tm time.Time
	if isNotNil(game) {
		if isNotNil(game.GameDate) {
			tm, err := time.Parse(string(design.DateTimeFormat), *game.GameDate)
			if err != nil {
				return tm, errors.New(fmt.Sprintf("error parsing time format: %s", *game.GameDate))
			}
			return tm, nil
		}
	}
	return tm, errors.New("seems date value is null")
}
