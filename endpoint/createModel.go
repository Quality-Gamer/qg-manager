package endpoint

import (
	"encoding/json"
	"github.com/labstack/echo"
	"github.com/Quality-Gamer/qg-manager/model"
)

func CreateGameModel(c echo.Context) (err error) {
	input := new(model.InputModel)

	if err = c.Bind(input); err != nil {
		return
	}

	gm := input.Body.GameModel

	gm.Id = model.GerenateHash(gm.BelongsTo)
	model.CreateGameModel(gm)

	c.Response().Header().Set("Access-Control-Allow-Origin","*")
	c.Response().Header().Set(echo.HeaderContentType,echo.MIMEApplicationJSONCharsetUTF8)
	return json.NewEncoder(c.Response()).Encode(gm)
}
