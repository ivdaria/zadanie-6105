package main

import (
	"context"
	"database/sql"
	"fmt"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725731644-team-78845/zadanie-6105/internal/gateway"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725731644-team-78845/zadanie-6105/internal/repository/bids"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725731644-team-78845/zadanie-6105/internal/repository/employee"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725731644-team-78845/zadanie-6105/internal/repository/feedback"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725731644-team-78845/zadanie-6105/internal/repository/organization"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725731644-team-78845/zadanie-6105/internal/repository/organization_responsible"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725731644-team-78845/zadanie-6105/internal/repository/tenders"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725731644-team-78845/zadanie-6105/internal/usecase/check_can_edit_bid"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725731644-team-78845/zadanie-6105/pkg/api"
	"github.com/caarlos0/env/v11"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	echomiddleware "github.com/oapi-codegen/echo-middleware"
	"github.com/pressly/goose/v3"
)

//POSTGRES_CONN=postgres://postgres:postgres@localhost:5432/zadanie

type Config struct {
	DSN           string `env:"POSTGRES_CONN"`
	ServerAddress string `env:"SERVER_ADDRESS" envDefault:"0.0.0.0:8080"`
}

func main() {
	ctx := context.Background()
	var cfg Config

	if err := env.Parse(&cfg); err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", cfg)

	upMigrations(ctx, cfg)

	conn, err := pgx.Connect(ctx, cfg.DSN)
	if err != nil {
		panic(err)
	}
	defer conn.Close(ctx)

	employeeRepo := employee.NewRepo(conn)
	tendersRepo := tenders.NewRepo(conn)
	bidsRepo := bids.NewRepo(conn)
	organizationRepo := organization.NewRepo(conn)
	organizationResponsibleRepo := organization_responsible.NewRepo(conn)
	feedbackRepo := feedback.NewRepo(conn)
	userCanEditBidCheckerUseCase := check_can_edit_bid.NewUseCase(organizationResponsibleRepo)

	swagger, err := api.GetSwagger()
	if err != nil {
		panic(fmt.Errorf("get swagger: %w", err))
	}
	//Обход проверки валидации сервера в сваггере
	swagger.Servers = openapi3.Servers{
		{
			URL: "/api",
		},
	}

	e := echo.New()
	e.Use(middleware.Logger())
	validator := echomiddleware.OapiRequestValidator(swagger)

	e.Pre(validator)
	handlers := gateway.NewServer(
		tendersRepo,
		employeeRepo,
		organizationRepo,
		bidsRepo,
		organizationResponsibleRepo,
		feedbackRepo,
		userCanEditBidCheckerUseCase,
	)

	api.RegisterHandlersWithBaseURL(e, handlers, "/api")

	e.Logger.Fatal(e.Start(cfg.ServerAddress))
}

func upMigrations(ctx context.Context, cfg Config) {
	db, err := sql.Open("postgres", cfg.DSN)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		panic(err)
	}

	if err := goose.RunContext(ctx, "up", db, "./migrations"); err != nil {
		panic(err)
	}
}
