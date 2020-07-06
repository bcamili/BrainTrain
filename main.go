package main

import (
	"BrainTrain/app/controller"
	"BrainTrain/app/model"
	"BrainTrain/app/route"
)

func main() {

	model.Init()
	controller.Init()
	routes.Routes()

}
