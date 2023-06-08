package pg

import (
	"context"
	"errors"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/spf13/viper"
	"os"
	"ozon-test-unzhakov/internal/config"
	"ozon-test-unzhakov/internal/model"
	"ozon-test-unzhakov/internal/storage/storage"
	"path/filepath"
)

type linkStorage struct {
	s *pg.DB
}

func NewLinkStorage() (storage.LinkStorage, error) {
	err := config.InitConfig(filepath.Join("..", "..", "..", "config"), "config", "yaml")
	if err != nil {
		return nil, err
	}
	connectOptions, err := pg.ParseURL("postgres://" +
		os.Getenv("POSTGRES_USER") + ":" +
		os.Getenv("POSTGRES_PASSWORD") + "@" +
		viper.GetString("postgres.host") + ":" +
		viper.GetString("postgres.port") + "/" +
		viper.GetString("postgres.name") + "?sslmode=disable")
	if err != nil {
		return nil, err
	}
	s := pg.Connect(connectOptions)
	err = s.Ping(context.Background())
	if err != nil {
		return nil, err
	}
	return &linkStorage{s: s}, nil
}

func (ls *linkStorage) Migrate() error {
	if ls.s != nil {
		err := ls.s.Model((*model.Link)(nil)).CreateTable(&orm.CreateTableOptions{
			IfNotExists: true,
		})
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("linkStoragePg: initialized incorrectly")
}

func (ls *linkStorage) GetLink(l string) (*model.Link, error) {
	if ls.s != nil {
		link := &model.Link{}
		err := ls.s.Model(link).
			Where("link = ?", l).
			Select(link)
		if err != nil {
			return nil, err
		}
		return link, err
	}
	return nil, errors.New("linkStoragePg: initialized incorrectly")
}

func (ls *linkStorage) CreateLink(l *model.Link) (*model.Link, error) {
	if ls.s != nil {
		_, err := ls.s.Model(l).Insert(l)
		if err != nil {
			return nil, err
		}
		return l, err
	}
	return nil, errors.New("linkStoragePg: initialized incorrectly")
}

func (ls *linkStorage) UpdateLink(l *model.Link) (*model.Link, error) {
	if ls.s != nil {
		_, err := ls.s.Model(l).Where("id = ?", l.Id).Update(l)
		if err != nil {
			return nil, err
		}
		return l, err
	}
	return nil, errors.New("linkStoragePg: initialized incorrectly")
}

func (ls *linkStorage) DeleteLink(l string) error {
	if ls.s != nil {
		link := &model.Link{}
		_, err := ls.s.Model(link).
			Where("link = ?", l).
			Delete()
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("linkStoragePg: initialized incorrectly")
}
