package charts

import "bytes"

type ChartType string

const (
	UnknownChart ChartType = "UnknownChart"
	SimpleBar    ChartType = "SimpleBar"
	SimpleLine   ChartType = "SimpleLine"
)

type DataType string

type DataTypes []DataType

const (
	UnknownType DataType = "UnknownType"
	Text        DataType = "Text"
	Number      DataType = "Number"
	Time        DataType = "Time"
)

var charts map[string]ChartType

func sequence(types DataTypes) string {
	var seq bytes.Buffer

	for _, t := range types {
		seq.WriteString(t.String())
	}

	return seq.String()
}

func Detect(types DataTypes) ChartType {
	if chart, ok := charts[sequence(types)]; ok {
		return chart
	}

	return UnknownChart
}

func (ct ChartType) String() string {
	return string(ct)
}

func (t DataType) String() string {
	return string(t)
}

func init() {
	charts = make(map[string]ChartType)

	charts[sequence(DataTypes{Text, Number})] = SimpleBar
	charts[sequence(DataTypes{Number, Number})] = SimpleBar
	charts[sequence(DataTypes{Time, Number})] = SimpleLine
}
