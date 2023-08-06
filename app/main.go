package main

import (
	"fmt"

	"github.com/thearjnep/rag-bot/bot"
	"github.com/thearjnep/rag-bot/config"
)

func main() {
	err := config.ReadConfig()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	bot.Initialize()

	<-make(chan struct{})
	return
}
