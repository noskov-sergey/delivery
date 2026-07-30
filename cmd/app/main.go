package main

import (
	"database/sql"
	"delivery/cmd"
	"delivery/internal/generated/servers"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/gommon/log"
	"github.com/pressly/goose/v3"
	"github.com/robfig/cron/v3"

	server "delivery/internal/adapters/in/http"
)

func main() {
	cfg := getConfigs()

	dsn, err := makeConnectionString(cfg.DbHost, cfg.DbPort, cfg.DbUser, cfg.DbPassword, cfg.DbName, cfg.DbSslMode)
	if err != nil {
		log.Fatal(err.Error())
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = goose.Up(db, "migrations")
	if err != nil {
		log.Fatal(err)
	}

	compositionRoot := cmd.NewCompositionRoot(
		cfg, db,
	)
	defer compositionRoot.CloseAll()

	startCron(compositionRoot)

	startKafkaConsumer(compositionRoot)

	startWebServer(compositionRoot, cfg.HttpPort)
}

func getConfigs() cmd.Config {
	config := cmd.Config{
		HttpPort:               goDotEnvVariable("HTTP_PORT"),
		DbHost:                 goDotEnvVariable("DB_HOST"),
		DbPort:                 goDotEnvVariable("DB_PORT"),
		DbUser:                 goDotEnvVariable("DB_USER"),
		DbPassword:             goDotEnvVariable("DB_PASSWORD"),
		DbName:                 goDotEnvVariable("DB_NAME"),
		DbSslMode:              goDotEnvVariable("DB_SSLMODE"),
		GeoServiceGrpcHost:     goDotEnvVariable("GEO_SERVICE_GRPC_HOST"),
		KafkaHost:              goDotEnvVariable("KAFKA_HOST"),
		KafkaConsumerGroup:     goDotEnvVariable("KAFKA_CONSUMER_GROUP"),
		KafkaBasketEventsTopic: goDotEnvVariable("KAFKA_BASKET_EVENTS_TOPIC"),
		KafkaOrderEventsTopic:  goDotEnvVariable("KAFKA_ORDER_EVENTS_TOPIC"),
	}
	return config
}

func goDotEnvVariable(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	return os.Getenv(key)
}

func startWebServer(cr *cmd.CompositionRoot, port string) {
	// create a type that satisfies the `api.ServerInterface`, which contains an implementation of every operation from the generated code
	server, err := server.NewServer(
		cr.NewCreateOrderCommandHandler(),
		cr.NewGetIncompletedOrdersQueryHandler(),
		cr.NewCreateCourierCommandHandler(),
		cr.NewGetAllCouriersQueryHandler())
	if err != nil {
		log.Fatal(err)
	}

	r := http.NewServeMux()

	// get an `http.Handler` that we can use
	h := servers.HandlerFromMux(server, r)

	s := &http.Server{
		Handler: h,
		Addr:    fmt.Sprintf("0.0.0.0:%s", port),
	}

	// And we serve HTTP until the world ends.

	log.Info("starting server on port " + port + " ...")
	log.Fatal(s.ListenAndServe(), h)
}

func makeConnectionString(host string, port string, user string,
	password string, dbName string, sslMode string,
) (string, error) {
	if host == "" {
		return "", fmt.Errorf("host is required")
	}
	if port == "" {
		return "", fmt.Errorf("port is required")
	}
	if user == "" {
		return "", fmt.Errorf("user is required")
	}
	if password == "" {
		return "", fmt.Errorf("password is required")
	}
	if dbName == "" {
		return "", fmt.Errorf("dbName is required")
	}
	if sslMode == "" {
		return "", fmt.Errorf("sslMode is required")
	}
	return fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=%v",
		host, port, user, password, dbName, sslMode), nil
}

func startCron(compositionRoot *cmd.CompositionRoot) {
	c := cron.New()
	_, err := c.AddJob("@every 1s", compositionRoot.NewAssignOrdersJob())
	if err != nil {
		log.Fatalf("ошибка при добавлении задачи: %v", err)
	}
	_, err = c.AddJob("@every 1s", compositionRoot.NewMoveCouriersJob())
	if err != nil {
		log.Fatalf("ошибка при добавлении задачи: %v", err)
	}
	c.Start()
}

func startKafkaConsumer(compositionRoot *cmd.CompositionRoot) {
	go func() {
		if err := compositionRoot.NewBasketConsumer().Consume(); err != nil {
			log.Fatalf("ERROR: kafka consumer got an error: %v", err)
		}
	}()
}
