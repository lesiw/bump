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
	}, {
		in:  "pre",
		seg: 3,
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
	}, {
		in:  "1.2.3",
		seg: 3,
		out: "1.2.4-rc.1",
	}, {
		in:  "1.2.4-rc.1",
		seg: 3,
		out: "1.2.4-rc.2",
	}, {
		in:  "1.2.4-rc",
		seg: 3,
		out: "1.2.4-rc.1",
	}, {
		in:  "1.2.4-rc1",
		seg: 3,
		out: "1.2.4-rc2",
	}, {
		in:  "1.2.4-rc.42",
		seg: 2,
		out: "1.2.4",
	}, {
		in:  "1.2.3+sometag",
		seg: 2,
		out: "1.2.4+sometag",
	}, {
		in:  "1.2.3-rc.4",
		seg: 1,
		out: "1.3.0",
	}, {
		in:  "1.2.3-alpha.6",
		seg: 0,
		out: "2.0.0",
	}, {
		in:  "version 100.200.300-release-candidate---42+long.tag.with+symbols",
		seg: 0,
		out: "version 101.0.0+long.tag.with+symbols",
	}, {
		in:  "v1.2.3-rc.1.2.3",
		seg: 3,
		out: "v1.2.3-rc.1.2.4",
	}, {
		in:  "1.2-rc.002",
		seg: 2,
		out: "1.2-rc.003",
	}, {
		in:  "version 1---rc.042.1",
		seg: 1,
		out: "version 1---rc.042.2",
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
	_, err := bumpVersion("1.2.3", 4)
	if err == nil {
		t.Fatalf("expected error")
	}
	want := "segment index out of range: 4"
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
