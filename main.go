package main

import (
	"github.com/Quality-Gamer/qg-manager/endpoint"
	"github.com/labstack/echo"
	"os"
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
	e.POST("/api/create", endpoint.Create)

	//Find
	e.POST("/api/find", endpoint.Find)

	//Store
	e.POST("/api/store", endpoint.Store)

	//Transaction
	e.POST("/api/transaction", endpoint.Transaction)

	//Next
	e.POST("/api/next", endpoint.Next)

	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}
