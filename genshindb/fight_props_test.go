package genshindb

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestParseTemplate(t *testing.T) {
	NewWithT(t).Expect(ParseParamName(FightPropParamNames[FIGHT_PROP_HP_PERCENT], 0.1115)).To(Equal([]string{
		"生命值", "11.2%",
	}))
}
