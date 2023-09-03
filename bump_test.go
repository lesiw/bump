package main

import (
	"fmt"
	"testing"
)

func TestParseSegment(t *testing.T) {
	type testCase struct {
		in  string
		seg int
		err bool
	}
	testCases := []testCase{{
		in:  "1",
		seg: 1,
		err: false,
	}, {
		in:  "",
		seg: -1,
		err: false,
	}, {
		in:  "x",
		seg: 0,
		err: true,
	}, {
		in:  "major",
		seg: 0,
		err: false,
	}, {
		in:  "minor",
		seg: 1,
		err: false,
	}, {
		in:  "patch",
		seg: 2,
		err: false,
	}}

	for _, tc := range testCases {
		test := func(t *testing.T) {
			res, err := parseSegment(tc.in)
			if !checkErr(err, tc.err) {
				t.Errorf("wanted error: %t, but got %t", tc.err, !tc.err)
			}
			if tc.seg != res {
				t.Errorf("want %d, got %d", tc.seg, res)
			}
		}
		var inStr string
		if tc.in == "" {
			inStr = "EMPTY"
		} else {
			inStr = fmt.Sprintf("'%s'", tc.in)
		}
		t.Run(fmt.Sprintf("parseseg_%s_%d", inStr, tc.seg), test)
	}
}

func TestBumpVersion(t *testing.T) {
	type testCase struct {
		in  string
		seg int
		out string
	}
	testCases := []testCase{{
		in:  "1.0.0",
		seg: 1,
		out: "1.1.0",
	}, {
		in:  "1.0.0",
		seg: 0,
		out: "2.0.0",
	}, {
		in:  "1.0.0",
		seg: 2,
		out: "1.0.1",
	}, {
		in:  "1.2.3",
		seg: 1,
		out: "1.3.0",
	}, {
		in:  "1.2.3",
		seg: 0,
		out: "2.0.0",
	}, {
		in:  "100.18.42",
		seg: 1,
		out: "100.19.0",
	}, {
		in:  "v1.2.3",
		seg: 2,
		out: "v1.2.4",
	}, {
		in:  "bigprefix13.17.19",
		seg: 0,
		out: "bigprefix14.0.0",
	}, {
		in:  "v1",
		seg: 0,
		out: "v2",
	}, {
		in:  "1.2.3.4",
		seg: 1,
		out: "1.3.0.0",
	}, {
		in:  "1.2.3",
		seg: -1,
		out: "1.2.4",
	}}

	for _, tc := range testCases {
		test := func(t *testing.T) {
			res, err := bumpVersion(tc.in, tc.seg)
			if err != nil {
				t.Errorf("err: %s", err)
			}
			if tc.out != res {
				t.Errorf("want %s, got %s", tc.out, res)
			}
		}
		t.Run(fmt.Sprintf("%s_bump_%d", tc.in, tc.seg), test)
	}
}

func TestBumpVersionOutOfRange(t *testing.T) {
	_, err := bumpVersion("1.2.3", 3)
	if err == nil {
		t.Fatalf("expected error")
	}
	want := "segment index out of range: 3"
	if err.Error() != want {
		t.Errorf("want '%s', got '%s'", want, err.Error())
	}
}

func checkErr(err error, present bool) bool {
	if err == nil && !present {
		return true
	} else if err != nil && present {
		return true
	}
	return false
}
