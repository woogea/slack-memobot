package mention

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type MemoryStorage struct {
	Db        map[string][]string
	Namespace string
}

func (m *MemoryStorage) restore() {
	// restore from namespace based backup file.
	bytes, err := ioutil.ReadFile("./data/" + m.Namespace + ".json")
	if err != nil {
		log.Println(err)
	}
	err = json.Unmarshal(bytes, m)
	if err != nil {
		log.Println(err)
	}
}

//Backup is required after updating the db.
func (m *MemoryStorage) backup() {
	// backup to namespace based backup file.
	file, err := os.Create("./data/" + m.Namespace + ".json")
	if err != nil {
		log.Println(err)
	}
	defer file.Close()
	err = json.NewEncoder(file).Encode(m)
	if err != nil {
		log.Println(err)
	}
}

func NewMemoryStorage(namespace string) *MemoryStorage {
	// need to restore
	m := &MemoryStorage{Db: map[string][]string{}, Namespace: namespace}
	m.restore()
	return m
}

func (m *MemoryStorage) Get(key string) (string, error) {
	d, ok := m.Db[key]
	if ok {
		return d[0], nil
	}
	return "", &NotfoundError{}
}

func (m *MemoryStorage) List(key string) ([]string, error) {
	d, ok := m.Db[key]
	if ok {
		return d, nil
	}
	return []string{}, &NotfoundError{}
}

func (m *MemoryStorage) Add(key string, value string) error {
	m.Db[key] = append(m.Db[key], value)
	m.backup()
	return nil
}
func (m *MemoryStorage) Remove(key string, value string) error {
	d, ok := m.Db[key]
	if !ok {
		return nil
	}
	var r = -1
	for i := range d {
		if d[i] == value {
			r = i
		}
	}
	if r != -1 {
		m.Db[key] = append(d[:r], d[r+1:]...)
	}
	m.backup()
	return nil
}

func (m *MemoryStorage) RemoveAll(key string) error {
	m.Db[key] = []string{}
	m.backup()
	return nil
}

func (m *MemoryStorage) Rotate(key string) error {
	m.Db[key] = append(m.Db[key][1:], m.Db[key][0])
	m.backup()
	return nil
}
