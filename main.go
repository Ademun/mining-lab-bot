package main

import (
	"fmt"
	"log"

	"github.com/Ademun/mining-lab-bot/scraper"
	"github.com/Ademun/mining-lab-bot/service"
)

func main() {
	err := scraper.UpdateServiceIDs()
	if err != nil {
		log.Fatal(err)
	}
	data, err := service.CheckAvailableLabs()
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range data {
		fmt.Println(v)
	}
}
