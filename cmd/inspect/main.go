package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/captncraig/sci"
	"github.com/captncraig/sci/resources"
)

var dir = flag.String("d", "games/sci0/SierraCard1988", "directory containing resource files")
var print = flag.Bool("p", false, "print a bunch of stuff")

func main() {
	flag.Parse()
	loader := sci.NewFromDir(*dir)
	rMap, err := loader.GetFile("RESOURCE.MAP")
	if err != nil {
		log.Fatalf("Couldn't read resource map: %s", err)
	}
	rez, err := resources.ReadMap(rMap, loader)
	if err != nil {
		log.Fatal(err)
	}
	if !*print {
		return
	}
	for _, r := range rez.AllHeaders {
		fmt.Println(r)
	}
}
