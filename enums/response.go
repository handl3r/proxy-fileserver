package enums

type Response interface {
	GetCode() int
	GetMessage() string
}
