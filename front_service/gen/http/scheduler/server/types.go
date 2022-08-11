// Code generated by goa v3.8.2, DO NOT EDIT.
//
// Scheduler HTTP server types
//
// Command:
// $ goa gen
// github.com/mlb/mlb-ballpark-segregation-service/front_service/design

package server

import (
	scheduler "github.com/mlb/mlb-ballpark-segregation-service/front_service/gen/scheduler"
	schedulerviews "github.com/mlb/mlb-ballpark-segregation-service/front_service/gen/scheduler/views"
	goa "goa.design/goa/v3/pkg"
)

// IndexRequestBody is the type of the "Scheduler" service "index" endpoint
// HTTP request body.
type IndexRequestBody struct {
	// The date (YYYY-mm-dd) used to get all games scheduled on that day
	Date *string `form:"date,omitempty" json:"date,omitempty" xml:"date,omitempty"`
}

// IndexResponseBody is the type of the "Scheduler" service "index" endpoint
// HTTP response body.
type IndexResponseBody struct {
	// mlb copyright for this service.
	Copyright *string `form:"copyright,omitempty" json:"copyright,omitempty" xml:"copyright,omitempty"`
	// total items in a day
	TotalItems *uint `form:"totalItems,omitempty" json:"totalItems,omitempty" xml:"totalItems,omitempty"`
	// total events in a day
	TotalEvents *uint `form:"totalEvents,omitempty" json:"totalEvents,omitempty" xml:"totalEvents,omitempty"`
	// total games in a day
	TotalGames *uint `form:"totalGames,omitempty" json:"totalGames,omitempty" xml:"totalGames,omitempty"`
	// total games in progress
	TotalGamesInProgress *uint `form:"totalGamesInProgress,omitempty" json:"totalGamesInProgress,omitempty" xml:"totalGamesInProgress,omitempty"`
	// List of dates with detailed schedule of games.
	Dates []*DateResponseBody `form:"dates,omitempty" json:"dates,omitempty" xml:"dates,omitempty"`
}

// IndexInternalErrorResponseBody is the type of the "Scheduler" service
// "index" endpoint HTTP response body for the "internal_error" error.
type IndexInternalErrorResponseBody struct {
	// Name is the name of this class of errors.
	Name string `form:"name" json:"name" xml:"name"`
	// ID is a unique identifier for this particular occurrence of the problem.
	ID string `form:"id" json:"id" xml:"id"`
	// Message is a human-readable explanation specific to this occurrence of the
	// problem.
	Message string `form:"message" json:"message" xml:"message"`
	// Is the error temporary?
	Temporary bool `form:"temporary" json:"temporary" xml:"temporary"`
	// Is the error a timeout?
	Timeout bool `form:"timeout" json:"timeout" xml:"timeout"`
	// Is the error a server-side fault?
	Fault bool `form:"fault" json:"fault" xml:"fault"`
}

// IndexBadGatewayResponseBody is the type of the "Scheduler" service "index"
// endpoint HTTP response body for the "bad_gateway" error.
type IndexBadGatewayResponseBody struct {
	// Name is the name of this class of errors.
	Name string `form:"name" json:"name" xml:"name"`
	// ID is a unique identifier for this particular occurrence of the problem.
	ID string `form:"id" json:"id" xml:"id"`
	// Message is a human-readable explanation specific to this occurrence of the
	// problem.
	Message string `form:"message" json:"message" xml:"message"`
	// Is the error temporary?
	Temporary bool `form:"temporary" json:"temporary" xml:"temporary"`
	// Is the error a timeout?
	Timeout bool `form:"timeout" json:"timeout" xml:"timeout"`
	// Is the error a server-side fault?
	Fault bool `form:"fault" json:"fault" xml:"fault"`
}

// IndexBadRequestResponseBody is the type of the "Scheduler" service "index"
// endpoint HTTP response body for the "bad_request" error.
type IndexBadRequestResponseBody struct {
	// Name is the name of this class of errors.
	Name string `form:"name" json:"name" xml:"name"`
	// ID is a unique identifier for this particular occurrence of the problem.
	ID string `form:"id" json:"id" xml:"id"`
	// Message is a human-readable explanation specific to this occurrence of the
	// problem.
	Message string `form:"message" json:"message" xml:"message"`
	// Is the error temporary?
	Temporary bool `form:"temporary" json:"temporary" xml:"temporary"`
	// Is the error a timeout?
	Timeout bool `form:"timeout" json:"timeout" xml:"timeout"`
	// Is the error a server-side fault?
	Fault bool `form:"fault" json:"fault" xml:"fault"`
}

// IndexNotFoundResponseBody is the type of the "Scheduler" service "index"
// endpoint HTTP response body for the "not_found" error.
type IndexNotFoundResponseBody struct {
	// Name is the name of this class of errors.
	Name string `form:"name" json:"name" xml:"name"`
	// ID is a unique identifier for this particular occurrence of the problem.
	ID string `form:"id" json:"id" xml:"id"`
	// Message is a human-readable explanation specific to this occurrence of the
	// problem.
	Message string `form:"message" json:"message" xml:"message"`
	// Is the error temporary?
	Temporary bool `form:"temporary" json:"temporary" xml:"temporary"`
	// Is the error a timeout?
	Timeout bool `form:"timeout" json:"timeout" xml:"timeout"`
	// Is the error a server-side fault?
	Fault bool `form:"fault" json:"fault" xml:"fault"`
}

// DateResponseBody is used to define fields on response body types.
type DateResponseBody struct {
	// official date of the game
	Date *string `form:"date,omitempty" json:"date,omitempty" xml:"date,omitempty"`
	// total items in a day
	TotalItems *uint `form:"totalItems,omitempty" json:"totalItems,omitempty" xml:"totalItems,omitempty"`
	// total events in a day
	TotalEvents *uint `form:"totalEvents,omitempty" json:"totalEvents,omitempty" xml:"totalEvents,omitempty"`
	// total games in a day
	TotalGames *uint `form:"totalGames,omitempty" json:"totalGames,omitempty" xml:"totalGames,omitempty"`
	// total games in progress
	TotalGamesInProgress *uint `form:"totalGamesInProgress,omitempty" json:"totalGamesInProgress,omitempty" xml:"totalGamesInProgress,omitempty"`
	// list of games on this date
	Games []*GameResponseBody `form:"games,omitempty" json:"games,omitempty" xml:"games,omitempty"`
	// list of events on this date
	Events []interface{} `form:"events,omitempty" json:"events,omitempty" xml:"events,omitempty"`
}

// GameResponseBody is used to define fields on response body types.
type GameResponseBody struct {
	// Unique identifier for the game
	GamePk *uint64 `form:"gamePk,omitempty" json:"gamePk,omitempty" xml:"gamePk,omitempty"`
	// live feed link for the game
	Link *string `form:"link,omitempty" json:"link,omitempty" xml:"link,omitempty"`
	// type of the game
	GameType *string `form:"gameType,omitempty" json:"gameType,omitempty" xml:"gameType,omitempty"`
	// season of the game
	Season *string `form:"season,omitempty" json:"season,omitempty" xml:"season,omitempty"`
	// date of the game
	GameDate *string `form:"gameDate,omitempty" json:"gameDate,omitempty" xml:"gameDate,omitempty"`
	// official date of the game
	OfficialDate *string `form:"officialDate,omitempty" json:"officialDate,omitempty" xml:"officialDate,omitempty"`
	// status details of the game
	Status *StatusResponseBody `form:"status,omitempty" json:"status,omitempty" xml:"status,omitempty"`
	// details of the two teams of a game
	Teams *TeamsResponseBody `form:"teams,omitempty" json:"teams,omitempty" xml:"teams,omitempty"`
	// venue of the game
	Venue *VenueResponseBody `form:"venue,omitempty" json:"venue,omitempty" xml:"venue,omitempty"`
	// content of the game
	Content *ContentResponseBody `form:"content,omitempty" json:"content,omitempty" xml:"content,omitempty"`
	// is it tie game
	IsTie *bool `form:"isTie,omitempty" json:"isTie,omitempty" xml:"isTie,omitempty"`
	// game number
	GameNumber *uint `form:"gameNumber,omitempty" json:"gameNumber,omitempty" xml:"gameNumber,omitempty"`
	// is game public facing
	PublicFacing *bool `form:"publicFacing,omitempty" json:"publicFacing,omitempty" xml:"publicFacing,omitempty"`
	// double header situation
	DoubleHeader *string `form:"doubleHeader,omitempty" json:"doubleHeader,omitempty" xml:"doubleHeader,omitempty"`
	// type of the game day
	GamedayType *string `form:"gamedayType,omitempty" json:"gamedayType,omitempty" xml:"gamedayType,omitempty"`
	// tie breaker
	Tiebreaker *string `form:"tiebreaker,omitempty" json:"tiebreaker,omitempty" xml:"tiebreaker,omitempty"`
	// game calender event id
	CalendarEventID *string `form:"calendarEventID,omitempty" json:"calendarEventID,omitempty" xml:"calendarEventID,omitempty"`
	// game season display
	SeasonDisplay *string `form:"seasonDisplay,omitempty" json:"seasonDisplay,omitempty" xml:"seasonDisplay,omitempty"`
	// is game day or night
	DayNight *string `form:"dayNight,omitempty" json:"dayNight,omitempty" xml:"dayNight,omitempty"`
	// scheduled innings of the game
	ScheduledInnings *uint `form:"scheduledInnings,omitempty" json:"scheduledInnings,omitempty" xml:"scheduledInnings,omitempty"`
	// reverse home status of the game
	ReverseHomeAwayStatus *bool `form:"reverseHomeAwayStatus,omitempty" json:"reverseHomeAwayStatus,omitempty" xml:"reverseHomeAwayStatus,omitempty"`
	// inning Break Length of the game
	InningBreakLength *uint `form:"inningBreakLength,omitempty" json:"inningBreakLength,omitempty" xml:"inningBreakLength,omitempty"`
	// game In series
	GamesInSeries *uint `form:"gamesInSeries,omitempty" json:"gamesInSeries,omitempty" xml:"gamesInSeries,omitempty"`
	// series Number of the game
	SeriesGameNumber *uint `form:"seriesGameNumber,omitempty" json:"seriesGameNumber,omitempty" xml:"seriesGameNumber,omitempty"`
	// series description of the game
	SeriesDescription *string `form:"seriesDescription,omitempty" json:"seriesDescription,omitempty" xml:"seriesDescription,omitempty"`
	// record source of the game
	RecordSource *string `form:"recordSource,omitempty" json:"recordSource,omitempty" xml:"recordSource,omitempty"`
	// is necessary
	IfNecessary *string `form:"ifNecessary,omitempty" json:"ifNecessary,omitempty" xml:"ifNecessary,omitempty"`
	// description of the game
	IfNecessaryDescription *string `form:"ifNecessaryDescription,omitempty" json:"ifNecessaryDescription,omitempty" xml:"ifNecessaryDescription,omitempty"`
}

// StatusResponseBody is used to define fields on response body types.
type StatusResponseBody struct {
	// abstract state of the game
	AbstractGameState *string `form:"abstractGameState,omitempty" json:"abstractGameState,omitempty" xml:"abstractGameState,omitempty"`
	// coded game of the game
	CodedGameState *string `form:"codedGameState,omitempty" json:"codedGameState,omitempty" xml:"codedGameState,omitempty"`
	// detailed state of the game
	DetailedState *string `form:"detailedState,omitempty" json:"detailedState,omitempty" xml:"detailedState,omitempty"`
	// status code of the game
	StatusCode *string `form:"statusCode,omitempty" json:"statusCode,omitempty" xml:"statusCode,omitempty"`
	// start time to be determined for this game
	StartTimeTBD *bool `form:"startTimeTBD,omitempty" json:"startTimeTBD,omitempty" xml:"startTimeTBD,omitempty"`
	// abstract code of the game
	AbstractGameCode *string `form:"abstractGameCode,omitempty" json:"abstractGameCode,omitempty" xml:"abstractGameCode,omitempty"`
}

// TeamsResponseBody is used to define fields on response body types.
type TeamsResponseBody struct {
	// non home team information
	Away *TeamInfoResponseBody `form:"away,omitempty" json:"away,omitempty" xml:"away,omitempty"`
	// home team information
	Home *TeamInfoResponseBody `form:"home,omitempty" json:"home,omitempty" xml:"home,omitempty"`
}

// TeamInfoResponseBody is used to define fields on response body types.
type TeamInfoResponseBody struct {
	// score of the team
	Score *uint `form:"score,omitempty" json:"score,omitempty" xml:"score,omitempty"`
	// is this team won this game
	IsWinner *bool `form:"isWinner,omitempty" json:"isWinner,omitempty" xml:"isWinner,omitempty"`
	// splitSquad for the team
	SplitSquad *bool `form:"splitSquad,omitempty" json:"splitSquad,omitempty" xml:"splitSquad,omitempty"`
	// seriesNumber for the team
	SeriesNumber *uint `form:"seriesNumber,omitempty" json:"seriesNumber,omitempty" xml:"seriesNumber,omitempty"`
	// leagueRecord of the team
	LeagueRecord *LeagueRecordResponseBody `form:"LeagueRecord,omitempty" json:"LeagueRecord,omitempty" xml:"LeagueRecord,omitempty"`
	// team basic information
	Team *TeamResponseBody `form:"Team,omitempty" json:"Team,omitempty" xml:"Team,omitempty"`
}

// LeagueRecordResponseBody is used to define fields on response body types.
type LeagueRecordResponseBody struct {
	// Number of wins
	Wins *uint `form:"wins,omitempty" json:"wins,omitempty" xml:"wins,omitempty"`
	// Number of losses
	Losses *uint `form:"losses,omitempty" json:"losses,omitempty" xml:"losses,omitempty"`
	// win percentage. no of wins/total no of matches
	Pct *string `form:"pct,omitempty" json:"pct,omitempty" xml:"pct,omitempty"`
}

// TeamResponseBody is used to define fields on response body types.
type TeamResponseBody struct {
	// unique team identifier
	ID *uint `form:"id,omitempty" json:"id,omitempty" xml:"id,omitempty"`
	// team name
	Name *string `form:"name,omitempty" json:"name,omitempty" xml:"name,omitempty"`
	// rest show endpoint to get team details
	Link *string `form:"link,omitempty" json:"link,omitempty" xml:"link,omitempty"`
}

// VenueResponseBody is used to define fields on response body types.
type VenueResponseBody struct {
	// unique Venue identifier
	ID *uint `form:"id,omitempty" json:"id,omitempty" xml:"id,omitempty"`
	// Venue name
	Name *string `form:"name,omitempty" json:"name,omitempty" xml:"name,omitempty"`
	// rest show endpoint to get Venue details
	Link *string `form:"link,omitempty" json:"link,omitempty" xml:"link,omitempty"`
}

// ContentResponseBody is used to define fields on response body types.
type ContentResponseBody struct {
	// rest show endpoint to get Content details
	Link *string `form:"link,omitempty" json:"link,omitempty" xml:"link,omitempty"`
}

// NewIndexResponseBody builds the HTTP response body from the result of the
// "index" endpoint of the "Scheduler" service.
func NewIndexResponseBody(res *schedulerviews.ScheduleView) *IndexResponseBody {
	body := &IndexResponseBody{
		Copyright:            res.Copyright,
		TotalItems:           res.TotalItems,
		TotalEvents:          res.TotalEvents,
		TotalGames:           res.TotalGames,
		TotalGamesInProgress: res.TotalGamesInProgress,
	}
	if res.Dates != nil {
		body.Dates = make([]*DateResponseBody, len(res.Dates))
		for i, val := range res.Dates {
			body.Dates[i] = marshalSchedulerviewsDateViewToDateResponseBody(val)
		}
	}
	return body
}

// NewIndexInternalErrorResponseBody builds the HTTP response body from the
// result of the "index" endpoint of the "Scheduler" service.
func NewIndexInternalErrorResponseBody(res *goa.ServiceError) *IndexInternalErrorResponseBody {
	body := &IndexInternalErrorResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewIndexBadGatewayResponseBody builds the HTTP response body from the result
// of the "index" endpoint of the "Scheduler" service.
func NewIndexBadGatewayResponseBody(res *goa.ServiceError) *IndexBadGatewayResponseBody {
	body := &IndexBadGatewayResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewIndexBadRequestResponseBody builds the HTTP response body from the result
// of the "index" endpoint of the "Scheduler" service.
func NewIndexBadRequestResponseBody(res *goa.ServiceError) *IndexBadRequestResponseBody {
	body := &IndexBadRequestResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewIndexNotFoundResponseBody builds the HTTP response body from the result
// of the "index" endpoint of the "Scheduler" service.
func NewIndexNotFoundResponseBody(res *goa.ServiceError) *IndexNotFoundResponseBody {
	body := &IndexNotFoundResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewIndexPayload builds a Scheduler service index endpoint payload.
func NewIndexPayload(body *IndexRequestBody, id uint) *scheduler.IndexPayload {
	v := &scheduler.IndexPayload{
		Date: *body.Date,
	}
	v.ID = id

	return v
}

// ValidateIndexRequestBody runs the validations defined on IndexRequestBody
func ValidateIndexRequestBody(body *IndexRequestBody) (err error) {
	if body.Date == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("date", "body"))
	}
	if body.Date != nil {
		err = goa.MergeErrors(err, goa.ValidateFormat("body.date", *body.Date, goa.FormatDate))
	}
	return
}
