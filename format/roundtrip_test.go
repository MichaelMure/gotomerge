package format

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	ioutil "gotomerge/utils/io"
)

// TestRoundTrip reads each listed testdata file, re-encodes every chunk, and
// verifies the output is byte-for-byte identical to the original. Failures
// indicate either a missing write path (document chunks) or an encoding
// discrepancy that needs fixing.
func TestRoundTrip(t *testing.T) {
	cases := []struct {
		name string
	}{
		{name: "counter_value_is_ok.automerge"},
		// Synthetic: val_meta claims 2 bytes for counter 16, but canonical encoding is 1 byte.
		// We correctly decode and re-encode canonically; exact bytes cannot be preserved.
		// {name: "counter_value_has_incorrect_meta.automerge"},
		{name: "list_sequential.automerge"},
		{name: "list_concurrent_inserts.automerge"},
		{name: "list_with_delete.automerge"},
		{name: "list_insert_after_deleted.automerge"},
		{name: "text_sequential.automerge"},
		{name: "map_conflict.automerge"},
		{name: "map_delete.automerge"},
		{name: "two_change_chunks.automerge"},
		// Synthetic: change payload < 256 bytes but stored compressed; real encoders never do this.
		// We correctly decompress on read and write uncompressed (DEFLATE_MIN_SIZE = 256).
		// {name: "two_change_chunks_compressed.automerge"},
		{name: "two_change_chunks_out_of_order.automerge"},
		{name: "exemplar"},
		{name: "text-edits.amrg"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.name == "text-edits.amrg" && testing.Short() {
				t.Skip("skipping large file under -short")
			}

			original, err := os.ReadFile("../testdata/" + tc.name)
			require.NoError(t, err)

			r := ioutil.NewSubReader(original)
			var out bytes.Buffer

			for !r.Empty() {
				chunk, toSkip, err := ReadChunk(r)
				require.NoError(t, err)
				require.NoError(t, r.Skip(toSkip))

				switch cc := chunk.(type) {
				case *ChangeChunk:
					enc := NewChangeOpsWriter()
					// Operations already use local actor indices (0 = own actor,
					// 1..N = OtherActors). Identity localOf preserves them as-is.
					n := uint32(1 + len(cc.OtherActors))
					localOf := make(map[uint32]uint32, n)
					for i := uint32(0); i < n; i++ {
						localOf[i] = i
					}
					for op, err := range cc.Operations() {
						require.NoError(t, err)
						enc.Append(op.Object, op.Key, op.Insert, op.Action, op.Predecessors, localOf)
					}
					require.NoError(t, enc.Finalise())

					cc2 := &ChangeChunk{
						Dependencies: cc.Dependencies,
						Actor:        cc.Actor,
						SeqNum:       cc.SeqNum,
						StartOp:      cc.StartOp,
						Time:         cc.Time,
						Message:      cc.Message,
						OtherActors:  cc.OtherActors,
					}
					require.NoError(t, WriteChange(&out, cc2, enc))

				case *DocumentChunk:
					require.NoError(t, WriteDocument(&out, cc))
				}
			}

			if !bytes.Equal(original, out.Bytes()) {
				printChunks := func(data []byte) string {
					r := ioutil.NewSubReader(data)
					var res strings.Builder
					for !r.Empty() {
						chunk, toSkip, err := ReadChunk(r)
						if err != nil {
							t.Logf("  parse error: %v", err)
							break
						}
						res.WriteString(fmt.Sprint(chunk))
						_ = r.Skip(toSkip)
					}
					return res.String()
				}
				require.Equal(t, printChunks(original), printChunks(out.Bytes()))
			}
			require.Equal(t, original, out.Bytes())
		})
	}
}
