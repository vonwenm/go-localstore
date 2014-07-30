package localstore

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"path"
)

var ErrNotFound = errors.New("Key not found.")

// Creates a new JSON store. It will create a directory named appDir in the user's home directory.
func New(appDir string, defaultStoreName string) (*JsonStore, error) {

	u, err := user.Current()

	if err != nil {
		return nil, err
	}

	dirPath := path.Join(u.HomeDir, appDir)

	if _, err := os.Stat(dirPath); err != nil {

		if os.IsNotExist(err) {

			err = os.Mkdir(dirPath, os.ModePerm)

			if err != nil {
				return nil, err
			}

		} else {

			return nil, err

		}
	}

	return &JsonStore{dirPath, defaultStoreName}, nil
}

type JsonStore struct {
	path             string
	defaultStoreName string
}

// Loads the content of the storeName JSON file into value.
func (js JsonStore) Load(storeName string, value interface{}) error {

	return read(js.getPath(storeName), value)
}

// Loads the content of the default JSON file into value.
func (js JsonStore) LoadDefault(value interface{}) error {

	return js.Load(js.defaultStoreName, value)
}

// Stores the content of value in the storeName JSON file.
func (js JsonStore) Store(storeName string, value interface{}) error {

	return write(js.getPath(storeName), value)
}

// Stores the content of value in the default JSON file.
func (js JsonStore) StoreDefault(value interface{}) error {

	return js.Store(js.defaultStoreName, value)
}

// Returns the value for key in the storeName JSON file. Returns ErrNotFound is key hasn't been found.
func (js JsonStore) Get(storeName string, key string) (interface{}, error) {

	m := map[string]interface{}{}

	err := read(js.getPath(storeName), &m)

	if err != nil {

		return nil, err
	}

	if value, ok := m[key]; ok {

		return value, nil

	}

	return nil, ErrNotFound
}

// Returns the value for key in the default JSON file. Returns ErrNotFound is key hasn't been found.
func (js JsonStore) GetDefault(key string) (interface{}, error) {
	return js.Get(js.defaultStoreName, key)
}

// Sets value as value for key in the storeName JSON file.
func (js JsonStore) Set(storeName string, key string, value interface{}) error {

	m := map[string]interface{}{}

	err := read(js.getPath(storeName), &m)

	// Can be ignored if EOF - simply use new map with new key
	if err != nil && err != io.EOF {

		return err
	}

	m[key] = value

	return write(js.getPath(storeName), m)
}

// Sets value as value for key in the default JSON file.
func (js JsonStore) SetDefault(key string, value interface{}) error {
	return js.Set(js.defaultStoreName, key, value)
}

func (js JsonStore) getPath(name string) string {
	return path.Join(js.path, name+".json")
}

func getFile(path string) (*os.File, error) {

	f, err := os.OpenFile(path, os.O_RDWR, 0666)

	if err != nil {

		if _, ok := err.(*os.PathError); ok {

			return os.Create(path)

		}

		return nil, err
	}

	return f, nil
}

func read(path string, value interface{}) error {

	f, err := getFile(path)
	defer f.Close()

	if err != nil {
		return err
	}

	return json.NewDecoder(bufio.NewReader(f)).Decode(&value)
}

func write(path string, value interface{}) error {

	content, err := json.Marshal(value)

	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, content, os.ModePerm)
}
