package main

import (
	//"github.com/Quality-Gamer/qg-manager/endpoint"
	"github.com/labstack/echo"
	"os"
	"qg-manager/endpoint"
)

//var configuration = conf.Configuration{}

//func init() {
//	err := gonfig.GetConf("./conf/conf.json", &configuration)
//
//	if err != nil {
//		panic(err)
//	}
//}

//main contains all API endpoints
func main() {
	e := echo.New()

	//Create
	e.GET("/api/create", endpoint.Create)

	//Find
	e.GET("/api/find", endpoint.Find)

	//Store
	e.GET("/api/store", endpoint.Store)

	//Transaction
	e.GET("/api/transaction", endpoint.Transaction)

	//Next
	e.GET("/api/next", endpoint.Next)

	//Below there are new endpoints with the new business rules

	//Create/Model
	e.POST("/api/create/model", endpoint.CreateGameModel)

	//Create/Match
	e.POST("/api/create/match", endpoint.StartGame)

	//Get/Store
	e.GET("/api/get/store", endpoint.StoreModel)

	//Create/Transaction
	e.POST("/api/create/transaction", endpoint.TransactionModel)

	//Get/Match
	e.GET("/api/get/match", endpoint.FindMatch)

	//Debug
	e.POST("/api/debug", endpoint.Debug)
	e.GET("/api/debug", endpoint.Debug)

	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}
