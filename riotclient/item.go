package riotclient

import "time"

// Item contains the data for one Item
type Item struct {
	Key uint16 `json:"key"`

	Name        string   `json:"name"`
	Description string   `json:"description"`
	Colloq      string   `json:"colloq"`
	Plaintext   string   `json:"plaintext"`
	Into        []string `json:"into"`
	Image       struct {
		Full   string `json:"full"`
		Sprite string `json:"sprite"`
		Group  string `json:"group"`
		X      int    `json:"x"`
		Y      int    `json:"y"`
		W      int    `json:"w"`
		H      int    `json:"h"`
	} `json:"image"`
	Gold struct {
		Base        int  `json:"base"`
		Purchasable bool `json:"purchasable"`
		Total       int  `json:"total"`
		Sell        int  `json:"sell"`
	} `json:"gold"`
	Tags  []string        `json:"tags"`
	Maps  map[string]bool `json:"maps"`
	Stats Stats           `json:"stats"`

	Timestamp time.Time `json:"timestamp"`
}

// BasicItem contains the basic item stats provided from Riot (not exactly sure for what it is used, yet)
type BasicItem struct {
	Name string `json:"name"`
	Rune struct {
		Isrune bool   `json:"isrune"`
		Tier   int    `json:"tier"`
		Type   string `json:"type"`
	} `json:"rune"`
	Gold struct {
		Base        int  `json:"base"`
		Total       int  `json:"total"`
		Sell        int  `json:"sell"`
		Purchasable bool `json:"purchasable"`
	} `json:"gold"`
	Group            string          `json:"group"`
	Description      string          `json:"description"`
	Colloq           string          `json:"colloq"`
	Plaintext        string          `json:"plaintext"`
	Consumed         bool            `json:"consumed"`
	Stacks           int             `json:"stacks"`
	Depth            int             `json:"depth"`
	ConsumeOnFull    bool            `json:"consumeOnFull"`
	From             []interface{}   `json:"from"`
	Into             []interface{}   `json:"into"`
	SpecialRecipe    int             `json:"specialRecipe"`
	InStore          bool            `json:"inStore"`
	HideFromAll      bool            `json:"hideFromAll"`
	RequiredChampion string          `json:"requiredChampion"`
	RequiredAlly     string          `json:"requiredAlly"`
	Stats            Stats           `json:"stats"`
	Tags             []interface{}   `json:"tags"`
	Maps             map[string]bool `json:"maps"`
}

// ItemList is used to pass around a list of items
type ItemList map[uint16]Item

//////////////////////////////////////////////
// Subelementtype definitions follow bellow //
//////////////////////////////////////////////

// Stats contains all stats for a specific item
type Stats struct {
	FlatHPPoolMod                       float32 `json:"FlatHPPoolMod"`
	RFlatHPModPerLevel                  float32 `json:"rFlatHPModPerLevel"`
	FlatMPPoolMod                       float32 `json:"FlatMPPoolMod"`
	RFlatMPModPerLevel                  float32 `json:"rFlatMPModPerLevel"`
	PercentHPPoolMod                    float32 `json:"PercentHPPoolMod"`
	PercentMPPoolMod                    float32 `json:"PercentMPPoolMod"`
	FlatHPRegenMod                      float32 `json:"FlatHPRegenMod"`
	RFlatHPRegenModPerLevel             float32 `json:"rFlatHPRegenModPerLevel"`
	PercentHPRegenMod                   float32 `json:"PercentHPRegenMod"`
	FlatMPRegenMod                      float32 `json:"FlatMPRegenMod"`
	RFlatMPRegenModPerLevel             float32 `json:"rFlatMPRegenModPerLevel"`
	PercentMPRegenMod                   float32 `json:"PercentMPRegenMod"`
	FlatArmorMod                        float32 `json:"FlatArmorMod"`
	RFlatArmorModPerLevel               float32 `json:"rFlatArmorModPerLevel"`
	PercentArmorMod                     float32 `json:"PercentArmorMod"`
	RFlatArmorPenetrationMod            float32 `json:"rFlatArmorPenetrationMod"`
	RFlatArmorPenetrationModPerLevel    float32 `json:"rFlatArmorPenetrationModPerLevel"`
	RPercentArmorPenetrationMod         float32 `json:"rPercentArmorPenetrationMod"`
	RPercentArmorPenetrationModPerLevel float32 `json:"rPercentArmorPenetrationModPerLevel"`
	FlatPhysicalDamageMod               float32 `json:"FlatPhysicalDamageMod"`
	RFlatPhysicalDamageModPerLevel      float32 `json:"rFlatPhysicalDamageModPerLevel"`
	PercentPhysicalDamageMod            float32 `json:"PercentPhysicalDamageMod"`
	FlatMagicDamageMod                  float32 `json:"FlatMagicDamageMod"`
	RFlatMagicDamageModPerLevel         float32 `json:"rFlatMagicDamageModPerLevel"`
	PercentMagicDamageMod               float32 `json:"PercentMagicDamageMod"`
	FlatMovementSpeedMod                float32 `json:"FlatMovementSpeedMod"`
	RFlatMovementSpeedModPerLevel       float32 `json:"rFlatMovementSpeedModPerLevel"`
	PercentMovementSpeedMod             float32 `json:"PercentMovementSpeedMod"`
	RPercentMovementSpeedModPerLevel    float32 `json:"rPercentMovementSpeedModPerLevel"`
	FlatAttackSpeedMod                  float32 `json:"FlatAttackSpeedMod"`
	PercentAttackSpeedMod               float32 `json:"PercentAttackSpeedMod"`
	RPercentAttackSpeedModPerLevel      float32 `json:"rPercentAttackSpeedModPerLevel"`
	RFlatDodgeMod                       float32 `json:"rFlatDodgeMod"`
	RFlatDodgeModPerLevel               float32 `json:"rFlatDodgeModPerLevel"`
	PercentDodgeMod                     float32 `json:"PercentDodgeMod"`
	FlatCritChanceMod                   float32 `json:"FlatCritChanceMod"`
	RFlatCritChanceModPerLevel          float32 `json:"rFlatCritChanceModPerLevel"`
	PercentCritChanceMod                float32 `json:"PercentCritChanceMod"`
	FlatCritDamageMod                   float32 `json:"FlatCritDamageMod"`
	RFlatCritDamageModPerLevel          float32 `json:"rFlatCritDamageModPerLevel"`
	PercentCritDamageMod                float32 `json:"PercentCritDamageMod"`
	FlatBlockMod                        float32 `json:"FlatBlockMod"`
	PercentBlockMod                     float32 `json:"PercentBlockMod"`
	FlatSpellBlockMod                   float32 `json:"FlatSpellBlockMod"`
	RFlatSpellBlockModPerLevel          float32 `json:"rFlatSpellBlockModPerLevel"`
	PercentSpellBlockMod                float32 `json:"PercentSpellBlockMod"`
	FlatEXPBonus                        float32 `json:"FlatEXPBonus"`
	PercentEXPBonus                     float32 `json:"PercentEXPBonus"`
	RPercentCooldownMod                 float32 `json:"rPercentCooldownMod"`
	RPercentCooldownModPerLevel         float32 `json:"rPercentCooldownModPerLevel"`
	RFlatTimeDeadMod                    float32 `json:"rFlatTimeDeadMod"`
	RFlatTimeDeadModPerLevel            float32 `json:"rFlatTimeDeadModPerLevel"`
	RPercentTimeDeadMod                 float32 `json:"rPercentTimeDeadMod"`
	RPercentTimeDeadModPerLevel         float32 `json:"rPercentTimeDeadModPerLevel"`
	RFlatGoldPer10Mod                   float32 `json:"rFlatGoldPer10Mod"`
	RFlatMagicPenetrationMod            float32 `json:"rFlatMagicPenetrationMod"`
	RFlatMagicPenetrationModPerLevel    float32 `json:"rFlatMagicPenetrationModPerLevel"`
	RPercentMagicPenetrationMod         float32 `json:"rPercentMagicPenetrationMod"`
	RPercentMagicPenetrationModPerLevel float32 `json:"rPercentMagicPenetrationModPerLevel"`
	FlatEnergyRegenMod                  float32 `json:"FlatEnergyRegenMod"`
	RFlatEnergyRegenModPerLevel         float32 `json:"rFlatEnergyRegenModPerLevel"`
	FlatEnergyPoolMod                   float32 `json:"FlatEnergyPoolMod"`
	RFlatEnergyModPerLevel              float32 `json:"rFlatEnergyModPerLevel"`
	PercentLifeStealMod                 float32 `json:"PercentLifeStealMod"`
	PercentSpellVampMod                 float32 `json:"PercentSpellVampMod"`
}
