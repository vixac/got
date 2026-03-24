package engine_util

import (
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"
	"time"

	"vixac.com/got/engine"
)

const Line = "  ---------------------------------------------------\n"

// takes a commented out string, and a filename for a temp file, and returns the users
// uncommented modifications to the file.
func OpenTextEditorWithCommentedOutString(commented string) (*string, error) {

	// 2. Temp file
	tmp, err := os.CreateTemp("", "got-note-*.txt")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmp.Name())

	// 3. Write existing content
	if _, err := tmp.WriteString(commented); err != nil {
		return nil, err
	}
	tmp.Close()

	// 4. Launch editor
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}

	cmd := exec.Command(editor, tmp.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	// 5. Read edited content
	updated, err := os.ReadFile(tmp.Name())
	if err != nil {
		return nil, err
	}

	// 6. Save back to DB
	updatedString := string(updated)
	withRemovedComments := IgnoreCommentedOut(updatedString)

	whiteSpaceCheck := strings.ReplaceAll(withRemovedComments, "\n", "")
	whiteSpaceCheck = strings.ReplaceAll(whiteSpaceCheck, " ", "")
	if whiteSpaceCheck == "" {
		return nil, nil
	}
	fmt.Printf("VX: the entire comment was '%s'\n", withRemovedComments)
	return &withRemovedComments, nil

}
func ConsolidateBlocksIntoCommentedString(blockResult *engine.LongFormBlockResult) string {

	var note = ""

	if blockResult != nil {
		allStrings := ""
		for _, v := range slices.Backward(blockResult.Blocks) {

			allStrings += Line
			allStrings += "  " + datePrefix(v.Edited) + "  " + v.Id.ToString() + "\n"
			allStrings += Line
			allStrings += "\n  " + v.Content + "\n\n"
		}
		note = allStrings
	}
	commentedOutNotes := "\n\n" + CommentOut(note)
	return commentedOutNotes
}

func datePrefix(date time.Time) string {
	formatted := date.Format("Mon 2 Jan 2006 15:04:05 MST")
	return formatted
}

func CommentOut(text string) string {
	var result = ""
	var currentLine = ""
	for _, r := range text {
		if r == '\n' {
			result += "# " + currentLine + "\n"
			currentLine = ""
		} else {
			currentLine += string(r)
		}
	}
	if currentLine != "" {
		result += "# " + currentLine
	}
	return result
}

func IgnoreCommentedOut(text string) string {
	var kept []string
	for _, line := range strings.Split(text, "\n") {
		if len(line) == 0 || line[0] != '#' {
			kept = append(kept, line)
		}
	}
	return strings.Join(kept, "\n")
}
