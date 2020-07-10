package endpoint

import (
	"encoding/json"
	"github.com/labstack/echo"
	"github.com/Quality-Gamer/qg-manager/conf"
	"github.com/Quality-Gamer/qg-manager/database"
	"github.com/Quality-Gamer/qg-manager/model"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func Next(c echo.Context) error {
	var res model.Response
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)

	if len(c.FormValue("user_id")) > 0 && len(c.FormValue("manager_id")) > 0 {
		userId, _ := strconv.Atoi(c.FormValue("user_id"))
		managerId := c.FormValue("manager_id")
		week, _ := strconv.Atoi(database.GetKey(conf.GetKeyOccurrence(userId,managerId) + ":" + conf.CurrentWeek))
		week += 1
		match,end := generateNextWeek(userId,week,managerId)
		res.Status = conf.SuccessCode
		res.Message = conf.SuccessMessage
		res.Response = append(res.Response,match)

		endGame := make(map[string]int)

		if end {
			endGame["end"] = 1
		} else {
			endGame["end"] = 0
		}

		res.Response = append(res.Response, endGame)

		c.Response().WriteHeader(http.StatusOK)
		return json.NewEncoder(c.Response()).Encode(res)
	} else {
		res.Status = conf.ErrorCode
		res.Message = conf.ErrorInputMessage
		c.Response().WriteHeader(http.StatusBadRequest)
		return json.NewEncoder(c.Response()).Encode(res)
	}
}

func generateNextWeek(userId,week int, managerId string) (model.ManagerMatch,bool) {
	var new model.ManagerMatch
	old := loadOldWeek(userId,week,managerId)

	if week > 8 {
		return old,true
	}

	new = old
	new.Week = week
	new.Money += conf.WeekMoney
	new.Time += conf.WeekTime

	keyMatch := conf.GetKeyManager(new.UserId,new.Week,new.Id)
	keyProjectId := keyMatch + ":" + conf.ProjectId
	keyProgress := keyMatch + ":" + conf.Progress

	database.SetKey(keyProjectId,new.ProjectId)
	database.SetKey(keyProgress,new.Progress)

	level := 7

	if hasOccurrence(level) {
		var occurrence model.ManagerOccurrence
		occurrence.Occurrence = getOccurrence()
		occurrence.Status = "O"
		new.Occurrence = append(new.Occurrence, occurrence)
		keyOcurrences := conf.GetKeyOccurrence(new.UserId,new.Id)
		numberOfOccurrences,_ := strconv.Atoi(database.GetKey(keyOcurrences + ":" + conf.Occurrence + ":" + conf.NumberOccurrences))
		newNumberOfOccurrences := numberOfOccurrences + 1

		database.SetKey(keyOcurrences + ":" + conf.Occurrence + ":" + conf.NumberOccurrences, newNumberOfOccurrences)
		database.SetKey(keyOcurrences + ":" + conf.Occurrence + ":" + strconv.Itoa(newNumberOfOccurrences) + ":" + conf.OccurrenceId,occurrence.Occurrence.Id)
		database.SetKey(keyOcurrences + ":" + conf.Occurrence + ":" + strconv.Itoa(newNumberOfOccurrences) + ":" + conf.Description,occurrence.Occurrence.Description)
		database.SetKey(keyOcurrences + ":" + conf.Occurrence + ":" + strconv.Itoa(newNumberOfOccurrences) + ":" + conf.Status,occurrence.Status)
	}

	new.UserOccurrence = hasUserOccurrence(new)

	for _,value := range new.UserOccurrence {
		database.SetKey(keyMatch + ":" + conf.UserOccurrence + ":" +  strconv.Itoa(value.Occurrence.Id) + ":" + conf.Description,value.Occurrence.Description)
		database.SetKey(keyMatch + ":" + conf.UserOccurrence + ":" +  strconv.Itoa(value.Occurrence.Id) + ":" + conf.Status,value.Status)
	}

	progress, progressStatus := model.GameLogic(new)

	new.Progress = progress
	new.ProgressStatus = progressStatus

	nextWeekSaveManagerMatch(new)

	return new,false
}

func loadOldWeek(userId,week int, managerId string) model.ManagerMatch {
	oldWeek := week - 1
	manager,_ := findManagerMatch(userId, oldWeek, managerId)
	return manager
}

func getOccurrence() model.Occurrence{
	list := model.LoadOccurrenceList()
	len := len(list)
	id := random(0, len - 1)
	return list[id]
}

func hasOccurrence(level int) bool{
	if random(0,100) > (100 - getProbabilty(level) ) {
		return true
	} else {
		return false
	}
}

func random(min,max int) int{
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max - min + 1) + min
}

func getProbabilty(level int) int{
	// 7 levels
	//each level add 10%
	return conf.OccurrenceProbability * level
}

func hasUserOccurrence(match model.ManagerMatch) []model.UserOccurrence{
	list := model.LoadUserOccurrenceList()
	var listUserOccurrence []model.UserOccurrence
	programmers := match.Team.Backend + match.Team.Frontend
	designers := match.Team.Designer

	if programmers > match.License.Ide {
		var uo model.UserOccurrence
		uo.Status = "O"
		uo.Occurrence = list[0]
		new := strconv.Itoa(programmers - match.License.Ide)
		uo.Occurrence.Description = strings.Replace(uo.Occurrence.Description,"x", new,1)
		listUserOccurrence = append(listUserOccurrence,uo)
	}

	if designers > match.License.DesignSoftware {
		var uo model.UserOccurrence
		uo.Status = "O"
		uo.Occurrence = list[1]
		new := strconv.Itoa(designers - match.License.DesignSoftware)
		uo.Occurrence.Description = strings.Replace(uo.Occurrence.Description,"x", new,1)
		listUserOccurrence = append(listUserOccurrence,uo)
	}

	return listUserOccurrence
}

func nextWeekSaveManagerMatch(new model.ManagerMatch) {
	keyMatch := conf.GetKeyManager(new.UserId,new.Week,new.Id)
	keyProjectId := keyMatch + ":" + conf.ProjectId
	keyProgress := keyMatch + ":" + conf.Progress

	database.SetKey(keyProjectId,new.ProjectId)
	database.SetKey(keyProgress,new.Progress)
	database.SetKey(keyMatch + ":" + conf.Team + ":" + conf.Tester,new.Team.Tester)
	database.SetKey(keyMatch + ":" + conf.Team + ":" + conf.RequirementAnalyst,new.Team.RiskAnalyst)
	database.SetKey(keyMatch + ":" + conf.Team + ":" + conf.ProductOwner,new.Team.ProductOwner)
	database.SetKey(keyMatch + ":" + conf.Team + ":" + conf.Backend,new.Team.Backend)
	database.SetKey(keyMatch + ":" + conf.Team + ":" + conf.Frontend,new.Team.Frontend)
	database.SetKey(keyMatch + ":" + conf.Team + ":" + conf.Designer,new.Team.Designer)
	database.SetKey(keyMatch + ":" + conf.License + ":" + conf.Ide,new.License.Ide)
	database.SetKey(keyMatch + ":" + conf.License + ":" + conf.DesignSoftware,new.License.DesignSoftware)
	database.SetKey(keyMatch + ":" + conf.Action + ":" + conf.Scrum,new.Action.Scrum)
	database.SetKey(keyMatch + ":" + conf.Action + ":" + conf.CustomerContact,new.Action.CustomerContact)
	database.SetKey(keyMatch + ":" + conf.Action + ":" + conf.Delivery,new.Action.Delivery)
	database.SetKey(keyMatch + ":" + conf.Action + ":" + conf.RiskAnalysis,new.Action.RiskAnalysis)
	database.SetKey(conf.GetKeyOccurrence(new.UserId,new.Id) + ":" + conf.CurrentWeek, new.Week)
	database.SetKey(conf.GetKeyOccurrence(new.UserId,new.Id) + ":" + conf.CurrentMoney, new.Money)
	database.SetKey(conf.GetKeyOccurrence(new.UserId,new.Id) + ":" + conf.CurrentTime, new.Time)
}
