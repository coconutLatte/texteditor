package main

import (
	"fmt"
	"github.com/coconutLatte/texteditor"
	"os"
)

var staticData = `Hello
World!
`

func main() {
	d, err := texteditor.EditorStatic([]byte(staticData))
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(d))

	os.WriteFile("res", d, 0644)
}
