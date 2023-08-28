package goprefs

import (
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path"
	"runtime"
)

const (
	DOCTYPE = `<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
	` + "\n"
)

type FileType int

const (
	PreferenceXML FileType = iota
	ConfigXML
	ConfigJSON
)

type Prefs struct {
	contents PList
	data     map[string]interface{}
	FileType FileType
}

type PList struct {
	Plist   string                 `xml:"plist"`
	Version string                 `xml:"version,attr"`
	Dict    map[string]interface{} `xml:"dict"`
}

func (id *Prefs) Load(fileName string) (bool, error) {
	path, err := id.prefsPath(fileName)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
		return false, err
	}
	id.contents, err = id.deserializeXML(data)
	if err != nil {
		log.Fatal(err)
		return false, err
	}
	return true, nil
}

func (id *Prefs) Save(fileName string) (bool, error) {
	path, err := id.prefsPath(fileName)

	data, err := id.serializeXML(path, id.contents)
	if err != nil {
		return false, err
	}
	ioutil.WriteFile(path, data, 644)
	return true, nil
}

func (id *Prefs) prefsPath(prefsFile string) (string, error) {
	if id.FileType != PreferenceXML {
		cwd, err := os.Getwd()
		prefsPath := path.Join(cwd, prefsFile+".plist")
		log.Println(prefsPath)
		if err != nil {
			return "", err
		}
		return prefsPath, nil
	}

	usr, err := user.Current()
	if err == nil {
		var prefsPath string
		switch runtime.GOOS {
		case "darwin":
			prefsPath = path.Join(usr.HomeDir, "Library", "Preferences", prefsFile+".plist")
		default:
			prefsPath = path.Join(usr.HomeDir, prefsFile+".plist")
		}

		_, err = os.Stat(prefsPath)
		if os.IsNotExist(err) {
			os.MkdirAll(prefsPath, os.ModePerm)
		}

		result := path.Join(prefsPath, prefsFile)

		return result, nil
	}

	return "", err
}

func (id *Prefs) deserializeXML(contents []byte) (PList, error) {
	var result PList
	err := xml.Unmarshal(contents, &result)
	if err != nil {
		return PList{}, err
	}

	return result, nil
}

func (id *Prefs) serializeXML(prefsPath string, content PList) ([]byte, error) {
	data, err := xml.MarshalIndent(content, "", "        ")
	data = []byte(xml.Header + DOCTYPE + string(data))

	err = ioutil.WriteFile(prefsPath, data, 644)
	return data, err
}

func (id *Prefs) deserializeJSON(contents []byte) (result map[string]interface{}, err error) {
	err = json.Unmarshal(contents, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (id *Prefs) serializeJSON(fileName string, content map[string]interface{}) ([]byte, error) {
	data, err := json.Marshal(content)
	return data, err
}
