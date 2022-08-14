package cli

import (
	"fmt"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/olekukonko/tablewriter"
	"github.com/skratchdot/open-golang/open"
	"github.com/urfave/cli/v2"
)

type App = cli.App
type Context = cli.Context

type Command = cli.Command

type Flag = cli.Flag
type IntFlag = cli.IntFlag
type IntSliceFlag = cli.IntSliceFlag
type StringFlag = cli.StringFlag
type StringSliceFlag = cli.StringSliceFlag
type BoolFlag = cli.BoolFlag
type PathFlag = cli.PathFlag

func Info(v ...interface{}) {
	os.Stdout.WriteString(fmt.Sprintln(v...))
}

func Infof(format string, a ...interface{}) {
	v := fmt.Sprintf(format, a...)
	Info(v)
}

func Warn(v ...interface{}) {
	os.Stderr.WriteString(fmt.Sprintln(v...))
}

func Warnf(format string, a ...interface{}) {
	v := fmt.Sprintf(format, a...)
	Warn(v)
}

func Error(v ...interface{}) {
	os.Stderr.WriteString(fmt.Sprintln(v...))
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

func OpenURL(url string) error {
	err := open.Run(url)

	if err != nil {
		Error("Unable to start your browser. try manually")
		Info(url)
	}

	return nil
}

func Select(label string, items []string) (int, string, error) {
	prompt := promptui.Select{
		Label: label,
		Items: items,
	}

	return prompt.Run()
}

func Prompt(label, placeholder string) (string, error) {
	prompt := promptui.Prompt{
		Label:   label,
		Default: placeholder,
	}

	return prompt.Run()
}

func Confirm(label string, placeholder bool) (bool, error) {
	value := "n"

	if placeholder {
		value = "y"
	}

	prompt := promptui.Prompt{
		Label: label,

		IsConfirm: true,
		Default:   value,
	}

	_, err := prompt.Run()

	if err != nil {
		return false, err
	}

	return true, nil
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
