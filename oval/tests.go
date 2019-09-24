package oval

import (
	"fmt"
	"sync"
)

// Lookup returns the kind of test and index into that kind-specific slice, if
// found.
func (t *Tests) Lookup(ref string) (kind string, index int, err error) {
	id, err := ParseID(ref)
	if err != nil {
		return "", -1, err
	}
	if id.Type != OvalTest {
		return "", -1, fmt.Errorf("oval: wrong identifier type %q", id.Type)
	}
	type result struct {
		kind  string
		index int
	}
	ch := make(chan *result)
	done := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		defer wg.Done()
		for i, t := range t.LineTests {
			if t.ID == ref {
				ch <- &result{t.XMLName.Local, i}
				break
			}
		}
	}()
	go func() {
		defer wg.Done()
		for i, t := range t.Version55Tests {
			if t.ID == ref {
				ch <- &result{t.XMLName.Local, i}
				break
			}
		}
	}()
	go func() {
		defer wg.Done()
		for i, t := range t.RPMInfoTests {
			if t.ID == ref {
				ch <- &result{t.XMLName.Local, i}
				break
			}
		}
	}()
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case r := <-ch:
		return r.kind, r.index, nil
	case <-done:
	}
	return "", -1, ErrNotFound
}
