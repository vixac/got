package bullet_engine

import (
	"fmt"
	"strconv"

	"vixac.com/got/console"
	"vixac.com/got/engine"
)

func NewTable(items []engine.GotItemDisplay) console.ConsoleTable {
	var rows []console.TableRow

	titleCells := []console.TableCell{
		console.NewTableCellFromStr("#", console.TokenTextTertiary{}),
		console.NewTableCellFromStr("Path", console.TokenTextTertiary{}),
		//console.NewTableCellFromStr("Alias", console.TokenPrimary{}),
		console.NewTableCellFromStr(engine.CompleteChar, console.TokenComplete{}),
		console.NewTableCellFromStr(engine.NoteChar, console.TokenTextTertiary{}),
		console.NewTableCellFromStr(engine.ActiveChar, console.TokenPrimary{}),
		//console.NewTableCellFromStr("ID", console.TokenGid{}),
		console.NewTableCellFromStr("Title", console.TokenTextTertiary{}),
	}
	titleRow := console.NewCellTableRow(titleCells)
	rows = append(rows, console.NewDividerRow('-', console.TokenGid{}))
	rows = append(rows, titleRow)

	rows = append(rows, console.NewDividerRow('.', console.TokenGid{}))
	for _, item := range items {
		var cells []console.TableCell

		//number go
		numSnippets := []console.Snippet{
			console.NewSnippet("#"+strconv.Itoa(item.NumberGo), console.TokenTextTertiary{}),
		}
		cells = append(cells, console.NewTableCell(numSnippets))

		//path
		path := item.Path
		var pathSnippets []console.Snippet
		for i, node := range path.Ancestry {
			if i != 0 {
				pathSnippets = append(pathSnippets, console.NewSnippet("/", console.TokenTextTertiary{}))
			}
			if node.Alias != nil {
				pathSnippets = append(pathSnippets, console.NewSnippet(*node.Alias, console.TokenAlias{}))
			} else {
				if node.Id != "0" {
					pathSnippets = append(pathSnippets, console.NewSnippet("0"+node.Id, console.TokenGid{}))
				}
			}
		}
		//add the gid or alias of this item into it's own path.
		pathSnippets = append(pathSnippets, console.NewSnippet(":", console.TokenTextTertiary{}))
		pathSuffixShortcut, isAlias := item.Shortcut()
		var pathSuffixToken console.Token
		if isAlias {
			pathSuffixToken = console.TokenAlias{}
		} else {
			pathSuffixToken = console.TokenGid{}
		}
		endOfPathPadding := "  " //put some padding at the end of path to make summaries appear as one
		pathSnippets = append(pathSnippets, console.NewSnippet(pathSuffixShortcut+endOfPathPadding, pathSuffixToken))

		cells = append(cells, console.NewTableCell(pathSnippets))

		//alias
		//		cells = append(cells, console.NewTableCellFromStr(item.Alias, console.TokenPrimary{}))

		//summary
		if item.SummaryObj != nil && item.SummaryObj.Counts != nil {
			//complete
			cells = append(cells, console.NewTableCellFromStr(zeroIsEmpty(item.SummaryObj.Counts.Complete), console.TokenComplete{}))
			//note
			cells = append(cells, console.NewTableCellFromStr(zeroIsEmpty(item.SummaryObj.Counts.Notes), console.TokenTextTertiary{}))
			//active
			endOfSummaryPadding := "  " //we shove some padding on the final one to make these seem grouped
			cells = append(cells, console.NewTableCellFromStr(zeroIsEmpty(item.SummaryObj.Counts.Active)+endOfSummaryPadding, console.TokenPrimary{}))

			//cells = append(cells, console.NewTableCell(snippets))
		} else {
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
					//token = console.TokenPrimary{}
				} else if *state == engine.Note {
					noteCell = console.NewTableCellFromStr(state.ToStr(), console.TokenTextTertiary{})
					//token = console.TokenSecondary{}
				} else {

					completeCell = console.NewTableCellFromStr(state.ToStr(), console.TokenComplete{})
					//token = console.TokenComplete{}
				}
				//	snippet := []console.Snippet{console.NewSnippet(state.ToStr(), token)}
				cells = append(cells, completeCell)
				cells = append(cells, noteCell)
				cells = append(cells, activeCell)
			}
		}
		//gid cells = append(cells, console.NewTableCellFromStr(item.Gid, console.TokenGid{}))
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
