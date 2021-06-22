package main

type (
	base struct {
		Data map[string]interface{}
	}
)

func mainData() *base {
	return &base{
		Data: map[string]interface{}{
			"jsonrpc": "2.0",
			"id":      "",
			"method":  "",
			"params": map[string]interface{}{
				"sdk":                   "custom01",
				"sdkVersion":            "3.0.0",
				"appDefaultLanguage":    "en",
				"userPreferredLanguage": "en",
			},
		},
	}
}

func joinData(applicationId string, sessionToken string, cId string, csId string) *base {
	base := mainData()
	uParams := base.Data["params"].(map[string]interface{})
	uParams["applicationId"] = applicationId
	uParams["sessionToken"] = sessionToken
	uParams["campaignId"] = cId
	uParams["challengeStructureId"] = csId
	uParams["timezone"] = "China Standard Time"
	base.Data["id"] = applicationId
	base.Data["method"] = "mobile.challenges.v2.join"
	return base
}

func progressData(applicationId string, sessionToken string, strTime []string, eventName string) *base {
	base := mainData()
	uParams := base.Data["params"].(map[string]interface{})
	uParams["applicationId"] = applicationId
	uParams["sessionToken"] = sessionToken
	uParams["startedAt"] = strTime[0]
	uParams["endedAt"] = strTime[1]
	uParams["eventName"] = eventName
	uParams["value"] = 1
	uParams["signature"] = getSign(applicationId, sessionToken, eventName, strTime)
	base.Data["id"] = applicationId
	base.Data["method"] = "mobile.challenges.v2.progressEvent"
	return base
}

func currentChallengeListData(applicationId string, sessionToken string) *base {
	base := mainData()
	uParams := base.Data["params"].(map[string]interface{})
	uParams["applicationId"] = applicationId
	uParams["sessionToken"] = sessionToken
	uParams["page"] = 1
	uParams["pageSize"] = 10
	base.Data["id"] = applicationId
	base.Data["method"] = "mobile.challenges.v2.current"
	return base
}

func challengeListData(applicationId string, sessionToken string) *base {
	base := mainData()
	uParams := base.Data["params"].(map[string]interface{})
	uParams["applicationId"] = applicationId
	uParams["sessionToken"] = sessionToken
	uParams["onlyShowEligibleChallenges"] = true
	uParams["page"] = 1
	uParams["pageSize"] = 100
	base.Data["id"] = applicationId
	base.Data["method"] = "mobile.challenges.v4.list"
	return base
}

func handshakeData(applicationId string, userToken string) *base {
	base := mainData()
	uParams := base.Data["params"].(map[string]interface{})
	uParams["applicationId"] = applicationId
	uParams["sessionToken"] = nil
	uParams["userToken"] = userToken
	uParams["Birthdate"] = "2001-01-05"
	base.Data["id"] = applicationId
	base.Data["method"] = "mobile.accounts.v1.handshake"
	return base
}

func startData(applicationId string, accountToken string, HpidUserId string) *base {
	base := mainData()
	base.Data["id"] = applicationId
	base.Data["method"] = "mobile.sessions.v2.start"
	uParams := base.Data["params"].(map[string]interface{})
	uParams["accountToken"] = accountToken
	uParams["applicationId"] = applicationId
	uParams["externalPlayerId"] = HpidUserId
	uParams["eventNames"] = []string{
		"PLAY:OVERWATCH",
		"PLAY:HEROES_OF_THE_STORM",
		"PLAY:FORTNITE",
		"PLAY:THE_DIVISION",
		"PLAY:THE_DIVISION_2",
		"PLAY:PUBG",
		"PLAY:APEX_LEGENDS",
		"PLAY:CS_GO",
		"PLAY:LEAGUE_OF_LEGENDS",
		"PLAY:DOTA_2",
		"PLAY:SMITE",
		"PLAY:AGE_OF_EMPIRES_2",
		"PLAY:STARCRAFT_2",
		"PLAY:COMPANY_OF_HEROES_2",
		"PLAY:ASSASSINS_CREED_ODYSSEY",
		"PLAY:WORLD_OF_WARCRAFT",
		"PLAY:WORLD_OF_WARCRAFT_CLASSIC",
		"PLAY:SPOTIFY",
		"PLAY:RINGS_OF_ELYSIUM",
		"PLAY:HEARTHSTONE",
		"PLAY:GARRYS_MOD",
		"PLAY:GOLF_IT",
		"PLAY:DECEIT",
		"PLAY:SEVEN_DAYS_TO_DIE",
		"PLAY:DOOM_ETERNAL",
		"PLAY:STARWARS_JEDI_FALLEN_ORDER",
		"PLAY:MINECRAFT",
		"PLAY:DEAD_BY_DAYLIGHT",
		"PLAY:NETFLIX",
		"PLAY:HULU",
		"PLAY:PATH_OF_EXILE",
		"PLAY:WARTHUNDER",
		"PLAY:CALL_OF_DUTY_MODERN_WARFARE",
		"PLAY:ROCKET_LEAGUE",
		"PLAY:NBA_2K20",
		"PLAY:STREET_FIGHTER_V",
		"PLAY:DRAGON_BALL_FIGHTER_Z",
		"PLAY:GEARS_OF_WAR_5",
		"PLAY:FIFA_20",
		"PLAY:MASTER_CHIEF_COLLECTION",
		"PLAY:RAINBOW_SIX",
		"PLAY:UPLAY",
		"PLAY:ROBLOX",
		"VERSUS_GAME_API:TEAMFIGHT_TACTICS:GOLD_LEFT",
		"VERSUS_GAME_API:TEAMFIGHT_TACTICS:TIME_ELIMINATED",
		"VERSUS_GAME_API:TEAMFIGHT_TACTICS:THIRD_PLACE_OR_HIGHER",
		"VERSUS_GAME_API:TEAMFIGHT_TACTICS:SECOND_PLACE_OR_HIGHER",
		"VERSUS_GAME_API:TEAMFIGHT_TACTICS:PLAYERS_ELIMINATED",
		"VERSUS_GAME_API:TEAMFIGHT_TACTICS:TOTAL_DAMAGE_TO_PLAYERS",
		"PLAY:MONSTER_HUNTER_WORLD",
		"PLAY:WARFRAME",
		"PLAY:LEGENDS_OF_RUNETERRA",
		"PLAY:VALORANT",
		"PLAY:CROSSFIRE",
		"PLAY:PALADINS",
		"PLAY:TROVE",
		"PLAY:RIFT",
		"PLAY:ARCHEAGE",
		"PLAY:IRONSIGHT",
		"GAMIGO_PLACEHOLDER",
		"PLAY:TWINSAGA",
		"PLAY:AURA_KINGDOM",
		"PLAY:SHAIYA",
		"PLAY:SOLITAIRE",
		"PLAY:TONY_HAWK",
		"PLAY:AVENGERS",
		"PLAY:FALL_GUYS",
		"PLAY:QQ_SPEED",
		"PLAY:FIFA_ONLINE_3",
		"PLAY:NBA2KOL2",
		"PLAY:DESTINY2",
		"PLAY:AMONG_US",
		"PLAY:MAPLE_STORY",
		"PLAY:ASSASSINS_CREED_VALHALLA",
		"PLAY:FREESTYLE_STREET_BASKETBALL",
		"PLAY:CRAZY_RACING_KART_RIDER",
		"PLAY:COD_BLACK_OPS_COLD_WAR",
		"PLAY:CYBERPUNK_2077",
		"PLAY:HADES",
		"PLAY:RUST",
		"PLAY:GENSHIN_IMPACT",
		"PLAY:ESCAPE_FROM_TARKOV",
		"PLAY:RED_DEAD_REDEMPTION_2",
		"PLAY:CIVILIZATION_VI",
		"PLAY:VALHEIM",
		"PLAY:FINAL_FANTASY_XIV",
		"PLAY:OASIS",
		"PLAY:CASTLE_CRASHERS",
		"PLAY:GANG_BEASTS",
		"PLAY:SPEEDRUNNERS",
		"PLAY:OVERCOOKED_2",
		"PLAY:OVERCOOKED_ALL_YOU_CAN_EAT",
		"Launch OMEN Command Center",
		"Use OMEN Command Center",
		"OMEN Command Center Macro Created",
		"OMEN Command Center Macro Assigned",
		"Mindframe Adjust Cooling Option",
		"Connect 2 different OMEN accessories to your PC at the same time",
		"Use Omen Reactor",
		"Use Omen Photon",
		"Launch Game From GameLauncher",
		"Image like From ImageGallery",
		"Set as background From ImageGallery",
		"Download image From ImageGallery",
		"overwatch",
		"heroesofthestorm",
		"heroesofthestorm_x64",
		"FortniteClient-Win64-Shipping",
		"FortniteClient-Win64-Shipping_BE",
		"thedivision",
		"thedivision2",
		"TslGame",
		"r5apex",
		"csgo",
		"League of Legends",
		"dota2",
		"smite",
		"AoE2DE_s",
		"AoK HD",
		"AoE2DE",
		"sc2",
		"s2_x64",
		"RelicCoH2",
		"acodyssey",
		"wow",
		"wow64",
		"wow_classic",
		"wowclassic",
		"Spotify",
		"Europa_client",
		"hearthstone",
		"hl2",
		"GolfIt-Win64-Shipping",
		"GolfIt",
		"Deceit",
		"7DaysToDie",
		"DoomEternal_temp",
		"starwarsjedifallenorder",
		"Minecraft.Windows",
		"net.minecraft.client.main.Main",
		"DeadByDaylight-Win64-Shipping",
		"4DF9E0F8.Netflix",
		"HuluLLC.HuluPlus",
		"PathOfExileSteam",
		"PathOfExile_x64Steam",
		"aces",
		"modernwarfare",
		"RocketLeague",
		"NBA2K20",
		"StreetFighterV",
		"RED-Win64-Shipping",
		"Gears5",
		"fifa20",
		"MCC-Win64-Shipping",
		"MCC-Win64-Shipping-WinStore",
		"RainbowSix",
		"RainbowSix_BE",
		"RainbowSix_Vulkan",
		"upc",
		"ROBLOXCORPORATION.ROBLOX",
		"RobloxPlayerBeta",
		"VERSUS_GAME_API_TEAMFIGHT_TACTICS_GOLD_LEFT",
		"VERSUS_GAME_API_TEAMFIGHT_TACTICS_TIME_ELIMINATED",
		"VERSUS_GAME_API_TEAMFIGHT_TACTICS_THIRD_PLACE_OR_HIGHER",
		"VERSUS_GAME_API_TEAMFIGHT_TACTICS_SECOND_PLACE_OR_HIGHER",
		"VERSUS_GAME_API_TEAMFIGHT_TACTICS_PLAYERS_ELIMINATED",
		"VERSUS_GAME_API_TEAMFIGHT_TACTICS_TOTAL_DAMAGE_TO_PLAYERS",
		"MonsterHunterWorld",
		"Warframe.x64",
		"lor",
		"valorant-Win64-shipping",
		"valorant",
		"crossfire",
		"Paladins",
		"trove",
		"rift_64",
		"rift_x64",
		"archeage",
		"ironsight",
		"Game",
		"Game.bin",
		"glyph_twinsaga",
		"glyph_aurakingdom",
		"glyph_shaiya",
		"Solitaire",
		"THPS12",
		"avengers",
		"Fallguys_client_game",
		"GameApp",
		"fifazf",
		"NBA2KOL2",
		"destiny2",
		"Among Us",
		"MapleStory",
		"ACValhalla",
		"FreeStyle",
		"KartRider",
		"BlackOpsColdWar",
		"Cyberpunk2077",
		"Hades",
		"RustClient",
		"GenshinImpact",
		"EscapeFromTarkov",
		"EscapeFromTarkov_BE",
		"RDR2",
		"CivilizationVI",
		"valheim",
		"ffxiv_dx11",
		"AD2F1837.OMENSpectate",
		"castle",
		"Gang Beasts",
		"SpeedRunners",
		"Overcooked2",
		"Overcooked All You Can Eat",
	}
	latitude := 30.5832367
	longitude := 103.982384
	res, _ := Requests().Get("https://api.bilibili.com/x/web-interface/zone")
	var zoneRet map[string]interface{}
	res.Json(&zoneRet)
	latitude = zoneRet["data"].(map[string]float64)["latitude"]
	longitude = zoneRet["data"].(map[string]float64)["longitude"]
	uParams["location"] = map[string]float64{
		"latitude":  latitude,
		"longitude": longitude,
	}
	return base
}
