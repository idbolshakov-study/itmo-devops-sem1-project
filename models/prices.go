package models
import (
  "io"
  "log"
  "fmt"
  "time"
  "bytes"
  "context"
  "strings"
  "strconv"
  "archive/zip"
  "encoding/csv"

  "github.com/jackc/pgx/v5/pgxpool"
)

var PgxPool *pgxpool.Pool

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

func StorePricesFromBody(body io.ReadCloser) (int, error) {
  bodyBytes, err := io.ReadAll(body)
  if err != nil {
    log.Println("Error reading request body:", err)
    return 0, err
  }

  zipReader, err := zip.NewReader(bytes.NewReader(bodyBytes), int64(len(bodyBytes)))
  if err != nil {
    log.Println("Error reading zip archive:", err)
    return 0, err
  }

  var csvReader *csv.Reader = nil
  for _, file := range zipReader.File {
    if file.Name != "data.csv" {
      log.Println("data.csv not found in zip archive", file.Name)
      return 0, fmt.Errorf("data.csv no found in zip archive")
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
    return 0, fmt.Errorf("Error while reading csv header: %w", err)
  }

  // read csv rows
  for {
    record, err := csvReader.Read()

    if err == io.EOF {
      break
    }

    if err != nil {
      return 0, fmt.Errorf("Error while reading csv record: %w", err)
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

  var rows []string
  for _, p := range products {
    rows = append(
      rows,
      fmt.Sprintf("('%s','%s',%6.1f,'%s')", p.Name, p.Category, p.Price, p.CreateDate.Format("2006-01-02")),
    )
  }

  PgxPool.QueryRow(context.Background(), fmt.Sprintf(
    "%s %s",
    "INSERT INTO prices (name,category,price,create_date) VALUES",
    strings.Join(rows[:], ","),
  ))

  return len(products), nil
}


func SelectPricesSummary(totalItems int) *PricesSummary {
  ps := PricesSummary{}
  ps.TotalItems = totalItems

  query := "SELECT COUNT(DISTINCT category) AS total_categories, SUM(price) AS total_price FROM prices;"
  PgxPool.QueryRow(context.Background(), query).Scan(&ps.TotalCategories, &ps.TotalPrice)

  return &ps
}
