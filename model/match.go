package model

import (
	"manager/conf"
	"manager/database"
	"math"
	"strconv"
	"strings"
)

type ManagerMatch struct {
	Id             string              `json:"id"`
	ChallengeId    int                 `json:"challenge_id"`
	UserId         int                 `json:"user_id"`
	ProjectId      int                 `json:"project_id"`
	Week           int                 `json:"weak"`
	Progress       float64             `json:"progress"`
	ProgressStatus string			   `json:"progress_status"`
	Money          int                 `json:"money"`
	Time           int                 `json:"time"`
	Team           Team                `json:"team"`
	License        License             `json:"license"`
	Action         Action              `json:"action"`
	Event          Event               `json:"event"`
	Occurrence     []ManagerOccurrence `json:"manager_occurrence"`
	UserOccurrence []UserOccurrence    `json:"user_occurrence"`
}

type Team struct {
	Backend      int `json:"backend"`
	Frontend     int `json:"frontend"`
	Designer     int `json:"designer"`
	Tester       int `json:"tester"`
	ProductOwner int `json:"product_owner"`
	RiskAnalyst  int `json:"requirement_analyst"`
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

type Event struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

func GameLogic(manager ManagerMatch) (float64,string){
	initialPerc := balancedTeam(manager)
	initialPerc += hasOccurrence(manager)
	initialPerc += hasUserOccurrence(manager)
	thisWeekPerc := initialPerc * conf.TotalWeek

	progress := saveWeekPerc(manager,thisWeekPerc)
	status := getProgressStatus(manager,progress)

	return progress,status
}

func getProgressStatus(manager ManagerMatch,progress float64) string{
	min := float64(manager.Week) * conf.MinGoalWeek
	max := float64(manager.Week) * conf.MaxGoalWeek

	if progress < min {
		return "L" //late
	}

	if progress > max {
		return "A" //ahead
	}

	return "N" //normal
}

func saveWeekPerc(manager ManagerMatch, progress float64) float64 {
	oldWeek := manager.Week - 1
	key := conf.GetKeyManager(manager.UserId,oldWeek,manager.Id)
	pg,_ := strconv.Atoi(database.GetKey(key + ":" + conf.Progress))
	floatPG := float64(pg)
	newPG := floatPG + progress
	database.SetKey(key + ":" + conf.Progress,newPG)
	return newPG
}

func balancedTeam(manager ManagerMatch) float64 {
	var devTeam, devTeamAverage int
	var perc float64
	perc = 0.0

	devTeam = manager.Team.Frontend + manager.Team.Backend + manager.Team.Designer
	devTeamAverage = devTeam / 3

	if manager.Team.ProductOwner == 1 {
		perc += conf.GeneralTeam
	} else {
		diffPO, _ := strconv.ParseFloat(strconv.Itoa(manager.Team.ProductOwner-1), 64)
		perc += diffPO * conf.GeneralLoss
	}

	if manager.Team.RiskAnalyst == devTeamAverage/conf.PropRA {
		perc += conf.GeneralTeam
	} else {
		diffRA, _ := strconv.ParseFloat(strconv.Itoa(manager.Team.RiskAnalyst-(devTeamAverage/conf.PropRA)), 64)
		perc += diffRA * conf.GeneralLoss
	}

	if manager.Team.Tester == devTeamAverage/conf.PropTT {
		perc += conf.GeneralTeam
	} else {
		diffRA, _ := strconv.ParseFloat(strconv.Itoa(manager.Team.Tester-(devTeamAverage/conf.PropTT)), 64)
		perc += diffRA * conf.TesterLoss
	}

	if manager.Team.Frontend == manager.Team.Backend && manager.Team.Backend == manager.Team.Designer {
		perc += conf.DevTeam
	} else {
		max := math.Max(float64(manager.Team.Frontend), float64(manager.Team.Backend))
		max = math.Max(max, float64(manager.Team.Designer))
		min := math.Min(float64(manager.Team.Frontend), float64(manager.Team.Backend))
		min = math.Min(min, float64(manager.Team.Designer))

		diff := math.Abs(float64(float64(devTeamAverage) - max)) + math.Abs(float64(float64(devTeamAverage) - min))
		perc += diff * conf.DevTeamLoss
	}

	return perc
}

func hasOccurrence(manager ManagerMatch) float64 {
	var perc float64
	perc = 0.0
	keyOccurrences := conf.GetKeyOccurrence(manager.UserId,manager.Id)
	numberOfOccurrences, _ := strconv.Atoi(database.GetKey(keyOccurrences + ":" + conf.Occurrence + ":" + conf.NumberOccurrences))

	if numberOfOccurrences == 0 {
		return conf.OccurrenceFalse
	}

	perc += float64(numberOfOccurrences) * conf.OccurrenceTrue

	return perc
}

func hasUserOccurrence(manager ManagerMatch) float64 {
	var perc float64

	userOccurrence := getUserOccurrences(manager)
	var index int
	index = 0

	for _, value := range userOccurrence {
		if value.Status != "C" {
			index += 1
		}
	}

	if index == 0{
		return conf.OccurrenceUserFalse
	}

	perc = float64(index) * conf.OccurrenceUserTrue

	return perc
}

func getUserOccurrences(match ManagerMatch) []UserOccurrence{
	list := LoadUserOccurrenceList()
	var listUserOccurrence []UserOccurrence
	programmers := match.Team.Backend + match.Team.Frontend
	designers := match.Team.Designer

	if programmers > match.License.Ide {
		var uo UserOccurrence
		uo.Status = "O"
		uo.Occurrence = list[0]
		new := strconv.Itoa(programmers - match.License.Ide)
		uo.Occurrence.Description = strings.Replace(uo.Occurrence.Description,"x", new,1)
		listUserOccurrence = append(listUserOccurrence,uo)
	}

	if designers > match.License.DesignSoftware {
		var uo UserOccurrence
		uo.Status = "O"
		uo.Occurrence = list[1]
		new := strconv.Itoa(designers - match.License.DesignSoftware)
		uo.Occurrence.Description = strings.Replace(uo.Occurrence.Description,"x", new,1)
		listUserOccurrence = append(listUserOccurrence,uo)
	}

	return listUserOccurrence
}