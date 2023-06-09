package pg

import (
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"ozon-test-unzhakov/internal/model"
	"ozon-test-unzhakov/internal/storage/pg/errors"
	"ozon-test-unzhakov/internal/storage/storage"
)

func init() {

}

type linkStorage struct {
	s *pg.DB
}

func NewLinkStorage(s *pg.DB) (storage.LinkStorage, error) {
	if s == nil {
		return nil, errors.InitializationErr{Err: errors.IncorrectInitializationErr{}}
	}
	ls := &linkStorage{s: s}
	if err := ls.Migrate(); err != nil {
		return nil, err
	}
	return ls, nil
}

func (ls *linkStorage) Migrate() error {
	if ls == nil || ls.s != nil {
		err := ls.s.Model((*model.Link)(nil)).CreateTable(&orm.CreateTableOptions{
			IfNotExists: true,
		})
		if err != nil {
			return errors.ModelErr{Err: err}
		}
		return nil
	}
	return errors.InitializationErr{Err: errors.InitializedIncorrectlyErr{}}
}

func (ls *linkStorage) Get(l *model.Link) (*model.Link, error) {
	if ls.s != nil {
		var err error
		link := &model.Link{}
		if l.Link != "" && l.Code != "" {
			err = ls.s.Model(link).
				Where("link = ? and code = ?", l.Link, l.Code).
				Select(link)
		} else {
			err = ls.s.Model(link).
				Where("link = ? or code = ?", l.Link, l.Code).
				Select(link)
		}
		if err != nil {
			return nil, errors.ModelErr{Err: err}
		}
		return link, err
	}
	return nil, errors.InitializationErr{Err: errors.InitializedIncorrectlyErr{}}
}

func (ls *linkStorage) CreateLink(l *model.Link) (*model.Link, error) {
	if ls.s != nil {
		_, err := ls.s.Model(l).Insert(l)
		if err != nil {
			return nil, errors.ModelErr{Err: err}
		}
		return l, err
	}
	return nil, errors.InitializationErr{Err: errors.InitializedIncorrectlyErr{}}
}

func (ls *linkStorage) UpdateLink(l *model.Link) (*model.Link, error) {
	if ls.s != nil {
		var err error
		_, err = ls.s.Model(l).Where("id = ?", l.Id).Update(l)
		if err != nil {
			return nil, errors.ModelErr{Err: err}
		}
		return l, err
	}
	return nil, errors.InitializationErr{Err: errors.InitializedIncorrectlyErr{}}
}

func (ls *linkStorage) DeleteLink(l *model.Link) error {
	if ls.s != nil {
		var err error
		link := &model.Link{}
		err = ls.s.Model(l).Select()
		if err != nil {
			return errors.ModelErr{Err: err}
		}
		if l.Link != "" && l.Code != "" {
			_, err = ls.s.Model(link).
				Where("link = ? and code = ?", l.Link, l.Code).
				Delete()
		} else {
			_, err = ls.s.Model(link).
				Where("link = ? or code = ?", l.Link, l.Code).
				Delete()
		}
		if err != nil {
			return errors.ModelErr{Err: err}
		}
		return nil
	}
	return errors.InitializationErr{Err: errors.InitializedIncorrectlyErr{}}
}
