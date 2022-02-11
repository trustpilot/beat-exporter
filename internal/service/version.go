package service

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
)

var (
	releaseVersion = "None"
	buildDate      = "None"
	gitHash        = "None"
	gitBranch      = "None"
)

type Info struct {
	ReleaseVersion string `json:"release_version"`
	GitHash        string `json:"git_hash"`
	GitBranch      string `json:"git_branch"`
	BuildDate      string `json:"build_date"`
	GoVersion      string `json:"go_version"`
	Compiler       string `json:"compiler"`
	Platform       string `json:"platform"`
}

func newVersionInfo() *Info {
	return &Info{
		ReleaseVersion: fmt.Sprintf("beat-exporter[%s]", releaseVersion),
		GitHash:        gitHash,
		GitBranch:      gitBranch,
		BuildDate:      buildDate,
		GoVersion:      runtime.Version(),
		Compiler:       runtime.Compiler,
		Platform:       fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}

func PrintVersion(versionFlag bool) bool {
	if versionFlag {
		marshaled, err := json.MarshalIndent(newVersionInfo(), "", "  ")
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		fmt.Println(string(marshaled))
		return true
	}
	return false
}
