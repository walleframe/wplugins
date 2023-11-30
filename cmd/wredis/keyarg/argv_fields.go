package keyarg

// type KeyStructer interface {
// 	GetFieldName(idx int) string
// 	GetFieldType(field string) *buildpb.FieldType
// }

// // 结构体字段参数
// type FieldArg struct {
// 	EmptyArg
// 	field string
// }

// func (x *FieldArg) DynamicCode() (code string) {
// 	return "b.WriteString()"
// }

// // 结构体取模字段参数
// type FieldModArg struct {
// 	EmptyArg
// 	field string
// 	mod   int64
// }

// func (x *FieldModArg) Imports() []string {
// 	return []string{"strconv"}
// }

// func (x *FieldModArg) DynamicCode() (code string) {
// 	return "x." + x.field + "%" + strconv.FormatInt(x.mod, 10)
// }

// type FieldMatch struct{}

// func (FieldMatch) Match(arg string, st KeyStructer) (KeyArg, error) {
// 	if !strings.HasPrefix(arg, "$") {
// 		return nil, nil
// 	}
// 	field := strings.TrimPrefix(arg, "$")
// 	if strings.Contains(field, "%") {
// 		// 字段取模
// 		list := strings.SplitN(field, "%", 2)
// 		if len(list) != 2 {
// 			return nil, fmt.Errorf("[%s] invalid! need $field_name%%number", arg)
// 		}
// 		mod, err := strconv.ParseInt(list[1], 10, 64)

// 		if err != nil {
// 			return nil, fmt.Errorf("[%s] invalid! %+v", arg, err)
// 		}
// 		real, err := checkField(list[0], st)
// 		if err != nil {
// 			return nil, err
// 		}
// 		// 字段非数值类型
// 		if !st.Integer(real) {
// 			return nil, fmt.Errorf("[%s] field is not integer,can use mod operation.", real)
// 		}
// 		return &FieldModArg{field: real, mod: mod}, nil
// 	}
// 	real, err := checkField(field, st)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &FieldArg{field: real}, nil
// }

// func checkField(field string, st KeyStructer) (string, error) {
// 	if st.CheckField(field) {
// 		return field, nil
// 	}
// 	num, err := strconv.ParseInt(field, 10, 64)
// 	if err != nil {
// 		return "", fmt.Errorf("[%s] invalid field number.", field)
// 	}

// 	field = st.GetField(int(num))
// 	if field == "" {
// 		return "", fmt.Errorf("[%d] invalid field number.", num)
// 	}
// 	return field, nil
// }
