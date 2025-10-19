package main

import (
	"fmt"
	"testing"

	"github.com/weaweawe01/ParserSpel/ast"
)

func TestMainCHeck(t *testing.T) {
	parser := ast.NewSpelExpressionParser()
	// 测试普通表达式
	normalExpressions := []string{
		"\"Hello ' World\"",
		"testMap['monday']",
		"testMap.get('monday')",
		"new org.springframework.expression.spel.OperatorTests$SubComparable(2) > new org.springframework.expression.spel.OperatorTests$OtherSubComparable(1)",
		"new java.math.BigDecimal('5') eq new java.math.BigDecimal('5')",
		"new java.math.BigDecimal('5') != new java.math.BigDecimal('3')",
		"new java.math.BigDecimal('5') ne new java.math.BigDecimal('5')",
		"new org.springframework.expression.spel.OperatorTests$SubComparable() ne new org.springframework.expression.spel.OperatorTests$OtherSubComparable()",
		"3.0d < new java.math.BigDecimal('3.0')",
		"new java.math.BigDecimal('5') <= new java.math.BigDecimal('5')",
		"#service.findFruitsByColor(null)?.[1]",
		"#service.findFruitsByColor(null)?.^[#this.length > 5]",
		"#service.findJediByName(null) ?: 'unknown'",
		"#service.findJediByName('Yoda').name",
		"#service.findJediByName('Yoda').present",
		"#service.findJediByName('').orElse('Luke')",
		"0xB0BG",                 // 错误的
		"true or ",               //错误
		"1 + ",                   //错误
		"null instanceof T('a')", //错误
		"#var1.methodOne().methodTwo(42)",
		"#func1().methodOne().methodTwo(42)",
		"property1[0][1].property2['key'][42].methodTwo()",
		"property1?.[0]?.[1]?.property2?.['key']?.[42]?.methodTwo()",
		"#_$_='value'",
		"__age",
		"Person_1.get__age()",
		"null.?[#this < 5]",
		"mapOfNumbersUpToTen.?['hello']",
		"{1,2,3,4,5,6}?.?[#this between {2, 4}]",
		"integers.^[#this < 5]",
		"testMap.keySet().?[#this matches '.*o.*']",
		"publicName='Andy'",
		"arrayContainer.booleans[1]",
		"#aMap['one'] eq 1",
		"printDouble(T(java.math.BigDecimal).valueOf(14.35))",
		"printDouble(14.35)",
		"printDoubles(getDoublesAsStringList())",
		"new Spr5899Class(null,'a','b').toString()",
		"#varargsFunction2(9, new String[0])",
		"#varargsObjectFunction('a',null,'b')",
	}
	for i, expr := range normalExpressions {
		result, err := parser.ParseExpressionWithContext(expr, nil)
		if err != nil {
			fmt.Printf("\n=== 词法解析: %d  表达式:[%s]===", i+1, expr)
			fmt.Printf("❌ 解析错误: %v\n", err)
		} else {
			ast.PrintASTWithTitle(result.AST, "完整 AST 树形结构")
		}
	}
}
