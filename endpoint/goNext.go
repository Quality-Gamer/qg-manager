package endpoint

import (
	"encoding/json"
	"fmt"
	"qg-manager/conf"
	"qg-manager/database"
	"qg-manager/model"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
)

func GoNext(c echo.Context) error {
	var res model.Response
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)

	if len(c.FormValue("user_id")) > 0 && len(c.FormValue("match_id")) > 0 {
		userId, _ := strconv.Atoi(c.FormValue("user_id"))
		managerId := c.FormValue("match_id")
		week, _ := strconv.Atoi(database.GetKey(conf.GetKeyOccurrence(userId,managerId) + ":" + conf.CurrentWeek))
		week += 1
		match,end := goToNext(userId,week,managerId)
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

func goToNext(userId,week int, managerId string) (model.ManagerMatch,bool) {
	m, exists := model.FindManagerMatch(userId,managerId)

	if week > 8 || !exists {
		return model.ManagerMatch{},true
	}

	fmt.Println(m.RunGame())

	return m,false
}
