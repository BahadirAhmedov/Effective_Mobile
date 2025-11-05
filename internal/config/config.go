package config

import (
	"fmt"
	"log"
	"os"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct{
	Env string `yaml:"env" env-default:"local"`
	Storage StorageCredentials `yaml:"storage-credentials"`
	//TODO: Define config fields
}

type StorageCredentials struct{
	Host string `yaml:"host"`
	Port int `yaml:"port"`
	User string `yaml:"user"`
	Password string `yaml:"password"`
	DbName string `yaml:"dbname"`
}

func MustLoad() (*Config) {

	err := godotenv.Load("local.env")
	if err != nil {
		fmt.Println(err)
		log.Fatal("unable to load .env file values")		
	}

	configPath := os.Getenv("CONFIG_PATH")
	fmt.Println(configPath)

	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}
	
	if _, err :=  os.Stat(configPath); os.IsNotExist(err){
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil{
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
