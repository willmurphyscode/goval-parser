package oval

import (
	"encoding/xml"
	"github.com/google/go-cmp/cmp"
	"os"
	"testing"
)

func TestAdvisory(t *testing.T) {
	f, err := os.Open("../testdata/RHEL-8-including-unpatched-test.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	var root Root
	if err := xml.NewDecoder(f).Decode(&root); err != nil {
		t.Fatal(err)
	}

	arr := root.Definitions.Definitions
	m := len(arr)
	if !cmp.Equal(m, 4) {
		t.Error("Definition list length is incorrect")
	}
	stateExample := [4]string{"Will not fix", "Will not fix", "Will not fix", "Affected"}
	componentLength := [4]int{1, 9, 1, 3}
	for i := 0; i < m; i++ {
		if !cmp.Equal(len(arr[i].Advisory.Affected.Resolution), 1) {
			t.Fatal("Affected resolution list length is incorrect")
		}
		resolution := arr[i].Advisory.Affected.Resolution[0]
		name := resolution.State
		if !cmp.Equal(name, stateExample[i]) {
			t.Fatal("Resolution state is incorrect")
		}
		if !cmp.Equal(len(resolution.Components), componentLength[i]) {
			t.Fatal("Component list length is incorrect")
		}
	}

}
