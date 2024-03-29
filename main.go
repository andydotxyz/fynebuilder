package main

import (
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
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
				c.c.Objects = append(c.c.Objects, wrapContent(selected.create(), c.c))
				c.c.Refresh()
			}
			return
		}
		log.Println("Please select a container")
	}), nil, nil, list)
}

func buildUI(win fyne.Window) fyne.CanvasObject {
	content := previewUI().(*fyne.Container)
	overlay := wrapContent(content, nil)
	wrap := container.NewMax(overlay)

	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.FolderOpenIcon(), func() {
			d := dialog.NewFileOpen(func(r fyne.URIReadCloser, err error) {
				if err != nil {
					dialog.ShowError(err, win)
				}
				if r == nil {
					return
				}

				obj := DecodeJSON(r)
				_ = r.Close()

				overlay = wrapContent(obj, nil)
				wrap.Objects[0] = overlay
				wrap.Refresh()
			}, win)
			d.SetFilter(storage.NewExtensionFileFilter([]string{".json"}))
			d.Show()
		}),
		widget.NewToolbarAction(theme.DocumentSaveIcon(), func() {
			d := dialog.NewFileSave(func(w fyne.URIWriteCloser, err error) {
				if err != nil {
					dialog.ShowError(err, win)
				}
				if w == nil {
					return
				}

				err = EncodeJSON(overlay, w)
				if err != nil {
					dialog.ShowError(err, win)
				}
				_ = w.Close()
			}, win)
			d.SetFilter(storage.NewExtensionFileFilter([]string{".json"}))
			d.SetFileName("main.ui.json")
			d.Show()
		}),
		widget.NewToolbarAction(theme.DownloadIcon(), func() {
			packagesList := packagesRequired(overlay)
			code := exportCode(packagesList, overlay)
			fmt.Println(code)
		}),
		widget.NewToolbarAction(theme.MailForwardIcon(), func() {
			packagesList := append(packagesRequired(overlay), "app")
			code := exportCode(packagesList, overlay)
			code += `
func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Hello")
	myWindow.SetContent(makeUI())
	myWindow.ShowAndRun()
}
`
			path := filepath.Join(os.TempDir(), "fynebuilder")
			os.MkdirAll(path, 0711)
			path = filepath.Join(path, "main.go")
			_ = ioutil.WriteFile(path, []byte(code), 0600)

			cmd := exec.Command("go", "run", path)
			cmd.Stderr = os.Stderr
			cmd.Stdout = os.Stdout
			cmd.Start()
		}))

	widType = widget.NewLabelWithStyle("(None Selected)", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	paletteList = container.NewVBox()
	palette := container.NewBorder(widType, nil, nil, nil,
		container.NewGridWithRows(2, widget.NewCard("Properties", "", paletteList),
			widget.NewCard("Component List", "", buildLibrary()),
		))

	split := container.NewHSplit(wrap, palette)
	split.Offset = 0.8
	return container.New(layout.NewBorderLayout(toolbar, nil, nil, nil), toolbar,
		split)
}

func packagesRequired(obj fyne.CanvasObject) []string {
	if w, ok := obj.(*overlayWidget); ok {
		return w.Packages()
	}

	ret := []string{"container"}
	var objs []fyne.CanvasObject
	if c, ok := obj.(*fyne.Container); ok {
		objs = c.Objects
	} else if c, ok := obj.(*overlayContainer); ok {
		objs = c.c.Objects
	}
	for _, w := range objs {
		for _, p := range packagesRequired(w) {
			added := false
			for _, exists := range ret {
				if p == exists {
					added = true
					break
				}
			}
			if !added {
				ret = append(ret, p)
			}
		}
	}
	return ret
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
	remove := widget.NewButton("Remove", func() {
		var parent *fyne.Container
		var obj fyne.CanvasObject
		if c, ok := current.(*overlayContainer); ok {
			parent = c.parent
			obj = c
		} else if w, ok := current.(*overlayWidget); ok {
			parent = w.parent
			for _, o := range parent.Objects { // match our widget in the container wrapping us
				if c, ok := o.(*fyne.Container); ok && c.Objects[0] == w.child {
					obj = c
					break
				}
			}
		}
		if parent == nil {
			log.Println("Nothing to remove")
			return
		}

		parent.Remove(obj)
		parent.Refresh()
	})
	paletteList.Objects = []fyne.CanvasObject{editForm, remove}
	paletteList.Refresh()
}

func exportCode(pkgs []string, obj fyne.CanvasObject) string {
	for i := 0; i < len(pkgs); i++ {
		pkgs[i] = fmt.Sprintf(`	"fyne.io/fyne/v2/%s"`, pkgs[i])
	}
	code := fmt.Sprintf(`
package main

import (
	"fyne.io/fyne/v2"
%s
)

func makeUI() fyne.CanvasObject {
	return %#v
}
`,
		strings.Join(pkgs, "\n"),
		obj)

	formatted, err := format.Source([]byte(code))
	if err != nil {
		log.Fatal(err)
	}
	return string(formatted)
}

func main() {
	a := app.NewWithID("xyz.andy.fynebuilder")
	initIcons()
	initWidgets()

	w := a.NewWindow("Fyne Builder")
	w.SetContent(buildUI(w))
	w.Resize(fyne.NewSize(600, 400))
	w.ShowAndRun()
}

func previewUI() fyne.CanvasObject {
	return container.New(layout.NewVBoxLayout(),
		widget.NewIcon(theme.ContentAddIcon()),
		widget.NewLabel("label"),
		widget.NewButton("Button", func() {}))
}
