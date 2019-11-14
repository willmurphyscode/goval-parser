package oval

import (
	"encoding/xml"
	"fmt"
	"sync"
)

// ErrNotFound is returned by Lookup methods when the specified identifier is
// not found.
type ErrNotFound string

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("oval: identifier %q not found", string(e))
}

// Root : root object
type Root struct {
	XMLName     xml.Name    `xml:"oval_definitions"`
	Generator   Generator   `xml:"generator"`
	Definitions Definitions `xml:"definitions"`
	Tests       Tests       `xml:"tests"`
	Objects     Objects     `xml:"objects"`
	States      States      `xml:"states"`
}

// Generator : >generator
type Generator struct {
	XMLName        xml.Name `xml:"generator"`
	ProductName    string   `xml:"product_name"`
	ProductVersion string   `xml:"product_version"`
	SchemaVersion  string   `xml:"schema_version"`
	Timestamp      string   `xml:"timestamp"`
}

// Definitions : >definitions
type Definitions struct {
	XMLName     xml.Name     `xml:"definitions"`
	Definitions []Definition `xml:"definition"`
}

// Definition : >definitions>definition
type Definition struct {
	XMLName     xml.Name    `xml:"definition"`
	ID          string      `xml:"id,attr"`
	Class       string      `xml:"class,attr"`
	Title       string      `xml:"metadata>title"`
	Affecteds   []Affected  `xml:"metadata>affected"`
	References  []Reference `xml:"metadata>reference"`
	Description string      `xml:"metadata>description"`
	Advisory    Advisory    `xml:"metadata>advisory"` // RedHat, Oracle, Ubuntu
	Debian      Debian      `xml:"metadata>debian"`   // Debian
	Criteria    Criteria    `xml:"criteria"`
}

// Criteria : >definitions>definition>criteria
type Criteria struct {
	XMLName    xml.Name    `xml:"criteria"`
	Operator   string      `xml:"operator,attr"`
	Criterias  []Criteria  `xml:"criteria"`
	Criterions []Criterion `xml:"criterion"`
}

// Criterion : >definitions>definition>criteria>*>criterion
type Criterion struct {
	XMLName xml.Name `xml:"criterion"`
	Negate  bool     `xml:"negate,attr"`
	TestRef string   `xml:"test_ref,attr"`
	Comment string   `xml:"comment,attr"`
}

// Affected : >definitions>definition>metadata>affected
type Affected struct {
	XMLName   xml.Name `xml:"affected"`
	Family    string   `xml:"family,attr"`
	Platforms []string `xml:"platform"`
}

// Reference : >definitions>definition>metadata>reference
type Reference struct {
	XMLName xml.Name `xml:"reference"`
	Source  string   `xml:"source,attr"`
	RefID   string   `xml:"ref_id,attr"`
	RefURL  string   `xml:"ref_url,attr"`
}

// Advisory : >definitions>definition>metadata>advisory
// RedHat and Ubuntu OVAL
type Advisory struct {
	XMLName         xml.Name   `xml:"advisory"`
	Severity        string     `xml:"severity"`
	Cves            []Cve      `xml:"cve"`
	Bugzillas       []Bugzilla `xml:"bugzilla"`
	AffectedCPEList []string   `xml:"affected_cpe_list>cpe"`
	Refs            []Ref      `xml:"ref"` // Ubuntu Only
	Bugs            []Bug      `xml:"bug"` // Ubuntu Only
	Issued          struct {
		Date string `xml:"date,attr"`
	} `xml:"issued"`
	Updated struct {
		Date string `xml:"date,attr"`
	} `xml:"updated"`
}

// Ref : >definitions>definition>metadata>advisory>ref
// Ubuntu OVAL
type Ref struct {
	XMLName xml.Name `xml:"ref"`
	URL     string   `xml:",chardata"`
}

// Bug : >definitions>definition>metadata>advisory>bug
// Ubuntu OVAL
type Bug struct {
	XMLName xml.Name `xml:"bug"`
	URL     string   `xml:",chardata"`
}

// Cve : >definitions>definition>metadata>advisory>cve
// RedHat OVAL
type Cve struct {
	XMLName xml.Name `xml:"cve"`
	CveID   string   `xml:",chardata"`
	Cvss2   string   `xml:"cvss2,attr"`
	Cvss3   string   `xml:"cvss3,attr"`
	Cwe     string   `xml:"cwe,attr"`
	Impact  string   `xml:"impact,attr"`
	Href    string   `xml:"href,attr"`
	Public  string   `xml:"public,attr"`
}

// Bugzilla : >definitions>definition>metadata>advisory>bugzilla
// RedHat OVAL
type Bugzilla struct {
	XMLName xml.Name `xml:"bugzilla"`
	ID      string   `xml:"id,attr"`
	URL     string   `xml:"href,attr"`
	Title   string   `xml:",chardata"`
}

// Debian : >definitions>definition>metadata>debian
type Debian struct {
	XMLName  xml.Name `xml:"debian"`
	MoreInfo string   `xml:"moreinfo"`
	Date     string   `xml:"date"`
}

// Tests : >tests
type Tests struct {
	once           sync.Once
	XMLName        xml.Name        `xml:"tests"`
	LineTests      []LineTest      `xml:"line_test"`
	Version55Tests []Version55Test `xml:"version55_test"`
	RPMInfoTests   []RPMInfoTest   `xml:"rpminfo_test"`
	lineMemo       map[string]int
	version55Memo  map[string]int
	rpminfoMemo    map[string]int
}

// ObjectRef : >tests>line_test>object-object_ref
//           : >tests>version55_test>object-object_ref
type ObjectRef struct {
	XMLName   xml.Name `xml:"object"`
	ObjectRef string   `xml:"object_ref,attr"`
}

// StateRef : >tests>line_test>state-state_ref
//          : >tests>version55_test>state-state_ref
type StateRef struct {
	XMLName  xml.Name `xml:"state"`
	StateRef string   `xml:"state_ref,attr"`
}

// Objects : >objects
type Objects struct {
	once             sync.Once
	XMLName          xml.Name          `xml:"objects"`
	LineObjects      []LineObject      `xml:"line_object"`
	Version55Objects []Version55Object `xml:"version55_object"`
	RPMInfoObjects   []RPMInfoObject   `xml:"rpminfo_object"`
	lineMemo         map[string]int
	version55Memo    map[string]int
	rpminfoMemo      map[string]int
}

// States : >states
type States struct {
	once            sync.Once
	XMLName         xml.Name         `xml:"states"`
	LineStates      []LineState      `xml:"line_state"`
	Version55States []Version55State `xml:"version55_state"`
	RPMInfoStates   []RPMInfoState   `xml:"rpminfo_state"`
	lineMemo        map[string]int
	version55Memo   map[string]int
	rpminfoMemo     map[string]int
}
