package main

import (
	"log"
	"os"

	"github.com/Apricot-S/mjai-manue-go/internal/agent"
	"github.com/Apricot-S/mjai-manue-go/internal/client"
)

func main() {
	name := "tsumogiri"
	room := "default"

	agent := agent.NewTsumogiriAgent(name, room)

	client := client.NewClient(os.Stdin, os.Stdout, true, agent)

	if err := client.Run(); err != nil {
		log.Fatalf("error running client: %v", err)
	}
}
