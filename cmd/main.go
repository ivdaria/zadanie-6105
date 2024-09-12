package main

import (
	"context"
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
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

//POSTGRES_CONN=postgres://postgres:postgres@localhost:5432/zadanie

type Config struct {
	DSN string `env:"POSTGRES_CONN"`
}

func main() {
	ctx := context.Background()
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		panic(err)
	}

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

	e := echo.New()
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

	e.Logger.Fatal(e.Start("0.0.0.0:8080"))

}
