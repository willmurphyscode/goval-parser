package oval

import "encoding/xml"

// UnameTest : >tests>uname_test
type UnameTest struct {
	XMLName    xml.Name    `xml:"uname_test"`
	ID         string      `xml:"id,attr"`
	Comment    string      `xml:"comment,attr"`
	Check      string      `xml:"check,attr"`
	Version    int         `xml:"version,attr"`
	ObjectRefs []ObjectRef `xml:"object"`
	StateRefs  []StateRef  `xml:"state"`
}
