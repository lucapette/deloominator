package charts

import "bytes"

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

func init() {
	charts = make(map[string]ChartType)

	charts[sequence(DataTypes{Text, Number})] = SimpleBar
	charts[sequence(DataTypes{Number, Number})] = SimpleBar
	charts[sequence(DataTypes{Time, Number})] = SimpleLine
	charts[sequence(DataTypes{Text, Text, Number})] = GroupedBar
	charts[sequence(DataTypes{Time, Text, Number})] = MultiLine
}
