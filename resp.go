package resp

const (
	bufferSize         = 32 * 1024
	simpleStringPrefix = '+'
	errorPrefix        = '-'
	integerPrefix      = ':'
	bulkStringPrefix   = '$'
)

var objectSuffix = []byte("\r\n")
