package api

import (
	"fmt"
	"math/rand"
	"testing"
)

// TestCommands tests the handleCommand() function
func TestCommands(t *testing.T) {
	type Test struct {
		in, out  string
		hasError bool
	}

	rng := rand.New(rand.NewSource(1234))

	tests := [...]Test{
		{"/dice", fmt.Sprintf("Rolled a %d", rng.Intn(6)+1), false},
		{"/dice 10000", fmt.Sprintf("Rolled a %d", rng.Intn(10000)+1), false},
		{"/dice -10", "", true},
		{"/dice 0", "", true},
		{"/dice 10notanumber", "", true},
		{"/dice 598 more args", fmt.Sprintf("Rolled a %d", rng.Intn(598)+1), false},

		{"/coin", "Tails", false},
		{"/coin", "Tails", false},
		{"/coin", "Tails", false}, // that's rngesus for you

		{"/list", "Segmentation fault", false},
		{"/list a b c", "It's c", false},
		{"/list a b c", "Couldn't not be a", false},

		{"not a command", "", true},
		{"/doesntexist arg1 arg2", "", true},
	}

	rng = rand.New(rand.NewSource(1234))

	for _, test := range tests {
		out, err := processCommand(test.in, rng)
		if (err != nil) != test.hasError {
			if err == nil {
				t.Errorf("In <%s>, expected err!=nil, found err=%v\n", test.in, err)
			} else {
				t.Errorf("In <%s>, expected err==nil, found err=%v\n", test.in, err)
			}
		}
		if out != test.out {
			t.Errorf("In <%s>, expected out=%s, found out=%s\n", test.in, test.out, out)
		}
	}
}
