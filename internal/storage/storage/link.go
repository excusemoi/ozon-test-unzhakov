package storage

import "ozon-test-unzhakov/internal/model"

type LinkStorage interface {
	GetLink(link string) (*model.Link, error)
	CreateLink(link *model.Link) (*model.Link, error)
	UpdateLink(link *model.Link) (*model.Link, error)
	DeleteLink(id string) error
}
