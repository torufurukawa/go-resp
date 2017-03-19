package resp

const (
	bufferSize         = 32 * 1024
	simpleStringPrefix = '+'
	errorPrefix        = '-'
	integerPrefix      = ':'
	bulkStringPrefix   = '$'
	arrayPrefix        = '*'
)

var objectSuffix = []byte("\r\n")
