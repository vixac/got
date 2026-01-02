package bullet_engine

import (
	"strconv"

	"vixac.com/got/console"
	"vixac.com/got/engine"
)

const (
	separatorChar = "─"
)

type TableRenderOptions struct {
	FlatPaths         bool
	ShowCreatedColumn bool
	ShowUpdatedColumn bool
	SortByPath        bool
}

func renderPathFlat(item *engine.GotItemDisplay) console.TableCell {

	path := item.Path
	var pathSnippets []console.Snippet
	for i, node := range path.Ancestry {

		if i > 1 { //first node is 0 (i think, and we don't want to start with a slash on the second, so the first 2 items have no slash)
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
	if len(pathSnippets) > 0 { //we don't want the ":" at top level items because it screws alignment, but otherwise we use ":" instead of "/" to show that this is the id of this node.
		pathSnippets = append(pathSnippets, console.NewSnippet("/ ", console.TokenTextTertiary{}))
	}

	pathSuffixShortcut, isAlias := item.Shortcut()
	var pathSuffixToken console.Token
	if isAlias {
		pathSuffixToken = console.TokenAlias{}
	} else {
		pathSuffixToken = console.TokenGid{}
	}
	endOfPathPadding := "  " //put some padding at the end of path to make summaries appear as one
	pathSnippets = append(pathSnippets, console.NewSnippet(pathSuffixShortcut+endOfPathPadding, pathSuffixToken))

	return console.NewTableCell(pathSnippets)
}

func NewTable(items []engine.GotItemDisplay, options TableRenderOptions) (console.ConsoleTable, error) {

	if len(items) == 0 {
		return console.ConsoleTable{}, nil
	}

	var rows []console.TableRow

	mediumPadding := "  "
	smallPadding := " "
	emptyCell := console.NewTableCellFromStr("", console.TokenPrimary{})

	titleCells := []console.TableCell{}
	titleCells = append(titleCells, console.NewTableCellFromStr("#", console.TokenTextTertiary{}))
	if options.ShowCreatedColumn {
		titleCells = append(titleCells, console.NewTableCellFromStr("Created ", console.TokenTextTertiary{}))
	}
	if options.ShowUpdatedColumn {
		titleCells = append(titleCells, console.NewTableCellFromStr("Updated ", console.TokenTextTertiary{}))
	}

	titleCells = append(titleCells, console.NewTableCellFromStr("Path", console.TokenTextTertiary{}))
	titleCells = append(titleCells, console.NewTableCellFromStr("[", console.TokenTextTertiary{})) //"[" placeholder title
	titleCells = append(titleCells, console.NewTableCellFromStr(engine.CompleteChar+smallPadding, console.TokenComplete{}))
	titleCells = append(titleCells, console.NewTableCellFromStr(engine.NoteChar+smallPadding, console.TokenNote{}))
	titleCells = append(titleCells, console.NewTableCellFromStr(engine.ActiveChar, console.TokenPrimary{}))
	titleCells = append(titleCells, console.NewTableCellFromStr("]", console.TokenTextTertiary{})) //"]" placeholder title
	titleCells = append(titleCells, console.NewTableCellFromStr("Deadline ", console.TokenTextTertiary{}))

	titleCells = append(titleCells, console.NewTableCellFromStr("  ", console.TokenPrimary{})) //emptyCell, //leaf column has no title
	titleCells = append(titleCells, console.NewTableCellFromStr("Title", console.TokenTextTertiary{}))

	titleRow := console.NewCellTableRow(titleCells)
	rows = append(rows, console.NewDividerRow("─", console.TokenTextTertiary{}))
	rows = append(rows, titleRow)
	rows = append(rows, console.NewDividerRow("─", console.TokenTextTertiary{}))

	//unfortunately because of these 2 variables, the path rendering is contextual so we cant just do it line by line
	var lastParentId *string = nil
	var lastId *string = nil
	for _, item := range items {
		var cells []console.TableCell

		numSnippets := []console.Snippet{
			console.NewSnippet("#"+strconv.Itoa(item.NumberGo)+mediumPadding, console.TokenNote{}),
		}
		cells = append(cells, console.NewTableCell(numSnippets))
		if options.ShowCreatedColumn {
			cells = append(cells, console.NewTableCellFromStr(item.Created+" ", console.TokenGroup{}))
		}
		if options.ShowUpdatedColumn {
			cells = append(cells, console.NewTableCellFromStr(" "+item.Updated+" ", console.TokenNote{}))
		}

		if options.FlatPaths {
			pathCell := renderPathFlat(&item)
			cells = append(cells, pathCell)

		} else {
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

				if i == parentIndex {
					thisParent := "0" + item.Path.Ancestry[parentIndex].Id
					isSibling := lastParentId != nil && *lastParentId == thisParent
					isFirstChild := lastId != nil && *lastId == thisParent

					if !isFirstChild && !isSibling {
						rows = append(rows, console.NewDividerRow(" ", console.TokenTextTertiary{}))
					}

					lastParentId = &thisParent
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

			lastId = &item.Gid
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
		}
		//path

		if item.SummaryObj != nil && item.SummaryObj.Counts != nil {
			cells = append(cells, console.NewTableCellFromStr("[", console.TokenTextTertiary{}))
			cells = append(cells, console.NewTableCellFromStr(zeroIsEmpty(item.SummaryObj.Counts.Complete)+smallPadding, console.TokenComplete{}))
			cells = append(cells, console.NewTableCellFromStr(zeroIsEmpty(item.SummaryObj.Counts.Notes)+smallPadding, console.TokenNote{}))
			cells = append(cells, console.NewTableCellFromStr(zeroIsEmpty(item.SummaryObj.Counts.Active), console.TokenPrimary{}))
			cells = append(cells, console.NewTableCellFromStr("]"+mediumPadding, console.TokenTextTertiary{}))

		} else {
			cells = append(cells, emptyCell) //[ placeholder
			cells = append(cells, emptyCell) //complete placeholder
			cells = append(cells, emptyCell) //note placeholder
			cells = append(cells, emptyCell) //active placeholder
			cells = append(cells, emptyCell) //] plceholdere
		}

		cells = append(cells, console.NewTableCellFromStr(item.Deadline+" ", item.DeadlineToken))

		if item.HasTNote {
			cells = append(cells, console.NewTableCellFromStr(engine.TNoteChar, console.TokenGroup{}))
		} else {
			cells = append(cells, stateToCell(item.SummaryObj.State))
		}

		var titleToken console.Token
		var titlePrefix = ""
		if item.IsNote() {
			titleToken = console.TokenNote{}
		} else if item.SummaryObj.Counts != nil {
			titleToken = console.TokenGroup{}
		} else {
			titlePrefix = "  "
			titleToken = console.TokenSecondary{}
		}

		cells = append(cells, console.NewTableCellFromStr(titlePrefix+item.Title, titleToken))
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
		return console.TokenGid{}
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
