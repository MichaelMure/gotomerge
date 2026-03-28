// Package format handles reading and writing the Automerge binary file format.
//
// # File structure
//
// An Automerge file is a sequence of chunks. Each chunk starts with a 4-byte
// magic number, a 4-byte checksum, a type byte, and a LEB128-encoded length,
// followed by the payload. Use ReadChunk to parse one chunk at a time.
//
// There are two chunk types:
//
// A ChangeChunk records a single edit made by one peer. It is the unit of
// replication: peers exchange change chunks to converge on the same document
// state. Each change carries a set of operations (inserts, deletes, assignments),
// metadata (author, timestamp, message), and the hashes of the changes it
// directly depends on. Operations are listed in creation order. A change does
// not know whether its operations were later overwritten by other peers.
//
// A DocumentChunk is a fully merged snapshot of the entire document. Instead of
// replaying a history of edits, it stores every operation from every change in
// one place, with successor lists that indicate which operations were later
// overwritten. An operation with no successors is the current live value. A
// DocumentChunk also stores per-change metadata (author, sequence number, maxOp,
// dependencies) so the change history remains queryable, but it does not
// duplicate the operations — a change's operations are identified by their
// counter range, derived from consecutive maxOp values.
//
// # Typical lifecycle
//
// Every edit a peer makes produces a ChangeChunk. Peers replicate by exchanging
// these chunks — sending only the changes the other side hasn't seen yet.
//
// Periodically, the accumulated changes are compacted into a DocumentChunk. This
// merges all operations into a single columnar structure and records which values
// are still live. The individual ChangeChunks are no longer needed once a
// DocumentChunk covers them.
//
// A file saved to disk is usually a single DocumentChunk (the latest compaction)
// optionally followed by any ChangeChunks that arrived after it. A reader applies
// the ChangeChunks on top of the snapshot to reach the current state.
//
// # Hashes
//
// Every change is globally identified by a hash: the SHA-256 of its serialized
// binary representation (type byte + length + payload). Two peers that have the
// same change — whether stored compressed or uncompressed — agree on its hash.
// ReadChunk verifies the checksum and stores the full hash in ChangeChunk.Hash.
// DocumentChunks have no individual hash identity.
package format
