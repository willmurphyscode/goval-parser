package oval

import (
	"encoding/xml"
	"os"
	"strings"
	"testing"
)

func TestLookupRPMIntoObject(t *testing.T) {
	var tt = []struct{ Ref, Name string }{
		{"oval:com.redhat.rhba:obj:20070026001", "htdig"},        // This should be the first object.
		{"oval:com.redhat.rhsa:obj:20091206001", "libxml"},       // random one
		{"oval:com.redhat.rhsa:obj:20100720002", "mikmod-devel"}, // This should be the last object.
	}
	f, err := os.Open("../testdata/Red_Hat_Enterprise_Linux_3.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	var root Root
	if err := xml.NewDecoder(f).Decode(&root); err != nil {
		t.Fatal(err)
	}
	for _, tc := range tt {
		kind, i, err := root.Objects.Lookup(tc.Ref)
		if err != nil {
			t.Error(err)
		}
		if got, want := kind, "rpminfo_object"; got != want {
			t.Errorf("got: %q, want %q", got, want)
		}
		obj := &root.Objects.RPMInfoObjects[i]
		t.Logf("%s: %s (%#+v)", tc.Ref, obj.Name, obj)
		if got, want := obj.Name, tc.Name; got != want {
			t.Fatalf("got: %q, want: %q", got, want)
		}
	}
}

func TestLookupRPMTest(t *testing.T) {
	var tt = []struct{ Kind, Ref, Name string }{
		{"rpmverifyfile_test", "oval:com.redhat.rhsa:tst:20190966006", "Red Hat Enterprise Linux must be installed"},
		{"uname_test", "oval:com.redhat.rhsa:tst:20191167051", "kernel earlier than 0:4.18.0-80.1.2.el8_0 is currently running"},
		{"textfilecontent54_test", "oval:com.redhat.rhsa:tst:20191167052", "kernel earlier than 0:4.18.0-80.1.2.el8_0 is set to boot up on next boot"},
	}
	f, err := os.Open("../testdata/com.redhat.rhsa-RHEL8.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	var root Root
	if err := xml.NewDecoder(f).Decode(&root); err != nil {
		t.Fatal(err)
	}
	for _, tc := range tt {
		kind, i, err := root.Tests.Lookup(tc.Ref)
		if err != nil {
			t.Error(err)
		}
		if got, want := kind, tc.Kind; got != want {
			t.Errorf("got: %q, want %q", got, want)
		}
		var name string
		var test Test
		switch kind {
		case "rpmverifyfile_test":
			obj := &root.Tests.RPMVerifyFileTests[i]
			name = obj.Comment
			test = obj
			t.Logf("%s: %s (%#+v)", tc.Ref, obj.Comment, obj)
		case "uname_test":
			obj := &root.Tests.UnameTests[i]
			name = obj.Comment
			test = obj
			t.Logf("%s: %s (%#+v)", tc.Ref, obj.Comment, obj)
		case "textfilecontent54_test":
			obj := &root.Tests.TextfileContent54Tests[i]
			name = obj.Comment
			test = obj
			t.Logf("%s: %s (%#+v)", tc.Ref, obj.Comment, obj)
		default:
			t.Fatalf("unknown test kind %q", kind)
		}
		var ss, os strings.Builder
		ss.WriteByte('[')
		for i, ref := range test.StateRef() {
			if i != 0 {
				ss.WriteString(", ")
			}
			ss.WriteString(ref.StateRef)
		}
		ss.WriteByte(']')
		os.WriteByte('[')
		for i, ref := range test.ObjectRef() {
			if i != 0 {
				os.WriteString(", ")
			}
			os.WriteString(ref.ObjectRef)
		}
		os.WriteByte(']')
		t.Logf("states: %s; objects: %s", ss.String(), os.String())
		if got, want := name, tc.Name; got != want {
			t.Fatalf("got: %q, want: %q", got, want)
		}
	}
}
