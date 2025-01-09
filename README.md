# Termy :tv:

(_pronounced_ tˈɜːmi)

PKG **Termy** facilitates access to basic (terminal emulator's) functionality.

---

`termy.go` allows you to put your terminal in non-cooked modes.

```go
package main

import "github.com/mec-nyan/termy"


func main() {

	t := termy.New(int(os.Stdin.Fd()), false)
	err := t.Cbreaky()
	if err != nil {
		...
	}
	defer t.Restore()

	// Your code...

}
```

`escapy.go` provides function to send in-band control sequences to your terminal.

```go
func main() {

	// ...

	termy.SaveCurPos() // Save current cursor position.


	// Do stuff.

	termy.RestoreCurPos() // Go back to saved position.
	termy.ClearToEOS()    // Clean up.

	// ...
}
```

TODO: I'll soon add capabilities for styles, colours, etc.
