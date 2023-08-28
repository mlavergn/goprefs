package goprefs

import (
	"testing"
)

func TestPrefs(t *testing.T) {
	prefs := Prefs{ContainerType: XML, ContentType: Config}
	prefs.Load("demo/demo")
	if len(prefs.Contents.Dict.Key) == 0 {
		t.Fatal("Expected to find valid values")
	}
}
