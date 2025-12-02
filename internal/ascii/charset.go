package ascii

import (
	"fmt"
)

func ValidatePrintable(s string) error {
	for i := 0; i < len(s); i++ {
		chr := s[i]
		if chr =='\n' || chr == '\r'{
			continue
		}
		if chr < 32 || chr > 126 {
			return ErrNonPrintable
		}
	}
	return nil
}

var ErrNonPrintable = fmt.Errorf("input contains non-printable ASCII characters")
