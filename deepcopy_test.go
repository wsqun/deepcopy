package deepcopy

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

// just basic is this working stuff
func TestSimple(t *testing.T) {
	Strings := []string{"a", "b", "c"}
	cpyS := Copy[[]string](Strings)
	if len(cpyS) != len(Strings) {
		t.Errorf("[]string: len was %d; want %d", len(cpyS), len(Strings))
		goto CopyBools
	}
	for i, v := range Strings {
		if v != cpyS[i] {
			t.Errorf("[]string: got %v at index %d of the copy; want %v", cpyS[i], i, v)
		}
	}

CopyBools:
	Bools := []bool{true, true, false, false}
	cpyB := Copy[[]bool](Bools)
	if len(cpyB) != len(Bools) {
		t.Errorf("[]bool: len was %d; want %d", len(cpyB), len(Bools))
		goto CopyBytes
	}
	for i, v := range Bools {
		if v != cpyB[i] {
			t.Errorf("[]bool: got %v at index %d of the copy; want %v", cpyB[i], i, v)
		}
	}

CopyBytes:
	Bytes := []byte("hello")
	cpyBt := Copy[[]byte](Bytes)
	if len(cpyBt) != len(Bytes) {
		t.Errorf("[]byte: len was %d; want %d", len(cpyBt), len(Bytes))
		goto CopyInts
	}
	for i, v := range Bytes {
		if v != cpyBt[i] {
			t.Errorf("[]byte: got %v at index %d of the copy; want %v", cpyBt[i], i, v)
		}
	}

CopyInts:
	Ints := []int{42}
	cpyI := Copy[[]int](Ints)
	if len(cpyI) != len(Ints) {
		t.Errorf("[]int: len was %d; want %d", len(cpyI), len(Ints))
		goto CopyUints
	}
	for i, v := range Ints {
		if v != cpyI[i] {
			t.Errorf("[]int: got %v at index %d of the copy; want %v", cpyI[i], i, v)
		}
	}

CopyUints:
	Uints := []uint{1, 2, 3, 4, 5}
	cpyU := Copy[[]uint](Uints)
	if len(cpyU) != len(Uints) {
		t.Errorf("[]uint: len was %d; want %d", len(cpyU), len(Uints))
		goto CopyFloat32s
	}
	for i, v := range Uints {
		if v != cpyU[i] {
			t.Errorf("[]uint: got %v at index %d of the copy; want %v", cpyU[i], i, v)
		}
	}

CopyFloat32s:
	Float32s := []float32{3.14}
	cpyF := Copy[[]float32](Float32s)
	if len(cpyF) != len(Float32s) {
		t.Errorf("[]float32: len was %d; want %d", len(cpyF), len(Float32s))
		goto CopyInterfaces
	}
	for i, v := range Float32s {
		if v != cpyF[i] {
			t.Errorf("[]float32: got %v at index %d of the copy; want %v", cpyF[i], i, v)
		}
	}

CopyInterfaces:
	Interfaces := []interface{}{"a", 42, true, 4.32}
	cpyIf := Copy[[]interface{}](Interfaces)
	if len(cpyIf) != len(Interfaces) {
		t.Errorf("[]interface{}: len was %d; want %d", len(cpyIf), len(Interfaces))
		return
	}
	for i, v := range Interfaces {
		if v != cpyIf[i] {
			t.Errorf("[]interface{}: got %v at index %d of the copy; want %v", cpyIf[i], i, v)
		}
	}
}

type Basics struct {
	String      string
	Strings     []string
	StringArr   [4]string
	Bool        bool
	Bools       []bool
	Byte        byte
	Bytes       []byte
	Int         int
	Ints        []int
	Int8        int8
	Int8s       []int8
	Int16       int16
	Int16s      []int16
	Int32       int32
	Int32s      []int32
	Int64       int64
	Int64s      []int64
	Uint        uint
	Uints       []uint
	Uint8       uint8
	Uint8s      []uint8
	Uint16      uint16
	Uint16s     []uint16
	Uint32      uint32
	Uint32s     []uint32
	Uint64      uint64
	Uint64s     []uint64
	Float32     float32
	Float32s    []float32
	Float64     float64
	Float64s    []float64
	Complex64   complex64
	Complex64s  []complex64
	Complex128  complex128
	Complex128s []complex128
	Interface   interface{}
	Interfaces  []interface{}
}

// These tests test that all supported basic types are copied correctly.  This
// is done by copying a struct with fields of most of the basic types as []T.
func TestMostTypes(t *testing.T) {
	test := Basics{
		String:      "kimchi",
		Strings:     []string{"uni", "ika"},
		StringArr:   [4]string{"malort", "barenjager", "fernet", "salmiakki"},
		Bool:        true,
		Bools:       []bool{true, false, true},
		Byte:        'z',
		Bytes:       []byte("abc"),
		Int:         42,
		Ints:        []int{0, 1, 3, 4},
		Int8:        8,
		Int8s:       []int8{8, 9, 10},
		Int16:       16,
		Int16s:      []int16{16, 17, 18, 19},
		Int32:       32,
		Int32s:      []int32{32, 33},
		Int64:       64,
		Int64s:      []int64{64},
		Uint:        420,
		Uints:       []uint{11, 12, 13},
		Uint8:       81,
		Uint8s:      []uint8{81, 82},
		Uint16:      160,
		Uint16s:     []uint16{160, 161, 162, 163, 164},
		Uint32:      320,
		Uint32s:     []uint32{320, 321},
		Uint64:      640,
		Uint64s:     []uint64{6400, 6401, 6402, 6403},
		Float32:     32.32,
		Float32s:    []float32{32.32, 33},
		Float64:     64.1,
		Float64s:    []float64{64, 65, 66},
		Complex64:   complex64(-64 + 12i),
		Complex64s:  []complex64{complex64(-65 + 11i), complex64(66 + 10i)},
		Complex128:  complex128(-128 + 12i),
		Complex128s: []complex128{complex128(-128 + 11i), complex128(129 + 10i)},
		Interfaces:  []interface{}{42, true, "pan-galactic"},
	}

	cpy := Copy[Basics](test)

	// see if they point to the same location
	if fmt.Sprintf("%p", &cpy) == fmt.Sprintf("%p", &test) {
		t.Error("address of copy was the same as original; they should be different")
		return
	}

	// Go through each field and check to see it got copied properly
	if cpy.String != test.String {
		t.Errorf("String: got %v; want %v", cpy.String, test.String)
	}

	if len(cpy.Strings) != len(test.Strings) {
		t.Errorf("Strings: len was %d; want %d", len(cpy.Strings), len(test.Strings))
		goto StringArr
	}
	for i, v := range test.Strings {
		if v != cpy.Strings[i] {
			t.Errorf("Strings: got %v at index %d of the copy; want %v", cpy.Strings[i], i, v)
		}
	}

StringArr:
	for i, v := range test.StringArr {
		if v != cpy.StringArr[i] {
			t.Errorf("StringArr: got %v at index %d of the copy; want %v", cpy.StringArr[i], i, v)
		}
	}

	if cpy.Bool != test.Bool {
		t.Errorf("Bool: got %v; want %v", cpy.Bool, test.Bool)
	}

	if len(cpy.Bools) != len(test.Bools) {
		t.Errorf("Bools: len was %d; want %d", len(cpy.Bools), len(test.Bools))
		goto Bytes
	}
	for i, v := range test.Bools {
		if v != cpy.Bools[i] {
			t.Errorf("Bools: got %v at index %d of the copy; want %v", cpy.Bools[i], i, v)
		}
	}

Bytes:
	if cpy.Byte != test.Byte {
		t.Errorf("Byte: got %v; want %v", cpy.Byte, test.Byte)
	}

	if len(cpy.Bytes) != len(test.Bytes) {
		t.Errorf("Bytes: len was %d; want %d", len(cpy.Bytes), len(test.Bytes))
		goto Ints
	}
	for i, v := range test.Bytes {
		if v != cpy.Bytes[i] {
			t.Errorf("Bytes: got %v at index %d of the copy; want %v", cpy.Bytes[i], i, v)
		}
	}

Ints:
	if cpy.Int != test.Int {
		t.Errorf("Int: got %v; want %v", cpy.Int, test.Int)
	}

	if len(cpy.Ints) != len(test.Ints) {
		t.Errorf("Ints: len was %d; want %d", len(cpy.Ints), len(test.Ints))
		return
	}
	for i, v := range test.Ints {
		if v != cpy.Ints[i] {
			t.Errorf("Ints: got %v at index %d of the copy; want %v", cpy.Ints[i], i, v)
		}
	}
}

type I struct {
	A string
}

func (i *I) DeepCopy() I {
	return I{A: i.A + "_copy"}
}

type NestI struct {
	I *I
}

func TestInterface(t *testing.T) {
	i := &I{A: "test"}
	ni := &NestI{I: i}

	copy := Copy[*NestI](ni)
	if copy.I.A != "test_copy" {
		t.Errorf("Custom copy failed, got %v", copy.I.A)
	}
}

// not meant to be exhaustive
func TestComplexSlices(t *testing.T) {
	orig3Int := [][][]int{{{1, 2, 3}, {11, 22, 33}}, {{7, 8, 9}, {66, 77, 88, 99}}}
	cp3Int := Copy[[][][]int](orig3Int)
	if &orig3Int[0] == &cp3Int[0] {
		t.Error("address of copy was the same as original; they should be different")
		return
	}

	if len(orig3Int) != len(cp3Int) {
		t.Errorf("[][][]int: len was %d; want %d", len(cp3Int), len(orig3Int))
		return
	}

	for i, v := range orig3Int {
		if len(v) != len(cp3Int[i]) {
			t.Errorf("[][][]int: len of first slice was %d; want %d", len(cp3Int[i]), len(v))
			return
		}
		for j, vv := range v {
			if len(vv) != len(cp3Int[i][j]) {
				t.Errorf("[][][]int: len of second slice was %d; want %d", len(cp3Int[i][j]), len(vv))
				return
			}
			for k, vvv := range vv {
				if vvv != cp3Int[i][j][k] {
					t.Errorf("[][][]int: got %v; want %v", cp3Int[i][j][k], vvv)
				}
			}
		}
	}

	slMap := []map[int]string{{1: "one", 2: "two"}, {3: "three", 4: "four"}}
	cpSlMap := Copy[[]map[int]string](slMap)
	if &slMap[0] == &cpSlMap[0] {
		t.Error("address of copy was the same as original; they should be different")
		return
	}

	if len(slMap) != len(cpSlMap) {
		t.Errorf("[]map[int]string: len was %d; want %d", len(cpSlMap), len(slMap))
		return
	}
	for i, v := range slMap {
		if len(v) != len(cpSlMap[i]) {
			t.Errorf("[]map[int]string: len of map was %d; want %d", len(cpSlMap[i]), len(v))
			return
		}
		for k, vv := range v {
			val, ok := cpSlMap[i][k]
			if !ok {
				t.Errorf("[]map[int]string: key %d not found", k)
				return
			}
			if val != vv {
				t.Errorf("[]map[int]string: got %v; want %v", val, vv)
			}
		}
	}
}

type A struct {
	Int    int
	String string
	UintSl []uint
	NilSl  []string
	Map    map[string]int
	MapB   map[string]*B
	SliceB []B
	B      B
	T      time.Time
}

type B struct {
	Vals []string
}

var AStruct = A{
	Int:    42,
	String: "Konichiwa",
	UintSl: []uint{0, 1, 2, 3},
	Map:    map[string]int{"a": 1, "b": 2},
	MapB: map[string]*B{
		"hi":  &B{Vals: []string{"hello", "bonjour"}},
		"bye": &B{Vals: []string{"good-bye", "au revoir"}},
	},
	SliceB: []B{
		B{Vals: []string{"Ciao", "Aloha"}},
	},
	B: B{Vals: []string{"42"}},
	T: time.Now(),
}

func TestStructA(t *testing.T) {
	now := time.Now()
	original := A{
		Int:    42,
		String: "Konichiwa",
		UintSl: []uint{0, 1, 2, 3},
		Map:    map[string]int{"a": 1, "b": 2},
		MapB: map[string]*B{
			"hi":  {Vals: []string{"hello", "bonjour"}},
			"bye": {Vals: []string{"good bye", "au revoir"}},
		},
		SliceB: []B{
			{Vals: []string{"Ciao", "Aloha"}},
		},
		B: B{
			Vals: []string{"42"},
		},
		T: now,
	}

	copiedA := Copy[A](original)

	// 检查基本类型字段
	if copiedA.Int != original.Int {
		t.Errorf("Int: got %v; want %v", copiedA.Int, original.Int)
	}
	if copiedA.String != original.String {
		t.Errorf("String: got %v; want %v", copiedA.String, original.String)
	}

	// 检查切片
	if !reflect.DeepEqual(copiedA.UintSl, original.UintSl) {
		t.Errorf("UintSl: got %v; want %v", copiedA.UintSl, original.UintSl)
	}
	if copiedA.NilSl != nil {
		t.Error("NilSl: expected nil slice")
	}

	// 检查 map
	if !reflect.DeepEqual(copiedA.Map, original.Map) {
		t.Errorf("Map: got %v; want %v", copiedA.Map, original.Map)
	}

	// 检查嵌套结构体
	if !reflect.DeepEqual(copiedA.B, original.B) {
		t.Errorf("B: got %v; want %v", copiedA.B, original.B)
	}

	// 检查时间
	if !copiedA.T.Equal(original.T) {
		t.Errorf("Time: got %v; want %v", copiedA.T, original.T)
	}

	// 检查指针 map
	for k, v := range original.MapB {
		cv := copiedA.MapB[k]
		if cv == v {
			t.Errorf("MapB: expected different pointers for key %s", k)
		}
		if !reflect.DeepEqual(cv.Vals, v.Vals) {
			t.Errorf("MapB: got %v; want %v for key %s", cv.Vals, v.Vals, k)
		}
	}

	// 检查结构体切片
	if !reflect.DeepEqual(copiedA.SliceB, original.SliceB) {
		t.Errorf("SliceB: got %v; want %v", copiedA.SliceB, original.SliceB)
	}
}

type Unexported struct {
	A  string
	B  int
	aa string
	bb int
	cc []int
	dd map[string]string
}

func TestUnexportedFields(t *testing.T) {
	u := &Unexported{
		A:  "A",
		B:  42,
		aa: "aa",
		bb: 11,
		cc: []int{1, 2, 3},
		dd: map[string]string{"a": "b"},
	}

	copiedU := Copy[*Unexported](u)
	if copiedU == u {
		t.Error("expected different pointers")
		return
	}

	// 公开字段应该被复制
	if copiedU.A != u.A {
		t.Errorf("A: got %v; want %v", copiedU.A, u.A)
	}
	if copiedU.B != u.B {
		t.Errorf("B: got %v; want %v", copiedU.B, u.B)
	}

	// 未导出字段不应该被复制
	if copiedU.aa != "" {
		t.Errorf("aa: got %v; want %v", copiedU.aa, "")
	}
	if copiedU.bb != 0 {
		t.Errorf("bb: got %v; want %v", copiedU.bb, 0)
	}
	if copiedU.cc != nil {
		t.Errorf("cc: got %v; want %v", copiedU.cc, nil)
	}
	if copiedU.dd != nil {
		t.Errorf("dd: got %v; want %v", copiedU.dd, nil)
	}
}

// Note: this test will fail until https://github.com/golang/go/issues/15716 is
// fixed and the version it is part of gets released.
type T struct {
	time.Time
}

func TestTimeCopy(t *testing.T) {
	tests := []struct {
		name string
		t    T
	}{
		{
			name: "Test time copy",
			t:    T{Time: time.Now()},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			copied := Copy[T](tt.t)
			if !copied.Equal(tt.t.Time) {
				t.Errorf("Copy() = %v, want %v", copied, tt.t)
			}
		})
	}
}

func TestPointerToStruct(t *testing.T) {
	type Foo struct {
		Bar int
	}

	original := &Foo{Bar: 42}
	copied := Copy[*Foo](original)

	if copied == original {
		t.Error("expected different pointers")
	}
	if copied.Bar != original.Bar {
		t.Errorf("got %v, want %v", copied.Bar, original.Bar)
	}
}

func TestIssue9(t *testing.T) {
	type Foo struct {
		Alpha string
	}

	type Bar struct {
		Beta  string
		Gamma int
		Delta *Foo
	}

	type Biz struct {
		Epsilon map[int]*Bar
	}

	original := Biz{
		Epsilon: map[int]*Bar{
			0: {
				Beta:  "hello",
				Gamma: 1,
				Delta: &Foo{Alpha: "alpha"},
			},
			1: {
				Beta:  "world",
				Gamma: 2,
				Delta: &Foo{Alpha: "beta"},
			},
		},
	}

	copied := Copy[Biz](original)

	// 检查是否真的是深拷贝
	for k, v := range original.Epsilon {
		cv := copied.Epsilon[k]
		if cv == v {
			t.Errorf("expected different pointers for key %d", k)
		}
		if cv.Beta != v.Beta {
			t.Errorf("Beta: got %v; want %v for key %d", cv.Beta, v.Beta, k)
		}
		if cv.Gamma != v.Gamma {
			t.Errorf("Gamma: got %v; want %v for key %d", cv.Gamma, v.Gamma, k)
		}
		if cv.Delta == v.Delta {
			t.Errorf("expected different pointers for Delta at key %d", k)
		}
		if cv.Delta.Alpha != v.Delta.Alpha {
			t.Errorf("Delta.Alpha: got %v; want %v for key %d", cv.Delta.Alpha, v.Delta.Alpha, k)
		}
	}
}

// 基础测试结构体
type TestStruct struct {
	Int        int
	String     string
	Float      float64
	unexported string
}

// 实现了 Copier 接口的结构体
type CustomCopier struct {
	Value int
}

func (c CustomCopier) DeepCopy() CustomCopier {
	return CustomCopier{Value: c.Value * 2} // 自定义复制行为
}

// 嵌套结构体
type NestedStruct struct {
	Basic     TestStruct
	Pointer   *TestStruct
	Slice     []TestStruct
	Map       map[string]TestStruct
	Time      time.Time
	Interface interface{}
}

func TestCopyBasicTypes(t *testing.T) {
	tests := []struct {
		name string
		src  any
	}{
		{"int", 42},
		{"string", "hello"},
		{"float64", 3.14},
		{"bool", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch v := tt.src.(type) {
			case int:
				copied := Copy[int](v)
				if copied != v {
					t.Errorf("Copy() = %v, want %v", copied, v)
				}
			case string:
				copied := Copy[string](v)
				if copied != v {
					t.Errorf("Copy() = %v, want %v", copied, v)
				}
			case float64:
				copied := Copy[float64](v)
				if copied != v {
					t.Errorf("Copy() = %v, want %v", copied, v)
				}
			case bool:
				copied := Copy[bool](v)
				if copied != v {
					t.Errorf("Copy() = %v, want %v", copied, v)
				}
			}
		})
	}
}

func TestCopyStruct(t *testing.T) {
	original := TestStruct{
		Int:        42,
		String:     "test",
		Float:      3.14,
		unexported: "should not copy",
	}

	copied := Copy[TestStruct](original)

	// 检查导出字段是否正确拷贝
	if copied.Int != original.Int {
		t.Errorf("Int: got %v, want %v", copied.Int, original.Int)
	}
	if copied.String != original.String {
		t.Errorf("String: got %v, want %v", copied.String, original.String)
	}
	if copied.Float != original.Float {
		t.Errorf("Float: got %v, want %v", copied.Float, original.Float)
	}

	// 检查未导出字段是否为零值 (没有被拷贝)
	if copied.unexported != "" {
		t.Errorf("unexported field was copied: got %v, want empty string", copied.unexported)
	}

	// 确保是深拷贝
	if &copied == &original {
		t.Error("Copy() returned same address, want different")
	}
}

func TestCopyPointer(t *testing.T) {
	original := &TestStruct{
		Int:    42,
		String: "test",
	}

	copied := Copy[*TestStruct](original)

	if !reflect.DeepEqual(copied, original) {
		t.Errorf("Copy() = %+v, want %+v", copied, original)
	}

	// 确保是深拷贝
	if copied == original {
		t.Error("Copy() returned same pointer, want different")
	}
}

func TestCopySlice(t *testing.T) {
	original := []TestStruct{
		{Int: 1, String: "one"},
		{Int: 2, String: "two"},
	}

	copied := Copy[[]TestStruct](original)

	if !reflect.DeepEqual(copied, original) {
		t.Errorf("Copy() = %+v, want %+v", copied, original)
	}

	// 确保是深拷贝
	if &copied[0] == &original[0] {
		t.Error("Copy() returned slice with same element addresses")
	}
}

func TestCopyMap(t *testing.T) {
	original := map[string]TestStruct{
		"one": {Int: 1, String: "one"},
		"two": {Int: 2, String: "two"},
	}

	copied := Copy[map[string]TestStruct](original)

	if !reflect.DeepEqual(copied, original) {
		t.Errorf("Copy() = %+v, want %+v", copied, original)
	}

	// 确保是深拷贝
	originalValue := original["one"]
	copiedValue := copied["one"]
	if &originalValue == &copiedValue {
		t.Error("Copy() returned map with same element addresses")
	}
}

func TestCopyTime(t *testing.T) {
	original := time.Now()
	copied := Copy[time.Time](original)

	if !original.Equal(copied) {
		t.Errorf("Copy() = %v, want %v", copied, original)
	}
}

func TestCopyCustomCopier(t *testing.T) {
	original := CustomCopier{Value: 42}
	copied := Copy[CustomCopier](original)

	// 自定义拷贝会将值翻倍
	if copied.Value != original.Value*2 {
		t.Errorf("Copy() = %v, want %v", copied.Value, original.Value*2)
	}
}

func TestCopyNestedStruct(t *testing.T) {
	now := time.Now()
	original := NestedStruct{
		Basic:   TestStruct{Int: 1, String: "test"},
		Pointer: &TestStruct{Int: 2, String: "pointer"},
		Slice: []TestStruct{
			{Int: 3, String: "slice1"},
			{Int: 4, String: "slice2"},
		},
		Map: map[string]TestStruct{
			"key1": {Int: 5, String: "map1"},
			"key2": {Int: 6, String: "map2"},
		},
		Time:      now,
		Interface: &TestStruct{Int: 7, String: "interface"},
	}

	copied := Copy[NestedStruct](original)

	if !reflect.DeepEqual(copied, original) {
		t.Errorf("Copy() = %+v, want %+v", copied, original)
	}

	// 检查指针字段是否真的被深拷贝
	if copied.Pointer == original.Pointer {
		t.Error("Pointer field was not deep copied")
	}

	// 检查切片元素是否被深拷贝
	if &copied.Slice[0] == &original.Slice[0] {
		t.Error("Slice elements were not deep copied")
	}

	// 检查 map 值是否被深拷贝
	originalValue := original.Map["key1"]
	copiedValue := copied.Map["key1"]
	if &originalValue == &copiedValue {
		t.Error("Map values were not deep copied")
	}

	// 检查接口字段是否被深拷贝
	if reflect.ValueOf(copied.Interface).Pointer() == reflect.ValueOf(original.Interface).Pointer() {
		t.Error("Interface field was not deep copied")
	}
}

func TestCopyNilValues(t *testing.T) {
	var nilSlice []int
	var nilMap map[string]int
	var nilPtr *TestStruct
	var nilIface interface{}

	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			"nil slice",
			func(t *testing.T) {
				copied := Copy[[]int](nilSlice)
				if copied != nil {
					t.Errorf("Copy() = %v, want nil", copied)
				}
			},
		},
		{
			"nil map",
			func(t *testing.T) {
				copied := Copy[map[string]int](nilMap)
				if copied != nil {
					t.Errorf("Copy() = %v, want nil", copied)
				}
			},
		},
		{
			"nil pointer",
			func(t *testing.T) {
				copied := Copy[*TestStruct](nilPtr)
				if copied != nil {
					t.Errorf("Copy() = %v, want nil", copied)
				}
			},
		},
		{
			"nil interface",
			func(t *testing.T) {
				copied := Copy[interface{}](nilIface)
				if copied != nil {
					t.Errorf("Copy() = %v, want nil", copied)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}

func TestCopyCircularReference(t *testing.T) {
	type Node struct {
		Next  *Node
		Value int
	}

	original := &Node{Value: 1}
	original.Next = &Node{Value: 2}
	original.Next.Next = original // 创建循环引用

	copied := Copy[*Node](original)

	// 检查值是否正确复制
	if copied.Value != original.Value {
		t.Errorf("Copy() value = %v, want %v", copied.Value, original.Value)
	}
	if copied.Next.Value != original.Next.Value {
		t.Errorf("Copy() next value = %v, want %v", copied.Next.Value, original.Next.Value)
	}

	// 确保是深拷贝
	if copied == original {
		t.Error("Copy() returned same pointer for root node")
	}
	if copied.Next == original.Next {
		t.Error("Copy() returned same pointer for next node")
	}
}

// 基础类型测试用例
type basicTestCase[T any] struct {
	name     string
	input    T
	expected T
}

// 自定义拷贝接口测试
func TestCustomCopierInterface(t *testing.T) {
	original := CustomCopier{Value: 42}
	copied := Copy[CustomCopier](original)

	// CustomCopier.DeepCopy() 将值翻倍
	if copied.Value != original.Value*2 {
		t.Errorf("CustomCopier: got %v, want %v", copied.Value, original.Value*2)
	}
}

// 嵌套结构体测试
func TestNestedStructCopy(t *testing.T) {
	now := time.Now()
	original := NestedStruct{
		Basic: TestStruct{Int: 1, String: "basic"},
		Pointer: &TestStruct{
			Int:    2,
			String: "pointer",
		},
		Slice: []TestStruct{
			{Int: 3, String: "slice1"},
			{Int: 4, String: "slice2"},
		},
		Map: map[string]TestStruct{
			"key1": {Int: 5, String: "map1"},
			"key2": {Int: 6, String: "map2"},
		},
		Time:      now,
		Interface: "interface value",
	}

	copied := Copy[NestedStruct](original)

	// 验证基本字段
	if !reflect.DeepEqual(copied.Basic, original.Basic) {
		t.Error("Basic struct copy failed")
	}

	// 验证指针
	if copied.Pointer == original.Pointer {
		t.Error("Pointer was not deep copied")
	}
	if !reflect.DeepEqual(*copied.Pointer, *original.Pointer) {
		t.Error("Pointer value copy failed")
	}

	// 验证切片
	if len(copied.Slice) != len(original.Slice) {
		t.Error("Slice length mismatch")
	}
	for i := range original.Slice {
		if !reflect.DeepEqual(copied.Slice[i], original.Slice[i]) {
			t.Errorf("Slice element %d copy failed", i)
		}
	}

	// 验证 map
	if len(copied.Map) != len(original.Map) {
		t.Error("Map length mismatch")
	}
	for k, v := range original.Map {
		if !reflect.DeepEqual(copied.Map[k], v) {
			t.Errorf("Map element %s copy failed", k)
		}
	}

	// 验证时间
	if !copied.Time.Equal(original.Time) {
		t.Error("Time copy failed")
	}

	// 验证接口
	if copied.Interface != original.Interface {
		t.Error("Interface copy failed")
	}
}

// nil 值测试
func TestNilValues(t *testing.T) {
	var nilPtr *TestStruct
	var nilSlice []int
	var nilMap map[string]int
	var nilIface interface{}

	if Copy[*TestStruct](nilPtr) != nil {
		t.Error("nil pointer copy failed")
	}
	if Copy[[]int](nilSlice) != nil {
		t.Error("nil slice copy failed")
	}
	if Copy[map[string]int](nilMap) != nil {
		t.Error("nil map copy failed")
	}
	if Copy[interface{}](nilIface) != nil {
		t.Error("nil interface copy failed")
	}
}

// 循环引用测试
func TestCircularReference(t *testing.T) {
	// 创建一个循环链表
	original := &Node{Value: 1}
	original.Next = &Node{Value: 2}
	original.Next.Next = &Node{Value: 3}
	original.Next.Next.Next = original // 创建循环

	copied := Copy[*Node](original)

	// 验证值被正确复制
	if copied.Value != original.Value {
		t.Error("Node value copy failed")
	}
	if copied.Next.Value != original.Next.Value {
		t.Error("Next node value copy failed")
	}
	if copied.Next.Next.Value != original.Next.Next.Value {
		t.Error("Next next node value copy failed")
	}

	// 验证循环引用被正确处理
	if copied.Next.Next.Next != copied {
		t.Error("Circular reference was not properly copied")
	}

	// 验证是真正的深拷贝
	if copied == original {
		t.Error("Node was not deep copied")
	}
	if copied.Next == original.Next {
		t.Error("Next node was not deep copied")
	}
	if copied.Next.Next == original.Next.Next {
		t.Error("Next next node was not deep copied")
	}
}

// Node 用于测试循环引用
type Node struct {
	Next  *Node
	Value int
}
