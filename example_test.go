package errors_test

import (
	"fmt"
	"github.com/betNevS/errors"
)

func ExampleNew() {
	err := errors.New("whoops")
	fmt.Println(err)

	// Output: whoops
}
