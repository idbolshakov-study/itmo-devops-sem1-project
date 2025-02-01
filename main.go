package main
import (
  "net/http"
  "project_sem/controllers"
)


func main() {
  mux := http.NewServeMux()
  mux.Handle("/api/v0/prices", &controllers.PricesController{})
  http.ListenAndServe(":8080", mux)
}
