package bullet_engine

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"vixac.com/got/engine"
)

func commentOut(text string) string {
	var result = ""

	var currentLine = ""
	for _, r := range text {
		if r == '\n' {
			result += "# " + currentLine + "\n"
		} else {
			currentLine += string(r)
		}
	}
	return result
}
func ignoreCommentedOut(text string) string {
	var updatedString = ""
	var lines []string = strings.Split(string(text), "\n")
	for _, line := range lines {

		if len(line) > 0 && line[0] != '#' {
			updatedString += updatedString
		}
	}
	return updatedString
}

// opens the editor with the provided initial text, and returns edits from the user.
// commentedout text is ignored from the return if commentedOut mode is enabled.
func openEditor(initialText string, commentedOut bool) (string, error) {
	tmp, err := os.CreateTemp("", "got-note-*.txt")
	if err != nil {
		return "", err
	}
	defer os.Remove(tmp.Name())

	var note = ""
	if commentedOut {
		note = commentOut(initialText)
	} else {
		note = initialText
	}

	// 3. Write existing content
	if _, err := tmp.WriteString(note); err != nil {
		return "", err
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
		return "", err
	}

	// 5. Read edited content
	updated, err := os.ReadFile(tmp.Name())
	if err != nil {
		return "", err
	}

	if !commentedOut {
		return string(updated), nil
	}
	noComments := ignoreCommentedOut(string(updated))
	return noComments, nil
}

// VX:TODO rewrite to support commented out stuff etc.
func (e *EngineBullet) OpenThenTimestamp(lookup engine.GidLookup) error {

	gid, err := e.GidLookup.InputToGid(&lookup)
	if err != nil || gid == nil {
		return err
	}
	summaryId := engine.SummaryId(gid.IntValue)
	exists, err := e.SummaryStore.Fetch([]engine.SummaryId{summaryId})
	if err != nil {
		return err
	}
	if exists != nil {
		_, ok := exists[summaryId]
		if !ok {
			return errors.New("This gid does not exist.")
		}
	}

	var note = ""
	existing, err := e.LongFormStore.LongFormNotesFor(*gid)
	if err != nil {
		return err
	}
	if existing != nil {
		allStrings := ""
		for _, v := range existing.Blocks {
			allStrings += "---------\n"
			allStrings += v.Id.ToString() + " : " + v.Edited.GoString() + "\n"
			allStrings += v.Content + "\n\n"
		}
		note = allStrings
	}
	commentedOutNotes := commentOut(note)
	// 2. Temp file
	tmp, err := os.CreateTemp("", "got-note-*.txt")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())

	// 3. Write existing content
	if _, err := tmp.WriteString(commentedOutNotes); err != nil {
		return err
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
		return err
	}

	// 5. Read edited content
	updated, err := os.ReadFile(tmp.Name())
	if err != nil {
		return err
	}

	// 6. Save back to DB
	updatedString := string(updated)
	withRemovedComments := ignoreCommentedOut(updatedString)
	if withRemovedComments == "" {
		fmt.Printf("VX: No comments made.")
		return nil
	}

	newId, err := e.LongFormStore.AppendNote(*gid, withRemovedComments)
	if err != nil {
		return err
	}
	fmt.Printf("VX: New note: %s\n", newId.ToString())

	//we send the edit event so the update time gets changed
	e.publishEditEvent(EditItemEvent{Id: engine.SummaryId(gid.IntValue)})
	return nil
}

// VX:TODO not used.
func datePrefix() string {
	line := "\n\n----------------------------\n"

	now := time.Now().UTC()

	formatted := now.Format("Mon 2 Jan 2006 15:04:05 MST")
	return line + formatted + line + "\n"

}
