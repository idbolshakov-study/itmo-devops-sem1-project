package main
import (
  "os"
  "fmt"
  "context"

  "net/http"

  "github.com/jackc/pgx/v5/pgxpool"

  "project_sem/models"
  "project_sem/controllers"
)


func main() {
  dbpool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
  if err != nil {
    panic(fmt.Sprintf("Unable to create connection pool: %v\n", err))
  }
  defer dbpool.Close()

  models.PgxPool = dbpool

  mux := http.NewServeMux()
  mux.Handle("/api/v0/prices", &controllers.PricesController{})
  http.ListenAndServe(":8080", mux)
}
