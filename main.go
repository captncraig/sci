package main

import (
	"path/filepath"

	"github.com/captncraig/sci/resources"
)

func main() {
	dir := filepath.Join("games", "sci0", "SierraCard1988")
	resources.ReadMap(dir)
}
