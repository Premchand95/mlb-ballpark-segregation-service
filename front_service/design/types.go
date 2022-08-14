package design

import (
	. "goa.design/goa/v3/dsl"
)

var Schedule = ResultType("application/mlb.schedule", "Schedule", func() {
	Description("Schedule is the detailed custom report for a given date.")
	Attributes(func() {
		Attribute("copyright", String, "mlb copyright for this service.")
		Attribute("totalItems", UInt, "total items in a day", func() {
			Example(15)
		})
		Attribute("totalEvents", UInt, "total events in a day", func() {
			Example(0)
		})
		Attribute("totalGames", UInt, "total games in a day", func() {
			Example(15)
		})
		Attribute("totalGamesInProgress", UInt, "total games in progress", func() {
			Example(0)
		})
		Attribute("dates", ArrayOf(Date), "List of dates with detailed schedule of games.")
	})
})

var Date = Type("Date", func() {
	Attribute("date", String, "official date of the game", func() {
		Example("2021-09-19")
		Format(FormatDate)
	})
	Attribute("totalItems", UInt, "total items in a day", func() {
		Example(15)
	})
	Attribute("totalEvents", UInt, "total events in a day", func() {
		Example(0)
	})
	Attribute("totalGames", UInt, "total games in a day", func() {
		Example(15)
	})
	Attribute("totalGamesInProgress", UInt, "total games in progress", func() {
		Example(0)
	})
	Attribute("games", ArrayOf(Game), "list of games on this date")
	Attribute("events", ArrayOf(Any), "list of events on this date")
	Required("games", "events")
})

var Game = Type("Game", func() {
	Attribute("gamePk", UInt64, "Unique identifier for the game", func() {
		Example(632438)
	})
	Attribute("link", String, "live feed link for the game", func() {
		Example("/api/v1.1/game/632438/feed/live")
	})
	Attribute("gameType", String, "type of the game", func() {
		Example("R")
	})
	Attribute("season", String, "season of the game", func() {
		Example("2021")
	})
	Attribute("gameDate", String, "date of the game", func() {
		Example("2021-09-19T17:05:00Z")
		Format(FormatDateTime)
	})
	Attribute("officialDate", String, "official date of the game", func() {
		Example("2021-09-19")
		Format(FormatDate)
	})
	Attribute("rescheduledFrom", String, "if this game is rescheduled, it's original date", func() {
		Example("2021-07-08T23:05:00Z")
		Format(FormatDateTime)
	})
	Attribute("rescheduledFromDate", String, "official date of the game", func() {
		Example("2021-07-08")
		Format(FormatDate)
	})
	Attribute("description", String, "description of the game", func() {
		Example("Makeup of 7/8 PPD")
	})
	Attribute("status", Status, "status details of the game")
	Attribute("teams", Teams, "details of the two teams of a game")
	Attribute("venue", Venue, "venue of the game")
	Attribute("content", Content, "content of the game")
	Attribute("isTie", Boolean, "is it tie game", func() {
		Example(false)
	})
	Attribute("gameNumber", UInt, "game number", func() {
		Example(1)
	})
	Attribute("publicFacing", Boolean, "is game public facing", func() {
		Example(true)
	})
	Attribute("doubleHeader", String, "double header situation", func() {
		Example("N")
	})
	Attribute("gamedayType", String, "type of the game day", func() {
		Example("P")
	})
	Attribute("tiebreaker", String, "tie breaker", func() {
		Example("N")
	})
	Attribute("calendarEventID", String, "game calender event id", func() {
		Example("14-632438-2021-09-19")
	})
	Attribute("seasonDisplay", String, "game season display", func() {
		Example("2021")
	})
	Attribute("dayNight", String, "is game day or night", func() {
		Example("day")
	})
	Attribute("scheduledInnings", UInt, "scheduled innings of the game", func() {
		Example(9)
	})
	Attribute("reverseHomeAwayStatus", Boolean, "reverse home status of the game", func() {
		Example(false)
	})
	Attribute("inningBreakLength", UInt, "inning Break Length of the game", func() {
		Example(120)
	})
	Attribute("gamesInSeries", UInt, "game In series", func() {
		Example(3)
	})
	Attribute("seriesGameNumber", UInt, "series Number of the game", func() {
		Example(3)
	})
	Attribute("seriesDescription", String, "series description of the game", func() {
		Example("Regular Season")
	})
	Attribute("recordSource", String, "record source of the game", func() {
		Example("S")
	})
	Attribute("ifNecessary", String, "is necessary", func() {
		Example("N")
	})
	Attribute("ifNecessaryDescription", String, "description of the game", func() {
		Example("Normal Game")
	})
})

var Status = Type("Status", func() {
	Attribute("abstractGameState", String, "abstract state of the game", func() {
		Example("Final")
	})
	Attribute("codedGameState", String, "coded game of the game", func() {
		Example("F")
	})
	Attribute("detailedState", String, "detailed state of the game", func() {
		Example("Final")
	})
	Attribute("statusCode", String, "status code of the game", func() {
		Example("F")
	})
	Attribute("startTimeTBD", Boolean, "start time to be determined for this game", func() {
		Example(false)
	})
	Attribute("abstractGameCode", String, "abstract code of the game", func() {
		Example("F")
	})
	Attribute("reason", String, "reason of the game", func() {
		Example("Tied")
	})
})

var Teams = Type("Teams", func() {
	Attribute("away", TeamInfo, "non home team information")
	Attribute("home", TeamInfo, "home team information")
})

var TeamInfo = Type("TeamInfo", func() {
	Attribute("score", UInt, "score of the team", func() {
		Example(3)
	})
	Attribute("isWinner", Boolean, "is this team won this game", func() {
		Example(true)
	})
	Attribute("splitSquad", Boolean, "splitSquad for the team", func() {
		Example(false)
	})
	Attribute("seriesNumber", UInt, "seriesNumber for the team", func() {
		Example(48)
	})
	Attribute("leagueRecord", LeagueRecord, "leagueRecord of the team")
	Attribute("team", Team, "team basic information")
})

var LeagueRecord = Type("LeagueRecord", func() {
	Attribute("wins", UInt, "Number of wins", func() {
		Example(61)
	})
	Attribute("losses", UInt, "Number of losses", func() {
		Example(88)
	})
	Attribute("pct", String, "win percentage. no of wins/total no of matches", func() {
		Example(".409")
	})
})

var Team = Type("Team", func() {
	Attribute("id", UInt, "unique team identifier", func() {
		Example(120)
	})
	Attribute("name", String, "team name", func() {
		Example("Washington Nationals")
	})
	Attribute("link", String, "rest show endpoint to get team details", func() {
		Example("/api/v1/teams/120")
	})
})

var Venue = Type("Venue", func() {
	Attribute("id", UInt, "unique Venue identifier", func() {
		Example(3309)
	})
	Attribute("name", String, "Venue name", func() {
		Example("Nationals Park")
	})
	Attribute("link", String, "rest show endpoint to get Venue details", func() {
		Example("/api/v1/venues/3309")
	})
})

var Content = Type("Content", func() {
	Attribute("link", String, "rest show endpoint to get Content details", func() {
		Example("/api/v1/game/632438/content")
	})
})
