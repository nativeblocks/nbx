package model

type FrameJson struct {
	Id             string              `json:"id"`
	Name           string              `json:"name"`
	Route          string              `json:"route"`
	RouteArguments []RouteArgumentJson `json:"routeArguments"`
	Type           string              `json:"type"`
	IsStarter      bool                `json:"isStarter"`
	ProjectId      string              `json:"projectId"`
	Checksum       string              `json:"checksum"`
	Variables      []VariableJson      `json:"variables"`
	Blocks         []BlockJson         `json:"blocks"`
	Actions        []ActionJson        `json:"actions"`
}

type RouteArgumentJson struct {
	Name string `json:"name"`
}

type VariableJson struct {
	Id      string `json:"id"`
	FrameId string `json:"frameId"`
	Key     string `json:"key"`
	Value   string `json:"value"`
	Type    string `json:"type"`
}

type BlockJson struct {
	Id                          string              `json:"id"`
	FrameId                     string              `json:"frameId"`
	KeyType                     string              `json:"keyType"`
	Key                         string              `json:"key"`
	VisibilityKey               string              `json:"visibilityKey"`
	Position                    int                 `json:"position"`
	Slot                        string              `json:"slot"`
	IntegrationVersion          int                 `json:"integrationVersion"`
	ParentId                    string              `json:"parentId"`
	Data                        []BlockDataJson     `json:"data"`
	Properties                  []BlockPropertyJson `json:"properties"`
	Slots                       []BlockSlotJson     `json:"slots"`
	IntegrationDeprecated       bool                `json:"integrationDeprecated"`
	IntegrationDeprecatedReason string              `json:"integrationDeprecatedReason"`
}

type BlockPropertyJson struct {
	BlockId            string `json:"blockId"`
	Key                string `json:"key"`
	ValueMobile        string `json:"valueMobile"`
	ValueTablet        string `json:"valueTablet"`
	ValueDesktop       string `json:"valueDesktop"`
	Type               string `json:"type"`
	Description        string `json:"description"`
	ValuePicker        string `json:"valuePicker"`
	ValuePickerGroup   string `json:"valuePickerGroup"`
	ValuePickerOptions string `json:"valuePickerOptions"`
	Deprecated         bool   `json:"deprecated"`
	DeprecatedReason   string `json:"deprecatedReason"`
}

type BlockDataJson struct {
	BlockId          string `json:"blockId"`
	Key              string `json:"key"`
	Value            string `json:"value"`
	Type             string `json:"type"`
	Description      string `json:"description"`
	Deprecated       bool   `json:"deprecated"`
	DeprecatedReason string `json:"deprecatedReason"`
}

type BlockSlotJson struct {
	BlockId          string `json:"blockId"`
	Slot             string `json:"slot"`
	Description      string `json:"description"`
	Deprecated       bool   `json:"deprecated"`
	DeprecatedReason string `json:"deprecatedReason"`
}

type ActionJson struct {
	Id       string              `json:"id"`
	FrameId  string              `json:"frameId"`
	Key      string              `json:"key"`
	Event    string              `json:"event"`
	Triggers []ActionTriggerJson `json:"triggers"`
}

type ActionTriggerJson struct {
	Id                          string                `json:"id"`
	ActionId                    string                `json:"actionId"`
	ParentId                    string                `json:"parentId"`
	KeyType                     string                `json:"keyType"`
	Then                        string                `json:"then"`
	Name                        string                `json:"name"`
	IntegrationVersion          int                   `json:"integrationVersion"`
	Properties                  []TriggerPropertyJson `json:"properties"`
	Data                        []TriggerDataJson     `json:"data"`
	IntegrationDeprecated       bool                  `json:"integrationDeprecated"`
	IntegrationDeprecatedReason string                `json:"integrationDeprecatedReason"`
}

type TriggerPropertyJson struct {
	ActionTriggerId    string `json:"actionTriggerId"`
	Key                string `json:"key"`
	Value              string `json:"value"`
	Type               string `json:"type"`
	Description        string `json:"description"`
	ValuePicker        string `json:"valuePicker"`
	ValuePickerGroup   string `json:"valuePickerGroup"`
	ValuePickerOptions string `json:"valuePickerOptions"`
	Deprecated         bool   `json:"deprecated"`
	DeprecatedReason   string `json:"deprecatedReason"`
}

type TriggerDataJson struct {
	ActionTriggerId  string `json:"actionTriggerId"`
	Key              string `json:"key"`
	Value            string `json:"value"`
	Type             string `json:"type"`
	Description      string `json:"description"`
	Deprecated       bool   `json:"deprecated"`
	DeprecatedReason string `json:"deprecatedReason"`
}
