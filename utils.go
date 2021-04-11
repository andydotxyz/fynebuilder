package main

import "regexp"

func removeUnexportedProps(code string) string {
	areSimplePropsPresent := regexp.MustCompile(`[{ ]+[a-z][a-zA-Z]*:`)
	simpleProps1 := regexp.MustCompile(`([{ ]+)([a-z][a-zA-Z]*:[a-zA-Z0-9]*,)`)
	simpleProps2 := regexp.MustCompile(`([{ ]+)([a-z][a-zA-Z]*:[a-zA-Z0-9]*)([},]+)`)
	simpleProps3 := regexp.MustCompile(`([{ ]+)([a-z][a-zA-Z]*:[ ]*[a-z.A-Z]*{[a-zA-Z:0-9.,\- ]+}[,]*)`)
	simpleProps4 := regexp.MustCompile(`([{ ]+)([a-z][a-zA-Z]*:[ ]*[a-zA-Z0-9\.\*]*\([0-9a-zA-Z]+\)[,]*)`)
	simpleProps5 := regexp.MustCompile(`([{ ]+)([a-z][a-zA-Z]*:[ ]*\([a-zA-Z0-9\.\*]*\)\([0-9a-zA-Z]+\)[,]*)`)
	simpleProps6 := regexp.MustCompile(`([{ ]+)([a-z][a-zA-Z]*:[ ]*\[\][a-z]+\([a-zA-Z0-9, "]+\)[,]*)`)
	simpleProps7 := regexp.MustCompile(`([{ ]+)([a-z][a-zA-Z]*:[ ]*\(func\([a-zA-Z.]+\)\)\([0-9a-zA-Z]+\)[,]*)`)
	simpleProps8 := regexp.MustCompile(`([{ ]+)([a-z][a-zA-Z]*:[ ]*func\(\){fmt.Println\("Hello there"\)}[,]*)`)
	simpleProps9 := regexp.MustCompile(`([{ ]+)([a-z][a-zA-Z]*:[ ]*map\[[a-z]+\][a-z0-9]+[\({]+[a-zA-Z0-9:, "]+[\)}]+[,]*)`)
	simpleProps10 := regexp.MustCompile(`([{ ]+)(shortcut:[ ]*fyne.ShortcutHandler{.*?Copy.*?Cut.*?Paste.*?SelectAll.*?}}[,]*)`)

	for areSimplePropsPresent.MatchString(code) {
		code = simpleProps1.ReplaceAllString(code, "$1")
		code = simpleProps2.ReplaceAllString(code, "$1$3")
		code = simpleProps3.ReplaceAllString(code, "$1")
		code = simpleProps4.ReplaceAllString(code, "$1")
		code = simpleProps5.ReplaceAllString(code, "$1")
		code = simpleProps6.ReplaceAllString(code, "$1")
		code = simpleProps7.ReplaceAllString(code, "$1")
		code = simpleProps8.ReplaceAllString(code, "$1")
		code = simpleProps9.ReplaceAllString(code, "$1")
		code = simpleProps10.ReplaceAllString(code, "$1")

	}
	// fmt.Println(areSimplePropsPresent.MatchString(code))
	return code
}
