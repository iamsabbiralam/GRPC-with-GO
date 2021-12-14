package main

import (
	"log"
	"net"
	"strconv"
	"strings"

	coreCategory "grpc-blog/blog/core/category"
	"grpc-blog/blog/services/category"
	"grpc-blog/blog/storage/postgres"
	protoCategory "grpc-blog/gunk/v1/category"

	"github.com/spf13/viper"

	"google.golang.org/grpc"
)

func main () {
	lis, err := net.Listen("tcp", ":3000")

	if err != nil {
		log.Fatalf("Failed to run server: %s", err)
	}

	config := viper.NewWithOptions(
		viper.EnvKeyReplacer(
			strings.NewReplacer(".", "_"),
		),
	)
	config.SetConfigFile("env/config")
	config.SetConfigType("ini")
	config.AutomaticEnv()
	if err := config.ReadInConfig(); err != nil {
		log.Printf("error loading configuration: %v", err)
	}

	grpcServer := grpc.NewServer()
	store, err := newDBFromConfig(config)
	if err != nil {
		log.Fatalf("failed to connect to database: %s", err)
	}
	cs := coreCategory.NewCoreSvc(store)
	s := category.NewCategoryServer(cs)

	protoCategory.RegisterCategoryServiceServer(grpcServer, s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed too serve: %s", err)
	}
}

func newDBFromConfig(config *viper.Viper) (*postgres.Storage, error) {
	cf := func(c string) string { return config.GetString("database." + c) }
	ci := func(c string) string { return strconv.Itoa(config.GetInt("database." + c)) }
	dbParams := " " + "user=" + cf("user")
	dbParams += " " + "host=" + cf("host")
	dbParams += " " + "port=" + cf("port")
	dbParams += " " + "dbname=" + cf("dbname")
	if password := cf("password"); password != "" {
		dbParams += " " + "password=" + password
	}
	dbParams += " " + "sslmode=" + cf("sslMode")
	dbParams += " " + "connect_timeout=" + ci("connectionTimeout")
	dbParams += " " + "statement_timeout=" + ci("statementTimeout")
	dbParams += " " + "idle_in_transaction_session_timeout=" + ci("idleTransacionTimeout")
	db, err := postgres.NewStorage(dbParams)
	if err != nil {
		return nil, err
	}
	return db, nil
}