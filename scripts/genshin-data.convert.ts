import {groupBy, mapValues, reduce, Dictionary, includes} from "lodash"
import textMapCHS from "../GenshinData/TextMap/TextMapCHS.json"
import materialExcelConfigData from "../GenshinData/ExcelBinOutput/MaterialExcelConfigData.json"
import proudSkillExcelConfigData from "../GenshinData/ExcelBinOutput/ProudSkillExcelConfigData.json"
import equipAffixExcelConfigData from "../GenshinData/ExcelBinOutput/EquipAffixExcelConfigData.json"

// Avatar
import avatarExcelConfigData from "../GenshinData/ExcelBinOutput/AvatarExcelConfigData.json"
import avatarPromoteExcelConfigData from "../GenshinData/ExcelBinOutput/AvatarPromoteExcelConfigData.json"
import avatarTalentExcelConfigData from "../GenshinData/ExcelBinOutput/AvatarTalentExcelConfigData.json"
import avatarSkillDepotExcelConfigData from "../GenshinData/ExcelBinOutput/AvatarSkillDepotExcelConfigData.json"
import avatarCurveExcelConfigData from "../GenshinData/ExcelBinOutput/AvatarCurveExcelConfigData.json"
import avatarSkillExcelConfigData from "../GenshinData/ExcelBinOutput/AvatarSkillExcelConfigData.json"

// Weapon
import weaponExcelConfigData from "../GenshinData/ExcelBinOutput/WeaponExcelConfigData.json"
import weaponPromoteExcelConfigData from "../GenshinData/ExcelBinOutput/WeaponPromoteExcelConfigData.json"
import weaponCurveExcelConfigData from "../GenshinData/ExcelBinOutput/WeaponCurveExcelConfigData.json"

// Reliquary
import reliquaryExcelConfigData from "../GenshinData/ExcelBinOutput/ReliquaryExcelConfigData.json"
import reliquaryMainPropExcelConfigData from "../GenshinData/ExcelBinOutput/ReliquaryMainPropExcelConfigData.json"
import reliquaryAffixExcelConfigData from "../GenshinData/ExcelBinOutput/ReliquaryAffixExcelConfigData.json"
import reliquarySetExcelConfigData from "../GenshinData/ExcelBinOutput/ReliquarySetExcelConfigData.json"

// LevelExp
import avatarLevelExcelConfigData from "../GenshinData/ExcelBinOutput/AvatarLevelExcelConfigData.json"
import weaponLevelExcelConfigData from "../GenshinData/ExcelBinOutput/WeaponLevelExcelConfigData.json"
import reliquaryLevelExcelConfigData from "../GenshinData/ExcelBinOutput/ReliquaryLevelExcelConfigData.json"

import {writeFileSync} from "fs";

const text = (hash: number): string => {
    return (textMapCHS as any)[hash]
}

const material = (id: number) => {
    const m = materialExcelConfigData.find((a) => a.Id == id)!

    return ({
        Name: text(m.NameTextMapHash),
        RankLevel: m.RankLevel,
    })
}


const addPropSet = (addProps: { PropType?: string, Value?: number }[]) => {
    return addProps.filter((p) => p.PropType).reduce((r: any, p: any) => ({
        ...r,
        [p.PropType]: p.Value || 0
    }), {})
}

const avatarPropGrowCurve = (growCurve: string) => {
    return avatarCurveExcelConfigData
        .filter(a => a.Level <= 90)
        .map((a) => {
            return a.CurveInfos.find((i) => i.Type == growCurve)?.Value || 0
        })
}

const weaponPropGrowCurve = (growCurve: string) => {
    return weaponCurveExcelConfigData
        .filter(a => a.Level <= 90)
        .map((a) => {
            return a.CurveInfos.find((i) => i.Type == growCurve)?.Value || 0
        })
}


const talent = (talentId: number) => {
    const talent = avatarTalentExcelConfigData.find(a => a.TalentId == talentId)
    if (talent) {
        return {
            Name: text(talent.NameTextMapHash),
            Desc: text(talent.DescTextMapHash),
            AddProps: addPropSet(talent.AddProps),
        }
    }
    return null
}


const proudSkills = (proudSkillGroupId: number) => {
    return proudSkillExcelConfigData
        .filter((skill) => skill.ProudSkillGroupId == proudSkillGroupId)
        .map((skill) => ({
            Name: text(skill.NameTextMapHash),
            Desc: text(skill.DescTextMapHash),
            BreakLevel: skill.BreakLevel || 0,
            ParamNames: skill.ParamDescList.map(text).filter(s => s),
            Params: skill.ParamList
        }))
}

const skill = (skillID: number) => {
    const skill = avatarSkillExcelConfigData.find(a => a.Id == skillID)

    if (skill) {
        const proudSkills = proudSkillExcelConfigData
            .filter((item) => item.ProudSkillGroupId == skill.ProudSkillGroupId)

        const paramNames = proudSkills.length > 0 ? proudSkills[0].ParamDescList.map(text).filter(s => s) : []

        return {
            Name: text(skill.NameTextMapHash),
            Desc: text(skill.DescTextMapHash),
            CdTime: skill.CdTime,
            CostElemType: skill.CostElemType,
            CostElemVal: skill.CostElemVal,
            BreakLevels: proudSkills.map((item) => item.BreakLevel || 0),
            CoinCosts: proudSkills.map((item) => item.CoinCost || 0),
            MaterialCosts: proudSkills.map((item) => item.CostItems
                .filter((item: any) => item.Id)
                .map((item: any) => ({
                    ...material(item.Id),
                    Count: item.Count,
                }))),
            ProudSkills: paramNames.length > 0 ? {
                ParamNames: paramNames,
                Params: proudSkills.map((item) => item.ParamList),
            } : {}
        }
    }

    return null
}

const avatarSkillAndTalents = (skillDepotId: number) => {
    const skillDepot = avatarSkillDepotExcelConfigData.find(a => a.Id == skillDepotId)!

    const ret = {
        ElementType: "" as any,
        Skills: [...skillDepot.Skills.filter(id => id > 0), skillDepot.EnergySkill as number]
            .map((id) => skill(id)).filter(s => s),
        InherentSkills: skillDepot.InherentProudSkillOpens
            .filter(s => (s as any).ProudSkillGroupId)
            .map((s: any) => proudSkills(s.ProudSkillGroupId)).flat(),
        Talents: skillDepot.Talents.map((id) => talent(id))
    }


    ret.Skills.forEach((s: any) => {
        const matched = s.Desc.match(/([冰水火雷岩风草])元素/)
        if (matched) {
            ret.ElementType = matched[1]
        }
    })


    return ret
}

const avatarPromotes = (avatarPromoteId: number) => {
    const p = avatarPromoteExcelConfigData
        .filter((a) => (a.AvatarPromoteId == avatarPromoteId))

    return p.map((a, i) => ({
        MinLevel: i > 0 ? p[i - 1]?.UnlockMaxLevel : 1,
        UnlockMaxLevel: a.UnlockMaxLevel,
        CoinCost: a.ScoinCost || 0,
        MaterialCosts: a.CostItems
            .filter((item: any) => item.Id)
            .map((item: any) => ({
                ...material(item.Id),
                Count: item.Count,
            })),
        AddProps: addPropSet(a.AddProps),
    }))
}


const weaponPromotes = (weaponPromoteId: number) => {
    const p = weaponPromoteExcelConfigData
        .filter((a) => (a.WeaponPromoteId == weaponPromoteId))

    return p.map((a, i) => ({
        MinLevel: i > 0 ? p[i - 1]?.UnlockMaxLevel : 1,
        UnlockMaxLevel: a.UnlockMaxLevel,
        CoinCost: a.CoinCost || 0,
        MaterialCosts: a.CostItems
            .filter((item: any) => item.Id)
            .map((item: any) => ({
                ...material(item.Id),
                Count: item.Count,
            })),
        AddProps: addPropSet(a.AddProps),
    }))
}


const propGrowCurveFor = (fn: (type: string) => number[]) => (type: string, base: number,) => {
    return {
        Base: base,
        Values: fn(type)
    }
}


const avatars = () => avatarExcelConfigData
    .filter(a => a.UseType && text(a.DescTextMapHash))
    .map((a) => {
        return {
            Id: a.Id,
            Name: text(a.NameTextMapHash),
            Desc: text(a.InfoDescTextMapHash),

            RankLevel: a.QualityType === "QUALITY_ORANGE" ? 5 : 4,
            WeaponType: a.WeaponType,

            ChargeEfficiency: a.ChargeEfficiency,
            StaminaRecoverSpeed: a.StaminaRecoverSpeed,
            Critical: a.Critical,
            CriticalHurt: a.CriticalHurt,

            ...avatarSkillAndTalents(a.SkillDepotId),

            ...({
                Promotes: avatarPromotes(a.AvatarPromoteId),
                PropGrowCurves: a.PropGrowCurves.reduce((ret, curve) => ({
                    ...ret,
                    [curve.Type]: propGrowCurveFor(avatarPropGrowCurve)(curve.GrowCurve, {
                        FIGHT_PROP_BASE_HP: a.HpBase,
                        FIGHT_PROP_BASE_ATTACK: a.AttackBase,
                        FIGHT_PROP_BASE_DEFENSE: a.DefenseBase,
                    }[curve.Type] || 0)
                }), {}),
            })
        }
    })

const equipAffixes = (id: number) => {
    return equipAffixExcelConfigData.filter((a) => a.Id == id).map((skillAffix) => ({
        Name: text(skillAffix.NameTextMapHash),
        Desc: text(skillAffix.DescTextMapHash),
        Level: skillAffix.Level,
        AddProps: addPropSet(skillAffix.AddProps),
        ParamList: skillAffix.ParamList
    }))
}

const weapons = () => weaponExcelConfigData
    .map((a) => {
        return {
            Name: text(a.NameTextMapHash),
            Desc: text(a.DescTextMapHash),
            RankLevel: a.RankLevel,
            Promotes: weaponPromotes(a.WeaponPromoteId),
            Affixes: a.SkillAffix
                .filter((v) => v !== 0).map((id) => equipAffixes(id)).filter((v) => !!v),
            PropGrowCurves: (a.WeaponProp as any || [])
                .filter((a: any) => a.PropType)
                .reduce((ret: any, prop: any) => ({
                    ...ret,
                    [prop.PropType]: propGrowCurveFor(weaponPropGrowCurve)(prop.Type, prop.InitValue)
                }), {}),
        }
    })

const reliquarySets = () => {
    return reliquarySetExcelConfigData.reduce((ret, reliquarySet) => ({
        ...ret,
        [reliquarySet.SetId]: (() => {
            const affixes = equipAffixes(reliquarySet.EquipAffixId || 0)
            if (affixes.length > 0) {
                return {
                    Name: affixes[0].Name,
                    EquipAffixes: reliquarySet.SetNeedNum.map((v, i) => ({
                        NeedNum: v,
                        Desc: affixes[i].Desc,
                        AddProps: affixes[i].AddProps,
                        ParamList: affixes[i].ParamList,
                    }))
                }
            }
            return {}
        })(),
    }), {})
}

const reliquaryAffixDepots = () => mapValues(
    groupBy(reliquaryAffixExcelConfigData, (a) => a.DepotId),
    (reliquaryAffixes) => reduce(reliquaryAffixes, (ret, c) => ({
        ...ret,
        [c.PropType]: [...(ret[c.PropType] || []), c.PropValue],
    }), {} as Record<string, number[]>)
)

const reliquaries = () => ({
    ReliquaryAddProps: reliquaryAddProps(),
    ReliquaryAffixDepots: reliquaryAffixDepots(),
    ReliquarySets: reliquarySets(),
    Reliquaries: reliquaryExcelConfigData
        .map((a) => ({
            Name: text(a.NameTextMapHash),
            RankLevel: a.RankLevel,
            Desc: text(a.DescTextMapHash),
            EquipType: a.EquipType,
            MaxLevel: a.MaxLevel,
            AddPropLevels: a.AddPropLevels,
            AppendPropNum: a.AppendPropNum,
            MainPropTypes: reliquaryMainPropExcelConfigData
                .filter((c) => c.PropDepotId == a.MainPropDepotId)
                .map((c) => c.PropType),
            AppendPropDepotId: a.AppendPropDepotId,
            SetId: a.SetId,
        }))
})

const levelExpAvatar = () => avatarLevelExcelConfigData.map((v) => v.Exp)

const levelExpWeapon = () => new Array(5).fill(0).map((_, i) => {
    return weaponLevelExcelConfigData.map((v) => v.RequiredExps[i])
})

const levelExpReliquary = () => new Array(5).fill(0).map((_, idx) => reliquaryLevelExcelConfigData
    .filter((a) => a.Rank == idx + 1)
    .map((v) => v.Exp || 0))


const levelExps = () => {
    return {
        Avatar: levelExpAvatar(),
        Weapon: levelExpWeapon(),
        Reliquary: levelExpReliquary(),
    }
}

const reliquaryAddProps = () => {
    return new Array(5).fill(0).map((_, idx) => reliquaryLevelExcelConfigData
        .filter((a) => a.Rank == idx + 1)
        .map((v) => addPropSet(v.AddProps))
    )
}

writeFileSync("./genshindb/avatars.json", JSON.stringify(avatars(), null, 2))
writeFileSync("./genshindb/weapons.json", JSON.stringify(weapons(), null, 2))
writeFileSync("./genshindb/reliquaries.json", JSON.stringify(reliquaries(), null, 2))
writeFileSync("./genshindb/level_exps.json", JSON.stringify(levelExps(), null, 2))


