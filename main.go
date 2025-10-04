package main

import (
	"fmt"
	"go/token"
	"go/types"
)

func main() {
	intT := types.Typ[types.Int]
	int64T := types.Typ[types.Int64]
	strT := types.Typ[types.String]
	floatT := types.Typ[types.Float64]
	boolT := types.Typ[types.Bool]

	// 合成タイプ
	sliceInt := types.NewSlice(intT)                             // []int
	arrayInt := types.NewArray(intT, 3)                          // [3]int
	arraySliceInt := types.NewArray(sliceInt, 3)                 // [3][]int
	mapStrInt := types.NewMap(strT, intT)                        // map[string]int
	ptrInt := types.NewPointer(intT)                             // *int
	chInt := types.NewChan(types.SendRecv, intT)                 // chan int
	fn := types.NewSignatureType(nil, nil, nil, nil, nil, false) // func()

	// 構造体: フィールドすべて比較可能なら構造体も比較可能
	fieldA := types.NewVar(token.NoPos, nil, "a", intT)
	fieldB := types.NewVar(token.NoPos, nil, "b", strT)
	structOK := types.NewStruct([]*types.Var{fieldA, fieldB}, []string{"", ""}) // 比較可能

	fieldB2 := types.NewVar(token.NoPos, nil, "b", types.NewSlice(types.Typ[types.Byte]))
	structNG := types.NewStruct([]*types.Var{fieldA, fieldB2}, []string{"", ""}) // []byte があるので不可

	// 自己参照構造体（ポインタ経由なので比較可能）
	nodeName := types.NewTypeName(token.NoPos, nil, "Node", nil)
	node := types.NewNamed(nodeName, nil, nil)
	nextField := types.NewVar(token.NoPos, nil, "next", types.NewPointer(node))
	nodeStruct := types.NewStruct([]*types.Var{nextField}, []string{""})
	node.SetUnderlying(nodeStruct)

	// 空インターフェース（= any）
	iface := types.NewInterfaceType(nil, nil)
	iface.Complete() // インターフェースは Complete を呼ぶのが慣例

	cases := []struct {
		label string
		typ   types.Type
	}{
		{"int", intT},
		{"int64", int64T},
		{"string", strT},
		{"float64", floatT},
		{"bool", boolT},
		{"[]int", sliceInt},
		{"[3]int", arrayInt},
		{"[3][]int", arraySliceInt},
		{"map[string]int", mapStrInt},
		{"*int", ptrInt},
		{"chan int", chInt},
		{"func()", fn},
		{"struct{a int; b string}", structOK},
		{"struct{a int; b []byte}", structNG},
		{"type Node struct { next *Node }", node},
		{"interface{}", iface},
	}

	for _, c := range cases {
		fmt.Printf("%-40s -> Comparable=%v\n", c.label, types.Comparable(c.typ))
	}
}
