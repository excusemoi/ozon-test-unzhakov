package cache

import (
	"container/list"
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"ozon-test-unzhakov/internal/config"
	"ozon-test-unzhakov/internal/model"
	"path/filepath"
	"sync"
)

type linkStorage struct {
	sync.RWMutex
	hash     map[string]*list.Element
	list     list.List
	capacity int64
}

func NewLinkStorage() (interfaces.LinkStorage, error) {
	err := config.InitConfig(filepath.Join("..", "..", "..", "config"), "config", "yaml")
	if err != nil {
		return nil, err
	}
	ls := &linkStorage{
		hash:     map[string]*list.Element{},
		list:     list.List{},
		capacity: viper.GetInt64("cache.capacity"),
	}
	return ls, nil
}

func (ls *linkStorage) GetLink(l string) (*model.Link, error) {
	ls.RLock()
	defer ls.RUnlock()
	if link, in := ls.hash[l]; in {
		v := ls.list.Remove(link).(*model.Link)
		ls.hash[v.Link] = ls.list.PushBack(v)
		return v, nil
	}
	return nil, errors.New(fmt.Sprintf("linkStorageCache: can't find link %s", l))
}

func (ls *linkStorage) CreateLink(l *model.Link) (*model.Link, error) {
	ls.Lock()
	defer ls.Unlock()
	if link, in := ls.hash[l.Link]; in {
		ls.list.Remove(link)
	}
	if ls.capacity != 0 {
		if int64(ls.list.Len()) == ls.capacity {
			f := ls.list.Front()
			ls.list.Remove(f)
		}
	} else {
		return nil, errors.New("linkStorageCache: cache has 0 capacity")
	}
	ls.hash[l.Link] = ls.list.PushBack(l)
	return ls.hash[l.Link].Value.(*model.Link), nil
}

func (ls *linkStorage) UpdateLink(l *model.Link) (*model.Link, error) {
	return ls.CreateLink(l)
}

func (ls *linkStorage) DeleteLink(l string) error {
	ls.Lock()
	defer ls.Unlock()
	if link, in := ls.hash[l]; in {
		ls.list.Remove(link)
		delete(ls.hash, l)
		return nil
	}
	return errors.New(fmt.Sprintf("linkStorageCache: can't find link %s", l))
}
