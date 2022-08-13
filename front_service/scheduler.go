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
		// len of nil array is 0
		if len(res.Dates) <= 0 {
			s.logger.Printf("No games scheduled on this date: %s", p.Date)

		}

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

	favTeamGames := s.getFavoriteTeamGames(ctx, p.ID, games)
	if len(favTeamGames) == 0 {
		// there is no games for our fav team on given day, so return response as it is.
		s.logger.Print("there is no games for our fav team on given day", len(favTeamGames))
		return res, nil
	}

	favIndices := GetIndexOfGames(games, favTeamGames)

	// we are sure one team plays only two games for a day
	if len(favTeamGames) == 2 {
		// we have to deal with doubleheader situation & rearrange the index in custom order.
		favIndices, err = CreateCustomIndex(ctx, favIndices, favTeamGames)
		if err != nil {
			return
		}
	}
	s.logger.Printf("the indices of the fav games: %v", favIndices)
	games = recursiveRearrange(games, favIndices)
	res.Dates[0].Games = games
	s.logger.Print("favGames len", len(favTeamGames))
	return
}
