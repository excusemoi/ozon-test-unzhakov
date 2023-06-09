package cache

import (
	"container/list"
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"os"
	"ozon-test-unzhakov/internal/config"
	"ozon-test-unzhakov/internal/model"
	"ozon-test-unzhakov/internal/storage/storage"
	"path/filepath"
	"sync"
)

type linkStorage struct {
	sync.RWMutex
	hashByLink map[string]*list.Element
	hashByCode map[string]*list.Element
	list       list.List
	capacity   int64
}

func NewLinkStorage() (storage.LinkStorage, error) {
	err := config.InitConfig(filepath.Join("..", "..", "..", "config"), os.Getenv("CONFIG_NAME"), "yaml")
	if err != nil {
		return nil, err
	}
	ls := &linkStorage{
		hashByLink: map[string]*list.Element{},
		hashByCode: map[string]*list.Element{},
		list:       list.List{},
		capacity:   viper.GetInt64("cache.capacity"),
	}
	return ls, nil
}

func (ls *linkStorage) Get(l *model.Link) (*model.Link, error) {
	ls.RLock()
	defer ls.RUnlock()
	link, inByLinkHash := ls.hashByLink[l.Link]
	link, inByCodeHash := ls.hashByCode[l.Code]
	if inByLinkHash || inByCodeHash {
		v := ls.list.Remove(link).(*model.Link)
		ls.hashByLink[v.Link] = ls.list.PushBack(v)
		ls.hashByLink[v.Code] = ls.hashByLink[v.Link]
		return v, nil
	}
	return nil, errors.New(fmt.Sprintf("linkStorageCache: can't find link %v", l))
}

func (ls *linkStorage) CreateLink(l *model.Link) (*model.Link, error) {
	ls.Lock()
	defer ls.Unlock()
	if _, in := ls.hashByLink[l.Link]; in {
		//ls.list.Remove(link)
		return nil, errors.New(fmt.Sprintf("linkStorageCache: can't duplicate link %v", l))
	}
	if ls.capacity != 0 {
		if int64(ls.list.Len()) == ls.capacity {
			f := ls.list.Front()
			ls.list.Remove(f)
		}
	} else {
		return nil, errors.New("linkStorageCache: cache has 0 capacity")
	}
	ls.hashByLink[l.Link] = ls.list.PushBack(l)
	ls.hashByCode[l.Code] = ls.hashByLink[l.Link]
	return ls.hashByLink[l.Link].Value.(*model.Link), nil
}

func (ls *linkStorage) UpdateLink(l *model.Link) (*model.Link, error) {
	return ls.CreateLink(l)
}

func (ls *linkStorage) DeleteLink(l *model.Link) error {
	ls.Lock()
	defer ls.Unlock()
	if link, in := ls.hashByLink[l.Link]; in {
		ls.list.Remove(link)
		delete(ls.hashByLink, l.Link)
		delete(ls.hashByCode, l.Code)
		return nil
	}
	return errors.New(fmt.Sprintf("linkStorageCache: can't find link %v", l))
}
