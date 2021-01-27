package endpoint

import (
	"encoding/json"
	"github.com/labstack/echo"
	"github.com/Quality-Gamer/qg-manager/model"
)

func Debug(c echo.Context) (err error) {
	input := new(model.InputGame)

	if err = c.Bind(input); err != nil {
		return
	}

	id := c.FormValue("id")
	gm := model.GetGameModel(id)
	//mm := input.Body.ManagerMatch
	//mm.Level = 1
	//mm.RunGame()


	c.Response().Header().Set("Access-Control-Allow-Origin","*")
	c.Response().Header().Set(echo.HeaderContentType,echo.MIMEApplicationJSONCharsetUTF8)
	return json.NewEncoder(c.Response()).Encode(gm)
}