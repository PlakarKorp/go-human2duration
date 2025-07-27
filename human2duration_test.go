package human2duration

import (
	"testing"
	"time"
)

var tests = []struct {
	input    string
	expected time.Duration
}{
	// Fuzzy expressions
	{"half an hour", 30 * time.Minute},
	{"an hour and a half", 90 * time.Minute},
	{"half a day", 12 * time.Hour},
	{"couple of minutes", 2 * time.Minute},
	{"couple of hours", 2 * time.Hour},
	{"couple of days", 48 * time.Hour},
	{"an hour", time.Hour},
	{"a minute", time.Minute},
	{"a second", time.Second},
	{"a day", 24 * time.Hour},
	{"a week", 7 * 24 * time.Hour},
	{"a month", 30 * 24 * time.Hour},

	// Standard units
	{"1s", time.Second},
	{"2m", 2 * time.Minute},
	{"3h", 3 * time.Hour},
	{"4d", 96 * time.Hour},
	{"5w", 5 * 7 * 24 * time.Hour},
	{"6y", 6 * 365 * 24 * time.Hour},
	{"1mo", 30 * 24 * time.Hour},
	{"2mo", 60 * 24 * time.Hour},

	// Decimal inputs
	{"1.5h", 90 * time.Minute},
	{"2.25d", time.Duration(2.25 * float64(24*time.Hour))},
	{"3.75w", time.Duration(3.75 * float64(7*24*time.Hour))},

	// Mixed combinations
	{"1 day 4h 30m", 28*time.Hour + 30*time.Minute},
	{"2d 3h 15min", 2*24*time.Hour + 3*time.Hour + 15*time.Minute},
	{"1y 1month", 365*24*time.Hour + 30*24*time.Hour},
	{"3w2d", 3*7*24*time.Hour + 2*24*time.Hour},
	{"1y 3mo", 365*24*time.Hour + 90*24*time.Hour},
	{"2y6mo3w", 2*365*24*time.Hour + 6*30*24*time.Hour + 3*7*24*time.Hour},
	{"1d 12h", 36 * time.Hour},

	// Native duration fallback
	{"90m", 90 * time.Minute},
	{"2h45m", 2*time.Hour + 45*time.Minute},
}

func TestHuman2Duration(t *testing.T) {
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			got, err := ParseDuration(test.input)
			if err != nil {
				t.Errorf("ParseDuration(%q) returned error: %v", test.input, err)
				return
			}
			if got != test.expected {
				t.Errorf("ParseDuration(%q) = %v; want %v", test.input, got, test.expected)
			}
		})
	}
}

func TestHuman2Duration_Errors(t *testing.T) {
	invalid := []string{
		"",
		"nonsense",
		"1lightyear",
		"two hours",
		"half banana",
		"123abc",
	}

	for _, input := range invalid {
		t.Run("invalid:"+input, func(t *testing.T) {
			_, err := ParseDuration(input)
			if err == nil {
				t.Errorf("expected error for input %q but got none", input)
			}
		})
	}
}

func TestParseSinceDuration(t *testing.T) {
	testsAgo := []struct {
		input    string
		expected time.Duration
		wantErr  bool
	}{
		{"1h ago", -time.Hour, false},
		{"half an hour ago", -30 * time.Minute, false},
		{"ago", 0, true}, // invalid
	}

	for _, test := range testsAgo {
		t.Run("since-"+test.input, func(t *testing.T) {
			got, err := ParseSinceDuration(test.input)
			if test.wantErr && err == nil {
				t.Errorf("expected error for %q but got none", test.input)
			}
			if !test.wantErr && got != test.expected {
				t.Errorf("expected %v, got %v", test.expected, got)
			}
		})
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			got, err := ParseSinceDuration(test.input)
			if err != nil {
				t.Errorf("ParseDuration(%q) returned error: %v", test.input, err)
				return
			}
			if got != -test.expected {
				t.Errorf("ParseDuration(%q) = %v; want %v", test.input, got, test.expected)
			}
		})
	}
}

func TestParseAfterDuration(t *testing.T) {
	testsAfter := []struct {
		input    string
		expected time.Duration
		wantErr  bool
	}{
		{"in 2h", 2 * time.Hour, false},
		{"after 1 day", 24 * time.Hour, false},
		{"after", 0, true}, // invalid
	}

	for _, test := range testsAfter {
		t.Run("after-"+test.input, func(t *testing.T) {
			got, err := ParseAfterDuration(test.input)
			if test.wantErr && err == nil {
				t.Errorf("expected error for %q but got none", test.input)
			}
			if !test.wantErr && got != test.expected {
				t.Errorf("expected %v, got %v", test.expected, got)
			}
		})
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			got, err := ParseAfterDuration(test.input)
			if err != nil {
				t.Errorf("ParseDuration(%q) returned error: %v", test.input, err)
				return
			}
			if got != test.expected {
				t.Errorf("ParseDuration(%q) = %v; want %v", test.input, got, test.expected)
			}
		})
	}
}
