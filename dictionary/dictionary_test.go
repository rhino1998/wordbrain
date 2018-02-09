package dictionary

import (
	"fmt"
	"testing"
)

func TestFunction(t *testing.T) {
	d := New("hello", "hell", "welcome")
	fmt.Println(d.Contains("hello"))
	fmt.Println(d.Contains("hell"))
	fmt.Println(d.Contains("jell"))
	fmt.Println(d.Contains("wel"))
	fmt.Println(d.Contains("wek"))
}
