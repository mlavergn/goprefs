package main

import (
	. "goprefs"
	"log"
)

func main() {
	prefs := Prefs{}
	prefs.FileType = ConfigXML
	prefs.Load("demo")
	log.Println(prefs)
}
