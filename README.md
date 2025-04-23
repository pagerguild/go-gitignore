# go-gitignore

A gitignore parser for `Go`

> This is a fork of the original [go-gitignore](https://github.com/sabhiram/go-gitignore) by Shaba Abhiram. Licensed under MIT.

## Install

```shell
go get github.com/pagerguild/go-gitignore
```

## Usage

For a quick sample of how to use this library, check out the tests under `ignore_test.go`.

### Basic Usage

```go
// Compile gitignore patterns from strings
ignoreObj := ignore.CompileIgnoreLines([]string{
    "node_modules",
    "*.out", 
    "foo/*.c"
}...)

// Check if a path matches any pattern
if ignoreObj.MatchesPath("node_modules/test/foo.js") {
    // This path would be ignored
}

// Check with explanation of which pattern matched
matchesPath, reason := ignoreObj.MatchesPathHow("node_modules/test/foo.js")
if matchesPath {
    fmt.Printf("Path matched pattern: %s (line %d)\n", reason.Line, reason.LineNo)
}

// Load patterns from a .gitignore file
ignoreFromFile, err := ignore.CompileIgnoreFile("path/to/.gitignore")
if err != nil {
    log.Fatal(err)
}

// Load patterns from both file and additional lines
ignoreObj, err := ignore.CompileIgnoreFileAndLines("path/to/.gitignore", "**/foo", "bar")
```

### Working with Directory Trees

```go
// Create a gitignore tree that respects all .gitignore files in a directory structure
ignoreTree, err := ignore.NewIgnoreMap("/path/to/project")
if err != nil {
    log.Fatal(err)
}

// Check if a path should be ignored
if ignoreTree.Ignore("/path/to/project/node_modules/file.js") {
    // Path is ignored
}

// Iterate over all non-ignored files
for err, path := range ignoreTree.RegularFilesSeq() {
    if err != nil {
        // that means there was an issue in the underlying directory traversal
    }
    process(path)
}
```
