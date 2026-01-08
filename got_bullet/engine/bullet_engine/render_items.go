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
	State         console.TableCell
	Tags          console.TableCell
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
		State:         emptyCell,
		Tags:          emptyCell,
		Title:         emptyCell,
	}
}
func (g GotRow) TableRow() console.TableRow {
	cells := []console.TableCell{
		g.ItemNumber, g.Created, g.Updated, g.Path, g.GroupStart, g.CompleteCount, g.ActiveCount, g.GroupEnd, g.Deadline, g.LongForm, g.State, g.Tags, g.Title,
	}
	return console.NewCellTableRow(cells)
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
	titleRow := NewGotRow()
	titleRow.Path = console.NewTableCellFromStr("Path", console.TokenTextTertiary{})
	titleRow.GroupStart = console.NewTableCellFromStr("[", console.TokenTextTertiary{})
	titleRow.GroupEnd = console.NewTableCellFromStr("] ", console.TokenTextTertiary{})
	titleRow.CompleteCount = console.NewTableCellFromStr(engine.CompleteChar+smallPadding, console.TokenComplete{})
	titleRow.ActiveCount = console.NewTableCellFromStr(engine.ActiveChar, console.TokenPrimary{})
	titleRow.Deadline = console.NewTableCellFromStr("Deadline ", console.TokenTextTertiary{})
	titleRow.Tags = console.NewTableCellFromStr("Tags ", console.TokenPrimary{})
	titleRow.Title = console.NewTableCellFromStr("Title ", console.TokenTextTertiary{})
	titleRow.LongForm = console.NewTableCellFromStr(engine.TNoteChar, console.TokenTextTertiary{})
	titleRow.State = console.NewTableCellFromStr("~ ~", console.TokenTextTertiary{})
	rows = append(rows, console.NewDividerRow("─", console.TokenTextTertiary{}))
	rows = append(rows, titleRow.TableRow())
	rows = append(rows, console.NewDividerRow("=", console.TokenTextTertiary{}))

	if fetched.Parent != nil {
		parentRow := NewGotRow()
		if options.ShowCreatedColumn {
			parentRow.Created = console.NewTableCellFromStr(fetched.Parent.Created+" ", console.TokenGroup{})
		}
		if options.ShowUpdatedColumn {
			parentRow.Updated = console.NewTableCellFromStr(fetched.Parent.Created+" ", console.TokenGroup{})
		}
		pathCell := renderPathFlat(fetched.Parent)
		parentRow.Path = pathCell

		if fetched.Parent.SummaryObj != nil && fetched.Parent.SummaryObj.Counts != nil {
			parentRow.GroupStart = console.NewTableCellFromStr("[", console.TokenTextTertiary{})
			parentRow.CompleteCount = console.NewTableCellFromStr(zeroIsEmpty(fetched.Parent.SummaryObj.Counts.Complete)+smallPadding, console.TokenComplete{})
			parentRow.ActiveCount = console.NewTableCellFromStr(zeroIsEmpty(fetched.Parent.SummaryObj.Counts.Active), console.TokenPrimary{})
			parentRow.GroupEnd = console.NewTableCellFromStr("]", console.TokenTextTertiary{})
		}
		parentRow.Deadline = console.NewTableCellFromStr(fetched.Parent.Deadline+" ", fetched.Parent.DeadlineToken)

		if fetched.Parent.HasTNote {
			parentRow.LongForm = console.NewTableCellFromStr(engine.TNoteChar+" ", console.TokenGroup{})
		}
		parentRow.State = stateToCell(fetched.Parent.SummaryObj.State)
		//VX:TODO title token and truncations.
		parentRow.Title = console.NewTableCellFromStr(fetched.Parent.Title, console.TokenSecondary{})
		rows = append(rows, parentRow.TableRow())
		rows = append(rows, console.NewDividerRow("─", console.TokenTextTertiary{}))

	}

	//unfortunately because of these 2 variables, the path rendering is contextual so we cant just do it line by line
	var lastParentId *string = nil
	var lastId *string = nil

	for _, item := range fetched.Result {
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
		if item.HasTNote {
			itemRow.LongForm = console.NewTableCellFromStr(engine.TNoteChar+" ", console.TokenGroup{})
		}
		itemRow.State = stateToCell(item.SummaryObj.State)
		//tags
		if item.SummaryObj.Tags != nil && len(item.SummaryObj.Tags) == 0 {
			//VX:TODO invert the if
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
			itemRow.Tags = console.NewTableCellFromStr(tagStr, console.TokenAlert{})
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
			itemRow.Title = console.NewTableCellFromStr(titlePrefix+firstLine+dotDotDot, titleToken)
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
			itemRow.Title = console.NewTableCellFromStr(titlePrefix+item.Title, titleToken)
			rows = append(rows, itemRow.TableRow())
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
