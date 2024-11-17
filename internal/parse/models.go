package parse

// Root is the root object that contains all the data
type Root struct {
	MatchData    MatchData     `json:"matchData"`
	LeagueTables []LeagueTable `json:"leagueTable"`
	User         User          `json:"user"`
}

// Badge contains details about the badge of a user
type Badge struct {
	URL     string `json:"url"`
	BadgeID int    `json:"badgeId"`
	UserID  int    `json:"userId"`
}

// User contains details about the user
type User struct {
	Username string `json:"username"`
	UID      int    `json:"uid"`
	LeagueID int    `json:"leagueId"`
	League   League `json:"league"` // Nested league details
	Badge    Badge  `json:"badgeData"`
}

// League contains some details about the league the user is in
type League struct {
	LeagueID    int    `json:"leagueID"`
	Level       int    `json:"level"`
	Name        string `json:"name"`
	Teams       int    `json:"teams"`
	R           int    `json:"r"`
	G           int    `json:"g"`
	B           int    `json:"b"`
	ID          int    `json:"ID"`
	MatchDay    int    `json:"matchDay"`
	StateId     int    `json:"stateId"`
	CommunityId int    `json:"communityId"`
}

// MatchData is the root object that contains data about the last match, next match, etc.
type MatchData struct {
	LastMatch Match `json:"lastMatch"`
}

// Match contains details about a match (nested in MatchData)
type Match struct {
	MatchID             int       `json:"matchId"`
	LeagueID            int       `json:"leagueId"`
	Matchday            int       `json:"matchday"`
	Player1             int       `json:"player1"`
	Player2             int       `json:"player2"`
	GoalsPlayer1        int       `json:"goals_player1"`
	GoalsPlayer2        int       `json:"goals_player2"`
	OlTimestamp         int       `json:"olTimestamp"`
	Timestamp           string    `json:"timestamp"`
	MatchTypeID         int       `json:"matchTypeId"`
	GoalsFirstHalfUser1 int       `json:"goals_first_half_user1"`
	GoalsFirstHalfUser2 int       `json:"goals_first_half_user2"`
	BallPossession1     int       `json:"ballPossession1"`
	BallPossession2     int       `json:"ballPossession2"`
	UserA               MatchUser `json:"userA"`
	UserB               MatchUser `json:"userB"`
}

// MatchUser contains details about a user in a match
type MatchUser struct {
	UID         int    `json:"uid"`
	Username    string `json:"username"`
	TeamName    string `json:"teamName"`
	CommunityID int    `json:"communityId"`
	LeagueID    int    `json:"leagueId"`
	Badge       Badge  `json:"badge"`
}

// LeagueTable contains details about the league table
type LeagueTable struct {
	UserID         int    `json:"userId"`
	LeagueID       int    `json:"leagueId"`
	Matchday       int    `json:"matchday"`
	Type           string `json:"type"`
	Rank           int    `json:"rank"`
	Points         int    `json:"points"`
	ConcedingGoals int    `json:"concedingGoals"`
	ScoredGoals    int    `json:"scoredGoals"`
	Win            int    `json:"win"`
	Lost           int    `json:"lost"`
	Draw           int    `json:"draw"`
	MatchCount     int    `json:"matchCount"`
	TeamName       string `json:"teamName"`
	BadgeID        int    `json:"badgeId"`
	UID            int    `json:"uid"`
}

// MatchResult is the struct that contains the formatted match result that
// will be displayed in the output
type MatchResult struct {
	LeagueInfo     League `json:"leagueInfo"`
	LeagueLevel    string `json:"leagueLevel"`
	BadgeURL       string `json:"badgeURL"`
	LeaguePosition string `json:"leaguePosition"`
	HomeTeam       string `json:"homeTeam"`
	MatchResult    string `json:"matchResult"`
	MatchState     string `json:"matchState"`
	AwayTeam       string `json:"awayTeam"`
	Points         string `json:"points"`
}
