package main

import (
    "context"
    "log"

    "github.com/plus/taskflow/backend/internal/app"
)

func main() {
    if err := app.Run(context.Background()); err != nil {
        log.Fatal(err)
    }
}
