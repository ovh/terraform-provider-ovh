// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package tfaddr

import (
	"fmt"
)

type ParserError struct {
	Summary string
	Detail  string
}

func (pe *ParserError) Error() string {
	return fmt.Sprintf("%s: %s", pe.Summary, pe.Detail)
}
