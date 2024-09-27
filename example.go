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
	var tests []oval.RPMInfoTest
	var states []oval.RPMInfoState
	var objects []oval.RPMInfoObject
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
			tests = append(tests, t)
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
			states = append(states, state)
		}
	}
	for _, obj := range root.Objects.RPMInfoObjects {
		if _, ok := objRefs[obj.ID]; ok {
			objects = append(objects, obj)
		}
	}
	newDocument := &oval.Root{
		XMLName:   root.XMLName,
		Generator: root.Generator,
		Definitions: oval.Definitions{
			Definitions: keep,
		},
		Tests: oval.Tests{
			RPMInfoTests: tests,
		},
		States: oval.States{
			RPMInfoStates: states,
		},
		Objects: oval.Objects{
			RPMInfoObjects: objects,
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
		_, err := w.Write([]byte("  " + indent + describeTest(here.TestRef, root) + "\n"))
		if err != nil {
			panic(err)
		}
	}
}

func describeTest(testRef string, root *oval.Root) string {
	var test oval.RPMInfoTest
	var state oval.RPMInfoState
	var object oval.RPMInfoObject
	for _, t := range root.Tests.RPMInfoTests {
		if t.ID == testRef {
			test = t
		}
	}
	if len(test.StateRefs) == 0 {
		panic("malformed test")
	}
	if len(test.ObjectRefs) == 0 {
		panic("malformed test")
	}
	for _, s := range root.States.RPMInfoStates {
		if s.ID == test.StateRefs[0].StateRef {
			state = s
		}
	}
	for _, o := range root.Objects.RPMInfoObjects {
		if o.ID == test.ObjectRefs[0].ObjectRef {
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
			if _, ok := ids[def.Title]; !ok {
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
