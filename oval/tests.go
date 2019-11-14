package oval

import (
	"fmt"
	"sync"
)

// Init sets up the memoization maps.
func (t *Tests) init() {
	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		t.lineMemo = make(map[string]int, len(t.LineTests))
		for i, v := range t.LineTests {
			t.lineMemo[v.ID] = i
		}
	}()

	go func() {
		defer wg.Done()
		t.version55Memo = make(map[string]int, len(t.Version55Tests))
		for i, v := range t.Version55Tests {
			t.version55Memo[v.ID] = i
		}
	}()

	go func() {
		defer wg.Done()
		t.rpminfoMemo = make(map[string]int, len(t.RPMInfoTests))
		for i, v := range t.RPMInfoTests {
			t.rpminfoMemo[v.ID] = i
		}
	}()

	wg.Wait()
}

// Lookup returns the kind of test and index into that kind-specific slice, if
// found.
func (t *Tests) Lookup(ref string) (kind string, index int, err error) {
	// Pay to construct an index, once.
	t.once.Do(t.init)

	if i, ok := t.lineMemo[ref]; ok {
		return t.LineTests[i].XMLName.Local, i, nil
	}
	if i, ok := t.version55Memo[ref]; ok {
		return t.Version55Tests[i].XMLName.Local, i, nil
	}
	if i, ok := t.rpminfoMemo[ref]; ok {
		return t.RPMInfoTests[i].XMLName.Local, i, nil
	}

	// We didn't find it, maybe we can say why.
	id, err := ParseID(ref)
	if err != nil {
		return "", -1, err
	}
	if id.Type != OvalTest {
		return "", -1, fmt.Errorf("oval: wrong identifier type %q", id.Type)
	}
	return "", -1, ErrNotFound(ref)
}
