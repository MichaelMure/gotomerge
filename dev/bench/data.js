window.BENCHMARK_DATA = {
  "lastUpdate": 1776760611582,
  "repoUrl": "https://github.com/MichaelMure/gotomerge",
  "entries": {
    "Go Benchmark": [
      {
        "commit": {
          "author": {
            "email": "batolettre@gmail.com",
            "name": "Michael Muré",
            "username": "MichaelMure"
          },
          "committer": {
            "email": "batolettre@gmail.com",
            "name": "Michael Muré",
            "username": "MichaelMure"
          },
          "distinct": true,
          "id": "a8897014ae0d64c699e76395813a38e4700d514d",
          "message": "add GHA workflow, license",
          "timestamp": "2026-04-21T10:27:42+02:00",
          "tree_id": "f1bb57e0c0e6d0dbe9ef114fe4aec27035ef6b2c",
          "url": "https://github.com/MichaelMure/gotomerge/commit/a8897014ae0d64c699e76395813a38e4700d514d"
        },
        "date": 1776760610758,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkReadMetadata (github.com/MichaelMure/gotomerge/column)",
            "value": 179.1,
            "unit": "ns/op\t      80 B/op\t       2 allocs/op",
            "extra": "6710313 times\n4 procs"
          },
          {
            "name": "BenchmarkReadMetadata (github.com/MichaelMure/gotomerge/column) - ns/op",
            "value": 179.1,
            "unit": "ns/op",
            "extra": "6710313 times\n4 procs"
          },
          {
            "name": "BenchmarkReadMetadata (github.com/MichaelMure/gotomerge/column) - B/op",
            "value": 80,
            "unit": "B/op",
            "extra": "6710313 times\n4 procs"
          },
          {
            "name": "BenchmarkReadMetadata (github.com/MichaelMure/gotomerge/column) - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "6710313 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteMetadata (github.com/MichaelMure/gotomerge/column)",
            "value": 83.63,
            "unit": "ns/op\t      40 B/op\t       5 allocs/op",
            "extra": "14157045 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteMetadata (github.com/MichaelMure/gotomerge/column) - ns/op",
            "value": 83.63,
            "unit": "ns/op",
            "extra": "14157045 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteMetadata (github.com/MichaelMure/gotomerge/column) - B/op",
            "value": 40,
            "unit": "B/op",
            "extra": "14157045 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteMetadata (github.com/MichaelMure/gotomerge/column) - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "14157045 times\n4 procs"
          },
          {
            "name": "BenchmarkReadSpecification (github.com/MichaelMure/gotomerge/column)",
            "value": 170,
            "unit": "ns/op\t      80 B/op\t       2 allocs/op",
            "extra": "7058050 times\n4 procs"
          },
          {
            "name": "BenchmarkReadSpecification (github.com/MichaelMure/gotomerge/column) - ns/op",
            "value": 170,
            "unit": "ns/op",
            "extra": "7058050 times\n4 procs"
          },
          {
            "name": "BenchmarkReadSpecification (github.com/MichaelMure/gotomerge/column) - B/op",
            "value": 80,
            "unit": "B/op",
            "extra": "7058050 times\n4 procs"
          },
          {
            "name": "BenchmarkReadSpecification (github.com/MichaelMure/gotomerge/column) - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "7058050 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteSpecification (github.com/MichaelMure/gotomerge/column)",
            "value": 39.2,
            "unit": "ns/op\t      16 B/op\t       2 allocs/op",
            "extra": "30051612 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteSpecification (github.com/MichaelMure/gotomerge/column) - ns/op",
            "value": 39.2,
            "unit": "ns/op",
            "extra": "30051612 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteSpecification (github.com/MichaelMure/gotomerge/column) - B/op",
            "value": 16,
            "unit": "B/op",
            "extra": "30051612 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteSpecification (github.com/MichaelMure/gotomerge/column) - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "30051612 times\n4 procs"
          },
          {
            "name": "BenchmarkUint64Reader (github.com/MichaelMure/gotomerge/column/rle)",
            "value": 127.8,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "9412824 times\n4 procs"
          },
          {
            "name": "BenchmarkUint64Reader (github.com/MichaelMure/gotomerge/column/rle) - ns/op",
            "value": 127.8,
            "unit": "ns/op",
            "extra": "9412824 times\n4 procs"
          },
          {
            "name": "BenchmarkUint64Reader (github.com/MichaelMure/gotomerge/column/rle) - B/op",
            "value": 80,
            "unit": "B/op",
            "extra": "9412824 times\n4 procs"
          },
          {
            "name": "BenchmarkUint64Reader (github.com/MichaelMure/gotomerge/column/rle) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "9412824 times\n4 procs"
          },
          {
            "name": "BenchmarkInt64Reader (github.com/MichaelMure/gotomerge/column/rle)",
            "value": 121.4,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "9729748 times\n4 procs"
          },
          {
            "name": "BenchmarkInt64Reader (github.com/MichaelMure/gotomerge/column/rle) - ns/op",
            "value": 121.4,
            "unit": "ns/op",
            "extra": "9729748 times\n4 procs"
          },
          {
            "name": "BenchmarkInt64Reader (github.com/MichaelMure/gotomerge/column/rle) - B/op",
            "value": 80,
            "unit": "B/op",
            "extra": "9729748 times\n4 procs"
          },
          {
            "name": "BenchmarkInt64Reader (github.com/MichaelMure/gotomerge/column/rle) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "9729748 times\n4 procs"
          },
          {
            "name": "BenchmarkStringReader (github.com/MichaelMure/gotomerge/column/rle)",
            "value": 193.7,
            "unit": "ns/op\t     104 B/op\t       4 allocs/op",
            "extra": "6172736 times\n4 procs"
          },
          {
            "name": "BenchmarkStringReader (github.com/MichaelMure/gotomerge/column/rle) - ns/op",
            "value": 193.7,
            "unit": "ns/op",
            "extra": "6172736 times\n4 procs"
          },
          {
            "name": "BenchmarkStringReader (github.com/MichaelMure/gotomerge/column/rle) - B/op",
            "value": 104,
            "unit": "B/op",
            "extra": "6172736 times\n4 procs"
          },
          {
            "name": "BenchmarkStringReader (github.com/MichaelMure/gotomerge/column/rle) - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "6172736 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/repeated_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 45237,
            "unit": "ns/op\t   61600 B/op\t     354 allocs/op",
            "extra": "26474 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/repeated_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 45237,
            "unit": "ns/op",
            "extra": "26474 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/repeated_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 61600,
            "unit": "B/op",
            "extra": "26474 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/repeated_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 354,
            "unit": "allocs/op",
            "extra": "26474 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/repeated_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 613681,
            "unit": "ns/op\t  708237 B/op\t    2967 allocs/op",
            "extra": "1789 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/repeated_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 613681,
            "unit": "ns/op",
            "extra": "1789 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/repeated_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 708237,
            "unit": "B/op",
            "extra": "1789 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/repeated_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 2967,
            "unit": "allocs/op",
            "extra": "1789 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/repeated_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 9146489,
            "unit": "ns/op\t 9102222 B/op\t   30069 allocs/op",
            "extra": "130 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/repeated_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 9146489,
            "unit": "ns/op",
            "extra": "130 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/repeated_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 9102222,
            "unit": "B/op",
            "extra": "130 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/repeated_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 30069,
            "unit": "allocs/op",
            "extra": "130 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/increasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 150461,
            "unit": "ns/op\t  161429 B/op\t     479 allocs/op",
            "extra": "8862 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/increasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 150461,
            "unit": "ns/op",
            "extra": "8862 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/increasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 161429,
            "unit": "B/op",
            "extra": "8862 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/increasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 479,
            "unit": "allocs/op",
            "extra": "8862 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/increasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 965350,
            "unit": "ns/op\t 1100746 B/op\t    4946 allocs/op",
            "extra": "1194 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/increasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 965350,
            "unit": "ns/op",
            "extra": "1194 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/increasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 1100746,
            "unit": "B/op",
            "extra": "1194 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/increasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 4946,
            "unit": "allocs/op",
            "extra": "1194 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/increasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 13527519,
            "unit": "ns/op\t11785149 B/op\t   50067 allocs/op",
            "extra": "92 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/increasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 13527519,
            "unit": "ns/op",
            "extra": "92 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/increasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 11785149,
            "unit": "B/op",
            "extra": "92 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/increasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 50067,
            "unit": "allocs/op",
            "extra": "92 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/decreasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 156429,
            "unit": "ns/op\t  166900 B/op\t     479 allocs/op",
            "extra": "7461 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/decreasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 156429,
            "unit": "ns/op",
            "extra": "7461 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/decreasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 166900,
            "unit": "B/op",
            "extra": "7461 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/decreasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 479,
            "unit": "allocs/op",
            "extra": "7461 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/decreasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 988822,
            "unit": "ns/op\t 1122948 B/op\t    4947 allocs/op",
            "extra": "1215 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/decreasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 988822,
            "unit": "ns/op",
            "extra": "1215 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/decreasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 1122948,
            "unit": "B/op",
            "extra": "1215 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/decreasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 4947,
            "unit": "allocs/op",
            "extra": "1215 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/decreasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 13821879,
            "unit": "ns/op\t11785286 B/op\t   50072 allocs/op",
            "extra": "88 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/decreasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 13821879,
            "unit": "ns/op",
            "extra": "88 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/decreasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 11785286,
            "unit": "B/op",
            "extra": "88 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/decreasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 50072,
            "unit": "allocs/op",
            "extra": "88 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/repeated_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 25249,
            "unit": "ns/op\t   9.51 MB/s\t    5952 B/op\t     167 allocs/op",
            "extra": "47276 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/repeated_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 25249,
            "unit": "ns/op",
            "extra": "47276 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/repeated_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 9.51,
            "unit": "MB/s",
            "extra": "47276 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/repeated_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 5952,
            "unit": "B/op",
            "extra": "47276 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/repeated_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 167,
            "unit": "allocs/op",
            "extra": "47276 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/repeated_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 257048,
            "unit": "ns/op\t   6.35 MB/s\t  221202 B/op\t    1079 allocs/op",
            "extra": "4502 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/repeated_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 257048,
            "unit": "ns/op",
            "extra": "4502 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/repeated_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 6.35,
            "unit": "MB/s",
            "extra": "4502 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/repeated_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 221202,
            "unit": "B/op",
            "extra": "4502 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/repeated_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 1079,
            "unit": "allocs/op",
            "extra": "4502 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/repeated_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 1690476,
            "unit": "ns/op\t  11.20 MB/s\t 1065116 B/op\t   10100 allocs/op",
            "extra": "691 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/repeated_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 1690476,
            "unit": "ns/op",
            "extra": "691 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/repeated_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 11.2,
            "unit": "MB/s",
            "extra": "691 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/repeated_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 1065116,
            "unit": "B/op",
            "extra": "691 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/repeated_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 10100,
            "unit": "allocs/op",
            "extra": "691 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/increasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 80530,
            "unit": "ns/op\t   4.77 MB/s\t   68475 B/op\t     181 allocs/op",
            "extra": "14908 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/increasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 80530,
            "unit": "ns/op",
            "extra": "14908 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/increasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 4.77,
            "unit": "MB/s",
            "extra": "14908 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/increasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 68475,
            "unit": "B/op",
            "extra": "14908 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/increasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 181,
            "unit": "allocs/op",
            "extra": "14908 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/increasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 375952,
            "unit": "ns/op\t   8.91 MB/s\t  439336 B/op\t    1104 allocs/op",
            "extra": "2944 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/increasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 375952,
            "unit": "ns/op",
            "extra": "2944 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/increasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 8.91,
            "unit": "MB/s",
            "extra": "2944 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/increasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 439336,
            "unit": "B/op",
            "extra": "2944 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/increasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 1104,
            "unit": "allocs/op",
            "extra": "2944 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/increasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 2936430,
            "unit": "ns/op\t  12.30 MB/s\t 2032541 B/op\t   10137 allocs/op",
            "extra": "415 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/increasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 2936430,
            "unit": "ns/op",
            "extra": "415 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/increasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 12.3,
            "unit": "MB/s",
            "extra": "415 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/increasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 2032541,
            "unit": "B/op",
            "extra": "415 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/increasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 10137,
            "unit": "allocs/op",
            "extra": "415 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/decreasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 83096,
            "unit": "ns/op\t   4.62 MB/s\t   71598 B/op\t     181 allocs/op",
            "extra": "14274 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/decreasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 83096,
            "unit": "ns/op",
            "extra": "14274 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/decreasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 4.62,
            "unit": "MB/s",
            "extra": "14274 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/decreasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 71598,
            "unit": "B/op",
            "extra": "14274 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/decreasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 181,
            "unit": "allocs/op",
            "extra": "14274 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/decreasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 386717,
            "unit": "ns/op\t   8.67 MB/s\t  463115 B/op\t    1104 allocs/op",
            "extra": "3286 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/decreasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 386717,
            "unit": "ns/op",
            "extra": "3286 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/decreasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 8.67,
            "unit": "MB/s",
            "extra": "3286 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/decreasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 463115,
            "unit": "B/op",
            "extra": "3286 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/decreasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 1104,
            "unit": "allocs/op",
            "extra": "3286 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/decreasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 3017199,
            "unit": "ns/op\t  11.87 MB/s\t 2032567 B/op\t   10137 allocs/op",
            "extra": "399 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/decreasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 3017199,
            "unit": "ns/op",
            "extra": "399 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/decreasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 11.87,
            "unit": "MB/s",
            "extra": "399 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/decreasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 2032567,
            "unit": "B/op",
            "extra": "399 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/decreasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 10137,
            "unit": "allocs/op",
            "extra": "399 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/repeated_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 30649,
            "unit": "ns/op\t   7.83 MB/s\t   13048 B/op\t     208 allocs/op",
            "extra": "39234 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/repeated_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 30649,
            "unit": "ns/op",
            "extra": "39234 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/repeated_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 7.83,
            "unit": "MB/s",
            "extra": "39234 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/repeated_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 13048,
            "unit": "B/op",
            "extra": "39234 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/repeated_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 208,
            "unit": "allocs/op",
            "extra": "39234 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/repeated_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 336823,
            "unit": "ns/op\t   4.84 MB/s\t  166360 B/op\t    2062 allocs/op",
            "extra": "3508 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/repeated_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 336823,
            "unit": "ns/op",
            "extra": "3508 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/repeated_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 4.84,
            "unit": "MB/s",
            "extra": "3508 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/repeated_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 166360,
            "unit": "B/op",
            "extra": "3508 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/repeated_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 2062,
            "unit": "allocs/op",
            "extra": "3508 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/repeated_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 3367689,
            "unit": "ns/op\t   5.62 MB/s\t 1205668 B/op\t   21690 allocs/op",
            "extra": "366 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/repeated_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 3367689,
            "unit": "ns/op",
            "extra": "366 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/repeated_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 5.62,
            "unit": "MB/s",
            "extra": "366 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/repeated_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 1205668,
            "unit": "B/op",
            "extra": "366 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/repeated_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 21690,
            "unit": "allocs/op",
            "extra": "366 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/increasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 59370,
            "unit": "ns/op\t   6.47 MB/s\t   65024 B/op\t     506 allocs/op",
            "extra": "19512 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/increasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 59370,
            "unit": "ns/op",
            "extra": "19512 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/increasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 6.47,
            "unit": "MB/s",
            "extra": "19512 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/increasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 65024,
            "unit": "B/op",
            "extra": "19512 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/increasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 506,
            "unit": "allocs/op",
            "extra": "19512 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/increasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 573607,
            "unit": "ns/op\t   5.84 MB/s\t  413073 B/op\t    5077 allocs/op",
            "extra": "2084 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/increasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 573607,
            "unit": "ns/op",
            "extra": "2084 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/increasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 5.84,
            "unit": "MB/s",
            "extra": "2084 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/increasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 413073,
            "unit": "B/op",
            "extra": "2084 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/increasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 5077,
            "unit": "allocs/op",
            "extra": "2084 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/increasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 5410974,
            "unit": "ns/op\t   6.67 MB/s\t 3003308 B/op\t   51767 allocs/op",
            "extra": "225 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/increasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 5410974,
            "unit": "ns/op",
            "extra": "225 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/increasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 6.67,
            "unit": "MB/s",
            "extra": "225 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/increasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 3003308,
            "unit": "B/op",
            "extra": "225 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/increasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 51767,
            "unit": "allocs/op",
            "extra": "225 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/decreasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 56455,
            "unit": "ns/op\t   6.80 MB/s\t   65024 B/op\t     506 allocs/op",
            "extra": "21206 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/decreasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 56455,
            "unit": "ns/op",
            "extra": "21206 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/decreasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 6.8,
            "unit": "MB/s",
            "extra": "21206 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/decreasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 65024,
            "unit": "B/op",
            "extra": "21206 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/decreasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 506,
            "unit": "allocs/op",
            "extra": "21206 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/decreasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 538181,
            "unit": "ns/op\t   6.23 MB/s\t  413073 B/op\t    5077 allocs/op",
            "extra": "2190 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/decreasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 538181,
            "unit": "ns/op",
            "extra": "2190 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/decreasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 6.23,
            "unit": "MB/s",
            "extra": "2190 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/decreasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 413073,
            "unit": "B/op",
            "extra": "2190 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/decreasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 5077,
            "unit": "allocs/op",
            "extra": "2190 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/decreasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 5358743,
            "unit": "ns/op\t   6.68 MB/s\t 3003508 B/op\t   51767 allocs/op",
            "extra": "224 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/decreasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 5358743,
            "unit": "ns/op",
            "extra": "224 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/decreasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 6.68,
            "unit": "MB/s",
            "extra": "224 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/decreasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 3003508,
            "unit": "B/op",
            "extra": "224 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/decreasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 51767,
            "unit": "allocs/op",
            "extra": "224 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/repeated_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 32525,
            "unit": "ns/op\t   6.15 MB/s\t   33984 B/op\t     197 allocs/op",
            "extra": "36981 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/repeated_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 32525,
            "unit": "ns/op",
            "extra": "36981 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/repeated_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 6.15,
            "unit": "MB/s",
            "extra": "36981 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/repeated_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 33984,
            "unit": "B/op",
            "extra": "36981 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/repeated_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 197,
            "unit": "allocs/op",
            "extra": "36981 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/repeated_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 330466,
            "unit": "ns/op\t   4.87 MB/s\t  331506 B/op\t    1899 allocs/op",
            "extra": "3577 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/repeated_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 330466,
            "unit": "ns/op",
            "extra": "3577 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/repeated_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 4.87,
            "unit": "MB/s",
            "extra": "3577 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/repeated_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 331506,
            "unit": "B/op",
            "extra": "3577 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/repeated_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 1899,
            "unit": "allocs/op",
            "extra": "3577 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/repeated_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 4809312,
            "unit": "ns/op\t   3.93 MB/s\t 4492032 B/op\t   19978 allocs/op",
            "extra": "240 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/repeated_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 4809312,
            "unit": "ns/op",
            "extra": "240 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/repeated_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 3.93,
            "unit": "MB/s",
            "extra": "240 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/repeated_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 4492032,
            "unit": "B/op",
            "extra": "240 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/repeated_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 19978,
            "unit": "allocs/op",
            "extra": "240 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/increasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 60556,
            "unit": "ns/op\t   6.14 MB/s\t   86088 B/op\t     496 allocs/op",
            "extra": "19790 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/increasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 60556,
            "unit": "ns/op",
            "extra": "19790 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/increasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 6.14,
            "unit": "MB/s",
            "extra": "19790 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/increasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 86088,
            "unit": "B/op",
            "extra": "19790 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/increasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 496,
            "unit": "allocs/op",
            "extra": "19790 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/increasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 546015,
            "unit": "ns/op\t   7.09 MB/s\t  542635 B/op\t    4950 allocs/op",
            "extra": "2121 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/increasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 546015,
            "unit": "ns/op",
            "extra": "2121 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/increasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 7.09,
            "unit": "MB/s",
            "extra": "2121 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/increasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 542635,
            "unit": "B/op",
            "extra": "2121 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/increasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 4950,
            "unit": "allocs/op",
            "extra": "2121 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/increasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 7411879,
            "unit": "ns/op\t   5.11 MB/s\t 6217730 B/op\t   50038 allocs/op",
            "extra": "160 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/increasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 7411879,
            "unit": "ns/op",
            "extra": "160 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/increasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 5.11,
            "unit": "MB/s",
            "extra": "160 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/increasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 6217730,
            "unit": "B/op",
            "extra": "160 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/increasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 50038,
            "unit": "allocs/op",
            "extra": "160 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/decreasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 62145,
            "unit": "ns/op\t   6.05 MB/s\t   86088 B/op\t     496 allocs/op",
            "extra": "19155 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/decreasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 62145,
            "unit": "ns/op",
            "extra": "19155 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/decreasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 6.05,
            "unit": "MB/s",
            "extra": "19155 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/decreasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 86088,
            "unit": "B/op",
            "extra": "19155 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/decreasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 496,
            "unit": "allocs/op",
            "extra": "19155 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/decreasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 555394,
            "unit": "ns/op\t   6.96 MB/s\t  542634 B/op\t    4950 allocs/op",
            "extra": "2172 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/decreasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 555394,
            "unit": "ns/op",
            "extra": "2172 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/decreasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 6.96,
            "unit": "MB/s",
            "extra": "2172 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/decreasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 542634,
            "unit": "B/op",
            "extra": "2172 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/decreasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 4950,
            "unit": "allocs/op",
            "extra": "2172 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/decreasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 7473688,
            "unit": "ns/op\t   4.93 MB/s\t 6217996 B/op\t   50046 allocs/op",
            "extra": "159 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/decreasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 7473688,
            "unit": "ns/op",
            "extra": "159 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/decreasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 4.93,
            "unit": "MB/s",
            "extra": "159 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/decreasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 6217996,
            "unit": "B/op",
            "extra": "159 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/decreasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 50046,
            "unit": "allocs/op",
            "extra": "159 times\n4 procs"
          },
          {
            "name": "BenchmarkSaveLoad/big_paste (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 53088,
            "unit": "ns/op\t   75236 B/op\t     162 allocs/op",
            "extra": "21908 times\n4 procs"
          },
          {
            "name": "BenchmarkSaveLoad/big_paste (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 53088,
            "unit": "ns/op",
            "extra": "21908 times\n4 procs"
          },
          {
            "name": "BenchmarkSaveLoad/big_paste (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 75236,
            "unit": "B/op",
            "extra": "21908 times\n4 procs"
          },
          {
            "name": "BenchmarkSaveLoad/big_paste (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 162,
            "unit": "allocs/op",
            "extra": "21908 times\n4 procs"
          },
          {
            "name": "BenchmarkSaveLoad/maps_in_maps (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 1800526,
            "unit": "ns/op\t 1335943 B/op\t    6493 allocs/op",
            "extra": "657 times\n4 procs"
          },
          {
            "name": "BenchmarkSaveLoad/maps_in_maps (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 1800526,
            "unit": "ns/op",
            "extra": "657 times\n4 procs"
          },
          {
            "name": "BenchmarkSaveLoad/maps_in_maps (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 1335943,
            "unit": "B/op",
            "extra": "657 times\n4 procs"
          },
          {
            "name": "BenchmarkSaveLoad/maps_in_maps (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 6493,
            "unit": "allocs/op",
            "extra": "657 times\n4 procs"
          },
          {
            "name": "BenchmarkSaveLoad/deep_history (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 1516832,
            "unit": "ns/op\t  803277 B/op\t   10133 allocs/op",
            "extra": "783 times\n4 procs"
          },
          {
            "name": "BenchmarkSaveLoad/deep_history (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 1516832,
            "unit": "ns/op",
            "extra": "783 times\n4 procs"
          },
          {
            "name": "BenchmarkSaveLoad/deep_history (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 803277,
            "unit": "B/op",
            "extra": "783 times\n4 procs"
          },
          {
            "name": "BenchmarkSaveLoad/deep_history (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 10133,
            "unit": "allocs/op",
            "extra": "783 times\n4 procs"
          },
          {
            "name": "BenchmarkSaveLoad/poorly_simulated_typing (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 532917,
            "unit": "ns/op\t  270165 B/op\t    4434 allocs/op",
            "extra": "2214 times\n4 procs"
          },
          {
            "name": "BenchmarkSaveLoad/poorly_simulated_typing (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 532917,
            "unit": "ns/op",
            "extra": "2214 times\n4 procs"
          },
          {
            "name": "BenchmarkSaveLoad/poorly_simulated_typing (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 270165,
            "unit": "B/op",
            "extra": "2214 times\n4 procs"
          },
          {
            "name": "BenchmarkSaveLoad/poorly_simulated_typing (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 4434,
            "unit": "allocs/op",
            "extra": "2214 times\n4 procs"
          },
          {
            "name": "BenchmarkEditTrace/apply (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 981412656,
            "unit": "ns/op\t  29.80 MB/s\t1072416760 B/op\t12966370 allocs/op",
            "extra": "2 times\n4 procs"
          },
          {
            "name": "BenchmarkEditTrace/apply (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 981412656,
            "unit": "ns/op",
            "extra": "2 times\n4 procs"
          },
          {
            "name": "BenchmarkEditTrace/apply (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 29.8,
            "unit": "MB/s",
            "extra": "2 times\n4 procs"
          },
          {
            "name": "BenchmarkEditTrace/apply (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 1072416760,
            "unit": "B/op",
            "extra": "2 times\n4 procs"
          },
          {
            "name": "BenchmarkEditTrace/apply (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 12966370,
            "unit": "allocs/op",
            "extra": "2 times\n4 procs"
          },
          {
            "name": "BenchmarkEditTrace/save (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 1106540079,
            "unit": "ns/op\t   0.13 MB/s\t1110107696 B/op\t13346988 allocs/op",
            "extra": "1 times\n4 procs"
          },
          {
            "name": "BenchmarkEditTrace/save (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 1106540079,
            "unit": "ns/op",
            "extra": "1 times\n4 procs"
          },
          {
            "name": "BenchmarkEditTrace/save (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 0.13,
            "unit": "MB/s",
            "extra": "1 times\n4 procs"
          },
          {
            "name": "BenchmarkEditTrace/save (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 1110107696,
            "unit": "B/op",
            "extra": "1 times\n4 procs"
          },
          {
            "name": "BenchmarkEditTrace/save (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 13346988,
            "unit": "allocs/op",
            "extra": "1 times\n4 procs"
          },
          {
            "name": "BenchmarkEditTrace/load (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 87332079,
            "unit": "ns/op\t   1.62 MB/s\t31288012 B/op\t  761472 allocs/op",
            "extra": "13 times\n4 procs"
          },
          {
            "name": "BenchmarkEditTrace/load (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 87332079,
            "unit": "ns/op",
            "extra": "13 times\n4 procs"
          },
          {
            "name": "BenchmarkEditTrace/load (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 1.62,
            "unit": "MB/s",
            "extra": "13 times\n4 procs"
          },
          {
            "name": "BenchmarkEditTrace/load (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 31288012,
            "unit": "B/op",
            "extra": "13 times\n4 procs"
          },
          {
            "name": "BenchmarkEditTrace/load (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 761472,
            "unit": "allocs/op",
            "extra": "13 times\n4 procs"
          },
          {
            "name": "BenchmarkRange (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 38264636,
            "unit": "ns/op\t48206950 B/op\t  100559 allocs/op",
            "extra": "36 times\n4 procs"
          },
          {
            "name": "BenchmarkRange (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 38264636,
            "unit": "ns/op",
            "extra": "36 times\n4 procs"
          },
          {
            "name": "BenchmarkRange (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 48206950,
            "unit": "B/op",
            "extra": "36 times\n4 procs"
          },
          {
            "name": "BenchmarkRange (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 100559,
            "unit": "allocs/op",
            "extra": "36 times\n4 procs"
          },
          {
            "name": "BenchmarkGet/scalar/keys=100 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 94.21,
            "unit": "ns/op\t      48 B/op\t       1 allocs/op",
            "extra": "12831223 times\n4 procs"
          },
          {
            "name": "BenchmarkGet/scalar/keys=100 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 94.21,
            "unit": "ns/op",
            "extra": "12831223 times\n4 procs"
          },
          {
            "name": "BenchmarkGet/scalar/keys=100 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 48,
            "unit": "B/op",
            "extra": "12831223 times\n4 procs"
          },
          {
            "name": "BenchmarkGet/scalar/keys=100 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "12831223 times\n4 procs"
          },
          {
            "name": "BenchmarkGet/scalar/keys=1000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 97.26,
            "unit": "ns/op\t      48 B/op\t       1 allocs/op",
            "extra": "12755914 times\n4 procs"
          },
          {
            "name": "BenchmarkGet/scalar/keys=1000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 97.26,
            "unit": "ns/op",
            "extra": "12755914 times\n4 procs"
          },
          {
            "name": "BenchmarkGet/scalar/keys=1000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 48,
            "unit": "B/op",
            "extra": "12755914 times\n4 procs"
          },
          {
            "name": "BenchmarkGet/scalar/keys=1000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "12755914 times\n4 procs"
          },
          {
            "name": "BenchmarkGet/scalar/keys=10000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 101,
            "unit": "ns/op\t      48 B/op\t       1 allocs/op",
            "extra": "12981700 times\n4 procs"
          },
          {
            "name": "BenchmarkGet/scalar/keys=10000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 101,
            "unit": "ns/op",
            "extra": "12981700 times\n4 procs"
          },
          {
            "name": "BenchmarkGet/scalar/keys=10000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 48,
            "unit": "B/op",
            "extra": "12981700 times\n4 procs"
          },
          {
            "name": "BenchmarkGet/scalar/keys=10000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "12981700 times\n4 procs"
          },
          {
            "name": "BenchmarkGet/nested/keys=100 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 90.42,
            "unit": "ns/op\t      48 B/op\t       1 allocs/op",
            "extra": "13352162 times\n4 procs"
          },
          {
            "name": "BenchmarkGet/nested/keys=100 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 90.42,
            "unit": "ns/op",
            "extra": "13352162 times\n4 procs"
          },
          {
            "name": "BenchmarkGet/nested/keys=100 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 48,
            "unit": "B/op",
            "extra": "13352162 times\n4 procs"
          },
          {
            "name": "BenchmarkGet/nested/keys=100 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "13352162 times\n4 procs"
          },
          {
            "name": "BenchmarkGet/nested/keys=1000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 97.51,
            "unit": "ns/op\t      48 B/op\t       1 allocs/op",
            "extra": "12724249 times\n4 procs"
          },
          {
            "name": "BenchmarkGet/nested/keys=1000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 97.51,
            "unit": "ns/op",
            "extra": "12724249 times\n4 procs"
          },
          {
            "name": "BenchmarkGet/nested/keys=1000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 48,
            "unit": "B/op",
            "extra": "12724249 times\n4 procs"
          },
          {
            "name": "BenchmarkGet/nested/keys=1000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "12724249 times\n4 procs"
          },
          {
            "name": "BenchmarkGet/nested/keys=10000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 93.66,
            "unit": "ns/op\t      48 B/op\t       1 allocs/op",
            "extra": "13151428 times\n4 procs"
          },
          {
            "name": "BenchmarkGet/nested/keys=10000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 93.66,
            "unit": "ns/op",
            "extra": "13151428 times\n4 procs"
          },
          {
            "name": "BenchmarkGet/nested/keys=10000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 48,
            "unit": "B/op",
            "extra": "13151428 times\n4 procs"
          },
          {
            "name": "BenchmarkGet/nested/keys=10000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "13151428 times\n4 procs"
          },
          {
            "name": "BenchmarkAs/int64 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 150.6,
            "unit": "ns/op\t      64 B/op\t       3 allocs/op",
            "extra": "7914273 times\n4 procs"
          },
          {
            "name": "BenchmarkAs/int64 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 150.6,
            "unit": "ns/op",
            "extra": "7914273 times\n4 procs"
          },
          {
            "name": "BenchmarkAs/int64 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 64,
            "unit": "B/op",
            "extra": "7914273 times\n4 procs"
          },
          {
            "name": "BenchmarkAs/int64 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "7914273 times\n4 procs"
          },
          {
            "name": "BenchmarkAs/string (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 120.3,
            "unit": "ns/op\t      80 B/op\t       2 allocs/op",
            "extra": "10115610 times\n4 procs"
          },
          {
            "name": "BenchmarkAs/string (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 120.3,
            "unit": "ns/op",
            "extra": "10115610 times\n4 procs"
          },
          {
            "name": "BenchmarkAs/string (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 80,
            "unit": "B/op",
            "extra": "10115610 times\n4 procs"
          },
          {
            "name": "BenchmarkAs/string (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "10115610 times\n4 procs"
          },
          {
            "name": "BenchmarkAs/struct (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 863.5,
            "unit": "ns/op\t     320 B/op\t      14 allocs/op",
            "extra": "1377698 times\n4 procs"
          },
          {
            "name": "BenchmarkAs/struct (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 863.5,
            "unit": "ns/op",
            "extra": "1377698 times\n4 procs"
          },
          {
            "name": "BenchmarkAs/struct (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 320,
            "unit": "B/op",
            "extra": "1377698 times\n4 procs"
          },
          {
            "name": "BenchmarkAs/struct (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 14,
            "unit": "allocs/op",
            "extra": "1377698 times\n4 procs"
          },
          {
            "name": "BenchmarkAs/slice_string/len=100 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 16665,
            "unit": "ns/op\t   20432 B/op\t     304 allocs/op",
            "extra": "67425 times\n4 procs"
          },
          {
            "name": "BenchmarkAs/slice_string/len=100 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 16665,
            "unit": "ns/op",
            "extra": "67425 times\n4 procs"
          },
          {
            "name": "BenchmarkAs/slice_string/len=100 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 20432,
            "unit": "B/op",
            "extra": "67425 times\n4 procs"
          },
          {
            "name": "BenchmarkAs/slice_string/len=100 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 304,
            "unit": "allocs/op",
            "extra": "67425 times\n4 procs"
          },
          {
            "name": "BenchmarkMerge/changes=10 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 43591,
            "unit": "ns/op\t   33808 B/op\t     442 allocs/op",
            "extra": "27682 times\n4 procs"
          },
          {
            "name": "BenchmarkMerge/changes=10 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 43591,
            "unit": "ns/op",
            "extra": "27682 times\n4 procs"
          },
          {
            "name": "BenchmarkMerge/changes=10 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 33808,
            "unit": "B/op",
            "extra": "27682 times\n4 procs"
          },
          {
            "name": "BenchmarkMerge/changes=10 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 442,
            "unit": "allocs/op",
            "extra": "27682 times\n4 procs"
          },
          {
            "name": "BenchmarkMerge/changes=100 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 312811,
            "unit": "ns/op\t  320624 B/op\t    4312 allocs/op",
            "extra": "4077 times\n4 procs"
          },
          {
            "name": "BenchmarkMerge/changes=100 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 312811,
            "unit": "ns/op",
            "extra": "4077 times\n4 procs"
          },
          {
            "name": "BenchmarkMerge/changes=100 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 320624,
            "unit": "B/op",
            "extra": "4077 times\n4 procs"
          },
          {
            "name": "BenchmarkMerge/changes=100 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 4312,
            "unit": "allocs/op",
            "extra": "4077 times\n4 procs"
          },
          {
            "name": "BenchmarkMerge/changes=1000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 3239802,
            "unit": "ns/op\t 3533105 B/op\t   43780 allocs/op",
            "extra": "370 times\n4 procs"
          },
          {
            "name": "BenchmarkMerge/changes=1000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 3239802,
            "unit": "ns/op",
            "extra": "370 times\n4 procs"
          },
          {
            "name": "BenchmarkMerge/changes=1000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 3533105,
            "unit": "B/op",
            "extra": "370 times\n4 procs"
          },
          {
            "name": "BenchmarkMerge/changes=1000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 43780,
            "unit": "allocs/op",
            "extra": "370 times\n4 procs"
          },
          {
            "name": "BenchmarkList/build/len=100 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 67674,
            "unit": "ns/op\t   78912 B/op\t     488 allocs/op",
            "extra": "17631 times\n4 procs"
          },
          {
            "name": "BenchmarkList/build/len=100 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 67674,
            "unit": "ns/op",
            "extra": "17631 times\n4 procs"
          },
          {
            "name": "BenchmarkList/build/len=100 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 78912,
            "unit": "B/op",
            "extra": "17631 times\n4 procs"
          },
          {
            "name": "BenchmarkList/build/len=100 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 488,
            "unit": "allocs/op",
            "extra": "17631 times\n4 procs"
          },
          {
            "name": "BenchmarkList/build/len=1000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 784359,
            "unit": "ns/op\t  721079 B/op\t    4007 allocs/op",
            "extra": "1621 times\n4 procs"
          },
          {
            "name": "BenchmarkList/build/len=1000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 784359,
            "unit": "ns/op",
            "extra": "1621 times\n4 procs"
          },
          {
            "name": "BenchmarkList/build/len=1000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 721079,
            "unit": "B/op",
            "extra": "1621 times\n4 procs"
          },
          {
            "name": "BenchmarkList/build/len=1000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 4007,
            "unit": "allocs/op",
            "extra": "1621 times\n4 procs"
          },
          {
            "name": "BenchmarkList/build/len=10000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 8022561,
            "unit": "ns/op\t10017465 B/op\t   40146 allocs/op",
            "extra": "147 times\n4 procs"
          },
          {
            "name": "BenchmarkList/build/len=10000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 8022561,
            "unit": "ns/op",
            "extra": "147 times\n4 procs"
          },
          {
            "name": "BenchmarkList/build/len=10000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 10017465,
            "unit": "B/op",
            "extra": "147 times\n4 procs"
          },
          {
            "name": "BenchmarkList/build/len=10000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 40146,
            "unit": "allocs/op",
            "extra": "147 times\n4 procs"
          },
          {
            "name": "BenchmarkList/save/len=100 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 26134,
            "unit": "ns/op\t  10.75 MB/s\t    6976 B/op\t     179 allocs/op",
            "extra": "49983 times\n4 procs"
          },
          {
            "name": "BenchmarkList/save/len=100 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 26134,
            "unit": "ns/op",
            "extra": "49983 times\n4 procs"
          },
          {
            "name": "BenchmarkList/save/len=100 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 10.75,
            "unit": "MB/s",
            "extra": "49983 times\n4 procs"
          },
          {
            "name": "BenchmarkList/save/len=100 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 6976,
            "unit": "B/op",
            "extra": "49983 times\n4 procs"
          },
          {
            "name": "BenchmarkList/save/len=100 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 179,
            "unit": "allocs/op",
            "extra": "49983 times\n4 procs"
          },
          {
            "name": "BenchmarkList/save/len=1000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 192012,
            "unit": "ns/op\t   8.71 MB/s\t   29214 B/op\t    1087 allocs/op",
            "extra": "6230 times\n4 procs"
          },
          {
            "name": "BenchmarkList/save/len=1000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 192012,
            "unit": "ns/op",
            "extra": "6230 times\n4 procs"
          },
          {
            "name": "BenchmarkList/save/len=1000 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 8.71,
            "unit": "MB/s",
            "extra": "6230 times\n4 procs"
          },
          {
            "name": "BenchmarkList/save/len=1000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 29214,
            "unit": "B/op",
            "extra": "6230 times\n4 procs"
          },
          {
            "name": "BenchmarkList/save/len=1000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 1087,
            "unit": "allocs/op",
            "extra": "6230 times\n4 procs"
          },
          {
            "name": "BenchmarkList/save/len=10000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 1354653,
            "unit": "ns/op\t  14.01 MB/s\t  310147 B/op\t   10096 allocs/op",
            "extra": "878 times\n4 procs"
          },
          {
            "name": "BenchmarkList/save/len=10000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 1354653,
            "unit": "ns/op",
            "extra": "878 times\n4 procs"
          },
          {
            "name": "BenchmarkList/save/len=10000 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 14.01,
            "unit": "MB/s",
            "extra": "878 times\n4 procs"
          },
          {
            "name": "BenchmarkList/save/len=10000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 310147,
            "unit": "B/op",
            "extra": "878 times\n4 procs"
          },
          {
            "name": "BenchmarkList/save/len=10000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 10096,
            "unit": "allocs/op",
            "extra": "878 times\n4 procs"
          },
          {
            "name": "BenchmarkList/load/len=100 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 29649,
            "unit": "ns/op\t   9.48 MB/s\t   12336 B/op\t     221 allocs/op",
            "extra": "39613 times\n4 procs"
          },
          {
            "name": "BenchmarkList/load/len=100 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 29649,
            "unit": "ns/op",
            "extra": "39613 times\n4 procs"
          },
          {
            "name": "BenchmarkList/load/len=100 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 9.48,
            "unit": "MB/s",
            "extra": "39613 times\n4 procs"
          },
          {
            "name": "BenchmarkList/load/len=100 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 12336,
            "unit": "B/op",
            "extra": "39613 times\n4 procs"
          },
          {
            "name": "BenchmarkList/load/len=100 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 221,
            "unit": "allocs/op",
            "extra": "39613 times\n4 procs"
          },
          {
            "name": "BenchmarkList/load/len=1000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 285385,
            "unit": "ns/op\t   5.86 MB/s\t  149232 B/op\t    2116 allocs/op",
            "extra": "4256 times\n4 procs"
          },
          {
            "name": "BenchmarkList/load/len=1000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 285385,
            "unit": "ns/op",
            "extra": "4256 times\n4 procs"
          },
          {
            "name": "BenchmarkList/load/len=1000 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 5.86,
            "unit": "MB/s",
            "extra": "4256 times\n4 procs"
          },
          {
            "name": "BenchmarkList/load/len=1000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 149232,
            "unit": "B/op",
            "extra": "4256 times\n4 procs"
          },
          {
            "name": "BenchmarkList/load/len=1000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 2116,
            "unit": "allocs/op",
            "extra": "4256 times\n4 procs"
          },
          {
            "name": "BenchmarkList/load/len=10000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 2736284,
            "unit": "ns/op\t   6.94 MB/s\t 1016936 B/op\t   22165 allocs/op",
            "extra": "446 times\n4 procs"
          },
          {
            "name": "BenchmarkList/load/len=10000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 2736284,
            "unit": "ns/op",
            "extra": "446 times\n4 procs"
          },
          {
            "name": "BenchmarkList/load/len=10000 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 6.94,
            "unit": "MB/s",
            "extra": "446 times\n4 procs"
          },
          {
            "name": "BenchmarkList/load/len=10000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 1016936,
            "unit": "B/op",
            "extra": "446 times\n4 procs"
          },
          {
            "name": "BenchmarkList/load/len=10000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 22165,
            "unit": "allocs/op",
            "extra": "446 times\n4 procs"
          },
          {
            "name": "BenchmarkList/iterate/len=100 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 7095,
            "unit": "ns/op\t   12992 B/op\t     101 allocs/op",
            "extra": "148599 times\n4 procs"
          },
          {
            "name": "BenchmarkList/iterate/len=100 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 7095,
            "unit": "ns/op",
            "extra": "148599 times\n4 procs"
          },
          {
            "name": "BenchmarkList/iterate/len=100 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 12992,
            "unit": "B/op",
            "extra": "148599 times\n4 procs"
          },
          {
            "name": "BenchmarkList/iterate/len=100 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 101,
            "unit": "allocs/op",
            "extra": "148599 times\n4 procs"
          },
          {
            "name": "BenchmarkList/iterate/len=1000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 71756,
            "unit": "ns/op\t  121794 B/op\t    1001 allocs/op",
            "extra": "16980 times\n4 procs"
          },
          {
            "name": "BenchmarkList/iterate/len=1000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 71756,
            "unit": "ns/op",
            "extra": "16980 times\n4 procs"
          },
          {
            "name": "BenchmarkList/iterate/len=1000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 121794,
            "unit": "B/op",
            "extra": "16980 times\n4 procs"
          },
          {
            "name": "BenchmarkList/iterate/len=1000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 1001,
            "unit": "allocs/op",
            "extra": "16980 times\n4 procs"
          },
          {
            "name": "BenchmarkList/iterate/len=10000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 739124,
            "unit": "ns/op\t 1206681 B/op\t   10019 allocs/op",
            "extra": "1620 times\n4 procs"
          },
          {
            "name": "BenchmarkList/iterate/len=10000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 739124,
            "unit": "ns/op",
            "extra": "1620 times\n4 procs"
          },
          {
            "name": "BenchmarkList/iterate/len=10000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 1206681,
            "unit": "B/op",
            "extra": "1620 times\n4 procs"
          },
          {
            "name": "BenchmarkList/iterate/len=10000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 10019,
            "unit": "allocs/op",
            "extra": "1620 times\n4 procs"
          },
          {
            "name": "BenchmarkList/as_slice/len=1000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 186004,
            "unit": "ns/op\t  170148 B/op\t    4004 allocs/op",
            "extra": "6272 times\n4 procs"
          },
          {
            "name": "BenchmarkList/as_slice/len=1000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 186004,
            "unit": "ns/op",
            "extra": "6272 times\n4 procs"
          },
          {
            "name": "BenchmarkList/as_slice/len=1000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 170148,
            "unit": "B/op",
            "extra": "6272 times\n4 procs"
          },
          {
            "name": "BenchmarkList/as_slice/len=1000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 4004,
            "unit": "allocs/op",
            "extra": "6272 times\n4 procs"
          },
          {
            "name": "BenchmarkTextRead/poorly_simulated_typing/n=1000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 4145,
            "unit": "ns/op\t    1064 B/op\t       3 allocs/op",
            "extra": "294626 times\n4 procs"
          },
          {
            "name": "BenchmarkTextRead/poorly_simulated_typing/n=1000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 4145,
            "unit": "ns/op",
            "extra": "294626 times\n4 procs"
          },
          {
            "name": "BenchmarkTextRead/poorly_simulated_typing/n=1000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 1064,
            "unit": "B/op",
            "extra": "294626 times\n4 procs"
          },
          {
            "name": "BenchmarkTextRead/poorly_simulated_typing/n=1000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "294626 times\n4 procs"
          },
          {
            "name": "BenchmarkTextRead/edit_trace (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 479819,
            "unit": "ns/op\t  106520 B/op\t       2 allocs/op",
            "extra": "2437 times\n4 procs"
          },
          {
            "name": "BenchmarkTextRead/edit_trace (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 479819,
            "unit": "ns/op",
            "extra": "2437 times\n4 procs"
          },
          {
            "name": "BenchmarkTextRead/edit_trace (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 106520,
            "unit": "B/op",
            "extra": "2437 times\n4 procs"
          },
          {
            "name": "BenchmarkTextRead/edit_trace (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "2437 times\n4 procs"
          },
          {
            "name": "BenchmarkSaveIncremental/incremental (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 3599,
            "unit": "ns/op\t 268.98 MB/s\t    2512 B/op\t       6 allocs/op",
            "extra": "322792 times\n4 procs"
          },
          {
            "name": "BenchmarkSaveIncremental/incremental (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 3599,
            "unit": "ns/op",
            "extra": "322792 times\n4 procs"
          },
          {
            "name": "BenchmarkSaveIncremental/incremental (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 268.98,
            "unit": "MB/s",
            "extra": "322792 times\n4 procs"
          },
          {
            "name": "BenchmarkSaveIncremental/incremental (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 2512,
            "unit": "B/op",
            "extra": "322792 times\n4 procs"
          },
          {
            "name": "BenchmarkSaveIncremental/incremental (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 6,
            "unit": "allocs/op",
            "extra": "322792 times\n4 procs"
          },
          {
            "name": "BenchmarkSaveIncremental/full (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 17206,
            "unit": "ns/op\t   9.71 MB/s\t    6480 B/op\t      94 allocs/op",
            "extra": "71383 times\n4 procs"
          },
          {
            "name": "BenchmarkSaveIncremental/full (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 17206,
            "unit": "ns/op",
            "extra": "71383 times\n4 procs"
          },
          {
            "name": "BenchmarkSaveIncremental/full (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 9.71,
            "unit": "MB/s",
            "extra": "71383 times\n4 procs"
          },
          {
            "name": "BenchmarkSaveIncremental/full (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 6480,
            "unit": "B/op",
            "extra": "71383 times\n4 procs"
          },
          {
            "name": "BenchmarkSaveIncremental/full (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 94,
            "unit": "allocs/op",
            "extra": "71383 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteChange/ops=100 (github.com/MichaelMure/gotomerge/format)",
            "value": 67454,
            "unit": "ns/op\t   5.46 MB/s\t   14807 B/op\t     181 allocs/op",
            "extra": "17857 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteChange/ops=100 (github.com/MichaelMure/gotomerge/format) - ns/op",
            "value": 67454,
            "unit": "ns/op",
            "extra": "17857 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteChange/ops=100 (github.com/MichaelMure/gotomerge/format) - MB/s",
            "value": 5.46,
            "unit": "MB/s",
            "extra": "17857 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteChange/ops=100 (github.com/MichaelMure/gotomerge/format) - B/op",
            "value": 14807,
            "unit": "B/op",
            "extra": "17857 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteChange/ops=100 (github.com/MichaelMure/gotomerge/format) - allocs/op",
            "value": 181,
            "unit": "allocs/op",
            "extra": "17857 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteChange/ops=1000 (github.com/MichaelMure/gotomerge/format)",
            "value": 413067,
            "unit": "ns/op\t   9.37 MB/s\t  106504 B/op\t    1094 allocs/op",
            "extra": "2620 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteChange/ops=1000 (github.com/MichaelMure/gotomerge/format) - ns/op",
            "value": 413067,
            "unit": "ns/op",
            "extra": "2620 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteChange/ops=1000 (github.com/MichaelMure/gotomerge/format) - MB/s",
            "value": 9.37,
            "unit": "MB/s",
            "extra": "2620 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteChange/ops=1000 (github.com/MichaelMure/gotomerge/format) - B/op",
            "value": 106504,
            "unit": "B/op",
            "extra": "2620 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteChange/ops=1000 (github.com/MichaelMure/gotomerge/format) - allocs/op",
            "value": 1094,
            "unit": "allocs/op",
            "extra": "2620 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteChange/ops=10000 (github.com/MichaelMure/gotomerge/format)",
            "value": 4218519,
            "unit": "ns/op\t   9.81 MB/s\t 1898680 B/op\t   10226 allocs/op",
            "extra": "285 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteChange/ops=10000 (github.com/MichaelMure/gotomerge/format) - ns/op",
            "value": 4218519,
            "unit": "ns/op",
            "extra": "285 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteChange/ops=10000 (github.com/MichaelMure/gotomerge/format) - MB/s",
            "value": 9.81,
            "unit": "MB/s",
            "extra": "285 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteChange/ops=10000 (github.com/MichaelMure/gotomerge/format) - B/op",
            "value": 1898680,
            "unit": "B/op",
            "extra": "285 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteChange/ops=10000 (github.com/MichaelMure/gotomerge/format) - allocs/op",
            "value": 10226,
            "unit": "allocs/op",
            "extra": "285 times\n4 procs"
          },
          {
            "name": "BenchmarkReadExamplar (github.com/MichaelMure/gotomerge/format)",
            "value": 9608,
            "unit": "ns/op\t  42.25 MB/s\t         1.000 chunks\t    2888 B/op\t      46 allocs/op",
            "extra": "122702 times\n4 procs"
          },
          {
            "name": "BenchmarkReadExamplar (github.com/MichaelMure/gotomerge/format) - ns/op",
            "value": 9608,
            "unit": "ns/op",
            "extra": "122702 times\n4 procs"
          },
          {
            "name": "BenchmarkReadExamplar (github.com/MichaelMure/gotomerge/format) - MB/s",
            "value": 42.25,
            "unit": "MB/s",
            "extra": "122702 times\n4 procs"
          },
          {
            "name": "BenchmarkReadExamplar (github.com/MichaelMure/gotomerge/format) - chunks",
            "value": 1,
            "unit": "chunks",
            "extra": "122702 times\n4 procs"
          },
          {
            "name": "BenchmarkReadExamplar (github.com/MichaelMure/gotomerge/format) - B/op",
            "value": 2888,
            "unit": "B/op",
            "extra": "122702 times\n4 procs"
          },
          {
            "name": "BenchmarkReadExamplar (github.com/MichaelMure/gotomerge/format) - allocs/op",
            "value": 46,
            "unit": "allocs/op",
            "extra": "122702 times\n4 procs"
          },
          {
            "name": "BenchmarkReadLarge (github.com/MichaelMure/gotomerge/format)",
            "value": 388506293,
            "unit": "ns/op\t  75.29 MB/s\t    259779 chunks\t462359680 B/op\t 5792673 allocs/op",
            "extra": "3 times\n4 procs"
          },
          {
            "name": "BenchmarkReadLarge (github.com/MichaelMure/gotomerge/format) - ns/op",
            "value": 388506293,
            "unit": "ns/op",
            "extra": "3 times\n4 procs"
          },
          {
            "name": "BenchmarkReadLarge (github.com/MichaelMure/gotomerge/format) - MB/s",
            "value": 75.29,
            "unit": "MB/s",
            "extra": "3 times\n4 procs"
          },
          {
            "name": "BenchmarkReadLarge (github.com/MichaelMure/gotomerge/format) - chunks",
            "value": 259779,
            "unit": "chunks",
            "extra": "3 times\n4 procs"
          },
          {
            "name": "BenchmarkReadLarge (github.com/MichaelMure/gotomerge/format) - B/op",
            "value": 462359680,
            "unit": "B/op",
            "extra": "3 times\n4 procs"
          },
          {
            "name": "BenchmarkReadLarge (github.com/MichaelMure/gotomerge/format) - allocs/op",
            "value": 5792673,
            "unit": "allocs/op",
            "extra": "3 times\n4 procs"
          },
          {
            "name": "BenchmarkCompressed (github.com/MichaelMure/gotomerge/format)",
            "value": 27401,
            "unit": "ns/op\t   7.01 MB/s\t         2.000 chunks\t   85768 B/op\t      66 allocs/op",
            "extra": "42055 times\n4 procs"
          },
          {
            "name": "BenchmarkCompressed (github.com/MichaelMure/gotomerge/format) - ns/op",
            "value": 27401,
            "unit": "ns/op",
            "extra": "42055 times\n4 procs"
          },
          {
            "name": "BenchmarkCompressed (github.com/MichaelMure/gotomerge/format) - MB/s",
            "value": 7.01,
            "unit": "MB/s",
            "extra": "42055 times\n4 procs"
          },
          {
            "name": "BenchmarkCompressed (github.com/MichaelMure/gotomerge/format) - chunks",
            "value": 2,
            "unit": "chunks",
            "extra": "42055 times\n4 procs"
          },
          {
            "name": "BenchmarkCompressed (github.com/MichaelMure/gotomerge/format) - B/op",
            "value": 85768,
            "unit": "B/op",
            "extra": "42055 times\n4 procs"
          },
          {
            "name": "BenchmarkCompressed (github.com/MichaelMure/gotomerge/format) - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "42055 times\n4 procs"
          },
          {
            "name": "BenchmarkTextEdits (github.com/MichaelMure/gotomerge/opset)",
            "value": 631151,
            "unit": "ns/op\t    104852 chars\t  183316 B/op\t     321 allocs/op",
            "extra": "1720 times\n4 procs"
          },
          {
            "name": "BenchmarkTextEdits (github.com/MichaelMure/gotomerge/opset) - ns/op",
            "value": 631151,
            "unit": "ns/op",
            "extra": "1720 times\n4 procs"
          },
          {
            "name": "BenchmarkTextEdits (github.com/MichaelMure/gotomerge/opset) - chars",
            "value": 104852,
            "unit": "chars",
            "extra": "1720 times\n4 procs"
          },
          {
            "name": "BenchmarkTextEdits (github.com/MichaelMure/gotomerge/opset) - B/op",
            "value": 183316,
            "unit": "B/op",
            "extra": "1720 times\n4 procs"
          },
          {
            "name": "BenchmarkTextEdits (github.com/MichaelMure/gotomerge/opset) - allocs/op",
            "value": 321,
            "unit": "allocs/op",
            "extra": "1720 times\n4 procs"
          }
        ]
      }
    ]
  }
}