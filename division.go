package division

import (
	_ `embed`
	`encoding/json`
	`strings`
)

const (
	// codeUnknown 未知的区划码
	codeUnknown code = 0
	// codeProvince 省级区划码
	codeProvince code = 1
	// codeCity 市级区划码
	codeCity code = 2
	// codeArea 县级区划码
	codeArea code = 3
)

type (
	// code 区划码种类
	code int8

	// Data 缓存的区划码数据
	Data struct {
		// 区划码
		Code string `json:"code"`
		// 区划名字
		Name string `json:"name"`
	}

	// Division 缓存的区划码对象
	Division struct {
		// 省级区划码与区划名字的对应关系
		province map[string]string
		// 市级区划码与区划名字的对应关系
		city map[string]string
		// 县级区划码与区划名字的对应关系
		area map[string]string

		// 省级区划码数据
		provinces []Data
		// 市级区划码数据
		cities map[string][]Data
		// 县级区划码数据
		areas map[string][]Data

		// 省级区划码的二进制数据
		provinceBytes []byte
		// 市级区划码的二进制数据
		cityBytes map[string][]byte
		// 县级区划码的二进制数据
		areaBytes map[string][]byte
	}
)

//go:embed asset/division.json
var divisionJson []byte

func newDivision() (d *Division, err error) {
	var data []Data
	if err = json.Unmarshal(divisionJson, &data); nil != err {
		return
	}

	d = &Division{
		province: make(map[string]string, 40),
		city:     make(map[string]string, 30),
		area:     make(map[string]string, 30),

		cities: make(map[string][]Data, 40),
		areas:  make(map[string][]Data, 30),

		cityBytes: make(map[string][]byte, 40),
		areaBytes: make(map[string][]byte, 30),
	}

	for _, v := range data {
		switch codeType(v.Code) {
		case codeProvince:
			d.province[v.Code[:2]] = v.Name
			d.provinces = append(d.provinces, v)
		case codeCity:
			d.city[v.Code[:4]] = v.Name
			d.cities[v.Code[:2]+"0000"] = append(d.cities[v.Code[:2]+"0000"], v)
		case codeArea:
			d.area[v.Code] = v.Name
			d.areas[v.Code[:4]+"00"] = append(d.areas[v.Code[:4]+"00"], v)
		}
	}

	if d.provinceBytes, err = json.Marshal(d.provinces); nil != err {
		return
	}

	for k, v := range d.cities {
		if d.cityBytes[k], err = json.Marshal(v); nil != err {
			return
		}
	}

	for k, v := range d.areas {
		if d.areaBytes[k], err = json.Marshal(v); nil != err {
			return
		}
	}

	return
}

func (d *Division) GetChildren(code string) []byte {
	if "" == code || "000000" == code {
		return d.provinceBytes
	}

	switch codeType(code) {
	case codeProvince:
		return d.cityBytes[code]
	case codeCity:
		return d.areaBytes[code]
	default:
		return nil
	}
}

func (d *Division) GetName(code string, seps ...string) string {
	var sep string
	if 0 != len(seps) {
		sep = seps[0]
	}

	return strings.Join(d.getName(code), sep)
}

func (d *Division) getName(code string) []string {
	var rsp []string

	switch codeType(code) {
	case codeProvince:
		rsp = []string{d.province[code[:2]]}
	case codeCity:
		rsp = []string{d.province[code[:2]], d.city[code[:4]]}
	case codeArea:
		rsp = []string{d.province[code[:2]], d.city[code[:4]], d.area[code]}
	default:
		return nil
	}

	if "" == rsp[len(rsp)-1] {
		rsp = nil
	}

	return rsp
}

func codeType(code string) code {
	if 6 != len(code) {
		return codeUnknown
	}

	switch {
	case "00" == code[:2]:
		return codeUnknown
	case "0000" == code[2:]:
		return codeProvince
	case "00" == code[4:]:
		return codeCity
	default:
		return codeArea
	}
}
