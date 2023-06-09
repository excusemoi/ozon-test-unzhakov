package storage

import "ozon-test-unzhakov/internal/model"

type LinkStorage interface {
	Get(link *model.Link) (*model.Link, error)
	CreateLink(link *model.Link) (*model.Link, error)
	UpdateLink(link *model.Link) (*model.Link, error)
	DeleteLink(link *model.Link) error
}
