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
	FlatPaths          bool
	ShowCreatedColumn  bool
	ShowUpdatedColumn  bool
	SortByPath         bool
	GroupByTimeFrame   bool
	HideUnderCollapsed bool
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
	if len(pathSnippets) > 0 { // here is for the leaf node. We currently print it withthe same prefix as the rest of the path, but we could ":" if we wanted to.
		pathSnippets = append(pathSnippets, console.NewSnippet("/", console.TokenTextTertiary{}))
	}

	pathSuffixShortcut, isAlias := item.Shortcut()
	var pathSuffixToken console.Token
	if isAlias {
		pathSuffixToken = console.TokenAlias{}
	} else {
		pathSuffixToken = console.TokenGid{}
	}
	endOfPathPadding := "  " //put some padding at the end of path to make summaries appear as one
	pathSnippets = append(pathSnippets, console.NewSnippet(pathSuffixShortcut.Display+endOfPathPadding, pathSuffixToken))

	return console.NewTableCell(pathSnippets)
}

// all the cells that make up a got row
type GotRow struct {
	ItemNumber    console.TableCell
	Created       console.TableCell
	Updated       console.TableCell
	Path          console.TableCell
	GroupStart    console.TableCell
	CompleteCount console.TableCell
	ActiveCount   console.TableCell
	GroupEnd      console.TableCell
	Deadline      console.TableCell
	LongForm      console.TableCell
	Collapsed     console.TableCell
	State         console.TableCell
	Title         console.TableCell
}

func NewGotRow() GotRow {
	emptyCell := console.NewTableCellFromStr(" ", console.TokenPrimary{})
	return GotRow{
		ItemNumber:    emptyCell,
		Created:       emptyCell,
		Updated:       emptyCell,
		Path:          emptyCell,
		GroupStart:    emptyCell,
		CompleteCount: emptyCell,
		ActiveCount:   emptyCell,
		GroupEnd:      emptyCell,
		Deadline:      emptyCell,
		LongForm:      emptyCell,
		Collapsed:     emptyCell,
		State:         emptyCell,
		Title:         emptyCell,
	}
}

// Each section gets a divider between them
type GotTableSections struct {
	Sections [][]engine.GotItemDisplay
}

func (g GotRow) TableRow() console.TableRow {
	cells := []console.TableCell{
		g.ItemNumber, g.Created, g.Updated, g.Path, g.GroupStart, g.CompleteCount, g.ActiveCount, g.GroupEnd, g.Deadline, g.Collapsed, g.LongForm, g.State, g.Title,
	}
	return console.NewCellTableRow(cells)
}

func NewTable(sections *GotTableSections, options TableRenderOptions) (console.ConsoleTable, error) {

	if len(sections.Sections) == 0 {
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
	titleRow := NewGotRow()
	titleRow.Path = console.NewTableCellFromStr("Path", console.TokenTextTertiary{})
	titleRow.GroupStart = console.NewTableCellFromStr("[", console.TokenTextTertiary{})
	titleRow.GroupEnd = console.NewTableCellFromStr("] ", console.TokenTextTertiary{})
	titleRow.CompleteCount = console.NewTableCellFromStr(engine.CompleteChar+smallPadding, console.TokenComplete{})
	titleRow.ActiveCount = console.NewTableCellFromStr(engine.ActiveChar, console.TokenPrimary{})
	titleRow.Deadline = console.NewTableCellFromStr("Deadline ", console.TokenTextTertiary{})
	titleRow.Title = console.NewTableCellFromStr("Title ", console.TokenTextTertiary{})
	rows = append(rows, console.NewDividerRow("─", console.TokenTextTertiary{}))
	rows = append(rows, titleRow.TableRow())
	rows = append(rows, console.NewDividerRow("=", console.TokenTextTertiary{}))

	//unfortunately because of these 2 variables, the path rendering is contextual so we cant just do it line by line
	//var lastParentId *string = nil
	//var lastId *string = nil

	for i, section := range sections.Sections {

		if i != 0 {
			//section divider
			rows = append(rows, console.NewDividerRow("-", console.TokenTextTertiary{}))
		}
		for _, item := range section {
			itemRow := NewGotRow()
			numSnippets := []console.Snippet{
				console.NewSnippet("#"+strconv.Itoa(item.NumberGo)+mediumPadding, console.TokenNote{}),
			}
			itemRow.ItemNumber = console.NewTableCell(numSnippets)
			if options.ShowCreatedColumn {
				itemRow.Created = console.NewTableCellFromStr(item.Created+" ", console.TokenGroup{})
			}
			if options.ShowUpdatedColumn {
				itemRow.Updated = console.NewTableCellFromStr(" "+item.Updated+" ", console.TokenNote{})
			}

			if options.FlatPaths {
				pathCell := renderPathFlat(&item)
				itemRow.Path = pathCell

			} else {
				var pathSnippets []console.Snippet
				var parentIndex int = -1
				var ancestryPath []engine.PathItem // paths are optional so we support empty array
				if item.Path != nil {
					parentIndex = len(item.Path.Ancestry) - 1
					ancestryPath = item.Path.Ancestry
				}

				var treePattern = ""

				for i, node := range ancestryPath {
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

					//VX:TODO this logic adds a space line each time
					/*
						if i == parentIndex {
							thisParent := "0" + ancestryPath[parentIndex].Id

							//						isSibling := lastParentId != nil && *lastParentId == thisParent
							//				isFirstChild := lastId != nil && *lastId == thisParent

							//	if !isFirstChild && !isSibling {
							//		rows = append(rows, console.NewDividerRow(" ", console.TokenTextTertiary{}))
							//	}

							lastParentId = &thisParent
						}
					*/
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

				//lastId = &item.DisplayGidi
				pathSnippets = append(pathSnippets, console.NewSnippet(treePattern, console.TokenTextTertiary{}))

				pathSuffixShortcut, isAlias := item.Shortcut()
				var pathSuffixToken console.Token
				if isAlias {
					pathSuffixToken = console.TokenAlias{}
				} else {
					pathSuffixToken = console.TokenGid{}
				}
				pathSnippets = append(pathSnippets, console.NewSnippet(pathSuffixShortcut.Display+mediumPadding, pathSuffixToken))
				itemRow.Path = console.NewTableCell(pathSnippets)
			}
			//path

			if item.SummaryObj != nil && item.SummaryObj.Counts != nil {
				itemRow.GroupStart = console.NewTableCellFromStr("[", console.TokenTextTertiary{})
				itemRow.GroupEnd = console.NewTableCellFromStr("]", console.TokenTextTertiary{})
				itemRow.ActiveCount = console.NewTableCellFromStr(zeroIsEmpty(item.SummaryObj.Counts.Active), console.TokenPrimary{})
				itemRow.CompleteCount = console.NewTableCellFromStr(zeroIsEmpty(item.SummaryObj.Counts.Complete)+smallPadding, console.TokenComplete{})
			}
			itemRow.Deadline = console.NewTableCellFromStr(item.Deadline+" ", item.DeadlineToken)
			if item.IsCollapsed() {
				itemRow.Collapsed = console.NewTableCellFromStr(engine.CollapsedChar+" ", console.TokenGroup{})
			}
			if item.HasTNote {
				itemRow.LongForm = console.NewTableCellFromStr(engine.TNoteChar+" ", console.TokenGroup{})
			}

			itemRow.State = stateToCell(item.SummaryObj.State)
			tagStr := ""
			//tags
			if item.SummaryObj.Tags != nil && len(item.SummaryObj.Tags) == 0 {
				//VX:TODO invert the if
			} else {

				for i, t := range item.SummaryObj.Tags {
					if i == 0 {
						tagStr = " ("
					} else {
						tagStr += ","
					}
					tagStr += t.Literal.Display
				}
				if tagStr != "" {
					tagStr += ")"
				}
				tagStr += " "
			}
			tagSnippet := console.NewSnippet(tagStr, console.TokenAlert{})

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

			totalTitle := tagStr + item.Title
			var truncationIndex = -1
			for j := maxTitleLen; j < len(totalTitle); j++ {
				if totalTitle[j] == ' ' {
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

				titleSnippet := console.NewSnippet(titlePrefix+firstLine+dotDotDot, titleToken)

				itemRow.Title = console.NewTableCell([]console.Snippet{tagSnippet, titleSnippet})
				totalEmptyCells := len(titleCells) - 1
				var overflowRowCells []console.TableCell
				for i := 0; i < totalEmptyCells; i++ {
					overflowRowCells = append(overflowRowCells, emptyCell)
				}
				overflowRowCells = append(overflowRowCells, console.NewTableCellFromStr(titlePrefix+paddedSecondString, titleToken))

				overFlowRow := NewGotRow()
				overFlowRow.Title = console.NewTableCellFromStr(titlePrefix+paddedSecondString, titleToken)

				rows = append(rows, itemRow.TableRow())
				rows = append(rows, overFlowRow.TableRow())

			} else { //no truncation
				titleSnippet := console.NewSnippet(titlePrefix+item.Title, titleToken)
				itemRow.Title = console.NewTableCell([]console.Snippet{tagSnippet, titleSnippet})
				rows = append(rows, itemRow.TableRow())
			}
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
