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
		"{1,2,3,4}",
	}
	for i, expr := range normalExpressions {
		tokenizer := ast.NewTokenizer(expr)
		tokens, err := tokenizer.Process()
		if err != nil {
			fmt.Errorf("tokenization failed: %v", err)
			return
		}
		fmt.Println("语法解析结果:")
		for count, token := range tokens {
			fmt.Printf("[%d] %s\n", count, token)
		}
		fmt.Printf("\n=== 词法解析: %d ===", i+1)
		result, err := parser.ParseExpressionWithContext(expr, nil)
		if err != nil {
			fmt.Printf("❌ 解析错误: %v\n", err)
		} else {
			ast.PrintASTWithTitle(result.AST, "完整 AST 树形结构")
		}
	}
}

```


## 文件结构

```
├── ast
│   ├── ast_analyzer.go # AST 分析器实现
│   ├── ast_base.go # AST 基础接口和类型定义
│   ├── ast.go  # AST 核心接口和类型定义
│   ├── ast_nodes.go  # AST 节点实现
│   ├── ast_operators.go  # 操作符实现
│   ├── parser.go  # 语法分析器实现
│   ├── token.go   # 词法令牌实现
│   ├── tokenizer.go  # 词法分析器实现
│   └── token_kind.go # 词法令牌类型定义
├── go.mod
├── lexical_test.go   # 词法分析测试文件
├── LICENSE  
├── main.go       # 主程序入口
├── parser_test.go  # 语法分析测试文件
└── README.md
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

## 与 Java 版本的差异
1. **空值处理**：Go 使用指针 `*string` 来模拟 Java 的 `@Nullable String`
2. **字符处理**：Go 使用 `[]rune` 来正确处理 Unicode 字符
3. **错误处理**：Go 使用显式的错误返回而不是异常
4. **集合**：Go 使用切片 `[]*Token` 代替 Java 的 `List<Token>`
5. **枚举**：Go 使用常量和 iota 来模拟 Java 枚举

