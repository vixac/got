package console

type ConsoleTable struct {
	Rows          []TableRow
	ColumnWidths  []int
	ColumnPadding string
}

func (c *ConsoleTable) Render(printer Messenger, scheme Theme) {

	for _, row := range c.Rows {
		var messages []Message

		for i, cell := range row.Cells {
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
			paddingStr := FitString("", paddingRequired, " ")
			messages = append(messages, Message{Message: paddingStr})
		}
		printer.PrintInLine(messages)
	}
}

func NewConsoleTable(rows []TableRow) ConsoleTable {
	var maxWidths []int
	for _, r := range rows {
		numCells := len(r.Cells)
		for col := 0; col < numCells; col++ {
			minColWidthThisCell := r.Cells[col].Length
			if len(maxWidths) <= col {
				maxWidths = append(maxWidths, 0)
			}
			if minColWidthThisCell > maxWidths[col] {
				maxWidths[col] = minColWidthThisCell
			}
		}
	}

	var padding = ""
	for range 3 {
		padding += " "
	}
	return ConsoleTable{
		Rows:          rows,
		ColumnWidths:  maxWidths,
		ColumnPadding: padding,
	}
}

// this becomes a bit of an enum eventually. cells or divider
type TableRow struct {
	Cells    []TableCell
	NumCells int
}

func NewTableRow(cells []TableCell) TableRow {
	return TableRow{
		Cells:    cells,
		NumCells: len(cells),
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
		Len:   len(text),
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
