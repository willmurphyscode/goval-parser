package oval

import "encoding/xml"

// TextfileContent54Test : >tests>textfilecontent54_test
type TextfileContent54Test struct {
	XMLName       xml.Name `xml:"textfilecontent54_test"`
	ID            string   `xml:"id,attr"`
	StateOperator string   `xml:"state_operator,attr"`
	Comment       string   `xml:"comment,attr"`
	testRef
}

var _ Test = (*TextfileContent54Test)(nil)

// TextfileContent54Object : >tests>textfilecontent54_object
type TextfileContent54Object struct {
	XMLName  xml.Name   `xml:"textfilecontent54_object"`
	ID       string     `xml:"id,attr"`
	Version  int        `xml:"version,attr"`
	Filepath string     `xml:"filepath"`
	Pattern  string     `xml:"pattern"`
	Instance TCInstance `xml:"instance"`
}

// TextfileContent54State : >states>textfilecontent54_state
type TextfileContent54State struct {
	XMLName xml.Name                   `xml:"textfilecontent54_state"`
	ID      string                     `xml:"id,attr"`
	Version string                     `xml:"version,attr"`
	Text    TextfileContent54StateText `xml:"text"`
}

type TextfileContent54StateText struct {
	XMLName   xml.Name `xml:"text"`
	Operation string   `xml:"operation,attr"`
}

type TCInstance struct {
	XMLName  xml.Name `xml:"instance"`
	Instance string   `xml:",chardata"`
	Kind     string   `xml:"datatype,attr"`
}
