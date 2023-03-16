package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/aageorg/altclient"
	"os"
	"regexp"
	"sync"
)

var quiet bool

func PrintInfo(s string) {
	if quiet == false {
		fmt.Fprintln(os.Stdout, s)

	}
}

type Response [2]Diff

type Diff struct {
	BranchName string `json:"branch"`

	// The map with a missing packages in the branch. Key is the arch name (e.g. "aarch64")

	Missing map[string][]altclient.Package `json:"missing"`

	// The map with the out of date packages in the branch. Key is the arch name (e.g. "aarch64")

	OutOfDate map[string][]altclient.Package `json:"out_of_date,omitempty"`
}

func main() {
	flag.BoolVar(&quiet, "q", false, "Result output without additional messages")
	flag.Parse()
	var args []string
	if quiet {
		args = flag.Args()
	} else {
		args = os.Args[1:]
	}
	if len(args) != 2 {
		fmt.Fprintln(os.Stdout, "Usage: "+os.Args[0]+" [-q] -b branch1 branch2")
		return
	}
	re := regexp.MustCompile(`^[a-zA-Z0-9\_\-]+$`)
	if !re.Match([]byte(args[0])) || !re.Match([]byte(args[1])) {
		fmt.Fprintln(os.Stderr, "Invalid branch name. Usage: "+os.Args[0]+" [-q] branch1 branch2")

		return
	}
	if args[0] == args[1] {
		fmt.Fprintln(os.Stderr, "Invalid branch name. Cannot compare branch "+args[0]+" with itself")
		return
	}
	var branch1 *altclient.Branch
	var branch2 *altclient.Branch
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		PrintInfo("Downloading of package list from branch " + args[0] + " is started")
		var err error
		branch1, err = altclient.NewBranch(args[0])
		if err != nil {
			fmt.Fprintln(os.Stderr, "Cannot get branch "+args[0]+": "+err.Error())
			os.Exit(1)
		}
		PrintInfo("Packages list from branch " + args[0] + " is downloaded")
	}()
	go func() {
		defer wg.Done()
		PrintInfo("Downloading of package list from branch " + args[1] + " is started")
		var err error
		branch2, err = altclient.NewBranch(args[1])
		if err != nil {
			fmt.Fprintln(os.Stderr, "Cannot get branch "+args[1]+": "+err.Error())
			os.Exit(1)
		}
		PrintInfo("Packages list from branch " + args[1] + " is downloaded")
	}()
	wg.Wait()
	wg.Add(2)

	var d1 = Diff{BranchName: args[0], Missing: make(map[string][]altclient.Package)}
	var d2 = Diff{BranchName: args[1], Missing: make(map[string][]altclient.Package), OutOfDate: make(map[string][]altclient.Package)}

	go func() {
		defer wg.Done()
		archs := branch1.GetArchs()
		for _, a := range archs {
			d1.Missing[a] = branch2.GetMissing(branch1, a)
		}
		PrintInfo("Search of missing packages in branch " + args[0] + " is finished")

	}()
	go func() {
		defer wg.Done()
		archs := branch2.GetArchs()
		PrintInfo("Search of missing and obsolete packages in branch " + args[1] + " is started")
		for _, a := range archs {
			d2.Missing[a] = branch1.GetMissing(branch2, a)
			d2.OutOfDate[a] = branch1.GetOutOfDate(branch2, a)
		}
		PrintInfo("Search of missing and obsolete packages in branch " + args[1] + " is finished")
	}()
	wg.Wait()

	var r Response
	r[0] = d1
	r[1] = d2

	js, err := json.Marshal(r)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	fmt.Fprintln(os.Stdout, string(js))
}
