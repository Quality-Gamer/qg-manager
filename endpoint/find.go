package endpoint

import (
	"encoding/json"
	"github.com/labstack/echo"
	"manager/conf"
	"manager/database"
	"manager/model"
	"net/http"
	"strconv"
)

func Find(c echo.Context) error {
	var res model.Response
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)

	if len(c.FormValue("user_id")) > 0 && len(c.FormValue("manager_id")) > 0 {
		userId, _ := strconv.Atoi(c.FormValue("user_id"))
		managerId := c.FormValue("manager_id")
		var week int

		if len(c.FormValue("week")) > 0 {
			week, _ = strconv.Atoi(c.FormValue("week"))
		} else {
			week, _ = strconv.Atoi(database.GetKey(conf.GetKeyOccurrence(userId,managerId) + ":" + conf.CurrentWeek))
		}

		match, err := findManagerMatch(userId, week, managerId)
		res.Message = conf.SuccessMessage
		res.Status = conf.SuccessCode
		res.Response = append(res.Response, match)

		if err {
			res.Message = conf.ErrorDoesNotExist
			res.Status = conf.ErrorCode
			res.Response = nil
		}

		c.Response().WriteHeader(http.StatusOK)
		return json.NewEncoder(c.Response()).Encode(res)
	} else {
		res.Status = conf.ErrorCode
		res.Message = conf.ErrorInputMessage
		c.Response().WriteHeader(http.StatusBadRequest)
		return json.NewEncoder(c.Response()).Encode(res)
	}
}

func findManagerMatch(userId, week int, managerId string) (model.ManagerMatch, bool) {
	var manager model.ManagerMatch
	keyMatch := conf.GetKeyManager(userId,week,managerId)
	manager.Id = managerId
	manager.UserId = userId
	manager.ChallengeId, _ = strconv.Atoi(database.GetKey(conf.GetKeyOccurrence(userId,managerId) + ":" + conf.ChallengeId))
	manager.Week = week
	manager.ProjectId, _ = strconv.Atoi(database.GetKey(conf.GetKeyOccurrence(manager.UserId,manager.Id) + ":" + conf.ProjectId))
	manager.Progress, _ = strconv.ParseFloat(database.GetKey(keyMatch+":"+conf.Progress), 64)
	manager.Team.Tester, _ = strconv.Atoi(database.GetKey(keyMatch + ":" + conf.Team + ":" + conf.Tester))
	manager.Team.RiskAnalyst, _ = strconv.Atoi(database.GetKey(keyMatch + ":" + conf.Team + ":" + conf.RiskAnalyst))
	manager.Team.ProductOwner, _ = strconv.Atoi(database.GetKey(keyMatch + ":" + conf.Team + ":" + conf.ProductOwner))
	manager.Team.Designer, _ = strconv.Atoi(database.GetKey(keyMatch + ":" + conf.Team + ":" + conf.Designer))
	manager.Team.Frontend, _ = strconv.Atoi(database.GetKey(keyMatch + ":" + conf.Team + ":" + conf.Frontend))
	manager.Team.Backend, _ = strconv.Atoi(database.GetKey(keyMatch + ":" + conf.Team + ":" + conf.Backend))
	manager.Action.CustomerContact, _ = strconv.Atoi(database.GetKey(keyMatch + ":" + conf.Action + ":" + conf.CustomerContact))
	manager.Action.Delivery, _ = strconv.Atoi(database.GetKey(keyMatch + ":" + conf.Action + ":" + conf.Delivery))
	manager.Action.RiskAnalysis, _ = strconv.Atoi(database.GetKey(keyMatch + ":" + conf.Action + ":" + conf.RiskAnalysis))
	manager.Action.Scrum, _ = strconv.Atoi(database.GetKey(keyMatch + ":" + conf.Action + ":" + conf.Scrum))
	manager.License.Ide, _ = strconv.Atoi(database.GetKey(keyMatch + ":" + conf.License + ":" + conf.Ide))
	manager.License.DesignSoftware, _ = strconv.Atoi(database.GetKey(keyMatch + ":" + conf.License + ":" + conf.DesignSoftware))
	manager.Event.Status = database.GetKey(keyMatch + ":" + conf.Event + ":" + conf.Status)
	manager.Event.Name = database.GetKey(keyMatch + ":" + conf.Event + ":" + conf.Name)

	oldMoney, _ := strconv.Atoi(database.GetKey(conf.GetKeyOccurrence(manager.UserId,manager.Id) + ":" + conf.CurrentMoney))
	oldTime, _ := strconv.Atoi(database.GetKey(conf.GetKeyOccurrence(manager.UserId,manager.Id) + ":" + conf.CurrentTime))

	manager.Money = oldMoney + conf.WeekMoney
	manager.Time = oldTime + conf.WeekTime

	keyOccurrences := conf.GetKeyOccurrence(userId,managerId)
	numberOfOccurrences,_ := strconv.Atoi(database.GetKey(keyOccurrences + ":" + conf.Occurrence + ":" + conf.NumberOccurrences))

	for n := 1 ; n <= numberOfOccurrences ; n++  {
		var occurrence model.ManagerOccurrence
		occurrence.Occurrence.Id, _ = strconv.Atoi(database.GetKey(keyOccurrences + ":" + conf.Occurrence + ":" + strconv.Itoa(n) + ":" + conf.OccurrenceId))
		occurrence.Occurrence.Description = database.GetKey(keyOccurrences + ":" + conf.Occurrence + ":" + strconv.Itoa(n) + ":" + conf.Description)
		occurrence.Status = database.GetKey(keyOccurrences + ":" + conf.Occurrence + ":" + strconv.Itoa(n) + ":" + conf.Status)
		manager.Occurrence = append(manager.Occurrence, occurrence)
	}

	updateUserOccurrence(keyMatch,manager)

	userOccurrenceDescriptionProgrammers := database.GetKey(keyMatch + ":" + conf.UserOccurrence + ":" + "1" + ":" + conf.Description)
	userOccurrenceStatusProgrammers := database.GetKey(keyMatch + ":" + conf.UserOccurrence + ":" + "1" + ":" + conf.Status)
	userOccurrenceDescriptionDesigners := database.GetKey(keyMatch + ":" + conf.UserOccurrence + ":" + "2" + ":" + conf.Description)
	userOccurrenceStatusDesigners := database.GetKey(keyMatch + ":" + conf.UserOccurrence + ":" + "2" + ":" + conf.Status)

	if userOccurrenceDescriptionProgrammers != "" && userOccurrenceStatusProgrammers != "" {
		var uo model.UserOccurrence
		uo.Occurrence.Id = 1
		uo.Status = userOccurrenceStatusProgrammers
		uo.Occurrence.Description = userOccurrenceDescriptionProgrammers
		manager.UserOccurrence = append(manager.UserOccurrence,uo)
	}

	if userOccurrenceDescriptionDesigners != "" && userOccurrenceStatusDesigners != "" {
		var uo model.UserOccurrence
		uo.Occurrence.Id = 2
		uo.Status = userOccurrenceStatusDesigners
		uo.Occurrence.Description = userOccurrenceDescriptionDesigners
		manager.UserOccurrence = append(manager.UserOccurrence,uo)
	}

	currentWeek, _ := strconv.Atoi(database.GetKey(conf.GetKeyOccurrence(manager.UserId,manager.Id) + ":" + conf.CurrentWeek))

	if manager.ProjectId == 0 || manager.Week > currentWeek {
		return manager, true
	}

	return manager, false
}

func updateUserOccurrence(key string, match model.ManagerMatch) {
	uo := hasUserOccurrence(match)
	currentWeek, _ := strconv.Atoi(database.GetKey(conf.GetKeyOccurrence(match.UserId,match.Id) + ":" + conf.CurrentWeek))

	if (match.Team.Backend + match.Team.Frontend) == match.License.Ide && !(currentWeek < match.Week) {
		database.SetKey(key + ":" + conf.UserOccurrence + ":" +  "1" + ":" + conf.Status,"C")
	}

	if match.Team.Designer == match.License.DesignSoftware && !(currentWeek < match.Week) {
		database.SetKey(key + ":" + conf.UserOccurrence + ":" +  "2" + ":" + conf.Status,"C")
	}

	for _,value := range uo {
		database.SetKey(key + ":" + conf.UserOccurrence + ":" +  strconv.Itoa(value.Occurrence.Id) + ":" + conf.Description,value.Occurrence.Description)
		database.SetKey(key + ":" + conf.UserOccurrence + ":" +  strconv.Itoa(value.Occurrence.Id) + ":" + conf.Status,value.Status)
	}

}