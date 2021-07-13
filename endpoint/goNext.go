package endpoint

import (
	"encoding/json"
	"github.com/cheekybits/genny/generic"
	"github.com/labstack/echo"
	"net/http"
	"github.com/Quality-Gamer/qg-manager/conf"
	"github.com/Quality-Gamer/qg-manager/database"
	"github.com/Quality-Gamer/qg-manager/model"
	"strconv"
)

type ResponseNext struct {
	Status   string   `json:"status"`
	DataNext DataNext `json:"response"`
	Message  string   `json:"message"`
}

type DataNext struct {
	End   int            `json:"end"`
	Next  int            `json:"next"`
	Match []generic.Type `json:"match"`
}

func GoNext(c echo.Context) error {
	var res ResponseNext
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)

	if len(c.FormValue("user_id")) > 0 && len(c.FormValue("match_id")) > 0 {
		userId, _ := strconv.Atoi(c.FormValue("user_id"))
		managerId := c.FormValue("match_id")
		week, _ := strconv.Atoi(database.GetKey(conf.GetKeyOccurrence(userId, managerId) + ":" + conf.CurrentWeek))
		week += 1
		match, next, end, err := goToNext(userId, week, managerId)

		res.Status = conf.SuccessCode
		res.Message = conf.SuccessMessage
		var response []generic.Type
		response = append(response, match)

		endGame := 0
		nextGame := 0

		if end {
			endGame = 1
		} else {
			endGame= 0
		}

		if next {
			nextGame = 1
		} else {
			nextGame = 0
		}

		if err {
			res.Status = conf.SuccessCode
			res.Message = conf.ErrorDoesNotExist
		}

		res.DataNext.End = endGame
		res.DataNext.Next = nextGame
		res.DataNext.Match = response

		c.Response().WriteHeader(http.StatusOK)
		return json.NewEncoder(c.Response()).Encode(res)
	} else {
		res.Status = conf.ErrorCode
		res.Message = conf.ErrorInputMessage
		c.Response().WriteHeader(http.StatusBadRequest)
		return json.NewEncoder(c.Response()).Encode(res)
	}
}

func goToNext(userId, week int, matchId string) (model.ManagerMatch, bool, bool, bool) {
	m, exists := model.FindManagerMatch(userId, matchId)

	if week > 8 {
		return m, false, true, false
	}

	if !exists {
		return model.ManagerMatch{}, false, true, true
	}

	success := m.RunGame()
	newMatch, _ := model.FindManagerMatch(userId, matchId)
	
	newMatch.Money = newMatch.Money + newMatch.GameModel.InitialMoney
	newMatch.Time = newMatch.Time + newMatch.GameModel.InitialTime
	keyNoWeek := conf.GetKeyOccurrence(m.UserId, m.Id)
	keyCurrentMoney := keyNoWeek + ":" + conf.CurrentMoney
	keyCurrentTime := keyNoWeek + ":" + conf.CurrentTime
	database.SetKey(keyCurrentMoney, newMatch.Money)
	database.SetKey(keyCurrentTime, newMatch.Time)
	
	end := false

	if newMatch.Week >= 8 {
		end = true
	}

	if success {
		return newMatch, true, false, end
	}

	return newMatch, false, false, end
}
