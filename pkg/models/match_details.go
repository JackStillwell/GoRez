package models

type MatchDetails struct {
	AccountLevel         *int64   `json:"Account_Level,omitempty"`
	ActiveId1            *int64   `json:"ActiveId1,omitempty"`
	ActiveId2            *int64   `json:"ActiveId2,omitempty"`
	ActiveId3            *int64   `json:"ActiveId3,omitempty"`
	ActiveId4            *int64   `json:"ActiveId4,omitempty"`
	ActivePlayerID       *string  `json:"ActivePlayerId,omitempty"`
	Assists              *int64   `json:"Assists,omitempty"`
	Ban1                 *string  `json:"Ban1,omitempty"`
	Ban10                *string  `json:"Ban10,omitempty"`
	Ban10ID              *int64   `json:"Ban10Id,omitempty"`
	Ban1ID               *int64   `json:"Ban1Id,omitempty"`
	Ban2                 *string  `json:"Ban2,omitempty"`
	Ban2ID               *int64   `json:"Ban2Id,omitempty"`
	Ban3                 *string  `json:"Ban3,omitempty"`
	Ban3ID               *int64   `json:"Ban3Id,omitempty"`
	Ban4                 *string  `json:"Ban4,omitempty"`
	Ban4ID               *int64   `json:"Ban4Id,omitempty"`
	Ban5                 *string  `json:"Ban5,omitempty"`
	Ban5ID               *int64   `json:"Ban5Id,omitempty"`
	Ban6                 *string  `json:"Ban6,omitempty"`
	Ban6ID               *int64   `json:"Ban6Id,omitempty"`
	Ban7                 *string  `json:"Ban7,omitempty"`
	Ban7ID               *int64   `json:"Ban7Id,omitempty"`
	Ban8                 *string  `json:"Ban8,omitempty"`
	Ban8ID               *int64   `json:"Ban8Id,omitempty"`
	Ban9                 *string  `json:"Ban9,omitempty"`
	Ban9ID               *int64   `json:"Ban9Id,omitempty"`
	CampsCleared         *int64   `json:"Camps_Cleared,omitempty"`
	ConquestLosses       *int64   `json:"Conquest_Losses,omitempty"`
	ConquestPoints       *int64   `json:"Conquest_Points,omitempty"`
	ConquestTier         *int64   `json:"Conquest_Tier,omitempty"`
	ConquestWINS         *int64   `json:"Conquest_Wins,omitempty"`
	DamageBot            *int64   `json:"Damage_Bot,omitempty"`
	DamageDoneInHand     *int64   `json:"Damage_Done_In_Hand,omitempty"`
	DamageDoneMagical    *int64   `json:"Damage_Done_Magical,omitempty"`
	DamageDonePhysical   *int64   `json:"Damage_Done_Physical,omitempty"`
	DamageMitigated      *int64   `json:"Damage_Mitigated,omitempty"`
	DamagePlayer         *int64   `json:"Damage_Player,omitempty"`
	DamageTaken          *int64   `json:"Damage_Taken,omitempty"`
	DamageTakenMagical   *int64   `json:"Damage_Taken_Magical,omitempty"`
	DamageTakenPhysical  *int64   `json:"Damage_Taken_Physical,omitempty"`
	Deaths               *int64   `json:"Deaths,omitempty"`
	DistanceTraveled     *int64   `json:"Distance_Traveled,omitempty"`
	DuelLosses           *int64   `json:"Duel_Losses,omitempty"`
	DuelPoints           *int64   `json:"Duel_Points,omitempty"`
	DuelTier             *int64   `json:"Duel_Tier,omitempty"`
	DuelWINS             *int64   `json:"Duel_Wins,omitempty"`
	EntryDatetime        *string  `json:"Entry_Datetime,omitempty"`
	FinalMatchLevel      *int64   `json:"Final_Match_Level,omitempty"`
	FirstBanSide         *string  `json:"First_Ban_Side,omitempty"`
	GodID                *int64   `json:"GodId,omitempty"`
	GoldEarned           *int64   `json:"Gold_Earned,omitempty"`
	GoldPerMinute        *int64   `json:"Gold_Per_Minute,omitempty"`
	Healing              *int64   `json:"Healing,omitempty"`
	HealingBot           *int64   `json:"Healing_Bot,omitempty"`
	HealingPlayerSelf    *int64   `json:"Healing_Player_Self,omitempty"`
	ItemId1              *int64   `json:"ItemId1,omitempty"`
	ItemId2              *int64   `json:"ItemId2,omitempty"`
	ItemId3              *int64   `json:"ItemId3,omitempty"`
	ItemId4              *int64   `json:"ItemId4,omitempty"`
	ItemId5              *int64   `json:"ItemId5,omitempty"`
	ItemId6              *int64   `json:"ItemId6,omitempty"`
	ItemActive1          *string  `json:"Item_Active_1,omitempty"`
	ItemActive2          *string  `json:"Item_Active_2,omitempty"`
	ItemActive3          *string  `json:"Item_Active_3,omitempty"`
	ItemActive4          *string  `json:"Item_Active_4,omitempty"`
	ItemPurch1           *string  `json:"Item_Purch_1,omitempty"`
	ItemPurch2           *string  `json:"Item_Purch_2,omitempty"`
	ItemPurch3           *string  `json:"Item_Purch_3,omitempty"`
	ItemPurch4           *string  `json:"Item_Purch_4,omitempty"`
	ItemPurch5           *string  `json:"Item_Purch_5,omitempty"`
	ItemPurch6           *string  `json:"Item_Purch_6,omitempty"`
	JoustLosses          *int64   `json:"Joust_Losses,omitempty"`
	JoustPoints          *int64   `json:"Joust_Points,omitempty"`
	JoustTier            *int64   `json:"Joust_Tier,omitempty"`
	JoustWINS            *int64   `json:"Joust_Wins,omitempty"`
	KillingSpree         *int64   `json:"Killing_Spree,omitempty"`
	KillsBot             *int64   `json:"Kills_Bot,omitempty"`
	KillsDouble          *int64   `json:"Kills_Double,omitempty"`
	KillsFireGiant       *int64   `json:"Kills_Fire_Giant,omitempty"`
	KillsFirstBlood      *int64   `json:"Kills_First_Blood,omitempty"`
	KillsGoldFury        *int64   `json:"Kills_Gold_Fury,omitempty"`
	KillsPenta           *int64   `json:"Kills_Penta,omitempty"`
	KillsPhoenix         *int64   `json:"Kills_Phoenix,omitempty"`
	KillsPlayer          *int64   `json:"Kills_Player,omitempty"`
	KillsQuadra          *int64   `json:"Kills_Quadra,omitempty"`
	KillsSiegeJuggernaut *int64   `json:"Kills_Siege_Juggernaut,omitempty"`
	KillsSingle          *int64   `json:"Kills_Single,omitempty"`
	KillsTriple          *int64   `json:"Kills_Triple,omitempty"`
	KillsWildJuggernaut  *int64   `json:"Kills_Wild_Juggernaut,omitempty"`
	MapGame              *string  `json:"Map_Game,omitempty"`
	MasteryLevel         *int64   `json:"Mastery_Level,omitempty"`
	Match                *int64   `json:"Match,omitempty"`
	MatchDuration        *int64   `json:"Match_Duration,omitempty"`
	Minutes              *int64   `json:"Minutes,omitempty"`
	MultiKillMax         *int64   `json:"Multi_kill_Max,omitempty"`
	ObjectiveAssists     *int64   `json:"Objective_Assists,omitempty"`
	PartyID              *int64   `json:"PartyId,omitempty"`
	RankStatConquest     *float64 `json:"Rank_Stat_Conquest,omitempty"`
	RankStatDuel         *int64   `json:"Rank_Stat_Duel,omitempty"`
	RankStatJoust        *int64   `json:"Rank_Stat_Joust,omitempty"`
	ReferenceName        *string  `json:"Reference_Name,omitempty"`
	Region               *string  `json:"Region,omitempty"`
	Skin                 *string  `json:"Skin,omitempty"`
	SkinID               *int64   `json:"SkinId,omitempty"`
	StructureDamage      *int64   `json:"Structure_Damage,omitempty"`
	Surrendered          *int64   `json:"Surrendered,omitempty"`
	TaskForce            *int64   `json:"TaskForce,omitempty"`
	Team1Score           *int64   `json:"Team1Score,omitempty"`
	Team2Score           *int64   `json:"Team2Score,omitempty"`
	TeamID               *int64   `json:"TeamId,omitempty"`
	TeamName             *string  `json:"Team_Name,omitempty"`
	TimeDeadSeconds      *int64   `json:"Time_Dead_Seconds,omitempty"`
	TimeInMatchSeconds   *int64   `json:"Time_In_Match_Seconds,omitempty"`
	TowersDestroyed      *int64   `json:"Towers_Destroyed,omitempty"`
	WardsPlaced          *int64   `json:"Wards_Placed,omitempty"`
	WinStatus            *string  `json:"Win_Status,omitempty"`
	WinningTaskForce     *int64   `json:"Winning_TaskForce,omitempty"`
	HasReplay            *string  `json:"hasReplay,omitempty"`
	HzPlayerName         *string  `json:"hz_player_name,omitempty"`
	MatchQueueID         *int64   `json:"match_queue_id,omitempty"`
	Name                 *string  `json:"name,omitempty"`
	PlayerID             *string  `json:"playerId,omitempty"`
	PlayerName           *string  `json:"playerName,omitempty"`
	RetMsg               *string  `json:"ret_msg,omitempty"`
}