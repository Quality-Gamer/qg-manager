package endpoint

import (
	"encoding/hex"
	"encoding/json"
	"github.com/labstack/echo"
	"manager/conf"
	"manager/model"
	"manager/database"
	"crypto/sha1"
	"net/http"
	"strconv"
	"time"
)

func Create(c echo.Context) error {
	var res model.Response
	c.Response().Header().Set("Access-Control-Allow-Origin","*")
	c.Response().Header().Set(echo.HeaderContentType,echo.MIMEApplicationJSONCharsetUTF8)

	if len(c.FormValue("user_id")) > 0 && len(c.FormValue("challenge_id")) > 0 {
		userId, _ := strconv.Atoi(c.FormValue("user_id"))
		challengeId, _ :=  strconv.Atoi(c.FormValue("challenge_id"))
		projectId := 1
		match := createManagerMatch(userId,challengeId,projectId)
		res.Message = conf.SuccessMessage
		res.Status = conf.SuccessCode
		res.Response = append(res.Response, match)
		c.Response().WriteHeader(http.StatusOK)
		return json.NewEncoder(c.Response()).Encode(res)
	} else {
		res.Status = conf.ErrorCode
		res.Message = conf.ErrorInputMessage
		c.Response().WriteHeader(http.StatusBadRequest)
		return json.NewEncoder(c.Response()).Encode(res)
	}
}

func createManagerMatch(userId,challengeId,projectId int) model.ManagerMatch{
	var new model.ManagerMatch
	new.Id = generateId(userId,projectId)
	new.UserId = userId
	new.ProjectId = projectId
	new.Progress = 0.0
	new.Week = 1
	new.Money = conf.WeekMoney
	new.Time = conf.WeekTime
	new.ChallengeId = challengeId
	new.ProgressStatus = "N"

	keyMatch := conf.GetKeyManager(new.UserId,new.Week,new.Id)
	keyProgress := keyMatch + ":" + conf.Progress
	keyNoWeek := conf.GetKeyOccurrence(new.UserId,new.Id)
	keyProjectId := keyNoWeek + ":" + conf.ProjectId
	keyChallengeId := keyNoWeek + ":" + conf.ChallengeId
	keyOccurrence := keyNoWeek + ":" + conf.Occurrence
	keyCurrentWeek :=  keyNoWeek + ":" + conf.CurrentWeek
	keyCurrentMoney :=  keyNoWeek + ":" + conf.CurrentMoney
	keyCurrentTime :=  keyNoWeek + ":" + conf.CurrentTime

	database.SetKey(keyProjectId,new.ProjectId)
	database.SetKey(keyChallengeId,new.ChallengeId)
	database.SetKey(keyProgress,new.Progress)
	database.SetKey(keyOccurrence,0)
	database.SetKey(keyCurrentWeek,new.Week)
	database.SetKey(keyCurrentMoney,new.Money)
	database.SetKey(keyCurrentTime,new.Time)

	return new
}

func generateId(userId,projectId int) string{
	string := time.Now().String()+ string(userId) + string(projectId)
	hash := sha1.New()
	hash.Write([]byte(string))
	return hex.EncodeToString(hash.Sum(nil))
}