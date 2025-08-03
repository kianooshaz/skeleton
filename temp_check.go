package main

import (
	"fmt"

	"github.com/kianooshaz/skeleton/foundation/derror"
)

func main() {
	fmt.Println("ErrUserNotFound:", derror.ErrUserNotFound.Error())
	fmt.Println("ErrUserIDRequired:", derror.ErrUserIDRequired.Error())
}
