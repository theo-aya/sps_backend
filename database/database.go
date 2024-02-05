package database

import (
    "fmt"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "os"
)

var DB *gorm.DB

func Connect() {
    var err error 
    host := os.Getenv("DB_HOST")
    username := os.Getenv("DB_USER")
    password := os.Getenv("DB_PASSWORD")
    databaseName := os.Getenv("DB_NAME")
    port := os.Getenv("DB_PORT")

    dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Africa/Accra ", host, username, password, databaseName, port)
    DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
        DisableAutomaticPing: true,
        PrepareStmt: false,
    })

    if err != nil {
        panic("fail to connect database")
 
    } else {
        fmt.Println("Successfully connected to the database")
    }

    

   
}