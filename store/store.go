package store

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	c "github.com/openstadia/openstadia/config"
)

const DbFile = "openstadia.db"

const (
	AppsBucketName  = "Apps"
	HubBucketName   = "Hub"
	LocalBucketName = "Local"
)

type Store struct {
	db *bolt.DB
}

func CreateStore() (*Store, error) {
	db, err := bolt.Open(DbFile, 0600, nil)
	if err != nil {
		return nil, err
	}

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

		if b == nil {
			return nil
		}

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
		b, _ := tx.CreateBucketIfNotExists([]byte(AppsBucketName))

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
		b := tx.Bucket([]byte(HubBucketName))

		if b == nil {
			return nil
		}

		addr := b.Get([]byte("Addr"))
		token := b.Get([]byte("Token"))

		hub = &c.Hub{
			Addr:  string(addr),
			Token: string(token),
		}

		return nil
	})

	return hub
}

func (s *Store) SetHub(hub *c.Hub) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		if hub == nil {
			tx.DeleteBucket([]byte(HubBucketName))
			return nil
		}

		b, _ := tx.CreateBucketIfNotExists([]byte(HubBucketName))

		b.Put([]byte("Addr"), []byte(hub.Addr))
		b.Put([]byte("Token"), []byte(hub.Token))

		return nil
	})
}

func (s *Store) Local() *c.Local {
	var local *c.Local

	s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(LocalBucketName))

		if b == nil {
			return nil
		}

		host := b.Get([]byte("Host"))
		port := b.Get([]byte("Port"))

		local = &c.Local{
			Host: string(host),
			Port: string(port),
		}

		return nil
	})

	return local
}

func (s *Store) SetLocal(local *c.Local) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		if local == nil {
			tx.DeleteBucket([]byte(LocalBucketName))
			return nil
		}

		b, _ := tx.CreateBucketIfNotExists([]byte(LocalBucketName))

		b.Put([]byte("Host"), []byte(local.Host))
		b.Put([]byte("Port"), []byte(local.Port))

		return nil
	})
}

func (s *Store) GetAppByName(name string) (*c.DbApp, error) {
	for _, app := range s.Apps() {
		if app.Name == name {
			return &app, nil
		}
	}

	return nil, fmt.Errorf("no such application: %s", name)
}

func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
