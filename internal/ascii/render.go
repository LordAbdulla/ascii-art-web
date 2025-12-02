package ascii

import (
	"strings"
)

func DoAsciiArt(input string, lines []string) string {
	norm := strings.ReplaceAll(input, "\r\n", "\n")
	norm = strings.ReplaceAll(norm, "\r", "\n")
	norm = strings.ReplaceAll(norm, `\n`, "\n")

	const height = 8
	const block = 9

	var out strings.Builder
	parts := strings.Split(norm, "\n")

	onlyLine := len(norm) == 0

	for i, word := range parts {
		if len(word) == 0 {
			if !(onlyLine && i == len(parts)-1) {
				out.WriteByte('\n')
			}
			continue
		}
		for row := 1; row <= height; row++ {
			for _, ch := range word {
				idx := int(ch-32) * block
				out.WriteString(lines[idx+row])
			}
			out.WriteByte('\n')
		}
	}
	return out.String()
}
