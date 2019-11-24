package object

import (
    "testing"
)

func TestStringHashKey(t *testing.T) {
    hello1 := &String{Value: "Hello"}
    hello2 := &String{Value: "Hello"}
    some1 := &String{Value: "Something"}
    some2 := &String{Value: "Something"}

    if hello1.HashKey() != hello2.HashKey() {
        t.Errorf("strings with same content have different hash")
    }

    if some1.HashKey() != some2.HashKey() {
        t.Errorf("strings with same content have different hash")
    }

    if hello1.HashKey() == some1.HashKey() {
        t.Errorf("strings with different content have same hash")
    }
}

func TestIntHashKey(t *testing.T) {
    hello1 := &Integer{Value: 1337}
    hello2 := &Integer{Value: 1337}
    some1 := &Integer{Value: 42}
    some2 := &Integer{Value: 42}

    if hello1.HashKey() != hello2.HashKey() {
        t.Errorf("integers with same content have different hash")
    }

    if some1.HashKey() != some2.HashKey() {
        t.Errorf("integers with same content have different hash")
    }

    if hello1.HashKey() == some1.HashKey() {
        t.Errorf("integers with different content have same hash")
    }
}

func TestBoolHashKey(t *testing.T) {
    hello1 := &Boolean{Value: true}
    hello2 := &Boolean{Value: true}
    some1 := &Boolean{Value: false}
    some2 := &Boolean{Value: false}

    if hello1.HashKey() != hello2.HashKey() {
        t.Errorf("booleans with same content have different hash")
    }

    if some1.HashKey() != some2.HashKey() {
        t.Errorf("booleans with same content have different hash")
    }

    if hello1.HashKey() == some1.HashKey() {
        t.Errorf("booleans with different content have same hash")
    }
}
