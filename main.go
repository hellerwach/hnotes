// HNotes is an HTTP file server that converts Markdown to HTML before serving
// it. It also provides directory viewing functionality.
package main // go install hellerwach.com/go/hnotes@latest

import "hellerwach.com/go/hnotes/cmd"

func main() {
	cmd.Execute()
}
