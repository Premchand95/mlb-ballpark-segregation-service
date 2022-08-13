package front_service_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	frontSvc "github.com/mlb/mlb-ballpark-segregation-service/front_service"
	"github.com/mlb/mlb-ballpark-segregation-service/front_service/design"
	endpoint "github.com/mlb/mlb-ballpark-segregation-service/front_service/gen/scheduler"
	mock_statsapi "github.com/mlb/mlb-ballpark-segregation-service/front_service/mocks/mlb_statsAPI"
)

type testSetupSchedulerObjs struct {
	statsAPIClientMock *mock_statsapi.MockClient
	//requestsClientMock *mock_requests.MockClientAPI
}

func createTestSetupSchedulerObjs(ctrl *gomock.Controller) testSetupSchedulerObjs {
	statsAPIClientMock := mock_statsapi.NewMockClient(ctrl)
	//requestsClientMock := mock_requests.NewMockClientAPI(ctrl)

	return testSetupSchedulerObjs{
		statsAPIClientMock: statsAPIClientMock,
		//requestsClientMock: &requestsClientMock,
	}
}

func returnPointer[T any](i T) *T {
	return &i
}

func TestIndexSchedulerBasic(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testSetup := createTestSetupSchedulerObjs(ctrl)
	testPayload := endpoint.IndexPayload{
		ID:   141,
		Date: "2021-09-11",
	}

	var logger *log.Logger = log.New(os.Stderr, fmt.Sprintf("[%s]- unit testing ", design.ServiceName), log.Ltime)

	testObject := createBasicSchedulerObject()

	testSetup.statsAPIClientMock.EXPECT().GetStatsAPISchedule(gomock.Any(), gomock.Any()).Return(&testObject, nil)

	svc, err := frontSvc.NewScheduler(&frontSvc.SchedulerParams{
		StatsC: testSetup.statsAPIClientMock,
		Logger: *logger,
	})
	if err != nil {
		t.Errorf("failed to create new scheduler: %v", err)
	}

	res, err := svc.Index(context.Background(), &testPayload)
	if err != nil {
		t.Errorf("failed to call index endpoint: %v", err)
	}

	if res == nil {
		t.Errorf("failed to retrieve schedule")
	}
}

func createBasicSchedulerObject() endpoint.Schedule {
	return endpoint.Schedule{
		Copyright:            returnPointer("Copyright 2022 MLB Advanced Media"),
		TotalItems:           returnPointer(uint(2)),
		TotalEvents:          returnPointer(uint(0)),
		TotalGames:           returnPointer(uint(2)),
		TotalGamesInProgress: returnPointer(uint(0)),
		Dates: []*endpoint.Date{
			{
				Date:                 returnPointer("2021-09-11"),
				TotalItems:           returnPointer(uint(2)),
				TotalEvents:          returnPointer(uint(0)),
				TotalGames:           returnPointer(uint(2)),
				TotalGamesInProgress: returnPointer(uint(0)),
				Games: []*endpoint.Game{
					{
						GamePk:       returnPointer(uint64(632539)),
						Link:         returnPointer("/api/v1.1/game/632539/feed/live"),
						GameType:     returnPointer("R"),
						Season:       returnPointer("2021"),
						GameDate:     returnPointer("2021-09-11T20:35:00Z"),
						OfficialDate: returnPointer("2021-09-11"),
						Status: &endpoint.Status{
							AbstractGameState: returnPointer("Final"),
							CodedGameState:    returnPointer("F"),
							DetailedState:     returnPointer("Final"),
							StatusCode:        returnPointer("F"),
							StartTimeTBD:      returnPointer(false),
							AbstractGameCode:  returnPointer("F"),
						},
						Teams: &endpoint.Teams{
							Away: &endpoint.TeamInfo{
								Score:        returnPointer(uint(11)),
								IsWinner:     returnPointer(true),
								SplitSquad:   returnPointer(false),
								SeriesNumber: returnPointer(uint(46)),
								LeagueRecord: &endpoint.LeagueRecord{
									Wins:   returnPointer(uint(78)),
									Losses: returnPointer(uint(63)),
									Pct:    returnPointer(".553"),
								},
								Team: &endpoint.Team{
									ID:   returnPointer(uint(141)),
									Name: returnPointer("Toronto Blue Jays"),
									Link: returnPointer("api/v1/teams/141"),
								},
							},
							Home: &endpoint.TeamInfo{
								Score:        returnPointer(uint(10)),
								IsWinner:     returnPointer(false),
								SplitSquad:   returnPointer(false),
								SeriesNumber: returnPointer(uint(46)),
								LeagueRecord: &endpoint.LeagueRecord{
									Wins:   returnPointer(uint(46)),
									Losses: returnPointer(uint(95)),
									Pct:    returnPointer(".326"),
								},
								Team: &endpoint.Team{
									ID:   returnPointer(uint(110)),
									Name: returnPointer("Baltimore Orioles"),
									Link: returnPointer("api/v1/teams/110"),
								},
							},
						},
						Venue: &endpoint.Venue{
							ID:   returnPointer(uint(2)),
							Name: returnPointer("Oriole Park at camden Yards"),
							Link: returnPointer("/api/v1/venues/2"),
						},
						Content: &endpoint.Content{
							Link: returnPointer("/api/v1/game/632539/content"),
						},
						IsTie:                  returnPointer(false),
						GameNumber:             returnPointer(uint(1)),
						PublicFacing:           returnPointer(true),
						DoubleHeader:           returnPointer("Y"),
						GamedayType:            returnPointer("P"),
						Tiebreaker:             returnPointer("N"),
						CalendarEventID:        returnPointer("14-632539-2021-09-11"),
						SeasonDisplay:          returnPointer("2021"),
						DayNight:               returnPointer("night"),
						ScheduledInnings:       returnPointer(uint(7)),
						ReverseHomeAwayStatus:  returnPointer(false),
						InningBreakLength:      returnPointer(uint(120)),
						GamesInSeries:          returnPointer(uint(4)),
						SeriesGameNumber:       returnPointer(uint(2)),
						SeriesDescription:      returnPointer("Regular Season"),
						RecordSource:           returnPointer("S"),
						IfNecessary:            returnPointer("N"),
						IfNecessaryDescription: returnPointer("Normal Game"),
					},
					{
						GamePk:       returnPointer(uint64(633331)),
						Link:         returnPointer("/api/v1.1/game/633331/feed/live"),
						GameType:     returnPointer("R"),
						Season:       returnPointer("2021"),
						GameDate:     returnPointer("2021-09-11T20:40:00Z"),
						OfficialDate: returnPointer("2021-09-11"),
						Status: &endpoint.Status{
							AbstractGameState: returnPointer("Final"),
							CodedGameState:    returnPointer("F"),
							DetailedState:     returnPointer("Final"),
							StatusCode:        returnPointer("F"),
							StartTimeTBD:      returnPointer(true),
							AbstractGameCode:  returnPointer("F"),
						},
						Teams: &endpoint.Teams{
							Away: &endpoint.TeamInfo{
								Score:        returnPointer(uint(11)),
								IsWinner:     returnPointer(true),
								SplitSquad:   returnPointer(false),
								SeriesNumber: returnPointer(uint(46)),
								LeagueRecord: &endpoint.LeagueRecord{
									Wins:   returnPointer(uint(79)),
									Losses: returnPointer(uint(63)),
									Pct:    returnPointer(".556"),
								},
								Team: &endpoint.Team{
									ID:   returnPointer(uint(141)),
									Name: returnPointer("Toronto Blue Jays"),
									Link: returnPointer("api/v1/teams/141"),
								},
							},
							Home: &endpoint.TeamInfo{
								Score:        returnPointer(uint(2)),
								IsWinner:     returnPointer(false),
								SplitSquad:   returnPointer(false),
								SeriesNumber: returnPointer(uint(46)),
								LeagueRecord: &endpoint.LeagueRecord{
									Wins:   returnPointer(uint(46)),
									Losses: returnPointer(uint(96)),
									Pct:    returnPointer(".324"),
								},
								Team: &endpoint.Team{
									ID:   returnPointer(uint(110)),
									Name: returnPointer("Baltimore Orioles"),
									Link: returnPointer("api/v1/teams/110"),
								},
							},
						},
						Venue: &endpoint.Venue{
							ID:   returnPointer(uint(2)),
							Name: returnPointer("Oriole Park at camden Yards"),
							Link: returnPointer("/api/v1/venues/2"),
						},
						Content: &endpoint.Content{
							Link: returnPointer("/api/v1/game/633331/content"),
						},
						IsTie:                  returnPointer(false),
						GameNumber:             returnPointer(uint(2)),
						PublicFacing:           returnPointer(true),
						DoubleHeader:           returnPointer("Y"),
						GamedayType:            returnPointer("P"),
						Tiebreaker:             returnPointer("N"),
						CalendarEventID:        returnPointer("14-633331-2021-09-11"),
						SeasonDisplay:          returnPointer("2021"),
						DayNight:               returnPointer("night"),
						ScheduledInnings:       returnPointer(uint(7)),
						ReverseHomeAwayStatus:  returnPointer(false),
						InningBreakLength:      returnPointer(uint(120)),
						GamesInSeries:          returnPointer(uint(4)),
						SeriesGameNumber:       returnPointer(uint(3)),
						SeriesDescription:      returnPointer("Regular Season"),
						RecordSource:           returnPointer("S"),
						IfNecessary:            returnPointer("N"),
						IfNecessaryDescription: returnPointer("Normal Game"),
					},
				},
				Events: []interface{}{},
			},
		},
	}
}
