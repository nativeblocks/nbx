package model

type FrameDSLModel struct {
	Name      string             `json:"name"`
	Route     string             `json:"route"`
	Type      string             `json:"type"`
	Variables []VariableDSLModel `json:"variables"`
	Blocks    []BlockDSLModel    `json:"blocks"`
}

type VariableDSLModel struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Type  string `json:"type"`
}

type BlockDSLModel struct {
	KeyType            string                  `json:"keyType"`
	Key                string                  `json:"key"`
	VisibilityKey      string                  `json:"visibilityKey"`
	Slot               string                  `json:"slot,omitempty"`
	IntegrationVersion int                     `json:"integrationVersion"`
	Data               []BlockDataDSLModel     `json:"data"`
	Properties         []BlockPropertyDSLModel `json:"properties"`
	Slots              []BlockSlotDSLModel     `json:"slots"`
	Blocks             []BlockDSLModel         `json:"blocks"`
	Actions            []ActionDSLModel        `json:"actions"`
}

type BlockPropertyDSLModel struct {
	Key          string `json:"key"`
	ValueMobile  string `json:"valueMobile"`
	ValueTablet  string `json:"valueTablet"`
	ValueDesktop string `json:"valueDesktop"`
	Type         string `json:"type"`
}

type BlockDataDSLModel struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Type  string `json:"type"`
}

type BlockSlotDSLModel struct {
	Slot string `json:"slot"`
}

type ActionDSLModel struct {
	Key      string                  `json:"key"`
	Event    string                  `json:"event"`
	Triggers []ActionTriggerDSLModel `json:"triggers"`
}

type ActionTriggerDSLModel struct {
	KeyType            string                    `json:"keyType"`
	Then               string                    `json:"then"`
	Name               string                    `json:"name"`
	IntegrationVersion int                       `json:"integrationVersion"`
	Properties         []TriggerPropertyDSLModel `json:"properties"`
	Data               []TriggerDataDSLModel     `json:"data"`
	Triggers           []ActionTriggerDSLModel   `json:"triggers"`
}

type TriggerPropertyDSLModel struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Type  string `json:"type"`
}

type TriggerDataDSLModel struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Type  string `json:"type"`
}
