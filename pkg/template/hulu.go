package template

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
