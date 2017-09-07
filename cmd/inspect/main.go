package main

import (
	"flag"
	"log"

	"github.com/captncraig/sci"
	"github.com/captncraig/sci/resources"
)

var dir = flag.String("d", "games/sci0/SierraCard1988", "directory containing resource files")

func main() {
	flag.Parse()
	loader := sci.NewFromDir(*dir)
	rMap, err := loader.GetFile("RESOURCE.MAP")
	if err != nil {
		log.Fatalf("Couldn't read resource map: %s", err)
	}
	resources.ReadMap(rMap, loader)
}
