package engine_util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// --- commentOut ---

func TestCommentOut_Empty(t *testing.T) {
	assert.Equal(t, "", CommentOut(""))
}

func TestCommentOut_SingleLineNoTrailingNewline(t *testing.T) {
	assert.Equal(t, "# hello", CommentOut("hello"))
}

func TestCommentOut_SingleLineWithTrailingNewline(t *testing.T) {
	assert.Equal(t, "# hello\n", CommentOut("hello\n"))
}

func TestCommentOut_MultipleLines(t *testing.T) {
	input := "line one\nline two\nline three\n"
	expected := "# line one\n# line two\n# line three\n"
	assert.Equal(t, expected, CommentOut(input))
}

func TestCommentOut_MultipleLines_NoTrailingNewline(t *testing.T) {
	input := "line one\nline two"
	expected := "# line one\n# line two"
	assert.Equal(t, expected, CommentOut(input))
}

func TestCommentOut_BlankLinesAreCommentedOut(t *testing.T) {
	// Blank lines within the text should also get the # prefix
	input := "first\n\nsecond\n"
	expected := "# first\n# \n# second\n"
	assert.Equal(t, expected, CommentOut(input))
}

func TestCommentOut_AlreadyCommentedLine(t *testing.T) {
	// Lines already starting with # get double-commented — the function is
	// not idempotent, and that is expected (it treats its input as plain text).
	assert.Equal(t, "# # already", CommentOut("# already"))
}

// --- ignoreCommentedOut ---

func TestIgnoreCommentedOut_Empty(t *testing.T) {
	assert.Equal(t, "", IgnoreCommentedOut(""))
}

func TestIgnoreCommentedOut_AllPlainText(t *testing.T) {
	input := "my note\nsecond line\n"
	assert.Equal(t, input, IgnoreCommentedOut(input))
}

func TestIgnoreCommentedOut_AllCommented(t *testing.T) {
	input := "# old note line 1\n# old note line 2\n"
	assert.Equal(t, "", IgnoreCommentedOut(input))
}

func TestIgnoreCommentedOut_MixedContent(t *testing.T) {
	// The user types new text at the top; old notes are below as # lines
	input := "new note\n# old note line 1\n# old note line 2\n"
	assert.Equal(t, "new note\n", IgnoreCommentedOut(input))
}

func TestIgnoreCommentedOut_PreservesBlankLinesBetweenParagraphs(t *testing.T) {
	// Blank lines in the user's new content should be preserved (git commit UX)
	input := "paragraph one\n\nparagraph two\n# old commented note\n"
	assert.Equal(t, "paragraph one\n\nparagraph two\n", IgnoreCommentedOut(input))
}

func TestIgnoreCommentedOut_HashMidLine_IsNotComment(t *testing.T) {
	// Only lines that START with # are stripped; # elsewhere is preserved
	input := "tag #work done\n# this is a comment\n"
	assert.Equal(t, "tag #work done\n", IgnoreCommentedOut(input))
}

// --- round-trip: the full git-commit-style UX ---

func TestRoundTrip_UserPrependsNewContent(t *testing.T) {
	// Existing notes are commented out and shown below the cursor.
	// The user types new text at the top. ignoreCommentedOut returns only the new text.
	existingNote := "old note line one\nold note line two\n"
	commentedExisting := CommentOut(existingNote)

	// Simulate the user prepending a new note above the commented block
	fileContents := "my brand new note\n" + commentedExisting

	result := IgnoreCommentedOut(fileContents)
	assert.Equal(t, "my brand new note\n", result)
}

func TestRoundTrip_UserWritesMultiParagraphNote(t *testing.T) {
	existingNote := "previous entry\n"
	commentedExisting := CommentOut(existingNote)

	fileContents := "first paragraph\n\nsecond paragraph\n" + commentedExisting

	result := IgnoreCommentedOut(fileContents)
	assert.Equal(t, "first paragraph\n\nsecond paragraph\n", result)
}

func TestRoundTrip_UserWritesNothing(t *testing.T) {
	// File contains only the commented-out existing notes; user wrote nothing.
	existingNote := "old stuff\n"
	fileContents := CommentOut(existingNote)

	result := IgnoreCommentedOut(fileContents)
	// Only the trailing empty string from the split remains, joined back to ""
	assert.Equal(t, "", result)
}
