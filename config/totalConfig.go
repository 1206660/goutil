package config

import (
	"github.com/fsnotify/fsnotify"
	"configPath"
	"github.com/name5566/leaf/log"
	"github.com/name5566/leaf/util"
	"strings"
)

var (
	//StringInfoMgr            ConfigMap
	HeroExpInfoMgr            ConfigMap
	DropInfoMgr               ConfigMap
	GameInfoMgr               ConfigMap
	PlayerExpInfoMgr          ConfigMap
	ItemCreateInfoMgr         ConfigMap
	EquipCreateInfoMgr        ConfigMap
	SkillInfoMgr              ConfigMap
	SkillAdvanceInfoMgr       ConfigMap
	HeroCreateInfoMgr         ConfigMap
	HeroLevelupInfoMgr        ConfigMap
	HeroAwakeInfoMgr          ConfigMap
	EquipCompositeInfoMgr     ConfigMap
	FragmentCreateInfoMgr     ConfigMap
	FragmentCompositeInfoMgr  ConfigMap
	EquipSetInfoMgr           ConfigMap
	MessageInfoMgr            ConfigMap
	QuestInfoMgr              ConfigMap
	EventInfoMgr              ConfigMap
	AchievementInfoMgr        ConfigMap
	RouletteInfoMgr           ConfigMap
	HighRouletteInfoMgr       ConfigMap
	NpcDataInfoMgr            ConfigMap
	NpcLevelupInfoMgr         ConfigMap
	ClimbingUpMgr             ConfigMap
	HHeroLotteryMgr           ConfigMap
	LHeroLotteryMgr           ConfigMap
	FHeroLotteryMgr           ConfigMap
	VipMgr                    ConfigMap
	ShopInfoMgr               ConfigMap
	SigninInfoMgr             ConfigMap
	PayInfoMgr                ConfigMap
	DecompositionMgr          ConfigMap
	VipInfoMgr                ConfigMap
	TargetInfoMgr             ConfigMap
	ActionInfoMgr             ConfigMap
	BuffInfoMgr               ConfigMap
	FriendSearchInfoMgr       ConfigMap
	PrefixNameInfoMgr         ConfigMap
	SuffixNameInfoMgr         ConfigMap
	GuildInfoMgr              ConfigMap
	GuildTalentInfoMgr        ConfigMap
	GuildBossInfoMgr          ConfigMap
	ClimbingDownInfoMgr       ConfigMap
	TacticsInfoMgr            ConfigMap
	OnlineRewardInfoMgr       ConfigMap
	GamePlayInfoMgr           ConfigMap
	ActivityTreasureMgr       ConfigMap
	ServerRuntimeInfoMgr      ConfigMap
	ProjectHeroLotteryMgr     ConfigMap
	ReplacementHeroLotteryMgr ConfigMap
	GrowthFundInfoMgr         ConfigMap
)

// 管理文件名对应管理器
var __loadMap map[string]*ConfigMap
var basePath string
var __reloadMap util.Map

const ITEM_EQUIP_SPLIT_RANGE = 10000
const EQUIP_FRAGMENT_SPLIT_RANGE = 20000

func init() {
	OnInit()
}

func OnInit() {
	// 为了文件可以重新加载，所以改为这样
	__loadMap = map[string]*ConfigMap{
		//"String.txt":              &StringInfoMgr,
		"HeroExpInfo.txt":            &HeroExpInfoMgr,
		"DropInfo.txt":               &DropInfoMgr,
		"GameInfo.txt":               &GameInfoMgr,
		"PlayerExpInfo.txt":          &PlayerExpInfoMgr,
		"ItemCreateInfo.txt":         &ItemCreateInfoMgr,
		"EquipCreateInfo.txt":        &EquipCreateInfoMgr,
		"SkillInfo.txt":              &SkillInfoMgr,
		"SkillAdvance.txt":           &SkillAdvanceInfoMgr,
		"HeroCreateInfo.txt":         &HeroCreateInfoMgr,
		"HeroLevelupInfo.txt":        &HeroLevelupInfoMgr,
		"HeroAwakeInfo.txt":          &HeroAwakeInfoMgr,
		"EquipCompositeInfo.txt":     &EquipCompositeInfoMgr,
		"FragmentCreateInfo.txt":     &FragmentCreateInfoMgr,
		"FragmentCompositeInfo.txt":  &FragmentCompositeInfoMgr,
		"EquipSetInfo.txt":           &EquipSetInfoMgr,
		"Message.txt":                &MessageInfoMgr,
		"QuestInfo.txt":              &QuestInfoMgr,
		"EventInfo.txt":              &EventInfoMgr,
		"AchievementInfo.txt":        &AchievementInfoMgr,
		"RouletteInfo1.txt":          &RouletteInfoMgr,
		"RouletteInfo2.txt":          &HighRouletteInfoMgr,
		"NpcDataInfo.txt":            &NpcDataInfoMgr,
		"NpcLevelupInfo.txt":         &NpcLevelupInfoMgr,
		"ClimbingUp.txt":             &ClimbingUpMgr,
		"HHeroLottery.txt":           &HHeroLotteryMgr,
		"LHeroLottery.txt":           &LHeroLotteryMgr,
		"FHeroLottery.txt":           &FHeroLotteryMgr,
		"Vip.txt":                    &VipMgr,
		"ShopInfo.txt":               &ShopInfoMgr,
		"SigninInfo.txt":             &SigninInfoMgr,
		"Pay.txt":                    &PayInfoMgr,
		"HeroDecompose.txt":          &DecompositionMgr,
		"VIPInfo.txt":                &VipInfoMgr,
		"TargetInfo.txt":             &TargetInfoMgr,
		"ActionInfo.txt":             &ActionInfoMgr,
		"BuffInfo.txt":               &BuffInfoMgr,
		"FriendSearchInfo.txt":       &FriendSearchInfoMgr,
		"PrefixNameInfo.txt":         &PrefixNameInfoMgr,
		"SuffixNameInfo.txt":         &SuffixNameInfoMgr,
		"GuildInfo.txt":              &GuildInfoMgr,
		"GuildTalentInfo.txt":        &GuildTalentInfoMgr,
		"GuildBossInfo.txt":          &GuildBossInfoMgr,
		"ClimbingDown.txt":           &ClimbingDownInfoMgr,
		"TacticsInfo.txt":            &TacticsInfoMgr,
		"OnlineRewardInfo.txt":       &OnlineRewardInfoMgr,
		"GamePlayInfo.txt":           &GamePlayInfoMgr,
		"ActivityTreasure.txt":       &ActivityTreasureMgr,
		"ServerRuntimeInfo.txt":      &ServerRuntimeInfoMgr,
		"ProjectHeroLottery.txt":     &ProjectHeroLotteryMgr,
		"ReplacementHeroLottery.txt": &ReplacementHeroLotteryMgr,
		"GrowthFund.txt":             &GrowthFundInfoMgr,
	}

	basePath = configPath.GetConfDir() + "excel/"
	for k, v := range __loadMap {
		if v.LoadConfigByFilePath(basePath+k) == false {
			log.Fatal("LoadConfigByFilePath %v fail", basePath+k)
		}
	}
	CheckAllInfoValied()
	go StartWatchDog(basePath)
	log.Debug("ConfigMgr OnInit Success %v", basePath)
}

func StartWatchDog(watchPath string) {
	//创建一个监控对象
	watch, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err.Error())
	}
	defer watch.Close()

	err = watch.Add(watchPath)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Debug("StartWatchDog Start %v", watchPath)

	for {
		select {
		case ev := <-watch.Events:
			{
				for k, mgr := range __loadMap {
					found := strings.Contains(ev.Name, k)
					if found == false {
						continue
					}
					__reloadMap.Set(mgr, true)
					log.Debug("StartWatchDog %v Contain:%v", ev.Name, k)
					break
				}
			}
		case err := <-watch.Errors:
			{
				log.Debug("error : ", err)
				return
			}
		}
	}
}

//func G_Msg(id int) string {
//	pInfo := MessageInfoMgr.GetInfo(id)
//	if pInfo != nil {
//		return pInfo["String"].GetAsString() + "(" + strconv.Itoa(id) + ")"
//	}
//	return "null" + "(" + strconv.Itoa(id) + ")"
//}

func G_GetItemInfo(id int) map[string]ConfigCell {
	if id < ITEM_EQUIP_SPLIT_RANGE {
		return ItemCreateInfoMgr.GetInfo(id)
	} else if id < EQUIP_FRAGMENT_SPLIT_RANGE {
		return EquipCreateInfoMgr.GetInfo(id)
	} else {
		return FragmentCompositeInfoMgr.GetInfo(id)
	}
	return nil
}

func CheckAllReloadConfigMgr() {
	__reloadMap.LockRange(func(k interface{}, v interface{}) {
		configMap := k.(*ConfigMap)
		configMap.Reload()
	})
	__reloadMap = util.Map{}
}

//func G_Str(id int) string {
//	pInfo := StringInfoMgr.GetInfo(id)
//	if pInfo != nil {
//		return pInfo["CN"].GetAsString()
//	}
//	return "null" + "(" + strconv.Itoa(id) + ")"
//}

func CheckAllInfoValied() {
	CheckNpcDataInfoValied()
}

// 配置表关联性检查
func CheckNpcDataInfoValied() {
	allInfos := NpcDataInfoMgr.GetAllInfo()
	for k, v := range allInfos {
		iHeroID := v["iHeroID"].GetAsInt()
		pHeroInfo := HeroCreateInfoMgr.GetInfo(iHeroID)
		if pHeroInfo == nil {
			log.Error("CheckNpcDataInfoValied failed NpcID:%v HeroID:%v", k, iHeroID)
		}
	}
}
