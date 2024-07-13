package cli

import (
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/charmbracelet/lipgloss"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/browser"
	"github.com/urfave/cli/v3"
)

type Command = cli.Command

type Flag = cli.Flag
type IntFlag = cli.IntFlag
type IntSliceFlag = cli.IntSliceFlag
type StringFlag = cli.StringFlag
type StringSliceFlag = cli.StringSliceFlag
type BoolFlag = cli.BoolFlag

func Info(v ...interface{}) {
	os.Stdout.WriteString(fmt.Sprintln(v...))
}

func Infof(format string, a ...interface{}) {
	v := fmt.Sprintf(format, a...)
	Info(v)
}

func Warn(v ...interface{}) {
	color := lipgloss.AdaptiveColor{Light: "#FFFDF5", Dark: "#FFFDF5"}

	var style = lipgloss.NewStyle().
		Foreground(color)

	s := style.Render(fmt.Sprintln(v...))
	os.Stderr.WriteString(s + "\n")
}

func Warnf(format string, a ...interface{}) {
	v := fmt.Sprintf(format, a...)
	Warn(v)
}

func Error(v ...interface{}) {
	color := lipgloss.AdaptiveColor{Light: "#FF4672", Dark: "#ED567A"}

	var style = lipgloss.NewStyle().
		Foreground(color)

	s := style.Render(fmt.Sprintln(v...))
	os.Stderr.WriteString(s + "\n")
}

func Errorf(format string, a ...interface{}) {
	v := fmt.Sprintf(format, a...)
	Error(v)
}

func Fatal(v ...interface{}) {
	Error(v...)
	os.Exit(1)
}

func Fatalf(format string, a ...interface{}) {
	v := fmt.Sprintf(format, a...)
	Fatal(v)
}

func OpenFile(path string) error {
	err := browser.OpenFile(path)

	if err != nil {
		Error("Unable to start file. try manually")
		Error(path)
	}

	return nil
}

func OpenURL(url string) error {
	err := browser.OpenURL(url)

	if err != nil {
		Error("Unable to start your browser. try manually.")
		Error(url)
	}

	return nil
}

func Select(label string, items []string) (int, string, error) {
	s := huh.NewSelect[int]()

	if label != "" {
		s.Title(label)
	}

	options := make([]huh.Option[int], 0)

	for i, item := range items {
		options = append(options, huh.NewOption(item, i))
	}

	var index int

	s.Value(&index)
	s.Options(options...)

	if err := s.Run(); err != nil {
		return 0, "", err
	}

	result := items[index]

	if result != "" {
		fmt.Println("> " + result)
	}

	return index, result, nil
}

func MustSelect(label string, items []string) (int, string) {
	index, value, err := Select(label, items)

	if err != nil {
		Fatal(err)
	}

	return index, value
}

func Prompt(label, placeholder string) (string, error) {
	i := huh.NewInput()

	if label != "" {
		i.Title(label)
	}

	if placeholder != "" {
		i.Placeholder(placeholder)
	}

	var result string
	i.Value(&result)

	if err := i.Run(); err != nil {
		return "", err
	}

	if result != "" {
		fmt.Println("> " + result)
	}

	return result, nil
}

func MustPrompt(label, placeholder string) string {
	value, err := Prompt(label, placeholder)

	if err != nil {
		Fatal(err)
	}

	return value
}

func Confirm(label string, placeholder bool) (bool, error) {
	c := huh.NewConfirm()

	if label != "" {
		c.Title(label)
	}

	var result bool
	c.Value(&result)

	return result, c.Run()
}

func MustConfirm(label string, placeholder bool) bool {
	value, err := Confirm(label, placeholder)

	if err != nil {
		Fatal(err)
	}

	return value
}

func Title(val string) {
	// Green
	color := lipgloss.AdaptiveColor{Light: "#02BA84", Dark: "#02BF87"}

	var style = lipgloss.NewStyle().
		Foreground(color).
		Bold(true).
		Underline(true)

	fmt.Println(style.Render(val))
}

func Table(header []string, rows [][]string) {
	table := tablewriter.NewWriter(os.Stdout)

	table.SetHeader(header)
	table.AppendBulk(rows)

	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t") // pad with tabs
	table.SetNoWhiteSpace(true)

	table.Render()
}

func Run(title string, action func() error) error {
	var err error

	spinner.New().
		Title(title).
		Action(func() {
			err = action()
		}).
		Run()

	return err
}

func MustRun(title string, action func() error) error {
	err := Run(title, action)

	if err != nil {
		Fatal(err)
	}

	return err
}
