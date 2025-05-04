// filepath: d:\VS code\WebSocket\EngageSync\Service\config.go
package service

import (
    "log"
    "os"

    "github.com/joho/godotenv"
)

var SCHEMA string

func InitConfig() {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }
    SCHEMA = os.Getenv("DB_SCHEMA")
}