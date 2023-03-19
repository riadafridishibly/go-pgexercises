package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"

	"github.com/riadafridishibly/go-pgexercises/pgexercise"
)

type PgConfig struct {
	Host     string
	Port     uint16
	User     string
	Password string
	DBName   string
}

func (cfg PgConfig) ConnString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName)
}

func main() {
	pgCfg := PgConfig{
		Host:     "127.0.0.1",
		Port:     5432,
		User:     "postgres",
		Password: "password",
		DBName:   "exercises",
	}

	conn, err := sql.Open("postgres", pgCfg.ConnString())
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.TODO()

	var q pgexercise.Queries
	q.Prepare(ctx, conn)

	v, err := q.GetAllMembersWithRecommender(ctx)
	if err != nil {
		log.Fatal(err)
	}

	for _, vv := range v {
		fmt.Printf("Member: %-30q Recommender: %q\n", vv.MemberName, vv.RecommenderName)
	}
}
