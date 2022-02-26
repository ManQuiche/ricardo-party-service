package main

import "ricardo/party-service/boot"

func main() {
	boot.LoadEnv()
	boot.LoadDb()
	boot.LoadServices()

	boot.ServeHTTP()
}
