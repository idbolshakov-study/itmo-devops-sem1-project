package models
import (
  "io"
  "log"
  "fmt"
  "time"
  "bytes"
  "strconv"
  "archive/zip"
  "encoding/csv"
)

type Product struct {
  Id         int
  Name       string
  Category   string
  Price      float32
  CreateDate time.Time
}

type PricesSummary struct {
  TotalItems      int     `json:"total_items"`
  TotalCategories int     `json:"total_categories"`
  TotalPrice      float32 `json:"total_price"`
}

func StorePricesFromBody(body io.ReadCloser) (error) {
  bodyBytes, err := io.ReadAll(body)
  if err != nil {
    log.Println("Error reading request body:", err)
    return err
  }

  zipReader, err := zip.NewReader(bytes.NewReader(bodyBytes), int64(len(bodyBytes)))
  if err != nil {
    log.Println("Error reading zip archive:", err)
    return err
  }

  var csvReader *csv.Reader = nil
  for _, file := range zipReader.File {
    if file.Name != "data.csv" {
      log.Println("data.csv not found in zip archive", file.Name)
      return fmt.Errorf("data.csv no found in zip archive")
    }

    fileReader, err := file.Open()
    if err != nil {
      log.Println("Error opening file inside zip:", err)
    }
    defer fileReader.Close()

    csvReader = csv.NewReader(fileReader)
  }

  var products []Product

  // read header first
  _, err = csvReader.Read()
  if err != nil {
    return fmt.Errorf("Error while reading csv header: %w", err)
  }

  // read csv rows
  for {
    record, err := csvReader.Read()

    if err == io.EOF {
      break
    }

    if err != nil {
      return fmt.Errorf("Error while reading csv record: %w", err)
    }

    if len(record) != 5 {
      log.Println("Skipping malformed CSV record (incorrect number of columns):", record)
      continue
    }

    price, err := strconv.ParseFloat(record[3], 32)
    if err != nil {
      log.Println("Skipping malformed CSV record (invalid Price):", record, err)
      continue
    }

    createDate, err := time.Parse("2006-01-02", record[4])
    if err != nil {
      log.Println("Skipping malformed CSV record (invalid Date):", record, err)
      continue
    }

    product := Product{
      Name:       record[1],
      Category:   record[2],
      Price:      float32(price),
      CreateDate: createDate,
    }

    products = append(products, product)
  }

  log.Println(len(products))

  return nil
}


func SelectPricesSummary() *PricesSummary {
  ps := PricesSummary{}

  ps.TotalItems = 12
  ps.TotalCategories = 10
  ps.TotalPrice = 1000.5

  return &ps
}
