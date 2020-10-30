package main

import (
	"context"
	"fmt"
	"log"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

func main() {
	// Create a new client using an InfluxDB server base URL and an authentication token
	client := influxdb2.NewClient("http://106.13.76.237:8086", "my-token")
	defer client.Close()

	//// Use blocking write client for writes to desired bucket
	//writeAPI := client.WriteAPIBlocking("my-org", "iotdata")
	//// Create point using full params constructor
	//p := influxdb2.NewPoint("stat",
	//	map[string]string{"unit": "temperature"},
	//	map[string]interface{}{"avg": 24.5, "max": 45.0},
	//	time.Now())
	//// write point immediately
	//writeAPI.WritePoint(context.Background(), p)
	//// Create point using fluent style
	//p = influxdb2.NewPointWithMeasurement("stat").
	//	AddTag("unit", "temperature").
	//	AddField("avg", 23.2).
	//	AddField("max", 45.0).
	//	SetTime(time.Now())
	//writeAPI.WritePoint(context.Background(), p)
	//
	//// Or write directly line protocol
	//line := fmt.Sprintf("stat,unit=temperature avg=%f,max=%f", 23.5, 45.0)
	//writeAPI.WriteRecord(context.Background(), line)

	// Get query client
	queryAPI := client.QueryAPI("my-org")
	// Get parser flux query result
	result, err := queryAPI.Query(context.Background(), `from(bucket:"iotdata")|> range(start: -1h) |> filter(fn: (r) => r._measurement == "stat")`)
	if err != nil {
		log.Printf("query failed: %s\n", err.Error())
		return
	}

	for result.Next() {
		// Observe when there is new grouping key producing new table
		//if result.TableChanged() {
		//	fmt.Printf("table: %s\n", result.TableMetadata().String())
		//}
		// read result
		//result.Record().ValueByKey("")
		//for k, v := range result.Record().Values() {
		//	fmt.Printf("key: %v,value:%v\t", k, v)
		//}

		result.Record().Values()
		fmt.Println(result.Record().Value())
		//fmt.Println(result.Record().ValueByKey("avg"))

		fmt.Println()

		//fmt.Printf("row: %s\n", result.Record().String())
	}
	if result.Err() != nil {
		fmt.Printf("Query error: %s\n", result.Err().Error())
	}
	// Ensures background processes finishes
}
