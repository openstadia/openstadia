package store

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	c "github.com/openstadia/openstadia/config"
	"github.com/openstadia/openstadia/utils"
	"path/filepath"
)

const DbFile = "openstadia.db"

const (
	AppsBucketName = "Apps"
	MetaBucketName = "Meta"
)

type Store struct {
	db *bolt.DB
}

func CreateStore() (*Store, error) {
	appConfigDir, err := utils.GetConfigDir()
	if err != nil {
		panic("can't get config directory")
	}

	dbPath := filepath.Join(appConfigDir, DbFile)

	db, err := bolt.Open(dbPath, 0600, nil)
	if err != nil {
		return nil, err
	}

	db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte(AppsBucketName))
		tx.CreateBucketIfNotExists([]byte(MetaBucketName))
		return nil
	})

	return &Store{db: db}, nil
}

func (s *Store) SetConfig(config *c.Openstadia) {
	if config == nil {
		return
	}

	if s.Apps() == nil {
		for _, app := range config.Apps {
			s.AddApp(&app)
		}
	}

	if s.Hub() == nil {
		s.SetHub(config.Hub)
	}

	if s.Local() == nil {
		s.SetLocal(config.Local)
	}
}

func (s *Store) Config() *c.DbOpenstadia {
	apps := s.Apps()
	hub := s.Hub()
	local := s.Local()

	config := c.DbOpenstadia{
		Openstadia: c.Openstadia{
			Hub:   hub,
			Local: local,
		},
		Apps: apps,
	}

	return &config
}

func (s *Store) Apps() []c.DbApp {
	var apps []c.DbApp

	s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(AppsBucketName))

		cursor := b.Cursor()

		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			var dbApp c.DbApp
			err := json.Unmarshal(v, &dbApp)
			if err != nil {
				continue
			}

			apps = append(apps, dbApp)
		}

		return nil
	})

	return apps
}

func (s *Store) AddApp(app *c.BaseApp) error {
	if app == nil {
		return nil
	}

	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(AppsBucketName))

		id, _ := b.NextSequence()

		dbApp := c.DbApp{
			Id:      int(id),
			BaseApp: *app,
		}

		buf, err := json.Marshal(dbApp)
		if err != nil {
			return err
		}

		return b.Put(itob(dbApp.Id), buf)
	})
}

func (s *Store) Hub() *c.Hub {
	var hub *c.Hub

	s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(MetaBucketName))

		value := b.Get([]byte("Hub"))
		if value == nil {
			return nil
		}

		err := json.Unmarshal(value, &hub)
		if err != nil {
			return err
		}

		return nil
	})

	return hub
}

func (s *Store) SetHub(hub *c.Hub) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(MetaBucketName))

		if hub == nil {
			b.Delete([]byte("Hub"))
			return nil
		}

		buf, err := json.Marshal(hub)
		if err != nil {
			return err
		}

		b.Put([]byte("Hub"), buf)

		return nil
	})
}

func (s *Store) Local() *c.Local {
	var local *c.Local

	s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(MetaBucketName))

		value := b.Get([]byte("Local"))
		if value == nil {
			return nil
		}

		err := json.Unmarshal(value, &local)
		if err != nil {
			return err
		}

		return nil
	})

	return local
}

func (s *Store) SetLocal(local *c.Local) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(MetaBucketName))

		if local == nil {
			b.Delete([]byte("Hub"))
			return nil
		}

		buf, err := json.Marshal(local)
		if err != nil {
			return err
		}

		b.Put([]byte("Local"), buf)

		return nil
	})
}

func (s *Store) GetAppById(id int) (*c.DbApp, error) {
	for _, app := range s.Apps() {
		if app.Id == id {
			return &app, nil
		}
	}

	return nil, fmt.Errorf("no such application: %d", id)
}

func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
