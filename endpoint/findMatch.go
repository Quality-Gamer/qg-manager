package endpoint

import (
	"encoding/json"
	"github.com/labstack/echo"
	"net/http"
	"github.com/Quality-Gamer/qg-manager/conf"
	"github.com/Quality-Gamer/qg-manager/model"
	"strconv"
)

func FindMatch(c echo.Context) error {
	var res model.Response
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)

	if len(c.FormValue("user_id")) > 0 && len(c.FormValue("match_id")) > 0 {
		userId, _ := strconv.Atoi(c.FormValue("user_id"))
		matchId := c.FormValue("match_id")

		match, err := findMatchById(userId, matchId)
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

func findMatchById(userId int, matchId string) (model.ManagerMatch, bool) {
	match,exists := model.FindManagerMatch(userId,matchId)

	if exists {
		return match,false
	}

	return match, true
}

