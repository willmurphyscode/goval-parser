package oval

import (
	"encoding/xml"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestAdvisoryDates(t *testing.T) {
	var tt = []struct{ Updated, Issued time.Time }{
		{time.Date(2019, time.Month(5), 7, 0, 0, 0, 0, time.UTC), time.Date(2019, time.Month(5), 7, 0, 0, 0, 0, time.UTC)},
		{time.Date(2019, time.Month(5), 7, 0, 0, 0, 0, time.UTC), time.Date(2019, time.Month(5), 7, 0, 0, 0, 0, time.UTC)},
		{time.Date(2019, time.Month(5), 7, 0, 0, 0, 0, time.UTC), time.Date(2019, time.Month(5), 7, 0, 0, 0, 0, time.UTC)},
		{time.Date(2019, time.Month(5), 7, 0, 0, 0, 0, time.UTC), time.Date(2019, time.Month(5), 7, 0, 0, 0, 0, time.UTC)},
		{time.Date(2019, time.Month(5), 7, 0, 0, 0, 0, time.UTC), time.Date(2019, time.Month(5), 7, 0, 0, 0, 0, time.UTC)},
		{time.Date(2019, time.Month(5), 7, 0, 0, 0, 0, time.UTC), time.Date(2019, time.Month(5), 7, 0, 0, 0, 0, time.UTC)},
		{time.Date(2019, time.Month(5), 7, 0, 0, 0, 0, time.UTC), time.Date(2019, time.Month(5), 7, 0, 0, 0, 0, time.UTC)},
		{time.Date(2019, time.Month(5), 7, 0, 0, 0, 0, time.UTC), time.Date(2019, time.Month(5), 7, 0, 0, 0, 0, time.UTC)},
		{time.Date(2019, time.Month(5), 7, 0, 0, 0, 0, time.UTC), time.Date(2019, time.Month(5), 7, 0, 0, 0, 0, time.UTC)},
		{time.Date(2019, time.Month(5), 7, 0, 0, 0, 0, time.UTC), time.Date(2019, time.Month(5), 7, 0, 0, 0, 0, time.UTC)},
		{time.Date(2019, time.Month(5), 7, 0, 0, 0, 0, time.UTC), time.Date(2019, time.Month(5), 7, 0, 0, 0, 0, time.UTC)},
		{time.Date(2019, time.Month(5), 7, 0, 0, 0, 0, time.UTC), time.Date(2019, time.Month(5), 7, 0, 0, 0, 0, time.UTC)},
		//
		{time.Date(2019, time.Month(5), 13, 0, 0, 0, 0, time.UTC), time.Date(2019, time.Month(5), 13, 0, 0, 0, 0, time.UTC)},
		{time.Date(2019, time.Month(5), 13, 0, 0, 0, 0, time.UTC), time.Date(2019, time.Month(5), 13, 0, 0, 0, 0, time.UTC)},
		{time.Date(2019, time.Month(5), 13, 0, 0, 0, 0, time.UTC), time.Date(2019, time.Month(5), 13, 0, 0, 0, 0, time.UTC)},
		{time.Date(2019, time.Month(5), 13, 0, 0, 0, 0, time.UTC), time.Date(2019, time.Month(5), 13, 0, 0, 0, 0, time.UTC)},
		{time.Date(2019, time.Month(5), 13, 0, 0, 0, 0, time.UTC), time.Date(2019, time.Month(5), 13, 0, 0, 0, 0, time.UTC)},
		{time.Date(2019, time.Month(5), 13, 0, 0, 0, 0, time.UTC), time.Date(2019, time.Month(5), 13, 0, 0, 0, 0, time.UTC)},
		//
		{time.Date(2019, time.Month(5), 14, 0, 0, 0, 0, time.UTC), time.Date(2019, time.Month(5), 14, 0, 0, 0, 0, time.UTC)},
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
	for i, tc := range tt {
		def := &root.Definitions.Definitions[i]
		a := def.Advisory
		if got, want := a.Updated.Date, tc.Updated; !cmp.Equal(got, want) {
			t.Errorf("def #%d: %s: Updated\n%v", i, def.Title, cmp.Diff(got, want))
		}
		if got, want := a.Issued.Date, tc.Issued; !cmp.Equal(got, want) {
			t.Errorf("def #%d: %s: Issued\n%v", i, def.Title, cmp.Diff(got, want))
		}
	}
}

func TestDebianDates(t *testing.T) {
	var tt = []time.Time{
		time.Date(2004, time.Month(10), 29, 0, 0, 0, 0, time.UTC),
		time.Date(2003, time.Month(6), 6, 0, 0, 0, 0, time.UTC),
		time.Date(2005, time.Month(2), 2, 0, 0, 0, 0, time.UTC),
		time.Date(2020, time.Month(4), 24, 0, 0, 0, 0, time.UTC),
		time.Date(2020, time.Month(4), 24, 0, 0, 0, 0, time.UTC),
		time.Date(2020, time.Month(4), 24, 0, 0, 0, 0, time.UTC),
		time.Date(2020, time.Month(4), 24, 0, 0, 0, 0, time.UTC),
		time.Date(2020, time.Month(4), 24, 0, 0, 0, 0, time.UTC),
		time.Date(2020, time.Month(4), 24, 0, 0, 0, 0, time.UTC),
		time.Date(2020, time.Month(4), 24, 0, 0, 0, 0, time.UTC),
		time.Date(2020, time.Month(4), 24, 0, 0, 0, 0, time.UTC),
		time.Date(2001, time.Month(7), 11, 0, 0, 0, 0, time.UTC),
		time.Date(2001, time.Month(8), 9, 0, 0, 0, 0, time.UTC),
	}

	f, err := os.Open("../testdata/oval-definitions-buster.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	var root Root
	if err := xml.NewDecoder(f).Decode(&root); err != nil {
		t.Fatal(err)
	}
	for i, tc := range tt {
		def := &root.Definitions.Definitions[i]
		a := def.Debian
		if got, want := a.Date.Date, tc; !cmp.Equal(got, want) {
			t.Errorf("def #%d: %s:\n%v", i, def.Title, cmp.Diff(got, want))
		}
	}
}

func TestUbuntuDates(t *testing.T) {
	var tt = []time.Time{
		time.Time{}, // First entry has no date attached.
		time.Date(2015, time.Month(2), 23, 0, 0, 0, 0, time.UTC),
		time.Date(2007, time.Month(1), 16, 23, 28, 0, 0, time.UTC),
		time.Date(2007, time.Month(9), 26, 23, 17, 0, 0, time.UTC),
		time.Date(2008, time.Month(11), 18, 16, 0, 0, 0, time.UTC),
		time.Date(2008, time.Month(11), 18, 16, 0, 0, 0, time.UTC),
		time.Date(2008, time.Month(11), 18, 16, 0, 0, 0, time.UTC),
		time.Date(2008, time.Month(11), 18, 16, 0, 0, 0, time.UTC),
		time.Date(2018, time.Month(11), 18, 19, 29, 0, 0, time.UTC),
		time.Date(2009, time.Month(4), 23, 19, 30, 0, 0, time.UTC),
	}

	f, err := os.Open("../testdata/com.ubuntu.bionic.cve.oval.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	var root Root
	if err := xml.NewDecoder(f).Decode(&root); err != nil {
		t.Fatal(err)
	}
	for i, tc := range tt {
		def := &root.Definitions.Definitions[i]
		a := def.Advisory
		if got, want := a.PublicDate.Date, tc; !cmp.Equal(got, want) {
			t.Errorf("def #%d: %s:\n%v", i, def.Title, cmp.Diff(got, want))
		}
	}
}

func TestPointlessElement(t *testing.T) {
	const doc = xml.Header + `<div><issued date=""/></div>`
	var got struct {
		Date Date `xml:"issued"`
	}

	rd := strings.NewReader(doc)
	if err := xml.NewDecoder(rd).Decode(&got); err != nil {
		t.Error(err)
	}
}
