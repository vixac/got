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
	emptyCell := console.NewTableCellFromStr("", console.TokenPrimary{})
	if item == nil {
		return emptyCell
	}
	if item.Path == nil {
		return emptyCell
	}

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

func NewTable(fetched *engine.GotFetchResult, options TableRenderOptions) (console.ConsoleTable, error) {

	if len(fetched.Result) == 0 {
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

	titleCells = append(titleCells, console.NewTableCellFromStr("  ", console.TokenPrimary{}))    //emptyCell, //leaf column has no title
	titleCells = append(titleCells, console.NewTableCellFromStr("Tags ", console.TokenPrimary{})) //emptyCell, //leaf column has no title
	titleCells = append(titleCells, console.NewTableCellFromStr("Title", console.TokenTextTertiary{}))

	titleRow := console.NewCellTableRow(titleCells)
	rows = append(rows, console.NewDividerRow("─", console.TokenTextTertiary{}))
	rows = append(rows, titleRow)
	rows = append(rows, console.NewDividerRow("=", console.TokenTextTertiary{}))

	//VX:TODO squash all this
	if fetched.Parent != nil {

		parentCells := []console.TableCell{}
		parentCells = append(parentCells, emptyCell) //#
		if options.ShowCreatedColumn {
			parentCells = append(parentCells, console.NewTableCellFromStr(fetched.Parent.Created+" ", console.TokenGroup{}))
		}
		if options.ShowUpdatedColumn {
			parentCells = append(parentCells, console.NewTableCellFromStr(fetched.Parent.Updated+" ", console.TokenGroup{}))
		}
		pathCell := renderPathFlat(fetched.Parent)
		parentCells = append(parentCells, pathCell)

		if fetched.Parent.SummaryObj != nil && fetched.Parent.SummaryObj.Counts != nil {
			parentCells = append(parentCells, console.NewTableCellFromStr("[", console.TokenTextTertiary{}))
			parentCells = append(parentCells, console.NewTableCellFromStr(zeroIsEmpty(fetched.Parent.SummaryObj.Counts.Complete)+smallPadding, console.TokenComplete{}))
			parentCells = append(parentCells, console.NewTableCellFromStr(zeroIsEmpty(fetched.Parent.SummaryObj.Counts.Notes)+smallPadding, console.TokenNote{}))
			parentCells = append(parentCells, console.NewTableCellFromStr(zeroIsEmpty(fetched.Parent.SummaryObj.Counts.Active), console.TokenPrimary{}))
			parentCells = append(parentCells, console.NewTableCellFromStr("]"+mediumPadding, console.TokenTextTertiary{}))

		} else {
			parentCells = append(parentCells, emptyCell) //[ placeholder
			parentCells = append(parentCells, emptyCell) //complete placeholder
			parentCells = append(parentCells, emptyCell) //note placeholder
			parentCells = append(parentCells, emptyCell) //active placeholder
			parentCells = append(parentCells, emptyCell) //] plceholdere
		}
		parentCells = append(parentCells, console.NewTableCellFromStr(fetched.Parent.Deadline+" ", fetched.Parent.DeadlineToken))

		if fetched.Parent.HasTNote {
			parentCells = append(parentCells, console.NewTableCellFromStr(engine.TNoteChar, console.TokenGroup{}))
		} else {
			parentCells = append(parentCells, stateToCell(fetched.Parent.SummaryObj.State))
		}
		parentCells = append(parentCells, emptyCell) //VX:TODO tag for parent
		//VX:TODO title token and truncations.
		parentCells = append(parentCells, console.NewTableCellFromStr(fetched.Parent.Title, console.TokenSecondary{}))
		rows = append(rows, console.NewCellTableRow(parentCells))
		rows = append(rows, console.NewDividerRow("─", console.TokenTextTertiary{}))

	}

	//unfortunately because of these 2 variables, the path rendering is contextual so we cant just do it line by line
	var lastParentId *string = nil
	var lastId *string = nil
	for _, item := range fetched.Result {
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

		//tags
		if item.SummaryObj.Tags != nil && len(item.SummaryObj.Tags) == 0 {
			cells = append(cells, emptyCell)
		} else {

			tagStr := ""
			for i, t := range item.SummaryObj.Tags {
				if i == 0 {
					tagStr = "("
				} else {
					tagStr += ","
				}
				tagStr += t.Literal.Display
			}
			if tagStr != "" {
				tagStr += ")"
			}

			cells = append(cells, console.NewTableCellFromStr(tagStr, console.TokenAlert{}))

		}

		//title
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
		maxTitleLen := 100

		//check if we need to truncate
		var truncationIndex = -1
		for j := maxTitleLen; j < len(item.Title); j++ {
			if item.Title[j] == ' ' {
				truncationIndex = j
				break
			}
		}
		//we need to truncate, so we append ... on first line, and then prefix ... on second, and right pad the second line.
		if truncationIndex > 0 {
			var firstLine = ""
			var secondLine = "" //only wrapping to 2 lines, not recursively. Because yagni
			for i := 0; i < len(item.Title); i++ {
				if i < truncationIndex {
					firstLine += string(item.Title[i])
				} else {
					secondLine += string(item.Title[i])
				}
			}
			dotDotDot := " ..."
			paddingRequired := len(firstLine) - len(secondLine)
			paddedSecondString := secondLine
			if paddingRequired > 0 {
				paddedSecondString = ""
				for i := 0; i < paddingRequired; i++ {
					paddedSecondString += " "
				}
				paddedSecondString += dotDotDot + secondLine
			}
			cells = append(cells, console.NewTableCellFromStr(titlePrefix+firstLine+dotDotDot, titleToken))
			totalEmptyCells := len(titleCells) - 1
			var overflowRowCells []console.TableCell
			for i := 0; i < totalEmptyCells; i++ {
				overflowRowCells = append(overflowRowCells, emptyCell)
			}
			overflowRowCells = append(overflowRowCells, console.NewTableCellFromStr(titlePrefix+paddedSecondString, titleToken))

			rows = append(rows, console.NewCellTableRow(cells))
			rows = append(rows, console.NewCellTableRow(overflowRowCells))

		} else { //no truncation
			cells = append(cells, console.NewTableCellFromStr(titlePrefix+item.Title, titleToken))
			rows = append(rows, console.NewCellTableRow(cells))
		}

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
