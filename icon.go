package main

import (
	"fmt"
	"sort"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

var (
	// iconNames is an array with the list of names of all the icons
	iconNames = extractIconNames()

	// iconsReverse Contains the key value pair where the key is the address of the icon and the value is the name
	iconReverse map[string]string

	// icons Has the hashmap of icons from the standard theme.
	// ToDo: Will have to look for a way to sync the list from `fyne_demo`
	icons map[string]fyne.Resource
)

func initIcons() {
	icons = map[string]fyne.Resource{
		"CancelIcon":        theme.CancelIcon(),
		"ConfirmIcon":       theme.ConfirmIcon(),
		"DeleteIcon":        theme.DeleteIcon(),
		"SearchIcon":        theme.SearchIcon(),
		"SearchReplaceIcon": theme.SearchReplaceIcon(),

		"CheckButtonIcon":        theme.CheckButtonIcon(),
		"CheckButtonCheckedIcon": theme.CheckButtonCheckedIcon(),
		"RadioButtonIcon":        theme.RadioButtonIcon(),
		"RadioButtonCheckedIcon": theme.RadioButtonCheckedIcon(),

		"ColorAchromaticIcon": theme.ColorAchromaticIcon(),
		"ColorChromaticIcon":  theme.ColorChromaticIcon(),
		"ColorPaletteIcon":    theme.ColorPaletteIcon(),

		"ContentAddIcon":    theme.ContentAddIcon(),
		"ContentRemoveIcon": theme.ContentRemoveIcon(),
		"ContentClearIcon":  theme.ContentClearIcon(),
		"ContentCutIcon":    theme.ContentCutIcon(),
		"ContentCopyIcon":   theme.ContentCopyIcon(),
		"ContentPasteIcon":  theme.ContentPasteIcon(),
		"ContentRedoIcon":   theme.ContentRedoIcon(),
		"ContentUndoIcon":   theme.ContentUndoIcon(),

		"InfoIcon":     theme.InfoIcon(),
		"ErrorIcon":    theme.ErrorIcon(),
		"QuestionIcon": theme.QuestionIcon(),
		"WarningIcon":  theme.WarningIcon(),

		"DocumentIcon":       theme.DocumentIcon(),
		"DocumentCreateIcon": theme.DocumentCreateIcon(),
		"DocumentPrintIcon":  theme.DocumentPrintIcon(),
		"DocumentSaveIcon":   theme.DocumentSaveIcon(),

		"FileIcon":            theme.FileIcon(),
		"FileApplicationIcon": theme.FileApplicationIcon(),
		"FileAudioIcon":       theme.FileAudioIcon(),
		"FileImageIcon":       theme.FileImageIcon(),
		"FileTextIcon":        theme.FileTextIcon(),
		"FileVideoIcon":       theme.FileVideoIcon(),
		"FolderIcon":          theme.FolderIcon(),
		"FolderNewIcon":       theme.FolderNewIcon(),
		"FolderOpenIcon":      theme.FolderOpenIcon(),
		"ComputerIcon":        theme.ComputerIcon(),
		"HomeIcon":            theme.HomeIcon(),
		"HelpIcon":            theme.HelpIcon(),
		"HistoryIcon":         theme.HistoryIcon(),
		"SettingsIcon":        theme.SettingsIcon(),
		"StorageIcon":         theme.StorageIcon(),
		"DownloadIcon":        theme.DownloadIcon(),
		// "UploadIcon":          theme.UploadIcon(),

		"ViewFullScreenIcon": theme.ViewFullScreenIcon(),
		"ViewRestoreIcon":    theme.ViewRestoreIcon(),
		"ViewRefreshIcon":    theme.ViewRefreshIcon(),
		"VisibilityIcon":     theme.VisibilityIcon(),
		"VisibilityOffIcon":  theme.VisibilityOffIcon(),
		"ZoomFitIcon":        theme.ZoomFitIcon(),
		"ZoomInIcon":         theme.ZoomInIcon(),
		"ZoomOutIcon":        theme.ZoomOutIcon(),

		"MoveDownIcon": theme.MoveDownIcon(),
		"MoveUpIcon":   theme.MoveUpIcon(),

		"NavigateBackIcon": theme.NavigateBackIcon(),
		"NavigateNextIcon": theme.NavigateNextIcon(),

		"MenuIcon":         theme.MenuIcon(),
		"MenuExpandIcon":   theme.MenuExpandIcon(),
		"MenuDropDownIcon": theme.MenuDropDownIcon(),
		"MenuDropUpIcon":   theme.MenuDropUpIcon(),

		"MailAttachmentIcon": theme.MailAttachmentIcon(),
		"MailComposeIcon":    theme.MailComposeIcon(),
		"MailForwardIcon":    theme.MailForwardIcon(),
		"MailReplyIcon":      theme.MailReplyIcon(),
		"MailReplyAllIcon":   theme.MailReplyAllIcon(),
		"MailSendIcon":       theme.MailSendIcon(),

		"MediaFastForwardIcon": theme.MediaFastForwardIcon(),
		"MediaFastRewindIcon":  theme.MediaFastRewindIcon(),
		"MediaPauseIcon":       theme.MediaPauseIcon(),
		"MediaPlayIcon":        theme.MediaPlayIcon(),
		// "MediaStopIcon":         theme.MediaStopIcon(),
		"MediaRecordIcon":       theme.MediaRecordIcon(),
		"MediaReplayIcon":       theme.MediaReplayIcon(),
		"MediaSkipNextIcon":     theme.MediaSkipNextIcon(),
		"MediaSkipPreviousIcon": theme.MediaSkipPreviousIcon(),

		"VolumeDownIcon": theme.VolumeDownIcon(),
		"VolumeMuteIcon": theme.VolumeMuteIcon(),
		"VolumeUpIcon":   theme.VolumeUpIcon(),
	}
	iconNames = extractIconNames()
	iconReverse = reverseIconMap()
}

// extractIconNames returns all the list of names of all the icons from the hashmap `icons`
func extractIconNames() []string {
	var iconNamesFromData = make([]string, len(icons))
	i := 0
	for k := range icons {
		iconNamesFromData[i] = k
		i++
	}

	sort.Strings(iconNamesFromData)
	return iconNamesFromData
}

// reverseIconMap returns all the list of icons and their addresses
func reverseIconMap() map[string]string {
	var iconReverseFromData = make(map[string]string, len(icons))
	for k, v := range icons {
		s := fmt.Sprintf("%p", v)
		iconReverseFromData[s] = k
	}
	return iconReverseFromData
}
