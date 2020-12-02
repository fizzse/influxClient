package metadata

import (
	"context"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

func NewCpuInfo() *CpuInfo {
	res := new(CpuInfo)
	res.Tag = new(CpuInfoTag)
	res.Filed = new(CpuInfoField)
	return res
}

type CpuInfoTag struct {
	Host string
}

type CpuInfoField struct {
	Avg float64
	Max float64
}

type CpuInfo struct {
	Timestamp int64
	Tag       *CpuInfoTag
	Filed     *CpuInfoField
}

func (*CpuInfo) TableName() string {
	return "cpu"
}

func (c *CpuInfo) QueryAll() {

}

func (c *CpuInfo) Insert(conn influxdb2.Client) error {
	point := influxdb2.NewPointWithMeasurement(c.TableName())

	if c.Tag.Host != "" {
		point.AddTag("host", c.Tag.Host)
	}

	point.AddField("avg", c.Filed.Avg)
	point.AddField("max", c.Filed.Max)
	writeAPI := conn.WriteAPIBlocking("my-org", "iotdata")

	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	err := writeAPI.WritePoint(ctx, point)
	return err
}
