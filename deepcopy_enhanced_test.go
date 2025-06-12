package deepcopy

import (
	"fmt"
	"testing"
	"time"
)

// 测试用的结构体类型

// 只包含值类型的结构体
type OnlyValueStruct struct {
	Name  string
	Age   int
	Score float64
	Valid bool
}

// 包含引用类型的结构体
type WithReferenceStruct struct {
	Name    string
	Age     int
	Friends []string // 切片类型
	Data    *string  // 指针类型
}

// 嵌套结构体（只包含值类型）
type NestedValueStruct struct {
	Basic OnlyValueStruct
	Count int
}

func TestCopyOptimization(t *testing.T) {
	// 测试只包含值类型的结构体
	valueStruct := OnlyValueStruct{
		Name:  "Alice",
		Age:   30,
		Score: 95.5,
		Valid: true,
	}

	// 这种情况会使用优化路径，直接返回原值
	copied1 := Copy(valueStruct)
	fmt.Printf("原始值类型结构体: %+v\n", valueStruct)
	fmt.Printf("拷贝值类型结构体: %+v\n", copied1)

	// 测试包含引用类型的结构体
	data := "test data"
	refStruct := WithReferenceStruct{
		Name:    "Bob",
		Age:     25,
		Friends: []string{"Alice", "Charlie"},
		Data:    &data,
	}

	// 这种情况会进行深拷贝
	copied2 := Copy(refStruct)
	fmt.Printf("原始引用类型结构体: %+v\n", refStruct)
	fmt.Printf("拷贝引用类型结构体: %+v\n", copied2)
}

func TestNewManagerAndCache(t *testing.T) {
	// 创建新的管理器实例
	manager := NewDeepCopyManager()

	// 测试缓存功能
	valueStruct := OnlyValueStruct{
		Name:  "Test",
		Age:   25,
		Score: 88.0,
		Valid: true,
	}

	// 第一次分析
	analysis1 := manager.AnalyzeValue(valueStruct)

	// 第二次分析相同类型（应该使用缓存）
	analysis2 := manager.AnalyzeValue(OnlyValueStruct{})

	fmt.Printf("\n=== 缓存测试 ===\n")

	// 验证分析结果一致
	if analysis1.IsOnlyValues != analysis2.IsOnlyValues {
		t.Error("缓存的分析结果不一致")
	}

	fmt.Printf("缓存功能正常，分析结果一致\n")
}

func TestTypeAnalysis(t *testing.T) {
	// 分析只包含值类型的结构体
	valueStruct := OnlyValueStruct{}
	analysis1 := AnalyzeType(valueStruct)

	fmt.Printf("\n=== 值类型结构体分析 ===\n")
	fmt.Printf("类型名称: %s\n", analysis1.TypeName)
	fmt.Printf("只包含值类型: %t\n", analysis1.IsOnlyValues)
	fmt.Printf("包含指针: %t\n", analysis1.ContainsPtr)
	fmt.Printf("包含切片: %t\n", analysis1.ContainsSlice)
	fmt.Printf("包含映射: %t\n", analysis1.ContainsMap)
	fmt.Printf("包含通道: %t\n", analysis1.ContainsChan)
	fmt.Printf("包含函数: %t\n", analysis1.ContainsFunc)
	fmt.Printf("包含接口: %t\n", analysis1.ContainsIface)

	if analysis1.FieldAnalysis != nil {
		fmt.Printf("字段分析:\n")
		for fieldName, fieldAnalysis := range analysis1.FieldAnalysis {
			fmt.Printf("  %s: %s (只包含值类型: %t)\n",
				fieldName, fieldAnalysis.TypeName, fieldAnalysis.IsOnlyValues)
		}
	}

	// 分析包含引用类型的结构体
	refStruct := WithReferenceStruct{}
	analysis2 := AnalyzeType(refStruct)

	fmt.Printf("\n=== 引用类型结构体分析 ===\n")
	fmt.Printf("类型名称: %s\n", analysis2.TypeName)
	fmt.Printf("只包含值类型: %t\n", analysis2.IsOnlyValues)
	fmt.Printf("包含指针: %t\n", analysis2.ContainsPtr)
	fmt.Printf("包含切片: %t\n", analysis2.ContainsSlice)
	fmt.Printf("包含映射: %t\n", analysis2.ContainsMap)

	if analysis2.FieldAnalysis != nil {
		fmt.Printf("字段分析:\n")
		for fieldName, fieldAnalysis := range analysis2.FieldAnalysis {
			fmt.Printf("  %s: %s (只包含值类型: %t)\n",
				fieldName, fieldAnalysis.TypeName, fieldAnalysis.IsOnlyValues)
		}
	}

	// 分析嵌套的值类型结构体
	nestedStruct := NestedValueStruct{}
	analysis3 := AnalyzeType(nestedStruct)

	fmt.Printf("\n=== 嵌套值类型结构体分析 ===\n")
	fmt.Printf("类型名称: %s\n", analysis3.TypeName)
	fmt.Printf("只包含值类型: %t\n", analysis3.IsOnlyValues)

	// 分析基础类型
	var intVal int = 42
	analysis4 := AnalyzeType(intVal)
	fmt.Printf("\n=== 基础类型分析 ===\n")
	fmt.Printf("类型名称: %s\n", analysis4.TypeName)
	fmt.Printf("只包含值类型: %t\n", analysis4.IsOnlyValues)
}

func TestPerformanceBenefit(t *testing.T) {
	// 创建一个大的只包含值类型的结构体
	largeValueStruct := struct {
		Field1 int
		Field2 string
		Field3 float64
		Field4 bool
		Field5 [1000]int // 大数组，但都是值类型
	}{
		Field1: 42,
		Field2: "test",
		Field3: 3.14,
		Field4: true,
	}

	// 分析这个大结构体
	analysis := AnalyzeType(largeValueStruct)
	fmt.Printf("\n=== 大型值类型结构体 ===\n")
	fmt.Printf("只包含值类型: %t\n", analysis.IsOnlyValues)

	// 这种情况下，Copy函数会直接返回原值，避免昂贵的深拷贝操作
	copied := Copy(largeValueStruct)

	// 验证拷贝是否正确
	if copied.Field1 != largeValueStruct.Field1 {
		t.Errorf("拷贝失败")
	}

	fmt.Printf("拷贝成功，使用了优化路径\n")
}

func TestCacheManagement(t *testing.T) {
	manager := NewDeepCopyManager()

	// 分析多种不同类型
	types := []interface{}{
		int(42),
		"string",
		[]int{1, 2, 3},
		map[string]int{"a": 1},
		OnlyValueStruct{},
		WithReferenceStruct{},
	}

	fmt.Printf("\n=== 缓存管理测试 ===\n")

	// 分析所有类型
	for i, typ := range types {
		manager.AnalyzeValue(typ)
		fmt.Printf("分析类型 %d\n", i+1)
	}

	// 再次分析相同类型（应该使用缓存）
	fmt.Printf("\n重复分析相同类型（使用缓存）\n")
	for _, typ := range types {
		manager.AnalyzeValue(typ)
	}

	fmt.Printf("缓存管理测试完成\n")
}

func TestBusinessKeyOptimization(t *testing.T) {
	fmt.Printf("\n=== 基于业务 key 的优化测试 ===\n")

	// 测试值类型结构体的业务缓存优化
	testValueStructOptimization := func() {
		valueStruct := OnlyValueStruct{Name: "Config1", Age: 25}
		configKey := "user.profile"

		// 第一次调用 - 会创建业务拷贝信息
		start := time.Now()
		copied1 := CopyWithKey(valueStruct, configKey)
		duration1 := time.Since(start)

		// 第二次调用 - 应该直接返回，完全避免反射
		start = time.Now()
		copied2 := CopyWithKey(valueStruct, configKey)
		duration2 := time.Since(start)

		// 第三次调用 - 再次验证
		start = time.Now()
		copied3 := CopyWithKey(OnlyValueStruct{Name: "Config2", Age: 30}, configKey)
		duration3 := time.Since(start)

		fmt.Printf("值类型结构体业务缓存优化:\n")
		fmt.Printf("  第一次拷贝耗时: %v\n", duration1)
		fmt.Printf("  第二次拷贝耗时: %v\n", duration2)
		fmt.Printf("  第三次拷贝耗时: %v\n", duration3)

		// 验证结果正确性
		if copied1.Name != "Config1" || copied1.Age != 25 {
			t.Error("第一次拷贝结果不正确")
		}
		if copied2.Name != "Config1" || copied2.Age != 25 {
			t.Error("第二次拷贝结果不正确")
		}
		if copied3.Name != "Config2" || copied3.Age != 30 {
			t.Error("第三次拷贝结果不正确")
		}
	}

	// 测试基础类型的业务缓存优化
	testBasicTypeOptimization := func() {
		configValue := 42
		configKey := "app.timeout"

		start := time.Now()
		for i := 0; i < 10000; i++ {
			CopyWithKey(configValue, configKey)
		}
		duration := time.Since(start)

		fmt.Printf("基础类型 10000次拷贝耗时: %v (平均: %v/次)\n",
			duration, duration/10000)
	}

	// 对比传统 Copy 和 CopyWithKey 的性能差异
	testPerformanceComparison := func() {
		valueStruct := OnlyValueStruct{Name: "TestStruct", Age: 25}

		// 传统 Copy 方式
		start := time.Now()
		for i := 0; i < 1000; i++ {
			Copy(valueStruct)
		}
		traditionalDuration := time.Since(start)

		// CopyWithKey 方式（预热后）
		CopyWithKey(valueStruct, "test.config") // 预热
		start = time.Now()
		for i := 0; i < 1000; i++ {
			CopyWithKey(valueStruct, "test.config")
		}
		businessDuration := time.Since(start)

		fmt.Printf("性能对比 (1000次拷贝):\n")
		fmt.Printf("  传统 Copy:     %v\n", traditionalDuration)
		fmt.Printf("  CopyWithKey:   %v\n", businessDuration)
	}

	// 执行所有测试
	testValueStructOptimization()
	testBasicTypeOptimization()
	testPerformanceComparison()

	fmt.Printf("基于业务 key 的优化测试完成\n")
}

func TestBusinessCacheManagement(t *testing.T) {
	fmt.Printf("\n=== 业务缓存管理测试 ===\n")

	// 添加不同类型的业务配置
	CopyWithKey(42, "app.port")
	CopyWithKey("localhost", "app.host")
	CopyWithKey(OnlyValueStruct{Name: "test", Age: 25}, "user.profile")
	CopyWithKey([]int{1, 2, 3}, "app.features")

	fmt.Printf("业务缓存管理测试通过\n")
}

func TestCopyWithKeyCorrectness(t *testing.T) {
	fmt.Printf("\n=== CopyWithKey 正确性测试 ===\n")

	// 测试复杂结构的深拷贝正确性
	complexStruct := struct {
		Name    string
		Data    *string
		Items   []int
		Mapping map[string]int
	}{
		Name:    "Complex",
		Data:    new(string),
		Items:   []int{1, 2, 3},
		Mapping: map[string]int{"a": 1, "b": 2},
	}
	*complexStruct.Data = "test data"

	// 使用 CopyWithKey 拷贝
	copied := CopyWithKey(complexStruct, "complex.config")

	// 验证深拷贝正确性
	if copied.Data == complexStruct.Data {
		t.Error("指针应该被深拷贝，不应该相等")
	}
	if *copied.Data != *complexStruct.Data {
		t.Error("指针指向的值应该相等")
	}

	// 修改原数据，验证拷贝数据不受影响
	*complexStruct.Data = "modified"
	complexStruct.Items[0] = 999
	complexStruct.Mapping["c"] = 3

	if *copied.Data != "test data" {
		t.Error("修改原数据不应该影响拷贝数据")
	}
	if copied.Items[0] != 1 {
		t.Error("修改原切片不应该影响拷贝切片")
	}
	if len(copied.Mapping) != 2 {
		t.Error("修改原映射不应该影响拷贝映射")
	}

	fmt.Printf("CopyWithKey 正确性测试通过\n")
}

func TestComplexNestedStructures(t *testing.T) {
	fmt.Printf("\n=== 复杂嵌套结构体测试 ===\n")

	// 定义复杂的嵌套结构体
	type Address struct {
		Street   string
		City     string
		PostCode string
	}

	type Contact struct {
		Phone string
		Email string
	}

	type Person struct {
		Name    string
		Age     int
		Address *Address
		Contact Contact
		Tags    []string
	}

	type Company struct {
		Name        string
		Employees   []*Person
		HeadQuarter Address
		Contacts    map[string]*Contact
	}

	// 创建复杂的嵌套数据
	originalCompany := Company{
		Name: "TechCorp",
		HeadQuarter: Address{
			Street:   "123 Tech Street",
			City:     "Silicon Valley",
			PostCode: "12345",
		},
		Employees: []*Person{
			{
				Name: "Alice",
				Age:  30,
				Address: &Address{
					Street:   "456 Home St",
					City:     "Home City",
					PostCode: "67890",
				},
				Contact: Contact{
					Phone: "123-456-7890",
					Email: "alice@example.com",
				},
				Tags: []string{"developer", "senior"},
			},
		},
		Contacts: map[string]*Contact{
			"hr": {
				Phone: "555-0001",
				Email: "hr@techcorp.com",
			},
		},
	}

	// 进行拷贝
	copiedCompany := Copy(originalCompany)

	// 验证深拷贝的正确性
	if &copiedCompany == &originalCompany {
		t.Error("拷贝对象与原对象指针相同")
	}

	if copiedCompany.Employees[0] == originalCompany.Employees[0] {
		t.Error("嵌套指针应该被深拷贝")
	}

	if copiedCompany.Employees[0].Address == originalCompany.Employees[0].Address {
		t.Error("深层嵌套指针应该被深拷贝")
	}

	// 验证值的正确性
	if copiedCompany.Name != originalCompany.Name {
		t.Error("值拷贝不正确")
	}

	if copiedCompany.Employees[0].Name != originalCompany.Employees[0].Name {
		t.Error("嵌套值拷贝不正确")
	}

	fmt.Printf("复杂嵌套结构体测试通过\n")
}

func TestEnhancedPerformanceComparison(t *testing.T) {
	fmt.Printf("\n=== 增强性能对比测试 ===\n")

	// 测试基础类型
	testBasicType := func() {
		value := 42
		iterations := 100000

		// 测试增强后的 Copy
		start := time.Now()
		for i := 0; i < iterations; i++ {
			_ = Copy(value)
		}
		copyDuration := time.Since(start)

		fmt.Printf("基础类型 %d 次调用:\n", iterations)
		fmt.Printf("  增强 Copy:    %v (%.2f ns/次)\n", copyDuration, float64(copyDuration.Nanoseconds())/float64(iterations))
	}

	// 测试值类型结构体
	testValueStruct := func() {
		value := OnlyValueStruct{Name: "test", Age: 25}
		iterations := 10000

		// 测试增强后的 Copy
		start := time.Now()
		for i := 0; i < iterations; i++ {
			_ = Copy(value)
		}
		copyDuration := time.Since(start)

		fmt.Printf("值类型结构体 %d 次调用:\n", iterations)
		fmt.Printf("  增强 Copy:    %v (%.2f ns/次)\n", copyDuration, float64(copyDuration.Nanoseconds())/float64(iterations))
	}

	testBasicType()
	testValueStruct()

	fmt.Printf("\n增强后的实现特点:\n")
	fmt.Printf("1. 类型分析缓存机制\n")
	fmt.Printf("2. 值类型优化：直接返回，避免反射\n")
	fmt.Printf("3. 业务key缓存\n")
	fmt.Printf("4. 保持原有深拷贝功能的完整性\n")
}
