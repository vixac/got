package bullet_engine

import (
	"os"
	"os/exec"
	"time"

	"vixac.com/got/engine"
)

func (e *EngineBullet) OpenThenTimestamp(lookup engine.GidLookup) error {
	gid, err := e.GidLookup.InputToGid(&lookup)
	if err != nil || gid == nil {
		return err
	}

	var note = ""
	existing, err := e.LongFormStore.LongFormFor(gid.IntValue)
	if err != nil {
		return err
	}
	if existing != nil {
		note = *existing
	}
	// 2. Temp file
	tmp, err := os.CreateTemp("", "got-note-*.txt")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())

	// 3. Write existing content
	if _, err := tmp.WriteString(note); err != nil {
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
	if updatedString == note { //no changes, don't save
		return nil
	}
	datedString := datePrefix() + updatedString
	return e.LongFormStore.UpsertItem(gid.IntValue, datedString)
}

func datePrefix() string {
	line := "\n\n----------------------------\n"

	now := time.Now().UTC()

	formatted := now.Format("Mon 2 Jan 2006 15:04:05 MST")
	return line + formatted + line + "\n"

}
