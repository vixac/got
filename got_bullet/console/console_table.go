package console

import (
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/rivo/uniseg"
)

type ConsoleTable struct {
	Rows          []TableRow
	ColumnWidths  []int
	ColumnPadding string
}

const (
	paddingSize = 0
)

func nchars(s string, targetWidth int) string {
	if s == "" || targetWidth <= 0 {
		return ""
	}

	r, _ := utf8.DecodeRuneInString(s)
	unit := string(r)
	unitWidth := uniseg.StringWidth(unit)

	var b strings.Builder
	current := 0

	for current < targetWidth {
		b.WriteString(unit)
		current += unitWidth
	}

	return b.String()
}
func (c *ConsoleTable) Render(printer Messenger, scheme Theme) {

	//the int is the index, and the message group is the rendered row.
	contentRows := make(map[int]MessageGroup)
	for rowNumber, row := range c.Rows {

		if row.CellRow != nil {

			var messages []Message

			for i, cell := range row.CellRow.Cells {

				if i != 0 {
					messages = append(messages, Message{Message: c.ColumnPadding})
				}

				for _, snippet := range cell.Content {
					message := Message{
						Color:   scheme.ColorFor(snippet.Token).Col(),
						Message: snippet.Text,
					}
					messages = append(messages, message)

				}
				paddingRequired := c.ColumnWidths[i] - cell.Length
				spaceCharacter := " " // //useful for debugging to make this and "X" or something.
				paddingStr := FitString("", paddingRequired, spaceCharacter)
				messages = append(messages, Message{Message: paddingStr})

			}
			var total = 0
			for _, m := range messages {
				total += uniseg.StringWidth(m.Message)
			}
			contentRows[rowNumber] = NewMessageGroup(messages)
		}
	}
	if len(contentRows) == 0 {
		//nothing to render besides maybe dividers. We should render nothin.
		return
	}

	var renderedRowLength = 0
	// we just look at one of these rows to get its length
	for _, val := range contentRows {
		renderedRowLength = val.TextLen
		break
	}

	//now we actually print
	for i, row := range c.Rows {
		if row.DividerRow != nil {
			dividerStr := nchars(row.DividerRow.Separator, renderedRowLength)
			dividerMessage := Message{
				Message: dividerStr,
				Color:   scheme.ColorFor(row.DividerRow.Token).Col(),
			}
			printer.Print(dividerMessage)
		} else {
			group := contentRows[i] //if this doesn't exist its a dev error, as we just populated this.
			printer.PrintInLine(group.Messages)
		}

	}
}

func NewConsoleTable(rows []TableRow) (ConsoleTable, error) {
	if len(rows) == 0 {
		return ConsoleTable{}, nil
	}
	var colCount = -1 //unset
	for _, r := range rows {
		if r.CellRow == nil {
			continue
		}
		if colCount == -1 {
			colCount = r.CellRow.NumCells //all future cells must be the same size as the first.
		} else {
			if r.CellRow.NumCells != colCount {
				fmt.Printf("Error, row is %d, instead of %dcells. \n", r.CellRow.NumCells, colCount)
				return ConsoleTable{}, errors.New("invalid cell count at row ")
			}
		}
	}

	var maxWidths []int
	for _, r := range rows {
		if r.CellRow == nil {
			continue
		}
		numCells := len(r.CellRow.Cells)
		for col := 0; col < numCells; col++ {
			minColWidthThisCell := r.CellRow.Cells[col].Length
			if len(maxWidths) <= col {
				maxWidths = append(maxWidths, 0)
			}
			if minColWidthThisCell > maxWidths[col] {
				maxWidths[col] = minColWidthThisCell
			}
		}
	}

	var padding = ""

	for range paddingSize {
		padding += " "
	}
	return ConsoleTable{
		Rows:          rows,
		ColumnWidths:  maxWidths,
		ColumnPadding: padding,
	}, nil
}

// this becomes a bit of an enum eventually. cells or divider
type CellRow struct {
	Cells     []TableCell
	NumCells  int
	RowLength int
}
type DividerRow struct {
	Separator string
	Token     Token
}
type TableRow struct {
	CellRow    *CellRow
	DividerRow *DividerRow
}

func NewDividerRow(separator string, token Token) TableRow {
	div := DividerRow{
		Separator: separator,
		Token:     token,
	}
	return TableRow{
		DividerRow: &div,
	}
}

func NewCellTableRow(cells []TableCell) TableRow {
	rowLength := 0
	for _, cell := range cells {
		rowLength += cell.Length
	}
	cellRow := CellRow{
		Cells:     cells,
		NumCells:  len(cells),
		RowLength: rowLength,
	}
	return TableRow{
		CellRow: &cellRow,
	}
}

type Snippet struct {
	Text  string
	Token Token
	Len   int
}

func NewSnippet(text string, token Token) Snippet {
	return Snippet{
		Text:  text,
		Token: token,
		Len:   uniseg.StringWidth(text),
	}
}

type TableCell struct {
	Content []Snippet
	Length  int
}

func NewTableCellFromStr(str string, token Token) TableCell {
	s := []Snippet{NewSnippet(str, token)}
	return NewTableCell(s)
}

func NewTableCell(snippets []Snippet) TableCell {
	var len = 0
	for _, s := range snippets {
		len += s.Len
	}
	return TableCell{
		Content: snippets,
		Length:  len,
	}
}
