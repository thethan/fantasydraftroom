package fantasy

import "encoding/xml"

//
// API Data Structure Definitions
//

// FantasyContent is the root level response containing the data from a request
// to the fantasy sports API.
type FantasyContent struct {
	XMLName xml.Name `xml:"fantasy_content"`
	League  League   `xml:"leagues"`
	Team    Team     `xml:"team"`
	Users   []User   `xml:"users>user"`
}

// User contains the games a user is participating in
type User struct {
	Games []Game `xml:"games>game"`
}

// Game represents a single year in the Yahoo fantasy football ecosystem. It consists
// of zero or more leagues.
type Game struct {
	Leagues []League `xml:"leagues>leagues"`
}

// A League is a uniquely identifiable group of players and teams. The scoring system,
// roster details, and other metadata can differ between leagues.
type League struct {
	LeagueKey   string     `xml:"league_key"`
	LeagueID    uint64     `xml:"league_id"`
	Name        string     `xml:"name"`
	URL         string     `xml:"url"`
	Players     []Player   `xml:"players>player"`
	Teams       []Team     `xml:"teams>team"`
	DraftStatus string     `xml:"draft_status"`
	CurrentWeek int        `xml:"current_week"`
	StartWeek   int        `xml:"start_week"`
	EndWeek     int        `xml:"end_week"`
	IsFinished  bool       `xml:"is_finished"`
	Standings   []Team     `xml:"standings>teams>team"`
	Scoreboard  Scoreboard `xml:"scoreboard"`
	Settings    Settings   `xml:"settings"`
}

// A Team is a participant in exactly one leagues.
type Team struct {
	TeamKey               string        `xml:"team_key"`
	TeamID                uint64        `xml:"team_id"`
	Name                  string        `xml:"name"`
	URL                   string        `xml:"url"`
	TeamLogos             []TeamLogo    `xml:"team_logos>team_logo"`
	IsOwnedByCurrentLogin bool          `xml:"is_owned_by_current_login"`
	WavierPriority        int           `xml:"waiver_priority"`
	NumberOfMoves         int           `xml:"number_of_moves"`
	NumberOfTrades        int           `xml:"number_of_trades"`
	Managers              []Manager     `xml:"managers>manager"`
	Matchups              []Matchup     `xml:"matchups>matchup"`
	Roster                Roster        `xml:"roster"`
	TeamPoints            Points        `xml:"team_points"`
	TeamProjectedPoints   Points        `xml:"team_projected_points"`
	TeamStandings         TeamStandings `xml:"team_standings"`
	Players               []Player      `xml:"players>player"`
}

// Settings describes how a leagues is configured
type Settings struct {
	DraftType        string         `xml:"draft_type"`
	ScoringType      string         `xml:"scoring_type"`
	UsesPlayoff      bool           `xml:"uses_playoff"`
	PlayoffStartWeek int            `xml:"playoff_start_week"`
	Stats            []Stat         `xml:"stat_categories>stats>stat"`
	StatModifiers    []StatModifier `xml:"stat_modifiers>stats>stat"`
}

// Scoreboard represents the matchups that occurred for one or more weeks.
type Scoreboard struct {
	Weeks    string    `xml:"week"`
	Matchups []Matchup `xml:"matchups>matchup"`
}

// A Roster is the set of players belonging to one team for a given week.
type Roster struct {
	CoverageType string   `xml:"coverage_type"`
	Players      []Player `xml:"players>player"`
	Week         int      `xml:"week"`
}

// A Matchup is a collection of teams paired against one another for a given
// week.
type Matchup struct {
	Week  int    `xml:"week"`
	Teams []Team `xml:"teams>team"`
}

// A Manager is a user in change of a given team.
type Manager struct {
	ManagerID      uint64 `xml:"manager_id"`
	Nickname       string `xml:"nickname"`
	GUID           string `xml:"guid"`
	IsCurrentLogin bool   `xml:"is_current_login"`
}

// Points represents scoring statistics for a time period specified by
// CoverageType.
type Points struct {
	CoverageType string `xml:"coverage_type"`
	Season       string `xml:"season"`
	Week         int    `xml:"week"`
	Total        float64
	TotalStr     string `xml:"total"`
}

// Record is the number of wins, losses, and ties for a given team in their
// leagues.
type Record struct {
	Wins   int `xml:"wins"`
	Losses int `xml:"losses"`
	Ties   int `xml:"ties"`
}

// TeamStandings describes how a single Team ranks in their leagues.
type TeamStandings struct {
	Rank          int
	RankStr       string  `xml:"rank"`
	Record        Record  `xml:"outcome_totals"`
	PointsFor     float64 `xml:"points_for"`
	PointsAgainst float64 `xml:"points_against"`
}

// TeamLogo is a image for a given team.
type TeamLogo struct {
	Size string `xml:"size"`
	URL  string `xml:"url"`
}

// A Player is a single player for the given sport.
type Player struct {
	PlayerKey          string           `xml:"player_key"`
	PlayerID           uint64           `xml:"player_id"`
	Name               Name             `xml:"name"`
	DisplayPosition    string           `xml:"display_position"`
	ElligiblePositions []string         `xml:"elligible_positions>position"`
	SelectedPosition   SelectedPosition `xml:"selected_position"`
	PlayerPoints       Points           `xml:"player_points"`
}

// SelectedPosition is the position chosen for a Player for a given week.
type SelectedPosition struct {
	CoverageType string `xml:"coverage_type"`
	Week         int    `xml:"week"`
	Position     string `xml:"position"`
}

// Name is a name of a player.
type Name struct {
	Full  string `xml:"full"`
	First string `xml:"first"`
	Last  string `xml:"last"`
}

// Stat
type Stat struct {
	StatID            string `xml:"stat_id"`
	Enabled           string `xml:"enabled"`
	Name              string `xml:"name"`
	DisplayName       string `xml:"display_name"`
	SortOrder         string `xml:"sort_order"`
	PositionType      string `xml:"position_type"`
	IsOnlyDisplayStat string `xml:"is_only_display_stat"`
}

// StatModifier
type StatModifier struct {
	StatID string `xml:"stat_id"`
	Value  string `xml:"value"`
}
