package packagename

import "fmt"

func TestMake() {
	var mappa map[string]string

	mappa = make(map[string]string)
	mappa["iluha"] = "anna"
	mappa = make(map[string]string)
	fmt.Print(mappa)
}
