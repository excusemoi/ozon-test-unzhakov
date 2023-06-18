package pg

import (
	"context"
	orm "github.com/go-pg/pg/v10"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log"
	"os"
	"ozon-test-unzhakov/internal/config"
	"ozon-test-unzhakov/internal/model"
	linkStorageErrors "ozon-test-unzhakov/internal/storage/pg/errors"
	"path/filepath"
	"strings"
	"testing"
)

var (
	s *orm.DB
)

func init() { //TODO better way
	err := godotenv.Load(filepath.Join("..", "..", "..", ".env"))
	if err != nil {
		log.Fatal(err)
	}
	err = config.InitConfig(filepath.Join("..", "..", "..", "config"), os.Getenv("CONFIG_NAME"), "yaml")
	if err != nil {
		log.Fatal(err)
	}
	connectOptions, err := orm.ParseURL("postgres://" +
		os.Getenv("POSTGRES_USER") + ":" +
		os.Getenv("POSTGRES_PASSWORD") + "@" +
		viper.GetString("postgres.host") + ":" +
		viper.GetString("postgres.port") + "/" +
		viper.GetString("postgres.name") + "?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	s = orm.Connect(connectOptions)
	err = s.Ping(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}

func TestNewLinkStorage(t *testing.T) {
	cases := []struct {
		input struct {
			pg *orm.DB
		}
		name        string
		expectedErr error
	}{
		{
			input: struct {
				pg *orm.DB
			}{
				pg: s,
			},
			name: "valid_case_1",
		},
		{
			input: struct {
				pg *orm.DB
			}{
				pg: nil,
			},
			name:        "nil_input_storage_case",
			expectedErr: linkStorageErrors.InitializationErr{Err: linkStorageErrors.IncorrectInitializationErr{}},
		},
	}

	for i := range cases {
		_, err := NewLinkStorage(cases[i].input.pg)
		assert.ErrorIs(t, err, cases[i].expectedErr)
	}
}

func TestLinkStorage_Migrate(t *testing.T) {
	cases := []struct {
		input struct {
			s *linkStorage
		}
		name        string
		expectedErr error
	}{
		{
			input: struct {
				s *linkStorage
			}{
				s: &linkStorage{s: s}},
			name: "valid_case_1",
		},
		{
			input: struct {
				s *linkStorage
			}{
				s: &linkStorage{}},
			name:        "empty_storage_migration_case",
			expectedErr: linkStorageErrors.InitializationErr{Err: linkStorageErrors.InitializedIncorrectlyErr{}},
		},
	}

	for i := range cases {
		err := cases[i].input.s.Migrate()
		assert.ErrorIs(t, err, cases[i].expectedErr)
	}

}

func TestLinkStorage_Create(t *testing.T) {

	cases := []struct {
		input struct {
			s         *linkStorage
			link      *model.Link
			duplicate bool
		}
		expectedErr error
		name        string
	}{
		{
			input: struct {
				s         *linkStorage
				link      *model.Link
				duplicate bool
			}{
				s: &linkStorage{s: s},
				link: &model.Link{
					Link: "https://www.google.com/maps",
					Code: "link_code_1",
				}},
			name: "valid_case_1",
		},
		{
			input: struct {
				s         *linkStorage
				link      *model.Link
				duplicate bool
			}{
				s: &linkStorage{s: s},
				link: &model.Link{
					Link: "https://github.com/grpc-ecosystem/grpc-gateway#readme",
					Code: "link_code_1",
				}},
			name: "valid_case_2",
		},
		{
			input: struct {
				s         *linkStorage
				link      *model.Link
				duplicate bool
			}{
				s: &linkStorage{},
				link: &model.Link{
					Link: "https://github.com/grpc-ecosystem/grpc-gateway#readme",
					Code: "link_code_1",
				}},
			name:        "empty_storage_creation_case",
			expectedErr: linkStorageErrors.InitializationErr{Err: linkStorageErrors.InitializedIncorrectlyErr{}},
		},
		{
			input: struct {
				s         *linkStorage
				link      *model.Link
				duplicate bool
			}{
				s:    &linkStorage{s: s},
				link: nil},
			name:        "nil_link_creation_case",
			expectedErr: linkStorageErrors.ModelErr{Err: linkStorageErrors.ModelNilErr{}},
		},
	}

	for i := range cases {
		t.Run(cases[i].name, func(t1 *testing.T) {
			_, err := cases[i].input.s.CreateLink(cases[i].input.link)
			if cases[i].expectedErr == nil {
				assert.ErrorIs(t1, err, cases[i].expectedErr)
			} else {
				assert.ErrorAs(t1, err, &cases[i].expectedErr)
			}
			if err == nil {
				err = cases[i].input.s.DeleteLink(cases[i].input.link)
				assert.ErrorIs(t1, err, cases[i].expectedErr)
			}
		})
	}
}

func TestLinkStorage_Delete(t *testing.T) {
	s, err := NewLinkStorage(s)
	require.NoError(t, err)

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
		assert.NoError(t, err)

		err = s.DeleteLink(cases[i].input.link)
		assert.NoError(t, err)
	}
}

func TestLinkStorage_Get(t *testing.T) {
	s, err := NewLinkStorage(s)
	require.NoError(t, err)

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
		_, err := s.CreateLink(cases[i].input.link)
		assert.NoError(t, err)

		result, err := s.Get(cases[i].input.link)
		assert.NoError(t, err)
		t.Log(*result)

		err = s.DeleteLink(cases[i].input.link)
		assert.NoError(t, err)
	}
}

func TestLinkStorage_Update(t *testing.T) {
	s, err := NewLinkStorage(s)
	require.NoError(t, err)

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
		assert.NoError(t, err)
		cases[i].input.link.Code += "updated"

		result, err := s.UpdateLink(cases[i].input.link)
		assert.NoError(t, err)
		result, err = s.Get(cases[i].input.link)
		if !strings.Contains(result.Code, "updated") {
			t.Errorf("link with %d was not updated", cases[i].input.link.Id)
		}
		t.Log(*result)

		err = s.DeleteLink(cases[i].input.link)
		assert.NoError(t, err)
	}
}
