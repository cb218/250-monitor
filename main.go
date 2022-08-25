package main

import (
	"250-monitor/probers"
	"fmt"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"strconv"
)

type Configuration struct {
	HostIP     string `yaml:"HOST_IP"`
	PingPeriod string `yaml:"PING_PERIOD"`
}

func loadConfiguration(path string) *Configuration {
	yfile, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("Could not open %s error: %s\n", path, err)
		conf := &Configuration{"127.0.0.1", "60"}
		fmt.Printf("Host Monitor will use default configuration: %v\n", conf)
		return conf
	}
	if yfile == nil {
		panic("There was no error but YFile was null")
	}
	conf := Configuration{}
	err2 := yaml.Unmarshal(yfile, &conf)
	if err2 != nil {
		fmt.Printf("Configuration file could not be parsed, error: %s\n", err2)
		panic(err2)
	}
	fmt.Printf("Found configuration: %v\n", conf)
	return &conf
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	conf := loadConfiguration("configuration.yaml")

	hostIP := conf.HostIP
	fmt.Printf("Host to be monitored: %s\n", hostIP)

	pingPeriod := 60
	if conf.PingPeriod != "" {
		pingPeriod, err = strconv.Atoi(conf.PingPeriod)
		if err != nil {
			fmt.Printf("Error converting %s to integer, ping period set to default (60)", conf.PingPeriod)
		}
		fmt.Printf("Host will be pinged every %v minutes\n", conf.PingPeriod)
	}

	monitor := probers.NewMonitor()
	monitor.Start(hostIP, pingPeriod)

}
