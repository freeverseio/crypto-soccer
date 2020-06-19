package leaderboard

type Entry struct {
	TeamId   string
	Points   int
	Position int
}

type Leaderboard [8]Entry
