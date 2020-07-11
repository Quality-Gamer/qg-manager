package main

import (
	"github.com/Quality-Gamer/qg-manager/endpoint"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
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

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		//AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

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

	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}
