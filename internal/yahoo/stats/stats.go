package stats

// Service Import stats for a league ID
type Service interface {
	ImportLeagueStats(leagueID int)
}