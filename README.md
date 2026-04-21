<div align="center">
  <h1 align="center">gotomerge</h1>

  <p>
    <a href="https://github.com/MichaelMure/gotomerge/tags">
        <img alt="GitHub Tag" src="https://img.shields.io/github/v/tag/MichaelMure/gotomerge">
    </a>
    <a href="https://github.com/MichaelMure/gotomerge/actions?query=">
      <img src="https://github.com/MichaelMure/gotomerge/actions/workflows/gotest.yml/badge.svg" alt="Build Status">
    </a>
    <a href="https://MichaelMure.github.io/gotomerge/dev/bench/">
        <img alt="Go benchmarks" src="https://img.shields.io/badge/Benchmarks-go-blue">
    </a>
    <a href="https://github.com/MichaelMure/gotomerge/blob/v1/LICENSE">
        <img alt="MIT License" src="https://img.shields.io/badge/License-MIT-green">
    </a>
    <a href="https://pkg.go.dev/github.com/MichaelMure/gotomerge">
      <img src="https://img.shields.io/badge/Docs-godoc-blue" alt="Docs">
    </a>
  </p>
</div>

# gotomerge

A Go implementation of the [Automerge](https://automerge.org) CRDT.

Automerge documents can be edited concurrently by multiple peers without coordination. Changes are merged automatically and deterministically — no conflict resolution code required. The binary format is interoperable with other the canonical Automerge implementation (Rust).

## Architecture

An overview of the internal components and data flows is available in [architecture.md](architecture.md).

## Install

```
go get github.com/MichaelMure/gotomerge
```

## Usage

```go
doc := gotomerge.NewDocument()

doc.Change(func(txn *gotomerge.Txn) error {
    txn.Set("title", "Hello")
    txn.Set("published", false)
    txn.Map("meta").Set("author", "alice")
    txn.List("tags").Append("go")
    txn.Text("body").Splice(0, 0, "First draft.")
    return nil
})

// Read scalars with As[T] — no type assertions.
title, _ := gotomerge.As[string](doc.Get("title"))

// Unmarshal a nested map into a struct.
// Field names default to lowercase, so no tags needed for the common case.
type Meta struct {
    Author string
}
meta, _ := gotomerge.As[Meta](doc.Get("meta"))

// Or work with container views directly.
if tags, ok := doc.List("tags"); ok {
    fmt.Println(tags.Len())
}
```

## Merging

Fork a document by saving and reloading, let each peer make changes independently, then merge:

```go
var snap bytes.Buffer
base.SaveIncremental(&snap)
b := snap.Bytes()

peer1, _ := gotomerge.LoadDocument(bytes.NewReader(b))
peer2, _ := gotomerge.LoadDocument(bytes.NewReader(b))

peer1.Change(func(txn *gotomerge.Txn) error {
    txn.Set("status", "published")
    return nil
})
peer2.Change(func(txn *gotomerge.Txn) error {
    txn.Set("author", "alice")
    return nil
})

peer1.Merge(peer2)
// peer1 now has both "status" and "author".
```

When two peers set the same key concurrently, Automerge picks a winner deterministically (highest actor ID). If you need both values, use `MapView.GetAll` (conflicts API — not yet implemented).

## Data types

| Go type     | Write                           | Read with `As[T]`                  |
|-------------|---------------------------------|------------------------------------|
| `string`    | `txn.Set(k, "…")`               | `As[string]`                       |
| `bool`      | `txn.Set(k, true)`              | `As[bool]`                         |
| `int64`     | `txn.Set(k, int64(n))`          | `As[int64]` (also int8…int32)      |
| `float64`   | `txn.Set(k, 3.14)`              | `As[float64]`                      |
| `[]byte`    | `txn.Set(k, b)`                 | `As[[]byte]`                       |
| `time.Time` | `txn.Set(k, types.FromTime(t))` | `As[time.Time]`                    |
| counter     | `txn.Set(k, types.Counter(0))`  | `As[int64]`                        |
| map         | `txn.Map(k)`                    | `As[MyStruct]`, `As[map[string]V]` |
| list        | `txn.List(k)`                   | `As[[]T]`                          |
| text        | `txn.Text(k)`                   | `As[string]`                       |

Counters differ from plain integers: concurrent increments from different peers add together. Use `txn.Increment(key, delta)` to increment.

## Persistence

`Save` writes a compact snapshot of the full document state. `SaveIncremental` writes only the changes since the last save — cheap enough to call after every `Change`.

```go
// Full snapshot — use on first write or after periodic compaction.
doc.Save(w)

// Delta — append to a file or stream to a peer.
doc.SaveIncremental(w)

// Both formats are read by the same loader.
doc2, _ := gotomerge.LoadDocument(r)
```

## Status

The core document model is complete:

- Full binary format compatibility (reads and writes `.automerge` files)
- All scalar types, nested maps, lists, and collaborative text
- Counters with concurrent increment/decrement
- Merge with last-write-wins conflict resolution
- Incremental and full-snapshot persistence

Not yet implemented: conflicts API (`GetAll`), path-based access, sync protocol.
