package main

import (
	"fmt"
	"go/format"
	"log"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var (
	editForm    *widget.Form
	widType     *widget.Label
	paletteList *fyne.Container
)

func buildLibrary() fyne.CanvasObject {
	var selected *widgetInfo
	tempNames := []string{}
	widgetLowerNames := []string{}
	for _, name := range widgetNames {
		widgetLowerNames = append(widgetLowerNames, strings.ToLower(name))
		tempNames = append(tempNames, name)
	}
	list := widget.NewList(func() int {
		return len(tempNames)
	}, func() fyne.CanvasObject {
		return widget.NewLabel("")
	}, func(i widget.ListItemID, obj fyne.CanvasObject) {
		obj.(*widget.Label).SetText(widgets[tempNames[i]].name)
	})
	list.OnSelected = func(i widget.ListItemID) {
		if match, ok := widgets[tempNames[i]]; ok {
			selected = &match
		}
	}
	list.OnUnselected = func(widget.ListItemID) {
		selected = nil
	}

	searchBox := widget.NewEntry()
	searchBox.SetPlaceHolder("Search Widgets")
	searchBox.OnChanged = func(s string) {
		s = strings.ToLower(s)
		tempNames = []string{}
		for i := 0; i < len(widgetLowerNames); i++ {
			if strings.Contains(widgetLowerNames[i], s) {
				tempNames = append(tempNames, widgetNames[i])
			}
		}
		list.Refresh()
		list.Select(0)   // Needed for new selection
		list.Unselect(0) // Without this (and with the above), list is behaving in a weird way
	}

	return container.NewBorder(searchBox, widget.NewButtonWithIcon("Insert", theme.ContentAddIcon(), func() {
		if c, ok := current.(*overlayContainer); ok {
			if selected != nil {
				c.c.Objects = append(c.c.Objects, wrapContent(selected.create()))
				c.c.Refresh()
			}
			return
		}
		log.Println("Please select a container")
	}), nil, nil, list)
}

func buildUI() fyne.CanvasObject {
	content := previewUI().(*fyne.Container)
	overlay := wrapContent(content)

	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.FolderOpenIcon(), func() {
			log.Println("TODO")
		}),
		widget.NewToolbarAction(theme.DocumentSaveIcon(), func() {
			log.Println("TODO")
		}),
		widget.NewToolbarAction(theme.MailForwardIcon(), func() {
			code := fmt.Sprintf("%#v", overlay)

			layoutReplace := regexp.MustCompile(`(layout.[a-zA-Z]+)`)
			code = layoutReplace.ReplaceAllString(code, "${1}Layout") // ToDo: should remove this line once the right layout is picked

			funcReplace := regexp.MustCompile(`\(func\(\)\)\([a-zA-Z0-9]*\)`)
			code = funcReplace.ReplaceAllString(code, "func(){fmt.Println(\"Hello there\")}")

			code = removeUnexportedProps(code)

			iconReplace := regexp.MustCompile(`\(\*theme.ThemedResource\)\((0x[a-zA-Z0-9]+)\)`)
			// fmt.Println("=================")

			cancelAddrFromTheme, _ := strconv.ParseUint(fmt.Sprintf("%p", theme.DefaultTheme().Icon("cancel")), 0, 64)

			for _, v := range iconReplace.FindAllStringSubmatch(code, -1) {
				iconAddr, _ := strconv.ParseUint(v[1], 0, 64)
				strAddr := "0x" + strconv.FormatUint(cancelAddr+iconAddr-cancelAddrFromTheme, 16)
				// fmt.Println(v, strings.Index(code, v[0]), iconAddr, strAddr, iconReverse[strAddr])
				code = strings.Replace(code, v[0], "theme."+iconReverse[strAddr]+"()", 1)
			}

			// baseWidgetRegex := regexp.MustCompile(`BaseWidget:widget.BaseWidget{size:fyne.Size{Width:[0-9]+, Height:[0-9]+}, position:fyne.Position{X:[0-9]+, Y:[0-9]+}, Hidden:false, impl:\(\*[a-zA-Z]+\.[a-zA-Z]+\)\([0-9a-zA-Z]+\), propertyLock:sync.RWMutex{w:sync.Mutex{state:[0-9]+, sema:[0-9a-zA-Z]+}, writerSem:[0-9a-zA-Z]+, readerSem:[0-9a-zA-Z]+, readerCount:[0-9]+, readerWait:[0-9]+}}, `)
			// code = baseWidgetRegex.ReplaceAllString(code, "")
			packagesList := []string{"", "/app", "/canvas", "/container", "/data/binding", "/layout", "/theme", "/widget"} //ToDo: Will fetch it dynamically later
			for i := 0; i < len(packagesList); i++ {
				packagesList[i] = fmt.Sprintf(`"fyne.io/fyne/v2%s"`, packagesList[i])
			}
			code = fmt.Sprintf(`
package main
import (
	%s
)

func makeUI() *fyne.Container {
	return %s
}
func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Hello")
	myWindow.SetContent(makeUI())

	myWindow.ShowAndRun()
}
			`,
				strings.Join(packagesList, "\n"),
				code)
			formatted, err := format.Source([]byte(code))
			if err != nil {
				log.Fatal(err)
			}
			code = string(formatted)
			fmt.Println(code)
		}))

	widType = widget.NewLabelWithStyle("(None Selected)", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	paletteList = container.NewVBox()
	palette := container.NewBorder(widType, nil, nil, nil,
		container.NewGridWithRows(2, widget.NewCard("Properties", "", paletteList),
			widget.NewCard("Component List", "", buildLibrary()),
		))

	split := container.NewHSplit(overlay, palette)
	split.Offset = 0.8
	return container.New(layout.NewBorderLayout(toolbar, nil, nil, nil), toolbar,
		split)
}

func choose(o fyne.CanvasObject) {
	typeName := reflect.TypeOf(o).Elem().Name()
	widName := reflect.TypeOf(o).String()
	l := reflect.ValueOf(o).Elem()
	if typeName == "Entry" {
		if l.FieldByName("Password").Bool() {
			typeName = "PasswordEntry"
		} else if l.FieldByName("MultiLine").Bool() {
			typeName = "MultiLineEntry"
		}
		widName = "*widget." + typeName
	}
	widType.SetText(typeName)

	var items []*widget.FormItem
	if match, ok := widgets[widName]; ok {
		items = match.edit(o)
	}

	editForm = widget.NewForm(items...)
	paletteList.Objects = []fyne.CanvasObject{editForm}
	paletteList.Refresh()
}

func main() {
	a := app.NewWithID("xyz.andy.fynebuilder")
	w := a.NewWindow("Fyne Builder")
	w.SetContent(buildUI())
	w.Resize(fyne.NewSize(600, 400))
	w.ShowAndRun()
}

func previewUI() fyne.CanvasObject {
	return container.New(layout.NewVBoxLayout(),
		widget.NewIcon(theme.ContentAddIcon()),
		widget.NewLabel("label"),
		widget.NewButton("Button", func() {}))
}
