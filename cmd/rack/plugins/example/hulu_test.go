// +build ignore
package main

import (
	"fmt"
	"rack/pkg/template"
	"testing"
)

func TestProcessor(t *testing.T) {
	input_1 := template.CarrierMap{Name: "temper_1", DataType: template.ValueTypeInt32, Value: 26}
	input_2 := template.CarrierMap{Name: "temper_2", DataType: template.ValueTypeInt32, Value: 100}
	inputs := []template.CarrierMap{input_1, input_2}

	outputs := RackContext.Processor(inputs)

	fmt.Printf("outputs: %v\n", outputs)
}
