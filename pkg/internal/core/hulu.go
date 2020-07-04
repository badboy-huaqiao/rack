package core

import (
	"encoding/json"
	"io/ioutil"
	"plugin"
	"rack/pkg/source"
	"rack/pkg/template"
)

const huluSourceFile = "res/hulu/hulu_task.json"

type huluMockContext struct {
}

func HuluInputOutputSource() (inputs, outputs []template.Carrier, err error) {

	source, err := ioutil.ReadFile(huluSourceFile)
	if err != nil {
		return nil, nil, err
	}
	mod := struct {
		Inputs  []template.Carrier `json:"inputs"`
		Outputs []template.Carrier `json:"outputs"`
	}{}
	if err = json.Unmarshal(source, &mod); err != nil {
		return nil, nil, err
	}
	inputs = mod.Inputs
	outputs = mod.Outputs
	return inputs, outputs, nil
}

func executeHulu(RackContextSym plugin.Symbol) {
	hulu := RackContextSym.(template.Hulu)
	inputs := hulu.SetupInputs()
	originInputs, originOuts, err := HuluInputOutputSource()
	if err != nil {
		return
	}
	edgeXSource := source.EdgeXSource{}
	for _, oriIn := range originInputs {
		for i, in := range inputs {
			if oriIn.Name == in.Name {
				result, err := edgeXSource.SendCommand(oriIn.Name, oriIn.Property, nil)
				if err != nil {
					inputs[i].Value = nil
				}
				inputs[i].Value = result.Readings[0].Value
			}
		}
	}

	outputs := hulu.Processor(inputs)
	for i, oriOut := range originOuts {
		for _, out := range outputs {
			if oriOut.Name == out.Name {
				originOuts[i].ParameterValue = out.Value
			}
		}
	}

	for _, carrier := range originOuts {
		param := map[string]interface{}{carrier.Property: carrier.ParameterValue}
		go edgeXSource.SendCommand(carrier.Device, carrier.Property, param)
	}

}
