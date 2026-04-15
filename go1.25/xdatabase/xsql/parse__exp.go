package xsql

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Deimvis/go-ext/go1.25/ext"
	"github.com/Deimvis/go-ext/go1.25/xcheck/xmust"
)

// ParseQueriesNaive parses sql queries separated with ";"
// Removes trailing ";" in the end of queries.
func ParseQueriesNaive(s string) []string {
	queries := strings.Split(s, ";")
	ext.MapIn(&queries, func(q string) string { return strings.TrimSpace(q) })
	ext.FilterIn(&queries, func(q string) bool { return q != "" })
	return queries
}

// ParseQueries parses sql queries separated with ";".
// Removes trailing ";" in the end of queries.
// Parses statements with dollar quoted string (like DO statement).
func ParseQueries(s string) ([]string, error) {
	var queries []string
	var buffer []string
	insideDO := false
	delimiterDO := ""
	insideDollarQuotedString := false
	for i, line := range strings.Split(s, "\n") {
		buffer = append(buffer, line)

		matchDollarQuotedString := DollarQuotedStringRegexp.FindStringSubmatch(line)
		if matchDollarQuotedString != nil {
			insideDollarQuotedString = !insideDollarQuotedString
		}
		if insideDollarQuotedString {
			continue
		}

		matchDO := DOStatementRegexp.FindStringSubmatch(line)
		if matchDO != nil { // DO statement start
			if insideDO {
				return nil, fmt.Errorf("syntax error: nested DO (line: %d)", i)
			}
			insideDO = true
			delimiterDO = matchDO[DOStatementRegexp.SubexpIndex("delimiter")]
		} else if insideDO { // inside DO statement
			if strings.Contains(line, delimiterDO) { // DO statement end
				queries = append(queries, removeTrailingSemicolon(strings.Join(buffer, "\n")))
				buffer = nil

				insideDO = false
				delimiterDO = ""
			}
		} else { // just line
			if strings.Contains(line, ";") { // query end
				queries = append(queries, removeTrailingSemicolon(strings.Join(buffer, "\n")))
				buffer = nil
			}
		}
	}
	if len(buffer) > 0 {
		queries = append(queries, removeTrailingSemicolon(strings.Join(buffer, "\n")))
		buffer = nil
	}
	ext.MapIn(&queries, func(q string) string { return strings.TrimSpace(q) })
	ext.FilterIn(&queries, func(q string) bool { return q != "" })
	return queries, nil
}

func removeTrailingSemicolon(s string) string {
	return strings.Trim(strings.TrimSpace(s), ";")
}

var DOStatementRegexp *regexp.Regexp
var DollarQuotedStringRegexp *regexp.Regexp

func init() {
	DOStatementRegexp = xmust.Do(regexp.Compile(`^DO\s+(?P<delimiter>[^\s]+)$`))
	DollarQuotedStringRegexp = xmust.Do(regexp.Compile(`.*\$\$.*`))
}
