package endpoint

import (
	"encoding/json"
	"github.com/cheekybits/genny/generic"
	"github.com/labstack/echo"
	"net/http"
	"github.com/Quality-Gamer/qg-manager/conf"
	"github.com/Quality-Gamer/qg-manager/database"
	"github.com/Quality-Gamer/qg-manager/model"
	"strconv"
)

func TransactionModel(c echo.Context) error {
	var res model.Response
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)

	if len(c.FormValue("user_id")) > 0 && len(c.FormValue("match_id")) > 0 && len(c.FormValue("item")) > 0 && len(c.FormValue("type")) > 0 {
		userId, _ := strconv.Atoi(c.FormValue("user_id"))
		managerId := c.FormValue("match_id")
		week, _ := strconv.Atoi(database.GetKey(conf.GetKeyOccurrence(userId,managerId) + ":" + conf.CurrentWeek))
		modelId := model.GetModelId(userId,managerId)
		item := c.FormValue("item")
		_type := c.FormValue("type")
		success, errMessage := makeTransactionModel(userId, week, modelId, managerId, item, _type)
		r := make(map[string]generic.Type)

		r["done"] = 0
		//1 - insufficient funds 2 - this item does not exist in this type of transaction
		r["message"] = errMessage

		if success {
			r["done"] = 1
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

func makeTransactionModel(userId, week int, modelId, matchId, item, _type string) (bool, string) {
	_, _, itemPrice := getItem(userId, week, modelId, item, _type)
	float64ItemPrice, _ := strconv.ParseFloat(itemPrice,64)
	intItemPrice, _ := strconv.Atoi(itemPrice)

	if intItemPrice == 0 {
		return false,getMessage(4)
	}

	if _type == "R" {
		value := model.GetCurrentMoney(userId,matchId)
		if value >= float64ItemPrice {
			newValue := value - float64ItemPrice
			key := conf.GetKeyOccurrence(userId,matchId) + ":" + conf.CurrentMoney
			database.SetKey(key,newValue)
			model.AddResource(userId,matchId,item)
			return true,getMessage(0)
		}
		return false,getMessage(1)
	} else if _type == "A" {
		value := model.GetCurrentTime(userId,matchId)
		if value >= intItemPrice {
			newValue := value - intItemPrice
			key := conf.GetKeyOccurrence(userId,matchId) + ":" + conf.CurrentTime
			database.SetKey(key,newValue)
			model.AddActivity(userId,matchId,item)
			return true,getMessage(0)
		}

		return false,getMessage(2)
	}

	return false,getMessage(3)
}

func getMessage(messageType int) string {
	var arrayMessage = [5]string{conf.SuccessMessage, "Saldo insuficiente", "Tempo insuficiente", "Tipo inválido", "Requisição inválida"}
	return arrayMessage[messageType]
}

func getItem(userId,week int, modelId, item, _type string) (string,string,string) {
	var tp,price string

	if _type == "A" {
		tp = conf.Activity
		price = conf.Time
	} else if _type == "R" {
		tp = conf.Resource
		price = conf.Price
	}

	lv := model.GetCurrentLevel(userId, week, modelId)
	strLv := strconv.Itoa(lv)
	gm := model.GetModel(modelId)
	indexMax := len(gm.Levels)

	for i := 0; i < indexMax; i++ {
		key := conf.Game + ":" + conf.Manager + ":" + conf.Model + ":" + modelId + ":" + conf.Level + ":" + strLv + ":" + conf.Process + ":" + strconv.Itoa(i) + ":" + tp + ":" + item + ":"
		keyName := key + conf.Name
		keyPrice := key + price
		keyId := key + conf.Identifier

		nm := database.GetKey(keyName)
		pc := database.GetKey(keyPrice)
		id := database.GetKey(keyId)

		intPc,_ := strconv.Atoi(pc)
		if intPc > 0 {
			return id,nm,pc
		}
	}


	return "","",""
}