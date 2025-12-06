package console

import "unicode/utf8"

type ConsoleTable struct {
	Rows          []TableRow
	ColumnWidths  []int
	ColumnPadding string
}

const (
	paddingSize = 0
)

func nchars(b byte, n int) string {
	s := make([]byte, n)
	for i := 0; i < n; i++ {
		s[i] = b
	}
	return string(s)
}
func (c *ConsoleTable) Render(printer Messenger, scheme Theme) {

	//the int is the index, and the message group is the rendered row.
	contentRows := make(map[int]MessageGroup)
	for rowNumber, row := range c.Rows {

		if row.CellRow == nil {

		} else {
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
				spaceCharacter := " " //useful for debugging to make this and X or something.
				paddingStr := FitString("", paddingRequired, spaceCharacter)
				messages = append(messages, Message{Message: paddingStr})

			}
			var total = 0
			for _, m := range messages {
				total += utf8.RuneCountInString(m.Message)
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
				Color:   scheme.ColorFor(TokenPrimary{}).Col(),
			}
			printer.Print(dividerMessage)
		} else {
			group := contentRows[i] //if this doesn't exist its a dev error, as we just populated this.
			printer.PrintInLine(group.Messages)
		}

	}
}

func NewConsoleTable(rows []TableRow) ConsoleTable {
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
	}
}

// this becomes a bit of an enum eventually. cells or divider
type CellRow struct {
	Cells     []TableCell
	NumCells  int
	RowLength int
}
type DividerRow struct {
	Separator byte
	Token     Token
}
type TableRow struct {
	CellRow    *CellRow
	DividerRow *DividerRow
}

func NewDividerRow(separator byte, token Token) TableRow {
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
		Len:   utf8.RuneCountInString(text),
	}
}

type TableCell struct {
	Content []Snippet
	Length  int
}

func (c *TableCell) PlainStr() string {
	var str = "{"
	for _, s := range c.Content {
		str += s.Text + ","
	}
	str += "}"
	return str
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
