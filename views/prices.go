package views
import (
  "bytes"
  "archive/zip"
  "encoding/json"
  "project_sem/models"
)


func CreatePricesSummaryJson(priceSummary *models.PricesSummary) ([]byte, error) {
  return json.Marshal(priceSummary)
}

func CreatePricesCsvZip() (*bytes.Buffer, error) {
  buf := new(bytes.Buffer)
  w := zip.NewWriter(buf)

  f, err := w.Create("data.csv")
  if err != nil {
    return nil, err
  }

  csvContent := "id,name,category,price,create_date\n"
  // TODO replace this loop by actual data from db
  for i := 0; i < 10; i++ {
    csvContent += "123,fsdfsf,bxcbcxb,15,12323\n"
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
