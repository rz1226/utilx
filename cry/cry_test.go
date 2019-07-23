package cry

import (
	"fmt"
	"testing"
)

func Test_cry(t *testing.T) {

	test := "abc"

	fmt.Println(GetMd5(test))

	test2 := "abcsdfsdfsdf234234234c242c43"

	fmt.Println(GetMd5(test2))

	test11 := "abc"

	fmt.Println(GetMd52(test11))

	test22 := "abcsdfsdfsdf234234234c242c43"

	fmt.Println(GetMd52(test22))
}
