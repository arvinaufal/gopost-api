package config

import (
    "log"
    "os"

    "github.com/jmoiron/sqlx"
    _ "github.com/lib/pq"
)

var SupabaseDB *sqlx.DB

func ConnectSupabaseDB() {
    var err error

    dsn := os.Getenv("POSTGRES_URL")
    if dsn == "" {
        log.Fatal("Missing POSTGRES_URL environment variable")
    }

    SupabaseDB, err = sqlx.Connect("postgres", dsn)
    if err != nil {
        log.Fatal("Error connecting to Supabase database:", err)
    }

    log.Println("Supabase database connected successfully!")
}
