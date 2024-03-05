package appenders

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/alessandrofascini/log4go/pkg"
)

const SpecialChar = ';'
const CompressExt = ".gz"

type Comparator struct {
	pattern string
	sep     string
}

func NewComparator(fullPath string, sep string, keepFileExt, isCompress bool) *Comparator {
	_, fullFileName := filepath.Split(fullPath)
	ext := filepath.Ext(fullFileName)
	filename := fullFileName[0 : len(fullFileName)-len(ext)]
	var p []string
	if keepFileExt {
		p = append(p, filename, string(SpecialChar), ext)
	} else {
		p = append(p, fullFileName, string(SpecialChar))
	}
	if isCompress {
		p = append(p, CompressExt)
	}
	return &Comparator{
		pattern: strings.Join(p, ""),
		sep:     sep,
	}
}

/*
Match
This method return -1 if no matching;
0 if it is the hot file
greater than 0 if there is a matching
*/
func (c *Comparator) Match(s string) int {
	i := 0
	for ; i < len(c.pattern) && c.pattern[i] != SpecialChar; i++ {
		if c.pattern[i] != s[i] {
			return -1
		}
	}
	j := len(s) - 1
	w := len(c.pattern) - 1
	for ; w > -1 && c.pattern[w] != SpecialChar; w-- {
		if c.pattern[w] != s[j] {
			return -1
		}
		j--
	}
	j++         // j is not included
	if i == j { // We found our hot file
		return 0
	}
	// Now indexing
	if j-i < len(c.sep) {
		return -1
	}
	w = 0
	for i < j && w < len(c.sep) {
		if s[i] != c.sep[w] {
			return -1
		}
		i++
		w++
	}
	if i == j { // This is useful for date file appender
		return 0
	}
	// between 'i' and 'j' we will find our number
	res := 0
	for i < j {
		res = res<<1 + res<<3
		b := int(s[i] - '0')
		res += b
		i++
	}
	return res
}

func (c *Comparator) Replace(i int) string {
	// Is better to create a custom function to do this?
	newStr := fmt.Sprintf("%s%d", c.sep, i)
	return strings.Replace(c.pattern, string(SpecialChar), newStr, -1)
}

// DateComparator struct
type DateComparator struct {
	sep         string
	pattern     string
	datePattern string
}

func NewDateComparator(fullPath string, sep string, keepFileExt, isCompress bool, datePattern string) DateComparator {
	_, fullFileName := filepath.Split(fullPath)
	ext := filepath.Ext(fullFileName)
	filename := fullFileName[0 : len(fullFileName)-len(ext)]
	var p []string
	if keepFileExt {
		p = append(p, filename, string(SpecialChar), ext)
	} else {
		p = append(p, fullFileName, string(SpecialChar))
	}
	if isCompress {
		p = append(p, CompressExt)
	}
	return DateComparator{
		sep:         sep,
		pattern:     strings.Join(p, ""),
		datePattern: datePattern,
	}
}

const days = 24 * 60 * 60 * 1000

// Match return the difference between now and the date saved
// -1
func (d *DateComparator) Match(s string) int {
	i := 0
	for ; i < len(d.pattern) && d.pattern[i] != SpecialChar; i++ {
		if d.pattern[i] != s[i] {
			return -1
		}
	}
	j := len(s) - 1
	w := len(d.pattern) - 1
	for ; w > -1 && d.pattern[w] != SpecialChar; w-- {
		if d.pattern[w] != s[j] {
			return -1
		}
		j--
	}
	j++         // j is not included
	if i == j { // We found our hot file
		return -1
	}
	// Now indexing
	if j-i < len(d.sep) {
		return -1
	}
	w = 0
	for i < j && w < len(d.sep) {
		if s[i] != d.sep[w] {
			return -1
		}
		i++
		w++
	}
	if i == j { // This is useful for date file appender
		return -1
	}
	// between 'i' and 'j' we will find our date
	if i+len(d.datePattern) > j {
		return -1
	}
	v := s[i : i+len(d.datePattern)]
	parsed, err := time.Parse(d.datePattern, v)
	if err != nil {
		return -1
	}
	now := time.Now()
	// in days
	diff := (now.UnixMilli() - parsed.UnixMilli()) / days
	if diff < 0 {
		return -1
	}
	return int(diff)
}

func (d *DateComparator) Replace(i int) string {
	date := fmt.Sprintf("%s%s", d.sep, pkg.GetDateByFormat(d.datePattern))
	if i != 0 {
		date = fmt.Sprintf("%s%s%d", date, d.sep, i)
	}
	return strings.Replace(d.pattern, string(SpecialChar), date, -1)
}

// Date File Comparator

type DateComparatorII struct {
	sep         string
	pattern     string
	datePattern string
}

func NewDateComparatorII(fullPath string, sep string, keepFileExt, isCompress bool, datePattern string) DateComparatorII {
	_, fullFileName := filepath.Split(fullPath)
	ext := filepath.Ext(fullFileName)
	filename := fullFileName[0 : len(fullFileName)-len(ext)]
	p := strings.Builder{}
	if keepFileExt {
		p.WriteString(filename)
		p.WriteString(string(SpecialChar))
		p.WriteString(ext)
	} else {
		p.WriteString(fullFileName)
		p.WriteString(string(SpecialChar))
	}
	if isCompress {
		p.WriteString(CompressExt)
	}
	return DateComparatorII{
		sep:         sep,
		pattern:     p.String(),
		datePattern: datePattern,
	}
}

// Match return the difference between now and the date saved
// -1
func (d *DateComparatorII) Match(s string) (int64, int) {
	i := 0
	for ; i < len(d.pattern) && d.pattern[i] != SpecialChar; i++ {
		if d.pattern[i] != s[i] {
			return -1, -1
		}
	}
	j := len(s) - 1
	w := len(d.pattern) - 1
	for ; w > -1 && d.pattern[w] != SpecialChar; w-- {
		if d.pattern[w] != s[j] {
			return -1, -1
		}
		j--
	}
	j++         // j is not included
	if i == j { // We found our hot file
		return -1, -1
	}
	// Now indexing
	if j-i < len(d.sep) {
		return -1, -1
	}
	w = 0
	for i < j && w < len(d.sep) {
		if s[i] != d.sep[w] {
			return -1, -1
		}
		i++
		w++
	}
	if i == j { // This is useful for date file appender
		return -1, -1
	}
	// between 'i' and 'j' we will find our date
	if i+len(d.datePattern) > j {
		return -1, -1
	}
	dateAndNumber := strings.Split(s[i:j], d.sep)
	if len(dateAndNumber) == 0 {
		return -1, -1
	}

	parsed, err := time.Parse(d.datePattern, dateAndNumber[0])
	if err != nil {
		return -1, -1
	}

	now := time.Now()
	midnightToday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	// in days
	diff := (midnightToday.UnixMilli() - parsed.UnixMilli()) / days
	if diff < 0 {
		return -1, -1
	}

	if len(dateAndNumber) < 2 {
		return diff, 0
	}

	// calculate offset => ${specialChar}${number}
	offset, err := strconv.Atoi(dateAndNumber[1])
	if err != nil {
		return -1, -1
	}
	return diff, offset
}

func (d *DateComparatorII) Replace(i int) string {
	date := fmt.Sprintf("%s%s", d.sep, pkg.GetDateByFormat(d.datePattern))
	if i != 0 {
		date = fmt.Sprintf("%s%s%d", date, d.sep, i)
	}
	return strings.Replace(d.pattern, string(SpecialChar), date, -1)
}

func (d *DateComparatorII) ReplaceWithDifferentDate(i int, date time.Time) string {
	dateStr := fmt.Sprintf("%s%s", d.sep, date.Format(d.datePattern))
	if i != 0 {
		dateStr = fmt.Sprintf("%s%s%d", dateStr, d.sep, i)
	}
	return strings.Replace(d.pattern, string(SpecialChar), dateStr, -1)
}
