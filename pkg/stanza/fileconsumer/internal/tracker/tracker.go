// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package tracker // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/fileconsumer/internal/tracker"

import (
	"context"
	"os"

	"go.opentelemetry.io/collector/component"
	"go.uber.org/zap"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/fileconsumer/internal/archive"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/fileconsumer/internal/fileset"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/fileconsumer/internal/fingerprint"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/fileconsumer/internal/reader"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/operator"
)

const (
	FileTracker    = "fileTracker"
	NoStateTracker = "noStateTracker"
)

// Interface for tracking files that are being consumed.
type Tracker interface {
	Name() string
	Add(reader *reader.Reader)
	GetCurrentFile(fp *fingerprint.Fingerprint) *reader.Reader
	GetOpenFile(fp *fingerprint.Fingerprint) *reader.Reader
	GetClosedFile(fp *fingerprint.Fingerprint) *reader.Metadata
	GetMetadata() []*reader.Metadata
	LoadMetadata(metadata []*reader.Metadata)
	CurrentPollFiles() []*reader.Reader
	PreviousPollFiles() []*reader.Reader
	ClosePreviousFiles() int
	EndPoll(context.Context)
	EndConsume() int
	TotalReaders() int
	AddUnmatched(*os.File, *fingerprint.Fingerprint)
	LookupArchive(context.Context) ([]*os.File, []*fingerprint.Fingerprint, []*reader.Metadata)
}

// fileTracker tracks known offsets for files that are being consumed by the manager.
type fileTracker struct {
	set component.TelemetrySettings

	maxBatchFiles int

	currentPollFiles  *fileset.Fileset[*reader.Reader]
	previousPollFiles *fileset.Fileset[*reader.Reader]
	knownFiles        []*fileset.Fileset[*reader.Metadata]

	unmatchedFiles []*os.File
	unmatchedFps   []*fingerprint.Fingerprint

	archive archive.Archive
}

func NewFileTracker(ctx context.Context, set component.TelemetrySettings, maxBatchFiles, pollsToArchive int, persister operator.Persister) Tracker {
	knownFiles := make([]*fileset.Fileset[*reader.Metadata], 3)
	for i := 0; i < len(knownFiles); i++ {
		knownFiles[i] = fileset.New[*reader.Metadata](maxBatchFiles)
	}
	set.Logger = set.Logger.With(zap.String("tracker", "fileTracker"))

	t := &fileTracker{
		set:               set,
		maxBatchFiles:     maxBatchFiles,
		currentPollFiles:  fileset.New[*reader.Reader](maxBatchFiles),
		previousPollFiles: fileset.New[*reader.Reader](maxBatchFiles),
		knownFiles:        knownFiles,
		archive:           archive.New(ctx, set.Logger.Named("archive"), pollsToArchive, persister),
	}
	return t
}

func (*fileTracker) Name() string {
	return FileTracker
}

func (t *fileTracker) Add(reader *reader.Reader) {
	// add a new reader for tracking
	t.currentPollFiles.Add(reader)
}

func (t *fileTracker) GetCurrentFile(fp *fingerprint.Fingerprint) *reader.Reader {
	return t.currentPollFiles.Match(fp, fileset.Equal)
}

func (t *fileTracker) GetOpenFile(fp *fingerprint.Fingerprint) *reader.Reader {
	return t.previousPollFiles.Match(fp, fileset.StartsWith)
}

func (t *fileTracker) GetClosedFile(fp *fingerprint.Fingerprint) *reader.Metadata {
	for i := 0; i < len(t.knownFiles); i++ {
		if oldMetadata := t.knownFiles[i].Match(fp, fileset.StartsWith); oldMetadata != nil {
			return oldMetadata
		}
	}
	return nil
}

func (t *fileTracker) AddUnmatched(file *os.File, fp *fingerprint.Fingerprint) {
	// exclude duplicate fingerprints
	for _, f := range t.unmatchedFps {
		if fp.Equal(f) {
			t.set.Logger.Debug("Skipping duplicate file", zap.String("path", file.Name()))
			return
		}
	}
	t.unmatchedFps = append(t.unmatchedFps, fp)
	t.unmatchedFiles = append(t.unmatchedFiles, file)
}

func (t *fileTracker) LookupArchive(ctx context.Context) ([]*os.File, []*fingerprint.Fingerprint, []*reader.Metadata) {
	// LookupArchive performs fingerprint matching against the archive and returns matched metadata, files and fingerprints.
	metadata := t.archive.FindFiles(ctx, t.unmatchedFps)
	return t.unmatchedFiles, t.unmatchedFps, metadata
}

func (t *fileTracker) GetMetadata() []*reader.Metadata {
	// return all known metadata for checkpoining
	allCheckpoints := make([]*reader.Metadata, 0, t.TotalReaders())
	for _, knownFiles := range t.knownFiles {
		allCheckpoints = append(allCheckpoints, knownFiles.Get()...)
	}

	for _, r := range t.previousPollFiles.Get() {
		allCheckpoints = append(allCheckpoints, r.Metadata)
	}
	return allCheckpoints
}

func (t *fileTracker) LoadMetadata(metadata []*reader.Metadata) {
	t.knownFiles[0].Add(metadata...)
}

func (t *fileTracker) CurrentPollFiles() []*reader.Reader {
	return t.currentPollFiles.Get()
}

func (t *fileTracker) PreviousPollFiles() []*reader.Reader {
	return t.previousPollFiles.Get()
}

func (t *fileTracker) ClosePreviousFiles() (filesClosed int) {
	// t.previousPollFiles -> t.knownFiles[0]
	for r, _ := t.previousPollFiles.Pop(); r != nil; r, _ = t.previousPollFiles.Pop() {
		t.knownFiles[0].Add(r.Close())
		filesClosed++
	}
	return
}

func (t *fileTracker) EndPoll(ctx context.Context) {
	// shift the filesets at end of every poll() call
	// t.knownFiles[0] -> t.knownFiles[1] -> t.knownFiles[2]

	// Instead of throwing it away, archive it.
	t.archive.WriteFiles(ctx, t.knownFiles[2])
	copy(t.knownFiles[1:], t.knownFiles)
	t.knownFiles[0] = fileset.New[*reader.Metadata](t.maxBatchFiles)
}

func (t *fileTracker) TotalReaders() int {
	total := t.previousPollFiles.Len()
	for i := 0; i < len(t.knownFiles); i++ {
		total += t.knownFiles[i].Len()
	}
	return total
}

func (t *fileTracker) archive(metadata *fileset.Fileset[*reader.Metadata]) {
	// We make use of a ring buffer, where each set of files is stored under a specific index.
	// Instead of discarding knownFiles[2], write it to the next index and eventually roll over.
	// Separate storage keys knownFilesArchive0, knownFilesArchive1, ..., knownFilesArchiveN, roll over back to knownFilesArchive0

	// Archiving:  ┌─────────────────────on-disk archive─────────────────────────┐
	//             |    ┌───┐     ┌───┐                     ┌──────────────────┐ |
	// index       | ▶  │ 0 │  ▶  │ 1 │  ▶      ...       ▶ │ polls_to_archive │ |
	//             | ▲  └───┘     └───┘                     └──────────────────┘ |
	//             | ▲    ▲                                                ▼     |
	//             | ▲    │ Roll over overriting older offsets, if any     ◀     |
	//             └──────│──────────────────────────────────────────────────────┘
	//                    │
	//                    │
	//                    │
	//                   start
	//                   index

	if t.pollsToArchive <= 0 || t.persister == nil {
		return
	}
	if err := t.writeArchive(t.archiveIndex, metadata); err != nil {
		t.set.Logger.Error("error faced while saving to the archive", zap.Error(err))
	}
	t.archiveIndex = (t.archiveIndex + 1) % t.pollsToArchive // increment the index
}

// readArchive loads data from the archive for a given index and returns a fileset.Filset.
func (t *fileTracker) readArchive(index int) (*fileset.Fileset[*reader.Metadata], error) {
	key := fmt.Sprintf("knownFiles%d", index)
	metadata, err := checkpoint.LoadKey(context.Background(), t.persister, key)
	if err != nil {
		return nil, err
	}
	f := fileset.New[*reader.Metadata](len(metadata))
	f.Add(metadata...)
	return f, nil
}

// writeArchive saves data to the archive for a given index and returns an error, if encountered.
func (t *fileTracker) writeArchive(index int, rmds *fileset.Fileset[*reader.Metadata]) error {
	key := fmt.Sprintf("knownFiles%d", index)
	return checkpoint.SaveKey(context.Background(), t.persister, rmds.Get(), key)
}

// FindFiles goes through archive, one fileset at a time and tries to match all fingerprints against that loaded set.
func (t *fileTracker) FindFiles(fps []*fingerprint.Fingerprint) []*reader.Metadata {
	// To minimize disk access, we first access the index, then review unmatched files and update the metadata, if found.
	// We exit if all fingerprints are matched.

	// Track number of matched fingerprints so we can exit if all matched.
	var numMatched int

	// Determine the index for reading archive, starting from the most recent and moving towards the oldest
	nextIndex := t.archiveIndex
	matchedMetadata := make([]*reader.Metadata, len(fps))

	// continue executing the loop until either all records are matched or all archive sets have been processed.
	for i := 0; i < t.pollsToArchive; i++ {
		// Update the mostRecentIndex
		nextIndex = (nextIndex - 1 + t.pollsToArchive) % t.pollsToArchive

		data, err := t.readArchive(nextIndex) // we load one fileset atmost once per poll
		if err != nil {
			t.set.Logger.Error("error while opening archive", zap.Error(err))
			continue
		}
		archiveModified := false
		for j, fp := range fps {
			if matchedMetadata[j] != nil {
				// we've already found a match for this index, continue
				continue
			}
			if md := data.Match(fp, fileset.StartsWith); md != nil {
				// update the matched metada for the index
				matchedMetadata[j] = md
				archiveModified = true
				numMatched++
			}
		}
		if !archiveModified {
			continue
		}
		// we save one fileset atmost once per poll
		if err := t.writeArchive(nextIndex, data); err != nil {
			t.set.Logger.Error("error while opening archive", zap.Error(err))
		}
		// Check if all metadata have been found
		if numMatched == len(fps) {
			return matchedMetadata
		}
	}
	return matchedMetadata
}

// noStateTracker only tracks the current polled files. Once the poll is
// complete and telemetry is consumed, the tracked files are closed. The next
// poll will create fresh readers with no previously tracked offsets.
type noStateTracker struct {
	set              component.TelemetrySettings
	maxBatchFiles    int
	currentPollFiles *fileset.Fileset[*reader.Reader]
	unmatchedFiles   []*os.File
	unmatchedFps     []*fingerprint.Fingerprint
}

func NewNoStateTracker(set component.TelemetrySettings, maxBatchFiles int) Tracker {
	set.Logger = set.Logger.With(zap.String("tracker", "noStateTracker"))
	return &noStateTracker{
		set:              set,
		maxBatchFiles:    maxBatchFiles,
		currentPollFiles: fileset.New[*reader.Reader](maxBatchFiles),
	}
}

func (*noStateTracker) Name() string {
	return NoStateTracker
}

func (t *noStateTracker) Add(reader *reader.Reader) {
	// add a new reader for tracking
	t.currentPollFiles.Add(reader)
}

func (t *noStateTracker) CurrentPollFiles() []*reader.Reader {
	return t.currentPollFiles.Get()
}

func (t *noStateTracker) GetCurrentFile(fp *fingerprint.Fingerprint) *reader.Reader {
	return t.currentPollFiles.Match(fp, fileset.Equal)
}

func (t *noStateTracker) EndConsume() (filesClosed int) {
	for r, _ := t.currentPollFiles.Pop(); r != nil; r, _ = t.currentPollFiles.Pop() {
		r.Close()
		filesClosed++
	}
	t.unmatchedFiles = make([]*os.File, 0)
	t.unmatchedFps = make([]*fingerprint.Fingerprint, 0)
	return
}

func (*noStateTracker) GetOpenFile(*fingerprint.Fingerprint) *reader.Reader { return nil }

func (*noStateTracker) GetClosedFile(*fingerprint.Fingerprint) *reader.Metadata { return nil }

func (*noStateTracker) GetMetadata() []*reader.Metadata { return nil }

func (*noStateTracker) LoadMetadata([]*reader.Metadata) {}

func (*noStateTracker) PreviousPollFiles() []*reader.Reader { return nil }

func (*noStateTracker) ClosePreviousFiles() int { return 0 }

func (*noStateTracker) EndPoll(context.Context) {}

func (*noStateTracker) TotalReaders() int { return 0 }

func (t *noStateTracker) AddUnmatched(file *os.File, fp *fingerprint.Fingerprint) {
	t.unmatchedFiles = append(t.unmatchedFiles, file)
	t.unmatchedFps = append(t.unmatchedFps, fp)
}

func (t *noStateTracker) LookupArchive(context.Context) ([]*os.File, []*fingerprint.Fingerprint, []*reader.Metadata) {
	return t.unmatchedFiles, t.unmatchedFps, make([]*reader.Metadata, len(t.unmatchedFps))
}
