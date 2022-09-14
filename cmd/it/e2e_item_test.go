package it_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
	"github.com/mercadolibre/fury_go-platform/pkg/fury"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/ory/dockertest/v3"
	"github.com/osalomon89/test-crud-api/internal/core/domain"
	"github.com/osalomon89/test-crud-api/internal/core/ports"
	mysqlrepository "github.com/osalomon89/test-crud-api/internal/infrastructure/repositories/mysql"
	"github.com/osalomon89/test-crud-api/internal/infrastructure/server"
	"github.com/osalomon89/test-crud-api/internal/infrastructure/server/handler/dto"
)

type e2eItemTestSuite struct {
	dbConnectionStr string
	port            int
	dbConn          *sqlx.DB
	repository      ports.ItemRepository
	dbMigration     *migrate.Migrate
}

var _ = Describe("Creating items in the marketplace", func() {
	var dbConn *sqlx.DB
	suite := new(e2eItemTestSuite)

	pool, err := dockertest.NewPool("")
	Expect(err).NotTo(HaveOccurred())

	resource, err := pool.Run("mysql", "5.7", []string{"MYSQL_ROOT_PASSWORD=secret"})
	Expect(err).NotTo(HaveOccurred())

	hostAndPort := resource.GetPort("3306/tcp")
	suite.dbConnectionStr = fmt.Sprintf("root:secret@(localhost:%s)/mysql?charset=utf8&parseTime=true", hostAndPort)
	resource.Expire(120) // Tell docker to hard kill the container in 120 seconds

	pool.MaxWait = 120 * time.Second

	err = pool.Retry(func() error {
		dbConn, err = sqlx.Connect("mysql", suite.dbConnectionStr)
		if err != nil {
			return err
		}

		suite.dbConn = dbConn
		return dbConn.Ping()
	})

	Expect(err).NotTo(HaveOccurred())

	suite.port = 8080
	driver, err := mysql.WithInstance(suite.dbConn.DB, &mysql.Config{})
	Expect(err).NotTo(HaveOccurred())

	suite.dbMigration, err = migrate.NewWithDatabaseInstance(
		"file://../../db/migration",
		"mysql",
		driver,
	)
	Expect(err).NotTo(HaveOccurred())

	suite.repository, err = mysqlrepository.NewItemRepository(suite.dbConn)
	Expect(err).NotTo(HaveOccurred())

	app, err := fury.NewWebApplication()
	Expect(err).NotTo(HaveOccurred())

	serverReady := make(chan bool)
	furyHandler, err := server.NewHTTPServer(app, suite.dbConn, serverReady)
	Expect(err).NotTo(HaveOccurred())

	furyHandler.SetupRouter()

	go furyHandler.Run()
	<-serverReady

	BeforeEach(func() {
		err := suite.dbMigration.Up()
		if err != nil && err != migrate.ErrNoChange {
			Expect(err).NotTo(HaveOccurred())
		}
	})

	AfterEach(func() {
		Expect(suite.dbMigration.Down()).NotTo(HaveOccurred())
	})

	When("the user send all parameters", func() {
		Context("and all parameters are ok", func() {
			It("is saved correctly in the database", func() {
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

				req, err := http.NewRequest(echo.POST, fmt.Sprintf("http://localhost:%d/v1/items", suite.port), strings.NewReader(reqStr))
				Expect(err).NotTo(HaveOccurred())

				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

				client := http.Client{}
				response, err := client.Do(req)
				Expect(err).NotTo(HaveOccurred())
				defer response.Body.Close()

				Expect(response.StatusCode).To(Equal(http.StatusCreated))

				var responseBody dto.Response
				err = json.NewDecoder(response.Body).Decode(&responseBody)
				Expect(err).NotTo(HaveOccurred())

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

				Expect(responseBody).To(Equal(want))
			})
		})
	})

	When("the user does not send all parameters", func() {

	})
})
