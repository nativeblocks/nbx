package model

import "encoding/xml"

type XMLFrame struct {
	XMLName   xml.Name      `xml:"frame"`
	Name      string        `xml:"name,attr"`
	Route     string        `xml:"route,attr"`
	Type      string        `xml:"type,attr"`
	Variables []XMLVariable `xml:"var"`
	Blocks    []XMLBlock    `xml:"block"`
}

type XMLVariable struct {
	Key   string `xml:"key,attr"`
	Type  string `xml:"type,attr"`
	Value string `xml:"value,attr"`
}

type XMLBlock struct {
	KeyType    string        `xml:"keyType,attr"`
	Key        string        `xml:"key,attr"`
	Visibility string        `xml:"visibility,attr"`
	Version    int           `xml:"version,attr"`
	Properties []XMLProperty `xml:"prop"`
	Data       []XMLData     `xml:"data"`
	Slots      []XMLSlot     `xml:"slot"`
	Actions    []XMLAction   `xml:"action"`
}

type XMLProperty struct {
	Key     string `xml:"key,attr"`
	Value   string `xml:"value,attr"`
	Mobile  string `xml:"mobile,attr"`
	Tablet  string `xml:"tablet,attr"`
	Desktop string `xml:"desktop,attr"`
}

type XMLData struct {
	Key   string `xml:"key,attr"`
	Value string `xml:"value,attr"`
}

type XMLSlot struct {
	Name   string     `xml:"name,attr"`
	Blocks []XMLBlock `xml:"block"`
}

type XMLAction struct {
	Event    string       `xml:"event,attr"`
	Triggers []XMLTrigger `xml:"trigger"`
}

type XMLTrigger struct {
	KeyType    string        `xml:"keyType,attr"`
	Name       string        `xml:"name,attr"`
	Version    int           `xml:"version,attr"`
	Properties []XMLProperty `xml:"prop"`
	Data       []XMLData     `xml:"data"`
	Then       []XMLThen     `xml:"then"`
}

type XMLThen struct {
	Value    string       `xml:"value,attr"`
	Triggers []XMLTrigger `xml:"trigger"`
}
