package object

type Environment struct {
    store map[string]Object
    constNames map[string]bool
    outer *Environment
}

func NewEnvironment() *Environment {
    return &Environment{store: make(map[string]Object), constNames: make(map[string]bool), outer: nil}
}

func NewEnclosingEnvironment(outer *Environment) *Environment {
    env := NewEnvironment()
    env.outer = outer

    return env
}

func (e *Environment) Get(name string) (Object, bool) {
    value, ok := e.store[name]
    if !ok && e.outer != nil {
        value, ok = e.outer.Get(name)
    }
    return value, ok
}

func (e *Environment) Set(name string, value Object) bool {
    if !e.hasEntry(name) {
        if e.outer != nil {
            return e.outer.Set(name, value)
        }
        return false
    }
    if e.isConst(name) {
        return false
    }
    e.store[name] = value
    return true
}

func (e *Environment) Add(name string, value Object) bool {
    if e.hasEntry(name) {
        return false
    }
    e.store[name] = value
    e.constNames[name] = false

    return true
}

func (e *Environment) AddConst(name string, value Object) bool {
    if e.hasEntry(name) {
        return false
    }
    e.store[name] = value
    e.constNames[name] = true

    return true
}

func (e *Environment) hasEntry(name string) bool {
    _, has := e.store[name]
    return has
}

func (e *Environment) isConst(name string) bool {
    isConst, ok := e.constNames[name]
    return ok && isConst
}
