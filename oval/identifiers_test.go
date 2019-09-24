package oval

import "testing"

func TestIDParse(t *testing.T) {
	ts := []string{
		"oval:com.redhat.rhsa:def:20100720002",
	}
	for _, tc := range ts {
		id, err := ParseID(tc)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("%q → %+#v", id, id)
	}

	tt := [][]byte{
		[]byte("oval:com.redhat.rhsa:def:20100720002"),
		[]byte("oval:com.redhat.rhsa:tst:20100720002"),
		[]byte("oval:com.redhat.rhsa:obj:20100720002"),
		[]byte("oval:com.redhat.rhsa:ste:20100720002"),
		[]byte("oval:com.redhat.rhsa:var:20100720002"),
	}
	for _, tc := range tt {
		var id ID
		if err := id.UnmarshalText(tc); err != nil {
			t.Fatal(err)
		}
		t.Logf("%q → %+#v", id, id)
	}

	tt = [][]byte{
		[]byte("oval:com.redhat.rhsa:obj:20100720002:extrasegment"),
		[]byte("rhombus:com.redhat.rhsa:obj:20100720002"),
		[]byte("oval:com.redhat.rhsa:nsa:20100720002"),
		[]byte("oval:com.redhat.rhsa:obj:201007200fa"),
	}
	for _, tc := range tt {
		err := (&ID{}).UnmarshalText(tc)
		if err == nil {
			t.Fatal("wanted error, got nil")
		}
		t.Logf("%q → %v", string(tc), err)
	}
}
