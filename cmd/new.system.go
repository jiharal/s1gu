package cmd

func CreateGeneralID() string {
	return `package system

	import uuid "github.com/satori/go.uuid"
	
	func DefaultID() uuid.UUID {
		return uuid.FromStringOrNil("5f1c7fe3-d3a6-4896-974a-3385185893f9")
	}
	`
}
