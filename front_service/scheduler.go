package front_service

import (
	"context"
	"errors"
	"fmt"
	"log"

	scheduler "github.com/mlb/mlb-ballpark-segregation-service/front_service/gen/scheduler"
	stats "github.com/mlb/mlb-ballpark-segregation-service/front_service/services/statsAPI"
)

// Scheduler service example implementation.
// The example methods log the requests and return zero values.
type schedulersrvc struct {
	logger log.Logger
	statsC stats.Client
}

type SchedulerParams struct {
	StatsC stats.Client
	Logger log.Logger
}

// checks weather all endpoints are implemented or not
var (
	_ scheduler.Service = (*schedulersrvc)(nil)
)

const (
	internalErrMessage = "please contact administrator for more details"
)

// NewScheduler returns the Scheduler service implementation.
func NewScheduler(params *SchedulerParams) (scheduler.Service, error) {
	if params.StatsC == nil {
		return nil, errors.New("all parameters are required for instantiating the scheduler client")
	}
	return &schedulersrvc{
		logger: params.Logger,
		statsC: params.StatsC,
	}, nil
}

// Retrieves a schedule of games
func (s *schedulersrvc) Index(ctx context.Context, p *scheduler.IndexPayload) (*scheduler.Schedule, error) {
	// function parameters
	var (
		res   *scheduler.Schedule
		err   error
		games []*scheduler.Game
	)

	// call stats API to get raw schedule
	res, err = s.statsC.GetStatsAPISchedule(ctx, p.Date)
	if err != nil {
		s.logger.Print("error while making call to stats API: ", err)
		return nil, scheduler.MakeInternalError(fmt.Errorf(internalErrMessage))
	}

	if isNotNil(res) {
		// len of nil array is 0
		if len(res.Dates) <= 0 {
			s.logger.Printf("No games scheduled on this date: %s", p.Date)
			return res, nil
		}

		// checking date to avoid NIL segmentation errors
		var date scheduler.Date
		for _, d := range res.Dates {
			if isNotNil(d.Date) {
				if *d.Date == p.Date {
					date = *d
				}
			}
		}
		games = append(games, date.Games...)
	}
	// get the favorite games for the given team
	favTeamGames := s.GetFavoriteTeamGames(ctx, p.ID, games)
	if len(favTeamGames) == 0 {
		// there is no games for our fav team on given day, so return response as it is.
		s.logger.Print("there is no games for our fav team on given day ", len(favTeamGames))
		return res, nil
	}

	// we are sure one team plays only two games for a day
	if len(favTeamGames) == 2 {
		// we have to deal with doubleheader situation & rearrange the index in custom order.
		// this function also handles past, future dates & rearrange the games in custom order
		favTeamGames, err = s.CreateCustomIndex(ctx, p.Date, favTeamGames)
		if err != nil {
			return nil, scheduler.MakeInternalError(fmt.Errorf(internalErrMessage))
		}
	}
	// get the indices of the games
	favIndices := GetIndexOfGames(games, favTeamGames)
	s.logger.Printf("the indices of the fav games: %v", favIndices)
	// rearrange the order of games array
	games = recursiveRearrange(games, favIndices)
	// attach the final custom order of games to result
	res.Dates[0].Games = games
	return res, nil
}
