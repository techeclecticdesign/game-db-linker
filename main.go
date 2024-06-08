package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/brandenc40/romannumeral"
)

/* Converts an integer string to a roman numeral.  Only works for values less than 20.
 * Returns original string on failure. */
func int2roman(numstr string) (string, bool) {
	num, err := strconv.Atoi(numstr)
	if err != nil {
		return numstr, false
	}

	if num > 20 {
		return numstr, false
	}

	result, err := romannumeral.IntToString(num)
	if err != nil {
		return numstr, false
	}

	return result, true
}

func _stripEnclosed(s string, regex string) string {
	r := regexp.MustCompile(regex)
	return r.ReplaceAllString(s, "")
}

/* Strip all parentheses and brackets from a string. */
func stripEnclosed(s string) string {
	result := s
	// Strip all parentheses and brackets
	result = _stripEnclosed(s, `\(.*?\)`)
	result = _stripEnclosed(result, `\[.*?\]`)
	result = _stripEnclosed(result, `\{.*?\}`)
	result = _stripEnclosed(result, `\<.*?\>`)
	return result
}

func stripSeparators(s string) string {
	// Strip all separators
	result := s
	result = regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(result, "")
	return result
}

/* Strip article (a, an, the) prefixes from a string. */
func stripArticlePrefix(s string) string {
	// Strip article prefixes
	if strings.HasPrefix(s, "the ") {
		return strings.TrimPrefix(s, "the ")
	} else if strings.HasPrefix(s, "a ") {
		return strings.TrimPrefix(s, "a ")
	} else if strings.HasPrefix(s, "an ") {
		return strings.TrimPrefix(s, "an ")
	}
	if strings.HasSuffix(s, ", the") {
		return strings.TrimSuffix(s, ", the")
	} else if strings.HasSuffix(s, ", a") {
		return strings.TrimSuffix(s, ", a")
	} else if strings.HasSuffix(s, ", an") {
		return strings.TrimSuffix(s, ", an")
	}
	return s
}

/* Try to find numerical string corresponding to a releases sequel number,
 * and convert it to a roman numeral. */
func sequelConvert(s string, sepLocs [][]int) string {
	for i := len(sepLocs) - 1; i >= 0; i-- {
		last := strings.LastIndex(s[:sepLocs[i][0]], " ") + 1

		lastWord := s[last:sepLocs[i][0]]
		if strnum, ok := int2roman(lastWord); ok {
			return s[:last] + strnum + s[sepLocs[i][0]:]
		}
	}

	return s
}

func main() {
	var sepLocs [][]int // ending positions of the last word in the string
	/* Simplify string to reduce stylistic differences, to facilitate more accurate comparison. */
	parseString := func(input string) string {
		result := stripEnclosed(input)
		result = stripArticlePrefix(result)
		/* find all occurances of : or -.  Convert to roman numeral if it makes sense to.  */
		sepLocs = regexp.MustCompile(`[-/:]`).FindAllStringIndex(result, -1)
		sepLocs = append(sepLocs, []int{len(result)})
		result = sequelConvert(result, sepLocs)
		result = stripSeparators(result)
		result = strings.Join(strings.Fields(result), " ") // get rid of extra white space
		result = strings.ToLower(result)
		return result
	}
	fmt.Println(parseString("The Game 12 - The Test (USA) [translation by Team X] {v1.0}"))
}
