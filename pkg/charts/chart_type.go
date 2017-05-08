package charts

type ChartType int

//go:generate stringer -type=ChartType -output=chart_type_string.go
const (
	UnknownChart ChartType = iota
	SimpleBar
	SimpleLine
	GroupedBar
	MultiLine
)
