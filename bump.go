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
	prefix, segments, err := parseVersion(old)
	if err != nil {
		return "", err
	}

	if len(segments) == 0 {
		return "", fmt.Errorf("no version segments found in '%s'", old)
	} else if index >= len(segments) {
		return "", fmt.Errorf("segment index out of range: %d", index)
	} else if index < 0 {
		index = len(segments) - 1
	}

	segments[index]++
	for {
		index++
		if index >= len(segments) {
			break
		}
		segments[index] = 0
	}

	ret := strings.Builder{}
	ret.WriteString(prefix)
	for i, seg := range segments {
		ret.WriteString(strconv.Itoa(seg))
		if i < len(segments)-1 {
			ret.WriteRune('.')
		}
	}

	return ret.String(), nil
}

func parseVersion(s string) (string, []int, error) {
	var prefixDone bool
	var prefix []rune
	var segment []rune
	var segments []int

	for i, r := range s {
		if !prefixDone && unicode.IsNumber(r) {
			prefixDone = true
		} else if !prefixDone {
			prefix = append(prefix, r)
			continue
		}

		if unicode.IsNumber(r) {
			segment = append(segment, r)
		} else if r != '.' {
			return "", segments,
				fmt.Errorf("version parse failed: unexpected character: %s",
					strconv.QuoteRune(r))
		}

		if r == '.' || i == len(s)-1 {
			if len(segment) == 0 {
				return "", segments, fmt.Errorf("parse failed: unexpected '.'")
			}
			int, err := strconv.Atoi(string(segment))
			if err != nil {
				return "", segments, fmt.Errorf("parse failed: %w", err)
			}
			segments = append(segments, int)
			segment = []rune{}
		}
	}

	return string(prefix), segments, nil
}
