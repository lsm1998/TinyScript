package interpreter

import (
	"fmt"
	"os"
	"reflect"
	"strconv"

	"tiny-script/ast"
	"tiny-script/errors"
	"tiny-script/resolver"
	"tiny-script/token"
	"tiny-script/valuer"
)

var (
	True  = &valuer.Boolean{Value: true}
	False = &valuer.Boolean{Value: false}
	Nil   = &valuer.Nil{}
)

// potential value is empty or "repl".
var evalEnv string

var (
	env     *valuer.Environment
	globals *valuer.Environment
)

func init() {
	initEnv()
}

func initEnv() {
	globals = valuer.NewEnv()
	env = globals
}

func Interpret(statements []ast.Stmt) {
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(errors.RuntimeError); ok {
				fmt.Fprintln(os.Stderr, err.Error())
			} else {
				panic(r)
			}
		}
	}()
	for _, stmt := range statements {
		resolver.Resolve(stmt)
	}
	//var v valuer.Valuer
	for _, stmt := range statements {
		val := Eval(stmt)
		if val != nil {
			if val.Type() == valuer.ReturnType {
				fmt.Fprintf(os.Stderr, "Unexpected return statement %v\n", val)
			} else {
				//v = val
			}
		}
	}

	//if v != nil && evalEnv == "repl" {
	//	fmt.Printf("%s %s\n", black(v.Type().String()), v)
	//}
}

func Eval(node ast.Node) valuer.Valuer {
	switch n := node.(type) {
	default:
		panic(fmt.Sprintf("unknown ast type %#v.", n))
	case *ast.Literal:
		return evalLiteral(n)
	case *ast.BinaryExpr:
		return evalBinaryExpr(n)
	case *ast.UnaryExpr:
		return evalUnaryExpr(n)
	case *ast.GroupingExpr:
		return Eval(n.Expression)
	case *ast.VariableExpr:
		return evalVariableExpr(n)
	case *ast.AssignExpr:
		return evalAssignExpr(n)
	case *ast.ArrayAssignExpr:
		return evalArrayAssignExpr(n)
	case *ast.LogicalExpr:
		return evalLogicalExpr(n)
	case *ast.CallExpr:
		return evalCallExpr(n)
	case *ast.GetExpr:
		return evalGetExpr(n)
	case *ast.SetExpr:
		return evalSetExpr(n)
	case *ast.ThisExpr:
		return evalThisExpr(n)
	case *ast.VarStmt:
		evalVarStmt(n)
		return nil
	case *ast.LetStmt:
		evalLetStmt(n)
		return nil
	case *ast.FunctionStmt:
		evalFunctionStmt(n)
		return nil
	case *ast.PrintStmt:
		evalPrintStmt(n)
		return nil
	case *ast.BlockStmt:
		return evalBlockStmt(n)
	case *ast.ExprStmt:
		return evalExprStmt(n)
	case *ast.IfStmt:
		return evalIfStmt(n)
	case *ast.WhileStmt:
		return evalWhileStmt(n)
	case *ast.ReturnStmt:
		return evalReturnStmt(n)
	case *ast.ClassStmt:
		evalClassStmt(n)
		return nil
	case *ast.ImportStmt:
		evalImportStmt(n)
		return nil
	case *ast.ArrayLiteralExpr:
		return evalArrayLiteralExpr(n)
	case *ast.IndexLiteralExpr:
		return evalIndexLiteralExpr(n)
	case *ast.IndexVariableExpr:
		return evalIndexVariableExpr(n)
	}
}

func evalArrayAssignExpr(n *ast.ArrayAssignExpr) valuer.Valuer {
	expr := evalVariableExpr(n.Left)

	var index valuer.Valuer
	switch n.Index.(type) {
	case *ast.IndexLiteralExpr:
		index = Eval(n.Index.(*ast.IndexLiteralExpr).Index)
	case *ast.IndexVariableExpr:
		t := n.Index.(*ast.IndexVariableExpr)
		index = evalVariableExpr(t.Index.(*ast.VariableExpr))
	}

	if index.Type() != valuer.NumberType {
		panic("Index must be number.")
	}

	array, ok := expr.(*valuer.Array)
	if !ok {
		panic("Only array can be indexed.")
	}

	indexVal := int(index.(*valuer.Number).Value)

	if indexVal >= len(array.Elements) || indexVal < 0 {
		panic("Index out of range.")
	}

	array.Elements[indexVal] = Eval(n.Value)
	return nil
}

func evalIndexVariableExpr(n *ast.IndexVariableExpr) valuer.Valuer {
	v, ok := env.Get(n.Name)
	if !ok {
		errors.Error(token.Identifier, fmt.Sprintf("Undefined variable %s.", n.Name))
	}

	if v.Type() != valuer.NumberType {
		panic("Index must be number.")
	}

	left := Eval(n.Left)

	array, ok := left.(*valuer.Array)
	if !ok {
		panic("Only array can be indexed.")
	}

	index := int(v.(*valuer.Number).Value)
	if index >= len(array.Elements) || index < 0 {
		panic("Index out of range.")
	}
	return array.Elements[index]
}

func evalIndexLiteralExpr(n *ast.IndexLiteralExpr) valuer.Valuer {
	variableExpr, ok := n.Left.(*ast.VariableExpr)
	if !ok {
		panic("Only variable can be indexed.")
	}
	var v valuer.Valuer
	if variableExpr.Distance >= 0 {
		v, ok = env.GetAt(variableExpr.Distance, variableExpr.Name)
	} else {
		v, ok = globals.Get(variableExpr.Name)
	}
	if !ok {
		panic("Undefined variable.")
	}
	array, ok := v.(*valuer.Array)
	if !ok {
		panic("Only array can be indexed.")
	}
	index := Eval(n.Index)
	if index.Type() != valuer.NumberType {
		panic("Index must be number.")
	}
	if int(index.(*valuer.Number).Value) >= len(array.Elements) {
		panic("Index out of range.")
	}
	return array.Elements[int(index.(*valuer.Number).Value)]
}

func evalArrayLiteralExpr(expr *ast.ArrayLiteralExpr) valuer.Valuer {
	var elements = make([]valuer.Valuer, 0, len(expr.Elements))
	for _, e := range expr.Elements {
		elements = append(elements, Eval(e))
	}
	return &valuer.Array{Elements: elements}
}

func evalLiteral(lit *ast.Literal) valuer.Valuer {
	switch lit.Token {
	case token.True:
		return True
	case token.False:
		return False
	case token.String:
		return &valuer.String{Value: lit.Value}
	case token.Number:
		v, err := strconv.ParseFloat(lit.Value, 64)
		if err != nil {
			panic(err)
		}
		return &valuer.Number{Value: v}
	case token.Nil:
		return Nil
	default:
		panic("unhandled default case")
	}
}

func evalBinaryExpr(expr *ast.BinaryExpr) valuer.Valuer {
	left := Eval(expr.Left)
	right := Eval(expr.Right)

	switch op := expr.Operator; op {
	case token.EqualEqual:
		t := isEqual(left, right)
		return toBooleanValuer(t)
	case token.BangEqual:
		t := !isEqual(left, right)
		return toBooleanValuer(t)
	case token.Greater:
		a, b := checkNumberOperands(op, left, right)
		t := a > b
		return toBooleanValuer(t)
	case token.GreaterEqual:
		a, b := checkNumberOperands(op, left, right)
		t := a >= b
		return toBooleanValuer(t)
	case token.Less:
		a, b := checkNumberOperands(op, left, right)
		t := a < b
		return toBooleanValuer(t)
	case token.LessEqual:
		a, b := checkNumberOperands(op, left, right)
		t := a <= b
		return toBooleanValuer(t)
	case token.Minus:
		a, b := checkNumberOperands(op, left, right)
		v := a - b
		return &valuer.Number{Value: v}
	case token.Plus:
		return doPlusOperation(left, right)
	case token.Slash:
		a, b := checkNumberOperands(op, left, right)
		if b == float64(0) {
			errors.Error(op, "Divisor can't be 0.")
		}
		v := a / b
		return &valuer.Number{Value: v}
	case token.Star:
		a, b := checkNumberOperands(op, left, right)
		v := a * b
		return &valuer.Number{Value: v}
	default:
		panic("unhandled default case")
	}
}

func evalUnaryExpr(expr *ast.UnaryExpr) valuer.Valuer {
	right := Eval(expr.Right)
	switch op := expr.Operator; op {
	case token.Bang:
		t := !isTruthy(right)
		return toBooleanValuer(t)
	case token.Minus:
		v := checkNumberOperand(op, right)
		return &valuer.Number{Value: -v}
	default:
		panic("unhandled default case")
	}
}

func evalVariableExpr(expr *ast.VariableExpr) valuer.Valuer {
	if expr.Distance >= 0 {
		if v, ok := env.GetAt(expr.Distance, expr.Name); ok {
			return v
		}
	} else {
		if v, ok := globals.Get(expr.Name); ok {
			return v
		}
	}

	errors.Error(token.Identifier, fmt.Sprintf("Undefined variable %s.", expr.Name))
	return nil
}

func evalAssignExpr(expr *ast.AssignExpr) valuer.Valuer {
	v := Eval(expr.Value)
	name, distance := expr.Left.Name, expr.Left.Distance
	if distance >= 0 {
		if ok := env.AssignAt(distance, name, v); ok {
			return v
		}
	} else {
		if ok := globals.Assign(name, v); ok {
			return v
		}
	}
	errors.Error(token.Equal, fmt.Sprintf("Undefined variable %s.", expr.Left))
	return nil
}

func evalLogicalExpr(expr *ast.LogicalExpr) valuer.Valuer {
	left := Eval(expr.Left)
	switch expr.Operator {
	default:
		panic(fmt.Sprintf("unknown operator %s", expr.Operator))
	case token.Or:
		if isTruthy(left) {
			return left
		}
	case token.And:
		if !isTruthy(left) {
			return left
		}
	}
	return Eval(expr.Right)
}

func evalCallExpr(expr *ast.CallExpr) valuer.Valuer {
	callee := Eval(expr.Callee)
	callableValue, ok := callee.(valuer.Callable)
	if !ok {
		errors.Error(token.LeftParen, "Can only call functions and classes.")
		return nil
	}
	if l, l1 := callableValue.Arity(), len(expr.Arguments); l != l1 {
		errors.Error(token.LeftParen, fmt.Sprintf("Expected %d arguments but got %d", l, l1))
		return nil
	}

	switch n := callee.(type) {
	default:
		panic("invaid type")
	case *valuer.Function:
		return callFunction(n, expr.Arguments)
	case *valuer.ClassValue:
		return constructInstance(n, expr.Arguments)
	}
}

func constructInstance(c *valuer.ClassValue, arguments []ast.Expr) *valuer.Instance {
	instance := &valuer.Instance{Klass: c}
	initializer := c.FindMethod("init")
	if initializer != nil {
		callFunction(initializer.Bind(instance), arguments)
	}
	return instance
}

func callNativeFunc(function *valuer.Function, arguments []ast.Expr) valuer.Valuer {
	var values = make([]reflect.Value, 0, len(arguments))
	for _, v := range arguments {
		switch v.(type) {
		case *ast.Literal:
			if v.(*ast.Literal).Token == token.String {
				values = append(values, reflect.ValueOf(v.(*ast.Literal).Value))
				continue
			} else if v.(*ast.Literal).Token == token.Number {
				// 暂不支持浮点类型
				val, err := strconv.ParseInt(v.(*ast.Literal).Value, 10, 32)
				if err != nil {
					errors.Error(token.Function, "Invalid number.")
				}
				values = append(values, reflect.ValueOf(val))
			}
		}
	}
	result := function.NativeFunc.Call(values)
	if len(result) == 0 {
		return Nil
	} else if function.IsErr && result[len(result)-1].Interface() != nil { // native函数最后的返回值为error
		errors.Error(token.Function, result[len(result)-1].Interface().(error).Error())
	} else {
		switch result[0].Kind() {
		case reflect.Bool:
			return &valuer.Boolean{Value: result[0].Bool()}
		case reflect.Int:
			return &valuer.Number{Value: float64(result[0].Int())}
		case reflect.Int64:
			return &valuer.Number{Value: float64(result[0].Int())}
		case reflect.Float64:
			return &valuer.Number{Value: result[0].Float()}
		case reflect.String:
			return &valuer.String{Value: result[0].String()}
		case reflect.Float32:
			return &valuer.Number{Value: result[0].Float()}
		default:
			iface := result[0].Interface()
			if iface == nil {
				return Nil
			}
			return &valuer.String{Value: iface.(string)}
		}
	}
	return Nil
}

func callFunction(function *valuer.Function, arguments []ast.Expr) valuer.Valuer {
	if function.NativeFunc.IsValid() { // 是否是内置函数
		return callNativeFunc(function, arguments)
	}
	environment := function.Closure
	environment = valuer.NewEnclosing(function.Closure)
	for i, param := range function.Params {
		environment.Define(param.Name, Eval(arguments[i]))
	}
	v := executeBlock(function.Body, environment)
	if function.IsInitializer {
		// lookup this in function.Closure
		if v, ok := function.Closure.GetAt(0, "this"); ok {
			return v
		}
		errors.Error(token.Function, "Cann't get this in currrent enviroment.")
		return nil
	}
	if returnValue, ok := v.(*valuer.ReturnValue); ok {
		return returnValue.Value
	}
	return v
}

func evalGetExpr(expr *ast.GetExpr) valuer.Valuer {
	object := Eval(expr.Object)

	switch object.(type) {
	case *valuer.Instance:
		instance, _ := object.(*valuer.Instance)
		if v, ok := instance.Get(expr.Name); ok {
			return v
		}
		errors.Error(token.Identifier, fmt.Sprintf("Undefined propterty %s.", expr.Name))
	case *valuer.Array: // 为数组添加length属性
		array, _ := object.(*valuer.Array)
		switch expr.Name {
		case "length":
			return &valuer.Number{Value: float64(len(array.Elements))}
		default:
			errors.Error(token.Identifier, fmt.Sprintf("Undefined propterty %s.", expr.Name))
		}
	default:
		errors.Error(token.Identifier, "Only instances or array have properties.")
	}
	return nil
}

func evalSetExpr(expr *ast.SetExpr) valuer.Valuer {
	object := Eval(expr.Object)
	instance, ok := object.(*valuer.Instance)
	if !ok {
		errors.Error(token.Identifier, "Only instances have properties.")
		return nil
	}
	v := Eval(expr.Value)
	instance.Set(expr.Name, v)
	return v
}

func evalThisExpr(expr *ast.ThisExpr) valuer.Valuer {
	if v, ok := env.Get("this"); ok {
		return v
	}
	errors.Error(token.This, "Cannot use 'this' outside of a class.")
	return nil
}

func evalExprStmt(stmt *ast.ExprStmt) valuer.Valuer {
	return Eval(stmt.Expression)
}

func evalVarStmt(stmt *ast.VarStmt) {
	name := stmt.Name.Name
	var v valuer.Valuer
	if stmt.Initializer != nil {
		v = Eval(stmt.Initializer)
	} else {
		v = Nil
	}
	env.Define(name, v)
}

func evalLetStmt(stmt *ast.LetStmt) {
	name := stmt.Name.Name
	var v valuer.Valuer
	if stmt.Initializer != nil {
		v = Eval(stmt.Initializer)
	} else {
		v = Nil
	}
	env.Define(name, v)
}

func evalPrintStmt(stmt *ast.PrintStmt) {
	v := Eval(stmt.Expression)
	fmt.Println(v)
}

func evalBlockStmt(block *ast.BlockStmt) valuer.Valuer {
	return executeBlock(block.Statements, valuer.NewEnclosing(env))
}

func executeBlock(statements []ast.Stmt, environment *valuer.Environment) valuer.Valuer {
	previous := env
	env = environment
	defer func() {
		env = previous
	}()
	for _, stmt := range statements {
		result := Eval(stmt)
		if result != nil {
			if rt := result.Type(); rt == valuer.ReturnType {
				return result
			}
		}
	}
	return Nil
}

func evalIfStmt(stmt *ast.IfStmt) valuer.Valuer {
	condition := Eval(stmt.Condition)
	if isTruthy(condition) {
		return Eval(stmt.ThenBranch)
	} else if stmt.ElseBranch != nil {
		return Eval(stmt.ElseBranch)
	}
	return Nil
}

func evalWhileStmt(stmt *ast.WhileStmt) valuer.Valuer {
	for isTruthy(Eval(stmt.Condition)) {
		result := Eval(stmt.Body)
		if result != nil {
			if rt := result.Type(); rt == valuer.ReturnType {
				return result
			}
		}
	}
	return Nil
}

func evalFunctionStmt(stmt *ast.FunctionStmt) {
	fn := &valuer.Function{
		Name:    stmt.Name,
		Params:  stmt.Params,
		Body:    stmt.Body,
		Closure: env,
	}
	env.Define(stmt.Name, fn)
}

func evalReturnStmt(stmt *ast.ReturnStmt) valuer.Valuer {
	var v valuer.Valuer = Nil
	if stmt.Value != nil {
		v = Eval(stmt.Value)
	}
	return &valuer.ReturnValue{
		Value: v,
	}
}

func evalImportStmt(stmt *ast.ImportStmt) {
	instance := valuer.GetNativeInstance(stmt.Name)
	if instance == nil {
		errors.Error(token.Identifier, "Cannot find native object.")
		return
	}
	globals.Define(stmt.Name, instance)
}

func evalClassStmt(stmt *ast.ClassStmt) {
	methods := make(map[string]*valuer.Function, len(stmt.Methods))
	for _, method := range stmt.Methods {
		fn := &valuer.Function{
			Name:          method.Name,
			Params:        method.Params,
			Body:          method.Body,
			Closure:       env,
			IsInitializer: method.IsInitializer,
		}
		methods[method.Name] = fn
	}
	cl := &valuer.ClassValue{
		Name:    stmt.Name,
		Mehtods: methods,
	}
	env.Define(stmt.Name, cl)
}

func checkNumberOperand(operator token.Token, right valuer.Valuer) float64 {
	a, ok := right.(*valuer.Number)
	if !ok {
		errors.Error(operator, "Operand must be a number.")
	}
	return a.Value
}

func checkNumberOperands(operator token.Token, left, right valuer.Valuer) (float64, float64) {
	a, ok := left.(*valuer.Number)
	b, ok1 := right.(*valuer.Number)
	if !(ok && ok1) {
		errors.Error(operator, "Operands must be numbers.")
	}
	return a.Value, b.Value
}

func doPlusOperation(left, right valuer.Valuer) valuer.Valuer {
	switch l := left.(type) {
	case *valuer.Number, *valuer.String:
		switch r := right.(type) {
		case *valuer.Number:
			if n, ok := l.(*valuer.Number); ok {
				return &valuer.Number{Value: n.Value + r.Value}
			}
			s, _ := l.(*valuer.String)
			return &valuer.String{
				Value: s.Value + r.String(),
			}
		case *valuer.String:
			if n, ok := l.(*valuer.Number); ok {
				return &valuer.String{
					Value: n.String() + r.Value,
				}
			}
			s, _ := l.(*valuer.String)
			return &valuer.String{Value: s.Value + r.Value}
		}
	}

	errors.Error(token.Plus, "Operands must be numbers or strings.")
	return nil
}

func isEqual(a, b valuer.Valuer) bool {
	_, ok := a.(*valuer.Boolean)
	_, ok1 := b.(*valuer.Boolean)
	if ok || ok1 {
		return isTruthy(a) == isTruthy(b)
	}

	switch a1 := a.(type) {
	case *valuer.Number:
		if b1, ok := b.(*valuer.Number); ok {
			return a1.Value == b1.Value
		}
	case *valuer.Nil:
		if _, ok := b.(*valuer.Nil); ok {
			return true
		}
	case *valuer.String:
		if b1, ok := b.(*valuer.String); ok {
			return a1.Value == b1.Value
		}
	}
	return false
}

func isTruthy(value valuer.Valuer) bool {
	if value == nil {
		return false
	}
	switch v := value.(type) {
	case *valuer.Boolean:
		return v.Value
	case *valuer.Number:
		return v.Value != float64(0)
	case *valuer.Nil:
		return false
	case *valuer.String:
		return v.Value != ""
	}
	return false
}

func toBooleanValuer(t bool) *valuer.Boolean {
	if t {
		return True
	}
	return False
}

func black(s string) string {
	return "\033[1;30m" + s + "\033[0m"
}

// SetEvalEnv specify eval env of Interpreter.
func SetEvalEnv(envConfig string) {
	evalEnv = envConfig
}
