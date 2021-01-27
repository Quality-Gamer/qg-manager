package model

import (
	"math/rand"
	"github.com/Quality-Gamer/qg-manager/conf"
	"github.com/Quality-Gamer/qg-manager/database"
	"strconv"
)

type GameModel struct {
	Id                 string              `json:"id"`
	Name               string              `json:"name"`
	Description        string              `json:"description"`
	InitialTime        int                 `json:"initial_time"`
	InitialMoney       float64             `json:"initial_money"`
	Type               string              `json:"type"`
	BelongsTo          int                 `json:"user_id"`
	Levels             []Level             `json:"levels"`
	ManagerOccurrences []ManagerOccurrence `json:"manager_occurrences"`
	UserOccurrences    []UserOccurrence    `json:"user_occurrences"`
}

type Level struct {
	Id      string    `json:"id"`
	Level   string    `json:"level"`
	Process []Process `json:"process"`
}

type Process struct {
	Id         string     `json:"id"`
	Name       string     `json:"name"`
	Activities []Activity `json:"activities"`
	Resources  []Resource `json:"resources"`
	Score      int        `json:"score"`
}

type Score struct {
	Id    string `json:"id"`
	Score string `json:"score"`
}

type Activity struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
	Time     int    `json:"time"`
	Score    int    `json:"score"`
}

type Resource struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
	Price    int    `json:"price"`
	Type     string `json:"type"`
	Score    int    `json:"score"`
}

type SolveOccurrence struct {
	Id         string     `json:"id"`
	Activities []Activity `json:"activities"`
	Resources  []Resource `json:"resources"`
}

func CreateGameModel(gm GameModel) {
	database.SetKey(conf.GetGameModelKey(gm.Id, conf.Identifier), gm.Id)
	database.SetKey(conf.GetGameModelKey(gm.Id, conf.Name), gm.Name)
	database.SetKey(conf.GetGameModelKey(gm.Id, conf.Description), gm.Description)
	database.SetKey(conf.GetGameModelKey(gm.Id, conf.InitialTime), gm.InitialTime)
	database.SetKey(conf.GetGameModelKey(gm.Id, conf.InitialMoney), gm.InitialMoney)
	database.SetKey(conf.GetGameModelKey(gm.Id, conf.Type), gm.Type)
	database.SetKey(conf.GetGameModelKey(gm.Id, conf.BelongsTo), gm.BelongsTo)
	//TODO: save all struct array

	for key, value := range gm.Levels {
		k := strconv.Itoa(key)
		database.SetKey(conf.GetLevelKey(gm.Id, k,key), value.Id)
		CreateLevel(gm, value, key)
	}

	for key, value := range gm.ManagerOccurrences {
		k := conf.Manager + ":" + strconv.Itoa(key)
		database.SetKey(conf.GetGameOccurrenceKey(gm.Id, k), value.Id)
	}

	for key, value := range gm.UserOccurrences {
		k := conf.UserOccurrence + ":" + strconv.Itoa(key)
		database.SetKey(conf.GetGameOccurrenceKey(gm.Id, k), value.Id)
	}
}

func GetGameModel(id string) GameModel {
	var gm GameModel
	gm.Id = database.GetKey(conf.GetGameModelKey(id, conf.Identifier))
	gm.Name = database.GetKey(conf.GetGameModelKey(id, conf.Name))
	gm.Description = database.GetKey(conf.GetGameModelKey(id, conf.Description))
	it,_ := strconv.Atoi(database.GetKey(conf.GetGameModelKey(id, conf.InitialTime)))
	gm.InitialTime = it
	im,_ := strconv.ParseFloat(database.GetKey(conf.GetGameModelKey(id, conf.InitialMoney)),64)
	gm.InitialMoney = im
	gm.Type = database.GetKey(conf.GetGameModelKey(id, conf.Type))
	bt,_ := strconv.Atoi(database.GetKey(conf.GetGameModelKey(id, conf.BelongsTo)))
	gm.BelongsTo = bt
	//TODO: save all struct array

	gm.Levels = GetLevels(id)
	gm.UserOccurrences,gm.ManagerOccurrences = GetOcurrences(gm)

	return gm
}

func CreateOcurrences(gm GameModel) {
	keyNoIndex := conf.GetGameOccurrenceKey(gm.Id,conf.Count)
	for k, value := range gm.UserOccurrences {
		key := conf.GetGameOccurrenceKeyType(gm.Id,conf.UserOccurrence,strconv.Itoa(k))
		hashString := strconv.Itoa(rand.Int()) + gm.Id + conf.UserOccurrence
		hashId := GerenateHashByString(hashString)
		database.SetKey(key + ":" + conf.Identifier,hashId)
		database.SetKey(key + ":" + conf.Level,value.Level)

		hashString2 := strconv.Itoa(rand.Int()) + gm.Id + conf.Occurrence
		hashId2 := GerenateHashByString(hashString2)
		database.SetKey(key + ":" + conf.Occurrence + ":" + conf.Identifier ,hashId2)
		database.SetKey(key + ":" + conf.Occurrence + ":" + conf.Description ,value.Occurrence.Description)

		database.IncrValue(keyNoIndex + ":" + conf.UserOccurrence)

		for _,i := range value.Occurrence.SolveOccurrences {
			for z,j := range i.Activities {
				database.HSetKey(key + ":" + conf.Occurrence + ":" + conf.Solve + ":" + conf.Activity,strconv.Itoa(z),j.Id)
			}
		}

		for _,i := range value.Occurrence.SolveOccurrences {
			for z,j := range i.Activities {
				database.HSetKey(key + ":" + conf.Occurrence + ":" + conf.Solve + ":" + conf.Resource,strconv.Itoa(z),j.Id)
			}
		}
	}

	for k, value := range gm.ManagerOccurrences {
		key := conf.GetGameOccurrenceKeyType(gm.Id,conf.ManagerOccurrence,strconv.Itoa(k))
		hashString := strconv.Itoa(rand.Int()) + gm.Id + conf.ManagerOccurrence
		hashId := GerenateHashByString(hashString)
		database.SetKey(key + ":" + conf.Identifier,hashId)


		hashString2 := strconv.Itoa(rand.Int()) + gm.Id + conf.Occurrence
		hashId2 := GerenateHashByString(hashString2)
		database.SetKey(key + ":" + conf.Occurrence + ":" + conf.Identifier ,hashId2)
		database.SetKey(key + ":" + conf.Occurrence + ":" + conf.Description ,value.Occurrence.Description)

		database.IncrValue(keyNoIndex + ":" + conf.ManagerOccurrence)

		for _,i := range value.Occurrence.SolveOccurrences {
			for z,j := range i.Activities {
				database.HSetKey(key + ":" + conf.Occurrence + ":" + conf.Solve + ":" + conf.Activity,strconv.Itoa(z),j.Id)
			}
		}

		for _,i := range value.Occurrence.SolveOccurrences {
			for z,j := range i.Activities {
				database.HSetKey(key + ":" + conf.Occurrence + ":" + conf.Solve + ":" + conf.Resource,strconv.Itoa(z),j.Id)
			}
		}
	}
}

func GetOcurrences(gm GameModel) ([]UserOccurrence,[]ManagerOccurrence){
	keyNoIndex := conf.GetGameOccurrenceKey(gm.Id,conf.Count)
	countUser,_ := strconv.Atoi(database.GetKey(keyNoIndex + ":" + conf.UserOccurrence))
	countManager,_ := strconv.Atoi(database.GetKey(keyNoIndex + ":" + conf.ManagerOccurrence))
	var user []UserOccurrence
	var manager []ManagerOccurrence

	for k := 0; k <= countUser; k++ {
		key := conf.GetGameOccurrenceKeyType(gm.Id,conf.UserOccurrence,strconv.Itoa(k))
		var uo UserOccurrence
		uo.Id = database.GetKey(key + ":" + conf.Identifier)
		uo.Level, _ = strconv.Atoi(database.GetKey(key + ":" + conf.Level))


		uo.Occurrence.Id = database.GetKey(key + ":" + conf.Occurrence + ":" + conf.Identifier)
		uo.Occurrence.Description = database.GetKey(key + ":" + conf.Occurrence + ":" + conf.Description)

		ac := database.HGetAllKey(key + ":" + conf.Occurrence + ":" + conf.Solve + ":" + conf.Activity)
		rc := database.HGetAllKey(key + ":" + conf.Occurrence + ":" + conf.Solve + ":" + conf.Resource)
		numberOfActivities := len(ac)
		numberOfResources := len(rc)
		var solves []SolveOccurrence
		var as []Activity
		var rs []Resource

		for i := 0; i <= numberOfActivities; i++ {
			var a Activity
			a.Id = database.HGetKey(key + ":" + conf.Occurrence + ":" + conf.Solve + ":" + conf.Activity,strconv.Itoa(i))
			as = append(as,a)
		}

		for i := 0; i <= numberOfResources; i++ {
			var r Resource
			r.Id = database.HGetKey(key + ":" + conf.Occurrence + ":" + conf.Solve + ":" + conf.Resource,strconv.Itoa(i))
			rs = append(rs,r)
		}

		var solve SolveOccurrence
		solve.Activities = as
		solve.Resources = rs
		solves = append(solves,solve)

		uo.Occurrence.SolveOccurrences = solves

		user = append(user,uo)
	}

	for k := 0; k <= countManager; k++ {
		key := conf.GetGameOccurrenceKeyType(gm.Id,conf.ManagerOccurrence,strconv.Itoa(k))
		var mo ManagerOccurrence

		mo.Id = database.GetKey(key + ":" + conf.Identifier)

		mo.Occurrence.Id = database.GetKey(key + ":" + conf.Occurrence + ":" + conf.Identifier)
		mo.Occurrence.Description = database.GetKey(key + ":" + conf.Occurrence + ":" + conf.Description)

		ac := database.HGetAllKey(key + ":" + conf.Occurrence + ":" + conf.Solve + ":" + conf.Activity)
		rc := database.HGetAllKey(key + ":" + conf.Occurrence + ":" + conf.Solve + ":" + conf.Resource)
		numberOfActivities := len(ac)
		numberOfResources := len(rc)
		var solves []SolveOccurrence
		var as []Activity
		var rs []Resource


		for i := 0; i <= numberOfActivities; i++ {
			var a Activity
			a.Id = database.HGetKey(key + ":" + conf.Occurrence + ":" + conf.Solve + ":" + conf.Activity,strconv.Itoa(i))
			as = append(as,a)
		}

		for i := 0; i <= numberOfResources; i++ {
			var r Resource
			r.Id = database.HGetKey(key + ":" + conf.Occurrence + ":" + conf.Solve + ":" + conf.Resource,strconv.Itoa(i))
			rs = append(rs,r)
		}

		var solve SolveOccurrence
		solve.Activities = as
		solve.Resources = rs
		solves = append(solves,solve)

		mo.Occurrence.SolveOccurrences = solves

		manager = append(manager,mo)
	}

	return user,manager
}

func GetLevels(id string) []Level {
	var list []Level

	for i := 0; i < 7 ; i++ {
		keyId :=  conf.Identifier
		keyLevel := conf.Level
		var lv Level

		lv.Id = database.GetKey(conf.GetLevelKey(id,keyId,i))
		lv.Level = database.GetKey(conf.GetLevelKey(id,keyLevel,i))
		lv.Process = GetProcess(id,i)
		list = append(list,lv)
	}

	return list
}

func GetProcess(id string, level int) []Process {
	var list []Process
	keyCount := conf.GetLevelKey(id, conf.Process,level)
	max, _ := strconv.Atoi(database.GetKey(keyCount))

	for i := 0; i < max ; i++ {
		keyId := conf.Identifier
		keyName := conf.Name
		keyScore := conf.Score

		var pc Process
		pc.Name = database.GetKey(conf.GetProcessKey(id,keyName,level,i))
		pc.Id = database.GetKey(conf.GetProcessKey(id,keyId,level,i))
		sc,_ := strconv.Atoi(database.GetKey(conf.GetProcessKey(id,keyScore,level,i)))
		pc.Score = sc
		pc.Activities = GetActivities(id,level,i)
		pc.Resources = GetResources(id,level,i)
		list = append(list,pc)

	}

	return list
}

func GetActivities(id string, level,process int) []Activity {
	var list []Activity
	baseKey := conf.GetProcessKey(id, conf.Activity, level,process)
	max, _ := strconv.Atoi(database.GetKey(baseKey + ":" + conf.Count))

	for i := 0; i < max ; i++ {
		keyId := conf.Identifier
		keyName := conf.Name
		keyScore := conf.Score
		keyQuantity := conf.Quantity
		keyTime := conf.Time

		hashIndex := database.HGetKey(baseKey + ":" + conf.ActivityIds,strconv.Itoa(i))

		var ac Activity
		ac.Name = database.GetKey(conf.GetActivityKey(id,keyName,level,process,hashIndex))
		ac.Id = database.GetKey(conf.GetActivityKey(id,keyId,level,process,hashIndex))
		sc,_ := strconv.Atoi(database.GetKey(conf.GetActivityKey(id,keyScore,level,process,hashIndex)))
		qt,_ := strconv.Atoi(database.GetKey(conf.GetActivityKey(id,keyQuantity,level,process,hashIndex)))
		tm,_ := strconv.Atoi(database.GetKey(conf.GetActivityKey(id,keyTime,level,process,hashIndex)))
		ac.Score = sc
		ac.Quantity = qt
		ac.Time = tm
		list = append(list,ac)
	}

	return list
}

func GetResources(id string, level,process int) []Resource {
	var list []Resource
	baseKey := conf.GetProcessKey(id, conf.Resource, level,process)
	max, _ := strconv.Atoi(database.GetKey(baseKey + ":" + conf.Count))

	for i := 0; i < max ; i++ {
		keyId :=  conf.Identifier
		keyName := conf.Name
		keyScore :=  conf.Score
		keyQuantity :=  conf.Quantity
		keyPrice :=  conf.Price
		keyType :=  conf.Type

		hashIndex := database.HGetKey(baseKey + ":" + conf.ResourcesIds,strconv.Itoa(i))

		var rc Resource
		rc.Name = database.GetKey(conf.GetResourceKey(id,keyName,level,process,hashIndex))
		rc.Id = database.GetKey(conf.GetResourceKey(id,keyId,level,process,hashIndex))
		sc,_ := strconv.Atoi(database.GetKey(conf.GetResourceKey(id,keyScore,level,process,hashIndex)))
		qt,_ := strconv.Atoi(database.GetKey(conf.GetResourceKey(id,keyQuantity,level,process,hashIndex)))
		pr,_ := strconv.Atoi(database.GetKey(conf.GetResourceKey(id,keyPrice,level,process,hashIndex)))
		rc.Score = sc
		rc.Quantity = qt
		rc.Price = pr
		rc.Type = database.GetKey(conf.GetResourceKey(id,keyType,level,process,hashIndex))

		list = append(list,rc)
	}

	return list
}

func CreateLevel(gm GameModel, lv Level, level int) {
	database.SetKey(conf.GetLevelKey(gm.Id, conf.Identifier,level), lv.Id)
	database.SetKey(conf.GetLevelKey(gm.Id, conf.Level,level), lv.Level)
	//TODO: save all struct array

	for key, value := range lv.Process {
		k := strconv.Itoa(key)
		database.SetKey(conf.GetProcessKey(gm.Id, k,level,key), value.Id)
		CreateProcess(gm, value, level,key)
		database.IncrValue(conf.GetLevelKey(gm.Id, conf.Process,level))
	}
}

func CreateProcess(gm GameModel, pc Process, level,process int) {
	database.SetKey(conf.GetProcessKey(gm.Id, conf.Identifier,level,process), pc.Id)
	database.SetKey(conf.GetProcessKey(gm.Id, conf.Name,level,process), pc.Name)
	database.SetKey(conf.GetProcessKey(gm.Id, conf.Score,level,process), pc.Score)
	//TODO: save all struct array

	for key, value := range pc.Activities {
		hashId := strconv.Itoa(level) + gm.Id + conf.Activity
		value.Id = GerenateHashByString(hashId)
		//k := conf.Activity + ":" + value.Id
		//baseKey := conf.GetProcessKey(gm.Id, k, level,process)
		aKey := conf.GetProcessKey(gm.Id, conf.Activity, level,process)
		database.IncrValue(aKey + ":" + conf.Count)
		database.HSetKey(aKey + ":" + conf.ActivityIds, strconv.Itoa(key), value.Id)
		CreateActivity(gm, value, level, process,value.Id)
	}

	for key, value := range pc.Resources {
		hashId := strconv.Itoa(level) + gm.Id + conf.Resource
		value.Id = GerenateHashByString(hashId)
		//k := conf.Resource + ":" + value.Id
		//baseKey := conf.GetProcessKey(gm.Id, k, level,process)
		rKey := conf.GetProcessKey(gm.Id, conf.Resource, level,process)
		database.IncrValue(rKey + ":" + conf.Count)
		database.HSetKey(rKey + ":" + conf.ResourcesIds, strconv.Itoa(key),value.Id)
		CreateResource(gm, value, level, process,value.Id)
	}

	CreateScore(gm, pc.Score)
}

func CreateScore(gm GameModel, sc int) {
	//database.SetKey(conf.GetScoreKey(gm.Id, conf.Identifier), "sc-id")
	database.SetKey(conf.GetScoreKey(gm.Id, conf.Score), sc)
	//TODO: save all struct array
}

func CreateActivity(gm GameModel, at Activity, level,process int , unit string) {
	database.SetKey(conf.GetActivityKey(gm.Id, conf.Identifier,level,process,unit), at.Id)
	database.SetKey(conf.GetActivityKey(gm.Id, conf.Name,level,process,unit), at.Name)
	database.SetKey(conf.GetActivityKey(gm.Id, conf.Time,level,process,unit), at.Time)
	database.SetKey(conf.GetActivityKey(gm.Id, conf.Quantity,level,process,unit), at.Quantity)
	database.SetKey(conf.GetActivityKey(gm.Id, conf.Score,level,process,unit), at.Score)
	//TODO: save all struct array
	CreateScore(gm, at.Score)
}

func CreateResource(gm GameModel, rc Resource, level,process int , unit string) {
	database.SetKey(conf.GetResourceKey(gm.Id, conf.Identifier,level,process,unit), rc.Id)
	database.SetKey(conf.GetResourceKey(gm.Id, conf.Name,level,process,unit), rc.Name)
	database.SetKey(conf.GetResourceKey(gm.Id, conf.Price,level,process,unit), rc.Price)
	database.SetKey(conf.GetResourceKey(gm.Id, conf.Type,level,process,unit), rc.Type)
	database.SetKey(conf.GetResourceKey(gm.Id, conf.Quantity,level,process,unit), rc.Quantity)
	database.SetKey(conf.GetResourceKey(gm.Id, conf.Score,level,process,unit), rc.Score)
	//TODO: save all struct array
	CreateScore(gm, rc.Score)
}
