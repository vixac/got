package bullet_engine

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"
	"time"

	"vixac.com/got/engine"
)

const line = "  ---------------------------------------------------\n"

func commentOut(text string) string {
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

func ignoreCommentedOut(text string) string {
	var kept []string
	for _, line := range strings.Split(text, "\n") {
		if len(line) == 0 || line[0] != '#' {
			kept = append(kept, line)
		}
	}
	return strings.Join(kept, "\n")
}

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
		for _, v := range slices.Backward(existing.Blocks) {

			allStrings += line
			allStrings += "  " + datePrefix(v.Edited) + "  " + v.Id.ToString() + "\n"
			allStrings += line
			allStrings += "\n  " + v.Content + "\n\n"
		}
		note = allStrings
	}
	commentedOutNotes := "\n\n" + commentOut(note)
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

	whiteSpaceCheck := strings.ReplaceAll(withRemovedComments, "\n", "")
	whiteSpaceCheck = strings.ReplaceAll(whiteSpaceCheck, " ", "")
	if whiteSpaceCheck == "" {
		fmt.Printf("VX: No comments made.")
		return nil
	}
	fmt.Printf("VX: the entire comment was '%s'\n", withRemovedComments)

	newId, err := e.LongFormStore.AppendNote(*gid, withRemovedComments)
	if err != nil {
		return err
	}
	fmt.Printf("VX: New note: %s\n", newId.ToString())

	//we send the edit event so the update time gets changed
	e.publishEditEvent(EditItemEvent{Id: engine.SummaryId(gid.IntValue)})
	return nil
}

func datePrefix(date time.Time) string {
	formatted := date.Format("Mon 2 Jan 2006 15:04:05 MST")
	return formatted
}
