package luar

import (
	"fmt"
	"reflect"
	"unicode"
	"unicode/utf8"

	"github.com/lunixbochs/luaish"
)

func check(L *lua.LState, idx int, kind reflect.Kind) (ref reflect.Value, mt *Metatable, isPtr bool) {
	ud := L.CheckUserData(idx)
	ref = reflect.ValueOf(ud.Value)
	if ref.Kind() != kind {
		if ref.Kind() != reflect.Ptr || ref.Elem().Kind() != kind {
			s := kind.String()
			L.ArgError(idx, "expecting "+s+" or "+s+" pointer")
		}
		isPtr = true
	}
	mt = &Metatable{LTable: ud.Metatable.(*lua.LTable)}
	return
}

func tostring(L *lua.LState) int {
	ud := L.CheckUserData(1)
	value := ud.Value
	if stringer, ok := value.(fmt.Stringer); ok {
		L.Push(lua.LString(stringer.String()))
	} else {
		L.Push(lua.LString(fmt.Sprintf("userdata (luar): %p", ud)))
	}
	return 1
}

func getUnexportedName(name string) string {
	first, n := utf8.DecodeRuneInString(name)
	if n == 0 {
		return name
	}
	return string(unicode.ToLower(first)) + name[n:]
}

func raiseInvalidArg(L *lua.LState, arg int, input lua.LValue, hint reflect.Type) {
	receivedTypeName := input.Type().String()
	if udArg, ok := input.(*lua.LUserData); ok {
		refArg, ok := udArg.Value.(reflect.Value)
		if ok {
			receivedTypeName = refArg.Type().String()
		}
	}
	L.ArgError(arg, hint.String()+" expected, got "+receivedTypeName)
}
