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

type ContentType int

const (
	Preferences ContentType = iota
	Config
)

type ContainerType int

const (
	XML ContainerType = iota
	JSON
)

type Prefs struct {
	Contents      PList
	ContentType   ContentType
	ContainerType ContainerType
	KV            map[string]any
}

type PList struct {
	XMLName xml.Name `xml:"plist"`
	Version string   `xml:"version,attr"`
	Dict    Dict     `xml:"dict"`
}

type Dict struct {
	XMLName xml.Name `xml:"dict"`
	Key     []string `xml:"key"`
	Date    []string `xml:"date"`
}

func (id *Prefs) Load(fileName string) (bool, error) {
	path, err := id.prefsPath(fileName)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal("Load read failure", err)
		return false, err
	}
	// log.Println(string(data))
	id.Contents, err = id.deserializeXML(data)
	if err != nil {
		log.Fatal("Load deserialize failure", err)
		return false, err
	}
	return true, nil
}

func (id *Prefs) Save(fileName string) (bool, error) {
	path, err := id.prefsPath(fileName)

	data, err := id.serializeXML(path, id.Contents)
	if err != nil {
		log.Fatal("Save serialize failure", err)
		return false, err
	}
	err = ioutil.WriteFile(path, data, 644)
	if err != nil {
		log.Fatal("Save write failure", err)
		return false, err
	}
	return true, nil
}

func (id *Prefs) prefsPath(prefsFile string) (string, error) {
	if id.ContentType != Preferences {
		cwd, err := os.Getwd()
		prefsPath := path.Join(cwd, prefsFile+".plist")
		if err != nil {
			log.Fatal("prefsPath cwd failure", err)
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

	log.Fatal("prefsPath failure", err)
	return "", err
}

func (id *Prefs) deserializeXML(contents []byte) (PList, error) {
	var result PList
	err := xml.Unmarshal(contents, &result)
	if err != nil {
		log.Fatal("deserializeXML failure", err)
		return PList{}, err
	}

	return result, nil
}

func (id *Prefs) serializeXML(prefsPath string, content PList) ([]byte, error) {
	data, err := xml.MarshalIndent(content, "", "        ")
	if err != nil {
		log.Fatal("serializeXML failure", err)
		return nil, err
	}

	data = []byte(xml.Header + DOCTYPE + string(data))
	return data, err
}

func (id *Prefs) deserializeJSON(contents []byte) (PList, error) {
	var result PList
	err := json.Unmarshal(contents, &result)
	if err != nil {
		log.Fatal("deserializeJSON failure", err)
		return PList{}, err
	}

	return result, nil
}

func (id *Prefs) serializeJSON(fileName string, content PList) ([]byte, error) {
	data, err := json.Marshal(content)
	if err != nil {
		log.Fatal("serializeJSON failure", err)
		return nil, err
	}
	return data, err
}
