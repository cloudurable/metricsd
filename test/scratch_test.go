package test

import (
    "testing"
    "fmt"
    "hash/fnv"
)

func TestScratch(test *testing.T) {
    Hash("foobar", 256)
    Hash("foobar", 256)
    Hash("foobar", 64)
    Hash("foobar", 64)
    Hash("foobarx", 256)
    Hash("foobary", 256)
}

func Hash(s string, partitions int32) int32 {
    hasher := fnv.New32()
    if _, err := hasher.Write([]byte(s)); err != nil {
        panic(err)
    }

    hash := int32(hasher.Sum32())
    if hash < 0 {
        hash = -hash
    }
    mod := hash % partitions

    fmt.Printf("%d %d %s\n", mod, hash, s)

    return mod
}
