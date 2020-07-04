// +build ignore
package main

import "rack/pkg/template"

type Hulu struct {
}

// RackContext 所有实现接口的plugin导出变量必须以RackContext命名
var RackContext Hulu

func (hulu Hulu) SetupInputs() (inputs []template.CarrierMap) {

	carrier := template.CarrierMap{
		Name:     "temper_1",
		DataType: template.ValueTypeInt32,
	}
	inputs = append(inputs, carrier)
	carrier = template.CarrierMap{
		Name:     "temper_2",
		DataType: template.ValueTypeInt32,
	}

	inputs = append(inputs, carrier)
	return inputs
}

func (hulu Hulu) Processor(inputs []template.CarrierMap) (outputs []template.CarrierMap) {
	for _, input := range inputs {
		if input.Name == "temper_1" && input.Value.(int) > 25 {
			outputs = append(outputs, template.CarrierMap{
				Name:     "humidity_1",
				DataType: template.ValueTypeInt32,
				Value:    40,
			})

		}

		if input.Name == "temper_2" && input.Value.(int) > 35 {
			outputs = append(outputs, template.CarrierMap{
				Name:     "humidity_2",
				DataType: template.ValueTypeInt32,
				Value:    50,
			})
		}

	}
	return outputs
}
