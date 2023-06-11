package pg

import (
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
	"log"
	"ozon-test-unzhakov/internal/model"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewLinkStorage(t *testing.T) {
	err := godotenv.Load(filepath.Join("..", "..", "..", ".env"))
	if err != nil {
		t.Error(err)
	}
	_, err = NewLinkStorage(nil) //TODO mock
	if err != nil {
		t.Error(err)
	}
}

func TestLinkStorage_Migrate(t *testing.T) {
	err := godotenv.Load(filepath.Join("..", "..", "..", ".env"))
	if err != nil {
		t.Error(err)
	}
	s, err := NewLinkStorage(nil) //TODO mock
	if err != nil {
		t.Fatal(err)
	}
	sn := s.(*linkStorage)
	err = sn.Migrate()
	if err != nil {
		t.Log(err)
	}
}

func TestLinkStorage_Create(t *testing.T) {
	err := godotenv.Load(filepath.Join("..", "..", "..", ".env"))
	if err != nil {
		t.Error(err)
	}
	s, err := NewLinkStorage(nil) //TODO mock
	if err != nil {
		t.Fatal(err)
	}
	err = s.(*linkStorage).Migrate()
	if err != nil {
		log.Fatal(err)
	}
	cases := []struct {
		input struct {
			link *model.Link
		}
		expectedErr error
		name        string
	}{
		{
			input: struct {
				link *model.Link
			}{
				link: &model.Link{
					Link: "https://www.google.com/maps",
					Code: "link_code_1",
				}},
			name: "valid_case_1",
		},
		{
			input: struct {
				link *model.Link
			}{
				link: &model.Link{
					Link: "https://github.com/grpc-ecosystem/grpc-gateway#readme",
					Code: "link_code_1",
				}},
			name: "valid_case_2",
		},
	}
	for i := range cases {
		t.Run(cases[i].name, func(t1 *testing.T) {
			result, err := s.CreateLink(cases[i].input.link)
			require.ErrorIs(t1, err, cases[i].expectedErr, result)
			err = s.DeleteLink(cases[i].input.link)
			require.ErrorIs(t1, err, cases[i].expectedErr, "")
		})
	}
}

func TestLinkStorage_Delete(t *testing.T) {
	err := godotenv.Load(filepath.Join("..", "..", "..", ".env"))
	if err != nil {
		t.Error(err)
	}
	s, err := NewLinkStorage(nil) //TODO mock
	if err != nil {
		t.Fatal(err)
	}
	err = s.(*linkStorage).Migrate()
	if err != nil {
		log.Fatal(err)
	}
	cases := []struct {
		input struct{ link *model.Link }
	}{
		{input: struct{ link *model.Link }{link: &model.Link{
			Link: "https://www.google.com/maps",
			Code: "link_code_1",
		}}},
		{input: struct{ link *model.Link }{link: &model.Link{
			Link: "https://github.com/grpc-ecosystem/grpc-gateway#readme",
			Code: "link_code_1",
		}}},
	}
	for i := range cases {
		_, err := s.CreateLink(cases[i].input.link)
		if err != nil {
			t.Log(err)
		}
		err = s.DeleteLink(cases[i].input.link)
		if err != nil {
			t.Log(err)
		}
	}
}

func TestLinkStorage_Get(t *testing.T) {
	err := godotenv.Load(filepath.Join("..", "..", "..", ".env"))
	if err != nil {
		t.Error(err)
	}
	s, err := NewLinkStorage(nil) //TODO mock
	if err != nil {
		t.Fatal(err)
	}
	err = s.(*linkStorage).Migrate()
	if err != nil {
		log.Fatal(err)
	}
	cases := []struct {
		input       struct{ link *model.Link }
		expectedErr error
		name        string
	}{
		{input: struct{ link *model.Link }{link: &model.Link{
			Link: "https://www.google.com/maps",
			Code: "link_code_1",
		}}},
		{input: struct{ link *model.Link }{link: &model.Link{
			Link: "https://github.com/grpc-ecosystem/grpc-gateway#readme",
			Code: "link_code_1",
		}}},
	}
	for i := range cases {
		result, err := s.CreateLink(cases[i].input.link)
		if err != nil {
			t.Log(err)
		}
		result, err = s.Get(cases[i].input.link)
		if err != nil {
			t.Log(err)
		}
		t.Log(*result)
		err = s.DeleteLink(cases[i].input.link)
		if err != nil {
			t.Log(err)
		}
	}
}

func TestLinkStorage_Update(t *testing.T) {
	err := godotenv.Load(filepath.Join("..", "..", "..", ".env"))
	if err != nil {
		t.Error(err)
	}
	s, err := NewLinkStorage(nil) //TODO mock
	if err != nil {
		t.Fatal(err)
	}
	err = s.(*linkStorage).Migrate()
	if err != nil {
		t.Fatal(err)
	}
	cases := []struct {
		input struct{ link *model.Link }
	}{
		{input: struct{ link *model.Link }{link: &model.Link{
			Link: "https://www.google.com/maps",
			Code: "link_code_1",
		}}},
		{input: struct{ link *model.Link }{link: &model.Link{
			Link: "https://github.com/grpc-ecosystem/grpc-gateway#readme",
			Code: "link_code_1",
		}}},
	}
	for i := range cases {
		result, err := s.CreateLink(cases[i].input.link)
		if err != nil {
			t.Log(err)
		}
		cases[i].input.link.Code += "updated"
		result, err = s.UpdateLink(cases[i].input.link)
		if err != nil {
			t.Log(err)
		}
		result, err = s.Get(cases[i].input.link)
		if !strings.Contains(result.Code, "updated") {
			t.Errorf("link with %d was not updated", cases[i].input.link.Id)
		}
		t.Log(*result)
		err = s.DeleteLink(cases[i].input.link)
		if err != nil {
			t.Log(err)
		}
	}
}
