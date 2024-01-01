package main

import (
	"flag"
)

func main() {
	var (
		dataManagerAddress string
		appPort            string
	)

	flag.StringVar(&dataManagerAddress, "data_manager_address", "localhost:8093", "address of data manager")
	flag.StringVar(&appPort, "app_port", "8094", "application port")

	flag.Parse()
}
