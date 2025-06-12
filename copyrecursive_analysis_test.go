package deepcopy

import (
	"reflect"
	"runtime"
	"testing"
)

// 测试结构体定义
type ComplexStruct struct {
	// 基本类型
	Int    int
	String string
	Bool   bool
	Float  float64

	// 复合类型
	Slice []int
	Map   map[string]int
	Array [3]int

	// 指针类型
	IntPtr    *int
	StringPtr *string

	// 嵌套结构体
	Nested *ComplexStruct

	// 接口类型
	Interface interface{}

	// 未导出字段
	private string
}

// 循环引用结构体
type CircularStruct struct {
	Name string
	Next *CircularStruct
	Prev *CircularStruct
}

// 自定义拷贝结构体
type CustomCopyStruct struct {
	Value int
}

func (c CustomCopyStruct) DeepCopy() CustomCopyStruct {
	return CustomCopyStruct{Value: c.Value + 100}
}

type CustomCopyPtrStruct struct {
	Value int
}

func (c *CustomCopyPtrStruct) DeepCopy() CustomCopyPtrStruct {
	return CustomCopyPtrStruct{Value: c.Value + 200}
}

// 1. 基本类型正确性测试
func TestCopyRecursive_BasicTypes(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
	}{
		{"int", 42},
		{"string", "hello"},
		{"bool", true},
		{"float64", 3.14},
		{"byte", byte(255)},
		{"rune", '中'},
		{"complex128", complex(1, 2)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			original := reflect.ValueOf(tt.input)
			cpy := reflect.New(original.Type()).Elem()
			visited := make(map[uintptr]reflect.Value)

			copyRecursive(original, cpy, visited)

			if !reflect.DeepEqual(original.Interface(), cpy.Interface()) {
				t.Errorf("Basic type copy failed: got %v, want %v", cpy.Interface(), original.Interface())
			}
		})
	}
}

// 2. 指针类型安全性测试
func TestCopyRecursive_PointerSafety(t *testing.T) {
	// 测试指针拷贝后的地址不同
	original := 42
	ptr := &original

	originalVal := reflect.ValueOf(ptr)
	cpy := reflect.New(originalVal.Type()).Elem()
	visited := make(map[uintptr]reflect.Value)

	copyRecursive(originalVal, cpy, visited)

	copiedPtr := cpy.Interface().(*int)

	// 检查地址不同
	if ptr == copiedPtr {
		t.Error("Pointer addresses should be different")
	}

	// 检查值相同
	if *ptr != *copiedPtr {
		t.Errorf("Pointer values should be same: got %v, want %v", *copiedPtr, *ptr)
	}

	// 修改原始值，确保拷贝不受影响
	*ptr = 100
	if *copiedPtr == 100 {
		t.Error("Copy should not be affected by original modification")
	}
}

// 3. 循环引用处理测试
func TestCopyRecursive_CircularReference(t *testing.T) {
	// 创建循环引用
	a := &CircularStruct{Name: "A"}
	b := &CircularStruct{Name: "B"}
	c := &CircularStruct{Name: "C"}

	a.Next = b
	b.Next = c
	c.Next = a // 创建循环

	a.Prev = c
	b.Prev = a
	c.Prev = b

	// 拷贝
	originalVal := reflect.ValueOf(a)
	cpy := reflect.New(originalVal.Type()).Elem()
	visited := make(map[uintptr]reflect.Value)

	copyRecursive(originalVal, cpy, visited)

	copiedA := cpy.Interface().(*CircularStruct)

	// 验证循环引用被正确处理
	if copiedA.Next.Next.Next != copiedA {
		t.Error("Circular reference not properly maintained")
	}

	// 验证是深拷贝
	if copiedA == a {
		t.Error("Should be different instances")
	}

	// 验证值正确
	if copiedA.Name != "A" || copiedA.Next.Name != "B" || copiedA.Next.Next.Name != "C" {
		t.Error("Values not copied correctly")
	}
}

// 4. nil值处理测试
func TestCopyRecursive_NilHandling(t *testing.T) {
	tests := []struct {
		name  string
		setup func() reflect.Value
	}{
		{
			"nil pointer",
			func() reflect.Value {
				var ptr *int
				return reflect.ValueOf(&ptr).Elem()
			},
		},
		{
			"nil slice",
			func() reflect.Value {
				var slice []int
				return reflect.ValueOf(&slice).Elem()
			},
		},
		{
			"nil map",
			func() reflect.Value {
				var m map[string]int
				return reflect.ValueOf(&m).Elem()
			},
		},
		{
			"nil interface",
			func() reflect.Value {
				var iface interface{}
				return reflect.ValueOf(&iface).Elem()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			original := tt.setup()
			cpy := reflect.New(original.Type()).Elem()
			visited := make(map[uintptr]reflect.Value)

			// 这不应该panic
			copyRecursive(original, cpy, visited)

			// 验证nil值被正确处理
			if !cpy.IsNil() {
				t.Errorf("Nil value should remain nil, got %v", cpy.Interface())
			}
		})
	}
}

// 5. 复杂结构体测试
func TestCopyRecursive_ComplexStruct(t *testing.T) {
	intVal := 42
	stringVal := "test"

	original := &ComplexStruct{
		Int:       100,
		String:    "hello",
		Bool:      true,
		Float:     3.14,
		Slice:     []int{1, 2, 3},
		Map:       map[string]int{"a": 1, "b": 2},
		Array:     [3]int{10, 20, 30},
		IntPtr:    &intVal,
		StringPtr: &stringVal,
		Interface: "interface value",
		private:   "should not copy",
	}

	// 添加嵌套结构
	original.Nested = &ComplexStruct{
		Int:    200,
		String: "nested",
	}

	originalVal := reflect.ValueOf(original)
	cpy := reflect.New(originalVal.Type()).Elem()
	visited := make(map[uintptr]reflect.Value)

	copyRecursive(originalVal, cpy, visited)

	copied := cpy.Interface().(*ComplexStruct)

	// 验证基本字段
	if copied.Int != original.Int {
		t.Errorf("Int: got %v, want %v", copied.Int, original.Int)
	}
	if copied.String != original.String {
		t.Errorf("String: got %v, want %v", copied.String, original.String)
	}

	// 验证指针字段是深拷贝
	if copied.IntPtr == original.IntPtr {
		t.Error("IntPtr should be different addresses")
	}
	if *copied.IntPtr != *original.IntPtr {
		t.Errorf("IntPtr value: got %v, want %v", *copied.IntPtr, *original.IntPtr)
	}

	// 验证slice是深拷贝
	if len(copied.Slice) != len(original.Slice) {
		t.Error("Slice length mismatch")
	}
	for i, v := range original.Slice {
		if copied.Slice[i] != v {
			t.Errorf("Slice[%d]: got %v, want %v", i, copied.Slice[i], v)
		}
	}

	// 验证map是深拷贝
	if len(copied.Map) != len(original.Map) {
		t.Error("Map length mismatch")
	}
	for k, v := range original.Map {
		if copied.Map[k] != v {
			t.Errorf("Map[%s]: got %v, want %v", k, copied.Map[k], v)
		}
	}

	// 验证嵌套结构是深拷贝
	if copied.Nested == original.Nested {
		t.Error("Nested should be different addresses")
	}
	if copied.Nested.Int != original.Nested.Int {
		t.Errorf("Nested.Int: got %v, want %v", copied.Nested.Int, original.Nested.Int)
	}
}

// 6. 自定义拷贝方法测试
func TestCopyRecursive_CustomCopy(t *testing.T) {
	// 值接收者
	original1 := CustomCopyStruct{Value: 10}
	originalVal1 := reflect.ValueOf(original1)
	cpy1 := reflect.New(originalVal1.Type()).Elem()
	visited1 := make(map[uintptr]reflect.Value)

	copyRecursive(originalVal1, cpy1, visited1)

	copied1 := cpy1.Interface().(CustomCopyStruct)
	if copied1.Value != 110 { // 10 + 100
		t.Errorf("Custom copy (value receiver) failed: got %v, want %v", copied1.Value, 110)
	}

	// 指针接收者
	original2 := &CustomCopyPtrStruct{Value: 20}
	originalVal2 := reflect.ValueOf(original2)
	cpy2 := reflect.New(originalVal2.Type()).Elem()
	visited2 := make(map[uintptr]reflect.Value)

	copyRecursive(originalVal2, cpy2, visited2)

	copied2 := cpy2.Interface().(*CustomCopyPtrStruct)
	if copied2.Value != 220 { // 20 + 200
		t.Errorf("Custom copy (pointer receiver) failed: got %v, want %v", copied2.Value, 220)
	}
}

// 7. 内存使用分析
func TestCopyRecursive_MemoryUsage(t *testing.T) {
	// 创建一个大的结构体
	large := make([][]int, 100)
	for i := range large {
		large[i] = make([]int, 100)
		for j := range large[i] {
			large[i][j] = i*100 + j
		}
	}

	originalVal := reflect.ValueOf(large)
	cpy := reflect.New(originalVal.Type()).Elem()
	visited := make(map[uintptr]reflect.Value)

	// 测量内存使用前后
	var m1, m2 runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m1)

	copyRecursive(originalVal, cpy, visited)

	runtime.GC()
	runtime.ReadMemStats(&m2)

	copied := cpy.Interface().([][]int)

	// 验证拷贝正确性
	if len(copied) != len(large) {
		t.Error("Large structure copy failed")
	}

	// 验证是深拷贝
	if &copied[0][0] == &large[0][0] {
		t.Error("Should be different memory addresses")
	}

	// 打印内存使用情况（仅用于观察）
	t.Logf("Memory before: %d bytes", m1.Alloc)
	t.Logf("Memory after: %d bytes", m2.Alloc)
	t.Logf("Memory increase: %d bytes", m2.Alloc-m1.Alloc)
}

// 8. 边界条件测试
func TestCopyRecursive_EdgeCases(t *testing.T) {
	t.Run("empty slice", func(t *testing.T) {
		original := []int{}
		originalVal := reflect.ValueOf(original)
		cpy := reflect.New(originalVal.Type()).Elem()
		visited := make(map[uintptr]reflect.Value)

		copyRecursive(originalVal, cpy, visited)

		copied := cpy.Interface().([]int)
		if len(copied) != 0 {
			t.Error("Empty slice should remain empty")
		}
	})

	t.Run("empty map", func(t *testing.T) {
		original := make(map[string]int)
		originalVal := reflect.ValueOf(original)
		cpy := reflect.New(originalVal.Type()).Elem()
		visited := make(map[uintptr]reflect.Value)

		copyRecursive(originalVal, cpy, visited)

		copied := cpy.Interface().(map[string]int)
		if len(copied) != 0 {
			t.Error("Empty map should remain empty")
		}
	})

	t.Run("zero value struct", func(t *testing.T) {
		original := ComplexStruct{}
		originalVal := reflect.ValueOf(original)
		cpy := reflect.New(originalVal.Type()).Elem()
		visited := make(map[uintptr]reflect.Value)

		copyRecursive(originalVal, cpy, visited)

		copied := cpy.Interface().(ComplexStruct)
		if !reflect.DeepEqual(copied, original) {
			t.Error("Zero value struct should be copied correctly")
		}
	})
}
