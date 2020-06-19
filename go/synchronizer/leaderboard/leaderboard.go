package leaderboard

type Entry struct {
	teamId   string
	points   int
	position int
}

type Leaderboard [8]Entry

type LeaderboardService interface {
	Compute(timezone int, country int, league int) (*Leaderboard, error)
}
