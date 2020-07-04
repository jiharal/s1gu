package cmd

func (n S1GU) createGeneralID() string {
	return `package system

	import uuid "github.com/satori/go.uuid"
	
	func DefaultID() uuid.UUID {
		return uuid.FromStringOrNil("5f1c7fe3-d3a6-4896-974a-3385185893f9")
	}
	`
}

func (n S1GU) createValidate() string {
	return `
	package system

	import (
		uuid "github.com/satori/go.uuid"
		"github.com/shopspring/decimal"
	)
	
	// TextValidate is used to validate text
	func TextValidate(param, val string) string {
		var str string
		if str != param {
			return param
		}
		return val
	}
	
	// BoolValidate is used to validate bool.
	func BoolValidate(param, val bool) bool {
		var bl bool
		if bl != param {
			return param
		}
		return val
	}
	
	// Float64Validate is used to validate val of float64.
	func Float64Validate(param, val float64) float64 {
		var flt64 float64
		if param != flt64 {
			return param
		}
		return val
	}
	
	// UUIDValidate is used to validate uuid.
	func UUIDValidate(param, val uuid.UUID) uuid.UUID {
		var uid uuid.UUID
		if uid != param {
			return param
		}
		return val
	}
	
	// IntValidate is used to validate int.
	func IntValidate(param, val int) int {
		var v int
		if v != param {
			return param
		}
		return val
	}
	
	// DecimalValidate is used to validate decimal.
	func DecimalValidate(param, val decimal.Decimal) decimal.Decimal {
		var decim decimal.Decimal
		if decim != param {
			return param
		}
		return val
	}	
	`
}
