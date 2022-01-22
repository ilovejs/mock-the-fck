package tool

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"runtime"
	"sort"
	"strings"

	"golang.org/x/mod/semver"
)

// vscode keep old version plugin after updates complete.
// this script will save space by remove the old one.
// e.g.
//42 ms-vscode.vscode-typescript-next-4.6.20220120
//43 ms-vscode.vscode-typescript-next-4.6.20220121
//44 pranaygp.vscode-css-peek-4.2.0
//45 rvest.vs-code-prettier-eslint-3.1.0

func list(relativePath string) (res []string) {
	// "./" as path
	files, err := ioutil.ReadDir(relativePath)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		res = append(res, f.Name())
	}
	return
}

func currentUser() string {
	currentUser, err := user.Current()
	if err != nil {
		log.Fatalf(err.Error())
	}
	return currentUser.Username
}

func printPlugins(dirName string) {
	for i, fileName := range list(dirName) {
		fmt.Println(i, fileName)
	}
}

func compareVer(a, b string) int {
	return semver.Compare("v"+a, "v"+b)
}

// not efficient..
// https://thedeveloperblog.com/remove-duplicates-slice
func removeDuplicates(elements []string) []string {
	var result []string

	for i := 0; i < len(elements); i++ {

		exists := false
		for v := 0; v < i; v++ {
			if elements[v] == elements[i] {
				exists = true
				break
			}
		}

		if !exists {
			result = append(result, elements[i])
		}
	}
	return result
}

// RunClean must handle a case where plugin has 3+ version. It needs to delete oldest 2 version.
func RunClean() {

	// TODO: dryRun first
	dryRun := true
	//dryRun := false

	if runtime.GOOS != "darwin" {
		log.Fatal("Should only run under MacOS")
	}

	userName := currentUser()
	dirName := fmt.Sprintf("/Users/%s/.vscode/extensions", userName)
	// map filename to the latest version
	mFileName := map[string]string{}

	fileNames := list(dirName)
	sort.Strings(fileNames)

	var dups []string

	for _, fileName := range fileNames {
		lastSepIdx := strings.LastIndex(fileName, "-")
		if lastSepIdx == -1 {
			// reject if not a plugin format
			continue
		}

		// etl
		pluginVer := fileName[lastSepIdx+1:]
		pluginName := fileName[0:lastSepIdx]

		// map algo
		curVer, ok := mFileName[pluginName]
		if !ok {
			mFileName[pluginName] = pluginVer
		} else {
			if compareVer(curVer, pluginVer) < 0 {
				// save duplicated old value
				dups = append(dups, pluginName+"-"+curVer)
				// keep larger version
				mFileName[pluginName] = pluginVer
			}
		}
	}

	// [debug] unique list
	//for k, v := range mFileName {
	//	fmt.Println(k, v)
	//}

	// we can do safe rm by moving to tmp
	// ps. validate by
	// echo 2 > cc, run, cat /tmp/cc
	// ls -al /private/tmp

	// Knowledge:
	// 1. not like zsh `mv cc /tmp` can infer it's a dir then place under.
	// 2. you can rename folder as well as file
	// Correct way:
	// os.Rename("./cc", "/tmp/cc")

	if len(dups) > 0 {
		fmt.Println("\nCleaning------->")
	} else {
		fmt.Println("No Dups :)")
		return
	}

	for _, filename := range dups {
		if dryRun {
			fmt.Println(dirName+"/"+filename, "=====>", "/tmp/"+filename)
		} else {
			// note: change to remove later !
			err := os.Rename(dirName+"/"+filename, "/tmp/"+filename)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	fmt.Println("Success!")
}
