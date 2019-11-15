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
	Arch           *RPMArch           `xml:"arch"`
	Epoch          *RPMEpoch          `xml:"epoch"`
	Release        *RPMRelease        `xml:"release"`
	RPMVersion     *RPMVersion        `xml:"version"`
	EVR            *RPMEVR            `xml:"evr"`
	SignatureKeyID *RPMSignatureKeyID `xml:"signature_keyid"`
}

type RPMArch struct {
	XMLName   xml.Name  `xml:"arch"`
	Operation Operation `xml:"operation,attr"`
	Body      string    `xml:",chardata"`
}
type RPMEpoch struct {
	XMLName   xml.Name  `xml:"epoch"`
	Operation Operation `xml:"operation,attr"`
	Body      string    `xml:",chardata"`
}
type RPMRelease struct {
	XMLName   xml.Name  `xml:"release"`
	Operation Operation `xml:"operation,attr"`
	Body      string    `xml:",chardata"`
}
type RPMVersion struct {
	XMLName   xml.Name  `xml:"version"`
	Operation Operation `xml:"operation,attr"`
	Body      string    `xml:",chardata"`
}
type RPMEVR struct {
	XMLName   xml.Name  `xml:"evr"`
	Operation Operation `xml:"operation,attr"`
	Body      string    `xml:",chardata"`
}
type RPMSignatureKeyID struct {
	XMLName   xml.Name  `xml:"signature_keyid"`
	Operation Operation `xml:"operation,attr"`
	Body      string    `xml:",chardata"`
}
