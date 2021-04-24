package main

import (
	"sort"
	"strings"
)

func encodeDoubleQuote(inStr string) (outStr string) {
	outStr = strings.ReplaceAll(inStr, "\"", "\\\"")
	return
}

var packageList = map[string]bool{"widget": true, "theme": true}

var packageNameMap = map[string]string{
	"fmt":        "fmt",
	"fyne":       "fyne.io/fyne/v2",
	"app":        "fyne.io/fyne/v2/app",
	"canvas":     "fyne.io/fyne/v2/canvas",
	"container":  "fyne.io/fyne/v2/container",
	"binding":    "fyne.io/fyne/v2/data/binding",
	"validation": "fyne.io/fyne/v2/data/validation",
	"dialog":     "fyne.io/fyne/v2/dialog",
	"desktop":    "fyne.io/fyne/v2/driver/desktop",
	"mobile":     "fyne.io/fyne/v2/driver/mobile",
	"layout":     "fyne.io/fyne/v2/layout",
	"storage":    "fyne.io/fyne/v2/storage",
	"test":       "fyne.io/fyne/v2/driver/test",
	"theme":      "fyne.io/fyne/v2/theme",
	"widget":     "fyne.io/fyne/v2/widget",
}

func getPackages() (finalPackageList []string) {
	for k, _ := range packageList {
		finalPackageList = append(finalPackageList, packageNameMap[k])
	}
	sort.Slice(finalPackageList, func(i, j int) bool {
		return finalPackageList[i] < finalPackageList[j]
	})
	return
}
