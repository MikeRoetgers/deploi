package protobuf

type Response interface {
	GetHeader() *ResponseHeader
}
