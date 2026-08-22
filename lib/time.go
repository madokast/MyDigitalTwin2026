package lib

import "time"

var UTC8 = time.FixedZone("CST", 8*3600)

const RFC3339 = time.RFC3339
const RFC3339Milli = "2006-01-02T15:04:05.000Z07:00"
const RFC3339Nano = time.RFC3339Nano
