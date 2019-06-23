structil  [![GoDoc](https://godoc.org/github.com/goldeneggg/structil?status.png)](https://godoc.org/github.com/goldeneggg/structil)
==========

[![Build Status](https://travis-ci.org/goldeneggg/structil.svg?branch=master)](https://travis-ci.org/goldeneggg/structil)
[![Go Report Card](https://goreportcard.com/badge/github.com/goldeneggg/structil)](https://goreportcard.com/report/github.com/goldeneggg/structil)
[![GolangCI](https://golangci.com/badges/github.com/goldeneggg/gat.svg)](https://golangci.com/r/github.com/goldeneggg/structil)
[![Codecov](https://codecov.io/github/goldeneggg/structil/coverage.svg?branch=master)](https://codecov.io/github/goldeneggg/structil?branch=master)
[![MIT License](http://img.shields.io/badge/license-MIT-lightgrey.svg)](https://github.com/goldeneggg/structil/blob/master/LICENSE)

Struct Utilities for runtime and dynamic environment in Go.

__Table of Contents__

<!-- TOC depthFrom:1 -->

- [`Finder`](#finder)
  - [With config? use `FinderKeys`](#with-config-use-finderkeys)
- [`Getter`](#getter)
  - [`MapGet` method](#mapget-method)
- [`DynamicStruct`](#dynamicstruct)
  - [JSON unmershal with `DynamicStructMap`](#json-unmershal-with-dynamicstruct)

<!-- /TOC -->

## `Finder`
We can access usefully nested struct fields using field name string.

[Sample script on playground](https://play.golang.org/p/AcF5c7Prf3z).

```go
package main

import (
	"fmt"

	"github.com/goldeneggg/structil"
)

type group struct {
	Name string
	Boss string
}

type company struct {
	Name    string
	Address string
	Period  int
	Group   *group
}

type school struct {
	Name          string
	GraduatedYear int
}

type person struct {
	Name    string
	Age     int
	Company *company
	School  *school
}

func main() {
	i := &person{
		Name: "Lisa Mary",
		Age:  34,
		Company: &company{
			Name:    "ZZZ Air inc.",
			Address: "Boston",
			Period:  11,
			Group: &group{
				Name: "ZZZZZZ Holdings",
				Boss: "Donald Mac",
			},
		},
		School: &school{
			Name:          "XYZ College",
			GraduatedYear: 2008,
		},
	}

	finder, err := structil.NewFinder(i)
	if err != nil {
		panic(err)
	}

	// We can use method chain for Find and Into methods.
	//  - FindTop returns a Finder that top level fields in struct are looked up and held named "names".
	//  - Into returns a Finder that nested struct fields are looked up and held named "names".
	//  - Find returns a Finder that fields in struct are looked up and held named "names".
	// And finally, we can call ToMap method for converting from struct to map.
	m, err := finder.
		FindTop("Name", "School").
		Into("Company").Find("Address").
		Into("Company", "Group").Find("Name", "Boss").
		ToMap()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v", m)
}
```

Result as follows.

```
map[string]interface {}{"Company.Address":"Boston", "Company.Group.Boss":"Donald Mac", "Company.Group.Name":"ZZZZZZ Holdings", "Name":"Lisa Mary", "School":main.school{Name:"XYZ College", GraduatedYear:2008}}
```

### With config? use `FinderKeys`
We can create a Finder from the configuration file that have some finding target keys.

We support some file format of configuration file such as `yaml`, `json`, `toml` and more.

Thanks for the awesome configuration management library [spf13/viper](https://github.com/spf13/viper).

```go
package main

import (
  "fmt"

  "github.com/goldeneggg/structil"
)

type person struct {
  Name    string
  Age     int
  Company *company
  Schools []*school
}

type school struct {
  Name          string
  GraduatedYear int
}

type company struct {
  Name    string
  Address string
  Period  int
  Group   *group
}

type group struct {
  Name string
  Boss string
}

func main() {
  i := &person{
    Name: "Lisa Mary",
    Age:  34,
    Company: &company{
      Name:    "ZZZ Air inc.",
      Address: "Boston",
      Period:  11,
      Group: &group{
        Name: "ZZZZZZ Holdings",
        Boss: "Donald Mac",
      },
    },
    Schools: []*school{
      {
        Name:          "STU High School",
        GraduatedYear: 2005,
      },
      {
        Name:          "XYZ College",
        GraduatedYear: 2008,
      },
    },
  }

  json(i)
  yml(i)
}

func json(i *person) {
  fks, err := structil.NewFinderKeysFromConf("examples/finder_from_conf", "ex_json")
  if err != nil {
    fmt.Printf("error: %v\n", err)
    return
  }
  fmt.Printf("fks.Keys(json): %#v\n", fks.Keys())

  finder, err := structil.NewFinder(i)
  if err != nil {
    fmt.Printf("error: %v\n", err)
    return
  }

  m, err := finder.FromKeys(fks).ToMap()
  fmt.Printf("Found Map(json): %#v, err: %v\n", m, err)
}

func yml(i *person) {
  fks, err := structil.NewFinderKeysFromConf("examples/finder_from_conf", "ex_yml")
  if err != nil {
    fmt.Printf("error: %v\n", err)
    return
  }
  fmt.Printf("fks.Keys(yml): %#v\n", fks.Keys())

  finder, err := structil.NewFinder(i)
  if err != nil {
    fmt.Printf("error: %v\n", err)
    return
  }

  m, err := finder.FromKeys(fks).ToMap()
  fmt.Printf("Found Map(yml): %#v, err: %v\n", m, err)
}
```

File `examples/finder_from_conf/ex_json.json` as follows:

```json
{
  "Keys":[
    {
      "Company":[
        {
          "Group":[
            "Name",
            "Boss"
          ]
        },
        "Address",
        "Period"
      ]
    },
    "Name",
    "Age"
  ]
}
```

File `examples/finder_from_conf/ex_yml.yml` as follows:

```yml
Keys:
  - Company:
    - Group:
      - Name
      - Boss
    - Address
    - Period
  - Name
  - Age
```

Result as follows.

```
fks.Keys(json): []string{"Company.Group.Name", "Company.Group.Boss", "Company.Address", "Company.Period", "Name", "Age"}
Found Map(json): map[string]interface {}{"Age":34, "Company.Address":"Boston", "Company.Group.Boss":"Donald Mac", "Company.Group.Name":"ZZZZZZ Holdings", "Company.Period":11, "Name":"Lisa Mary"}, err: <nil>
fks.Keys(yml): []string{"Company.Group.Name", "Company.Group.Boss", "Company.Address", "Company.Period", "Name", "Age"}
Found Map(yml): map[string]interface {}{"Age":34, "Company.Address":"Boston", "Company.Group.Boss":"Donald Mac", "Company.Group.Name":"ZZZZZZ Holdings", "Company.Period":11, "Name":"Lisa Mary"}, err: <nil>
```

## `Getter`
We can access a struct using field name string, like map.

[Sample script on playground](https://play.golang.org/p/3CNDJpW3UmN).

```go
package main

import (
	"fmt"

	"github.com/goldeneggg/structil"
)

type company struct {
	Name    string
	Address string
	Period  int
}

type person struct {
	Name    string
	Age     int
	Company *company
}

func main() {
	i := &person{
		Name: "Mike Davis",
		Age:  27,
		Company: &company{
			Name:    "Scott inc.",
			Address: "Osaka",
			Period:  2,
		},
	}

	getter, err := structil.NewGetter(i)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Name: %+v, Age: %+v, Company: %+v\n", getter.String("Name"), getter.Int("Age"), getter.Get("Company"))
}
```

Result as follows.

```
Name: "Mike Davis", Age: 27, Company: main.company{Name:"Scott inc.", Address:"Osaka", Period:2}
```

### `MapGet` method
`MapGet` method provides the __Map__ collection function for slice of struct

[Sample script on playground](https://play.golang.org/p/98wCWCrs0vf).

```go
package main

import (
	"fmt"

	"github.com/goldeneggg/structil"
)

type company struct {
	Name    string
	Address string
	Period  int
}

type person struct {
	Name      string
	Age       int
	Companies []*company
}

func main() {
	i := &person{
		Name: "John",
		Age:  28,
		Companies: []*company{
			{
				Name:    "Dragons inc.",
				Address: "Nagoya",
				Period:  3,
			},
			{
				Name:    "Swallows inc.",
				Address: "Tokyo",
				Period:  2,
			},
		},
	}

	getter, err := structil.NewGetter(i)
	if err != nil {
		panic(err)
	}

	fn := func(i int, g structil.Getter) (interface{}, error) {
		return fmt.Sprintf(
			"You worked for %d years since you joined the company %s",
			g.Int("Period"),
			g.String("Name"),
		), nil
	}

	intfs, err := getter.MapGet("Companies", fn)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v\n", intfs)
}
```

Result as follows.

```
[]interface {}{"You worked for 3 years since you joined the company Dragons inc.", "You worked for 2 years since you joined the company Swallows inc."}
```

## `DynamicStruct`
We can create dynamic and runtime struct.

```go
package main

import (
	"fmt"

	"github.com/goldeneggg/structil"
	"github.com/goldeneggg/structil/dynamicstruct"
)

// Hoge is test struct
type Hoge struct {
	Key   string
	Value interface{}
}

var (
	hoge    Hoge
	hogePtr *Hoge
)

func main() {
  // Add fields using Builder. We can use AddXXX method chain.
	b := dynamicstruct.NewBuilder().
		AddString("StringField").
		AddInt("IntField").
		AddFloat("FloatField").
		AddBool("BoolField").
		AddMap("MapField", dynamicstruct.SampleString, dynamicstruct.SampleFloat).
		AddChanBoth("ChanBothField", dynamicstruct.SampleInt).
		AddStructPtr("StructPtrField", hogePtr).
		AddSlice("SliceField", hogePtr)

  // Available for remove field by Remove method
	b = b.Remove("FloatField")

  // Build method generates a DynamicStruct
	ds := b.Build()

	// Decode from map to DynamicStruct
	input := map[string]interface{}{
		"StringField": "Test String Field",
		"IntField":    12345,
		"BoolField":   true,
	}
	dec, err := ds.DecodeMap(input)
	if err != nil {
		panic(err)
	}

  // confirm decoded result using Getter
	g, err := structil.NewGetter(dec)
	if err != nil {
		panic(err)
	}
	fmt.Printf("String: %v, Int: %v, Bool: %v\n", g.String("StringField"), g.Int("IntField"), g.Get("BoolField"))
}
```

Result as follows.

```
String: Test String Field, Int: 12345, Bool: true
```

### JSON unmershal with `DynamicStruct`

```go
package main

import (
	"encoding/json"
	"fmt"

	"github.com/goldeneggg/structil"
	"github.com/goldeneggg/structil/dynamicstruct"
)

// Hoge is test struct
type Hoge struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

var (
	hoge    Hoge
	hogePtr *Hoge
)

func main() {
	b := dynamicstruct.NewBuilder().
		AddStringWithTag("StringField", `json:"string_field"`).
		AddIntWithTag("IntField", `json:"int_field"`).
		AddFloatWithTag("FloatField", `json:"float_field"`).
		AddBoolWithTag("BoolField", `json:"bool_field"`).
		AddStructPtrWithTag("StructPtrField", hogePtr, `json:"struct_ptr_field"`)

	// Get interface of DynamicStruct using Interface() method
	ds := b.Build()
	intf := ds.Interface()

	// try json unmarshal
	input := []byte(`
{
	"string_field":"あいうえお",
	"int_field":9876,
	"float_field":5.67,
	"bool_field":true,
	"struct_ptr_field":{
		"key":"hogekey",
		"value":"hogevalue"
	}
}
`)

	err := json.Unmarshal(input, &intf)
	if err != nil {
		panic(err)
	}

	g, err := structil.NewGetter(intf)
	if err != nil {
		panic(err)
	}
	fmt.Printf("String: %v, Float: %v, StructPtr: %+v\n", g.String("StringField"), g.Float64("FloatField"), g.Get("StructPtrField"))
}
```

Result as follows.

```
String: あいうえお, Float: 5.67, StructPtr: {Key:hogekey Value:hogevalue}
```
