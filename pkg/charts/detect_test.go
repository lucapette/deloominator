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
		{expected: charts.PieChart, dataTypes: charts.DataTypes{charts.Number, charts.Text}},
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
