package endpoint

import (
	"encoding/json"
	"github.com/labstack/echo"
	"github.com/Quality-Gamer/qg-manager/conf"
	"github.com/Quality-Gamer/qg-manager/database"
	"github.com/Quality-Gamer/qg-manager/model"
	"net/http"
	"strconv"
)

func Transaction(c echo.Context) error {
	var res model.Response
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)

	if len(c.FormValue("user_id")) > 0 && len(c.FormValue("manager_id")) > 0 && len(c.FormValue("item")) > 0 && len(c.FormValue("type")) > 0 {
		userId, _ := strconv.Atoi(c.FormValue("user_id"))
		managerId := c.FormValue("manager_id")
		week, _ := strconv.Atoi(database.GetKey(conf.GetKeyOccurrence(userId,managerId) + ":" + conf.CurrentWeek))
		item := c.FormValue("item")
		_type := c.FormValue("type")
		err, errCode := makeTransaction(userId, week, managerId, item, _type)
		r := make(map[string]int)

		r["done"] = 1

		if err {
			r["done"] = 0
			//1 - insufficient funds 2 - this item does not exist in this type of transaction
			r["errCode"] = errCode
		}

		res.Response = append(res.Response, r)
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

func makeTransaction(userId, week int, managerId, item, _type string) (bool, int) {
	items := loadAllItems()
	var newUnit int
	var moneyOrTime string

	if _type == "A" {
		moneyOrTime = conf.CurrentTime
	} else {
		moneyOrTime = conf.CurrentMoney
	}

	unit, _ := strconv.Atoi(database.GetKey(conf.GetKeyOccurrence(userId,managerId) + ":" + moneyOrTime))

	if _type == "A" && !isAction(item) {
		return true,2
	} else if _type == "L" && !isLicense(item) {
		return true,2
	} else if _type == "T" && !isTeam(item) {
		return true,2
	}

	if unit >= items[item] {
		newUnit = unit - items[item]
		database.SetKey(conf.GetKeyOccurrence(userId,managerId) + ":" + moneyOrTime,newUnit)
		key := getItemKey(userId, week, managerId, item, _type)
		beforeUpdate,_ := strconv.Atoi(database.GetKey(key))
		afterUpdate := beforeUpdate + 1
		database.SetKey(key,afterUpdate)
		return false,0
	}

	return true,1
}

func getItemKey(userId, week int, managerId, item, _type string) string {
	var tp string = conf.Team

	if _type == "A"{
		tp = conf.Action
	} else if _type == "L" {
		tp = conf.License
	}

	return conf.GetKeyManager(userId,week,managerId) + ":" + tp + ":" + item
}

func isTeam(item string) bool {
	team := make(map[string]int)
	team[conf.Backend] = 1
	team[conf.Frontend] = 1
	team[conf.Designer] = 1
	team[conf.ProductOwner] = 1
	team[conf.RequirementAnalyst] = 1
	team[conf.Tester] = 1

	if _, ok := team[item]; ok {
		return true
	}

	return false
}

func isLicense(item string) bool {
	license := make(map[string]int)
	license[conf.Ide] = 1
	license[conf.DesignSoftware] = 1


	if _, ok := license[item]; ok {
		return true
	}

	return false
}

func isAction(item string) bool {
	action := make(map[string]int)
	action[conf.Scrum] = 1
	action[conf.CustomerContact] = 1
	action[conf.RiskAnalysis] = 1
	action[conf.Delivery] = 1

	if _, ok := action[item]; ok {
		return true
	}

	return false
}

