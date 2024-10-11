package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/quay/goval-parser/oval"
)

func main() {
	// if len(os.Args) < 2 {
	// 	fmt.Printf("usage: %s file1 cve-id-1 ...\n", os.Args[0])
	// 	return
	// }
	criteriaRegex := flag.String("criteria-regex", "", "Regex pattern for criteria, pass multiple like /foo/,/bar/ for AND, use regex | for OR")
	maxCriteriaCount := flag.Int("max-criteria-count", 0, "Maximum number of criteria")
	ids := flag.String("ids", "", "comma-separate list of IDs")
	max := flag.Int("max", 0, "most total results to return. If more are found, random subset of this size will be returned")
	human := flag.Bool("h", false, "output results in human readable format instead of XML")

	// Parse the flags
	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		// TODO: halt
		panic("exactly 1 arg required (file to scan). For stdin pass `-'")
	}
	fmt.Println("args: crit regex: ", *criteriaRegex, "crit count: ", *maxCriteriaCount, "ids: ", *ids, "max defs: ", *max)

	// Validate criteriaRegex if provided
	var res []*regexp.Regexp
	var err error
	if *criteriaRegex != "" {
		res = criteriaRegexFromFlagValue(*criteriaRegex)
	}
	// root, err := readOval(os.Args[1])
	var r io.Reader
	if args[0] == "-" {
		r = os.Stdin
	} else {
		f, err := os.Open(args[0])
		if err != nil {
			panic(err)
		}
		r = f
	}
	root, err := read(r)
	if err != nil {
		panic(err)
	}
	allowedIDs := make(map[string]struct{})
	if len(*ids) > 0 {
		idsStrs := strings.Split(*ids, ",")
		for _, s := range idsStrs {
			allowedIDs[strings.TrimSpace(s)] = struct{}{}
		}
	}
	keep := filterDefinitions(root, allowedIDs, *maxCriteriaCount, res)
	// for _, f := range os.Args[2:len(os.Args)] {
	// 	for _, d := range root.Definitions.Definitions {
	// 		if d.Title == f {
	// 			keep = append(keep, d)
	// 		}
	// 	}
	// }
	if *max > 0 && len(keep) > *max {
		keep = keep[0:*max]
	}
	var rpmInfoTests []oval.RPMInfoTest
	var rpmFileTests []oval.RPMVerifyFileTest
	var textfileContent54Tests []oval.TextfileContent54Test
	var rpmInfoStates []oval.RPMInfoState
	var rpmFilestates []oval.RPMVerifyFileState
	var textFileStates []oval.TextfileContent54State
	var dpkgInfoStates []oval.DpkgInfoState
	var version55States []oval.Version55State
	var lineStates []oval.LineState
	var rpmInfoObjects []oval.RPMInfoObject
	var rpmFileObjects []oval.RPMVerifyFileObject
	var textfileObjects []oval.TextfileContent54Object
	testRefs := make(map[string]struct{})
	for _, k := range keep {
		// TODO: walk criteria and figure out what we need
		accumulateTestRefs(testRefs, k.Criteria)
	}
	// TODO: write output
	stateRefs := make(map[string]struct{})
	objRefs := make(map[string]struct{})
	for _, t := range root.Tests.RPMInfoTests {
		if _, ok := testRefs[t.ID]; ok {
			rpmInfoTests = append(rpmInfoTests, t)
			for _, sr := range t.StateRef() {
				stateRefs[sr.StateRef] = struct{}{}
			}
			for _, objRef := range t.ObjectRef() {
				objRefs[objRef.ObjectRef] = struct{}{}
			}
		}
	}
	// TODO: there should be filter methods to reduce duplication here
	for _, t := range root.Tests.RPMVerifyFileTests {
		if _, ok := testRefs[t.ID]; ok {
			rpmFileTests = append(rpmFileTests, t)
			for _, sr := range t.StateRef() {
				stateRefs[sr.StateRef] = struct{}{}
			}
			for _, objRef := range t.ObjectRef() {
				objRefs[objRef.ObjectRef] = struct{}{}
			}
		}
	}
	for _, t := range root.Tests.TextfileContent54Tests {
		if _, ok := testRefs[t.ID]; ok {
			textfileContent54Tests = append(textfileContent54Tests, t)
			for _, sr := range t.StateRef() {
				stateRefs[sr.StateRef] = struct{}{}
			}
			for _, objRef := range t.ObjectRef() {
				objRefs[objRef.ObjectRef] = struct{}{}
			}
		}
	}
	for _, state := range root.States.RPMInfoStates {
		if _, ok := stateRefs[state.ID]; ok {
			rpmInfoStates = append(rpmInfoStates, state)
		}
	}
	for _, state := range root.States.DpkgInfoStates {
		if _, ok := stateRefs[state.ID]; ok {
			dpkgInfoStates = append(dpkgInfoStates, state)
		}
	}
	for _, state := range root.States.Version55States {
		if _, ok := stateRefs[state.ID]; ok {
			version55States = append(version55States, state)
		}
	}
	for _, state := range root.States.LineStates {
		if _, ok := stateRefs[state.ID]; ok {
			lineStates = append(lineStates, state)
		}
	}
	for _, state := range root.States.RPMVerifyFileStates {
		if _, ok := stateRefs[state.ID]; ok {
			rpmFilestates = append(rpmFilestates, state)
		}
	}
	for _, state := range root.States.TextfileContent54States {
		if _, ok := stateRefs[state.ID]; ok {
			textFileStates = append(textFileStates, state)
		}
	}
	for _, obj := range root.Objects.RPMInfoObjects {
		if _, ok := objRefs[obj.ID]; ok {
			rpmInfoObjects = append(rpmInfoObjects, obj)
		}
	}
	// TODO:
	//for _, obj := range root.Objects.DpkgInfoObjects {
	//
	//}
	//for _, obj := range root.Objects.Version55Objects {
	//
	//}
	//for _, obj := range root.Objects.LineObjects {
	//
	//}
	for _, obj := range root.Objects.TextfileContent54Objects {
		if _, ok := objRefs[obj.ID]; ok {
			textfileObjects = append(textfileObjects, obj)
		}
	}
	for _, obj := range root.Objects.RPMVerifyFileObjects {
		if _, ok := objRefs[obj.ID]; ok {
			rpmFileObjects = append(rpmFileObjects, obj)
		}
	}
	newDocument := &oval.Root{
		XMLName:   root.XMLName,
		Generator: root.Generator,
		Definitions: oval.Definitions{
			Definitions: keep,
		},
		Tests: oval.Tests{
			RPMInfoTests:           rpmInfoTests,
			RPMVerifyFileTests:     rpmFileTests,
			TextfileContent54Tests: textfileContent54Tests,
		},
		States: oval.States{
			RPMInfoStates:           rpmInfoStates,
			RPMVerifyFileStates:     rpmFilestates,
			DpkgInfoStates:          dpkgInfoStates,
			Version55States:         version55States,
			LineStates:              lineStates,
			TextfileContent54States: textFileStates,
		},
		Objects: oval.Objects{
			RPMInfoObjects:           rpmInfoObjects,
			RPMVerifyFileObjects:     rpmFileObjects,
			TextfileContent54Objects: textfileObjects,
		},
	}

	if *human {
		for _, d := range newDocument.Definitions.Definitions {
			fmt.Println(d.Title)
			printCriteria(os.Stdout, "", d.Criteria, newDocument)
		}
		return
	}

	enc := xml.NewEncoder(os.Stdout)
	enc.Indent("", "  ")
	err = enc.Encode(newDocument)
	if err != nil {
		panic(err)
	}
}

func printCriteria(w io.Writer, indent string, c oval.Criteria, root *oval.Root) {
	_, err := w.Write([]byte(indent + c.Operator + "\n"))
	if err != nil {
		panic(err)
	}
	for _, inner := range c.Criterias {
		printCriteria(w, "  "+indent, inner, root)
	}
	for _, here := range c.Criterions {
		_, err := w.Write([]byte("  " + indent + describeTestGood(here.TestRef, root) + "\n"))
		if err != nil {
			panic(err)
		}
	}
}

func describeTest(testRef string, root *oval.Root) string {
	var test oval.Test
	var state oval.RPMInfoState
	var object oval.RPMInfoObject
	for _, t := range root.Tests.RPMInfoTests {
		if t.ID == testRef {
			test = &t
			break
		}
	}
	for _, t := range root.Tests.RPMVerifyFileTests {
		if t.ID == testRef {
			test = &t
			break
		}
	}
	if test == nil {
		panic("test not found: " + testRef)
	}
	if len(test.StateRef()) == 0 {
		panic("malformed test: no states")
	}
	if len(test.ObjectRef()) == 0 {
		panic("malformed test: no objects")
	}
	sr := test.StateRef()[0]
	for _, s := range root.States.RPMInfoStates {
		if s.ID == sr.StateRef {
			state = s
			break
		}
	}
	obj_ref := test.ObjectRef()[0]
	for _, o := range root.Objects.RPMInfoObjects {
		if o.ID == obj_ref.ObjectRef {
			object = o
		}
	}
	if state.EVR != nil {
		return fmt.Sprintf("%s %s %s", object.Name, state.EVR.Operation, state.EVR.Body)
	}
	if state.RPMVersion != nil {
		return fmt.Sprintf("%s %s %s", object.Name, state.RPMVersion.Operation, state.RPMVersion.Body)
	}
	panic("surprise! not an EVR or an RPM version.")
}

func describeTestGood(testRef string, root *oval.Root) string {
	var test oval.Test
	testKind, testIx, err := root.Tests.Lookup(testRef)
	if err != nil {
		panic(err)
	}
	switch testKind {
	case "rpmverifyfile_test":
		test = &root.Tests.RPMVerifyFileTests[testIx]
	case "rpminfo_test":
		test = &root.Tests.RPMInfoTests[testIx]
	case "textfilecontent54_test":
		test = &root.Tests.TextfileContent54Tests[testIx]
	default:
		panic("not implemented for test kind: " + testKind)
	}
	var srKind, objKind string
	var srIx, objIx int
	for _, sr := range test.StateRef() {
		var err error
		srKind, srIx, err = root.States.Lookup(sr.StateRef)
		if err != nil {
			panic("state lookup failed: " + err.Error())
		}
	}
	for _, objRef := range test.ObjectRef() {
		var err error
		objKind, objIx, err = root.Objects.Lookup(objRef.ObjectRef)
		if err != nil {
			panic(err)
		}
	}
	var stateStr, objectStr string

	switch srKind {
	case "rpminfo_state":
		stateStr = formatRpmInfoState(root.States.RPMInfoStates[srIx])
	case "rpmverifyfile_state":
		stateStr = formatRpmVerifyFileState(root.States.RPMVerifyFileStates[srIx])
	case "textfilecontent54_state":
		stateStr = formatTextFileState(root.States.TextfileContent54States[srIx])
	default:
		panic("not implemented for: " + srKind)
	}

	switch objKind {
	case "rpminfo_object":
		objectStr = formatRpmInfoObj(root.Objects.RPMInfoObjects[objIx])
	case "rpmverifyfile_object":
		objectStr = formatRpmVerifyFileObject(root.Objects.RPMVerifyFileObjects[objIx])
	case "textfilecontent54_object":
		objectStr = formatTextFileObject(root.Objects.TextfileContent54Objects[objIx])
	default:
		panic("not implemented for: " + objKind)
	}
	result := fmt.Sprintf(stateStr, objectStr)
	return result
}

func formatTextFileObject(object oval.TextfileContent54Object) string {
	return fmt.Sprintf("%s has %s", object.Filepath, object.Pattern)
}

func formatTextFileState(state oval.TextfileContent54State) string {
	return fmt.Sprintf("%%s %s ", state.Text.Operation)
}

func formatRpmVerifyFileObject(object oval.RPMVerifyFileObject) string {
	return fmt.Sprintf("For path %s", object.Filepath)
}

func formatRpmVerifyFileState(state oval.RPMVerifyFileState) string {
	//if state.Version == nil {
	//	panic("cannot print nil version state with ID: " + state.ID)
	//}
	//return fmt.Sprintf("%%s %s %s", state.Version.Operation, state.Version.Body)
	return fmt.Sprintf("%%s %s %s", state.Name.Operation, state.Name.Body)
}

func formatRpmInfoObj(object oval.RPMInfoObject) string {
	return object.Name
}

func formatRpmInfoState(state oval.RPMInfoState) string {
	if state.EVR != nil {
		return fmt.Sprintf("%%s %s %s", state.EVR.Operation, state.EVR.Body)
	}
	if state.RPMVersion != nil {
		return fmt.Sprintf("%%s %s %s", state.EVR.Operation, state.EVR.Body)
	}

	if state.SignatureKeyID != nil {
		return fmt.Sprintf("%%s signature %s %s", state.SignatureKeyID.Operation, state.SignatureKeyID.Body)
	}
	panic(fmt.Sprintf("do not know how to format %+v", state))
}

func criteriaRegexFromFlagValue(v string) []*regexp.Regexp {
	var result []*regexp.Regexp
	if !(strings.HasPrefix(v, "/") && strings.HasSuffix(v, "/")) {
		return []*regexp.Regexp{
			regexp.MustCompile(v),
		}
	}
	all := strings.Trim(v, "/")
	parts := strings.Split(all, "/,/")
	for _, p := range parts {
		result = append(result, regexp.MustCompile(p))
	}

	return result
}

func filterDefinitions(root *oval.Root, ids map[string]struct{}, maxCriteriaCount int, criteriaRegexes []*regexp.Regexp) []oval.Definition {
	var filterIDs, filterCount bool
	if len(ids) > 0 {
		filterIDs = true
	}
	if maxCriteriaCount > 0 {
		filterCount = true
	}
	var result []oval.Definition
	for _, def := range root.Definitions.Definitions {
		if filterIDs {
			idMatch := false
			if _, ok := ids[def.Title]; ok {
				idMatch = true
			}
			rhelId := rhelAdvisoryFromTitle(def.Title)
			if _, ok := ids[rhelId]; ok {
				idMatch = true
			}
			if strings.Contains(def.Title, ":") {
				prefix := strings.Split(def.Title, ":")[0]
				if _, ok := ids[prefix]; ok {
					idMatch = true
				}
			}
			if !idMatch {
				continue
			}
		}
		if filterCount {
			tests := make(map[string]struct{})
			accumulateTestRefs(tests, def.Criteria)
			if len(tests) > maxCriteriaCount {
				continue
			}
		}
		if len(criteriaRegexes) > 0 {
			regexMatch := true
			for _, c := range criteriaRegexes {
				if !anyCriteriaMatch(c, def.Criteria) {
					regexMatch = false
					break
				}
			}
			if !regexMatch {
				continue
			}
		}

		result = append(result, def)
	}

	return result
}

func rhelAdvisoryFromTitle(title string) string {
	if !strings.Contains(title, ":") {
		return ""
	}
	parts := strings.Split(title, ":")
	if len(parts) > 2 {
		return fmt.Sprintf("%s-%s", parts[0], parts[1])
	}
	return ""
}

func accumulateTestRefs(acc map[string]struct{}, criteria oval.Criteria) {
	for _, rec := range criteria.Criterias {
		accumulateTestRefs(acc, rec)
	}
	for _, here := range criteria.Criterions {
		acc[here.TestRef] = struct{}{}
	}
}

func anyCriteriaMatch(re *regexp.Regexp, criteria oval.Criteria) bool {
	for _, here := range criteria.Criterions {
		if re.MatchString(here.Comment) {
			return true
		}
	}
	for _, rec := range criteria.Criterias {
		if anyCriteriaMatch(re, rec) {
			return true
		}
	}
	return false
}

func read(r io.Reader) (*oval.Root, error) {
	root := &oval.Root{}
	err := xml.NewDecoder(r).Decode(root)
	return root, err
}
