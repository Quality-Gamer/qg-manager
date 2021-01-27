package endpoint

import (
	"encoding/json"
	"github.com/labstack/echo"
	"github.com/Quality-Gamer/qg-manager/model"
)

func CreateOccurrences(c echo.Context) (err error) {
	input := new(model.InputModel)

	if err = c.Bind(input); err != nil {
		return
	}

	gm := input.Body.GameModel
	gameModel := model.GetModel(gm.Id)
	gameModel.ManagerOccurrences = gm.ManagerOccurrences
	gameModel.UserOccurrences = gm.UserOccurrences
	model.CreateOcurrences(gameModel)



	c.Response().Header().Set("Access-Control-Allow-Origin","*")
	c.Response().Header().Set(echo.HeaderContentType,echo.MIMEApplicationJSONCharsetUTF8)
	return json.NewEncoder(c.Response()).Encode(gameModel)
}
