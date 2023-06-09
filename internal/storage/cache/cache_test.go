package cache

import (
	"github.com/joho/godotenv"
	"ozon-test-unzhakov/internal/model"
	"path/filepath"
	"testing"
)

func TestNewLinkStorage(t *testing.T) {
	err := godotenv.Load(filepath.Join("..", "..", "..", ".env"))
	if err != nil {
		t.Error(err)
	}
	_, err = NewLinkStorage()
	if err != nil {
		t.Error(err)
	}
}

func TestLinkStorage_Create(t *testing.T) {
	err := godotenv.Load(filepath.Join("..", "..", "..", ".env"))
	if err != nil {
		t.Error(err)
	}
	s, err := NewLinkStorage()
	if err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		input struct {
			link *model.Link
		}
	}{
		{
			input: struct {
				link *model.Link
			}{
				link: &model.Link{Id: 1, Link: "link_1"},
			},
		},
		{
			input: struct {
				link *model.Link
			}{
				link: &model.Link{Id: 2, Link: "link_2"},
			},
		},
	}

	for i := range cases {
		if _, err := s.CreateLink(cases[i].input.link); err != nil {
			t.Log(err)
		}
	}
}

func TestLinkStorage_Get(t *testing.T) {
	err := godotenv.Load(filepath.Join("..", "..", "..", ".env"))
	if err != nil {
		t.Error(err)
	}

	s, _ := NewLinkStorage()

	cases := []struct {
		input struct {
			link *model.Link
		}
	}{
		{
			input: struct {
				link *model.Link
			}{
				link: &model.Link{Id: 1, Link: "link_1"},
			},
		},
		{
			input: struct {
				link *model.Link
			}{
				link: &model.Link{Id: 2, Link: "link_2"},
			},
		},
	}

	for i := range cases {
		s.CreateLink(cases[i].input.link)
		if _, err := s.Get(cases[i].input.link); err != nil {
			t.Error(err)
		}
	}
}

func TestLinkStorage_Delete(t *testing.T) {
	err := godotenv.Load(filepath.Join("..", "..", "..", ".env"))
	if err != nil {
		t.Error(err)
	}

	s, _ := NewLinkStorage()

	cases := []struct {
		input struct {
			link *model.Link
		}
	}{
		{
			input: struct {
				link *model.Link
			}{
				link: &model.Link{Id: 1, Link: "link_1"},
			},
		},
		{
			input: struct {
				link *model.Link
			}{
				link: &model.Link{Id: 2, Link: "link_2"},
			},
		},
	}

	for i := range cases {
		s.CreateLink(cases[i].input.link)
		if err := s.DeleteLink(cases[i].input.link); err != nil {
			t.Error(err)
		}
	}
}
