package main

import (
	"encoding/json"
	"fmt"
	"github.com/aageorg/altclient"
	"os"
	"regexp"
	"sync"
)

type Response [2]Diff

type Diff struct {
	BranchName string `json:"branch"`

	// The map with a missing packages in the branch. Key is the arch name (e.g. "aarch64")

	Missing map[string][]altclient.Package `json:"missing"`

	// The map with the out of date packages in the branch. Key is the arch name (e.g. "aarch64")

	OutOfDate map[string][]altclient.Package `json:"out_of_date,omitempty"`
}

func main() {

	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stdout, "Usage: branch1 branch2")
		return
	}
	re := regexp.MustCompile(`^[a-zA-Z0-9\_\-]+$`)
	if !re.Match([]byte(os.Args[1])) || !re.Match([]byte(os.Args[2])) {
		fmt.Fprintln(os.Stderr, "Invalid branch name. Usage: compare branch1 branch2")

		return
	}
	if os.Args[1] == os.Args[2] {
		fmt.Fprintln(os.Stderr, "Invalid branch name. Cannot compare branch "+os.Args[1]+" with itself")
		return
	}
	var branch1 *altclient.Branch
	var branch2 *altclient.Branch
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		var err error
		branch1, err = altclient.NewBranch(os.Args[1])
		if err != nil {
			fmt.Fprintln(os.Stderr, "Cannot get branch "+os.Args[1]+": "+err.Error())
			os.Exit(1)
		}
	}()
	go func() {
		defer wg.Done()
		var err error
		branch2, err = altclient.NewBranch(os.Args[2])
		if err != nil {
			fmt.Fprintln(os.Stderr, "Cannot get branch "+os.Args[2]+": "+err.Error())
			os.Exit(1)
		}
	}()
	wg.Wait()
	wg.Add(2)

	var d1 = Diff{BranchName: os.Args[1], Missing: make(map[string][]altclient.Package)}
	var d2 = Diff{BranchName: os.Args[2], Missing: make(map[string][]altclient.Package), OutOfDate: make(map[string][]altclient.Package)}

	go func() {
		defer wg.Done()
		archs := branch1.GetArchs()
		for _, a := range archs {
			d1.Missing[a] = branch2.GetMissing(branch1, a)
		}
	}()
	go func() {
		defer wg.Done()
		archs := branch2.GetArchs()
		for _, a := range archs {
			d2.Missing[a] = branch1.GetMissing(branch2, a)
			d2.OutOfDate[a] = branch1.GetOutOfDate(branch2, a)
		}
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
