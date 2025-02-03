package views
import (
  "fmt"
  "bytes"
  "archive/zip"
  "encoding/json"
  "project_sem/models"
)


func CreatePricesSummaryJson(priceSummary *models.PricesSummary) ([]byte, error) {
  return json.Marshal(priceSummary)
}

func CreatePricesCsvZip(products []models.Product) (*bytes.Buffer, error) {
  buf := new(bytes.Buffer)
  w := zip.NewWriter(buf)

  f, err := w.Create("data.csv")
  if err != nil {
    return nil, err
  }

  csvContent := "id,name,category,price,create_date\n"
  for _, p := range products {
    csvContent += fmt.Sprintf(
      "%d,%s,%s,%6.2f,%s\n",
      p.Id, p.Name, p.Category, p.Price, p.CreateDate.Format("2006-01-02"),
    )
  }

  _, err = f.Write([]byte(csvContent))
  if err != nil {
    return nil, err
  }

  err = w.Close()
  if err != nil {
    return nil, err
  }

  return buf, nil
}
