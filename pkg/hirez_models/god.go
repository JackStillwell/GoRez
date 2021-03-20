package hirez_models

type God struct {
	Ability1                   *string             `json:"Ability1,omitempty"`
	Ability2                   *string             `json:"Ability2,omitempty"`
	Ability3                   *string             `json:"Ability3,omitempty"`
	Ability4                   *string             `json:"Ability4,omitempty"`
	Ability5                   *string             `json:"Ability5,omitempty"`
	AbilityId1                 *int64              `json:"AbilityId1,omitempty"`
	AbilityId2                 *int64              `json:"AbilityId2,omitempty"`
	AbilityId3                 *int64              `json:"AbilityId3,omitempty"`
	AbilityId4                 *int64              `json:"AbilityId4,omitempty"`
	AbilityId5                 *int64              `json:"AbilityId5,omitempty"`
	GodAbility1                *Ability            `json:"Ability_1,omitempty"`
	GodAbility2                *Ability            `json:"Ability_2,omitempty"`
	GodAbility3                *Ability            `json:"Ability_3,omitempty"`
	GodAbility4                *Ability            `json:"Ability_4,omitempty"`
	GodAbility5                *Ability            `json:"Ability_5,omitempty"`
	AttackSpeed                *float64            `json:"AttackSpeed,omitempty"`
	AttackSpeedPerLevel        *float64            `json:"AttackSpeedPerLevel,omitempty"`
	AutoBanned                 *string             `json:"AutoBanned,omitempty"`
	Cons                       *string             `json:"Cons,omitempty"`
	HP5PerLevel                *float64            `json:"HP5PerLevel,omitempty"`
	Health                     *int64              `json:"Health,omitempty"`
	HealthPerFive              *int64              `json:"HealthPerFive,omitempty"`
	HealthPerLevel             *int64              `json:"HealthPerLevel,omitempty"`
	Lore                       *string             `json:"Lore,omitempty"`
	MP5PerLevel                *float64            `json:"MP5PerLevel,omitempty"`
	MagicProtection            *int64              `json:"MagicProtection,omitempty"`
	MagicProtectionPerLevel    *float64            `json:"MagicProtectionPerLevel,omitempty"`
	MagicalPower               *int64              `json:"MagicalPower,omitempty"`
	MagicalPowerPerLevel       *int64              `json:"MagicalPowerPerLevel,omitempty"`
	Mana                       *int64              `json:"Mana,omitempty"`
	ManaPerFive                *float64            `json:"ManaPerFive,omitempty"`
	ManaPerLevel               *int64              `json:"ManaPerLevel,omitempty"`
	Name                       *string             `json:"Name,omitempty"`
	OnFreeRotation             *string             `json:"OnFreeRotation,omitempty"`
	Pantheon                   *string             `json:"Pantheon,omitempty"`
	PhysicalPower              *int64              `json:"PhysicalPower,omitempty"`
	PhysicalPowerPerLevel      *int64              `json:"PhysicalPowerPerLevel,omitempty"`
	PhysicalProtection         *int64              `json:"PhysicalProtection,omitempty"`
	PhysicalProtectionPerLevel *int64              `json:"PhysicalProtectionPerLevel,omitempty"`
	Pros                       *string             `json:"Pros,omitempty"`
	Roles                      *string             `json:"Roles,omitempty"`
	Speed                      *int64              `json:"Speed,omitempty"`
	Title                      *string             `json:"Title,omitempty"`
	Type                       *string             `json:"Type,omitempty"`
	AbilityDescription1        *AbilityDescription `json:"abilityDescription1,omitempty"`
	AbilityDescription2        *AbilityDescription `json:"abilityDescription2,omitempty"`
	AbilityDescription3        *AbilityDescription `json:"abilityDescription3,omitempty"`
	AbilityDescription4        *AbilityDescription `json:"abilityDescription4,omitempty"`
	AbilityDescription5        *AbilityDescription `json:"abilityDescription5,omitempty"`
	BasicAttack                *AbilityDescription `json:"basicAttack,omitempty"`
	GodAbility1URL             *string             `json:"godAbility1_URL,omitempty"`
	GodAbility2URL             *string             `json:"godAbility2_URL,omitempty"`
	GodAbility3URL             *string             `json:"godAbility3_URL,omitempty"`
	GodAbility4URL             *string             `json:"godAbility4_URL,omitempty"`
	GodAbility5URL             *string             `json:"godAbility5_URL,omitempty"`
	GodCardURL                 *string             `json:"godCard_URL,omitempty"`
	GodIconURL                 *string             `json:"godIcon_URL,omitempty"`
	ID                         *int64              `json:"id,omitempty"`
	LatestGo                   *string             `json:"latestGo,omitempty"`
	RetMsg                     interface{}         `json:"ret_msg"`
}

type AbilityDescription struct {
	Description *Description `json:"itemDescription,omitempty"`
}

type Description struct {
	Cooldown    *string    `json:"cooldown,omitempty"`
	Cost        *string    `json:"cost,omitempty"`
	Description *string    `json:"description,omitempty"`
	MenuItems   []MenuItem `json:"menuitems,omitempty"`
	RankItems   []MenuItem `json:"rankitems,omitempty"`
}

type Ability struct {
	Description *AbilityDescription `json:"Description,omitempty"`
	ID          *int64              `json:"Id,omitempty"`
	Summary     *string             `json:"Summary,omitempty"`
	URL         *string             `json:"URL,omitempty"`
}
