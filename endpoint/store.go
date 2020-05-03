package endpoint

import (
	"encoding/json"
	"github.com/labstack/echo"
	"manager/conf"
	"manager/model"
	"net/http"
)

func Store(c echo.Context) error {
	var res model.Response
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	items := loadAllItems()
	res.Response = append(res.Response, items)
	res.Message = conf.SuccessMessage
	res.Status = conf.SuccessCode
	c.Response().WriteHeader(http.StatusOK)
	return json.NewEncoder(c.Response()).Encode(res)
}

func loadAllItems() map[string]int{
	items := make(map[string]int)
	items[conf.Tester] = conf.PriceTester
	items[conf.Frontend] = conf.PriceFrontend
	items[conf.Backend] = conf.PriceBackend
	items[conf.Designer] = conf.PriceDesigner
	items[conf.RequirementAnalyst] = conf.PriceRiskAnalyst
	items[conf.ProductOwner] = conf.PriceProductOwner
	items[conf.Ide] = conf.PriceIde
	items[conf.DesignSoftware] = conf.PriceDesignSoftware

	items[conf.Scrum] = conf.TimeScrum
	items[conf.CustomerContact] = conf.TimeCustomerContact
	items[conf.Delivery] = conf.TimeDelivery
	items[conf.RiskAnalysis] = conf.TimeRiskAnalysis

	return items
}

