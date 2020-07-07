package template

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"rack/pkg/source"
)

const huluSourceFile = "res/hulu/hulu_task.json"

type Carrier struct {
	ID             string      `json:"id"`
	Name           string      `json:"name"`
	Linkitem       string      `json:"linkitem"`
	Device         string      `json:"device"`
	Property       string      `json:"property"`
	Desc           string      `json:"desc"`
	DataType       interface{} `json:"dataType"`
	FixedValue     interface{} `json:"fixedValue"`
	ParameterValue interface{} `json:"parameterValue"`
}

type CarrierMap struct {
	Name     string
	DataType interface{}
	Value    interface{}
}

type Hulu interface {
	SetupInputs() (input []CarrierMap)
	Processor(input []CarrierMap) (output []CarrierMap)
}

type HuluContext struct {
	Ctx context.Context
}

type huluMockContext struct {
}

func NewHuluCtx(ctx context.Context) HuluContext {
	return HuluContext{Ctx: ctx}
}

func (huluCtx HuluContext) Run(hulu Hulu) {
	inputs := hulu.SetupInputs()
	originInputs, originOuts, err := HuluInputOutputSource()
	if err != nil {
		return
	}
	//device:property
	// origInputMap := make(map[string]string)

	// for _, oriIn := range originInputs {
	// 	origInputMap[oriIn.Name] = oriIn.Property
	// }

	edgeXSource := source.GetEdgeXSource()
	dataStream := edgeXSource.DataStream
	// count := len(origInputMap)
	select {
	case <-huluCtx.Ctx.Done():
		return
	case event := <-dataStream:
		fmt.Printf("Recv edgex data: %v\n", event)
	// for _, read := range event.Readings {
	// 	if count >= 0 {

	// 	}
	// }
	default:
		setupOriginInputValue(originInputs, inputs)
		process(inputs, originOuts, hulu)
	}

}

//主动采集
func setupOriginInputValue(originInputs []Carrier, inputs []CarrierMap) {
	for _, oriIn := range originInputs {
		for i, in := range inputs {
			if oriIn.Name == in.Name {
				result, err := source.GetEdgeXSource().SendCommand(oriIn.Name, oriIn.Property, nil)
				if err != nil {
					inputs[i].Value = nil
				}
				inputs[i].Value = result.Readings[0].Value
			}
		}
	}
}

//只发送put设置命令，支持get采集命令？
func process(inputs []CarrierMap, originOuts []Carrier, hulu Hulu) {
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
		go source.GetEdgeXSource().SendCommand(carrier.Device, carrier.Property, param)
	}
}

func HuluInputOutputSource() (inputs, outputs []Carrier, err error) {

	source, err := ioutil.ReadFile(huluSourceFile)
	if err != nil {
		return nil, nil, err
	}
	mod := struct {
		Inputs  []Carrier `json:"inputs"`
		Outputs []Carrier `json:"outputs"`
	}{}
	if err = json.Unmarshal(source, &mod); err != nil {
		return nil, nil, err
	}
	inputs = mod.Inputs
	outputs = mod.Outputs
	return inputs, outputs, nil
}
