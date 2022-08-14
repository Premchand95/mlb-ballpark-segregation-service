package front_service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	scheduler "github.com/mlb/mlb-ballpark-segregation-service/front_service/gen/scheduler"
)

const (
	DateFormat       = "2006-01-02"
	DateTimeFormat   = time.RFC3339
	AvgGameDuration  = 4
	TimeFramePast    = "PAST"
	TimeFrameFuture  = "FUTURE"
	TimeFramePresent = "PRESENT"
	EarlierGame      = "EARLIER"
	LaterGame        = "LATER"
)

var (
	SortingRules = map[string][]string{
		TimeFramePast:   []string{LaterGame, EarlierGame},
		TimeFrameFuture: []string{EarlierGame, LaterGame},
	}
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
func CreateCustomIndex(ctx context.Context, date string, favGames []*scheduler.Game) ([]*scheduler.Game, error) {
	var timeFrameGameMap map[string]*scheduler.Game
	if len(favGames) != 2 {
		/* if we got less than two games or more than two games, just return the order as it is. we are not handling more
		than two games played by one team in this API version. */
		return favGames, nil
	}
	g1 := favGames[0]
	g2 := favGames[1]

	timeFrameGameMap, err := GetTimeFrameGameMap(ctx, g1, g2)
	if err != nil {
		return favGames, err
	}
	// currentTime := time.Now()
	// window := time.Now().Add(AvgGameDuration * time.Hour)
	var res []*scheduler.Game
	timeFrame, err := getTimeFrame(date)
	if err != nil {
		return favGames, err
	}
	switch timeFrame {
	case TimeFramePresent:
		fmt.Println(timeFrame, "return the live game")
	case TimeFramePast, TimeFrameFuture:
		sortingRule := SortingRules[timeFrame]
		for _, v := range sortingRule {
			res = append(res, timeFrameGameMap[v])
		}
		return res, nil
	default:
		fmt.Println("something went wrong")
	}
	return favGames, nil
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
			tm, err := time.Parse(string(DateTimeFormat), *game.GameDate)
			if err != nil {
				return tm, errors.New(fmt.Sprintf("error parsing time format: %s", *game.GameDate))
			}
			return tm, nil
		}
	}
	return tm, errors.New("seems date value is null")
}

func getTimeFrame(date string) (string, error) {
	t1, err := time.Parse(string(DateFormat), time.Now().Format(DateFormat))
	if err != nil {
		return "", err
	}
	t2, err := parseDate(date)
	if err != nil {
		return "", err
	}
	if t2.Equal(t1) {
		return TimeFramePresent, nil
	}
	if t2.Before(t1) {
		return TimeFrameFuture, nil
	}
	return TimeFramePast, nil
}

func parseDate(date string) (time.Time, error) {
	tm, err := time.Parse(string(DateFormat), date)
	if err != nil {
		return tm, errors.New(fmt.Sprintf("error parsing date format: %s", date))
	}
	return tm, nil
}

func GetTimeFrameGameMap(ctx context.Context, g1 *scheduler.Game, g2 *scheduler.Game) (map[string]*scheduler.Game, error) {
	timeFrameGameMap := make(map[string]*scheduler.Game)
	// if it is single admission both games will gave doubleheader "Y"
	if isSingleAdmission(ctx, g1) || isSingleAdmission(ctx, g2) {
		if getStartTImeTBD(ctx, g1) && !getStartTImeTBD(ctx, g2) {
			timeFrameGameMap[LaterGame] = g1
			timeFrameGameMap[EarlierGame] = g2
		}
		timeFrameGameMap[LaterGame] = g2
		timeFrameGameMap[EarlierGame] = g1
		return timeFrameGameMap, nil
	}
	//even though if is it not split admission, we need to decide earlier & later games using gameDate
	timeFrameGameMap, err := RearrangeGamesByGameDate(ctx, g1, g2)
	if err != nil {
		return nil, err
	}
	return timeFrameGameMap, nil
}

func RearrangeGamesByGameDate(ctx context.Context, g1 *scheduler.Game, g2 *scheduler.Game) (map[string]*scheduler.Game, error) {
	timeFrameGameMap := make(map[string]*scheduler.Game)
	//even though if is it not split admission, we need to decide earlier & later games using gameDate
	g1Time, err := getGameDate(ctx, g1)
	if err != nil {
		return nil, err
	}
	g2Time, err := getGameDate(ctx, g2)
	if err != nil {
		return nil, err
	}
	if g1Time.Before(g2Time) {
		timeFrameGameMap[EarlierGame] = g1
		timeFrameGameMap[LaterGame] = g2
	} else {
		timeFrameGameMap[EarlierGame] = g2
		timeFrameGameMap[LaterGame] = g1
	}
	return timeFrameGameMap, nil
}
