package gameinfo

type CharacterBase struct {
	// ID 角色 id
	ID int `json:"id"`
	// Name 名称
	Name string `json:"name"`
	// Rarity 星级
	Rarity uint `json:"rarity"`
	// Constellations 命座
	Constellations []Constellation `json:"constellations"`

	CharacterLevels
	CharacterIcons
}

type CharacterLevels struct {
	// Element 属性
	Element string `json:"element"`
	// Level 等级
	Level uint `json:"level"`
	// Fetter 好感度
	Fetter int `json:"fetter"`
	// ActiveConstellationNum 激活命座数
	ActiveConstellationNum int `json:"actived_constellation_num"`
}

type CharacterIcons struct {
	// Icon 头像
	Icon string `json:"icon,omitempty"`
	// Image 立绘
	Image string `json:"image,omitempty"`
	// 元素 Icon
	ElementIcon string `json:"elementIcon,omitempty"`
}

// Character 角色
type Character struct {
	CharacterBase
	// Weapon 武器
	Weapon Weapon `json:"weapon"`
	// Reliquaries 圣遗物
	Reliquaries []Reliquary `json:"reliquaries"`
}

// Weapon 武器
type Weapon struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Type     int    `json:"type,omitempty"`
	TypeName string `json:"type_name"`
	Desc     string `json:"desc"`
	Rarity   uint   `json:"rarity"`

	WeaponIcons
	WeaponLevels
}

type WeaponIcons struct {
	Icon string `json:"icon"`
}

type WeaponLevels struct {
	// 等级
	Level uint `json:"level"`
	// PromoteLevel 突破等级
	PromoteLevel uint `json:"promote_level"`
	// AffixLevel 精炼等级
	AffixLevel uint `json:"affix_level"`
}

type ReliquaryIcons struct {
	Icon string `json:"icon"`
}

type ReliquaryLevels struct {
	Level uint `json:"level"`
}

// Reliquary 圣遗物
type Reliquary struct {
	ReliquaryIcons
	ReliquaryLevels

	ID      int    `json:"id"`
	Name    string `json:"name"`
	Pos     int    `json:"pos"`
	PosName string `json:"pos_name"`
	Rarity  uint   `json:"rarity"`

	Set struct {
		ID      int    `json:"id"`
		Name    string `json:"name"`
		Affixes []struct {
			ActivationNumber int    `json:"activation_number"`
			Effect           string `json:"effect"`
		} `json:"affixes"`
	} `json:"set"`
}

// Constellation 命座
type Constellation struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Icon     string `json:"icon"`
	Effect   string `json:"effect"`
	IsActive bool   `json:"is_actived"`
	Pos      int    `json:"pos"`
}

// UserInfo 用户信息
type UserInfo struct {
	Role    []Character `json:"role,omitempty"`
	Avatars []Character `json:"avatars,omitempty"`
}
