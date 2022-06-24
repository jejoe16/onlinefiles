package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/cors"
	"github.com/sendgrid/sendgrid-go"
	"github.com/wes-kay/golang_asset_engine/core"
	"github.com/wes-kay/golang_asset_engine/model"
)

type App struct {
	Router   *mux.Router
	Model    *model.Model
	Core     *core.Core
	Sessions map[uuid.UUID]Session
	AWS      *session.Session
}

func Initialize() *App {
	a := &App{}
	a.Router = mux.NewRouter().StrictSlash(true)
	a.InitializeRoutes(a.Router)
	config, err := DecodeFile("config.yaml")
	if err != nil {
		fmt.Println("Error opening config connection")
		panic(err)
	}

	pool, err := NewDatabase(config.Database.UserName, config.Database.Password, config.Database.Name)

	if err != nil {
		fmt.Println("Error opening database connection")
		panic(err)
	}

	a.Model = &model.Model{
		Pool: pool,
	}

	a.Core = &core.Core{
		SendgridClient: sendgrid.NewSendClient(config.Sendgrid.ApiKey),
	}

	a.Sessions = map[uuid.UUID]Session{}

	a.AWS, err = core.NewAWSSession(config.AWS.AWSAccessKey, config.AWS.AWSSecretAccessKey, config.AWS.AWSRegion)
	if err != nil {
		fmt.Println("Error opening AWS connection")
		panic(err)
	}

	return a
}

func NewDatabase(user, password, dbname string) (*pgxpool.Pool, error) {
	connectionString := fmt.Sprintf("postgres://%s:%s@localhost:5432/%s", user, password, dbname)
	pool, err := pgxpool.Connect(context.TODO(), connectionString)
	if err != nil {
		return nil, err
	}
	return pool, err
}

func (a *App) Run() {
	c := cors.New(cors.Options{

		AllowedOrigins: []string{"https://www.idsure.io", "https://api.idsure.io", "http://localhost:8080", "http://localhost:9990"},
		AllowedMethods: []string{"GET", "PUT", "POST", "DELETE", "UPDATE", "OPTIONS", "post", "get"},
		MaxAge:         3600,
		// ExposedHeaders:   []string{"Authorization"},
		AllowedHeaders:   []string{"Access-Control-Allow-Origin", "Content-Type", "X-Auth-Token", "Auth", "Authorization"},
		AllowCredentials: true,
		Debug:            false,
	})

	log.Fatal(http.ListenAndServe(":9990", c.Handler(a.Router)))
}
