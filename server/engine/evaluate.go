package engine

import (
	"bytes"
	"demo-backend/server/def"
	"demo-backend/server/engine/formatter"
	"demo-backend/server/engine/marshal"
	"demo-backend/server/engine/stack"
	"demo-backend/server/io"
	"fmt"
	"regexp"
	"strings"

	"github.com/RoaringBitmap/roaring"
)

var operators = map[string]bool{
	"OR":  true,
	"AND": true,
	"NOT": true,
}

//TODO: add corresponding function to each operator ("AND, "OR", "NOT")
var execute = map[string]func(roaring.Bitmap, roaring.Bitmap) roaring.Bitmap{
	"AND": func(rb1, rb2 roaring.Bitmap) roaring.Bitmap {
		return *roaring.FastAnd(&rb1, &rb2)
	},
	"OR": func(rb1, rb2 roaring.Bitmap) roaring.Bitmap {

		return *roaring.FastOr(&rb1, &rb2)
	},
	//TODO: implement NOT
	//"NOT IN": func(rb1, rb2 roaring.Bitmap) roaring.Bitmap {
	//	return rb1.AndNot()
	//},
}

var arithmeticExecution = map[string]func(io.Store, string, string, []byte, []byte,
	[]byte, []byte, []byte, []byte) (roaring.Bitmap, error){

	"=": func(s io.Store, fieldName string, fieldType string, byteOrderedValue []byte,
		dbID []byte, namespaceID []byte, collectionID []byte, compositeIndexKey []byte, compositePrefix []byte) (roaring.Bitmap, error) {

		rb := roaring.New()
		indexKey := []byte{}
		if len(compositeIndexKey) == 0 {
			indexKey = []byte(def.IndexKey + string(dbID) + ":" + string(collectionID) + ":" + string(namespaceID) + ":" + fieldName + ":" + fieldType + ":" + string(byteOrderedValue))
		} else {
			indexKey = compositeIndexKey
		}

		uniqueIDBitmapArray, err := s.Get(indexKey)
		if len(uniqueIDBitmapArray) == 0 || err != nil {
			return roaring.Bitmap{}, err
		}
		err = rb.UnmarshalBinary(uniqueIDBitmapArray)

		if err != nil {
			return roaring.Bitmap{}, err
		}

		return *rb, nil

	},

	//TODO: discuss memory related issue here
	">": func(s io.Store, fieldName string, fieldType string, byteOrderedValue []byte,
		dbID []byte, namespaceID []byte, collectionID []byte, compositeIndexKey []byte, compositePrefix []byte) (roaring.Bitmap, error) {
		rb := roaring.New()

		startKey := []byte{}
		prefix := []byte{}

		if len(compositeIndexKey) != 0 && len(compositePrefix) != 0 {
			startKey = compositeIndexKey
			prefix = compositePrefix
		} else {
			startKey = []byte(def.IndexKey + string(dbID) + ":" + string(collectionID) + ":" + string(namespaceID) + ":" + fieldName + ":" + fieldType + ":" + string(byteOrderedValue))
			prefix = []byte(def.IndexKey + string(dbID) + ":" + string(collectionID) + ":" + string(namespaceID) + ":" + fieldName + ":" + fieldType + ":")
		}

		keys, values, err := s.PrefixScan(startKey, prefix, 0)

		if len(values) == 0 || err != nil {
			return roaring.Bitmap{}, err
		}

		if bytes.Compare(keys[0], startKey) == 0 {
			values = values[1:]
		}

		for _, v := range values {
			if len(rb.ToArray()) == 0 {
				err = rb.UnmarshalBinary(v)
				if err != nil {
					return roaring.Bitmap{}, err
				}
			}
			tempRb := roaring.New()
			err = tempRb.UnmarshalBinary(v)
			if err != nil {
				return roaring.Bitmap{}, err
			}

			rb = roaring.FastOr(rb, tempRb)
		}

		return *rb, err

	},

	"<": func(s io.Store, fieldName string, fieldType string, byteOrderedValue []byte,
		dbID []byte, namespaceID []byte, collectionID []byte, compositeIndexKey []byte, compositePrefix []byte) (roaring.Bitmap, error) {

		rb := roaring.New()

		endKey := []byte{}
		prefix := []byte{}

		if len(compositeIndexKey) != 0 && len(compositePrefix) != 0 {
			endKey = compositeIndexKey
			prefix = compositePrefix
		} else {
			endKey = []byte(def.IndexKey + string(dbID) + ":" + string(collectionID) + ":" + string(namespaceID) + ":" + fieldName + ":" + fieldType + ":" + string(byteOrderedValue))
			prefix = []byte(def.IndexKey + string(dbID) + ":" + string(collectionID) + ":" + string(namespaceID) + ":" + fieldName + ":" + fieldType + ":")
		}

		keys, values, err := s.ReversePrefixScan(endKey, prefix, 0)

		if len(values) == 0 || err != nil {
			return roaring.Bitmap{}, err
		}

		if bytes.Compare(keys[0], endKey) == 0 {
			values = values[1:]
		}

		for _, v := range values {
			if len(rb.ToArray()) == 0 {
				err = rb.UnmarshalBinary(v)
				if err != nil {
					return roaring.Bitmap{}, err
				}
			}
			tempRb := roaring.New()
			err = tempRb.UnmarshalBinary(v)
			if err != nil {
				return roaring.Bitmap{}, err
			}
			rb = roaring.FastOr(rb, tempRb)
		}

		return *rb, err

	},

	">=": func(s io.Store, fieldName string, fieldType string, byteOrderedValue []byte,
		dbID []byte, namespaceID []byte, collectionID []byte, compositeIndexKey []byte, compositePrefix []byte) (roaring.Bitmap, error) {
		rb := roaring.New()

		startKey := []byte{}
		prefix := []byte{}
		if len(compositeIndexKey) != 0 && len(compositePrefix) != 0 {
			startKey = compositeIndexKey
			prefix = compositePrefix

		} else {
			startKey = []byte(def.IndexKey + string(dbID) + ":" + string(collectionID) + ":" + string(namespaceID) + ":" + fieldName + ":" + fieldType + ":" + string(byteOrderedValue))
			prefix = []byte(def.IndexKey + string(dbID) + ":" + string(collectionID) + ":" + string(namespaceID) + ":" + fieldName + ":" + fieldType + ":")
		}

		_, uniqueIDBitmapArray, err := s.PrefixScan(startKey, prefix, 0)

		if len(uniqueIDBitmapArray) == 0 || err != nil {
			return roaring.Bitmap{}, err
		}

		for _, v := range uniqueIDBitmapArray {
			if len(rb.ToArray()) == 0 {
				err = rb.UnmarshalBinary(v)
				if err != nil {
					return roaring.Bitmap{}, err
				}
			}
			tempRb := roaring.New()
			err = tempRb.UnmarshalBinary(v)
			if err != nil {
				return roaring.Bitmap{}, err
			}
			rb = roaring.FastOr(rb, tempRb)
		}

		return *rb, err

	},

	"<=": func(s io.Store, fieldName string, fieldType string, byteOrderedValue []byte,
		dbID []byte, namespaceID []byte, collectionID []byte, compositeIndexKey []byte, compositePrefix []byte) (roaring.Bitmap, error) {

		//fmt.Println("[[evaluate.go/airthmeticExecution<]]")
		rb := roaring.New()

		endKey := []byte{}
		prefix := []byte{}

		if len(compositeIndexKey) != 0 && len(compositePrefix) != 0 {
			endKey = compositeIndexKey
			prefix = compositePrefix
		} else {
			endKey = []byte(def.IndexKey + string(dbID) + ":" + string(collectionID) + ":" + string(namespaceID) + ":" + fieldName + ":" + fieldType + ":" + string(byteOrderedValue))
			prefix = []byte(def.IndexKey + string(dbID) + ":" + string(collectionID) + ":" + string(namespaceID) + ":" + fieldName + ":" + fieldType + ":")
		}

		_, uniqueIDBitmapArray, err := s.ReversePrefixScan(endKey, prefix, 0)

		if len(uniqueIDBitmapArray) == 0 || err != nil {
			return roaring.Bitmap{}, err
		}

		for _, v := range uniqueIDBitmapArray {
			if len(rb.ToArray()) == 0 {
				err = rb.UnmarshalBinary(v)
				if err != nil {
					return roaring.Bitmap{}, err
				}
			}
			tempRb := roaring.New()
			err = tempRb.UnmarshalBinary(v)
			if err != nil {
				return roaring.Bitmap{}, err
			}
			rb = roaring.FastOr(rb, tempRb)
		}

		return *rb, err

	},

	"!=": func(s io.Store, fieldName string, fieldType string, byteOrderedValue []byte,
		dbID []byte, namespaceID []byte, collectionID []byte, compositeIndexKey []byte, compositePrefix []byte) (roaring.Bitmap, error) {

		rb := roaring.New()

		endKey := []byte{}
		prefix := []byte{}

		if len(compositeIndexKey) != 0 && len(compositePrefix) != 0 {
			endKey = compositeIndexKey
			prefix = compositePrefix
		} else {
			endKey = []byte(def.IndexKey + string(dbID) + ":" + string(collectionID) + ":" + string(namespaceID) + ":" + fieldName + ":" + fieldType + ":" + string(byteOrderedValue))
			prefix = []byte(def.IndexKey + string(dbID) + ":" + string(collectionID) + ":" + string(namespaceID) + ":" + fieldName + ":" + fieldType + ":")
		}

		keysleft, uniqueIDBitmapArray, err := s.ReversePrefixScan(endKey, prefix, 0)

		if len(uniqueIDBitmapArray) == 0 || err != nil {
			return roaring.Bitmap{}, err
		}

		if bytes.Compare(keysleft[0], endKey) == 0 {
			uniqueIDBitmapArray = uniqueIDBitmapArray[1:]
		}

		keysRight, valuesFor, err := s.PrefixScan(endKey, prefix, 0)

		if len(valuesFor) == 0 || err != nil {
			return roaring.Bitmap{}, err
		}

		if bytes.Compare(keysRight[0], endKey) == 0 {
			valuesFor = valuesFor[1:]
		}

		uniqueIDBitmapArray = append(uniqueIDBitmapArray, valuesFor...)

		for _, v := range uniqueIDBitmapArray {
			if len(rb.ToArray()) == 0 {
				err = rb.UnmarshalBinary(v)
				if err != nil {
					return roaring.Bitmap{}, err
				}
			}
			tempRb := roaring.New()
			err = tempRb.UnmarshalBinary(v)
			if err != nil {
				return roaring.Bitmap{}, err
			}
			rb = roaring.FastOr(rb, tempRb)
		}

		return *rb, err

	},
}

//EvaluatePostFix evaluates postfix expression returns result
func (e *Engine) EvaluatePostFix(px []string, collectionID []byte) (interface{}, error) {
	if len(px) == 1 {
		rb, err := e.EvaluateExpression(px[0], collectionID)
		if err != nil {
			var tmp interface{}
			return tmp, err
		}
		var result interface{}
		result = rb
		return result, nil
	}
	var tempStack stack.Stack
	for _, v := range px {
		if _, ok := operators[v]; !ok {
			tempStack.Push(v)
		} else {
			exp1 := tempStack.Pop()
			exp2 := tempStack.Pop()
			exp1Type := fmt.Sprintf("%T", exp1)
			exp2Type := fmt.Sprintf("%T", exp2)

			var rb1 roaring.Bitmap
			var rb2 roaring.Bitmap
			var err error

			if exp1Type == "string" {
				rb1, err = e.EvaluateExpression(exp1.(string), collectionID)
				if err != nil {
					var tmp interface{}
					return tmp, err
				}
			} else {
				rb1 = exp1.(roaring.Bitmap)
			}

			if exp2Type == "string" {
				rb2, err = e.EvaluateExpression(exp2.(string), collectionID)
				if err != nil {
					var tmp interface{}
					return tmp, err
				}

			} else {
				rb2 = exp2.(roaring.Bitmap)
			}

			tempStack.Push(execute[v](rb1, rb2))
		}
	}

	//fmt.Println("[[evaluate.go]] EVALUATE POSTFIX")
	result := tempStack.Pop()

	return result, nil

}

//EvaluateExpression takes in expression and returns roaring bitmap as result
func (e *Engine) EvaluateExpression(exp string, collectionID []byte) (roaring.Bitmap, error) {
	/*
		1. Parse expression to find fieldname, operator, fieldvalue, fieldtype
		2. Based on operator, carry out operations
	*/
	compositeCondArr := strings.Split(exp, ",")
	if len(compositeCondArr) == 1 {
		//parse fieldname,operator,fieldvalue
		fieldname, operator, fieldvalue := parseExpressionFields(exp)
		//get fieldtype with ordered value
		typeOfData, byteOrderedData := findTypeOfValue(fieldvalue)

		rb, err := arithmeticExecution[operator](e.Store, fieldname, typeOfData, byteOrderedData, e.DBID, e.NamespaceID, collectionID, []byte{}, []byte{})

		if err != nil {
			return roaring.Bitmap{}, err
		}

		return rb, nil
	} else {
		//for composite condition
		indexKey := def.IndexKey + string(e.DBID) + ":" + string(collectionID) + ":" + string(e.NamespaceID)
		prefix := indexKey
		lastOperator := ""
		lenCompositeArr := len(compositeCondArr)
		i := 0
		//create index key first
		for _, cond := range compositeCondArr {
			fieldname, operator, fieldvalue := parseExpressionFields(cond)
			typeOfData, byteOrderedData := findTypeOfValue(fieldvalue)
			if i == (lenCompositeArr - 1) {
				indexKey += ":" + fieldname + ":" + typeOfData + ":" + string(byteOrderedData)
				prefix += ":" + fieldname + ":" + typeOfData + ":"
			} else {
				indexKey += ":" + fieldname + ":" + typeOfData + ":" + string(byteOrderedData)
				prefix += ":" + fieldname + ":" + typeOfData + ":" + string(byteOrderedData)
			}
			i++
			lastOperator = operator
		}

		rb, err := arithmeticExecution[lastOperator](e.Store, "", "", []byte{}, e.DBID, e.NamespaceID, collectionID, []byte(indexKey), []byte(prefix))
		if err != nil {
			return roaring.Bitmap{}, err
		}

		return rb, nil
	}

}

func parseExpressionFields(exp string) (string, string, string) {
	re := regexp.MustCompile(`(!=|>=|>|<=|<|=)`)
	operator := re.Find([]byte(exp)) //get first operator that is matched
	strArr := strings.Split(exp, string(operator))
	return strArr[0], string(operator), (strArr[1])
}

func findTypeOfValue(value string) (string, []byte) {
	datatype, formattedData, err := formatter.FormatData(value)
	if err != nil {
		panic(err)
	}

	specificDataType := def.ApplicationSpecificType[datatype]

	return specificDataType, marshal.TypeMarshal(datatype, formattedData)
}
