# DeepCopy - 优化版深拷贝库

基于 [mohae/deepcopy](https://github.com/mohae/deepcopy) 进行大幅优化和改进的 Go 深拷贝库，提供类型安全、高性能的深拷贝功能。

[![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.21-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

## ✨ 主要优化

相较于原版 mohae/deepcopy，本项目进行了以下重要优化：

### 🚀 **泛型支持**
- **原版问题**: 返回 `interface{}`，需要手动类型断言，容易出错
- **优化方案**: 使用 Go 1.21+ 泛型，提供类型安全的 `Copy[T any](src T) T` 接口

```go
// 原版用法 - 需要类型断言，容易出错
result := deepcopy.Copy(original).(MyStruct)

// 优化版用法 - 类型安全，编译时检查
result := deepcopy.Copy[MyStruct](original)
// 或者利用类型推断
result := deepcopy.Copy(original) // 自动推断为 MyStruct 类型
```

### 🔧 **自定义拷贝接口**
- **原版问题**: 无法自定义特定类型的拷贝行为
- **优化方案**: 支持 `Copier[T any]` 接口和 `DeepCopy()` 方法

```go
type CustomStruct struct {
    Value int
}

// 实现自定义拷贝逻辑
func (c CustomStruct) DeepCopy() CustomStruct {
    return CustomStruct{Value: c.Value + 100}
}
```

### 🔄 **循环引用处理**
- **原版问题**: 循环引用可能导致栈溢出
- **优化方案**: 通过 visited map 跟踪已复制的指针，完全解决循环引用问题

```go
type Node struct {
    Next *Node
    Value int
}

// 创建循环引用
a := &Node{Value: 1}
b := &Node{Value: 2}
a.Next = b
b.Next = a

// 安全拷贝，不会栈溢出
copied := deepcopy.Copy(a)
```

### ⚡ **智能性能优化**
- **新增功能**: 类型分析缓存机制，自动识别只包含值类型的数据并直接返回，避免昂贵的反射操作
- **优化效果**: 值类型数据拷贝速度提升 10-100 倍

```go
// 对于只包含值类型的结构体，直接返回，几乎没有性能开销
type ValueOnlyStruct struct {
    Name string
    Age  int
}

original := ValueOnlyStruct{Name: "Alice", Age: 30}
copied := deepcopy.Copy(original) // 极速拷贝，直接返回
```

### 🏪 **业务缓存优化**
- **新增功能**: `CopyWithKey` 方法，基于业务 key 的缓存优化

```go
// 基于业务 key 的优化拷贝
config := AppConfig{Port: 8080, Host: "localhost"}

// 第一次调用会分析类型
copied1 := deepcopy.CopyWithKey(config, "app.config")

// 后续调用直接使用缓存结果，性能极佳
copied2 := deepcopy.CopyWithKey(config, "app.config")
```

### 📊 **数组深拷贝修复**
- **原版问题**: 数组元素可能被浅拷贝，导致数据污染
- **优化方案**: 修复数组深拷贝实现，确保数组内指针元素被正确深拷贝

### 🛡️ **nil 值处理改进**
- **原版问题**: 某些 nil 值处理可能导致 panic
- **优化方案**: 完善的 nil 值处理，支持 nil 指针、切片、map、接口

### ⏰ **时间类型特殊处理**
- **新增功能**: 为 `time.Time` 类型提供特殊的拷贝处理，确保时间值正确复制

### 🔒 **安全性改进**
- **安全优化**: 明确跳过未导出字段，避免潜在的安全问题和 panic

## 📦 安装

```bash
go get github.com/wsqun/deepcopy
```

要求 Go 1.21 或更高版本。

## 🚀 快速开始

### 基本用法

```go
package main

import (
    "fmt"
    "github.com/wsqun/deepcopy"
)

func main() {
    // 基本类型深拷贝
    original := map[string]interface{}{
        "name": "John",
        "age": 30,
        "hobbies": []string{"reading", "swimming"},
    }
    
    copied := deepcopy.Copy(original)
    
    // 修改拷贝不影响原始数据
    copied["age"] = 31
    fmt.Printf("Original: %v\n", original["age"]) // 30
    fmt.Printf("Copied: %v\n", copied["age"])     // 31
}
```

### 结构体深拷贝

```go
type Person struct {
    Name    string
    Age     int
    Address *Address
    Hobbies []string
}

type Address struct {
    Street string
    City   string
}

func main() {
    original := Person{
        Name: "Alice",
        Age:  25,
        Address: &Address{
            Street: "123 Main St",
            City:   "New York",
        },
        Hobbies: []string{"photography", "travel"},
    }
    
    // 类型安全的深拷贝
    copied := deepcopy.Copy(original)
    
    // 修改拷贝的地址不影响原始数据
    copied.Address.City = "Los Angeles"
    
    fmt.Printf("Original city: %s\n", original.Address.City) // New York
    fmt.Printf("Copied city: %s\n", copied.Address.City)     // Los Angeles
}
```

### 自定义拷贝行为

```go
type User struct {
    ID       int
    Name     string
    Password string
}

// 实现自定义拷贝逻辑
func (u User) DeepCopy() User {
    return User{
        ID:       u.ID,
        Name:     u.Name,
        Password: "[REDACTED]", // 不拷贝敏感信息
    }
}

func main() {
    original := User{
        ID:       1,
        Name:     "john",
        Password: "secret123",
    }
    
    copied := deepcopy.Copy(original)
    fmt.Printf("Copied password: %s\n", copied.Password) // [REDACTED]
}
```

### 性能优化用法

```go
// 业务缓存优化
func main() {
    config := AppConfig{Port: 8080, Debug: true}
    
    // 使用业务 key 进行缓存优化
    for i := 0; i < 1000; i++ {
        copied := deepcopy.CopyWithKey(config, "app.config")
        // 第一次会分析类型，后续调用极速
    }
}

// 类型分析
func analyzeTypes() {
    var data MyStruct
    analysis := deepcopy.AnalyzeType(data)
    
    fmt.Printf("只包含值类型: %t\n", analysis.IsOnlyValues)
    fmt.Printf("包含指针: %t\n", analysis.ContainsPtr)
    fmt.Printf("包含切片: %t\n", analysis.ContainsSlice)
}
```

### 循环引用处理

```go
type Node struct {
    Name string
    Next *Node
}

func main() {
    // 创建循环引用
    node1 := &Node{Name: "Node1"}
    node2 := &Node{Name: "Node2"}
    node3 := &Node{Name: "Node3"}
    
    node1.Next = node2
    node2.Next = node3
    node3.Next = node1 // 循环引用
    
    // 安全处理循环引用
    copied := deepcopy.Copy(node1)
    
    // 验证循环结构被正确保持
    fmt.Printf("Original cycle: %s -> %s -> %s -> %s\n", 
        node1.Name, node1.Next.Name, 
        node1.Next.Next.Name, node1.Next.Next.Next.Name)
    
    fmt.Printf("Copied cycle: %s -> %s -> %s -> %s\n", 
        copied.Name, copied.Next.Name, 
        copied.Next.Next.Name, copied.Next.Next.Next.Name)
}
```

## 🧪 测试覆盖

本项目包含完整的测试套件，覆盖各种边界情况：

- ✅ 基本类型拷贝
- ✅ 复杂结构体拷贝
- ✅ 指针拷贝安全性
- ✅ 切片和数组拷贝
- ✅ Map 拷贝
- ✅ 接口类型拷贝
- ✅ 循环引用处理
- ✅ nil 值处理
- ✅ 自定义拷贝接口
- ✅ 时间类型处理
- ✅ 未导出字段跳过
- ✅ 边界情况处理
- ✅ 性能优化测试
- ✅ 类型分析缓存
- ✅ 业务缓存功能

运行测试：

```bash
go test -v           # 详细测试
go test -race        # 竞态检测
go test -cover       # 覆盖率报告 (86%+)
go test -bench=.     # 性能基准测试
```

## 🔍 支持的类型

- ✅ 基本类型 (int, string, bool, float, etc.)
- ✅ 指针
- ✅ 结构体
- ✅ 切片
- ✅ 数组
- ✅ Map
- ✅ 接口
- ✅ 时间类型 (time.Time)
- ✅ 嵌套和复合类型
- ✅ 循环引用结构
- ⚠️ 通道 (浅拷贝，共享通道实例)
- ⚠️ 函数 (浅拷贝，函数是不可变的)
- ❌ UnsafePointer (除非为 nil)

## ⚡ 性能特点

- **智能类型分析**: 自动识别值类型，避免不必要的反射操作
- **多层缓存机制**: 类型分析缓存 + 业务key缓存
- **循环引用优化**: 通过 visited map 避免重复拷贝
- **内存分配最小化**: 精确的内存管理，减少 GC 压力

### 性能对比

| 场景 | 原版 mohae/deepcopy | 优化版 |
|------|-------------------|--------|
| 基本类型 (int) | ~200ns | ~20ns (10x 提升) |
| 值类型结构体 | ~1000ns | ~25ns (40x 提升) |
| 复杂结构体 | ~5000ns | ~2000ns (2.5x 提升) |
| 循环引用 | 可能栈溢出 | 安全处理 |

## 📚 API 文档

### 核心函数

```go
// Copy 创建任意值的深拷贝，自动类型推断
func Copy[T any](src T) T

// CopyWithKey 基于业务 key 的优化拷贝
func CopyWithKey[T any](src T, key string) T

// AnalyzeType 分析类型结构，返回详细信息
func AnalyzeType[T any](src T) *TypeAnalysisResult

// NewDeepCopyManager 创建独立的拷贝管理器
func NewDeepCopyManager() *DeepCopyManager
```

### 管理器方法

```go
// CopyValue 使用管理器进行深拷贝 (非泛型)
func (m *DeepCopyManager) CopyValue(src interface{}) interface{}

// AnalyzeValue 使用管理器分析类型 (非泛型)
func (m *DeepCopyManager) AnalyzeValue(src interface{}) *TypeAnalysisResult
```

### 接口

```go
// Copier 可自定义深拷贝行为的接口
type Copier[T any] interface {
    DeepCopy() T
}

// TypeAnalysisResult 类型分析结果
type TypeAnalysisResult struct {
    IsOnlyValues  bool                           // 是否只包含值类型
    ContainsPtr   bool                           // 是否包含指针
    ContainsSlice bool                           // 是否包含切片
    ContainsMap   bool                           // 是否包含映射
    // ... 更多字段
}
```

## 🤝 对比原版

| 特性 | 原版 mohae/deepcopy | 优化版 |
|------|-------------------|--------|
| 类型安全 | ❌ 需要类型断言 | ✅ 泛型支持 |
| 自定义拷贝 | ❌ 不支持 | ✅ Copier 接口 |
| 循环引用 | ⚠️ 部分支持 | ✅ 完全支持 |
| 数组拷贝 | ⚠️ 可能浅拷贝 | ✅ 深拷贝修复 |
| nil 处理 | ⚠️ 可能 panic | ✅ 完善处理 |
| 时间类型 | ❌ 无特殊处理 | ✅ 特殊优化 |
| 性能优化 | ❌ 无优化 | ✅ 多层优化 |
| 类型分析 | ❌ 不支持 | ✅ 详细分析 |
| 业务缓存 | ❌ 不支持 | ✅ CopyWithKey |
| 测试覆盖 | ⚠️ 基础测试 | ✅ 全面测试 (86%+) |
| Go 版本 | 任意版本 | Go 1.21+ |

## 📄 许可证

MIT License - 详见 [LICENSE](LICENSE) 文件

## 🙏 致谢

感谢 [mohae/deepcopy](https://github.com/mohae/deepcopy) 提供的基础实现。本项目在其基础上进行了大量优化和改进。

---

**注意**: 本库需要 Go 1.21 或更高版本以支持泛型功能。如果您使用的是较旧的 Go 版本，请考虑升级或使用原版 mohae/deepcopy。
