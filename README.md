<img align=right width="168" src="docs/gouef_logo.png">

# Finder
`Finder` is a simple Go package for searching files and directories based on patterns in a specified directory. It also supports excluding certain files or directories and provides flexible configuration options for the search.

[![Static Badge](https://img.shields.io/badge/Github-gouef%2Ffinder-blue?style=for-the-badge&logo=github&link=github.com%2Fgouef%2Ffinder)](https://github.com/gouef/finder)

[![GoDoc](https://pkg.go.dev/badge/github.com/gouef/finder.svg)](https://pkg.go.dev/github.com/gouef/finder)
[![GitHub stars](https://img.shields.io/github/stars/gouef/finder?style=social)](https://github.com/gouef/finder/stargazers)
[![Go Report Card](https://goreportcard.com/badge/github.com/gouef/finder)](https://goreportcard.com/report/github.com/gouef/finder)
[![codecov](https://codecov.io/github/gouef/finder/branch/main/graph/badge.svg?token=YUG8EMH6Q8)](https://codecov.io/github/gouef/finder)

## Versions
![Stable Version](https://img.shields.io/github/v/release/gouef/finder?label=Stable&labelColor=green)
![GitHub Release](https://img.shields.io/github/v/release/gouef/finder?label=RC&include_prereleases&filter=*rc*&logoSize=diago)
![GitHub Release](https://img.shields.io/github/v/release/gouef/finder?label=Beta&include_prereleases&filter=*beta*&logoSize=diago)

## Installation

If you are using Go modules, you can add this package to your project by running:

```bash
go get -u github.com/gouef/finder
```

## Usage

### Basic FIle and Directory Search

```go
package main

import (
	"fmt"
	"github.com/gouef/finder"
)

func main() {
	// Search for all files in the current directory and its subdirectories
	files := finder.New().In(".").FindFiles("*.go").Get()

	// Print the found files
	for _, file := range files {
		fmt.Println(file.Name())
	}
}

```

### Filtering Files by Pattern
You can search for files or directories by specifying patterns using `FindFiles` or `FindDirectories`:

```go
// Search for .txt files in the current directory and subdirectories
files := finder.New().In(".").FindFiles("*.txt").Get()
```

### Excluding Files and Directories
You can exclude specific files or directories from the search results using the `Exclude` function:

```go
// Search for .txt files, but exclude "test1.txt"
files := finder.New().In(".").FindFiles("*.txt").Exclude("test1.txt").Get()
```

### Searching for Directories
If you want to search for directories instead of files:

```go
// Search for directories in the current directory and subdirectories
dirs := finder.New().In(".").FindDirectories("*").Get()
```

## Functions
- `New()` – Creates a new instance of Finder.
- `In(dirs ...string)` – Specifies the directories to search in.
- `Find(patterns ...string)` – Searches for both files and directories based on the given patterns.
- `FindFiles(patterns ...string)` – Searches only for files matching the given patterns.
- `FindDirectories(patterns ...string)` – Searches only for directories matching the given patterns.
- `Exclude(patterns ...string)` – Excludes files and directories matching the given patterns from the search results.
- `Get()` – Retrieves the search results.


## Example

```go
package main

import (
	"fmt"
	"github.com/gouef/finder"
)

func main() {
	// Search for .go files and exclude "main.go"
	files := finder.FindFiles("*.go").In(".").Exclude("main.go").Get()

	// Print the found files
	for _, file := range files {
		fmt.Println(file.Name())
	}
}

```


## Contributing

Read [Contributing](CONTRIBUTING.md)

## Contributors

<div>
<span>
  <a href="https://github.com/actions-user"><img src="https://raw.githubusercontent.com/gouef/finder/refs/heads/contributors-svg/.github/contributors/actions-user.svg" alt="actions-user" /></a>
</span>
<span>
  <a href="https://github.com/JanGalek"><img src="https://raw.githubusercontent.com/gouef/finder/refs/heads/contributors-svg/.github/contributors/JanGalek.svg" alt="JanGalek" /></a>
</span>
</div>

## Join our Discord Community! 🎉

[![Discord](https://img.shields.io/discord/1334331501462163509?style=for-the-badge&logo=discord&logoColor=white&logoSize=auto&label=Community%20discord&labelColor=blue&link=https%3A%2F%2Fdiscord.gg%2FwjGqeWFnqK
)](https://discord.gg/wjGqeWFnqK)

Click above to join our community on Discord!
