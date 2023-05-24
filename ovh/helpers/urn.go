package helpers

import "fmt"

const (
	VPSkind string = "vps"
)

func ServiceURN(plate, kind, name string) string {
	return fmt.Sprintf("urn:v1:%s:resource:%s:%s", plate, kind, name)
}
