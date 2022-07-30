# HNotes

`hnotes` is an HTTP file server that converts Markdown to HTML before serving
it. It also provides directory viewing functionality.

## Install

`go install hellerwach.com/go/hnotes@latest` tested with the latest Go version
(1.18.3).

## Setup

You might consider making `hnotes` available system-wide:
`export PATH=~/go/bin:$PATH`. Now you can call it anywhere by typing `hnotes`
in your terminal.

Then copy the directory `.config/hnotes` from the repository to
`~/.config/hnotes`. It contains all the files needed to run the program.

Now you are ready to run `hnotes` anywhere in your system, preferably in your
home directory.

## Usage

`hnotes` is an HTTP server run by default on port `4545`, you can change this
behavior with the `-p` or `--port` flag. You can access it in your browser
on `http://localhost:4545/` and you will see a directory view of the directory
it was started in. Clicking on a file or directory opens it.

_However_, the special part, Markdown files (`*.md`) are converted to HTML and
inserted into a template before serving. This allows you to write them in your
editor of choice and preview or convert them to a PDF in your browser.

If you are working in a CLI, you can also call `hnotes new` followed by a
variable number of filenames to create these files from a template.

## Customization

`hnotes` is making heavy usage of Go's `html/template` library and is thus very
customizable. All the files are located at `~/.config/hnotes`.

`hnotes` provides additional functions for use in the templates:

| Name | Documentation |
|------|---------------|
|filepathJoin|[filepath.Join](https://pkg.go.dev/filepath#Join)|
|hasPrefix|[strings.HasPrefix](https://pkg.go.dev/strings#HasPrefix)|
|hasSuffix|[strings.HasSuffix](https://pkg.go.dev/strings#HasSuffix)|

### Markdown to HTML template

The Markdown to HTML template is by default located in the config directory at
`templates/single.html`. Due to the already huge documentation of Go's
templating system (for links see down below), I will not give any tutorial
here, but just tell you the data you are working with:

[embedmd]:# (note/note.go go /^type Metadata/ /^}$/)
```go
type Metadata map[string]interface{}

type Note struct {
	Metadata
	Content template.HTML
}
```

#### Documentation

- [template.HTML](https://pkg.go.dev/html/template#HTML)

### Directory view template

The directory view template is by default located in the config directory at
`templates/dir.html`. The same as [above](#markdown-to-html-template) applies
here.

[embedmd]:# (server/server.go go /^type dir struct {$/ /^}$/)
```go
type dir struct {
	Path    string
	Entries []fs.DirEntry
	note.Note
}
```

#### Documentation

- [fs.DirEntry](https://pkg.go.dev/io/fs#DirEntry)
- [note.Note](#markdown-to-html-template)

### Markdown template

The Markdown template **is** located in the config directory at
`templates/single.md`. It will be copied uncoditionally by `hnotes new`.

### Extensions
`hnotes` only provides minimal functionality, but it is extensible. You can
put shell scripts or executables in `.config/hnotes/extensions` and if the
first command line argument is not a predefined sub command or a flag, it will
try to execute the shell script or executable and pass the arguments given.

Example with the `hello` script in the [`.config/hnotes/extensions`](.config/hnotes/extensions/hello):

```
$ hnotes hello world
hello world
```

# Go's templating system

Some tutorials about Go's templating system can be found here:

- <https://pkg.go.dev/html/template>
- <https://learn.hashicorp.com/tutorials/nomad/go-template-syntax>
- <https://zetcode.com/golang/template/>
- <https://blog.gopheracademy.com/advent-2017/using-go-templates/>
