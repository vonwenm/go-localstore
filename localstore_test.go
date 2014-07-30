package localstore

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/user"
	"path"
	"testing"
	"time"
)

var testFolder string

func TestNew(t *testing.T) {

	// Make it random so we make sure it doesn't exist and don't delete any possibly important folders
	testFolder = fmt.Sprintf(".testapp%d", time.Now().Unix())

	_, err := New(testFolder, "config")

	if err != nil {
		t.Error(err)
	}
}

func TestSet(t *testing.T) {

	l, _ := New(testFolder, "config")

	err := l.SetDefault("hello", "world")

	if err != nil {
		t.Error(err)
	}

	usr, err := user.Current()

	content, err := ioutil.ReadFile(path.Join(usr.HomeDir, testFolder, "config.json"))

	if err != nil {
		t.Error(err)
	}

	var m map[string]interface{}

	json.Unmarshal(content, &m)

	if m["hello"] != "world" {

		t.Error("Expected hello key to have world as value, has", m["hello"])

	}
}

func TestGet(t *testing.T) {

	l, _ := New(testFolder, "config")

	err := l.SetDefault("hello", "world")

	if err != nil {
		t.Error(err)
	}

	value, err := l.GetDefault("hello")

	if err != nil {
		t.Error(err)
	}

	valueStr, ok := value.(string)

	if !ok {
		t.Error("Value should be of type string")
	}

	if valueStr != "world" {
		t.Error("Value should equal world, now equals", valueStr)
	}

	// Test not existing key

	_, err = l.GetDefault("existsnot")

	if err != ErrNotFound {
		t.Error("Should return not exists error")
	}

}

type sampleConfig struct {
	Name string
	Host string
	Port int
}

func TestStore(t *testing.T) {

	origSample := sampleConfig{
		"Google",
		"www.google.com",
		80}

	l, _ := New(testFolder, "config")

	err := l.Store("testsave", origSample)

	if err != nil {
		t.Error(err)
	}

	usr, err := user.Current()

	content, err := ioutil.ReadFile(path.Join(usr.HomeDir, testFolder, "testsave.json"))

	if err != nil {
		t.Error(err)
	}

	var newSample sampleConfig

	json.Unmarshal(content, &newSample)

	eq := origSample.Host == newSample.Host && origSample.Name == newSample.Name && origSample.Port == newSample.Port

	if !eq {

		t.Error("Original sample does not equal stored sample.", newSample)
	}
}

func TestLoad(t *testing.T) {

	origSample := sampleConfig{
		"Google",
		"www.google.com",
		80}

	l, _ := New(testFolder, "config")

	err := l.StoreDefault(origSample)

	if err != nil {
		t.Error(err)
	}

	var newSample sampleConfig

	err = l.LoadDefault(&newSample)

	eq := origSample.Host == newSample.Host && origSample.Name == newSample.Name && origSample.Port == newSample.Port

	if !eq {

		t.Error("Original sample does not equal stored sample.", newSample)
	}
}
