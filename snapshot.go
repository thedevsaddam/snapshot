package snapshot

import (
	"encoding/gob"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const (
	extension = ".gob"
	prefix    = "_sc_"
)

var mutexList = make(map[string]*sync.RWMutex)

// Collection describes a collection of key value pairs
type Collection struct {
	mutex sync.Mutex
	dir   string
	items []string
}

// New return a instance of collection
func New(name string) (*Collection, error) {
	if len(name) <= 0 {
		return &Collection{}, errors.New("Collection name can not be empty")
	}
	//make file path correct
	dir := prefix + filepath.Clean(name)
	collection := Collection{
		dir: dir,
	}
	//create directory with collection name
	return &collection, os.MkdirAll(dir, 0755)
}

//Put store a new key with value in the collection
func (c *Collection) Put(key string, value interface{}) error {
	if len(key) <= 0 {
		return errors.New("Key can not be empty")
	}
	path := filepath.Join(c.dir, key+extension)
	// need get read lock and write lock
	m := c.getPathMutex(path)
	m.Lock()
	defer m.Unlock()
	file, err := os.Create(path)
	defer file.Close()
	if err == nil {
		encoder := gob.NewEncoder(file)
		encoder.Encode(value)
	}
	return err
}

//Get retrieve a value from collection by key
func (c *Collection) Get(key string, value interface{}) error {
	if len(key) <= 0 {
		return errors.New("Key can not be empty")
	}
	path := filepath.Join(c.dir, key+extension)
	m := c.getPathMutex(path)
	m.RLock()
	defer m.RUnlock()
	if !c.Has(key) {
		return fmt.Errorf("Key %s does not exist", key)
	}
	file, err := os.Open(path)
	defer file.Close()
	if err == nil {
		decoder := gob.NewDecoder(file)
		err = decoder.Decode(value)
	}
	return err
}

//Remove delete a key from collection
func (c *Collection) Remove(key string) error {
	if len(key) <= 0 {
		return errors.New("Key can not be empty")
	}
	path := filepath.Join(c.dir, key+extension)
	m := c.getPathMutex(path)
	m.Lock()
	defer m.Unlock()
	if c.Has(key) {
		return os.Remove(path)
	}
	return fmt.Errorf("Key %s does not exist", key)
}

//Flush delete a collection with its value
func (c *Collection) Flush() error {
	if _, err := os.Stat(c.dir); err == nil {
		os.RemoveAll(c.dir)
		return err
	}
	return nil
}

//Has check a key exist in the collection
func (c *Collection) Has(key string) bool {
	if len(key) <= 0 {
		return false
	}
	path := filepath.Join(c.dir, key+extension)
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return true
	}
	return false
}

//List fetch all items key in collection
func (c *Collection) List() ([]string, error) {
	var items []string
	files, err := ioutil.ReadDir(c.dir)
	if err != nil {
		return items, err
	}
	for _, f := range files {
		item := f.Name()
		item = strings.Trim(item, extension)
		items = append(items, item)
	}
	return items, err
}

//TotalItem return total item count
func (c *Collection) TotalItem() int {
	list, _ := c.List()
	return len(list)
}

// populate a package level mutex list
// with key of full path of an item
func (c *Collection) getPathMutex(path string) *sync.RWMutex {
	m, ok := mutexList[path]
	if !ok {
		c.mutex.Lock()
		defer c.mutex.Unlock()
		if m == nil {
			m = &sync.RWMutex{}
		}
		mutexList[path] = m
	}
	return m
}
