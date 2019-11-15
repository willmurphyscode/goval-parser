package oval

import "encoding/xml"

// TextfileContent54Test : >tests>textfilecontent54_test
type TextfileContent54Test struct {
	XMLName       xml.Name    `xml:"textfilecontent54_test"`
	ID            string      `xml:"id,attr"`
	StateOperator string      `xml:"state_operator,attr"`
	ObjectRefs    []ObjectRef `xml:"object"`
	StateRefs     []StateRef  `xml:"state"`
	Comment       string      `xml:"comment,attr"`
}
