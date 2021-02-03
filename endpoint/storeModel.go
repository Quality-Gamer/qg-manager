package endpoint

import (
	"encoding/json"
	"github.com/labstack/echo"
	"net/http"
	"github.com/Quality-Gamer/qg-manager/conf"
	"github.com/Quality-Gamer/qg-manager/model"
	"strconv"
	"strings"
)

type normalizedResponse struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Price    int    `json:"price"`
	Type     string `json:"type"`
}

func StoreModel(c echo.Context) error {
	var res model.Response

	if len(c.FormValue("user_id")) > 0 && len(c.FormValue("match_id")) > 0 {
		c.Response().Header().Set("Access-Control-Allow-Origin", "*")
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		user_id, _ := strconv.Atoi(c.FormValue("user_id"))
		matchId := c.FormValue("match_id")
		modelId := model.GetModelId(user_id,matchId)
		mm,_ := model.FindManagerMatch(user_id,matchId)

		// level will replace week in order to deliver the first version as soon as possible
		// the code has not been refactored
		level := mm.Level

		items := loadModelItems(modelId,matchId,user_id,level)
		res.Response = append(res.Response, items)
		res.Message = conf.SuccessMessage
		res.Status = conf.SuccessCode
		c.Response().WriteHeader(http.StatusOK)
		return json.NewEncoder(c.Response()).Encode(res)
	} else {
		res.Status = conf.ErrorCode
		res.Message = conf.ErrorInputMessage
		c.Response().WriteHeader(http.StatusBadRequest)
		return json.NewEncoder(c.Response()).Encode(res)
	}
}

func loadModelItems(modelId,matchId string, userId,week int) []*normalizedResponse{
	var items []*normalizedResponse
	_model := model.GetModel(modelId)
	levels := _model.Levels
	//index := week - 1
	index := model.GetCurrentLevel(userId,week,matchId)
	p := levels[index].Process

	for _, i := range p {
		for _, j := range i.Resources {
			if len(j.Name) > 0 {
				normalized := new(normalizedResponse)
				normalized.Id = j.Id
				normalized.Name = j.Name
				normalized.Type = j.Type
				normalized.Price = j.Price
				items = append(items,normalized)
			}
		}
		for _, k := range i.Activities {
			if len(k.Name) > 0 {
				normalized := new(normalizedResponse)
				normalized.Id = k.Id
				normalized.Name = k.Name
				normalized.Type = "A"
				normalized.Price = k.Time
				items = append(items,normalized)
			}
		}
	}

	var addedItems []string
	var noRepeatedItems []*normalizedResponse

	for _,i := range items {
		if !has(addedItems,i.Name) {
			noRepeatedItems = append(noRepeatedItems,i)
			addedItems = append(addedItems,i.Name)
		}
	}

	return noRepeatedItems
}

func has(slice []string, value string) bool {
	for _,i := range slice {
		if strings.Compare(i,value) == 0 {
			return true
		}
	}

	return false
}


