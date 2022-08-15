package front_service

import (
	"context"
	"errors"
	"fmt"
	"time"

	scheduler "github.com/mlb/mlb-ballpark-segregation-service/front_service/gen/scheduler"
)

// useful constant values in sorting logic
const (
	DateFormat       = "2006-01-02"
	DateTimeFormat   = time.RFC3339
	AvgGameDuration  = 3
	TimeFramePast    = "PAST"
	TimeFrameFuture  = "FUTURE"
	TimeFramePresent = "PRESENT"
	EarlierGame      = "EARLIER"
	LaterGame        = "LATER"
)

// useful sorting rules to sort the games in past & future dates
var (
	SortingRules = map[string][]string{
		TimeFramePast:   {LaterGame, EarlierGame},
		TimeFrameFuture: {EarlierGame, LaterGame},
	}
)

/*
* GetFavoriteTeamGames
* Description: returns the favorite games of a team, if team not found. It returns nil
 */
func (s *schedulersrvc) GetFavoriteTeamGames(ctx context.Context, teamID uint, games []*scheduler.Game) []*scheduler.Game {
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

/*
* extractTeams
* Description: returns the home & away teams of a game, if team not found. It returns nil
 */
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

/*
* rearrangeGameSlice
* Description: it rearranges the array, such that given input value sorts to zero index
*              Input:  [2 5 7 1 8 65 35](array) 8(val) 4(index)
*              Output: [8 2 5 7 1 65 35]
 */
func rearrangeGameSlice(s []*scheduler.Game, val scheduler.Game, gameIndex int) []*scheduler.Game {
	var temp *scheduler.Game
	s = append(s, temp) //Increase the size of array by 1
	copy(s[1:], s[:])   // move the array by index 1 & make space at 0 index
	s[0] = &val
	s = deleteGame(s, gameIndex+1) // delete the val at original place
	return s
}

/*
* recursiveRearrange
* Description: it rearranges the array, such that all values of given list of indices sorts to first
*              Input:  [2  5   7  1  8  65  35](array) [4,6](index array)
*               Index   0  1   2  3  4   5   6
*
*              Output: [8  35  2  5  7   1  65]
 */
func recursiveRearrange(s []*scheduler.Game, favGamesIndex []int) []*scheduler.Game {
	if len(favGamesIndex) == 0 {
		return s
	}
	index := len(favGamesIndex) - 1
	gameIndex := favGamesIndex[index]
	s = rearrangeGameSlice(s, *s[gameIndex], gameIndex)
	// if we deleted element is on left side of next deleting element, no need to increment, if you deleted on right side on next deleting element
	//  increment every value in favGamesIndex by one to make sure, we are rearranging correct value
	if len(favGamesIndex) > 1 {
		if gameIndex > favGamesIndex[index-1] {
			favGamesIndex = incrementByOne(favGamesIndex)
		}
	}
	return recursiveRearrange(s, favGamesIndex[:index])
}

/*
* CreateCustomIndex: returns the custom array of indices for fav games.
* Description: This function will handle the double header situation & finds the earlier & later games in a given day.
* 			   Also, It will take care of live game constrains & sort the games based future or past events
*		  doubleheader: Y (single admission)
*				second game is one with startTImeTBD value true
*         doubleheader: S (split admission)
*               decide first & second games based on gameDate
*         if current date and game date matches, check for live or do future & past sorting
 */
func (s *schedulersrvc) CreateCustomIndex(ctx context.Context, date string, favGames []*scheduler.Game) ([]*scheduler.Game, error) {
	var timeFrameGameMap map[string]*scheduler.Game
	var res []*scheduler.Game
	if len(favGames) != 2 {
		/* if we got less than two games or more than two games, just return the order as it is. we are not handling more
		than two games played by one team in this API version. */
		s.logger.Printf("not handling more than two games in this API version, no of games: %d", len(favGames))
		return favGames, nil
	}
	// get the games from array of size 2
	g1 := favGames[0]
	g2 := favGames[1]

	// get timeFrameMap for a give day
	timeFrameGameMap, err := s.GetTimeFrameGameMap(ctx, g1, g2)
	if err != nil {
		s.logger.Printf("error while getting time frame game map. err: %s", err.Error())
		return nil, err
	}
	// get the time frame to decide live,past & future
	timeFrame, err := getTimeFrame(date)
	if err != nil {
		s.logger.Printf("error while getting the time frame. err: %s", err.Error())
		return nil, err
	}
	// get the time frame
	switch timeFrame {
	case TimeFramePresent:
		res, err = s.RearrangeGamesByLiveStatus(ctx, g1, g2)
		if err != nil {
			s.logger.Printf("error while getting the RearrangeGamesByLiveStatus. err: %s", err.Error())
			return nil, err
		}
		return res, nil
	// sort the games based on time sorting rules
	case TimeFramePast, TimeFrameFuture:
		sortingRule := SortingRules[timeFrame]
		for _, v := range sortingRule {
			res = append(res, timeFrameGameMap[v])
		}
		return res, nil
	default:
		s.logger.Print("something went wrong")
	}
	return favGames, nil
}

/*
* getDoubleheader
* Description: returns the game double header value
 */
func getDoubleheader(ctx context.Context, game *scheduler.Game) *string {
	if isNotNil(game) {
		if isNotNil(game.DoubleHeader) {
			return game.DoubleHeader
		}
	}
	return nil
}

/*
* isSingleAdmission
* Description: returns true if a game is single admission
 */
func isSingleAdmission(ctx context.Context, game *scheduler.Game) bool {
	if s := getDoubleheader(ctx, game); s != nil {
		return *s == "Y"
	}
	return false
}

/*
* getGameDate
* Description: returns game Date of a game in Time format
 */
func (s *schedulersrvc) getGameDate(ctx context.Context, game *scheduler.Game) (time.Time, error) {
	var tm time.Time
	if isNotNil(game) {
		if isNotNil(game.GameDate) {
			tm, err := time.Parse(string(DateTimeFormat), *game.GameDate)
			if err != nil {
				s.logger.Printf("error while parsing time: %s", *game.GameDate)
				return tm, fmt.Errorf("error parsing time format: %s", *game.GameDate)
			}
			return tm, nil
		}
	}
	s.logger.Print("date value is null")
	return tm, errors.New("seems date value is null")
}

/*
* getTimeFrame
* Description: returns weather the game is in past, future or current day
 */
func getTimeFrame(date string) (string, error) {
	t1, err := time.Parse(string(DateFormat), time.Now().Format(DateFormat))
	if err != nil {
		return "", err
	}
	t2, err := parseDate(date)
	if err != nil {
		return "", err
	}
	// if times are equal, return present
	if t2.Equal(t1) {
		return TimeFramePresent, nil
	}
	// check weather given date is past event or not
	if t2.Before(t1) {
		return TimeFramePast, nil
	}
	return TimeFrameFuture, nil
}

/*
* GetTimeFrameGameMap
* Description: returns the map of earlier & later games.
*            Ex:- timeFrameMap : {earlier: "g1", "later": g2}
 */
func (s *schedulersrvc) GetTimeFrameGameMap(ctx context.Context, g1 *scheduler.Game, g2 *scheduler.Game) (map[string]*scheduler.Game, error) {
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
	timeFrameGameMap, err := s.RearrangeGamesByGameDate(ctx, g1, g2)
	if err != nil {
		s.logger.Printf("error while getting timeFrameMap, err: %s", err.Error())
		return nil, err
	}
	return timeFrameGameMap, nil
}

/*
* RearrangeGamesByGameDate
* Description: returns the array of games in chronology order
 */
func (s *schedulersrvc) RearrangeGamesByGameDate(ctx context.Context, g1 *scheduler.Game, g2 *scheduler.Game) (map[string]*scheduler.Game, error) {
	timeFrameGameMap := make(map[string]*scheduler.Game)
	//even though if is it not split admission, we need to decide earlier & later games using gameDate
	g1Time, err := s.getGameDate(ctx, g1)
	if err != nil {
		return nil, err
	}
	g2Time, err := s.getGameDate(ctx, g2)
	if err != nil {
		return nil, err
	}
	// is game1 start time before game 2
	if g1Time.Before(g2Time) {
		timeFrameGameMap[EarlierGame] = g1
		timeFrameGameMap[LaterGame] = g2
	} else {
		timeFrameGameMap[EarlierGame] = g2
		timeFrameGameMap[LaterGame] = g1
	}
	return timeFrameGameMap, nil
}

/*
* RearrangeGamesByLiveStatus
* Description: returns the array of games in chronology order or live game first
 */
func (s *schedulersrvc) RearrangeGamesByLiveStatus(ctx context.Context, g1 *scheduler.Game, g2 *scheduler.Game) ([]*scheduler.Game, error) {
	res := make([]*scheduler.Game, 0)
	// get the current time & make a window with avg time of a baseball game of 4 Hours
	// so, if current time falls between game start time & window time. the game is live
	currentTime := time.Now()
	window := time.Now().Add(AvgGameDuration * time.Hour)
	g1time, err := s.getGameDate(ctx, g1)
	if err != nil {
		s.logger.Printf("error parsing game date (RearrangeGamesByLiveStatus): %s", err.Error())
		return nil, err
	}
	g2time, err := s.getGameDate(ctx, g2)
	if err != nil {
		s.logger.Printf("error parsing game date (RearrangeGamesByLiveStatus): %s", err.Error())
		return nil, err
	}
	if currentTime.After(g1time) && g1time.Before(window) && !g2time.Before(currentTime) {
		res = append(res, g1, g2)
		return res, nil
	}
	if currentTime.After(g2time) && g2time.Before(window) && !g2time.Before(g1time) {
		res = append(res, g2, g1)
		return res, nil
	}
	// if both games are not live, arrange them in chronology
	if g1time.Before(g2time) {
		res = append(res, g1, g2)
		return res, nil
	} else {
		res = append(res, g2, g1)
		return res, nil
	}
}

/*
* getStartTImeTBD
* Description: returns StartTImeTBD of a game
 */
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
