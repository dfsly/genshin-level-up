package genshindb

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

const (
	FIGHT_PROP_BASE_HP           = "FIGHT_PROP_BASE_HP"
	FIGHT_PROP_BASE_ATTACK       = "FIGHT_PROP_BASE_ATTACK"
	FIGHT_PROP_BASE_DEFENSE      = "FIGHT_PROP_BASE_DEFENSE"
	FIGHT_PROP_HP                = "FIGHT_PROP_HP"
	FIGHT_PROP_HP_PERCENT        = "FIGHT_PROP_HP_PERCENT"
	FIGHT_PROP_ATTACK            = "FIGHT_PROP_ATTACK"
	FIGHT_PROP_ATTACK_PERCENT    = "FIGHT_PROP_ATTACK_PERCENT"
	FIGHT_PROP_DEFENSE           = "FIGHT_PROP_DEFENSE"
	FIGHT_PROP_DEFENSE_PERCENT   = "FIGHT_PROP_DEFENSE_PERCENT"
	FIGHT_PROP_CRITICAL          = "FIGHT_PROP_CRITICAL"
	FIGHT_PROP_CRITICAL_HURT     = "FIGHT_PROP_CRITICAL_HURT"
	FIGHT_PROP_CHARGE_EFFICIENCY = "FIGHT_PROP_CHARGE_EFFICIENCY"
	FIGHT_PROP_HEAL_ADD          = "FIGHT_PROP_HEAL_ADD"
	FIGHT_PROP_ELEMENT_MASTERY   = "FIGHT_PROP_ELEMENT_MASTERY"
	FIGHT_PROP_FIRE_ADD_HURT     = "FIGHT_PROP_FIRE_ADD_HURT"
	FIGHT_PROP_ELEC_ADD_HURT     = "FIGHT_PROP_ELEC_ADD_HURT"
	FIGHT_PROP_WATER_ADD_HURT    = "FIGHT_PROP_WATER_ADD_HURT"
	FIGHT_PROP_WIND_ADD_HURT     = "FIGHT_PROP_WIND_ADD_HURT"
	FIGHT_PROP_ROCK_ADD_HURT     = "FIGHT_PROP_ROCK_ADD_HURT"
	FIGHT_PROP_GRASS_ADD_HURT    = "FIGHT_PROP_GRASS_ADD_HURT"
	FIGHT_PROP_ICE_ADD_HURT      = "FIGHT_PROP_ICE_ADD_HURT"
	FIGHT_PROP_PHYSICAL_ADD_HURT = "FIGHT_PROP_PHYSICAL_ADD_HURT"
	FIGHT_PROP_FIRE_SUB_HURT     = "FIGHT_PROP_FIRE_SUB_HURT"
)

var FightPropKeys = []string{
	FIGHT_PROP_BASE_HP,
	FIGHT_PROP_HP,
	FIGHT_PROP_BASE_ATTACK,
	FIGHT_PROP_ATTACK,
	FIGHT_PROP_BASE_DEFENSE,
	FIGHT_PROP_DEFENSE,
	FIGHT_PROP_HP_PERCENT,
	FIGHT_PROP_ATTACK_PERCENT,
	FIGHT_PROP_DEFENSE_PERCENT,
	FIGHT_PROP_CRITICAL,
	FIGHT_PROP_CRITICAL_HURT,
	FIGHT_PROP_CHARGE_EFFICIENCY,
	FIGHT_PROP_HEAL_ADD,
	FIGHT_PROP_ELEMENT_MASTERY,
	FIGHT_PROP_FIRE_ADD_HURT,
	FIGHT_PROP_ELEC_ADD_HURT,
	FIGHT_PROP_WATER_ADD_HURT,
	FIGHT_PROP_WIND_ADD_HURT,
	FIGHT_PROP_ROCK_ADD_HURT,
	FIGHT_PROP_GRASS_ADD_HURT,
	FIGHT_PROP_ICE_ADD_HURT,
	FIGHT_PROP_PHYSICAL_ADD_HURT,
	FIGHT_PROP_FIRE_SUB_HURT,
}

var FightPropParamNames = map[string]string{
	FIGHT_PROP_BASE_HP:           "基础生命值|{param1:F}",
	FIGHT_PROP_BASE_ATTACK:       "基础攻击力|{param1:F}",
	FIGHT_PROP_BASE_DEFENSE:      "基础防御力|{param1:F}",
	FIGHT_PROP_HP:                "生命值|{param1:F}",
	FIGHT_PROP_HP_PERCENT:        "生命值|{param1:F1P}",
	FIGHT_PROP_ATTACK:            "攻击力|{param1:F}",
	FIGHT_PROP_ATTACK_PERCENT:    "攻击力|{param1:F1P}",
	FIGHT_PROP_DEFENSE:           "防御力|{param1:F}",
	FIGHT_PROP_DEFENSE_PERCENT:   "防御力|{param1:F1P}",
	FIGHT_PROP_CRITICAL:          "暴击率|{param1:F1P}",
	FIGHT_PROP_CRITICAL_HURT:     "暴击伤害|{param1:F1P}",
	FIGHT_PROP_CHARGE_EFFICIENCY: "元素充能|{param1:F1P}",
	FIGHT_PROP_HEAL_ADD:          "治疗加成|{param1:F1P}",
	FIGHT_PROP_ELEMENT_MASTERY:   "元素精通|{param1:F}",
	FIGHT_PROP_FIRE_ADD_HURT:     "火伤加成|{param1:F1P}",
	FIGHT_PROP_ELEC_ADD_HURT:     "雷伤加成|{param1:F1P}",
	FIGHT_PROP_WATER_ADD_HURT:    "水伤加成|{param1:F1P}",
	FIGHT_PROP_WIND_ADD_HURT:     "风伤加成|{param1:F1P}",
	FIGHT_PROP_ROCK_ADD_HURT:     "岩伤加成|{param1:F1P}",
	FIGHT_PROP_GRASS_ADD_HURT:    "草伤加成|{param1:F1P}",
	FIGHT_PROP_ICE_ADD_HURT:      "冰伤加成|{param1:F1P}",
	FIGHT_PROP_PHYSICAL_ADD_HURT: "物伤加成|{param1:F1P}",
	FIGHT_PROP_FIRE_SUB_HURT:     "火抗加成|{param1:F1P}",
}

var nonPercentFightProps = map[string]bool{
	FIGHT_PROP_BASE_HP:      true,
	FIGHT_PROP_BASE_ATTACK:  true,
	FIGHT_PROP_BASE_DEFENSE: true,
	FIGHT_PROP_HP:           true,
	FIGHT_PROP_ATTACK:       true,
	FIGHT_PROP_DEFENSE:      true,
}

var FightPropColors = map[string]string{
	FIGHT_PROP_FIRE_ADD_HURT:  "rgb(220, 20, 60)",
	FIGHT_PROP_ELEC_ADD_HURT:  "rgb(138, 43, 226)",
	FIGHT_PROP_WATER_ADD_HURT: "rgb(30, 144, 255)",
	FIGHT_PROP_ICE_ADD_HURT:   "rgb(0, 191, 255)",
	FIGHT_PROP_WIND_ADD_HURT:  "rgb(102, 205, 170)",
	FIGHT_PROP_ROCK_ADD_HURT:  "rgb(218, 165, 32)",
	FIGHT_PROP_GRASS_ADD_HURT: "green",
}

var FightPropElements = map[string]string{
	FIGHT_PROP_FIRE_ADD_HURT:  "火",
	FIGHT_PROP_ELEC_ADD_HURT:  "雷",
	FIGHT_PROP_WATER_ADD_HURT: "水",
	FIGHT_PROP_ICE_ADD_HURT:   "冰",
	FIGHT_PROP_WIND_ADD_HURT:  "风",
	FIGHT_PROP_ROCK_ADD_HURT:  "岩",
	FIGHT_PROP_GRASS_ADD_HURT: "草",
}

func isPercentProp(k string) bool {
	if _, ok := nonPercentFightProps[k]; !ok {
		return true
	}
	return false
}

type FightProps map[string]float64

func (fightProps FightProps) Get(fp string) float64 {
	if v, ok := fightProps[fp]; ok {
		return v
	}
	return 0
}

func (fightProps FightProps) Format() map[string][]string {
	m := map[string][]string{}

	for key, v := range fightProps {
		m[key] = ParseParamName(FightPropParamNames[key], v)
	}

	return m
}

func (fightProps FightProps) Compute() FightProps {
	computedFightProps := FightProps{}

	for key, v := range fightProps {
		if v == 0 {
			continue
		}

		switch key {
		case FIGHT_PROP_ATTACK_PERCENT:
			computedFightProps.Add(FIGHT_PROP_ATTACK, fightProps[FIGHT_PROP_BASE_ATTACK]*v)
		case FIGHT_PROP_HP_PERCENT:
			computedFightProps.Add(FIGHT_PROP_HP, fightProps[FIGHT_PROP_BASE_HP]*v)
		case FIGHT_PROP_DEFENSE_PERCENT:
			computedFightProps.Add(FIGHT_PROP_DEFENSE, fightProps[FIGHT_PROP_BASE_DEFENSE]*v)
		default:
			switch key {
			case FIGHT_PROP_BASE_ATTACK:
				computedFightProps.Add(FIGHT_PROP_ATTACK, v)
			case FIGHT_PROP_BASE_HP:
				computedFightProps.Add(FIGHT_PROP_HP, v)
			case FIGHT_PROP_BASE_DEFENSE:
				computedFightProps.Add(FIGHT_PROP_DEFENSE, v)
			}
			computedFightProps.Add(key, v)
		}
	}

	return computedFightProps
}

func (fightProps FightProps) Round() FightProps {
	rounded := FightProps{}

	for k := range fightProps {
		if isPercentProp(k) {
			rounded[k] = math.Round(fightProps[k]*1000) / 1000
		} else {
			rounded[k] = math.Round(fightProps[k])
		}
	}

	return rounded
}

func (fightProps FightProps) Clone() FightProps {
	cloned := FightProps{}

	for k := range fightProps {
		cloned[k] = fightProps[k]
	}

	return cloned
}

func (fightProps FightProps) Add(k string, v2 float64) {
	if v, ok := fightProps[k]; ok {
		fightProps[k] = v + v2
	} else {
		fightProps[k] = v2
	}
}

func (fightProps FightProps) AddFightProps(addPropsList ...FightProps) {

	for i := range addPropsList {
		addProps := addPropsList[i]

		for k := range addProps {
			if addProps[k] == 0 {
				continue
			}

			fightProps.Add(k, addProps[k])
		}
	}
}

func (fightProps FightProps) Sub(subPropsList ...FightProps) {
	for i := range subPropsList {
		subProps := subPropsList[i]

		for k := range subProps {
			if subProps[k] == 0 {
				continue
			}

			if v, ok := fightProps[k]; ok {
				fightProps[k] = v - subProps[k]
			} else {
				fightProps[k] = subProps[k]
			}
		}
	}
}

func ParseParamName(t string, params ...interface{}) []string {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err, t, params)
		}
	}()

	parts := strings.Split(t, "|")

	for i := range parts {
		parts[i] = funcs.Exec(parts[i], params...)
	}

	return parts
}

var funcs = Formatter(func(v interface{}, fnName string) string {
	var f float64
	switch x := v.(type) {
	case float64:
		f = x
	case int:
		f = float64(x)
	}

	switch fnName {
	case "F1":
		return fmt.Sprintf("%.1f", f)
	case "F2":
		return fmt.Sprintf("%.2f", f)
	case "F1P":
		return fmt.Sprintf("%.1f", f*100) + "%"
	case "F2P":
		return fmt.Sprintf("%.2f", f*100) + "%"
	case "P":
		return fmt.Sprintf("%.0f", f*100) + "%"
	case "F":
		return fmt.Sprintf("%.0f", f)
	case "I":
		return fmt.Sprintf("%.0f", f)
	}
	return ""
})

func SliceFloatToSliceInterface(list []float64) []interface{} {
	ret := make([]interface{}, len(list))
	for i := range ret {
		ret[i] = list[i]
	}
	return ret
}

type Formatter func(v interface{}, fnName string) string

var reParams = regexp.MustCompile(`({[^}]+})`)

func (f Formatter) Exec(t string, params ...interface{}) string {
	return reParams.ReplaceAllStringFunc(t, func(s string) string {
		pipe := strings.Split(strings.TrimSpace(s[1:len(s)-1]), ":")

		var v interface{}

		for i := range pipe {
			if i == 0 {
				idx, _ := strconv.ParseUint(strings.TrimLeft(pipe[0], "param"), 10, 64)
				if idx > 0 {
					v = params[idx-1]
				}
				continue
			}

			if v != nil {
				v = f(v, strings.ToUpper(pipe[i]))
			}
		}

		return fmt.Sprintf("%v", v)
	})
}
