package front_service

import (
	"context"
	"encoding/json"

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
