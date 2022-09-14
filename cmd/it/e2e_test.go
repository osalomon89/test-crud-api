package it_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
	"github.com/mercadolibre/fury_go-platform/pkg/fury"
	"github.com/ory/dockertest/v3"
	"github.com/osalomon89/test-crud-api/internal/core/domain"
	"github.com/osalomon89/test-crud-api/internal/core/ports"
	mysqlrepository "github.com/osalomon89/test-crud-api/internal/infrastructure/repositories/mysql"
	"github.com/osalomon89/test-crud-api/internal/infrastructure/server"
	"github.com/osalomon89/test-crud-api/internal/infrastructure/server/handler/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type e2eTestSuite struct {
	suite.Suite
	dbConnectionStr string
	port            int
	dbConn          *sqlx.DB
	repository      ports.ItemRepository
	dbMigration     *migrate.Migrate
}

// function will trigger the test, and the name of the function must start with Test (prefix)
func TestE2ETestSuite(t *testing.T) {
	suite.Run(t, &e2eTestSuite{})
}

// Run before starting all tests: will run first to Initialize the Suite data needed.
func (s *e2eTestSuite) SetupSuite() {
	var dbConn *sqlx.DB
	pool, err := dockertest.NewPool("")
	s.Require().NoError(err)

	resource, err := pool.Run("mysql", "5.7", []string{"MYSQL_ROOT_PASSWORD=secret"})
	s.Require().NoError(err)

	hostAndPort := resource.GetPort("3306/tcp")
	s.dbConnectionStr = fmt.Sprintf("root:secret@(localhost:%s)/mysql?charset=utf8&parseTime=true", hostAndPort)
	resource.Expire(120) // Tell docker to hard kill the container in 120 seconds

	pool.MaxWait = 120 * time.Second

	s.Require().NoError(pool.Retry(func() error {
		dbConn, err = sqlx.Connect("mysql", s.dbConnectionStr)
		if err != nil {
			return err
		}

		s.dbConn = dbConn
		return dbConn.Ping()
	}))

	s.port = 8080
	driver, err := mysql.WithInstance(s.dbConn.DB, &mysql.Config{})
	s.Require().NoError(err)

	s.dbMigration, err = migrate.NewWithDatabaseInstance(
		"file://../../db/migration",
		"mysql",
		driver,
	)
	s.Require().NoError(err)

	s.repository, err = mysqlrepository.NewItemRepository(s.dbConn)
	s.Require().NoError(err)

	app, err := fury.NewWebApplication()
	s.Require().NoError(err)

	serverReady := make(chan bool)
	furyHandler, err := server.NewHTTPServer(app, s.dbConn, serverReady)
	s.Require().NoError(err)

	furyHandler.SetupRouter()

	go furyHandler.Run()
	<-serverReady
}

// Run After All Test Done: will run the last, after all, tests are done.
func (s *e2eTestSuite) TearDownSuite() {
	s.Require().NoError(s.dbConn.Close())
}

// Run Before a Test: will run before a test.
func (s *e2eTestSuite) SetupTest() {
	if err := s.dbMigration.Up(); err != nil && err != migrate.ErrNoChange {
		s.Require().NoError(err)
	}
}

// Run After a Test: will run after a test.
func (s *e2eTestSuite) TearDownTest() {
	s.NoError(s.dbMigration.Down())
}

func (s *e2eTestSuite) Test_EndToEnd_CreateItem() {
	reqStr := `{"code":"sa4123", 
	"title": "my-title", 
	"description":"my-description", 
	"price":50, "stock":150, 
	"itemType":"SELLER", 
	"leader":true, 
	"leaderLevel":"PLATINUM",
	"photos": [
		"https://http2.mlstatic.com/D_NQ_NP_729539-MLA48049063325_102021-O.jpg",
		"https://http2.mlstatic.com/D_NQ_NP_879745-MLA48049070326_102021-O.jpg"
	]
	}`

	req, err := http.NewRequest(echo.POST, fmt.Sprintf("http://localhost:%d/v1/items", s.port), strings.NewReader(reqStr))
	assert.NoError(s.T(), err)

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	client := http.Client{}
	response, err := client.Do(req)
	assert.NoError(s.T(), err)
	defer response.Body.Close()

	assert.Equal(s.T(), http.StatusCreated, response.StatusCode)

	var responseBody dto.Response
	err = json.NewDecoder(response.Body).Decode(&responseBody)
	assert.NoError(s.T(), err)

	want := dto.Response{
		Status:  http.StatusCreated,
		Message: "Success",
		Data: &dto.ItemResponse{
			ID:          1,
			Code:        "sa4123",
			Title:       "my-title",
			Description: "my-description",
			Price:       50,
			Stock:       150,
			ItemType:    "SELLER",
			Leader:      true,
			LeaderLevel: domain.LeaderLevelPlatinum,
			Status:      domain.StatusActive,
			Photos: []string{
				"https://http2.mlstatic.com/D_NQ_NP_729539-MLA48049063325_102021-O.jpg",
				"https://http2.mlstatic.com/D_NQ_NP_879745-MLA48049070326_102021-O.jpg",
			},
			CreatedAt: responseBody.Data.CreatedAt,
			UpdatedAt: responseBody.Data.UpdatedAt,
		},
	}

	assert.Equal(s.T(), want, responseBody)
}

func (s *e2eTestSuite) Test_EndToEnd_GetItemByID() {
	item := domain.Item{
		Code:        "sa4123",
		Title:       "my-title",
		Description: "my-description",
		Price:       50,
		Stock:       15,
		ItemType:    "SELLER",
		Leader:      true,
		LeaderLevel: domain.LeaderLevelPlatinum,
		Status:      domain.StatusActive,
	}

	assert.NoError(s.T(), s.repository.SaveItem(context.TODO(), &item))

	req, err := http.NewRequest(echo.GET, fmt.Sprintf("http://localhost:%d/v1/items/1", s.port), nil)
	assert.NoError(s.T(), err)

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	client := http.Client{}
	response, err := client.Do(req)
	assert.NoError(s.T(), err)
	defer response.Body.Close()

	assert.Equal(s.T(), http.StatusOK, response.StatusCode)

	var responseBody dto.Response
	err = json.NewDecoder(response.Body).Decode(&responseBody)
	assert.NoError(s.T(), err)

	want := dto.Response{
		Status:  http.StatusOK,
		Message: "Success",
		Data: &dto.ItemResponse{
			ID:          1,
			Code:        "sa4123",
			Title:       "my-title",
			Description: "my-description",
			Price:       50,
			Stock:       15,
			ItemType:    "SELLER",
			Leader:      true,
			LeaderLevel: domain.LeaderLevelPlatinum,
			Status:      domain.StatusActive,
			CreatedAt:   responseBody.Data.CreatedAt,
			UpdatedAt:   responseBody.Data.UpdatedAt,
		},
	}

	assert.Equal(s.T(), want, responseBody)
}
