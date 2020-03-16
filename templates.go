package main

import (
	"text/template"
)

func tplString(i interface{}) string {
	if str, ok := i.([]byte); ok {
		return string(str)
	} else {
		return "<UNCOMPATIBLE DATA>"
	}
}

func tplRange2(block []*Block) [][]Block {
	var ret [][]Block
	tmp := make([]Block, 2)
	tmp[0] = *block[0]
	ret = append(ret, tmp)
	for i := 1; i < len(block); i += 2 {
		tmp := make([]Block, 2)
		tmp[0] = *block[i]
		if i+1 < len(block) && block != nil {
			tmp[1] = *block[i+1]
		}
		ret = append(ret, tmp)
	}
	return ret
}

func tplAdd(i, v int) int {
	ret := i + v
	return ret
}

func tplGetFuncMap() template.FuncMap {
	return template.FuncMap{
		"String": tplString,
		"Add":    tplAdd,
		// "Hex": tplHex,
		"Range2": tplRange2,
	}
}
