package hashcode

import (
	"bytes"
	"fmt"
	"hash/crc32"
)

// String hashes a string to a unique hashcode.
//
// Copy from https://github.com/hashicorp/terraform-plugin-sdk/blob/v1.17.2/helper/hashcode/hashcode.go
// following depracation comment
func String(s string) int {
	v := int(crc32.ChecksumIEEE([]byte(s)))
	if v >= 0 {
		return v
	}
	if -v >= 0 {
		return -v
	}
	// v == MinInt
	return 0
}

// Strings hashes a list of strings to a unique hashcode.
//
// Copy from https://github.com/hashicorp/terraform-plugin-sdk/blob/v1.17.2/helper/hashcode/hashcode.go
// following depracation comment
func Strings(strings []string) string {
	var buf bytes.Buffer

	for _, s := range strings {
		buf.WriteString(fmt.Sprintf("%s-", s))
	}

	return fmt.Sprintf("%d", String(buf.String()))
}
