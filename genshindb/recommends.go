package genshindb

type RecommendSet map[string]RecommendReliquaryProps

type RecommendReliquaryProps []string

func (r RecommendSet) For(name string, equipType string) (mainPropType string, affixPropTypes []string) {

	mainPropType = FIGHT_PROP_ATTACK_PERCENT

	switch pos := posForEquipType(equipType); pos {
	case 1:
		mainPropType = FIGHT_PROP_HP
	case 2:
		mainPropType = FIGHT_PROP_ATTACK
	default:
		if ret, ok := r[name]; ok {
			if len(ret) > 3 {
				mainPropType = ret[pos-3]
			}
		}
	}

	if ret, ok := r[name]; ok {
		if len(ret) > 3 {
			for _, p := range ret[3:] {
				if p != mainPropType {
					affixPropTypes = append(affixPropTypes, p)
				}
			}
		}
	}

	return
}
