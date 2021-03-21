package models

type Player struct {
	ActivePlayerID             *int64  `json:"ActivePlayerId,omitempty"`
	AvatarURL                  *string `json:"Avatar_URL,omitempty"`
	CreatedDatetime            *string `json:"Created_Datetime,omitempty"`
	HoursPlayed                *int64  `json:"HoursPlayed,omitempty"`
	ID                         *int64  `json:"Id,omitempty"`
	LastLoginDatetime          *string `json:"Last_Login_Datetime,omitempty"`
	Leaves                     *int64  `json:"Leaves,omitempty"`
	Level                      *int64  `json:"Level,omitempty"`
	Losses                     *int64  `json:"Losses,omitempty"`
	MasteryLevel               *int64  `json:"MasteryLevel,omitempty"`
	MinutesPlayed              *int64  `json:"MinutesPlayed,omitempty"`
	Name                       *string `json:"Name,omitempty"`
	PersonalStatusMessage      *string `json:"Personal_Status_Message,omitempty"`
	Platform                   *string `json:"Platform,omitempty"`
	RankStatConquest           *int64  `json:"Rank_Stat_Conquest,omitempty"`
	RankStatConquestController *int64  `json:"Rank_Stat_Conquest_Controller,omitempty"`
	RankStatDuel               *int64  `json:"Rank_Stat_Duel,omitempty"`
	RankStatDuelController     *int64  `json:"Rank_Stat_Duel_Controller,omitempty"`
	RankStatJoust              *int64  `json:"Rank_Stat_Joust,omitempty"`
	RankStatJoustController    *int64  `json:"Rank_Stat_Joust_Controller,omitempty"`
	RankedConquest             *Ranked `json:"RankedConquest,omitempty"`
	RankedConquestController   *Ranked `json:"RankedConquestController,omitempty"`
	RankedDuel                 *Ranked `json:"RankedDuel,omitempty"`
	RankedDuelController       *Ranked `json:"RankedDuelController,omitempty"`
	RankedJoust                *Ranked `json:"RankedJoust,omitempty"`
	RankedJoustController      *Ranked `json:"RankedJoustController,omitempty"`
	Region                     *string `json:"Region,omitempty"`
	TeamID                     *int64  `json:"TeamId,omitempty"`
	TeamName                   *string `json:"Team_Name,omitempty"`
	TierConquest               *int64  `json:"Tier_Conquest,omitempty"`
	TierDuel                   *int64  `json:"Tier_Duel,omitempty"`
	TierJoust                  *int64  `json:"Tier_Joust,omitempty"`
	TotalAchievements          *int64  `json:"Total_Achievements,omitempty"`
	TotalWorshippers           *int64  `json:"Total_Worshippers,omitempty"`
	WINS                       *int64  `json:"Wins,omitempty"`
	HzPlayerName               *string `json:"hz_player_name,omitempty"`
	RetMsg                     *string `json:"ret_msg,omitempty"`
}

type Ranked struct {
	Leaves       *int64  `json:"Leaves,omitempty"`
	Losses       *int64  `json:"Losses,omitempty"`
	Name         *string `json:"Name,omitempty"`
	Points       *int64  `json:"Points,omitempty"`
	PrevRank     *int64  `json:"PrevRank,omitempty"`
	Rank         *int64  `json:"Rank,omitempty"`
	RankStat     *int64  `json:"Rank_Stat,omitempty"`
	RankVariance *int64  `json:"Rank_Variance,omitempty"`
	Season       *int64  `json:"Season,omitempty"`
	Tier         *int64  `json:"Tier,omitempty"`
	Trend        *int64  `json:"Trend,omitempty"`
	WINS         *int64  `json:"Wins,omitempty"`
	PlayerID     *string `json:"player_id"`
	RetMsg       *string `json:"ret_msg"`
}
