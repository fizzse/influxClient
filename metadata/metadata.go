package metadata

/*
 * influx value number or string
 */

type UnitData struct {
	NumberValue *float64 // 数值类型
	StringValue *string  // 字符串类型
}

func (d *UnitData) GetNumber() *float64 {
	return d.NumberValue
}

func (d *UnitData) GetString() *string {
	return d.StringValue
}
