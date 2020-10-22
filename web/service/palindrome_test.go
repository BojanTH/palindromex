package service

import "testing"

func TestIsPalindrome(t *testing.T) {
	tables := []struct {
		text string
		expected bool
	} {
		{"racecar", true},
		{"madam", true},
		{"11/11/11 11:11", true},
		{"02/02/2020", true},
		{"a bccba!", true},
		{"a b1ba!", true},
		{"öö", true},
		{"Mus rev inuits öra, sa röst i universum.", true},

		{"foo", false},
		{"9922", false},
		{"ab1ca!", false},
		{"abca!", false},
		{"Some long text.", false},
	}

	for _, table := range tables {
		actual := IsPalindrome(table.text)
		if actual != table.expected {
			t.Errorf("Unexpected result for text: '%s'", table.text)
		}
	}
}
