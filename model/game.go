package model

import (
	"crypto/sha1"
	"encoding/hex"
	"math"
	"qg-manager/conf"
	"qg-manager/database"
	"strconv"
	"time"
)

const TeamType = "T"
const ProductType = "P"

type ManagerMatch struct {
	Id             string              `json:"id"`
	ChallengeId    int                 `json:"-"`
	GameModel      GameModel           `json:"-"`
	UserId         int                 `json:"user_id"`
	ProjectId      int                 `json:"-"`
	Week           int                 `json:"week"`
	Progress       float64             `json:"-"`
	ProgressStatus string              `json:"-"`
	Level          int                 `json:"level"`
	Money          float64             `json:"money"`
	Time           int                 `json:"time"`
	Team           Team                `json:"-"`
	Resources      MatchResource       `json:"resources"`
	License        License             `json:"-"`
	Action         Action              `json:"-"`
	Activities     []Activity     `json:"activities"`
	Event          Event               `json:"event"`
	Occurrence     []ManagerOccurrence `json:"manager_occurrence"`
	UserOccurrence []UserOccurrence    `json:"user_occurrence"`
}

type MatchResource struct {
	Team     Team      `json:"team"`
	Products []Product `json:"products"`
}

type Team struct {
	Members []Member `json:"members"`
}

type Product struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
}

type Member struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
}

type License struct {
	Ide            int `json:"ide"`
	DesignSoftware int `json:"design_software"`
}

type Action struct {
	Scrum           int `json:"scrum"`
	Delivery        int `json:"delivery"`
	CustomerContact int `json:"customer_contact"`
	RiskAnalysis    int `json:"risk_analysis"`
}

type MatchActivity struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
}

type Event struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

func NewMatch(modelId string,userId int) *ManagerMatch {
	game := GetGameModel(modelId)
	match := new(ManagerMatch)
	match.Id = GerenateHash(userId)
	match.UserId = userId
	match.GameModel = game
	match.Week = 1
	match.Level = 0
	match.Time = match.GameModel.InitialTime
	match.Money = match.GameModel.InitialMoney
	createMatch(match)
	return match
}

func FindManagerMatch(userId int, matchId string) (ManagerMatch,bool) {
	var match ManagerMatch
	//keyMatch := conf.GetKeyManager(userId,week,matchId)
	//keyLevel := keyMatch + ":" + conf.Level
	keyNoWeek := conf.GetKeyOccurrence(userId,matchId)
	keyGameModelId := keyNoWeek + ":" + conf.Model
	//keyOccurrence := keyNoWeek + ":" + conf.Occurrence
	keyCurrentWeek :=  keyNoWeek + ":" + conf.CurrentWeek
	keyCurrentMoney :=  keyNoWeek + ":" + conf.CurrentMoney
	keyCurrentTime :=  keyNoWeek + ":" + conf.CurrentTime
	keyCurrentLevel :=  keyNoWeek + ":" + conf.Level

	match.Id = matchId
	var exists bool
	exists = false

	//GameModel is not returned in the response
	match.GameModel = GetGameModel(database.GetKey(keyGameModelId))
	if len(match.GameModel.Id) > 0 { exists = true }

	match.UserId = userId
	match.Week,_ = strconv.Atoi(database.GetKey(keyCurrentWeek))
	match.Money,_ = strconv.ParseFloat(database.GetKey(keyCurrentMoney),64)
	match.Time,_ = strconv.Atoi(database.GetKey(keyCurrentTime))
	match.Level,_ = strconv.Atoi(database.GetKey(keyCurrentLevel))
	match.Resources.Team, match.Resources.Products = GetMatchResources(userId,match.Level,matchId)
	match.Activities = GetMatchActivities(userId,match.Level,matchId)

	//Event          Event               `json:"event"`
	//Occurrence     []ManagerOccurrence `json:"manager_occurrence"`
	//UserOccurrence []UserOccurrence    `json:"user_occurrence"`

	return match,exists
}

func GetMatchActivities(userId, level int, matchId string) []Activity {
	keyNoWeek := conf.GetKeyOccurrence(userId,matchId)
	key := keyNoWeek + ":" + conf.Team + ":" + conf.Member
	modelId := GetModelId(userId,matchId)
	keyCount := conf.GetLevelKey(modelId, conf.Level,level) + ":" + conf.Process
	nProcess, _ := strconv.Atoi(database.GetKey(keyCount))
	var activities []Activity

	for i := 0; i < nProcess; i++ {
		a := GetActivities(modelId,level,i)
		activities = append(activities,a...)
	}

	var act []Activity

	for _, activity := range activities {
		count, _ := strconv.Atoi(database.HGetKey(key,activity.Id))
		if count > 0 {
			var m Activity
			m.Id = activity.Id
			m.Name = activity.Name
			m.Quantity = count
			act = append(act,m)
		}
	}

	return act
}

func GetMatchResources(userId, level int, matchId string) (Team,[]Product) {
	keyNoWeek := conf.GetKeyOccurrence(userId,matchId)
	key := keyNoWeek + ":" + conf.Team + ":" + conf.Member
	modelId := GetModelId(userId,matchId)
	keyCount := conf.GetLevelKey(modelId, conf.Level,level) + ":" + conf.Process
	nProcess, _ := strconv.Atoi(database.GetKey(keyCount))
	var resources []Resource

	for i := 0; i < nProcess; i++ {
		r := GetResources(modelId,level,i)
		resources = append(resources,r...)
	}

	var team Team
	var products []Product

	for _, resource := range resources {
		if resource.Type == TeamType {
			count, _ := strconv.Atoi(database.HGetKey(key,resource.Id))
			if count > 0 {
				var m Member
				m.Id = resource.Id
				m.Name = resource.Name
				m.Quantity = count
				team.Members = append(team.Members,m)
			}
		} else if resource.Type == ProductType {
			count, _ := strconv.Atoi(database.HGetKey(key,resource.Id))
			if count > 0 {
				var p Product
				p.Id = resource.Id
				p.Name = resource.Name
				p.Quantity = count
				products = append(products,p)
			}
		}
	}

	return team,products
}

func GetModel(modelId string) GameModel {
	game := GetGameModel(modelId)
	return game
}

func GetModelId(user_id int, matchId string) string {
	key := conf.GetKeyOccurrence(user_id,matchId) + ":" + conf.Model
	return database.GetKey(key)
}

func GetCurrentLevel(user_id, week int, matchId string) int {
	key := conf.GetKeyManager(user_id,week,matchId) + ":" + conf.Level
	lv,_ := strconv.Atoi(database.GetKey(key))
	return lv
}

func GetCurrentTime(user_id int, matchId string) int {
	key := conf.GetKeyOccurrence(user_id,matchId) + ":" + conf.CurrentTime
	tm,_ := strconv.Atoi(database.GetKey(key))
	return tm
}

func GetCurrentMoney(user_id int, matchId string) float64 {
	key := conf.GetKeyOccurrence(user_id,matchId) + ":" + conf.CurrentMoney
	mn,_ := strconv.ParseFloat(database.GetKey(key),64)
	return mn
}


func createMatch(m *ManagerMatch) {
	keyMatch := conf.GetKeyManager(m.UserId,m.Week,m.Id)
	keyLevel := keyMatch + ":" + conf.Level
	keyNoWeek := conf.GetKeyOccurrence(m.UserId,m.Id)
	keyGameModelId := keyNoWeek + ":" + conf.Model
	keyOccurrence := keyNoWeek + ":" + conf.Occurrence
	keyCurrentWeek :=  keyNoWeek + ":" + conf.CurrentWeek
	keyCurrentMoney :=  keyNoWeek + ":" + conf.CurrentMoney
	keyCurrentTime :=  keyNoWeek + ":" + conf.CurrentTime
	keyCurrentLevel :=  keyNoWeek + ":" + conf.Level

	database.SetKey(keyGameModelId,m.GameModel.Id)
	database.SetKey(keyOccurrence,0)
	database.SetKey(keyCurrentWeek,m.Week)
	database.SetKey(keyCurrentMoney,m.Money)
	database.SetKey(keyCurrentTime,m.Time)
	database.SetKey(keyLevel,m.Level)
	database.SetKey(keyCurrentLevel,m.Level)
}

func GerenateHash(userId int) string{
	string := time.Now().String()+ string(userId)
	hash := sha1.New()
	hash.Write([]byte(string))
	return hex.EncodeToString(hash.Sum(nil))
}

func GerenateHashByString(hashId string) string{
	string := time.Now().String()+ hashId
	hash := sha1.New()
	hash.Write([]byte(string))
	return hex.EncodeToString(hash.Sum(nil))
}

//Deprecated: There is another better function to do it
func getFirstLayerAttr(modelId string) (string,string,string,int,float64,string,int) {
	id := database.GetKey(conf.GetGameModelKey(modelId, conf.Identifier))
	name := database.GetKey(conf.GetGameModelKey(modelId, conf.Name))
	desc := database.GetKey(conf.GetGameModelKey(modelId, conf.Description))
	it := database.GetKey(conf.GetGameModelKey(modelId, conf.InitialTime))
	im := database.GetKey(conf.GetGameModelKey(modelId, conf.InitialMoney))
	tp := database.GetKey(conf.GetGameModelKey(modelId, conf.Type))
	bt := database.GetKey(conf.GetGameModelKey(modelId, conf.BelongsTo))
	int_it, _ := strconv.Atoi(it)
	float_im, _ := strconv.ParseFloat(im,64)
	int_bt, _ := strconv.Atoi(bt)
	return id,name,desc,int_it,float_im,tp,int_bt
}

func AddResource(userId int, matchId, itemId string) {
	keyNoWeek := conf.GetKeyOccurrence(userId,matchId)
	keyTeam := keyNoWeek + ":" + conf.Team
	keyMember := keyTeam + ":" + conf.Member
	database.HSetIncrKey(keyMember,itemId)
}


func AddActivity(userId int, matchId, itemId string) {
	keyNoWeek := conf.GetKeyOccurrence(userId,matchId)
	keyAc := keyNoWeek + ":" + conf.Activity
	database.HSetIncrKey(keyAc,itemId)
}

func (mm *ManagerMatch) RunGame() bool {
	level := mm.Level
	levelModel := mm.GameModel.Levels[level]
	levelMax := len(mm.GameModel.Levels)
	level = int(math.Min(float64(level), float64(levelMax)))

	solveProcess := 0
	nProcess := len(levelModel.Process)

	for _, value := range levelModel.Process {
		scoreResourcesTeam := 0
		scoreResourcesProduct := 0
		scoreActivities := 0

		for _, j := range value.Resources {
			for _, v := range mm.Resources.Team.Members {
				if v.Id == j.Id && v.Quantity >= j.Quantity {
					scoreResourcesTeam += j.Score
				}
			}

			for _, v := range mm.Resources.Products {
				if v.Id == j.Id && v.Quantity >= j.Quantity {
					scoreResourcesProduct += j.Score
				}
			}
		}

		for _, j := range value.Activities {
			for _, v := range mm.Activities {
				if v.Id == j.Id && v.Quantity >= j.Quantity {
					scoreActivities += j.Score
				}
			}
		}

		if scoreActivities+scoreResourcesProduct+scoreResourcesTeam >= value.Score {
			solveProcess += 1
		}
	}

	if solveProcess == nProcess {
		mm.Level = mm.Level + 1
		return true
	}

	return false
}
