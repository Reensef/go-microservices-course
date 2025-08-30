package utils

type MultiValueType int32

const (
	MultiValueType_UNSPECIFIED MultiValueType = 0
	MultiValueType_INT64       MultiValueType = 1
	MultiValueType_FLOAT64     MultiValueType = 2
	MultiValueType_STRING      MultiValueType = 3
	MultiValueType_BOOL        MultiValueType = 4
)

type MultiValue struct {
	valueType MultiValueType `bson:"valueType"`
	value     any            `bson:"value"`
}

// Type возвращает тип хранимого значения
func (p *MultiValue) Type() MultiValueType {
	return p.valueType
}

// Int64 возвращает значение как int64, если тип совпадает
func (p *MultiValue) Int64() (int64, bool) {
	if val, ok := p.value.(int64); ok {
		return val, true
	}
	return 0, false
}

// MustInt64 возвращает значение как int64 без проверки
func (p *MultiValue) MustInt64() int64 {
	return p.value.(int64)
}

// Float64 возвращает значение как float64, если тип совпадает
func (p *MultiValue) Float64() (float64, bool) {
	if val, ok := p.value.(float64); ok {
		return val, true
	}
	return 0, false
}

// MustFloat64 возвращает значение как float64 без проверки
func (p *MultiValue) MustFloat64() float64 {
	return p.value.(float64)
}

// String возвращает значение как string, если тип совпадает
func (p *MultiValue) String() (string, bool) {
	if val, ok := p.value.(string); ok {
		return val, true
	}
	return "", false
}

// MustString возвращает значение как string без проверки
func (p *MultiValue) MustString() string {
	return p.value.(string)
}

// Bool возвращает значение как bool, если тип совпадает
func (p *MultiValue) Bool() (bool, bool) {
	if val, ok := p.value.(bool); ok {
		return val, true
	}
	return false, false
}

// MustBool возвращает значение как bool без проверки
func (p *MultiValue) MustBool() bool {
	return p.value.(bool)
}

// SetInt64 устанавливает значение типа int64
func (p *MultiValue) SetInt64(value int64) {
	p.valueType = MultiValueType_INT64
	p.value = value
}

// SetFloat64 устанавливает значение типа float64
func (p *MultiValue) SetFloat64(value float64) {
	p.valueType = MultiValueType_FLOAT64
	p.value = value
}

// SetString устанавливает значение типа string
func (p *MultiValue) SetString(value string) {
	p.valueType = MultiValueType_STRING
	p.value = value
}

// SetBool устанавливает значение типа bool
func (p *MultiValue) SetBool(value bool) {
	p.valueType = MultiValueType_BOOL
	p.value = value
}
