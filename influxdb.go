package checkup

import (
  "fmt"
  "strings"
	"net/http"
  "net/url"
	"log"
  "io/ioutil"
  "strconv"
)

type InfluxDB struct {
  Endpoint string `json:"endpoint"`
  DatabaseName string `json:"database_name"`
}

func (influx InfluxDB) Store(results []Result) error {
  db := influx.DatabaseName
  for _, r := range results {
      endpoint := r.Title
      timestamp := r.Timestamp
      // rtt := r.Times[0].RTT

      var status int
      if r.Healthy {
        status = 1
      } else if strings.Contains(r.Times[0].Error, "Timeout") {
        status = 2
      } else {
        // Degraded
        // Healthy
        status = 0
      }

      influx.Submit(fmt.Sprintf("%s,endpoint=\"%s\" status=%d %d", db, endpoint, status, timestamp))
    }

    return nil
}

func (influx InfluxDB) Submit(query string) {
  log.Printf(query)
  u, _ := url.ParseRequestURI(influx.Endpoint)
  q, _ := url.ParseQuery(u.RawQuery)
  q.Add("db", influx.DatabaseName)

  data := url.Values{}
  data.Set("q", query)
  log.Printf(data.Encode())
  u.Path = "/write"
  u.RawQuery = q.Encode()
  urlStr := u.String()
  log.Printf(urlStr)

  req, _ := http.NewRequest(http.MethodPost, urlStr, strings.NewReader(data.Encode()))
  req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
  req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

  client := &http.Client{}
  resp, _ := client.Do(req)
  log.Printf(resp.Status)

  bodyBytes, _ := ioutil.ReadAll(resp.Body)
	body := string(bodyBytes)
  log.Printf(body)

  resp.Body.Close()
}
