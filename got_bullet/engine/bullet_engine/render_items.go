package bullet_engine

import (
	"fmt"
	"strconv"

	"vixac.com/got/console"
	"vixac.com/got/engine"
)

const (
	separatorChar = "─"
)

func NewTable(items []engine.GotItemDisplay) console.ConsoleTable {
	var rows []console.TableRow

	mediumPadding := "  "
	smallPadding := " "

	titleCells := []console.TableCell{
		console.NewTableCellFromStr("#", console.TokenTextTertiary{}),
		console.NewTableCellFromStr("Path", console.TokenTextTertiary{}),
		console.NewTableCellFromStr("[", console.TokenTextTertiary{}), //"[" placeholder title
		console.NewTableCellFromStr(engine.CompleteChar, console.TokenComplete{}),
		console.NewTableCellFromStr(engine.NoteChar, console.TokenTextTertiary{}),
		console.NewTableCellFromStr(engine.ActiveChar, console.TokenPrimary{}),
		console.NewTableCellFromStr("]", console.TokenTextTertiary{}), //"]" placeholder title
		console.NewTableCellFromStr("Title", console.TokenTextTertiary{}),
	}
	titleRow := console.NewCellTableRow(titleCells)
	rows = append(rows, console.NewDividerRow("─", console.TokenTextTertiary{}))
	rows = append(rows, titleRow)

	rows = append(rows, console.NewDividerRow("─", console.TokenTextTertiary{}))
	for _, item := range items {
		var cells []console.TableCell

		//number go
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

			if wordLength == 1 {
				treePattern += "└"
			} else if wordLength == 2 {
				treePattern += "└ "
			} else if wordLength == 3 {
				treePattern += " └" + separatorChar
			} else {
				//3 parts
				halfWordLength := wordLength / 2
				treePattern += console.FitString("", halfWordLength-1, " ")
				treePattern += "└"
				treePattern += console.FitString("", halfWordLength-1, separatorChar)
				treePattern += " "
			}
		}
		//add the gid or alias of this item into it's own path.
		if len(pathSnippets) > 0 { //we don't want the ":" at top level items because it screws alignment, but otherwise we use ":" instead of "/" to show that this is the id of this node.
			pathSnippets = append(pathSnippets, console.NewSnippet(":", console.TokenTextTertiary{}))
		} else {
			pathSnippets = append(pathSnippets, console.NewSnippet(treePattern, console.TokenTextTertiary{}))
		}

		pathSuffixShortcut, isAlias := item.Shortcut()
		var pathSuffixToken console.Token
		if isAlias {
			pathSuffixToken = console.TokenAlias{}
		} else {
			pathSuffixToken = console.TokenGid{}
		}
		pathSnippets = append(pathSnippets, console.NewSnippet(pathSuffixShortcut+mediumPadding, pathSuffixToken))

		cells = append(cells, console.NewTableCell(pathSnippets))

		//summary
		if item.SummaryObj != nil && item.SummaryObj.Counts != nil {
			cells = append(cells, console.NewTableCellFromStr("[", console.TokenTextTertiary{}))
			//complete
			cells = append(cells, console.NewTableCellFromStr(zeroIsEmpty(item.SummaryObj.Counts.Complete)+smallPadding, console.TokenComplete{}))
			//note
			cells = append(cells, console.NewTableCellFromStr(zeroIsEmpty(item.SummaryObj.Counts.Notes)+smallPadding, console.TokenTextTertiary{}))

			//active
			cells = append(cells, console.NewTableCellFromStr(zeroIsEmpty(item.SummaryObj.Counts.Active), console.TokenPrimary{}))

			cells = append(cells, console.NewTableCellFromStr("]"+mediumPadding, console.TokenTextTertiary{}))

		} else {

			cells = append(cells, console.NewTableCellFromStr("", console.TokenSecondary{})) //"[" placeholder
			//complete
			state := item.SummaryObj.State
			if state == nil {
				fmt.Printf("VX: ERRORR should not happen. Either a count or a state.")
				snippet := []console.Snippet{console.NewSnippet("<VX:err>", console.TokenBrand{})}
				cells = append(cells, console.NewTableCell(snippet))
			} else {
				emptyCell := console.NewTableCellFromStr("", console.TokenPrimary{})

				noteCell := emptyCell
				activeCell := emptyCell
				completeCell := emptyCell

				if *state == engine.Active {
					activeCell = console.NewTableCellFromStr(state.ToStr(), console.TokenPrimary{})
				} else if *state == engine.Note {
					noteCell = console.NewTableCellFromStr(state.ToStr(), console.TokenTextTertiary{})
				} else {

					completeCell = console.NewTableCellFromStr(state.ToStr(), console.TokenComplete{})
				}
				cells = append(cells, completeCell)
				cells = append(cells, noteCell)
				cells = append(cells, activeCell)
				cells = append(cells, console.NewTableCellFromStr("", console.TokenSecondary{})) //"]" placeholder
			}
		}
		cells = append(cells, console.NewTableCellFromStr(item.Title, console.TokenSecondary{}))
		rows = append(rows, console.NewCellTableRow(cells))

	}
	table := console.NewConsoleTable(rows)
	return table
}

func zeroIsEmpty(input int) string {
	if input == 0 {
		return ""
	}
	return strconv.Itoa(input)
}
