package charts

import "bytes"

type ChartType int

const (
	UnknownChart ChartType = iota
	SimpleBar
	SimpleLine
	GroupedBar
	MultiLine
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
	case SimpleBar:
		return "SimpleBar"
	case SimpleLine:
		return "SimpleLine"
	case GroupedBar:
		return "GroupedBar"
	case MultiLine:
		return "MultiLine"
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

	charts[sequence(DataTypes{Text, Number})] = SimpleBar
	charts[sequence(DataTypes{Number, Number})] = SimpleBar
	charts[sequence(DataTypes{Time, Number})] = SimpleLine
	charts[sequence(DataTypes{Text, Text, Number})] = GroupedBar
	charts[sequence(DataTypes{Time, Text, Number})] = MultiLine
}
