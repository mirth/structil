package structil

import (
	"fmt"
	"strings"
)

const (
	defaultSep = "."
	rootKey    = "!"
)

type Finder interface {
	Struct(names ...string) Finder
	Find(names ...string) Finder
	ToMap() (map[string]interface{}, error)
	HasError() bool
	Error() string
	GetNameSeparator() string
	Reset() Finder
}

type fImpl struct {
	rootGetter Getter
	gMap       map[string]Getter
	fMap       map[string][]string
	eMap       map[string][]error
	ck         string
	sep        string
}

func NewFinder(i interface{}) (Finder, error) {
	return NewFinderWithSep(i, defaultSep)
}

func NewFinderWithSep(i interface{}, sep string) (Finder, error) {
	g, err := NewGetter(i)
	if err != nil {
		return nil, err
	}

	return NewFinderWithGetterAndSep(g, sep)
}

func NewFinderWithGetter(g Getter) (Finder, error) {
	return NewFinderWithGetterAndSep(g, defaultSep)
}

func NewFinderWithGetterAndSep(g Getter, sep string) (Finder, error) {
	if sep == "" {
		return nil, fmt.Errorf("sep [%s] is invalid", sep)
	}

	f := &fImpl{rootGetter: g, sep: sep}

	return f.Reset(), nil
}

func (f *fImpl) Reset() Finder {
	gMap := map[string]Getter{}
	gMap[rootKey] = f.rootGetter
	f.gMap = gMap

	fMap := map[string][]string{}
	fMap[rootKey] = []string{}
	f.fMap = fMap

	eMap := map[string][]error{}
	eMap[rootKey] = []error{}
	f.eMap = eMap

	f.ck = rootKey

	return f
}

func (f *fImpl) Struct(names ...string) Finder {
	if f.HasError() {
		return f
	}

	f.ck = rootKey

	var nextGetter Getter
	var ok bool
	var err error
	nextKey := ""

	for _, name := range names {
		if f.HasError() {
			break
		}

		if nextKey == "" {
			nextKey = name
		} else {
			nextKey = nextKey + f.sep + name
		}
		err = nil
		f.eMap[nextKey] = []error{}

		nextGetter, ok = f.gMap[nextKey]
		if !ok {
			nextGetter, err = NewGetter(f.gMap[f.ck].Get(name))
		}

		if err != nil {
			f.addError(nextKey, fmt.Errorf("Error in name: %s, key: %s. [%v]", name, nextKey, err))
		}

		f.gMap[nextKey] = nextGetter
		f.ck = nextKey
	}

	if !f.HasError() {
		f.fMap[f.ck] = []string{}
	}

	return f
}

func (f *fImpl) addError(key string, err error) Finder {
	if _, ok := f.eMap[key]; !ok {
		f.eMap[key] = []error{}
	}
	f.eMap[key] = append(f.eMap[key], err)

	return f
}

func (f *fImpl) Find(names ...string) Finder {
	if f.HasError() {
		return f
	}

	f.fMap[f.ck] = append(f.fMap[f.ck], names...)

	return f
}

func (f *fImpl) ToMap() (map[string]interface{}, error) {
	if f.HasError() {
		return nil, f
	}

	res := map[string]interface{}{}
	var key string

	for kg, getter := range f.gMap {
		for _, name := range f.fMap[kg] {
			if kg == rootKey {
				key = name
			} else {
				key = kg + f.sep + name
			}

			if !getter.Has(name) {
				f.addError(key, fmt.Errorf("field name %s does not exist", name))
				break
			}

			res[key] = getter.Get(name)
		}
	}

	if f.HasError() {
		return nil, f
	}

	return res, nil
}

func (f *fImpl) HasError() bool {
	for _, errs := range f.eMap {
		if len(errs) > 0 {
			return true
		}
	}

	return false
}

func (f *fImpl) Error() string {
	tmp := []string{}

	for _, errs := range f.eMap {
		for _, err := range errs {
			tmp = append(tmp, err.Error())
		}
	}

	// TODO: prettize
	return strings.Join(tmp, "\n")
}

func (f *fImpl) GetNameSeparator() string {
	return f.sep
}