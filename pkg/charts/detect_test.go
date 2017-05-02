package charts_test

import (
	"testing"

	"github.com/lucapette/deloominator/pkg/charts"
)

func TestDetect(t *testing.T) {
	tests := []struct {
		expected  charts.ChartType
		dataTypes charts.DataTypes
	}{
		{expected: charts.UnknownChart, dataTypes: charts.DataTypes{charts.UnknownType, charts.Number}},
		{expected: charts.SimpleBar, dataTypes: charts.DataTypes{charts.Number, charts.Number}},
		{expected: charts.SimpleBar, dataTypes: charts.DataTypes{charts.Text, charts.Number}},
		{expected: charts.SimpleLine, dataTypes: charts.DataTypes{charts.Time, charts.Number}},
		{expected: charts.GroupedBar, dataTypes: charts.DataTypes{charts.Text, charts.Text, charts.Number}},
		{expected: charts.MultiLine, dataTypes: charts.DataTypes{charts.Time, charts.Text, charts.Number}},
	}

	for _, test := range tests {
		t.Run(test.expected.String(), func(t *testing.T) {
			actual := charts.Detect(test.dataTypes)

			if actual != test.expected {
				t.Fatalf("Expected %s, but got %s", test.expected, actual)
			}
		})
	}
}
