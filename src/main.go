package main

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
)

func main() {
	app := iris.Default()

	api := app.Party("/api")
	{
		v1 := api.Party("/1.0")
		{
			hero.Register(&compositeSettingsRetrievalService{})
			settingsRetrievalHandler := hero.Handler(getSettings)
			v1.Get("/{path:path}", settingsRetrievalHandler)
		}
	}

	app.Run(iris.Addr(":8080"))
}
