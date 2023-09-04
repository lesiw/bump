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
	prefix   string
	segments []int
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

func bumpVersion(old string, index int) (string, error) {
	v, err := parseVersion(old)
	if err != nil {
		return "", err
	}

	if len(v.segments) == 0 {
		return "", fmt.Errorf("no version segments found in '%s'", old)
	} else if index >= len(v.segments) {
		return "", fmt.Errorf("segment index out of range: %d", index)
	} else if index < 0 {
		index = len(v.segments) - 1
	}

	v.segments[index]++
	for index++; index < len(v.segments); index++ {
		v.segments[index] = 0
	}

	ret := strings.Builder{}
	ret.WriteString(v.prefix)
	for i, seg := range v.segments {
		ret.WriteString(strconv.Itoa(seg))
		if i < len(v.segments)-1 {
			ret.WriteRune('.')
		}
	}

	return ret.String(), nil
}

func parseVersion(s string) (*version, error) {
	v := &version{}
	r := bufio.NewReader(strings.NewReader(s))
	parseState := parseVersionPrefix

	for {
		var err error
		parseState, err = parseState(v, r)
		if err != nil {
			return nil, err
		}
		if parseState == nil {
			break
		}
	}

	return v, nil
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
	for {
		r, _, err := reader.ReadRune()
		if err == io.EOF {
			return nil, nil
		}

		if unicode.IsNumber(r) {
			segment = append(segment, r)
			if !bufend(reader) {
				continue
			}
		}

		if r == '.' && len(segment) == 0 {
			return nil, fmt.Errorf("parse failed: unexpected '.'")
		} else if !unicode.IsNumber(r) && r != '.' {
			return nil, fmt.Errorf("version parse failed: unexpected character: %s",
				strconv.QuoteRune(r))
		}

		int, err := strconv.Atoi(string(segment))
		if err != nil {
			return nil, fmt.Errorf("parse failed: %w", err)
		}

		v.segments = append(v.segments, int)
		segment = []rune{}
	}
}

func bufend(reader *bufio.Reader) bool {
	_, _, nextErr := reader.ReadRune()
	_ = reader.UnreadRune()
	return (nextErr == io.EOF)
}
