package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"unicode"
)

type version struct {
	prefix     string
	segments   []int
	prerelease string
	tag        string
}

type versionParser func(*version, *bufio.Reader) (versionParser, error)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run() error {
	segstr := flag.String("s", "", "index of segment to bump")
	flag.Parse()

	seg, err := parseSegment(*segstr)
	if err != nil {
		return fmt.Errorf("error parsing segment: %w", err)
	}

	input, err := readInput(os.Stdin)
	if err != nil {
		return fmt.Errorf("error reading stdin: %w", err)
	}

	output, err := bumpVersion(input, seg)
	if err != nil {
		return err
	}

	fmt.Println(output)

	return nil
}

func parseSegment(s string) (int, error) {
	if s == "" {
		return -1, nil
	}

	ret, err := strconv.Atoi(s)
	if err != nil {
		switch s {
		case "major":
			return 0, nil
		case "minor":
			return 1, nil
		case "patch":
			return 2, nil
		case "pre":
			return 3, nil
		default:
			return 0, fmt.Errorf("unrecognized segment: '%s'", s)
		}
	}

	return ret, nil
}

func readInput(reader io.Reader) (string, error) {
	r := bufio.NewReader(reader)
	input, err := r.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(input), nil
}

func bumpVersion(s string, index int) (string, error) {
	v, err := newVersion(s)
	if err != nil {
		return "", err
	}

	if len(v.segments) == 0 {
		return "", fmt.Errorf("no version segments found in '%s'", s)
	} else if index > len(v.segments) {
		return "", fmt.Errorf("segment index out of range: %d", index)
	} else if index < 0 {
		index = len(v.segments) - 1
	}

	if index == len(v.segments) {
		if v.prerelease == "" {
			v.segments[len(v.segments)-1]++
			v.prerelease = "rc.1"
		} else {
			var ok bool
			v.prerelease, ok = bumpLastDigitRun(v.prerelease)
			if !ok {
				v.prerelease += ".1"
			}
		}
	} else if index == len(v.segments)-1 && v.prerelease != "" {
		v.prerelease = ""
	} else {
		v.segments[index]++
		for index++; index < len(v.segments); index++ {
			v.segments[index] = 0
		}
		v.prerelease = ""
	}

	return v.String(), nil
}

func newVersion(s string) (*version, error) {
	v := &version{}
	r := bufio.NewReader(strings.NewReader(s))
	parseState := parseVersionPrefix

	for parseState != nil {
		var err error
		parseState, err = parseState(v, r)
		if err != nil {
			return nil, err
		}
	}

	return v, nil
}

func (v *version) String() string {
	var b strings.Builder
	b.WriteString(v.prefix)
	for i, seg := range v.segments {
		b.WriteString(strconv.Itoa(seg))
		if i < len(v.segments)-1 {
			b.WriteRune('.')
		}
	}
	if v.prerelease != "" {
		b.WriteRune('-')
		b.WriteString(v.prerelease)
	}
	if v.tag != "" {
		b.WriteRune('+')
		b.WriteString(v.tag)
	}
	return b.String()
}

func parseVersionPrefix(v *version, reader *bufio.Reader) (versionParser, error) {
	var prefix []rune
	defer func() { v.prefix = string(prefix) }()

	for {
		r, _, err := reader.ReadRune()
		if err == io.EOF {
			return nil, nil
		}
		if unicode.IsNumber(r) {
			_ = reader.UnreadRune()
			return parseVersionSegments, nil
		}
		prefix = append(prefix, r)
	}
}

func parseVersionSegments(v *version, reader *bufio.Reader) (versionParser, error) {
	var segment []rune

	storeSegment := func() {
		if len(segment) == 0 {
			return
		}
		int, err := strconv.Atoi(string(segment))
		if err != nil {
			return
		}
		v.segments = append(v.segments, int)
		segment = []rune{}
	}
	defer storeSegment()

	for {
		r, _, err := reader.ReadRune()
		if err == io.EOF {
			return nil, nil
		}

		switch r {
		case '.':
			if len(segment) == 0 {
				return nil, fmt.Errorf("parse failed: unexpected '.'")
			}
			storeSegment()
			continue
		case '-':
			return parseVersionPrerelease, nil
		case '+':
			return parseVersionTag, nil
		}

		if unicode.IsNumber(r) {
			segment = append(segment, r)
		} else {
			return nil, fmt.Errorf("version parse failed: unexpected character: %s",
				strconv.QuoteRune(r))
		}
	}
}

func parseVersionPrerelease(v *version, reader *bufio.Reader) (versionParser, error) {
	var prerelease []rune
	defer func() { v.prerelease = string(prerelease) }()

	for {
		r, _, err := reader.ReadRune()
		if err == io.EOF {
			return nil, nil
		}
		if r == '+' {
			_ = reader.UnreadRune()
			return parseVersionSegments, nil
		}
		prerelease = append(prerelease, r)
	}
}

func parseVersionTag(v *version, reader *bufio.Reader) (versionParser, error) {
	var tag []rune
	defer func() { v.tag = string(tag) }()

	for {
		r, _, err := reader.ReadRune()
		if err == io.EOF {
			return nil, nil
		}
		tag = append(tag, r)
	}
}

func bumpLastDigitRun(s string) (string, bool) {
	var digits []rune
	pos := -1
	for i, r := range s {
		if unicode.IsDigit(r) {
			if i-len(digits) != pos {
				digits = []rune{}
				pos = i
			}
			digits = append(digits, r)
		}
	}
	if pos < 0 {
		return s, false
	}

	int, err := strconv.Atoi(string(digits))
	if err != nil {
		// Only digits are added to the run, so this should never happen.
		panic(fmt.Sprintf("failed to parse digit run: %s", err))
	}
	int++

	var b strings.Builder
	b.WriteString(string([]rune(s)[:pos]))
	b.WriteString(fmt.Sprintf("%0"+fmt.Sprint(len(digits))+"d", int))
	b.WriteString(string([]rune(s)[pos+len(digits):]))

	return b.String(), true
}
