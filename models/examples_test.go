package models

import "fmt"

func ExampleValidate() {
	a := &Article{
		Title:   "TestTitle",
		Content: "TestContent",
	}

	err := a.Validate()
	fmt.Println(err)
	// Output:
	// author: cannot be blank.
}
