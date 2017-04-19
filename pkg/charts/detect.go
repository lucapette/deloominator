package charts

import "bytes"

type ChartType int

const (
	UnknownChart ChartType = iota
	PieChart
)

type DataType int

type DataTypes []DataType

const (
	UnknownType DataType = iota
	Text
	Number
	Time
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
	switch ct {
	case PieChart:
		return "PieChart"
	}

	return "UnknownChart"
}

func (t DataType) String() string {
	switch t {
	case Text:
		return "Text"
	case Number:
		return "Number"
	case Time:
		return "Time"
	}

	return "Unknown"
}

func init() {
	charts = make(map[string]ChartType)

	charts[sequence(DataTypes{Number, Text})] = PieChart
}
