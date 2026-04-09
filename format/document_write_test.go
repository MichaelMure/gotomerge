package format

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"gotomerge/types"
	ioutil "gotomerge/utils/io"
)

// TestDocumentChunkRoundTrip reads each listed testdata file, re-encodes the
// document chunk, and verifies the output is byte-for-byte identical to the
// original. All test files contain a single document chunk with no per-column
// DEFLATE compression (columns are small), so byte-identical output is expected.
func TestDocumentChunkRoundTrip(t *testing.T) {
	cases := []string{
		"list_sequential.automerge",
		"list_concurrent_inserts.automerge",
		"list_with_delete.automerge",
		"list_insert_after_deleted.automerge",
		"text_sequential.automerge",
		"map_conflict.automerge",
		"map_delete.automerge",
		"exemplar",
	}

	for _, name := range cases {
		t.Run(name, func(t *testing.T) {
			original, err := os.ReadFile("../testdata/" + name)
			require.NoError(t, err)

			r := ioutil.NewSubReader(original)
			var out bytes.Buffer

			for !r.Empty() {
				chunk, toSkip, err := ReadChunk(r)
				require.NoError(t, err)
				require.NoError(t, r.Skip(toSkip))

				dc, ok := chunk.(*DocumentChunk)
				if !ok {
					continue // skip any change chunks (some files have both)
				}

				// Re-encode through the production path.
				mapper := types.IdentityActorMapper(uint32(len(dc.Actors)))
				// Re-encode change metadata.
				changes := NewChangeMetaWriter()
				for m, err := range dc.Changes() {
					require.NoError(t, err)
					changes.Append(m)
				}

				ops := NewDocOpsWriter()
				for op, err := range dc.Operations() {
					require.NoError(t, err)
					ops.Append(op.Object, op.Key, op.Id, op.Insert, op.Action, op.Successors, mapper)
				}
				require.NoError(t, WriteDocument(&out, dc.Actors, dc.Heads, dc.HeadIndexes, changes, ops))
			}

			require.Equal(t, original, out.Bytes())
		})
	}
}
