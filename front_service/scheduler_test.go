package front_service_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

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

/*
* TestNoFavTeam
* Description: No fav team in the raw schedule
 */
func TestNoFavTeam(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testSetup := createTestSetupSchedulerObjs(ctrl)
	testPayload := endpoint.IndexPayload{
		ID:   54354,
		Date: "2022-08-13",
	}

	var logger *log.Logger = log.New(os.Stderr, fmt.Sprintf("[%s]- unit testing ", design.ServiceName), log.Ltime)

	testObject := createBasicSchedulerObject()
	testObject.Dates[0].Date = &testPayload.Date

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

/*
* TestFavTeamOnlyOneGame
* Description: fav team playing only one game
 */
func TestFavTeamOnlyOneGame(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testSetup := createTestSetupSchedulerObjs(ctrl)
	testPayload := endpoint.IndexPayload{
		ID:   121,
		Date: "2022-08-13",
	}

	var logger *log.Logger = log.New(os.Stderr, fmt.Sprintf("[%s]- unit testing ", design.ServiceName), log.Ltime)

	testObject := createBasicSchedulerObject()
	testObject.Dates[0].Date = &testPayload.Date

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
	// checking first index game PK
	if *res.Dates[0].Games[0].GamePk != uint64(234) {
		t.Errorf("we didn't got our fab game on top")
	}
}

/*
* TestPastGamesOrderSingleAdmission
* Description: Past games order should be chronological. second game, first game
 */
func TestPastGamesOrderSingleAdmission(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testSetup := createTestSetupSchedulerObjs(ctrl)
	testPayload := endpoint.IndexPayload{
		ID:   11,
		Date: "2022-08-13", // past date
	}

	var logger *log.Logger = log.New(os.Stderr, fmt.Sprintf("[%s]- unit testing ", design.ServiceName), log.Ltime)

	testObject := createBasicSchedulerObject()
	testObject.Dates[0].Date = &testPayload.Date
	testStartTimeTBD := true
	testObject.Dates[0].Games[3].Status.StartTimeTBD = &testStartTimeTBD

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
	//second game should be first
	if *res.Dates[0].Games[0].Status.StartTimeTBD != testStartTimeTBD {
		t.Errorf("failed to test past single admission game")
	}
}

/*
* TestPastGamesOrderSplitAdmission
* Description: Past games order should be chronological. second game, first game
 */
func TestPastGamesOrderSplitAdmission(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testSetup := createTestSetupSchedulerObjs(ctrl)
	testPayload := endpoint.IndexPayload{
		ID:   11,
		Date: "2022-08-13", // past date
	}

	var logger *log.Logger = log.New(os.Stderr, fmt.Sprintf("[%s]- unit testing ", design.ServiceName), log.Ltime)

	testObject := createBasicSchedulerObject()
	testObject.Dates[0].Date = &testPayload.Date
	testObject.Dates[0].Games[1].DoubleHeader = returnPointer("S")
	testObject.Dates[0].Games[3].DoubleHeader = returnPointer("S")
	testObject.Dates[0].Games[1].GameDate = returnPointer("2022-08-13T09:35:00Z") // 9AM
	testObject.Dates[0].Games[3].GameDate = returnPointer("2022-08-13T20:35:00Z") // 8 PM

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
	//second game should be first
	if *res.Dates[0].Games[0].GamePk != uint64(3) {
		t.Errorf("failed to test past split admission game")
	}
}

/*
* TestFutureGamesOrderSingleAdmission
* Description: Future games order should be earlier game first. first game, second game
 */
func TestFutureGamesOrderSingleAdmission(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testSetup := createTestSetupSchedulerObjs(ctrl)
	testPayload := endpoint.IndexPayload{
		ID:   11,
		Date: "2023-09-13", // past date
	}

	var logger *log.Logger = log.New(os.Stderr, fmt.Sprintf("[%s]- unit testing ", design.ServiceName), log.Ltime)

	testObject := createBasicSchedulerObject()
	testObject.Dates[0].Date = &testPayload.Date
	testStartTimeTBD := true
	testObject.Dates[0].Games[3].Status.StartTimeTBD = &testStartTimeTBD

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
	//earlier game should be first
	if *res.Dates[0].Games[0].Status.StartTimeTBD != false {
		t.Errorf("failed to test past single admission game")
	}
}

/*
* TestFutureGamesOrderSPlitAdmission
* Description: Future games order should be earlier game first. first game, second game
 */
func TestFutureGamesOrderSPlitAdmission(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testSetup := createTestSetupSchedulerObjs(ctrl)
	testPayload := endpoint.IndexPayload{
		ID:   11,
		Date: "2023-08-13", // past date
	}

	var logger *log.Logger = log.New(os.Stderr, fmt.Sprintf("[%s]- unit testing ", design.ServiceName), log.Ltime)

	testObject := createBasicSchedulerObject()
	testObject.Dates[0].Date = &testPayload.Date
	testObject.Dates[0].Games[1].DoubleHeader = returnPointer("S")
	testObject.Dates[0].Games[3].DoubleHeader = returnPointer("S")
	testObject.Dates[0].Games[1].GameDate = returnPointer("2023-08-13T09:35:00Z") // 9AM
	testObject.Dates[0].Games[3].GameDate = returnPointer("2023-08-13T20:35:00Z") // 8 PM

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
	//first game should be top the list
	if *res.Dates[0].Games[0].GamePk != uint64(1) {
		t.Errorf("failed to test past split admission game")
	}
}

/*
* TestLiveGameSituation
* Description: rerun live game first
 */
func TestLiveGameSituation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	today := time.Now().Format(frontSvc.DateFormat)
	testSetup := createTestSetupSchedulerObjs(ctrl)
	testPayload := endpoint.IndexPayload{
		ID:   11,
		Date: today,
	}

	var logger *log.Logger = log.New(os.Stderr, fmt.Sprintf("[%s]- unit testing ", design.ServiceName), log.Ltime)

	testObject := createBasicSchedulerObject()
	testObject.Dates[0].Date = &testPayload.Date
	testObject.Dates[0].Games[1].DoubleHeader = returnPointer("S")
	testObject.Dates[0].Games[3].DoubleHeader = returnPointer("S")
	testDateTimeG1 := time.Now().Add(-5 * time.Minute).Format(frontSvc.DateTimeFormat)
	testDateTimeG2 := time.Now().Add(5 * time.Hour).Format(frontSvc.DateTimeFormat)
	testObject.Dates[0].Games[1].GameDate = returnPointer(testDateTimeG1) // live
	testObject.Dates[0].Games[3].GameDate = returnPointer(testDateTimeG2) // not-live

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
	//live game G1 should be top the list
	if *res.Dates[0].Games[0].GamePk != uint64(1) {
		t.Errorf("failed to test past split admission game")
	}
}

/*
* TestLiveGameSituationG2
* Description: rerun live game first, when second game is live
 */
func TestLiveGameSituationG2(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	today := time.Now().Format(frontSvc.DateFormat)
	testSetup := createTestSetupSchedulerObjs(ctrl)
	testPayload := endpoint.IndexPayload{
		ID:   11,
		Date: today,
	}

	var logger *log.Logger = log.New(os.Stderr, fmt.Sprintf("[%s]- unit testing ", design.ServiceName), log.Ltime)

	testObject := createBasicSchedulerObject()
	testObject.Dates[0].Date = &testPayload.Date
	testObject.Dates[0].Games[1].DoubleHeader = returnPointer("S")
	testObject.Dates[0].Games[3].DoubleHeader = returnPointer("S")
	testDateTimeG1 := time.Now().Add(-5 * time.Hour).Format(frontSvc.DateTimeFormat)
	testDateTimeG2 := time.Now().Add(-10 * time.Minute).Format(frontSvc.DateTimeFormat)
	testObject.Dates[0].Games[1].GameDate = returnPointer(testDateTimeG1) // not-live
	testObject.Dates[0].Games[3].GameDate = returnPointer(testDateTimeG2) // live

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
	//live game G1 should be top the list
	if *res.Dates[0].Games[0].GamePk != uint64(3) {
		t.Errorf("failed to test live game scenario for second game")
	}
}

/*
* returnPointer
* Description: returns pointer for a variable
 */
func returnPointer[T any](i T) *T {
	return &i
}

/*
* createBasicSchedulerObject
* Description: creates a basic scheduler object
 */
func createBasicSchedulerObject() endpoint.Schedule {
	return endpoint.Schedule{
		Copyright:            returnPointer("Copyright 2022 MLB Advanced Media"),
		TotalItems:           returnPointer(uint(2)),
		TotalEvents:          returnPointer(uint(0)),
		TotalGames:           returnPointer(uint(2)),
		TotalGamesInProgress: returnPointer(uint(0)),
		Dates: []*endpoint.Date{
			{
				Date:                 returnPointer("2022-08-13"),
				TotalItems:           returnPointer(uint(2)),
				TotalEvents:          returnPointer(uint(0)),
				TotalGames:           returnPointer(uint(2)),
				TotalGamesInProgress: returnPointer(uint(0)),
				Games: []*endpoint.Game{
					{
						GamePk:       returnPointer(uint64(67)),
						Link:         returnPointer("/api/v1.1/game/632539/feed/live"),
						GameType:     returnPointer("R"),
						Season:       returnPointer("2021"),
						GameDate:     returnPointer("2022-08-13T20:35:00Z"),
						OfficialDate: returnPointer("2022-08-13"),
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
									ID:   returnPointer(uint(1154534)),
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
									ID:   returnPointer(uint(12555)),
									Name: returnPointer("Baltimore Orioles"),
									Link: returnPointer("api/v1/teams/110"),
								},
							},
						},
						Venue: &endpoint.Venue{
							ID:   returnPointer(uint(12334)),
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
						GamePk:       returnPointer(uint64(1)),
						Link:         returnPointer("/api/v1.1/game/633331/feed/live"),
						GameType:     returnPointer("R"),
						Season:       returnPointer("2021"),
						GameDate:     returnPointer("2022-08-13T20:40:00Z"),
						OfficialDate: returnPointer("2022-08-13"),
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
									Wins:   returnPointer(uint(79)),
									Losses: returnPointer(uint(63)),
									Pct:    returnPointer(".556"),
								},
								Team: &endpoint.Team{
									ID:   returnPointer(uint(11)),
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
									ID:   returnPointer(uint(12)),
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
					{
						GamePk:       returnPointer(uint64(234)),
						Link:         returnPointer("/api/v1.1/game/633331/feed/live"),
						GameType:     returnPointer("R"),
						Season:       returnPointer("2021"),
						GameDate:     returnPointer("2022-08-13T20:40:00Z"),
						OfficialDate: returnPointer("2022-08-13"),
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
									Wins:   returnPointer(uint(79)),
									Losses: returnPointer(uint(63)),
									Pct:    returnPointer(".556"),
								},
								Team: &endpoint.Team{
									ID:   returnPointer(uint(111)),
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
									ID:   returnPointer(uint(121)),
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
					{
						GamePk:       returnPointer(uint64(3)),
						Link:         returnPointer("/api/v1.1/game/632539/feed/live"),
						GameType:     returnPointer("R"),
						Season:       returnPointer("2021"),
						GameDate:     returnPointer("2022-08-13T20:35:00Z"),
						OfficialDate: returnPointer("2022-08-13"),
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
									ID:   returnPointer(uint(11)),
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
									ID:   returnPointer(uint(12)),
									Name: returnPointer("Baltimore Orioles"),
									Link: returnPointer("api/v1/teams/110"),
								},
							},
						},
						Venue: &endpoint.Venue{
							ID:   returnPointer(uint(12)),
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
						GamePk:       returnPointer(uint64(34636)),
						Link:         returnPointer("/api/v1.1/game/632539/feed/live"),
						GameType:     returnPointer("R"),
						Season:       returnPointer("2021"),
						GameDate:     returnPointer("2022-08-13T20:35:00Z"),
						OfficialDate: returnPointer("2022-08-13"),
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
									ID:   returnPointer(uint(1154534)),
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
									ID:   returnPointer(uint(12555)),
									Name: returnPointer("Baltimore Orioles"),
									Link: returnPointer("api/v1/teams/110"),
								},
							},
						},
						Venue: &endpoint.Venue{
							ID:   returnPointer(uint(12334)),
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
				},
				Events: []interface{}{},
			},
		},
	}
}
