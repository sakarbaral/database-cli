package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/sakarbaral/database/models"

	"github.com/jcelliott/lumber"
)

type Driver struct {
	mutex   sync.Mutex
	mutexes map[string]*sync.Mutex
	dir     string
	log     models.Logger
}

func New(dir string, options *models.Options) (*Driver, error) {
	dir = filepath.Clean(dir)

	if options == nil {
		options = &models.Options{}
	}

	if options.Logger == nil {
		options.Logger = lumber.NewConsoleLogger(lumber.DEBUG)
	}

	driver := Driver{
		dir:     dir,
		mutexes: make(map[string]*sync.Mutex),
		log:     options.Logger,
	}

	if _, err := os.Stat(dir); err == nil {
		options.Logger.Debug("Using %s as it already exists", dir)
		return &driver, nil
	}
	options.Logger.Debug("Creating the db at %s", dir)
	return &driver, os.MkdirAll(dir, 0755)
}

func (d *Driver) Write(collection, resource string, v interface{}) error {
	if collection == "" || resource == "" {
		return fmt.Errorf("collection and resource must be provided")
	}

	mutex := d.getOrCreateMutex(collection)
	mutex.Lock()
	defer mutex.Unlock()

	dir := filepath.Join(d.dir, collection)
	finalPath := filepath.Join(dir, resource+".json")
	tmpPath := finalPath + ".tmp"

	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	b, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return err
	}

	b = append(b, byte('\n'))

	if err := os.WriteFile(tmpPath, b, 0644); err != nil {
		return err
	}

	return os.Rename(tmpPath, finalPath)
}

func (d *Driver) Read(collection, resource string, v interface{}) error {
	if collection == "" || resource == "" {
		return fmt.Errorf("collection and resource must be provided")
	}

	record := filepath.Join(d.dir, collection, resource+".json")
	b, err := os.ReadFile(record)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, v)
}

func (d *Driver) ReadAll(collection string) ([]string, error) {
	if collection == "" {
		return nil, fmt.Errorf("collection must be provided")
	}

	dir := filepath.Join(d.dir, collection)
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var records []string
	for _, file := range files {
		b, err := os.ReadFile(filepath.Join(dir, file.Name()))
		if err != nil {
			return nil, err
		}
		records = append(records, string(b))
	}

	return records, nil
}

func (d *Driver) Delete(collection, resource string) error {
	if collection == "" || resource == "" {
		return fmt.Errorf("collection and resource must be provided")
	}

	mutex := d.getOrCreateMutex(collection)
	mutex.Lock()
	defer mutex.Unlock()

	record := filepath.Join(d.dir, collection, resource+".json")

	if fi, err := os.Stat(record); err == nil {
		if fi.IsDir() {
			return os.RemoveAll(record)
		}
		return os.Remove(record)
	} else if os.IsNotExist(err) {
		return fmt.Errorf("record %s does not exist", record)
	}

	return nil
}

func (d *Driver) getOrCreateMutex(collection string) *sync.Mutex {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	mutex, ok := d.mutexes[collection]

	if !ok {
		mutex = &sync.Mutex{}
		d.mutexes[collection] = mutex
	}
	return mutex
}
