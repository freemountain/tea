package unidiff

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"

	"gitea.com/noerw/unidiff-comments/types"
)

const (
	stateStartOfDiff       = "stateStartOfDiff"
	stateDiffHeader        = "stateDiffHeader"
	stateHunkHeader        = "stateHunkHeader"
	stateHunkBody          = "stateHunkBody"
	stateComment           = "stateComment"
	stateCommentDelim      = "stateCommentDelim"
	stateCommentHeader     = "stateCommentHeader"
	stateDiffComment       = "stateDiffComment"
	stateDiffCommentDelim  = "stateDiffCommentDelim"
	stateDiffCommentHeader = "stateDiffCommentHeader"

	ignorePrefix = "###"
)

var (
	reDiffHeader = regexp.MustCompile(
		`^--- |^\+\+\+ `)

	reGitDiffHeader = regexp.MustCompile(
		`^diff |^index `)

	reFromFile = regexp.MustCompile(
		`^--- (\S+)(\s+(.*))`)

	reToFile = regexp.MustCompile(
		`^\+\+\+ (\S+)(\s+(.*))`)

	reHunk = regexp.MustCompile(
		`^@@ -(\d+)(,(\d+))? \+(\d+)(,(\d+))? @@`)

	reSegmentContext = regexp.MustCompile(
		`^ `)

	reSegmentAdded = regexp.MustCompile(
		`^\+`)

	reSegmentRemoved = regexp.MustCompile(
		`^-`)

	reCommentDelim = regexp.MustCompile(
		`^#\s+---`)

	reCommentHeader = regexp.MustCompile(
		`^#\s+\[(\d+)@(\d+)\]\s+\|([^|]+)\|(.*)`)

	reCommentText = regexp.MustCompile(
		`^#(\s*)(.*)\s*`)

	reIndent = regexp.MustCompile(
		`^#(\s+)`)

	reEmptyLine = regexp.MustCompile(
		`^\n$`)

	reIgnoredLine = regexp.MustCompile(
		`^` + ignorePrefix)
)

type parser struct {
	state      string
	changeset  types.Changeset
	diff       *types.Diff
	hunk       *types.Hunk
	segment    *types.Segment
	comment    *types.Comment
	line       *types.Line
	lineNumber int

	segmentType  string
	commentsList []*types.Comment
}

type Error struct {
	LineNumber int
	Message    string
}

func (err Error) Error() string {
	return fmt.Sprintf("line %d: %s", err.LineNumber, err.Message)
}

func ReadChangeset(r io.Reader) (types.Changeset, error) {
	buffer := bufio.NewReader(r)

	current := parser{}
	current.state = stateStartOfDiff

	for {
		current.lineNumber++

		line, err := buffer.ReadString('\n')
		if err != nil {
			break
		}

		if reIgnoredLine.MatchString(line) {
			continue
		}

		err = current.switchState(line)
		if err != nil {
			return current.changeset, err
		}

		err = current.createNodes(line)
		if err != nil {
			return current.changeset, err
		}

		err = current.locateNodes(line)
		if err != nil {
			return current.changeset, err
		}

		err = current.parseLine(line)
		if err != nil {
			return current.changeset, err
		}
	}

	for _, comment := range current.commentsList {
		comment.Text = strings.TrimSpace(comment.Text)
	}

	return current.changeset, nil
}

func (current *parser) switchState(line string) error {
	inComment := false

	switch current.state {
	case stateStartOfDiff:
		switch {
		case reDiffHeader.MatchString(line), reGitDiffHeader.MatchString(line):
			current.state = stateDiffHeader
		case reCommentText.MatchString(line):
			inComment = true
		case reEmptyLine.MatchString(line):
			// body intentionally left empty
		default:
			return Error{
				current.lineNumber,
				"expected diff header, but none found",
			}
		}
	case stateDiffHeader:
		switch {
		case reHunk.MatchString(line):
			current.state = stateHunkHeader
		}
	case stateDiffComment, stateDiffCommentDelim, stateDiffCommentHeader:
		switch {
		case reDiffHeader.MatchString(line), reGitDiffHeader.MatchString(line):
			current.state = stateDiffHeader
		case reCommentText.MatchString(line):
			inComment = true
		case reEmptyLine.MatchString(line):
			current.state = stateStartOfDiff
		}
	case stateHunkHeader:
		current.state = stateHunkBody
		fallthrough
	case stateHunkBody, stateComment, stateCommentDelim, stateCommentHeader:
		switch {
		case reSegmentContext.MatchString(line):
			current.state = stateHunkBody
			current.segmentType = types.SegmentTypeContext
		case reSegmentRemoved.MatchString(line):
			current.state = stateHunkBody
			current.segmentType = types.SegmentTypeRemoved
		case reSegmentAdded.MatchString(line):
			current.state = stateHunkBody
			current.segmentType = types.SegmentTypeAdded
		case reHunk.MatchString(line):
			current.state = stateHunkHeader
		case reCommentText.MatchString(line):
			inComment = true
		case reGitDiffHeader.MatchString(line):
			current.state = stateDiffHeader
			current.diff = nil
			current.hunk = nil
			current.segment = nil
			current.line = nil
		case reEmptyLine.MatchString(line):
			current.state = stateStartOfDiff
			current.diff = nil
			current.hunk = nil
			current.segment = nil
			current.line = nil
		}
	}

	if !inComment {
		current.comment = nil
	} else {
		switch current.state {
		case stateStartOfDiff:
			fallthrough
		case stateDiffComment, stateDiffCommentDelim, stateDiffCommentHeader:
			switch {
			case reCommentDelim.MatchString(line):
				current.state = stateDiffCommentDelim
			case reCommentHeader.MatchString(line):
				current.state = stateDiffCommentHeader
			case reCommentText.MatchString(line):
				current.state = stateDiffComment
			}
		case stateHunkBody:
			fallthrough
		case stateComment, stateCommentDelim, stateCommentHeader:
			switch {
			case reCommentDelim.MatchString(line):
				current.state = stateCommentDelim
			case reCommentHeader.MatchString(line):
				current.state = stateCommentHeader
			case reCommentText.MatchString(line):
				current.state = stateComment
			}
		}
	}

	// Uncomment for debug state switching
	// fmt.Printf("%20s : %#v\n", current.state, line)

	return nil
}

func (current *parser) createNodes(line string) error {
	switch current.state {
	case stateDiffComment:
		if current.comment != nil {
			break
		}
		fallthrough
	case stateDiffCommentDelim, stateDiffCommentHeader:
		current.comment = &types.Comment{}
		fallthrough
	case stateDiffHeader:
		if current.diff == nil {
			current.diff = &types.Diff{}
			current.changeset.Diffs = append(current.changeset.Diffs,
				current.diff)
		}
	case stateHunkHeader:
		current.hunk = &types.Hunk{}
		current.segment = &types.Segment{}
	case stateCommentDelim, stateCommentHeader:
		current.comment = &types.Comment{}
	case stateComment:
		if current.comment == nil {
			current.comment = &types.Comment{}
		}
	case stateHunkBody:
		if current.segment.Type != current.segmentType {
			current.segment = &types.Segment{Type: current.segmentType}
			current.hunk.Segments = append(current.hunk.Segments,
				current.segment)
		}

		current.line = &types.Line{}
		current.segment.Lines = append(current.segment.Lines, current.line)
	}

	return nil
}

func (current *parser) locateNodes(line string) error {
	switch current.state {
	case stateComment, stateDiffComment:
		current.locateComment(line)
	case stateHunkBody:
		current.locateLine(line)
	}

	return nil
}

func (current *parser) locateComment(line string) error {
	if current.comment.Parented || strings.TrimSpace(line) == "#" {
		return nil
	}

	current.commentsList = append(current.commentsList, current.comment)
	current.comment.Parented = true

	if current.hunk != nil {
		current.comment.Anchor.LineType = current.segment.Type
		current.comment.Anchor.Line = current.segment.GetLineNum(current.line)
		current.comment.Anchor.Path = current.diff.Destination.ToString
		current.comment.Anchor.SrcPath = current.diff.Source.ToString
	}

	current.comment.Indent = getIndentSize(line)

	parent := current.findParentComment(current.comment)
	if parent != nil {
		parent.Comments = append(parent.Comments, current.comment)
	} else {
		if current.line != nil {
			current.diff.LineComments = append(current.diff.LineComments,
				current.comment)
			current.line.Comments = append(current.line.Comments,
				current.comment)
		} else {
			current.diff.FileComments = append(current.diff.FileComments,
				current.comment)
		}
	}

	return nil
}

func (current *parser) locateLine(line string) error {
	sourceOffset := current.hunk.SourceLine - 1
	destinationOffset := current.hunk.DestinationLine - 1
	if len(current.hunk.Segments) > 1 {
		prevSegment := current.hunk.Segments[len(current.hunk.Segments)-2]
		lastLine := prevSegment.Lines[len(prevSegment.Lines)-1]
		sourceOffset = lastLine.Source
		destinationOffset = lastLine.Destination
	}
	hunkLength := int64(len(current.segment.Lines))
	switch current.segment.Type {
	case types.SegmentTypeContext:
		current.line.Source = sourceOffset + hunkLength
		current.line.Destination = destinationOffset + hunkLength
	case types.SegmentTypeAdded:
		current.line.Source = sourceOffset
		current.line.Destination = destinationOffset + hunkLength
	case types.SegmentTypeRemoved:
		current.line.Source = sourceOffset + hunkLength
		current.line.Destination = destinationOffset
	}

	return nil
}

func (current *parser) parseLine(line string) error {
	switch current.state {
	case stateDiffHeader:
		current.parseDiffHeader(line)
	case stateHunkHeader:
		current.parseHunkHeader(line)
	case stateHunkBody:
		current.parseHunkBody(line)
	case stateComment, stateDiffComment:
		current.parseComment(line)
	case stateCommentHeader, stateDiffCommentHeader:
		current.parseCommentHeader(line)
	}

	return nil
}

func (current *parser) parseDiffHeader(line string) error {
	switch {
	case reFromFile.MatchString(line):
		matches := reFromFile.FindStringSubmatch(line)
		current.changeset.Path = matches[1]
		current.diff.Source.ToString = matches[1]
		current.changeset.FromHash = matches[3]
		current.diff.Attributes.FromHash = []string{matches[3]}
	case reToFile.MatchString(line):
		matches := reToFile.FindStringSubmatch(line)
		current.diff.Destination.ToString = matches[1]
		current.changeset.ToHash = matches[3]
		current.diff.Attributes.ToHash = []string{matches[3]}
	default:
		return Error{
			current.lineNumber,
			"expected diff header, but not found",
		}
	}
	return nil
}

func (current *parser) parseHunkHeader(line string) error {
	matches := reHunk.FindStringSubmatch(line)
	current.hunk.SourceLine, _ = strconv.ParseInt(matches[1], 10, 64)
	current.hunk.SourceSpan, _ = strconv.ParseInt(matches[3], 10, 64)
	current.hunk.DestinationLine, _ = strconv.ParseInt(matches[4], 10, 64)
	current.hunk.DestinationSpan, _ = strconv.ParseInt(matches[6], 10, 64)
	current.diff.Hunks = append(current.diff.Hunks, current.hunk)

	return nil
}

func (current *parser) parseHunkBody(line string) error {
	current.line.Line = line[1 : len(line)-1]
	return nil
}

func (current *parser) parseCommentHeader(line string) error {
	matches := reCommentHeader.FindStringSubmatch(line)
	current.comment.Author.DisplayName = strings.TrimSpace(matches[3])
	current.comment.Id, _ = strconv.ParseInt(matches[1], 10, 64)
	updatedDate, _ := time.ParseInLocation(time.ANSIC,
		strings.TrimSpace(matches[4]),
		time.Local)
	current.comment.UpdatedDate = types.UnixTimestamp(updatedDate.Unix() * 1000)

	version, _ := strconv.ParseInt(matches[2], 10, 64)
	current.comment.Version = int(version)

	return nil
}

func (current *parser) parseComment(line string) error {
	matches := reCommentText.FindStringSubmatch(line)
	if len(matches[1]) < current.comment.Indent {
		return Error{
			LineNumber: current.lineNumber,
			Message: fmt.Sprintf(
				"unexpected indent, should be at least: %d",
				current.comment.Indent,
			),
		}
	}

	indentedLine := matches[1][current.comment.Indent:] + matches[2]
	current.comment.Text += "\n" + indentedLine

	return nil
}

func (current *parser) findParentComment(comment *types.Comment) *types.Comment {
	for i := len(current.commentsList) - 1; i >= 0; i-- {
		c := current.commentsList[i]
		if comment.Indent > c.Indent {
			return c
		}
	}

	return nil
}

func getIndentSize(line string) int {
	matches := reIndent.FindStringSubmatch(line)
	if len(matches) == 0 {
		return 0
	}

	return len(matches[1])
}
