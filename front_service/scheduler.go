package front_service

import (
	"context"
	"log"

	scheduler "github.com/mlb/mlb-ballpark-segregation-service/front_service/gen/scheduler"
	stats "github.com/mlb/mlb-ballpark-segregation-service/front_service/services/statsAPI"
)

// Scheduler service example implementation.
// The example methods log the requests and return zero values.
type schedulersrvc struct {
	logger *log.Logger
}

// checks weather all endpoints are implemented or not
var (
	_ scheduler.Service = (*schedulersrvc)(nil)
)

// NewScheduler returns the Scheduler service implementation.
func NewScheduler(logger *log.Logger) scheduler.Service {
	return &schedulersrvc{logger}
}

// Retrieves a schedule of games
func (s *schedulersrvc) Index(ctx context.Context, p *scheduler.IndexPayload) (res *scheduler.Schedule, err error) {
	res = &scheduler.Schedule{}

	statsC, err := stats.NewStatsClient(ctx)
	if err != nil {
		s.logger.Print("error", err)
		return nil, err
	}

	res, err = statsC.GetStatsAPISchedule(ctx, p.Date)
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
	//s.logger.Printf(prettyPrint(favTeamGames))

	if len(favTeamGames) == 0 {
		// there is no games for our fav team on given day, so return response as it is.
		s.logger.Print("there is no games for our fav team on given day", len(favTeamGames))
		return res, err
	}

	if len(favTeamGames) < 2 {
		index := indexOfGame(games, favTeamGames[0])
		in := []int{index}
		games = recursiveRearrange(games, in)
	}

	res.Dates[0].Games = games

	s.logger.Print("favGames len", len(favTeamGames))

	s.logger.Print("scheduler.index")
	return
}
