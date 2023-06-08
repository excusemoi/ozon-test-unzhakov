package service

import (
	"github.com/spf13/viper"
	"math/rand"
	"ozon-test-unzhakov/internal/dto"
	"ozon-test-unzhakov/internal/model"
	"ozon-test-unzhakov/internal/storage/storage"
	"path/filepath"
)

const bytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"

type LinkService interface {
	GetInitialLink(code string) (*dto.Link, error)
	CreateShortLink(link string) (*dto.Link, error)
	Create(link *dto.Link) (*dto.Link, error)
	Update(link *dto.Link) (*dto.Link, error)
	Delete(id string) error
}

type linkService struct {
	s           storage.LinkStorage
	linksLength int
}

func NewLinkService(s storage.LinkStorage) (LinkService, error) {
	viper.AddConfigPath(filepath.Join("..", "..", "config"))
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	return &linkService{s: s, linksLength: viper.GetInt("linksLength")}, nil
}

func (ls *linkService) GetInitialLink(c string) (*dto.Link, error) {
	link, err := ls.s.GetLink(c)
	if err != nil {
		return nil, err
	}
	return &dto.Link{
		Id:   link.Id,
		Link: link.Link,
		Code: link.Code,
	}, nil
}

func (ls *linkService) CreateShortLink(l string) (*dto.Link, error) {
	link, err := ls.s.CreateLink(&model.Link{
		Link: l,
		Code: ls.Short(),
	})
	if err != nil {
		return nil, err
	}
	return &dto.Link{
		Id:   link.Id,
		Link: link.Link,
		Code: link.Code,
	}, nil
}

func (ls *linkService) Create(l *dto.Link) (*dto.Link, error) {
	link, err := ls.s.CreateLink(&model.Link{
		Id:   l.Id,
		Link: l.Link,
		Code: l.Code,
	})
	if err != nil {
		return nil, err
	}
	return &dto.Link{
		Id:   link.Id,
		Link: link.Link,
		Code: link.Code,
	}, nil
}

func (ls *linkService) Update(l *dto.Link) (*dto.Link, error) {
	link, err := ls.s.UpdateLink(&model.Link{
		Id:   l.Id,
		Link: l.Link,
		Code: l.Code,
	})
	if err != nil {
		return nil, err
	}
	return &dto.Link{
		Id:   link.Id,
		Link: link.Link,
		Code: link.Code,
	}, nil
}

func (ls *linkService) Delete(l string) error {
	err := ls.s.DeleteLink(l)
	if err != nil {
		return err
	}
	return nil
}

func (ls *linkService) Short() string {
	sb := []byte(bytes)
	buff := make([]byte, ls.linksLength)
	rand.Shuffle(len(sb), func(i, j int) {
		sb[i], sb[j] = sb[j], sb[i]
	})
	for i := range buff {
		buff[i] = sb[i]
	}
	return string(buff)
}
