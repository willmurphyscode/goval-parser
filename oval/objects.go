package oval

import (
	"fmt"
	"sync"
)

func (o *Objects) init() {
	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		o.lineMemo = make(map[string]int, len(o.LineObjects))
		for i, v := range o.LineObjects {
			o.lineMemo[v.ID] = i
		}
	}()

	go func() {
		defer wg.Done()
		o.version55Memo = make(map[string]int, len(o.Version55Objects))
		for i, v := range o.Version55Objects {
			o.version55Memo[v.ID] = i
		}
	}()

	go func() {
		defer wg.Done()
		o.rpminfoMemo = make(map[string]int, len(o.RPMInfoObjects))
		for i, v := range o.RPMInfoObjects {
			o.rpminfoMemo[v.ID] = i
		}
	}()

	wg.Wait()
}

// Lookup returns the kind of object and index into that kind-specific slice, if
// found.
func (o *Objects) Lookup(ref string) (kind string, index int, err error) {
	o.once.Do(o.init)
	if i, ok := o.lineMemo[ref]; ok {
		return o.LineObjects[i].XMLName.Local, i, nil
	}

	if i, ok := o.version55Memo[ref]; ok {
		return o.Version55Objects[i].XMLName.Local, i, nil
	}
	if i, ok := o.rpminfoMemo[ref]; ok {
		return o.RPMInfoObjects[i].XMLName.Local, i, nil
	}

	// We didn't find it, maybe we can say why.
	id, err := ParseID(ref)
	if err != nil {
		return "", -1, err
	}
	if id.Type != OvalObject {
		return "", -1, fmt.Errorf("oval: wrong identifier type %q", id.Type)
	}
	return "", -1, ErrNotFound
}
