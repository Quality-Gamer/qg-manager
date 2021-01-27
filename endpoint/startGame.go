package endpoint

import (
	"encoding/json"
	"github.com/labstack/echo"
	"net/http"
	"github.com/Quality-Gamer/qg-manager/conf"
	"github.com/Quality-Gamer/qg-manager/model"
	"strconv"
)

func StartGame(c echo.Context) error {
	var res model.Response
	c.Response().Header().Set("Access-Control-Allow-Origin","*")
	c.Response().Header().Set(echo.HeaderContentType,echo.MIMEApplicationJSONCharsetUTF8)

	if len(c.FormValue("user_id")) > 0 && len(c.FormValue("model_id")) > 0 {
		uid, _ := strconv.Atoi(c.FormValue("user_id"))
		modelId := c.FormValue("model_id")
		match := model.NewMatch(modelId,uid)
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
