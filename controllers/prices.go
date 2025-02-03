package controllers
import (
  "log"
  "net/http"

  "project_sem/models"
  "project_sem/views"
)

type PricesController struct {}

func (h *PricesController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  switch {
    case r.Method == http.MethodGet:
      h.Get(w, r)
      return
    case r.Method == http.MethodPost:
      h.Create(w, r)
      return
    default:
      notFound(w, r)
      return
  }
}

func (h *PricesController) Get(w http.ResponseWriter, r *http.Request) {
  response, err := views.CreatePricesCsvZip()

  if err != nil {
    log.Println(err)
    internalServerError(w, r)
    return
  }

  w.Header().Set("Content-Type", "application/zip")
  w.Header().Set("Content-Disposition", "attachment; filename=data.zip")
  w.Write(response.Bytes())
}

func (h *PricesController) Create(w http.ResponseWriter, r *http.Request) {
  r.Body = http.MaxBytesReader(w, r.Body, 10<<20)
  defer r.Body.Close()

 totalItems, err := models.StorePricesFromBody(r.Body)
  if err != nil {
    badRequest(w, r)
    return
  }

  pricesSummary := models.SelectPricesSummary(totalItems)
  response, err := views.CreatePricesSummaryJson(pricesSummary)
  if err != nil {
    log.Println(err)
    internalServerError(w, r)
    return
  }

  w.WriteHeader(http.StatusCreated)
  w.Header().Set("Content-Type", "application/json")
  w.Write(response)
}

func internalServerError(w http.ResponseWriter, r *http.Request) {
  w.WriteHeader(http.StatusInternalServerError)
  w.Write([]byte("internal server error"))
}

func notFound(w http.ResponseWriter, r *http.Request) {
  w.WriteHeader(http.StatusNotFound)
  w.Write([]byte("not found"))
}

func badRequest(w http.ResponseWriter, r *http.Request) {
  w.WriteHeader(http.StatusBadRequest)
  w.Write([]byte("bad request"))
}
