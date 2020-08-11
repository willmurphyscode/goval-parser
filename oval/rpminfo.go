package oval

import (
	"encoding/xml"
)

// RPMInfoTest : >tests>rpminfo_test
type RPMInfoTest struct {
	XMLName    xml.Name    `xml:"rpminfo_test"`
	ID         string      `xml:"id,attr"`
	Comment    string      `xml:"comment,attr"`
	Check      string      `xml:"check,attr"`
	Version    int         `xml:"version,attr"`
	ObjectRefs []ObjectRef `xml:"object"`
	StateRefs  []StateRef  `xml:"state"`
}

// RPMVerifyFileTest : >tests>rpmverifyfile_test
type RPMVerifyFileTest struct {
	XMLName    xml.Name    `xml:"rpmverifyfile_test"`
	ID         string      `xml:"id,attr"`
	Comment    string      `xml:"comment,attr"`
	Check      string      `xml:"check,attr"`
	Version    int         `xml:"version,attr"`
	ObjectRefs []ObjectRef `xml:"object"`
	StateRefs  []StateRef  `xml:"state"`
}

// RPMInfoObject : >objects>RPMInfo_object
type RPMInfoObject struct {
	XMLName xml.Name `xml:"rpminfo_object"`
	ID      string   `xml:"id,attr"`
	Version int      `xml:"version,attr"`
	Name    string   `xml:"name"`
}

// RPMInfoState : >states>rpminfo_state
type RPMInfoState struct {
	XMLName        xml.Name           `xml:"rpminfo_state"`
	ID             string             `xml:"id,attr"`
	Version        int                `xml:"version,attr"`
	Arch           *Arch              `xml:"arch"`
	Epoch          *Epoch             `xml:"epoch"`
	Release        *Release           `xml:"release"`
	RPMVersion     *Version           `xml:"version"`
	EVR            *EVR               `xml:"evr"`
	SignatureKeyID *RPMSignatureKeyID `xml:"signature_keyid"`
}

type RPMSignatureKeyID struct {
	XMLName   xml.Name  `xml:"signature_keyid"`
	Operation Operation `xml:"operation,attr"`
	Body      string    `xml:",chardata"`
}
