package version

import (
	"fmt"
	"runtime"
)

const version = "v1.0.0"

// Sets values from ldflags
var (
	Commit string
	Branch string
	Tag    string
)

type Version struct {
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	Branch    string `json:"branch"`
	Tag       string `json:"tag"`
	GoVersion string `json:"goVersion"`
	Compiler  string `json:"compiler"`
	Platform  string `json:"platform"`
}

func ShortVersionInfo() string {
	if Tag != "" && Tag != "undefined" {
		return Tag
	}
	v := version
	if len(Commit) >= 7 {
		v += "+" + Commit[:7]
	}
	return v
}

func Get() Version {
	return Version{
		Version:   ShortVersionInfo(),
		Commit:    Commit,
		Branch:    Branch,
		GoVersion: runtime.Version(),
		Compiler:  runtime.Compiler,
		Platform:  fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}
