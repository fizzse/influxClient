package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

const (
	bucketSql = `from(bucket:"%s")`
	rangeSql  = `|> range(%s)`
	filterSql = `|> filter(fn: (r) => %s)`
)

func sqlBuild(bucket string, timeRange map[string]string, measurement string, tagFilter map[string]string) string {
	buf := &strings.Builder{}
	buf.WriteString(fmt.Sprintf(bucketSql, bucket))

	if len(timeRange) > 0 {
		i := 0
		var timeRangeStr string
		for k, v := range timeRange {
			if i != 0 {
				timeRangeStr += ","
			}

			timeRangeStr += k
			timeRangeStr += ":"
			timeRangeStr += v
			i++
		}

		buf.WriteString(fmt.Sprintf(rangeSql, timeRangeStr))
	}

	if measurement != "" {
		var filterStr string
		//filterStr += "r._measurement=="
		//filterStr += measurement

		filterStr = fmt.Sprintf(`r._measurement == "%s"`, measurement)

		if len(tagFilter) > 0 {
			for k, v := range tagFilter {
				subFilter := fmt.Sprintf(` and %s "%s"`, k, v)
				filterStr += subFilter
			}
		}

		buf.WriteString(fmt.Sprintf(filterSql, filterStr))
	}

	str := buf.String()
	log.Println(str)
	return str
}

func main() {
	// Create a new client using an InfluxDB server base URL and an authentication token
	client := influxdb2.NewClient("http://49.232.147.12:8086", "my-token")
	defer client.Close()

	//cpuInfo := metadata.NewCpuInfo()
	//cpuInfo.Tag.Host = "127.0.0.1"
	//cpuInfo.Filed.Avg = 56.1
	//cpuInfo.Filed.Max = 100.1
	//
	//if err := cpuInfo.Insert(client); err != nil {
	//	log.Fatal(err)
	//}

	// Get query client
	queryAPI := client.QueryAPI("my-org")
	// Get parser flux query result
	//sql := fmt.Sprintf("FROM %s LIMIT %d", "cpu", 10)

	sql := sqlBuild("iotdata", map[string]string{"start": "-1h"}, "cpu",
		map[string]string{
			"r.host == ": "127.0.0.1",
		})

	//sql := sqlBuild("iotdata", map[string]string{"start": "-1h"}, "cpu", nil)
	result, err := queryAPI.Query(context.Background(), sql)
	if err != nil {
		log.Printf("query failed: %s\n", err.Error())
		return
	}

	if result.Err() != nil {
		log.Fatal("result err", result.Err())
	}

	i := 0
	for result.Next() {

		fmt.Println("line:", i)
		i++

		fmt.Printf("time: %v\t", result.Record().Time())

		result.Record().Values()

		for k, v := range result.Record().Values() {
			fmt.Printf("key: %v,value: %v\t", k, v)
		}

		fmt.Println()
	}

}
