# Go SpEL Parser

这是完整的 Spring Expression Language (SpEL) 解析器的 Go 语言版本，从 Java 原版转换而来。包含完整的词法分析和语法分析功能。
基于（spring-expression-6.2.11.jar）版本开发


## 安装
```go
go get github.com/weaweawe01/ParserSpel/ast

```
## 使用案例
```go
package main

import (
	"fmt"
	"github.com/weaweawe01/ParserSpel/ast"
)

func main() {
	parser := ast.NewSpelExpressionParser()
	// 测试普通表达式
	normalExpressions := []string{
		"T(java.lang.Runtime).getRuntime().exec('id')",
	}
	for i, expr := range normalExpressions {
		fmt.Printf("\n=== 表达式测试 %d ===\n", i+1)
		fmt.Printf("表达式: %s\n", expr)
		result, err := parser.ParseExpressionWithContext(expr, nil)
		if err != nil {
			fmt.Printf("❌ 解析错误: %v\n", err)
		} else {
			fmt.Printf("   ✅ 解析成功!\n")
			ast.PrintASTWithTitle(result.AST, "完整 AST 树形结构")
		}
	}
}

```


## 文件结构

```
go_spel/
├── go.mod              # Go 模块文件
├── main.go            # 主程序和测试用例  
├── token_kind.go      # TokenKind 枚举和相关方法
├── token.go           # Token 结构体和方法
├── tokenizer.go       # 词法分析器
├── ast_base.go        # AST 基础结构和接口
├── ast_nodes.go       # 基本 AST 节点类型
├── ast_operators.go   # 操作符 AST 节点
├── parser.go          # 语法分析器主要逻辑
└── README.md          # 说明文档
```

## 主要组件

### 词法分析层

#### 1. TokenKind (token_kind.go)
- 定义了所有支持的令牌类型
- 包括字面量、操作符、标识符等
- 提供字符串表示和长度计算方法

#### 2. Token (token.go)
- 表示单个令牌，包含类型、数据和位置信息
- 支持各种令牌检查方法（IsIdentifier、IsNumericRelationalOperator等）
- 提供令牌转换方法

#### 3. Tokenizer (tokenizer.go)
- 主要的词法分析器
- 将输入字符串转换为令牌序列
- 支持所有 SpEL 语法元素

### 语法分析层

#### 4. AST Base (ast_base.go)
- 定义 SpelNode 接口和基础实现
- 提供 TypedValue、ExpressionState 等核心结构
- 管理求值上下文和配置

#### 5. AST Nodes (ast_nodes.go)
- 实现基础 AST 节点类型
- 包括字面量、标识符、属性访问、变量引用等
- 支持各种数据类型的解析和求值

#### 6. AST Operators (ast_operators.go)
- 实现所有操作符节点
- 包括算术、比较、逻辑、一元操作符
- 支持运算符优先级和求值逻辑

#### 7. Parser (parser.go)
- 主要的语法分析器
- 实现递归下降解析算法
- 构建抽象语法树并支持表达式求值

## 支持的 SpEL 语法

### 基本类型
- 整数：`123`, `123L`
- 十六进制：`0x1A2B`, `0x1A2BL`
- 浮点数：`3.14`, `3.14F`, `1.23e-4`
- 字符串：`'single'`, `"double"`
- 布尔值：`true`, `false`

### 操作符
- 算术：`+`, `-`, `*`, `/`, `%`, `^`
- 比较：`>`, `>=`, `<`, `<=`, `==`, `!=`
- 逻辑：`&&`, `||`, `!`
- 赋值：`=`, `+=`, `-=`
- 自增/自减：`++`, `--`

### 特殊操作符
- 条件：`?:`（Elvis）, `? :`（三元）
- 安全导航：`?.`
- 类型判断：`instanceof`
- 正则匹配：`matches`
- 范围：`between`

### 集合操作
- 选择：`?[...]`（过滤）
- 投影：`![...]`（映射）
- 首个：`^[...]`
- 最后：`$[...]`

### 访问符号
- 属性访问：`.`
- 索引访问：`[...]`
- Bean 引用：`@bean`
- 工厂 Bean：`&factory`
- 根对象：`#root`

## 功能特性

### 支持的表达式类型

#### 基本字面量
- 整数：`42`, `123L`, `0x1A2B`
- 浮点数：`3.14`, `3.14F`, `1.23e-4`
- 字符串：`'hello'`, `"world"`
- 布尔值：`true`, `false`
- 空值：`null`

#### 算术运算
- 基本运算：`+`, `-`, `*`, `/`, `%`
- 运算符优先级：`2 + 3 * 4` = `14`
- 括号表达式：`(2 + 3) * 4` = `20`

#### 比较运算
- 数值比较：`>`, `>=`, `<`, `<=`
- 相等比较：`==`, `!=`
- 正则匹配：`name matches '[A-Z].*'`

#### 逻辑运算
- 逻辑与：`&&`
- 逻辑或：`||`
- 逻辑非：`!`
- 短路求值支持

#### 属性访问
- 普通访问：`user.name`
- 安全导航：`obj?.property`
- 复合表达式：`user.address.city`

#### 特殊引用
- 变量引用：`#root`, `#variable`
- Bean 引用：`@myBean`, `@'complex.bean.name'`

## 运行示例

```bash
cd go_spel
go run .
```

## 测试结果

程序包含了完整的测试用例，涵盖从基础到复杂的各种 SpEL 表达式：

1. `name == 'John'` - 字符串比较
2. `age > 18 && active == true` - 逻辑表达式
3. `price * quantity + tax` - 算术运算
4. `user.name.length()` - 属性访问和方法调用
5. `items[0].value` - 索引访问
6. `#root.method()` - 根对象引用
7. `T(Math).PI * radius^2` - 类型引用和幂运算
8. `name matches '[A-Z].*'` - 正则匹配
9. `value instanceof T(java.lang.String)` - 类型判断
10. `list.?[#this > 10]` - 集合过滤
11. `map.![value * 2]` - 集合投影
12. `items.^[price > 100]` - 首个匹配
13. `items.$[price < 50]` - 最后匹配
14. `condition ? 'yes' : 'no'` - 三元操作符
15. `value ?: 'default'` - Elvis 操作符
16. `obj?.property` - 安全导航
17. `++counter` - 自增
18. `total += amount` - 复合赋值
19. `0x1A2B` - 十六进制
20. `3.14159F` - 浮点数
21. `1.23e-4` - 科学记数法
22. `"hello world"` - 双引号字符串
23. `'single quoted'` - 单引号字符串

## 与 Java 版本的差异

1. **空值处理**：Go 使用指针 `*string` 来模拟 Java 的 `@Nullable String`
2. **字符处理**：Go 使用 `[]rune` 来正确处理 Unicode 字符
3. **错误处理**：Go 使用显式的错误返回而不是异常
4. **集合**：Go 使用切片 `[]*Token` 代替 Java 的 `List<Token>`
5. **枚举**：Go 使用常量和 iota 来模拟 Java 枚举

## 扩展

这个词法分析器可以作为完整 SpEL 解析器的基础，后续可以添加：
- 语法分析器（Parser）
- 抽象语法树（AST）构建
- 表达式求值器（Evaluator）
- 上下文处理器（Context Handler）