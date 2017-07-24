package main

import (
    "github.com/hashicorp/hcl"
    "fmt"
    "encoding/json"
)

type Foo struct {
    Region      string
    CamelCase   string `hcl:"camel_case"`
    Bars        []Bar  `hcl:"bar"`
}

type Bar struct {
    Name      string
    CamelBaz  string   `hcl:"camel_baz"`
}

func main3() {
    fmt.Println(data, "\n")
    cfg := &Foo{}
    err := hcl.Decode(cfg, data)
    if err != nil {
        fmt.Println("Err", err)
    } else {
        b, err := json.MarshalIndent(cfg, "", "  ")
        if err == nil {
            fmt.Printf("%s\n", b)
        } else {
            fmt.Println("Err", err)
        }
    }
}

var data = `
region     = "us-west-2"
camel_case = "blah"

bar "" {
    name = "nameAprime"
    camel_baz = "bazA"
}
bar "" {
    name = "nameBprime"
    camel_baz = "bazB"
}
`