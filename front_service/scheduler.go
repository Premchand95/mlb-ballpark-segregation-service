package front_service

import (
	"context"
	"errors"
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
func (s *schedulersrvc) Index(ctx context.Context, p *scheduler.IndexPayload) (res *scheduler.Schedule, err error) {
	res = &scheduler.Schedule{}

	res, err = s.statsC.GetStatsAPISchedule(ctx, p.Date)
	if err != nil {
		s.logger.Print("error", err)
		return nil, err
	}

	var games []*scheduler.Game

	if isNotNil(res) {
		// dates is a required field, so no need to do nil check
		if len(res.Dates) > 0 {
			// As we query stats API scheduler with only one date, for now, we assume we only get one date
			if len(res.Dates[0].Games) > 0 {
				games = append(games, res.Dates[0].Games...)
			}
		}
	}

	favTeamGames := s.getFavoriteTeamGames(ctx, p.ID, games)
	if len(favTeamGames) == 0 {
		// there is no games for our fav team on given day, so return response as it is.
		s.logger.Print("there is no games for our fav team on given day", len(favTeamGames))
		return res, err
	}

	favIndices := GetIndexOfGames(games, favTeamGames)

	s.logger.Printf("the indices of the fav games: %v", favIndices)

	games = recursiveRearrange(games, favIndices)

	res.Dates[0].Games = games

	s.logger.Print("favGames len", len(favTeamGames))

	s.logger.Print("scheduler.index")
	return
}
