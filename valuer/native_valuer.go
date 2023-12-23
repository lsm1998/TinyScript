package valuer

import (
	"tiny-script/ast"
	"tiny-script/errors"
	"tiny-script/native"
	"tiny-script/token"
)

func GetNativeInstance(objName string) Valuer {
	// 是否支持此导入对象
	nativeObj := native.GetNativeObject(objName)
	if nativeObj == nil {
		errors.Error(token.Identifier, "Cannot find native object.")
		return nil
	}

	var methods = make(map[string]*Function)

	// 获取对象的所有方法
	for k, v := range nativeObj {
		var params []*ast.Ident
		for _, param := range v.Params {
			params = append(params, &ast.Ident{
				Name: param,
			})
		}
		methods[k] = &Function{
			Name:          k,
			Params:        params,
			Body:          nil,
			Closure:       nil,
			IsInitializer: false,
			NativeFunc:    v.Func,
			IsErr:         v.IsErr,
		}
	}
	return &Instance{
		Klass: &ClassValue{
			Name:    objName,
			Mehtods: methods,
		},
		Fileds: nil,
	}
}
