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

	//Run/Game
	e.POST("/api/run/game", endpoint.GoNext)

	//Create/Occurrence
	e.POST("/api/create/occurrences", endpoint.CreateOccurrences)

	//Debug
	//e.POST("/api/debug", endpoint.Debug)
	//e.GET("/api/debug", endpoint.Debug)

	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}
