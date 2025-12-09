package bullet_engine

import (
	"strconv"

	"vixac.com/got/console"
	"vixac.com/got/engine"
)

const (
	separatorChar = "─"
)

func NewTable(items []engine.GotItemDisplay) (console.ConsoleTable, error) {

	if len(items) == 0 {
		return console.ConsoleTable{}, nil
	}

	var rows []console.TableRow

	mediumPadding := "  "
	smallPadding := " "
	emptyCell := console.NewTableCellFromStr("", console.TokenPrimary{})

	titleCells := []console.TableCell{
		console.NewTableCellFromStr("#", console.TokenTextTertiary{}),
		console.NewTableCellFromStr("Path", console.TokenTextTertiary{}),
		console.NewTableCellFromStr("[", console.TokenTextTertiary{}), //"[" placeholder title
		console.NewTableCellFromStr(engine.CompleteChar+smallPadding, console.TokenComplete{}),
		console.NewTableCellFromStr(engine.NoteChar+smallPadding, console.TokenTextTertiary{}),
		console.NewTableCellFromStr(engine.ActiveChar, console.TokenPrimary{}),
		console.NewTableCellFromStr("]", console.TokenTextTertiary{}), //"]" placeholder title
		console.NewTableCellFromStr("  ", console.TokenPrimary{}),     //emptyCell, //leaf column has no title
		console.NewTableCellFromStr("Title", console.TokenTextTertiary{}),
	}

	titleRow := console.NewCellTableRow(titleCells)
	rows = append(rows, console.NewDividerRow("─", console.TokenTextTertiary{}))
	rows = append(rows, titleRow)
	rows = append(rows, console.NewDividerRow("─", console.TokenTextTertiary{}))

	for _, item := range items {
		var cells []console.TableCell

		numSnippets := []console.Snippet{
			console.NewSnippet("#"+strconv.Itoa(item.NumberGo)+mediumPadding, console.TokenTextTertiary{}),
		}
		cells = append(cells, console.NewTableCell(numSnippets))

		//path
		var pathSnippets []console.Snippet
		parentIndex := len(item.Path.Ancestry) - 1

		var treePattern = ""
		for i, node := range item.Path.Ancestry {
			var wordLength = 0
			if node.Alias != nil {
				wordLength = len(*node.Alias)
			} else if node.Id != "0" {
				wordLength = len(node.Id)
			}
			if i != parentIndex {
				treePattern += console.FitString("", wordLength, " ")
				continue
			}
			if wordLength == 0 {
				continue
			}
			switch wordLength {
			case 1:
				treePattern += "└"
			case 2:
				treePattern += "└ "
			case 3:
				treePattern += "└ "
			default:
				halfWordLength := wordLength / 2
				treePattern += console.FitString("", halfWordLength-1, " ")
				treePattern += "└"
				treePattern += console.FitString("", halfWordLength-1, separatorChar)
				treePattern += " "

			}
		}
		pathSnippets = append(pathSnippets, console.NewSnippet(treePattern, console.TokenTextTertiary{}))

		pathSuffixShortcut, isAlias := item.Shortcut()
		var pathSuffixToken console.Token
		if isAlias {
			pathSuffixToken = console.TokenAlias{}
		} else {
			pathSuffixToken = console.TokenGid{}
		}
		pathSnippets = append(pathSnippets, console.NewSnippet(pathSuffixShortcut+mediumPadding, pathSuffixToken))

		cells = append(cells, console.NewTableCell(pathSnippets))

		if item.SummaryObj != nil && item.SummaryObj.Counts != nil {
			cells = append(cells, console.NewTableCellFromStr("[", console.TokenTextTertiary{}))
			cells = append(cells, console.NewTableCellFromStr(zeroIsEmpty(item.SummaryObj.Counts.Complete)+smallPadding, console.TokenComplete{}))
			cells = append(cells, console.NewTableCellFromStr(zeroIsEmpty(item.SummaryObj.Counts.Notes)+smallPadding, console.TokenTextTertiary{}))
			cells = append(cells, console.NewTableCellFromStr(zeroIsEmpty(item.SummaryObj.Counts.Active), console.TokenPrimary{}))
			cells = append(cells, console.NewTableCellFromStr("]"+mediumPadding, console.TokenTextTertiary{}))

		} else {
			cells = append(cells, emptyCell) //[ placeholder
			cells = append(cells, emptyCell) //complete placeholder
			cells = append(cells, emptyCell) //note placeholder
			cells = append(cells, emptyCell) //active placeholder
			cells = append(cells, emptyCell) //] plceholdere
		}
		cells = append(cells, stateToCell(item.SummaryObj.State))
		var titleToken console.Token
		if item.IsNote() {
			titleToken = console.TokenTextTertiary{}
		} else {
			titleToken = console.TokenSecondary{}
		}
		cells = append(cells, console.NewTableCellFromStr(item.Title, titleToken))
		rows = append(rows, console.NewCellTableRow(cells))
	}
	return console.NewConsoleTable(rows)
}

func stateToToken(state *engine.GotState) console.Token {
	if state == nil {
		return console.TokenPrimary{}
	}
	switch *state {
	case engine.Active:
		return console.TokenPrimary{}
	case engine.Note:
		return console.TokenTextTertiary{}
	case engine.Complete:
		return console.TokenComplete{}
	}
	return console.TokenPrimary{}
}
func stateToStr(state *engine.GotState) string {
	if state == nil {
		return ""
	}
	return state.ToStr()
}

func stateToCell(state *engine.GotState) console.TableCell {
	return console.NewTableCellFromStr(stateToStr(state), stateToToken(state))
}

func zeroIsEmpty(input int) string {
	if input == 0 {
		return ""
	}
	return strconv.Itoa(input)
}
