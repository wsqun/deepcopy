package deepcopy

import (
	"reflect"
	"sync"
	"time"
)

// Copier 是一个可以自定义深拷贝行为的接口
type Copier[T any] interface {
	DeepCopy() T
}

// DeepCopyManager 深拷贝管理器，提供类型分析和深拷贝功能
// 使用缓存机制优化性能，避免重复的反射分析
type DeepCopyManager struct {
	// 类型分析结果缓存，key: reflect.Type, value: *TypeAnalysisResult
	analysisCache sync.Map
}

// TypeAnalysisResult 类型分析结果，包含所有必要的信息
type TypeAnalysisResult struct {
	IsOnlyValues  bool                           // 是否只包含值类型
	ContainsPtr   bool                           // 是否包含指针
	ContainsSlice bool                           // 是否包含切片
	ContainsMap   bool                           // 是否包含映射
	ContainsChan  bool                           // 是否包含通道
	ContainsFunc  bool                           // 是否包含函数
	ContainsIface bool                           // 是否包含接口
	FieldAnalysis map[string]*TypeAnalysisResult // 结构体字段分析（仅当类型为结构体时）
	TypeName      string                         // 类型名称
}

// BusinessCopyInfo 业务拷贝信息，基于配置 key 缓存的优化信息
type BusinessCopyInfo struct {
	IsOnlyValues   bool                // 是否只包含值类型
	analysisResult *TypeAnalysisResult // 类型分析结果
	rtype          reflect.Type        // 反射类型信息
	once           sync.Once           // 确保只初始化一次
}

// TypedCopyManager 泛型层面的拷贝管理器，为每个具体类型缓存分析结果
type TypedCopyManager[T any] struct {
	analysis *TypeAnalysisResult // 类型分析结果
	rtype    reflect.Type        // 反射类型信息
	once     sync.Once           // 确保只分析一次
}

// 全局默认管理器实例
var defaultManager = NewDeepCopyManager()

// 全局业务拷贝信息缓存，key: 业务配置key, value: *BusinessCopyInfo
var businessCopyCache sync.Map // map[string]*BusinessCopyInfo

// 全局的泛型管理器缓存
var typedManagers sync.Map // map[reflect.Type]*TypedCopyManager[any]

// NewDeepCopyManager 创建新的深拷贝管理器
func NewDeepCopyManager() *DeepCopyManager {
	return &DeepCopyManager{}
}

// getTypedManager 获取或创建特定类型的管理器
func getTypedManager[T any]() *TypedCopyManager[T] {
	var zero T
	rtype := reflect.TypeOf(zero)

	// 处理 nil 类型的特殊情况
	if rtype == nil {
		// 为 nil 类型创建一个特殊的管理器
		manager := &TypedCopyManager[T]{
			rtype: nil,
		}
		return manager
	}

	// 尝试从缓存获取
	if cached, ok := typedManagers.Load(rtype); ok {
		return cached.(*TypedCopyManager[T])
	}

	// 创建新的管理器
	manager := &TypedCopyManager[T]{
		rtype: rtype,
	}

	// 存入缓存
	typedManagers.Store(rtype, manager)

	return manager
}

// getOrAnalyzeType 获取或分析类型结果（使用 sync.Once 确保只分析一次）
func (tm *TypedCopyManager[T]) getOrAnalyzeType() *TypeAnalysisResult {
	tm.once.Do(func() {
		// 处理 nil 类型的特殊情况
		if tm.rtype == nil {
			tm.analysis = &TypeAnalysisResult{
				TypeName:     "nil",
				IsOnlyValues: true,
			}
			return
		}

		// 分析类型
		tm.analysis = defaultManager.getOrAnalyzeType(tm.rtype)
	})
	return tm.analysis
}

// hasDeepCopyMethod 检查值是否有 DeepCopy 方法
func hasDeepCopyMethod(v reflect.Value) (reflect.Method, bool) {
	if !v.IsValid() {
		return reflect.Method{}, false
	}

	method, found := v.Type().MethodByName("DeepCopy")
	if found && method.Func.IsValid() {
		// 检查方法签名：应该没有参数（除了接收者）且有一个返回值
		methodType := method.Type
		if methodType.NumIn() == 1 && methodType.NumOut() == 1 {
			return method, true
		}
	}

	return reflect.Method{}, false
}

// callDeepCopy 调用 DeepCopy 方法
func callDeepCopy(v reflect.Value, method reflect.Method) reflect.Value {
	results := method.Func.Call([]reflect.Value{v})
	if len(results) > 0 {
		return results[0]
	}
	return reflect.Value{}
}

// Copy 创建任意值的深拷贝并返回副本
// 如果类型实现了 DeepCopy 方法，将使用其自定义的拷贝方法
// 使用类型分析优化：对于只包含值类型的数据直接返回，避免昂贵的深拷贝操作
func Copy[T any](src T) T {
	// 处理零值情况
	srcVal := reflect.ValueOf(src)
	if !srcVal.IsValid() {
		var zero T
		return zero
	}

	// 首先检查是否有 DeepCopy 方法
	if method, found := hasDeepCopyMethod(srcVal); found {
		result := callDeepCopy(srcVal, method)
		if result.IsValid() {
			return result.Interface().(T)
		}
	}

	// 获取该类型的专用管理器
	manager := getTypedManager[T]()

	// 获取类型分析结果（只会分析一次）
	analysis := manager.getOrAnalyzeType()

	// 性能优化：如果只包含值类型，直接返回原值
	if analysis.IsOnlyValues {
		return src
	}

	// 需要深拷贝的情况，使用反射方式
	result := defaultManager.CopyValue(src)
	return result.(T)
}

// CopyWithKey 基于业务 key 的优化拷贝，避免重复反射调用
// 这个函数的核心目的是缓存反射类型信息，减少每次调用时的反射开销
func CopyWithKey[T any](src T, key string) T {
	// 获取或创建业务拷贝信息
	copyInfo := getOrCreateBusinessCopyInfo[T](key)

	// 性能优化：如果只包含值类型，直接返回原值，完全避免反射
	if copyInfo.IsOnlyValues {
		return src
	}

	// 需要深拷贝的情况，使用缓存的反射信息进行高效拷贝
	srcVal := reflect.ValueOf(src)
	if !srcVal.IsValid() {
		var zero T
		return zero
	}

	// 首先检查是否有自定义 DeepCopy 方法（这个检查很快，不影响缓存效果）
	if method, found := hasDeepCopyMethod(srcVal); found {
		result := callDeepCopy(srcVal, method)
		if result.IsValid() {
			return result.Interface().(T)
		}
	}

	// 使用缓存的类型信息进行深拷贝
	cpy := reflect.New(srcVal.Type()).Elem()
	visited := make(map[uintptr]reflect.Value)
	copyRecursiveWithCache(srcVal, cpy, visited, copyInfo.analysisResult)

	return cpy.Interface().(T)
}

// AnalyzeType 使用默认管理器分析类型
func AnalyzeType[T any](src T) *TypeAnalysisResult {
	return defaultManager.AnalyzeValue(src)
}

// CopyValue 执行深拷贝操作（非泛型方法）
func (m *DeepCopyManager) CopyValue(src interface{}) interface{} {
	// 获取源数据的反射值对象
	srcVal := reflect.ValueOf(src)

	// 检查反射值是否有效
	if !srcVal.IsValid() {
		return nil
	}

	// 获取类型分析结果（使用缓存）
	analysis := m.getOrAnalyzeType(srcVal.Type())

	// 性能优化：如果只包含值类型，直接返回原值
	if analysis.IsOnlyValues {
		return src
	}

	// 首先检查是否有 DeepCopy 方法
	if method, found := hasDeepCopyMethod(srcVal); found {
		result := callDeepCopy(srcVal, method)
		if result.IsValid() {
			return result.Interface()
		}
	}

	// 创建目标反射值对象
	cpy := reflect.New(srcVal.Type()).Elem()

	// 创建访问记录映射，处理循环引用
	visited := make(map[uintptr]reflect.Value)

	// 执行深拷贝
	copyRecursive(srcVal, cpy, visited)

	// 返回结果
	return cpy.Interface()
}

// AnalyzeValue 分析给定值的类型结构（非泛型方法）
func (m *DeepCopyManager) AnalyzeValue(src interface{}) *TypeAnalysisResult {
	t := reflect.TypeOf(src)
	if t == nil {
		return &TypeAnalysisResult{
			TypeName:     "nil",
			IsOnlyValues: true,
		}
	}

	return m.getOrAnalyzeType(t)
}

// getOrAnalyzeType 获取或分析类型，使用缓存机制
func (m *DeepCopyManager) getOrAnalyzeType(t reflect.Type) *TypeAnalysisResult {
	// 尝试从缓存获取
	if cached, ok := m.analysisCache.Load(t); ok {
		return cached.(*TypeAnalysisResult)
	}

	// 缓存未命中，进行分析
	result := m.analyzeTypeRecursive(t, make(map[reflect.Type]*TypeAnalysisResult))

	// 存入缓存
	m.analysisCache.Store(t, result)

	return result
}

// analyzeTypeRecursive 递归分析类型结构
func (m *DeepCopyManager) analyzeTypeRecursive(t reflect.Type, visited map[reflect.Type]*TypeAnalysisResult) *TypeAnalysisResult {
	// 检查循环引用
	if result, ok := visited[t]; ok {
		return result
	}

	// 创建结果对象
	result := &TypeAnalysisResult{
		TypeName: t.String(),
	}

	// 先放入visited，防止循环引用
	visited[t] = result

	// 根据类型进行分析
	switch t.Kind() {
	// 基础值类型
	case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128, reflect.String:
		result.IsOnlyValues = true

	// 数组类型
	case reflect.Array:
		elemResult := m.analyzeTypeRecursive(t.Elem(), visited)
		result.IsOnlyValues = elemResult.IsOnlyValues
		result.ContainsPtr = elemResult.ContainsPtr
		result.ContainsSlice = elemResult.ContainsSlice
		result.ContainsMap = elemResult.ContainsMap
		result.ContainsChan = elemResult.ContainsChan
		result.ContainsFunc = elemResult.ContainsFunc
		result.ContainsIface = elemResult.ContainsIface

	// 结构体类型
	case reflect.Struct:
		result.IsOnlyValues = true // 假设是值类型，遇到引用类型时修改
		result.FieldAnalysis = make(map[string]*TypeAnalysisResult)

		// 检查是否有未导出字段，如果有则不能使用值拷贝优化
		hasUnexportedFields := false
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			if field.PkgPath != "" {
				hasUnexportedFields = true
				break
			}
		}

		// 如果有未导出字段，则必须进行深拷贝以确保正确跳过这些字段
		if hasUnexportedFields {
			result.IsOnlyValues = false
		}

		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)

			// 跳过未导出字段
			if field.PkgPath != "" {
				continue
			}

			// 分析字段类型
			fieldResult := m.analyzeTypeRecursive(field.Type, visited)
			result.FieldAnalysis[field.Name] = fieldResult

			// 更新结构体的整体分析结果
			if !fieldResult.IsOnlyValues {
				result.IsOnlyValues = false
			}
			if fieldResult.ContainsPtr {
				result.ContainsPtr = true
			}
			if fieldResult.ContainsSlice {
				result.ContainsSlice = true
			}
			if fieldResult.ContainsMap {
				result.ContainsMap = true
			}
			if fieldResult.ContainsChan {
				result.ContainsChan = true
			}
			if fieldResult.ContainsFunc {
				result.ContainsFunc = true
			}
			if fieldResult.ContainsIface {
				result.ContainsIface = true
			}
		}

	// 引用类型
	case reflect.Ptr:
		result.IsOnlyValues = false
		result.ContainsPtr = true
		// 递归分析指针指向的类型
		elemResult := m.analyzeTypeRecursive(t.Elem(), visited)
		result.ContainsSlice = elemResult.ContainsSlice
		result.ContainsMap = elemResult.ContainsMap
		result.ContainsChan = elemResult.ContainsChan
		result.ContainsFunc = elemResult.ContainsFunc
		result.ContainsIface = elemResult.ContainsIface

	case reflect.Slice:
		result.IsOnlyValues = false
		result.ContainsSlice = true
		// 递归分析切片元素类型
		elemResult := m.analyzeTypeRecursive(t.Elem(), visited)
		result.ContainsPtr = elemResult.ContainsPtr
		result.ContainsMap = elemResult.ContainsMap
		result.ContainsChan = elemResult.ContainsChan
		result.ContainsFunc = elemResult.ContainsFunc
		result.ContainsIface = elemResult.ContainsIface

	case reflect.Map:
		result.IsOnlyValues = false
		result.ContainsMap = true
		// 分析键和值的类型
		keyResult := m.analyzeTypeRecursive(t.Key(), visited)
		valueResult := m.analyzeTypeRecursive(t.Elem(), visited)
		result.ContainsPtr = keyResult.ContainsPtr || valueResult.ContainsPtr
		result.ContainsSlice = keyResult.ContainsSlice || valueResult.ContainsSlice
		result.ContainsChan = keyResult.ContainsChan || valueResult.ContainsChan
		result.ContainsFunc = keyResult.ContainsFunc || valueResult.ContainsFunc
		result.ContainsIface = keyResult.ContainsIface || valueResult.ContainsIface

	case reflect.Chan:
		result.IsOnlyValues = false
		result.ContainsChan = true

	case reflect.Func:
		result.IsOnlyValues = false
		result.ContainsFunc = true

	case reflect.Interface:
		result.IsOnlyValues = false
		result.ContainsIface = true

	// 其他未知类型
	default:
		result.IsOnlyValues = false
	}

	return result
}

// getOrCreateBusinessCopyInfo 获取或创建业务拷贝信息
func getOrCreateBusinessCopyInfo[T any](key string) *BusinessCopyInfo {
	// 尝试从缓存获取
	if cached, ok := businessCopyCache.Load(key); ok {
		return cached.(*BusinessCopyInfo)
	}

	// 创建新的业务拷贝信息
	var zero T
	rtype := reflect.TypeOf(zero)

	copyInfo := &BusinessCopyInfo{
		rtype: rtype,
	}

	// 初始化拷贝信息
	copyInfo.once.Do(func() {
		copyInfo.initializeCopyInfo()
	})

	// 存入缓存
	businessCopyCache.Store(key, copyInfo)

	return copyInfo
}

// initializeCopyInfo 初始化拷贝信息
func (info *BusinessCopyInfo) initializeCopyInfo() {
	// 处理 nil 类型
	if info.rtype == nil {
		info.IsOnlyValues = true
		return
	}

	// 分析类型
	info.analysisResult = defaultManager.getOrAnalyzeType(info.rtype)
	info.IsOnlyValues = info.analysisResult.IsOnlyValues
}

// copyRecursive 使用反射递归地复制值
func copyRecursive(original, cpy reflect.Value, visited map[uintptr]reflect.Value) {
	// 处理不同的类型
	switch original.Kind() {
	case reflect.Ptr:
		if original.IsNil() {
			cpy.Set(reflect.Zero(original.Type()))
			return
		}

		// 检查是否已经复制过这个指针
		ptr := original.Pointer()
		if v, ok := visited[ptr]; ok {
			cpy.Set(v)
			return
		}

		// 首先检查指针本身是否有 DeepCopy 方法
		if method, found := hasDeepCopyMethod(original); found {
			result := callDeepCopy(original, method)
			if result.IsValid() {
				// 如果DeepCopy返回的是值类型，需要创建新指针
				if result.Type() != original.Type() {
					newPtr := reflect.New(result.Type())
					newPtr.Elem().Set(result)
					cpy.Set(newPtr)
				} else {
					cpy.Set(result)
				}
				visited[ptr] = cpy
				return
			}
		}

		originalValue := original.Elem()

		// 然后检查指针指向的值是否有 DeepCopy 方法
		if method, found := hasDeepCopyMethod(originalValue); found {
			result := callDeepCopy(originalValue, method)
			if result.IsValid() {
				newPtr := reflect.New(result.Type())
				newPtr.Elem().Set(result)
				cpy.Set(newPtr)
				visited[ptr] = cpy
				return
			}
		}

		cpy.Set(reflect.New(originalValue.Type()))
		// 保存新创建的指针
		visited[ptr] = cpy
		copyRecursive(originalValue, cpy.Elem(), visited)

	case reflect.Interface:
		if original.IsNil() {
			cpy.Set(reflect.Zero(original.Type()))
			return
		}
		originalValue := original.Elem()
		copyValue := reflect.New(originalValue.Type()).Elem()
		copyRecursive(originalValue, copyValue, visited)
		cpy.Set(copyValue)

	case reflect.Struct:
		// 特殊处理 time.Time
		if t, ok := original.Interface().(time.Time); ok {
			cpy.Set(reflect.ValueOf(t))
			return
		}

		// 检查结构体是否有 DeepCopy 方法
		if method, found := hasDeepCopyMethod(original); found {
			result := callDeepCopy(original, method)
			if result.IsValid() {
				cpy.Set(result)
				return
			}
		}

		// 复制结构体的每个导出字段
		for i := 0; i < original.NumField(); i++ {
			field := original.Type().Field(i)
			// 跳过未导出字段 (PkgPath 不为空表示未导出)
			if field.PkgPath != "" {
				continue
			}
			copyRecursive(original.Field(i), cpy.Field(i), visited)
		}

	case reflect.Slice:
		if original.IsNil() {
			cpy.Set(reflect.Zero(original.Type()))
			return
		}
		cpy.Set(reflect.MakeSlice(original.Type(), original.Len(), original.Cap()))
		for i := 0; i < original.Len(); i++ {
			copyRecursive(original.Index(i), cpy.Index(i), visited)
		}

	case reflect.Map:
		if original.IsNil() {
			cpy.Set(reflect.Zero(original.Type()))
			return
		}
		cpy.Set(reflect.MakeMap(original.Type()))
		for _, key := range original.MapKeys() {
			originalValue := original.MapIndex(key)
			copyValue := reflect.New(originalValue.Type()).Elem()
			copyRecursive(originalValue, copyValue, visited)
			// 对 map 的键也进行深拷贝
			copyKey := reflect.New(key.Type()).Elem()
			copyRecursive(key, copyKey, visited)
			cpy.SetMapIndex(copyKey, copyValue)
		}

	case reflect.Array:
		// 数组需要逐个元素进行深拷贝
		for i := 0; i < original.Len(); i++ {
			copyRecursive(original.Index(i), cpy.Index(i), visited)
		}

	case reflect.Chan, reflect.Func, reflect.UnsafePointer:
		// 这些类型直接复制（浅拷贝）
		// Chan: 通道是引用类型，通常需要共享
		// Func: 函数是不可变的，可以安全共享
		// UnsafePointer: 直接复制指针值
		cpy.Set(original)

	default:
		// 对于基本类型（int, string, bool, float等），直接设置值
		cpy.Set(original)
	}
}

// copyRecursiveWithCache 使用缓存的类型分析结果进行深拷贝，避免重复反射分析
func copyRecursiveWithCache(original, cpy reflect.Value, visited map[uintptr]reflect.Value, typeInfo *TypeAnalysisResult) {
	// 处理不同的类型
	switch original.Kind() {
	case reflect.Ptr:
		if original.IsNil() {
			cpy.Set(reflect.Zero(original.Type()))
			return
		}

		// 检查是否已经复制过这个指针
		ptr := original.Pointer()
		if v, ok := visited[ptr]; ok {
			cpy.Set(v)
			return
		}

		// 首先检查指针本身是否有 DeepCopy 方法
		if method, found := hasDeepCopyMethod(original); found {
			result := callDeepCopy(original, method)
			if result.IsValid() {
				if result.Type() != original.Type() {
					newPtr := reflect.New(result.Type())
					newPtr.Elem().Set(result)
					cpy.Set(newPtr)
				} else {
					cpy.Set(result)
				}
				visited[ptr] = cpy
				return
			}
		}

		originalValue := original.Elem()

		// 然后检查指针指向的值是否有 DeepCopy 方法
		if method, found := hasDeepCopyMethod(originalValue); found {
			result := callDeepCopy(originalValue, method)
			if result.IsValid() {
				newPtr := reflect.New(result.Type())
				newPtr.Elem().Set(result)
				cpy.Set(newPtr)
				visited[ptr] = cpy
				return
			}
		}

		cpy.Set(reflect.New(originalValue.Type()))
		visited[ptr] = cpy
		copyRecursiveWithCache(originalValue, cpy.Elem(), visited, nil) // 子类型分析信息暂时为nil

	case reflect.Interface:
		if original.IsNil() {
			cpy.Set(reflect.Zero(original.Type()))
			return
		}
		originalValue := original.Elem()
		copyValue := reflect.New(originalValue.Type()).Elem()
		copyRecursiveWithCache(originalValue, copyValue, visited, nil)
		cpy.Set(copyValue)

	case reflect.Struct:
		// 特殊处理 time.Time
		if t, ok := original.Interface().(time.Time); ok {
			cpy.Set(reflect.ValueOf(t))
			return
		}

		// 检查结构体是否有 DeepCopy 方法
		if method, found := hasDeepCopyMethod(original); found {
			result := callDeepCopy(original, method)
			if result.IsValid() {
				cpy.Set(result)
				return
			}
		}

		// 复制结构体的每个导出字段
		// 这里可以利用缓存的字段分析信息来优化
		for i := 0; i < original.NumField(); i++ {
			field := original.Type().Field(i)
			// 跳过未导出字段 (PkgPath 不为空表示未导出)
			if field.PkgPath != "" {
				continue
			}

			// 如果有字段分析信息，可以进一步优化
			var fieldTypeInfo *TypeAnalysisResult
			if typeInfo != nil && typeInfo.FieldAnalysis != nil {
				fieldTypeInfo = typeInfo.FieldAnalysis[field.Name]
			}

			copyRecursiveWithCache(original.Field(i), cpy.Field(i), visited, fieldTypeInfo)
		}

	case reflect.Slice:
		if original.IsNil() {
			cpy.Set(reflect.Zero(original.Type()))
			return
		}
		cpy.Set(reflect.MakeSlice(original.Type(), original.Len(), original.Cap()))
		for i := 0; i < original.Len(); i++ {
			copyRecursiveWithCache(original.Index(i), cpy.Index(i), visited, nil)
		}

	case reflect.Map:
		if original.IsNil() {
			cpy.Set(reflect.Zero(original.Type()))
			return
		}
		cpy.Set(reflect.MakeMap(original.Type()))
		for _, key := range original.MapKeys() {
			originalValue := original.MapIndex(key)
			copyValue := reflect.New(originalValue.Type()).Elem()
			copyRecursiveWithCache(originalValue, copyValue, visited, nil)
			// 对 map 的键也进行深拷贝
			copyKey := reflect.New(key.Type()).Elem()
			copyRecursiveWithCache(key, copyKey, visited, nil)
			cpy.SetMapIndex(copyKey, copyValue)
		}

	case reflect.Array:
		// 数组需要逐个元素进行深拷贝
		for i := 0; i < original.Len(); i++ {
			copyRecursiveWithCache(original.Index(i), cpy.Index(i), visited, nil)
		}

	case reflect.Chan, reflect.Func, reflect.UnsafePointer:
		// 这些类型直接复制（浅拷贝）
		cpy.Set(original)

	default:
		// 对于基本类型（int, string, bool, float等），直接设置值
		cpy.Set(original)
	}
}
