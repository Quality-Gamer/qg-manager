package endpoint

import (
	"encoding/json"
	"qg-manager/conf"
	"qg-manager/model"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
)

type normalizedResponse struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Price    int    `json:"price"`
	Type     string `json:"type"`
}

func StoreModel(c echo.Context) error {
	var res model.Response

	if len(c.FormValue("week")) > 0 && len(c.FormValue("model_id")) > 0 {
		c.Response().Header().Set("Access-Control-Allow-Origin", "*")
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		week, _ := strconv.Atoi(c.FormValue("week"))
		modelId := c.FormValue("model_id")
		items := loadModelItems(modelId,week)
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

func loadModelItems(modelId string, week int) []*normalizedResponse{
	var items []*normalizedResponse
	model := model.GetModel(modelId)
	levels := model.Levels
	index := week - 1
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

	return items
}


