package main

import (
	. "goprefs"
	"log"
)

func main() {
	prefs := Prefs{ContainerType: XML, ContentType: Config}
	prefs.Load("demo")
	log.Println(prefs.Contents.Dict)
}
