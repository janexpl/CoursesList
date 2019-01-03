package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/BurntSushi/toml"
	_ "github.com/lib/pq"
)

var DB *sql.DB

type Config struct {
	Server   string
	User     string
	Dbname   string
	Password string
	Sslmode  string
}

func readConfig() Config {
	var configfile = "config/config.toml"
	_, err := os.Stat(configfile)
	if err != nil {
		log.Fatal("Config file is missing: ", configfile)
	}
	var config Config
	if _, err := toml.DecodeFile(configfile, &config); err != nil {
		log.Fatal(err)
	}

	//log.Print(config.Index)
	return config
}

func init() {
	var err error
	config := readConfig()
	connStr := fmt.Sprintf("user=%v dbname=%v password=%v host=%v sslmode=%v", config.User, config.Dbname, config.Password, config.Server, config.Sslmode)
	//connStr := "user=certs dbname=courselist password=Acetamie2018 host=localhost sslmode=disable"
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	if err = DB.Ping(); err != nil {
		panic(err)
	}
	fmt.Println("You connected to your database.")
}
