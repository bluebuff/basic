package container

import (
	"reflect"
	"sync"
	"bytes"
	"fmt"
)

type ConcurrentMap interface {
	Map
}

type myConcurrentMap struct {
	ConcurrentMap
	m        map[interface{}]interface{}
	keyType  reflect.Type
	elemType reflect.Type
	rwmutex  sync.RWMutex
}

func (cmap *myConcurrentMap) Get(key interface{}) interface{} {
	cmap.rwmutex.RLock()
	defer cmap.rwmutex.Unlock()
	return cmap.m[key]
}

func (cmap *myConcurrentMap) isAcceptablePair(k, e interface{}) bool {
	if k == nil || reflect.TypeOf(k) != cmap.keyType {
		return false
	}
	if e == nil || reflect.TypeOf(e) != cmap.elemType {
		return false
	}
	return true
}

func (cmap *myConcurrentMap) Put(key interface{}, elem interface{}) (interface{}, bool) {
	if !cmap.isAcceptablePair(key, elem) {
		return nil, false
	}
	cmap.rwmutex.Lock()
	defer cmap.rwmutex.Unlock()
	oldEle := cmap.m[key]
	cmap.m[key] = elem
	return oldEle, true
}

func (cmap *myConcurrentMap) Remove(key interface{}) interface{} {
	cmap.rwmutex.Lock()
	defer cmap.rwmutex.Unlock()
	oldElem := cmap.m[key]
	delete(cmap.m, key)
	return oldElem
}

func (cmap *myConcurrentMap) Clear() {
	cmap.rwmutex.Lock()
	defer cmap.rwmutex.Unlock()
	cmap.m = make(map[interface{}]interface{})
}

func (cmap *myConcurrentMap) Len() int {
	cmap.rwmutex.RLock()
	defer cmap.rwmutex.RUnlock()
	return len(cmap.m)
}

func (cmap *myConcurrentMap) Contains(key interface{}) bool {
	cmap.rwmutex.RLock()
	defer cmap.rwmutex.RUnlock()
	_, ok := cmap.m[key]
	return ok
}

func (cmap *myConcurrentMap) Keys() []interface{} {
	cmap.rwmutex.RLock()
	defer cmap.rwmutex.RUnlock()
	keys := make([]interface{}, len(cmap.m))
	index := 0
	for k := range cmap.m {
		keys[index] = k
		index++
	}
	return keys
}

func (cmap *myConcurrentMap) Elems() []interface{} {
	cmap.rwmutex.RLock()
	defer cmap.rwmutex.RUnlock()
	elems := make([]interface{}, len(cmap.m))
	index := 0
	for _, v := range cmap.m {
		elems[index] = v
		index++
	}
	return elems
}

func (cmap *myConcurrentMap) ToMap() map[interface{}] interface{} {
	cmap.rwmutex.RLock()
	defer cmap.rwmutex.RUnlock()
	replica := make(map[interface{}] interface{})
	for k,v := range cmap.m {
		replica[k] = v
	}
	return replica
}


func (cmap *myConcurrentMap) KeyType() reflect.Type {
	return cmap.keyType
}

func (cmap *myConcurrentMap) ElemType() reflect.Type {
	return cmap.elemType
}

func (cmap *myConcurrentMap) String() string {
	var buf bytes.Buffer
	buf.WriteString("ConcurrentMap<")
	buf.WriteString(cmap.keyType.Kind().String())
	buf.WriteString(",")
	buf.WriteString(cmap.elemType.Kind().String())
	buf.WriteString(">{")
	first := true
	for k, v := range cmap.m {
		if first {
			first = false
		} else {
			buf.WriteString(" ")
		}
		buf.WriteString(fmt.Sprintf("%v", k))
		buf.WriteString(":")
		buf.WriteString(fmt.Sprintf("%v", v))
	}
	buf.WriteString("}")
	return buf.String()
}

func NewConcurrentMap(keyType, elemType reflect.Type) ConcurrentMap {
	return &myConcurrentMap{
		keyType:  keyType,
		elemType: elemType,
		m:        make(map[interface{}]interface{})}
}