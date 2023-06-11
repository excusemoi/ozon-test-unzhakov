package service

import (
	"math/rand"
	"net/url"
	"os"
	"ozon-test-unzhakov/internal/config"
	"ozon-test-unzhakov/internal/dto"
	"ozon-test-unzhakov/internal/model"
	"ozon-test-unzhakov/internal/storage/storage"
	"path/filepath"
)

const bytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"

type LinkService interface {
	GetInitialLink(link *dto.Link) (*dto.Link, error)
	CreateShortLink(link *dto.Link) (*dto.Link, error)
	Create(link *dto.Link) (*dto.Link, error)
	Update(link *dto.Link) (*dto.Link, error)
	Delete(link *dto.Link) error
}

type linkService struct {
	s           storage.LinkStorage
	linksLength int
}

func NewLinkService(s storage.LinkStorage, linksLength int) (LinkService, error) {
	err := config.InitConfig(filepath.Join("..", "..", "config"),
		os.Getenv("CONFIG_NAME"),
		"yaml")
	if err != nil {
		return nil, err
	}
	return &linkService{s: s, linksLength: linksLength}, nil
}

func (ls *linkService) GetInitialLink(l *dto.Link) (*dto.Link, error) {
	link, err := ls.s.Get(&model.Link{Code: l.Code, Link: l.Link})
	if err != nil {
		return nil, err
	}
	_, err = url.ParseRequestURI(link.Link)
	if err != nil {
		return nil, err
	}
	return &dto.Link{
		Id:   link.Id,
		Link: link.Link,
		Code: link.Code,
	}, nil
}

func (ls *linkService) CreateShortLink(l *dto.Link) (*dto.Link, error) {
	link, err := ls.s.CreateLink(&model.Link{
		Link: l.Link,
		Code: ls.Short(),
	})
	if err != nil {
		return nil, err
	}
	_, err = url.ParseRequestURI(link.Link)
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
	_, err = url.ParseRequestURI(link.Link)
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
	_, err = url.ParseRequestURI(link.Link)
	if err != nil {
		return nil, err
	}
	return &dto.Link{
		Id:   link.Id,
		Link: link.Link,
		Code: link.Code,
	}, nil
}

func (ls *linkService) Delete(l *dto.Link) error {
	err := ls.s.DeleteLink(&model.Link{
		Id:   l.Id,
		Link: l.Link,
		Code: l.Code,
	})
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
