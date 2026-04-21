window.BENCHMARK_DATA = {
  "lastUpdate": 1776761204694,
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
      },
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
          "id": "a18a3070bbc19e440bb7acd5e7d83244394db039",
          "message": "ci: update actions",
          "timestamp": "2026-04-21T10:37:37+02:00",
          "tree_id": "de60b4b1d81997080d6275ef86d6b6ba7307cd61",
          "url": "https://github.com/MichaelMure/gotomerge/commit/a18a3070bbc19e440bb7acd5e7d83244394db039"
        },
        "date": 1776760997994,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkReadMetadata (github.com/MichaelMure/gotomerge/column)",
            "value": 179.5,
            "unit": "ns/op\t      80 B/op\t       2 allocs/op",
            "extra": "6632394 times\n4 procs"
          },
          {
            "name": "BenchmarkReadMetadata (github.com/MichaelMure/gotomerge/column) - ns/op",
            "value": 179.5,
            "unit": "ns/op",
            "extra": "6632394 times\n4 procs"
          },
          {
            "name": "BenchmarkReadMetadata (github.com/MichaelMure/gotomerge/column) - B/op",
            "value": 80,
            "unit": "B/op",
            "extra": "6632394 times\n4 procs"
          },
          {
            "name": "BenchmarkReadMetadata (github.com/MichaelMure/gotomerge/column) - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "6632394 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteMetadata (github.com/MichaelMure/gotomerge/column)",
            "value": 84.11,
            "unit": "ns/op\t      40 B/op\t       5 allocs/op",
            "extra": "13942390 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteMetadata (github.com/MichaelMure/gotomerge/column) - ns/op",
            "value": 84.11,
            "unit": "ns/op",
            "extra": "13942390 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteMetadata (github.com/MichaelMure/gotomerge/column) - B/op",
            "value": 40,
            "unit": "B/op",
            "extra": "13942390 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteMetadata (github.com/MichaelMure/gotomerge/column) - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "13942390 times\n4 procs"
          },
          {
            "name": "BenchmarkReadSpecification (github.com/MichaelMure/gotomerge/column)",
            "value": 177.7,
            "unit": "ns/op\t      80 B/op\t       2 allocs/op",
            "extra": "6948130 times\n4 procs"
          },
          {
            "name": "BenchmarkReadSpecification (github.com/MichaelMure/gotomerge/column) - ns/op",
            "value": 177.7,
            "unit": "ns/op",
            "extra": "6948130 times\n4 procs"
          },
          {
            "name": "BenchmarkReadSpecification (github.com/MichaelMure/gotomerge/column) - B/op",
            "value": 80,
            "unit": "B/op",
            "extra": "6948130 times\n4 procs"
          },
          {
            "name": "BenchmarkReadSpecification (github.com/MichaelMure/gotomerge/column) - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "6948130 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteSpecification (github.com/MichaelMure/gotomerge/column)",
            "value": 39.67,
            "unit": "ns/op\t      16 B/op\t       2 allocs/op",
            "extra": "29552283 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteSpecification (github.com/MichaelMure/gotomerge/column) - ns/op",
            "value": 39.67,
            "unit": "ns/op",
            "extra": "29552283 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteSpecification (github.com/MichaelMure/gotomerge/column) - B/op",
            "value": 16,
            "unit": "B/op",
            "extra": "29552283 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteSpecification (github.com/MichaelMure/gotomerge/column) - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "29552283 times\n4 procs"
          },
          {
            "name": "BenchmarkUint64Reader (github.com/MichaelMure/gotomerge/column/rle)",
            "value": 126.6,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "9525066 times\n4 procs"
          },
          {
            "name": "BenchmarkUint64Reader (github.com/MichaelMure/gotomerge/column/rle) - ns/op",
            "value": 126.6,
            "unit": "ns/op",
            "extra": "9525066 times\n4 procs"
          },
          {
            "name": "BenchmarkUint64Reader (github.com/MichaelMure/gotomerge/column/rle) - B/op",
            "value": 80,
            "unit": "B/op",
            "extra": "9525066 times\n4 procs"
          },
          {
            "name": "BenchmarkUint64Reader (github.com/MichaelMure/gotomerge/column/rle) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "9525066 times\n4 procs"
          },
          {
            "name": "BenchmarkInt64Reader (github.com/MichaelMure/gotomerge/column/rle)",
            "value": 121.6,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "9694692 times\n4 procs"
          },
          {
            "name": "BenchmarkInt64Reader (github.com/MichaelMure/gotomerge/column/rle) - ns/op",
            "value": 121.6,
            "unit": "ns/op",
            "extra": "9694692 times\n4 procs"
          },
          {
            "name": "BenchmarkInt64Reader (github.com/MichaelMure/gotomerge/column/rle) - B/op",
            "value": 80,
            "unit": "B/op",
            "extra": "9694692 times\n4 procs"
          },
          {
            "name": "BenchmarkInt64Reader (github.com/MichaelMure/gotomerge/column/rle) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "9694692 times\n4 procs"
          },
          {
            "name": "BenchmarkStringReader (github.com/MichaelMure/gotomerge/column/rle)",
            "value": 194.7,
            "unit": "ns/op\t     104 B/op\t       4 allocs/op",
            "extra": "6149227 times\n4 procs"
          },
          {
            "name": "BenchmarkStringReader (github.com/MichaelMure/gotomerge/column/rle) - ns/op",
            "value": 194.7,
            "unit": "ns/op",
            "extra": "6149227 times\n4 procs"
          },
          {
            "name": "BenchmarkStringReader (github.com/MichaelMure/gotomerge/column/rle) - B/op",
            "value": 104,
            "unit": "B/op",
            "extra": "6149227 times\n4 procs"
          },
          {
            "name": "BenchmarkStringReader (github.com/MichaelMure/gotomerge/column/rle) - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "6149227 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/repeated_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 47657,
            "unit": "ns/op\t   61600 B/op\t     354 allocs/op",
            "extra": "26905 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/repeated_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 47657,
            "unit": "ns/op",
            "extra": "26905 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/repeated_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 61600,
            "unit": "B/op",
            "extra": "26905 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/repeated_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 354,
            "unit": "allocs/op",
            "extra": "26905 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/repeated_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 608533,
            "unit": "ns/op\t  699705 B/op\t    2967 allocs/op",
            "extra": "2037 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/repeated_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 608533,
            "unit": "ns/op",
            "extra": "2037 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/repeated_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 699705,
            "unit": "B/op",
            "extra": "2037 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/repeated_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 2967,
            "unit": "allocs/op",
            "extra": "2037 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/repeated_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 9209980,
            "unit": "ns/op\t 9102246 B/op\t   30069 allocs/op",
            "extra": "134 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/repeated_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 9209980,
            "unit": "ns/op",
            "extra": "134 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/repeated_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 9102246,
            "unit": "B/op",
            "extra": "134 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/repeated_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 30069,
            "unit": "allocs/op",
            "extra": "134 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/increasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 155888,
            "unit": "ns/op\t  160291 B/op\t     479 allocs/op",
            "extra": "8161 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/increasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 155888,
            "unit": "ns/op",
            "extra": "8161 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/increasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 160291,
            "unit": "B/op",
            "extra": "8161 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/increasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 479,
            "unit": "allocs/op",
            "extra": "8161 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/increasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 1023496,
            "unit": "ns/op\t 1091993 B/op\t    4946 allocs/op",
            "extra": "1158 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/increasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 1023496,
            "unit": "ns/op",
            "extra": "1158 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/increasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 1091993,
            "unit": "B/op",
            "extra": "1158 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/increasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 4946,
            "unit": "allocs/op",
            "extra": "1158 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/increasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 12768524,
            "unit": "ns/op\t11785122 B/op\t   50067 allocs/op",
            "extra": "97 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/increasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 12768524,
            "unit": "ns/op",
            "extra": "97 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/increasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 11785122,
            "unit": "B/op",
            "extra": "97 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/increasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 50067,
            "unit": "allocs/op",
            "extra": "97 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/decreasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 157350,
            "unit": "ns/op\t  162084 B/op\t     479 allocs/op",
            "extra": "8336 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/decreasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 157350,
            "unit": "ns/op",
            "extra": "8336 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/decreasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 162084,
            "unit": "B/op",
            "extra": "8336 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/decreasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 479,
            "unit": "allocs/op",
            "extra": "8336 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/decreasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 1047275,
            "unit": "ns/op\t 1137925 B/op\t    4947 allocs/op",
            "extra": "1124 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/decreasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 1047275,
            "unit": "ns/op",
            "extra": "1124 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/decreasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 1137925,
            "unit": "B/op",
            "extra": "1124 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/decreasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 4947,
            "unit": "allocs/op",
            "extra": "1124 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/decreasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 13759895,
            "unit": "ns/op\t11785301 B/op\t   50072 allocs/op",
            "extra": "84 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/decreasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 13759895,
            "unit": "ns/op",
            "extra": "84 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/decreasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 11785301,
            "unit": "B/op",
            "extra": "84 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/decreasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 50072,
            "unit": "allocs/op",
            "extra": "84 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/repeated_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 25178,
            "unit": "ns/op\t   9.53 MB/s\t    5952 B/op\t     167 allocs/op",
            "extra": "48921 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/repeated_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 25178,
            "unit": "ns/op",
            "extra": "48921 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/repeated_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 9.53,
            "unit": "MB/s",
            "extra": "48921 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/repeated_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 5952,
            "unit": "B/op",
            "extra": "48921 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/repeated_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 167,
            "unit": "allocs/op",
            "extra": "48921 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/repeated_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 254633,
            "unit": "ns/op\t   6.41 MB/s\t  226130 B/op\t    1079 allocs/op",
            "extra": "4797 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/repeated_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 254633,
            "unit": "ns/op",
            "extra": "4797 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/repeated_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 6.41,
            "unit": "MB/s",
            "extra": "4797 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/repeated_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 226130,
            "unit": "B/op",
            "extra": "4797 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/repeated_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 1079,
            "unit": "allocs/op",
            "extra": "4797 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/repeated_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 1675562,
            "unit": "ns/op\t  11.30 MB/s\t 1069243 B/op\t   10100 allocs/op",
            "extra": "715 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/repeated_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 1675562,
            "unit": "ns/op",
            "extra": "715 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/repeated_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 11.3,
            "unit": "MB/s",
            "extra": "715 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/repeated_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 1069243,
            "unit": "B/op",
            "extra": "715 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/repeated_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 10100,
            "unit": "allocs/op",
            "extra": "715 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/increasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 81151,
            "unit": "ns/op\t   4.73 MB/s\t   67386 B/op\t     181 allocs/op",
            "extra": "15247 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/increasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 81151,
            "unit": "ns/op",
            "extra": "15247 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/increasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 4.73,
            "unit": "MB/s",
            "extra": "15247 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/increasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 67386,
            "unit": "B/op",
            "extra": "15247 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/increasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 181,
            "unit": "allocs/op",
            "extra": "15247 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/increasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 397404,
            "unit": "ns/op\t   8.43 MB/s\t  431777 B/op\t    1104 allocs/op",
            "extra": "2799 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/increasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 397404,
            "unit": "ns/op",
            "extra": "2799 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/increasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 8.43,
            "unit": "MB/s",
            "extra": "2799 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/increasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 431777,
            "unit": "B/op",
            "extra": "2799 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/increasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 1104,
            "unit": "allocs/op",
            "extra": "2799 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/increasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 2927686,
            "unit": "ns/op\t  12.33 MB/s\t 2032545 B/op\t   10137 allocs/op",
            "extra": "415 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/increasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 2927686,
            "unit": "ns/op",
            "extra": "415 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/increasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 12.33,
            "unit": "MB/s",
            "extra": "415 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/increasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 2032545,
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
            "value": 83130,
            "unit": "ns/op\t   4.62 MB/s\t   68408 B/op\t     181 allocs/op",
            "extra": "14244 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/decreasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 83130,
            "unit": "ns/op",
            "extra": "14244 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/decreasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 4.62,
            "unit": "MB/s",
            "extra": "14244 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/decreasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 68408,
            "unit": "B/op",
            "extra": "14244 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/decreasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 181,
            "unit": "allocs/op",
            "extra": "14244 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/decreasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 406922,
            "unit": "ns/op\t   8.23 MB/s\t  425260 B/op\t    1103 allocs/op",
            "extra": "2972 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/decreasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 406922,
            "unit": "ns/op",
            "extra": "2972 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/decreasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 8.23,
            "unit": "MB/s",
            "extra": "2972 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/decreasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 425260,
            "unit": "B/op",
            "extra": "2972 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/decreasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 1103,
            "unit": "allocs/op",
            "extra": "2972 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/decreasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 2991139,
            "unit": "ns/op\t  11.97 MB/s\t 2032556 B/op\t   10137 allocs/op",
            "extra": "404 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/decreasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 2991139,
            "unit": "ns/op",
            "extra": "404 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/decreasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 11.97,
            "unit": "MB/s",
            "extra": "404 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/decreasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 2032556,
            "unit": "B/op",
            "extra": "404 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/decreasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 10137,
            "unit": "allocs/op",
            "extra": "404 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/repeated_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 30738,
            "unit": "ns/op\t   7.81 MB/s\t   13048 B/op\t     208 allocs/op",
            "extra": "38871 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/repeated_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 30738,
            "unit": "ns/op",
            "extra": "38871 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/repeated_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 7.81,
            "unit": "MB/s",
            "extra": "38871 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/repeated_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 13048,
            "unit": "B/op",
            "extra": "38871 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/repeated_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 208,
            "unit": "allocs/op",
            "extra": "38871 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/repeated_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 334544,
            "unit": "ns/op\t   4.88 MB/s\t  166360 B/op\t    2062 allocs/op",
            "extra": "3548 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/repeated_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 334544,
            "unit": "ns/op",
            "extra": "3548 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/repeated_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 4.88,
            "unit": "MB/s",
            "extra": "3548 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/repeated_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 166360,
            "unit": "B/op",
            "extra": "3548 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/repeated_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 2062,
            "unit": "allocs/op",
            "extra": "3548 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/repeated_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 3259301,
            "unit": "ns/op\t   5.81 MB/s\t 1205668 B/op\t   21690 allocs/op",
            "extra": "370 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/repeated_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 3259301,
            "unit": "ns/op",
            "extra": "370 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/repeated_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 5.81,
            "unit": "MB/s",
            "extra": "370 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/repeated_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 1205668,
            "unit": "B/op",
            "extra": "370 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/repeated_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 21690,
            "unit": "allocs/op",
            "extra": "370 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/increasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 57679,
            "unit": "ns/op\t   6.66 MB/s\t   65024 B/op\t     506 allocs/op",
            "extra": "20838 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/increasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 57679,
            "unit": "ns/op",
            "extra": "20838 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/increasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 6.66,
            "unit": "MB/s",
            "extra": "20838 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/increasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 65024,
            "unit": "B/op",
            "extra": "20838 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/increasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 506,
            "unit": "allocs/op",
            "extra": "20838 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/increasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 554992,
            "unit": "ns/op\t   6.04 MB/s\t  413073 B/op\t    5077 allocs/op",
            "extra": "2125 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/increasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 554992,
            "unit": "ns/op",
            "extra": "2125 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/increasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 6.04,
            "unit": "MB/s",
            "extra": "2125 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/increasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 413073,
            "unit": "B/op",
            "extra": "2125 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/increasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 5077,
            "unit": "allocs/op",
            "extra": "2125 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/increasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 5320605,
            "unit": "ns/op\t   6.79 MB/s\t 3003305 B/op\t   51767 allocs/op",
            "extra": "225 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/increasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 5320605,
            "unit": "ns/op",
            "extra": "225 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/increasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 6.79,
            "unit": "MB/s",
            "extra": "225 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/increasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 3003305,
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
            "value": 57263,
            "unit": "ns/op\t   6.71 MB/s\t   65024 B/op\t     506 allocs/op",
            "extra": "21148 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/decreasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 57263,
            "unit": "ns/op",
            "extra": "21148 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/decreasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 6.71,
            "unit": "MB/s",
            "extra": "21148 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/decreasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 65024,
            "unit": "B/op",
            "extra": "21148 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/decreasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 506,
            "unit": "allocs/op",
            "extra": "21148 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/decreasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 548906,
            "unit": "ns/op\t   6.10 MB/s\t  413073 B/op\t    5077 allocs/op",
            "extra": "2200 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/decreasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 548906,
            "unit": "ns/op",
            "extra": "2200 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/decreasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 6.1,
            "unit": "MB/s",
            "extra": "2200 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/decreasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 413073,
            "unit": "B/op",
            "extra": "2200 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/decreasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 5077,
            "unit": "allocs/op",
            "extra": "2200 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/decreasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 5301680,
            "unit": "ns/op\t   6.76 MB/s\t 3003508 B/op\t   51767 allocs/op",
            "extra": "224 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/decreasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 5301680,
            "unit": "ns/op",
            "extra": "224 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/decreasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 6.76,
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
            "value": 32785,
            "unit": "ns/op\t   6.10 MB/s\t   33984 B/op\t     197 allocs/op",
            "extra": "36525 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/repeated_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 32785,
            "unit": "ns/op",
            "extra": "36525 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/repeated_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 6.1,
            "unit": "MB/s",
            "extra": "36525 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/repeated_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 33984,
            "unit": "B/op",
            "extra": "36525 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/repeated_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 197,
            "unit": "allocs/op",
            "extra": "36525 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/repeated_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 335381,
            "unit": "ns/op\t   4.79 MB/s\t  331266 B/op\t    1900 allocs/op",
            "extra": "3505 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/repeated_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 335381,
            "unit": "ns/op",
            "extra": "3505 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/repeated_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 4.79,
            "unit": "MB/s",
            "extra": "3505 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/repeated_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 331266,
            "unit": "B/op",
            "extra": "3505 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/repeated_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 1900,
            "unit": "allocs/op",
            "extra": "3505 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/repeated_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 4771412,
            "unit": "ns/op\t   3.96 MB/s\t 4492033 B/op\t   19978 allocs/op",
            "extra": "244 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/repeated_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 4771412,
            "unit": "ns/op",
            "extra": "244 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/repeated_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 3.96,
            "unit": "MB/s",
            "extra": "244 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/repeated_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 4492033,
            "unit": "B/op",
            "extra": "244 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/repeated_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 19978,
            "unit": "allocs/op",
            "extra": "244 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/increasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 60155,
            "unit": "ns/op\t   6.18 MB/s\t   86088 B/op\t     496 allocs/op",
            "extra": "19687 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/increasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 60155,
            "unit": "ns/op",
            "extra": "19687 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/increasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 6.18,
            "unit": "MB/s",
            "extra": "19687 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/increasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 86088,
            "unit": "B/op",
            "extra": "19687 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/increasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 496,
            "unit": "allocs/op",
            "extra": "19687 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/increasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 540117,
            "unit": "ns/op\t   7.16 MB/s\t  542571 B/op\t    4949 allocs/op",
            "extra": "2169 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/increasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 540117,
            "unit": "ns/op",
            "extra": "2169 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/increasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 7.16,
            "unit": "MB/s",
            "extra": "2169 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/increasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 542571,
            "unit": "B/op",
            "extra": "2169 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/increasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 4949,
            "unit": "allocs/op",
            "extra": "2169 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/increasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 7336148,
            "unit": "ns/op\t   5.16 MB/s\t 6217742 B/op\t   50038 allocs/op",
            "extra": "159 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/increasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 7336148,
            "unit": "ns/op",
            "extra": "159 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/increasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 5.16,
            "unit": "MB/s",
            "extra": "159 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/increasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 6217742,
            "unit": "B/op",
            "extra": "159 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/increasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 50038,
            "unit": "allocs/op",
            "extra": "159 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/decreasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 61424,
            "unit": "ns/op\t   6.09 MB/s\t   86088 B/op\t     496 allocs/op",
            "extra": "20121 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/decreasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 61424,
            "unit": "ns/op",
            "extra": "20121 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/decreasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 6.09,
            "unit": "MB/s",
            "extra": "20121 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/decreasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 86088,
            "unit": "B/op",
            "extra": "20121 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/decreasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 496,
            "unit": "allocs/op",
            "extra": "20121 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/decreasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 541504,
            "unit": "ns/op\t   7.13 MB/s\t  542571 B/op\t    4949 allocs/op",
            "extra": "2180 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/decreasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 541504,
            "unit": "ns/op",
            "extra": "2180 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/decreasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 7.13,
            "unit": "MB/s",
            "extra": "2180 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/decreasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 542571,
            "unit": "B/op",
            "extra": "2180 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/decreasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 4949,
            "unit": "allocs/op",
            "extra": "2180 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/decreasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 7485226,
            "unit": "ns/op\t   4.92 MB/s\t 6218001 B/op\t   50046 allocs/op",
            "extra": "159 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/decreasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 7485226,
            "unit": "ns/op",
            "extra": "159 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/decreasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 4.92,
            "unit": "MB/s",
            "extra": "159 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/decreasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 6218001,
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
            "value": 52669,
            "unit": "ns/op\t   76092 B/op\t     162 allocs/op",
            "extra": "22480 times\n4 procs"
          },
          {
            "name": "BenchmarkSaveLoad/big_paste (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 52669,
            "unit": "ns/op",
            "extra": "22480 times\n4 procs"
          },
          {
            "name": "BenchmarkSaveLoad/big_paste (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 76092,
            "unit": "B/op",
            "extra": "22480 times\n4 procs"
          },
          {
            "name": "BenchmarkSaveLoad/big_paste (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 162,
            "unit": "allocs/op",
            "extra": "22480 times\n4 procs"
          },
          {
            "name": "BenchmarkSaveLoad/maps_in_maps (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 1786796,
            "unit": "ns/op\t 1332229 B/op\t    6493 allocs/op",
            "extra": "661 times\n4 procs"
          },
          {
            "name": "BenchmarkSaveLoad/maps_in_maps (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 1786796,
            "unit": "ns/op",
            "extra": "661 times\n4 procs"
          },
          {
            "name": "BenchmarkSaveLoad/maps_in_maps (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 1332229,
            "unit": "B/op",
            "extra": "661 times\n4 procs"
          },
          {
            "name": "BenchmarkSaveLoad/maps_in_maps (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 6493,
            "unit": "allocs/op",
            "extra": "661 times\n4 procs"
          },
          {
            "name": "BenchmarkSaveLoad/deep_history (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 1512504,
            "unit": "ns/op\t  856453 B/op\t   10134 allocs/op",
            "extra": "790 times\n4 procs"
          },
          {
            "name": "BenchmarkSaveLoad/deep_history (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 1512504,
            "unit": "ns/op",
            "extra": "790 times\n4 procs"
          },
          {
            "name": "BenchmarkSaveLoad/deep_history (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 856453,
            "unit": "B/op",
            "extra": "790 times\n4 procs"
          },
          {
            "name": "BenchmarkSaveLoad/deep_history (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 10134,
            "unit": "allocs/op",
            "extra": "790 times\n4 procs"
          },
          {
            "name": "BenchmarkSaveLoad/poorly_simulated_typing (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 523054,
            "unit": "ns/op\t  267377 B/op\t    4434 allocs/op",
            "extra": "2354 times\n4 procs"
          },
          {
            "name": "BenchmarkSaveLoad/poorly_simulated_typing (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 523054,
            "unit": "ns/op",
            "extra": "2354 times\n4 procs"
          },
          {
            "name": "BenchmarkSaveLoad/poorly_simulated_typing (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 267377,
            "unit": "B/op",
            "extra": "2354 times\n4 procs"
          },
          {
            "name": "BenchmarkSaveLoad/poorly_simulated_typing (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 4434,
            "unit": "allocs/op",
            "extra": "2354 times\n4 procs"
          },
          {
            "name": "BenchmarkEditTrace/apply (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 924871333,
            "unit": "ns/op\t  31.63 MB/s\t1072416688 B/op\t12966370 allocs/op",
            "extra": "2 times\n4 procs"
          },
          {
            "name": "BenchmarkEditTrace/apply (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 924871333,
            "unit": "ns/op",
            "extra": "2 times\n4 procs"
          },
          {
            "name": "BenchmarkEditTrace/apply (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 31.63,
            "unit": "MB/s",
            "extra": "2 times\n4 procs"
          },
          {
            "name": "BenchmarkEditTrace/apply (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 1072416688,
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
            "value": 1091015467,
            "unit": "ns/op\t   0.13 MB/s\t1109294032 B/op\t13346972 allocs/op",
            "extra": "1 times\n4 procs"
          },
          {
            "name": "BenchmarkEditTrace/save (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 1091015467,
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
            "value": 1109294032,
            "unit": "B/op",
            "extra": "1 times\n4 procs"
          },
          {
            "name": "BenchmarkEditTrace/save (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 13346972,
            "unit": "allocs/op",
            "extra": "1 times\n4 procs"
          },
          {
            "name": "BenchmarkEditTrace/load (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 85650593,
            "unit": "ns/op\t   1.65 MB/s\t31288001 B/op\t  761472 allocs/op",
            "extra": "13 times\n4 procs"
          },
          {
            "name": "BenchmarkEditTrace/load (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 85650593,
            "unit": "ns/op",
            "extra": "13 times\n4 procs"
          },
          {
            "name": "BenchmarkEditTrace/load (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 1.65,
            "unit": "MB/s",
            "extra": "13 times\n4 procs"
          },
          {
            "name": "BenchmarkEditTrace/load (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 31288001,
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
            "value": 36749216,
            "unit": "ns/op\t48200968 B/op\t  100559 allocs/op",
            "extra": "39 times\n4 procs"
          },
          {
            "name": "BenchmarkRange (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 36749216,
            "unit": "ns/op",
            "extra": "39 times\n4 procs"
          },
          {
            "name": "BenchmarkRange (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 48200968,
            "unit": "B/op",
            "extra": "39 times\n4 procs"
          },
          {
            "name": "BenchmarkRange (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 100559,
            "unit": "allocs/op",
            "extra": "39 times\n4 procs"
          },
          {
            "name": "BenchmarkGet/scalar/keys=100 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 93.47,
            "unit": "ns/op\t      48 B/op\t       1 allocs/op",
            "extra": "12848571 times\n4 procs"
          },
          {
            "name": "BenchmarkGet/scalar/keys=100 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 93.47,
            "unit": "ns/op",
            "extra": "12848571 times\n4 procs"
          },
          {
            "name": "BenchmarkGet/scalar/keys=100 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 48,
            "unit": "B/op",
            "extra": "12848571 times\n4 procs"
          },
          {
            "name": "BenchmarkGet/scalar/keys=100 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "12848571 times\n4 procs"
          },
          {
            "name": "BenchmarkGet/scalar/keys=1000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 97.35,
            "unit": "ns/op\t      48 B/op\t       1 allocs/op",
            "extra": "12799750 times\n4 procs"
          },
          {
            "name": "BenchmarkGet/scalar/keys=1000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 97.35,
            "unit": "ns/op",
            "extra": "12799750 times\n4 procs"
          },
          {
            "name": "BenchmarkGet/scalar/keys=1000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 48,
            "unit": "B/op",
            "extra": "12799750 times\n4 procs"
          },
          {
            "name": "BenchmarkGet/scalar/keys=1000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "12799750 times\n4 procs"
          },
          {
            "name": "BenchmarkGet/scalar/keys=10000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 95.49,
            "unit": "ns/op\t      48 B/op\t       1 allocs/op",
            "extra": "12953419 times\n4 procs"
          },
          {
            "name": "BenchmarkGet/scalar/keys=10000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 95.49,
            "unit": "ns/op",
            "extra": "12953419 times\n4 procs"
          },
          {
            "name": "BenchmarkGet/scalar/keys=10000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 48,
            "unit": "B/op",
            "extra": "12953419 times\n4 procs"
          },
          {
            "name": "BenchmarkGet/scalar/keys=10000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "12953419 times\n4 procs"
          },
          {
            "name": "BenchmarkGet/nested/keys=100 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 90.93,
            "unit": "ns/op\t      48 B/op\t       1 allocs/op",
            "extra": "13530848 times\n4 procs"
          },
          {
            "name": "BenchmarkGet/nested/keys=100 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 90.93,
            "unit": "ns/op",
            "extra": "13530848 times\n4 procs"
          },
          {
            "name": "BenchmarkGet/nested/keys=100 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 48,
            "unit": "B/op",
            "extra": "13530848 times\n4 procs"
          },
          {
            "name": "BenchmarkGet/nested/keys=100 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "13530848 times\n4 procs"
          },
          {
            "name": "BenchmarkGet/nested/keys=1000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 93.77,
            "unit": "ns/op\t      48 B/op\t       1 allocs/op",
            "extra": "13080481 times\n4 procs"
          },
          {
            "name": "BenchmarkGet/nested/keys=1000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 93.77,
            "unit": "ns/op",
            "extra": "13080481 times\n4 procs"
          },
          {
            "name": "BenchmarkGet/nested/keys=1000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 48,
            "unit": "B/op",
            "extra": "13080481 times\n4 procs"
          },
          {
            "name": "BenchmarkGet/nested/keys=1000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "13080481 times\n4 procs"
          },
          {
            "name": "BenchmarkGet/nested/keys=10000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 93.51,
            "unit": "ns/op\t      48 B/op\t       1 allocs/op",
            "extra": "13089325 times\n4 procs"
          },
          {
            "name": "BenchmarkGet/nested/keys=10000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 93.51,
            "unit": "ns/op",
            "extra": "13089325 times\n4 procs"
          },
          {
            "name": "BenchmarkGet/nested/keys=10000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 48,
            "unit": "B/op",
            "extra": "13089325 times\n4 procs"
          },
          {
            "name": "BenchmarkGet/nested/keys=10000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "13089325 times\n4 procs"
          },
          {
            "name": "BenchmarkAs/int64 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 148.8,
            "unit": "ns/op\t      64 B/op\t       3 allocs/op",
            "extra": "8265404 times\n4 procs"
          },
          {
            "name": "BenchmarkAs/int64 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 148.8,
            "unit": "ns/op",
            "extra": "8265404 times\n4 procs"
          },
          {
            "name": "BenchmarkAs/int64 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 64,
            "unit": "B/op",
            "extra": "8265404 times\n4 procs"
          },
          {
            "name": "BenchmarkAs/int64 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "8265404 times\n4 procs"
          },
          {
            "name": "BenchmarkAs/string (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 120.5,
            "unit": "ns/op\t      80 B/op\t       2 allocs/op",
            "extra": "9953690 times\n4 procs"
          },
          {
            "name": "BenchmarkAs/string (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 120.5,
            "unit": "ns/op",
            "extra": "9953690 times\n4 procs"
          },
          {
            "name": "BenchmarkAs/string (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 80,
            "unit": "B/op",
            "extra": "9953690 times\n4 procs"
          },
          {
            "name": "BenchmarkAs/string (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "9953690 times\n4 procs"
          },
          {
            "name": "BenchmarkAs/struct (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 861.1,
            "unit": "ns/op\t     320 B/op\t      14 allocs/op",
            "extra": "1389361 times\n4 procs"
          },
          {
            "name": "BenchmarkAs/struct (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 861.1,
            "unit": "ns/op",
            "extra": "1389361 times\n4 procs"
          },
          {
            "name": "BenchmarkAs/struct (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 320,
            "unit": "B/op",
            "extra": "1389361 times\n4 procs"
          },
          {
            "name": "BenchmarkAs/struct (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 14,
            "unit": "allocs/op",
            "extra": "1389361 times\n4 procs"
          },
          {
            "name": "BenchmarkAs/slice_string/len=100 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 16534,
            "unit": "ns/op\t   20432 B/op\t     304 allocs/op",
            "extra": "69652 times\n4 procs"
          },
          {
            "name": "BenchmarkAs/slice_string/len=100 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 16534,
            "unit": "ns/op",
            "extra": "69652 times\n4 procs"
          },
          {
            "name": "BenchmarkAs/slice_string/len=100 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 20432,
            "unit": "B/op",
            "extra": "69652 times\n4 procs"
          },
          {
            "name": "BenchmarkAs/slice_string/len=100 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 304,
            "unit": "allocs/op",
            "extra": "69652 times\n4 procs"
          },
          {
            "name": "BenchmarkMerge/changes=10 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 44971,
            "unit": "ns/op\t   33808 B/op\t     442 allocs/op",
            "extra": "27092 times\n4 procs"
          },
          {
            "name": "BenchmarkMerge/changes=10 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 44971,
            "unit": "ns/op",
            "extra": "27092 times\n4 procs"
          },
          {
            "name": "BenchmarkMerge/changes=10 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 33808,
            "unit": "B/op",
            "extra": "27092 times\n4 procs"
          },
          {
            "name": "BenchmarkMerge/changes=10 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 442,
            "unit": "allocs/op",
            "extra": "27092 times\n4 procs"
          },
          {
            "name": "BenchmarkMerge/changes=100 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 309687,
            "unit": "ns/op\t  320624 B/op\t    4312 allocs/op",
            "extra": "4058 times\n4 procs"
          },
          {
            "name": "BenchmarkMerge/changes=100 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 309687,
            "unit": "ns/op",
            "extra": "4058 times\n4 procs"
          },
          {
            "name": "BenchmarkMerge/changes=100 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 320624,
            "unit": "B/op",
            "extra": "4058 times\n4 procs"
          },
          {
            "name": "BenchmarkMerge/changes=100 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 4312,
            "unit": "allocs/op",
            "extra": "4058 times\n4 procs"
          },
          {
            "name": "BenchmarkMerge/changes=1000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 3192693,
            "unit": "ns/op\t 3533105 B/op\t   43780 allocs/op",
            "extra": "379 times\n4 procs"
          },
          {
            "name": "BenchmarkMerge/changes=1000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 3192693,
            "unit": "ns/op",
            "extra": "379 times\n4 procs"
          },
          {
            "name": "BenchmarkMerge/changes=1000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 3533105,
            "unit": "B/op",
            "extra": "379 times\n4 procs"
          },
          {
            "name": "BenchmarkMerge/changes=1000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 43780,
            "unit": "allocs/op",
            "extra": "379 times\n4 procs"
          },
          {
            "name": "BenchmarkList/build/len=100 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 68544,
            "unit": "ns/op\t   78912 B/op\t     488 allocs/op",
            "extra": "17422 times\n4 procs"
          },
          {
            "name": "BenchmarkList/build/len=100 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 68544,
            "unit": "ns/op",
            "extra": "17422 times\n4 procs"
          },
          {
            "name": "BenchmarkList/build/len=100 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 78912,
            "unit": "B/op",
            "extra": "17422 times\n4 procs"
          },
          {
            "name": "BenchmarkList/build/len=100 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 488,
            "unit": "allocs/op",
            "extra": "17422 times\n4 procs"
          },
          {
            "name": "BenchmarkList/build/len=1000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 784469,
            "unit": "ns/op\t  718097 B/op\t    4007 allocs/op",
            "extra": "1610 times\n4 procs"
          },
          {
            "name": "BenchmarkList/build/len=1000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 784469,
            "unit": "ns/op",
            "extra": "1610 times\n4 procs"
          },
          {
            "name": "BenchmarkList/build/len=1000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 718097,
            "unit": "B/op",
            "extra": "1610 times\n4 procs"
          },
          {
            "name": "BenchmarkList/build/len=1000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 4007,
            "unit": "allocs/op",
            "extra": "1610 times\n4 procs"
          },
          {
            "name": "BenchmarkList/build/len=10000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 8011291,
            "unit": "ns/op\t10012423 B/op\t   40146 allocs/op",
            "extra": "146 times\n4 procs"
          },
          {
            "name": "BenchmarkList/build/len=10000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 8011291,
            "unit": "ns/op",
            "extra": "146 times\n4 procs"
          },
          {
            "name": "BenchmarkList/build/len=10000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 10012423,
            "unit": "B/op",
            "extra": "146 times\n4 procs"
          },
          {
            "name": "BenchmarkList/build/len=10000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 40146,
            "unit": "allocs/op",
            "extra": "146 times\n4 procs"
          },
          {
            "name": "BenchmarkList/save/len=100 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 25893,
            "unit": "ns/op\t  10.85 MB/s\t    6976 B/op\t     179 allocs/op",
            "extra": "43467 times\n4 procs"
          },
          {
            "name": "BenchmarkList/save/len=100 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 25893,
            "unit": "ns/op",
            "extra": "43467 times\n4 procs"
          },
          {
            "name": "BenchmarkList/save/len=100 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 10.85,
            "unit": "MB/s",
            "extra": "43467 times\n4 procs"
          },
          {
            "name": "BenchmarkList/save/len=100 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 6976,
            "unit": "B/op",
            "extra": "43467 times\n4 procs"
          },
          {
            "name": "BenchmarkList/save/len=100 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 179,
            "unit": "allocs/op",
            "extra": "43467 times\n4 procs"
          },
          {
            "name": "BenchmarkList/save/len=1000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 192353,
            "unit": "ns/op\t   8.70 MB/s\t   29264 B/op\t    1087 allocs/op",
            "extra": "6128 times\n4 procs"
          },
          {
            "name": "BenchmarkList/save/len=1000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 192353,
            "unit": "ns/op",
            "extra": "6128 times\n4 procs"
          },
          {
            "name": "BenchmarkList/save/len=1000 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 8.7,
            "unit": "MB/s",
            "extra": "6128 times\n4 procs"
          },
          {
            "name": "BenchmarkList/save/len=1000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 29264,
            "unit": "B/op",
            "extra": "6128 times\n4 procs"
          },
          {
            "name": "BenchmarkList/save/len=1000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 1087,
            "unit": "allocs/op",
            "extra": "6128 times\n4 procs"
          },
          {
            "name": "BenchmarkList/save/len=10000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 1349947,
            "unit": "ns/op\t  14.06 MB/s\t  305252 B/op\t   10096 allocs/op",
            "extra": "885 times\n4 procs"
          },
          {
            "name": "BenchmarkList/save/len=10000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 1349947,
            "unit": "ns/op",
            "extra": "885 times\n4 procs"
          },
          {
            "name": "BenchmarkList/save/len=10000 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 14.06,
            "unit": "MB/s",
            "extra": "885 times\n4 procs"
          },
          {
            "name": "BenchmarkList/save/len=10000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 305252,
            "unit": "B/op",
            "extra": "885 times\n4 procs"
          },
          {
            "name": "BenchmarkList/save/len=10000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 10096,
            "unit": "allocs/op",
            "extra": "885 times\n4 procs"
          },
          {
            "name": "BenchmarkList/load/len=100 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 29431,
            "unit": "ns/op\t   9.55 MB/s\t   12336 B/op\t     221 allocs/op",
            "extra": "40015 times\n4 procs"
          },
          {
            "name": "BenchmarkList/load/len=100 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 29431,
            "unit": "ns/op",
            "extra": "40015 times\n4 procs"
          },
          {
            "name": "BenchmarkList/load/len=100 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 9.55,
            "unit": "MB/s",
            "extra": "40015 times\n4 procs"
          },
          {
            "name": "BenchmarkList/load/len=100 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 12336,
            "unit": "B/op",
            "extra": "40015 times\n4 procs"
          },
          {
            "name": "BenchmarkList/load/len=100 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 221,
            "unit": "allocs/op",
            "extra": "40015 times\n4 procs"
          },
          {
            "name": "BenchmarkList/load/len=1000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 287282,
            "unit": "ns/op\t   5.82 MB/s\t  149232 B/op\t    2116 allocs/op",
            "extra": "4146 times\n4 procs"
          },
          {
            "name": "BenchmarkList/load/len=1000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 287282,
            "unit": "ns/op",
            "extra": "4146 times\n4 procs"
          },
          {
            "name": "BenchmarkList/load/len=1000 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 5.82,
            "unit": "MB/s",
            "extra": "4146 times\n4 procs"
          },
          {
            "name": "BenchmarkList/load/len=1000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 149232,
            "unit": "B/op",
            "extra": "4146 times\n4 procs"
          },
          {
            "name": "BenchmarkList/load/len=1000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 2116,
            "unit": "allocs/op",
            "extra": "4146 times\n4 procs"
          },
          {
            "name": "BenchmarkList/load/len=10000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 2698405,
            "unit": "ns/op\t   7.04 MB/s\t 1016936 B/op\t   22165 allocs/op",
            "extra": "428 times\n4 procs"
          },
          {
            "name": "BenchmarkList/load/len=10000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 2698405,
            "unit": "ns/op",
            "extra": "428 times\n4 procs"
          },
          {
            "name": "BenchmarkList/load/len=10000 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 7.04,
            "unit": "MB/s",
            "extra": "428 times\n4 procs"
          },
          {
            "name": "BenchmarkList/load/len=10000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 1016936,
            "unit": "B/op",
            "extra": "428 times\n4 procs"
          },
          {
            "name": "BenchmarkList/load/len=10000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 22165,
            "unit": "allocs/op",
            "extra": "428 times\n4 procs"
          },
          {
            "name": "BenchmarkList/iterate/len=100 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 7284,
            "unit": "ns/op\t   12992 B/op\t     101 allocs/op",
            "extra": "149740 times\n4 procs"
          },
          {
            "name": "BenchmarkList/iterate/len=100 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 7284,
            "unit": "ns/op",
            "extra": "149740 times\n4 procs"
          },
          {
            "name": "BenchmarkList/iterate/len=100 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 12992,
            "unit": "B/op",
            "extra": "149740 times\n4 procs"
          },
          {
            "name": "BenchmarkList/iterate/len=100 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 101,
            "unit": "allocs/op",
            "extra": "149740 times\n4 procs"
          },
          {
            "name": "BenchmarkList/iterate/len=1000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 71430,
            "unit": "ns/op\t  121796 B/op\t    1001 allocs/op",
            "extra": "16629 times\n4 procs"
          },
          {
            "name": "BenchmarkList/iterate/len=1000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 71430,
            "unit": "ns/op",
            "extra": "16629 times\n4 procs"
          },
          {
            "name": "BenchmarkList/iterate/len=1000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 121796,
            "unit": "B/op",
            "extra": "16629 times\n4 procs"
          },
          {
            "name": "BenchmarkList/iterate/len=1000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 1001,
            "unit": "allocs/op",
            "extra": "16629 times\n4 procs"
          },
          {
            "name": "BenchmarkList/iterate/len=10000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 755003,
            "unit": "ns/op\t 1207942 B/op\t   10023 allocs/op",
            "extra": "1330 times\n4 procs"
          },
          {
            "name": "BenchmarkList/iterate/len=10000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 755003,
            "unit": "ns/op",
            "extra": "1330 times\n4 procs"
          },
          {
            "name": "BenchmarkList/iterate/len=10000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 1207942,
            "unit": "B/op",
            "extra": "1330 times\n4 procs"
          },
          {
            "name": "BenchmarkList/iterate/len=10000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 10023,
            "unit": "allocs/op",
            "extra": "1330 times\n4 procs"
          },
          {
            "name": "BenchmarkList/as_slice/len=1000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 190127,
            "unit": "ns/op\t  170149 B/op\t    4004 allocs/op",
            "extra": "6246 times\n4 procs"
          },
          {
            "name": "BenchmarkList/as_slice/len=1000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 190127,
            "unit": "ns/op",
            "extra": "6246 times\n4 procs"
          },
          {
            "name": "BenchmarkList/as_slice/len=1000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 170149,
            "unit": "B/op",
            "extra": "6246 times\n4 procs"
          },
          {
            "name": "BenchmarkList/as_slice/len=1000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 4004,
            "unit": "allocs/op",
            "extra": "6246 times\n4 procs"
          },
          {
            "name": "BenchmarkTextRead/poorly_simulated_typing/n=1000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 4198,
            "unit": "ns/op\t    1064 B/op\t       3 allocs/op",
            "extra": "295255 times\n4 procs"
          },
          {
            "name": "BenchmarkTextRead/poorly_simulated_typing/n=1000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 4198,
            "unit": "ns/op",
            "extra": "295255 times\n4 procs"
          },
          {
            "name": "BenchmarkTextRead/poorly_simulated_typing/n=1000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 1064,
            "unit": "B/op",
            "extra": "295255 times\n4 procs"
          },
          {
            "name": "BenchmarkTextRead/poorly_simulated_typing/n=1000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "295255 times\n4 procs"
          },
          {
            "name": "BenchmarkTextRead/edit_trace (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 486303,
            "unit": "ns/op\t  106520 B/op\t       2 allocs/op",
            "extra": "2454 times\n4 procs"
          },
          {
            "name": "BenchmarkTextRead/edit_trace (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 486303,
            "unit": "ns/op",
            "extra": "2454 times\n4 procs"
          },
          {
            "name": "BenchmarkTextRead/edit_trace (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 106520,
            "unit": "B/op",
            "extra": "2454 times\n4 procs"
          },
          {
            "name": "BenchmarkTextRead/edit_trace (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "2454 times\n4 procs"
          },
          {
            "name": "BenchmarkSaveIncremental/incremental (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 3176,
            "unit": "ns/op\t 304.82 MB/s\t    2512 B/op\t       6 allocs/op",
            "extra": "405868 times\n4 procs"
          },
          {
            "name": "BenchmarkSaveIncremental/incremental (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 3176,
            "unit": "ns/op",
            "extra": "405868 times\n4 procs"
          },
          {
            "name": "BenchmarkSaveIncremental/incremental (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 304.82,
            "unit": "MB/s",
            "extra": "405868 times\n4 procs"
          },
          {
            "name": "BenchmarkSaveIncremental/incremental (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 2512,
            "unit": "B/op",
            "extra": "405868 times\n4 procs"
          },
          {
            "name": "BenchmarkSaveIncremental/incremental (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 6,
            "unit": "allocs/op",
            "extra": "405868 times\n4 procs"
          },
          {
            "name": "BenchmarkSaveIncremental/full (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 17605,
            "unit": "ns/op\t   9.49 MB/s\t    6480 B/op\t      94 allocs/op",
            "extra": "73194 times\n4 procs"
          },
          {
            "name": "BenchmarkSaveIncremental/full (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 17605,
            "unit": "ns/op",
            "extra": "73194 times\n4 procs"
          },
          {
            "name": "BenchmarkSaveIncremental/full (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 9.49,
            "unit": "MB/s",
            "extra": "73194 times\n4 procs"
          },
          {
            "name": "BenchmarkSaveIncremental/full (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 6480,
            "unit": "B/op",
            "extra": "73194 times\n4 procs"
          },
          {
            "name": "BenchmarkSaveIncremental/full (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 94,
            "unit": "allocs/op",
            "extra": "73194 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteChange/ops=100 (github.com/MichaelMure/gotomerge/format)",
            "value": 67666,
            "unit": "ns/op\t   5.44 MB/s\t   15032 B/op\t     181 allocs/op",
            "extra": "17917 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteChange/ops=100 (github.com/MichaelMure/gotomerge/format) - ns/op",
            "value": 67666,
            "unit": "ns/op",
            "extra": "17917 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteChange/ops=100 (github.com/MichaelMure/gotomerge/format) - MB/s",
            "value": 5.44,
            "unit": "MB/s",
            "extra": "17917 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteChange/ops=100 (github.com/MichaelMure/gotomerge/format) - B/op",
            "value": 15032,
            "unit": "B/op",
            "extra": "17917 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteChange/ops=100 (github.com/MichaelMure/gotomerge/format) - allocs/op",
            "value": 181,
            "unit": "allocs/op",
            "extra": "17917 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteChange/ops=1000 (github.com/MichaelMure/gotomerge/format)",
            "value": 427517,
            "unit": "ns/op\t   9.05 MB/s\t  116620 B/op\t    1094 allocs/op",
            "extra": "2770 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteChange/ops=1000 (github.com/MichaelMure/gotomerge/format) - ns/op",
            "value": 427517,
            "unit": "ns/op",
            "extra": "2770 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteChange/ops=1000 (github.com/MichaelMure/gotomerge/format) - MB/s",
            "value": 9.05,
            "unit": "MB/s",
            "extra": "2770 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteChange/ops=1000 (github.com/MichaelMure/gotomerge/format) - B/op",
            "value": 116620,
            "unit": "B/op",
            "extra": "2770 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteChange/ops=1000 (github.com/MichaelMure/gotomerge/format) - allocs/op",
            "value": 1094,
            "unit": "allocs/op",
            "extra": "2770 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteChange/ops=10000 (github.com/MichaelMure/gotomerge/format)",
            "value": 4185362,
            "unit": "ns/op\t   9.89 MB/s\t 1901531 B/op\t   10226 allocs/op",
            "extra": "285 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteChange/ops=10000 (github.com/MichaelMure/gotomerge/format) - ns/op",
            "value": 4185362,
            "unit": "ns/op",
            "extra": "285 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteChange/ops=10000 (github.com/MichaelMure/gotomerge/format) - MB/s",
            "value": 9.89,
            "unit": "MB/s",
            "extra": "285 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteChange/ops=10000 (github.com/MichaelMure/gotomerge/format) - B/op",
            "value": 1901531,
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
            "value": 9444,
            "unit": "ns/op\t  42.99 MB/s\t         1.000 chunks\t    2888 B/op\t      46 allocs/op",
            "extra": "124832 times\n4 procs"
          },
          {
            "name": "BenchmarkReadExamplar (github.com/MichaelMure/gotomerge/format) - ns/op",
            "value": 9444,
            "unit": "ns/op",
            "extra": "124832 times\n4 procs"
          },
          {
            "name": "BenchmarkReadExamplar (github.com/MichaelMure/gotomerge/format) - MB/s",
            "value": 42.99,
            "unit": "MB/s",
            "extra": "124832 times\n4 procs"
          },
          {
            "name": "BenchmarkReadExamplar (github.com/MichaelMure/gotomerge/format) - chunks",
            "value": 1,
            "unit": "chunks",
            "extra": "124832 times\n4 procs"
          },
          {
            "name": "BenchmarkReadExamplar (github.com/MichaelMure/gotomerge/format) - B/op",
            "value": 2888,
            "unit": "B/op",
            "extra": "124832 times\n4 procs"
          },
          {
            "name": "BenchmarkReadExamplar (github.com/MichaelMure/gotomerge/format) - allocs/op",
            "value": 46,
            "unit": "allocs/op",
            "extra": "124832 times\n4 procs"
          },
          {
            "name": "BenchmarkReadLarge (github.com/MichaelMure/gotomerge/format)",
            "value": 388110045,
            "unit": "ns/op\t  75.36 MB/s\t    259779 chunks\t462359632 B/op\t 5792673 allocs/op",
            "extra": "3 times\n4 procs"
          },
          {
            "name": "BenchmarkReadLarge (github.com/MichaelMure/gotomerge/format) - ns/op",
            "value": 388110045,
            "unit": "ns/op",
            "extra": "3 times\n4 procs"
          },
          {
            "name": "BenchmarkReadLarge (github.com/MichaelMure/gotomerge/format) - MB/s",
            "value": 75.36,
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
            "value": 462359632,
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
            "value": 27452,
            "unit": "ns/op\t   6.99 MB/s\t         2.000 chunks\t   85768 B/op\t      66 allocs/op",
            "extra": "42775 times\n4 procs"
          },
          {
            "name": "BenchmarkCompressed (github.com/MichaelMure/gotomerge/format) - ns/op",
            "value": 27452,
            "unit": "ns/op",
            "extra": "42775 times\n4 procs"
          },
          {
            "name": "BenchmarkCompressed (github.com/MichaelMure/gotomerge/format) - MB/s",
            "value": 6.99,
            "unit": "MB/s",
            "extra": "42775 times\n4 procs"
          },
          {
            "name": "BenchmarkCompressed (github.com/MichaelMure/gotomerge/format) - chunks",
            "value": 2,
            "unit": "chunks",
            "extra": "42775 times\n4 procs"
          },
          {
            "name": "BenchmarkCompressed (github.com/MichaelMure/gotomerge/format) - B/op",
            "value": 85768,
            "unit": "B/op",
            "extra": "42775 times\n4 procs"
          },
          {
            "name": "BenchmarkCompressed (github.com/MichaelMure/gotomerge/format) - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "42775 times\n4 procs"
          },
          {
            "name": "BenchmarkTextEdits (github.com/MichaelMure/gotomerge/opset)",
            "value": 652392,
            "unit": "ns/op\t    104852 chars\t  177154 B/op\t     295 allocs/op",
            "extra": "1870 times\n4 procs"
          },
          {
            "name": "BenchmarkTextEdits (github.com/MichaelMure/gotomerge/opset) - ns/op",
            "value": 652392,
            "unit": "ns/op",
            "extra": "1870 times\n4 procs"
          },
          {
            "name": "BenchmarkTextEdits (github.com/MichaelMure/gotomerge/opset) - chars",
            "value": 104852,
            "unit": "chars",
            "extra": "1870 times\n4 procs"
          },
          {
            "name": "BenchmarkTextEdits (github.com/MichaelMure/gotomerge/opset) - B/op",
            "value": 177154,
            "unit": "B/op",
            "extra": "1870 times\n4 procs"
          },
          {
            "name": "BenchmarkTextEdits (github.com/MichaelMure/gotomerge/opset) - allocs/op",
            "value": 295,
            "unit": "allocs/op",
            "extra": "1870 times\n4 procs"
          }
        ]
      },
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
          "id": "f14d2dab98e976d9e916436a66c69308fd7a7f11",
          "message": "ci: more explicit names",
          "timestamp": "2026-04-21T10:41:06+02:00",
          "tree_id": "5fd934dd61c2d10b5303cab8523dc5a0b3f5d29d",
          "url": "https://github.com/MichaelMure/gotomerge/commit/f14d2dab98e976d9e916436a66c69308fd7a7f11"
        },
        "date": 1776761204336,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkReadMetadata (github.com/MichaelMure/gotomerge/column)",
            "value": 144.4,
            "unit": "ns/op\t      80 B/op\t       2 allocs/op",
            "extra": "8089293 times\n4 procs"
          },
          {
            "name": "BenchmarkReadMetadata (github.com/MichaelMure/gotomerge/column) - ns/op",
            "value": 144.4,
            "unit": "ns/op",
            "extra": "8089293 times\n4 procs"
          },
          {
            "name": "BenchmarkReadMetadata (github.com/MichaelMure/gotomerge/column) - B/op",
            "value": 80,
            "unit": "B/op",
            "extra": "8089293 times\n4 procs"
          },
          {
            "name": "BenchmarkReadMetadata (github.com/MichaelMure/gotomerge/column) - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "8089293 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteMetadata (github.com/MichaelMure/gotomerge/column)",
            "value": 78.42,
            "unit": "ns/op\t      40 B/op\t       5 allocs/op",
            "extra": "14838576 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteMetadata (github.com/MichaelMure/gotomerge/column) - ns/op",
            "value": 78.42,
            "unit": "ns/op",
            "extra": "14838576 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteMetadata (github.com/MichaelMure/gotomerge/column) - B/op",
            "value": 40,
            "unit": "B/op",
            "extra": "14838576 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteMetadata (github.com/MichaelMure/gotomerge/column) - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "14838576 times\n4 procs"
          },
          {
            "name": "BenchmarkReadSpecification (github.com/MichaelMure/gotomerge/column)",
            "value": 143.1,
            "unit": "ns/op\t      80 B/op\t       2 allocs/op",
            "extra": "8259171 times\n4 procs"
          },
          {
            "name": "BenchmarkReadSpecification (github.com/MichaelMure/gotomerge/column) - ns/op",
            "value": 143.1,
            "unit": "ns/op",
            "extra": "8259171 times\n4 procs"
          },
          {
            "name": "BenchmarkReadSpecification (github.com/MichaelMure/gotomerge/column) - B/op",
            "value": 80,
            "unit": "B/op",
            "extra": "8259171 times\n4 procs"
          },
          {
            "name": "BenchmarkReadSpecification (github.com/MichaelMure/gotomerge/column) - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "8259171 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteSpecification (github.com/MichaelMure/gotomerge/column)",
            "value": 39.03,
            "unit": "ns/op\t      16 B/op\t       2 allocs/op",
            "extra": "30326188 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteSpecification (github.com/MichaelMure/gotomerge/column) - ns/op",
            "value": 39.03,
            "unit": "ns/op",
            "extra": "30326188 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteSpecification (github.com/MichaelMure/gotomerge/column) - B/op",
            "value": 16,
            "unit": "B/op",
            "extra": "30326188 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteSpecification (github.com/MichaelMure/gotomerge/column) - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "30326188 times\n4 procs"
          },
          {
            "name": "BenchmarkUint64Reader (github.com/MichaelMure/gotomerge/column/rle)",
            "value": 117.1,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "10149675 times\n4 procs"
          },
          {
            "name": "BenchmarkUint64Reader (github.com/MichaelMure/gotomerge/column/rle) - ns/op",
            "value": 117.1,
            "unit": "ns/op",
            "extra": "10149675 times\n4 procs"
          },
          {
            "name": "BenchmarkUint64Reader (github.com/MichaelMure/gotomerge/column/rle) - B/op",
            "value": 80,
            "unit": "B/op",
            "extra": "10149675 times\n4 procs"
          },
          {
            "name": "BenchmarkUint64Reader (github.com/MichaelMure/gotomerge/column/rle) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "10149675 times\n4 procs"
          },
          {
            "name": "BenchmarkInt64Reader (github.com/MichaelMure/gotomerge/column/rle)",
            "value": 116.6,
            "unit": "ns/op\t      80 B/op\t       1 allocs/op",
            "extra": "10093276 times\n4 procs"
          },
          {
            "name": "BenchmarkInt64Reader (github.com/MichaelMure/gotomerge/column/rle) - ns/op",
            "value": 116.6,
            "unit": "ns/op",
            "extra": "10093276 times\n4 procs"
          },
          {
            "name": "BenchmarkInt64Reader (github.com/MichaelMure/gotomerge/column/rle) - B/op",
            "value": 80,
            "unit": "B/op",
            "extra": "10093276 times\n4 procs"
          },
          {
            "name": "BenchmarkInt64Reader (github.com/MichaelMure/gotomerge/column/rle) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "10093276 times\n4 procs"
          },
          {
            "name": "BenchmarkStringReader (github.com/MichaelMure/gotomerge/column/rle)",
            "value": 167.2,
            "unit": "ns/op\t     104 B/op\t       4 allocs/op",
            "extra": "7206819 times\n4 procs"
          },
          {
            "name": "BenchmarkStringReader (github.com/MichaelMure/gotomerge/column/rle) - ns/op",
            "value": 167.2,
            "unit": "ns/op",
            "extra": "7206819 times\n4 procs"
          },
          {
            "name": "BenchmarkStringReader (github.com/MichaelMure/gotomerge/column/rle) - B/op",
            "value": 104,
            "unit": "B/op",
            "extra": "7206819 times\n4 procs"
          },
          {
            "name": "BenchmarkStringReader (github.com/MichaelMure/gotomerge/column/rle) - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "7206819 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/repeated_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 41106,
            "unit": "ns/op\t   61600 B/op\t     354 allocs/op",
            "extra": "29040 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/repeated_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 41106,
            "unit": "ns/op",
            "extra": "29040 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/repeated_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 61600,
            "unit": "B/op",
            "extra": "29040 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/repeated_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 354,
            "unit": "allocs/op",
            "extra": "29040 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/repeated_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 580694,
            "unit": "ns/op\t  707786 B/op\t    2967 allocs/op",
            "extra": "2187 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/repeated_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 580694,
            "unit": "ns/op",
            "extra": "2187 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/repeated_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 707786,
            "unit": "B/op",
            "extra": "2187 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/repeated_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 2967,
            "unit": "allocs/op",
            "extra": "2187 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/repeated_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 9347978,
            "unit": "ns/op\t 9102213 B/op\t   30069 allocs/op",
            "extra": "133 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/repeated_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 9347978,
            "unit": "ns/op",
            "extra": "133 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/repeated_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 9102213,
            "unit": "B/op",
            "extra": "133 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/repeated_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 30069,
            "unit": "allocs/op",
            "extra": "133 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/increasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 153777,
            "unit": "ns/op\t  157001 B/op\t     479 allocs/op",
            "extra": "7962 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/increasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 153777,
            "unit": "ns/op",
            "extra": "7962 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/increasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 157001,
            "unit": "B/op",
            "extra": "7962 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/increasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 479,
            "unit": "allocs/op",
            "extra": "7962 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/increasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 1012842,
            "unit": "ns/op\t 1101749 B/op\t    4946 allocs/op",
            "extra": "1222 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/increasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 1012842,
            "unit": "ns/op",
            "extra": "1222 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/increasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 1101749,
            "unit": "B/op",
            "extra": "1222 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/increasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 4946,
            "unit": "allocs/op",
            "extra": "1222 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/increasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 13348568,
            "unit": "ns/op\t11785012 B/op\t   50066 allocs/op",
            "extra": "99 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/increasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 13348568,
            "unit": "ns/op",
            "extra": "99 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/increasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 11785012,
            "unit": "B/op",
            "extra": "99 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/increasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 50066,
            "unit": "allocs/op",
            "extra": "99 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/decreasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 147221,
            "unit": "ns/op\t  157780 B/op\t     479 allocs/op",
            "extra": "8035 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/decreasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 147221,
            "unit": "ns/op",
            "extra": "8035 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/decreasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 157780,
            "unit": "B/op",
            "extra": "8035 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/decreasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 479,
            "unit": "allocs/op",
            "extra": "8035 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/decreasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 953844,
            "unit": "ns/op\t 1117007 B/op\t    4947 allocs/op",
            "extra": "1207 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/decreasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 953844,
            "unit": "ns/op",
            "extra": "1207 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/decreasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 1117007,
            "unit": "B/op",
            "extra": "1207 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/decreasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 4947,
            "unit": "allocs/op",
            "extra": "1207 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/decreasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 13391761,
            "unit": "ns/op\t11785266 B/op\t   50072 allocs/op",
            "extra": "82 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/decreasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 13391761,
            "unit": "ns/op",
            "extra": "82 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/decreasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 11785266,
            "unit": "B/op",
            "extra": "82 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/build/decreasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 50072,
            "unit": "allocs/op",
            "extra": "82 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/repeated_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 19483,
            "unit": "ns/op\t  12.32 MB/s\t    5952 B/op\t     167 allocs/op",
            "extra": "61425 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/repeated_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 19483,
            "unit": "ns/op",
            "extra": "61425 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/repeated_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 12.32,
            "unit": "MB/s",
            "extra": "61425 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/repeated_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 5952,
            "unit": "B/op",
            "extra": "61425 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/repeated_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 167,
            "unit": "allocs/op",
            "extra": "61425 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/repeated_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 223481,
            "unit": "ns/op\t   7.30 MB/s\t  217086 B/op\t    1079 allocs/op",
            "extra": "5414 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/repeated_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 223481,
            "unit": "ns/op",
            "extra": "5414 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/repeated_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 7.3,
            "unit": "MB/s",
            "extra": "5414 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/repeated_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 217086,
            "unit": "B/op",
            "extra": "5414 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/repeated_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 1079,
            "unit": "allocs/op",
            "extra": "5414 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/repeated_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 1534504,
            "unit": "ns/op\t  12.34 MB/s\t 1077354 B/op\t   10101 allocs/op",
            "extra": "816 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/repeated_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 1534504,
            "unit": "ns/op",
            "extra": "816 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/repeated_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 12.34,
            "unit": "MB/s",
            "extra": "816 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/repeated_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 1077354,
            "unit": "B/op",
            "extra": "816 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/repeated_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 10101,
            "unit": "allocs/op",
            "extra": "816 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/increasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 73630,
            "unit": "ns/op\t   5.22 MB/s\t   68098 B/op\t     181 allocs/op",
            "extra": "16047 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/increasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 73630,
            "unit": "ns/op",
            "extra": "16047 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/increasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 5.22,
            "unit": "MB/s",
            "extra": "16047 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/increasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 68098,
            "unit": "B/op",
            "extra": "16047 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/increasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 181,
            "unit": "allocs/op",
            "extra": "16047 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/increasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 354919,
            "unit": "ns/op\t   9.44 MB/s\t  398612 B/op\t    1103 allocs/op",
            "extra": "3532 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/increasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 354919,
            "unit": "ns/op",
            "extra": "3532 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/increasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 9.44,
            "unit": "MB/s",
            "extra": "3532 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/increasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 398612,
            "unit": "B/op",
            "extra": "3532 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/increasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 1103,
            "unit": "allocs/op",
            "extra": "3532 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/increasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 2810924,
            "unit": "ns/op\t  12.85 MB/s\t 2032560 B/op\t   10137 allocs/op",
            "extra": "415 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/increasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 2810924,
            "unit": "ns/op",
            "extra": "415 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/increasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 12.85,
            "unit": "MB/s",
            "extra": "415 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/increasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 2032560,
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
            "value": 74634,
            "unit": "ns/op\t   5.15 MB/s\t   70157 B/op\t     181 allocs/op",
            "extra": "15532 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/decreasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 74634,
            "unit": "ns/op",
            "extra": "15532 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/decreasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 5.15,
            "unit": "MB/s",
            "extra": "15532 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/decreasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 70157,
            "unit": "B/op",
            "extra": "15532 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/decreasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 181,
            "unit": "allocs/op",
            "extra": "15532 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/decreasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 358275,
            "unit": "ns/op\t   9.35 MB/s\t  429430 B/op\t    1103 allocs/op",
            "extra": "3379 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/decreasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 358275,
            "unit": "ns/op",
            "extra": "3379 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/decreasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 9.35,
            "unit": "MB/s",
            "extra": "3379 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/decreasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 429430,
            "unit": "B/op",
            "extra": "3379 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/decreasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 1103,
            "unit": "allocs/op",
            "extra": "3379 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/decreasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 2851959,
            "unit": "ns/op\t  12.56 MB/s\t 2032572 B/op\t   10137 allocs/op",
            "extra": "421 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/decreasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 2851959,
            "unit": "ns/op",
            "extra": "421 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/decreasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 12.56,
            "unit": "MB/s",
            "extra": "421 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/decreasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 2032572,
            "unit": "B/op",
            "extra": "421 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/save/decreasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 10137,
            "unit": "allocs/op",
            "extra": "421 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/repeated_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 26895,
            "unit": "ns/op\t   8.92 MB/s\t   13048 B/op\t     208 allocs/op",
            "extra": "44438 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/repeated_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 26895,
            "unit": "ns/op",
            "extra": "44438 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/repeated_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 8.92,
            "unit": "MB/s",
            "extra": "44438 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/repeated_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 13048,
            "unit": "B/op",
            "extra": "44438 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/repeated_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 208,
            "unit": "allocs/op",
            "extra": "44438 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/repeated_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 287571,
            "unit": "ns/op\t   5.67 MB/s\t  166360 B/op\t    2062 allocs/op",
            "extra": "4268 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/repeated_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 287571,
            "unit": "ns/op",
            "extra": "4268 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/repeated_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 5.67,
            "unit": "MB/s",
            "extra": "4268 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/repeated_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 166360,
            "unit": "B/op",
            "extra": "4268 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/repeated_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 2062,
            "unit": "allocs/op",
            "extra": "4268 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/repeated_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 2718222,
            "unit": "ns/op\t   6.97 MB/s\t 1205668 B/op\t   21690 allocs/op",
            "extra": "438 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/repeated_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 2718222,
            "unit": "ns/op",
            "extra": "438 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/repeated_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 6.97,
            "unit": "MB/s",
            "extra": "438 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/repeated_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 1205668,
            "unit": "B/op",
            "extra": "438 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/repeated_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 21690,
            "unit": "allocs/op",
            "extra": "438 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/increasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 52005,
            "unit": "ns/op\t   7.38 MB/s\t   65024 B/op\t     506 allocs/op",
            "extra": "23288 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/increasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 52005,
            "unit": "ns/op",
            "extra": "23288 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/increasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 7.38,
            "unit": "MB/s",
            "extra": "23288 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/increasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 65024,
            "unit": "B/op",
            "extra": "23288 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/increasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 506,
            "unit": "allocs/op",
            "extra": "23288 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/increasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 494291,
            "unit": "ns/op\t   6.78 MB/s\t  413073 B/op\t    5077 allocs/op",
            "extra": "2360 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/increasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 494291,
            "unit": "ns/op",
            "extra": "2360 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/increasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 6.78,
            "unit": "MB/s",
            "extra": "2360 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/increasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 413073,
            "unit": "B/op",
            "extra": "2360 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/increasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 5077,
            "unit": "allocs/op",
            "extra": "2360 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/increasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 4705071,
            "unit": "ns/op\t   7.68 MB/s\t 3003306 B/op\t   51767 allocs/op",
            "extra": "254 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/increasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 4705071,
            "unit": "ns/op",
            "extra": "254 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/increasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 7.68,
            "unit": "MB/s",
            "extra": "254 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/increasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 3003306,
            "unit": "B/op",
            "extra": "254 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/increasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 51767,
            "unit": "allocs/op",
            "extra": "254 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/decreasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 51485,
            "unit": "ns/op\t   7.46 MB/s\t   65024 B/op\t     506 allocs/op",
            "extra": "23164 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/decreasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 51485,
            "unit": "ns/op",
            "extra": "23164 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/decreasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 7.46,
            "unit": "MB/s",
            "extra": "23164 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/decreasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 65024,
            "unit": "B/op",
            "extra": "23164 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/decreasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 506,
            "unit": "allocs/op",
            "extra": "23164 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/decreasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 492403,
            "unit": "ns/op\t   6.81 MB/s\t  413073 B/op\t    5077 allocs/op",
            "extra": "2390 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/decreasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 492403,
            "unit": "ns/op",
            "extra": "2390 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/decreasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 6.81,
            "unit": "MB/s",
            "extra": "2390 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/decreasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 413073,
            "unit": "B/op",
            "extra": "2390 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/decreasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 5077,
            "unit": "allocs/op",
            "extra": "2390 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/decreasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 4705978,
            "unit": "ns/op\t   7.61 MB/s\t 3003506 B/op\t   51767 allocs/op",
            "extra": "250 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/decreasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 4705978,
            "unit": "ns/op",
            "extra": "250 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/decreasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 7.61,
            "unit": "MB/s",
            "extra": "250 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/decreasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 3003506,
            "unit": "B/op",
            "extra": "250 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/load/decreasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 51767,
            "unit": "allocs/op",
            "extra": "250 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/repeated_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 29570,
            "unit": "ns/op\t   6.76 MB/s\t   33984 B/op\t     197 allocs/op",
            "extra": "40050 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/repeated_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 29570,
            "unit": "ns/op",
            "extra": "40050 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/repeated_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 6.76,
            "unit": "MB/s",
            "extra": "40050 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/repeated_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 33984,
            "unit": "B/op",
            "extra": "40050 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/repeated_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 197,
            "unit": "allocs/op",
            "extra": "40050 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/repeated_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 299564,
            "unit": "ns/op\t   5.36 MB/s\t  331265 B/op\t    1900 allocs/op",
            "extra": "3841 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/repeated_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 299564,
            "unit": "ns/op",
            "extra": "3841 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/repeated_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 5.36,
            "unit": "MB/s",
            "extra": "3841 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/repeated_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 331265,
            "unit": "B/op",
            "extra": "3841 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/repeated_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 1900,
            "unit": "allocs/op",
            "extra": "3841 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/repeated_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 4448633,
            "unit": "ns/op\t   4.25 MB/s\t 4492039 B/op\t   19978 allocs/op",
            "extra": "264 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/repeated_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 4448633,
            "unit": "ns/op",
            "extra": "264 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/repeated_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 4.25,
            "unit": "MB/s",
            "extra": "264 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/repeated_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 4492039,
            "unit": "B/op",
            "extra": "264 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/repeated_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 19978,
            "unit": "allocs/op",
            "extra": "264 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/increasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 55464,
            "unit": "ns/op\t   6.74 MB/s\t   86120 B/op\t     498 allocs/op",
            "extra": "21715 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/increasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 55464,
            "unit": "ns/op",
            "extra": "21715 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/increasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 6.74,
            "unit": "MB/s",
            "extra": "21715 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/increasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 86120,
            "unit": "B/op",
            "extra": "21715 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/increasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 498,
            "unit": "allocs/op",
            "extra": "21715 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/increasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 480920,
            "unit": "ns/op\t   8.04 MB/s\t  542635 B/op\t    4950 allocs/op",
            "extra": "2468 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/increasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 480920,
            "unit": "ns/op",
            "extra": "2468 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/increasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 8.04,
            "unit": "MB/s",
            "extra": "2468 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/increasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 542635,
            "unit": "B/op",
            "extra": "2468 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/increasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 4950,
            "unit": "allocs/op",
            "extra": "2468 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/increasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 6895396,
            "unit": "ns/op\t   5.49 MB/s\t 6217734 B/op\t   50038 allocs/op",
            "extra": "172 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/increasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 6895396,
            "unit": "ns/op",
            "extra": "172 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/increasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 5.49,
            "unit": "MB/s",
            "extra": "172 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/increasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 6217734,
            "unit": "B/op",
            "extra": "172 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/increasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 50038,
            "unit": "allocs/op",
            "extra": "172 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/decreasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 55510,
            "unit": "ns/op\t   6.74 MB/s\t   86088 B/op\t     496 allocs/op",
            "extra": "21253 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/decreasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 55510,
            "unit": "ns/op",
            "extra": "21253 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/decreasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 6.74,
            "unit": "MB/s",
            "extra": "21253 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/decreasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 86088,
            "unit": "B/op",
            "extra": "21253 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/decreasing_put/ops=100 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 496,
            "unit": "allocs/op",
            "extra": "21253 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/decreasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 483687,
            "unit": "ns/op\t   8.01 MB/s\t  542635 B/op\t    4950 allocs/op",
            "extra": "2443 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/decreasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 483687,
            "unit": "ns/op",
            "extra": "2443 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/decreasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 8.01,
            "unit": "MB/s",
            "extra": "2443 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/decreasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 542635,
            "unit": "B/op",
            "extra": "2443 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/decreasing_put/ops=1000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 4950,
            "unit": "allocs/op",
            "extra": "2443 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/decreasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 6933209,
            "unit": "ns/op\t   5.32 MB/s\t 6217912 B/op\t   50042 allocs/op",
            "extra": "176 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/decreasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 6933209,
            "unit": "ns/op",
            "extra": "176 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/decreasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 5.32,
            "unit": "MB/s",
            "extra": "176 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/decreasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 6217912,
            "unit": "B/op",
            "extra": "176 times\n4 procs"
          },
          {
            "name": "BenchmarkMap/apply/decreasing_put/ops=10000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 50042,
            "unit": "allocs/op",
            "extra": "176 times\n4 procs"
          },
          {
            "name": "BenchmarkSaveLoad/big_paste (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 49481,
            "unit": "ns/op\t   74140 B/op\t     162 allocs/op",
            "extra": "24500 times\n4 procs"
          },
          {
            "name": "BenchmarkSaveLoad/big_paste (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 49481,
            "unit": "ns/op",
            "extra": "24500 times\n4 procs"
          },
          {
            "name": "BenchmarkSaveLoad/big_paste (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 74140,
            "unit": "B/op",
            "extra": "24500 times\n4 procs"
          },
          {
            "name": "BenchmarkSaveLoad/big_paste (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 162,
            "unit": "allocs/op",
            "extra": "24500 times\n4 procs"
          },
          {
            "name": "BenchmarkSaveLoad/maps_in_maps (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 1678912,
            "unit": "ns/op\t 1308463 B/op\t    6493 allocs/op",
            "extra": "721 times\n4 procs"
          },
          {
            "name": "BenchmarkSaveLoad/maps_in_maps (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 1678912,
            "unit": "ns/op",
            "extra": "721 times\n4 procs"
          },
          {
            "name": "BenchmarkSaveLoad/maps_in_maps (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 1308463,
            "unit": "B/op",
            "extra": "721 times\n4 procs"
          },
          {
            "name": "BenchmarkSaveLoad/maps_in_maps (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 6493,
            "unit": "allocs/op",
            "extra": "721 times\n4 procs"
          },
          {
            "name": "BenchmarkSaveLoad/deep_history (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 1389920,
            "unit": "ns/op\t  813036 B/op\t   10133 allocs/op",
            "extra": "844 times\n4 procs"
          },
          {
            "name": "BenchmarkSaveLoad/deep_history (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 1389920,
            "unit": "ns/op",
            "extra": "844 times\n4 procs"
          },
          {
            "name": "BenchmarkSaveLoad/deep_history (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 813036,
            "unit": "B/op",
            "extra": "844 times\n4 procs"
          },
          {
            "name": "BenchmarkSaveLoad/deep_history (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 10133,
            "unit": "allocs/op",
            "extra": "844 times\n4 procs"
          },
          {
            "name": "BenchmarkSaveLoad/poorly_simulated_typing (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 476208,
            "unit": "ns/op\t  248636 B/op\t    4433 allocs/op",
            "extra": "2481 times\n4 procs"
          },
          {
            "name": "BenchmarkSaveLoad/poorly_simulated_typing (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 476208,
            "unit": "ns/op",
            "extra": "2481 times\n4 procs"
          },
          {
            "name": "BenchmarkSaveLoad/poorly_simulated_typing (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 248636,
            "unit": "B/op",
            "extra": "2481 times\n4 procs"
          },
          {
            "name": "BenchmarkSaveLoad/poorly_simulated_typing (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 4433,
            "unit": "allocs/op",
            "extra": "2481 times\n4 procs"
          },
          {
            "name": "BenchmarkEditTrace/apply (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 957412828,
            "unit": "ns/op\t  30.55 MB/s\t1072416680 B/op\t12966370 allocs/op",
            "extra": "2 times\n4 procs"
          },
          {
            "name": "BenchmarkEditTrace/apply (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 957412828,
            "unit": "ns/op",
            "extra": "2 times\n4 procs"
          },
          {
            "name": "BenchmarkEditTrace/apply (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 30.55,
            "unit": "MB/s",
            "extra": "2 times\n4 procs"
          },
          {
            "name": "BenchmarkEditTrace/apply (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 1072416680,
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
            "value": 1132659600,
            "unit": "ns/op\t   0.12 MB/s\t1109293920 B/op\t13346971 allocs/op",
            "extra": "1 times\n4 procs"
          },
          {
            "name": "BenchmarkEditTrace/save (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 1132659600,
            "unit": "ns/op",
            "extra": "1 times\n4 procs"
          },
          {
            "name": "BenchmarkEditTrace/save (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 0.12,
            "unit": "MB/s",
            "extra": "1 times\n4 procs"
          },
          {
            "name": "BenchmarkEditTrace/save (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 1109293920,
            "unit": "B/op",
            "extra": "1 times\n4 procs"
          },
          {
            "name": "BenchmarkEditTrace/save (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 13346971,
            "unit": "allocs/op",
            "extra": "1 times\n4 procs"
          },
          {
            "name": "BenchmarkEditTrace/load (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 80960547,
            "unit": "ns/op\t   1.75 MB/s\t31288002 B/op\t  761472 allocs/op",
            "extra": "14 times\n4 procs"
          },
          {
            "name": "BenchmarkEditTrace/load (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 80960547,
            "unit": "ns/op",
            "extra": "14 times\n4 procs"
          },
          {
            "name": "BenchmarkEditTrace/load (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 1.75,
            "unit": "MB/s",
            "extra": "14 times\n4 procs"
          },
          {
            "name": "BenchmarkEditTrace/load (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 31288002,
            "unit": "B/op",
            "extra": "14 times\n4 procs"
          },
          {
            "name": "BenchmarkEditTrace/load (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 761472,
            "unit": "allocs/op",
            "extra": "14 times\n4 procs"
          },
          {
            "name": "BenchmarkRange (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 41560829,
            "unit": "ns/op\t48200968 B/op\t  100559 allocs/op",
            "extra": "33 times\n4 procs"
          },
          {
            "name": "BenchmarkRange (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 41560829,
            "unit": "ns/op",
            "extra": "33 times\n4 procs"
          },
          {
            "name": "BenchmarkRange (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 48200968,
            "unit": "B/op",
            "extra": "33 times\n4 procs"
          },
          {
            "name": "BenchmarkRange (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 100559,
            "unit": "allocs/op",
            "extra": "33 times\n4 procs"
          },
          {
            "name": "BenchmarkGet/scalar/keys=100 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 96.34,
            "unit": "ns/op\t      48 B/op\t       1 allocs/op",
            "extra": "11881119 times\n4 procs"
          },
          {
            "name": "BenchmarkGet/scalar/keys=100 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 96.34,
            "unit": "ns/op",
            "extra": "11881119 times\n4 procs"
          },
          {
            "name": "BenchmarkGet/scalar/keys=100 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 48,
            "unit": "B/op",
            "extra": "11881119 times\n4 procs"
          },
          {
            "name": "BenchmarkGet/scalar/keys=100 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "11881119 times\n4 procs"
          },
          {
            "name": "BenchmarkGet/scalar/keys=1000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 100.6,
            "unit": "ns/op\t      48 B/op\t       1 allocs/op",
            "extra": "12410354 times\n4 procs"
          },
          {
            "name": "BenchmarkGet/scalar/keys=1000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 100.6,
            "unit": "ns/op",
            "extra": "12410354 times\n4 procs"
          },
          {
            "name": "BenchmarkGet/scalar/keys=1000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 48,
            "unit": "B/op",
            "extra": "12410354 times\n4 procs"
          },
          {
            "name": "BenchmarkGet/scalar/keys=1000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "12410354 times\n4 procs"
          },
          {
            "name": "BenchmarkGet/scalar/keys=10000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 98.52,
            "unit": "ns/op\t      48 B/op\t       1 allocs/op",
            "extra": "12592791 times\n4 procs"
          },
          {
            "name": "BenchmarkGet/scalar/keys=10000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 98.52,
            "unit": "ns/op",
            "extra": "12592791 times\n4 procs"
          },
          {
            "name": "BenchmarkGet/scalar/keys=10000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 48,
            "unit": "B/op",
            "extra": "12592791 times\n4 procs"
          },
          {
            "name": "BenchmarkGet/scalar/keys=10000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "12592791 times\n4 procs"
          },
          {
            "name": "BenchmarkGet/nested/keys=100 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 92.54,
            "unit": "ns/op\t      48 B/op\t       1 allocs/op",
            "extra": "13422706 times\n4 procs"
          },
          {
            "name": "BenchmarkGet/nested/keys=100 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 92.54,
            "unit": "ns/op",
            "extra": "13422706 times\n4 procs"
          },
          {
            "name": "BenchmarkGet/nested/keys=100 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 48,
            "unit": "B/op",
            "extra": "13422706 times\n4 procs"
          },
          {
            "name": "BenchmarkGet/nested/keys=100 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "13422706 times\n4 procs"
          },
          {
            "name": "BenchmarkGet/nested/keys=1000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 94.79,
            "unit": "ns/op\t      48 B/op\t       1 allocs/op",
            "extra": "11676180 times\n4 procs"
          },
          {
            "name": "BenchmarkGet/nested/keys=1000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 94.79,
            "unit": "ns/op",
            "extra": "11676180 times\n4 procs"
          },
          {
            "name": "BenchmarkGet/nested/keys=1000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 48,
            "unit": "B/op",
            "extra": "11676180 times\n4 procs"
          },
          {
            "name": "BenchmarkGet/nested/keys=1000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "11676180 times\n4 procs"
          },
          {
            "name": "BenchmarkGet/nested/keys=10000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 94.34,
            "unit": "ns/op\t      48 B/op\t       1 allocs/op",
            "extra": "12745552 times\n4 procs"
          },
          {
            "name": "BenchmarkGet/nested/keys=10000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 94.34,
            "unit": "ns/op",
            "extra": "12745552 times\n4 procs"
          },
          {
            "name": "BenchmarkGet/nested/keys=10000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 48,
            "unit": "B/op",
            "extra": "12745552 times\n4 procs"
          },
          {
            "name": "BenchmarkGet/nested/keys=10000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "12745552 times\n4 procs"
          },
          {
            "name": "BenchmarkAs/int64 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 144.1,
            "unit": "ns/op\t      64 B/op\t       3 allocs/op",
            "extra": "8597403 times\n4 procs"
          },
          {
            "name": "BenchmarkAs/int64 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 144.1,
            "unit": "ns/op",
            "extra": "8597403 times\n4 procs"
          },
          {
            "name": "BenchmarkAs/int64 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 64,
            "unit": "B/op",
            "extra": "8597403 times\n4 procs"
          },
          {
            "name": "BenchmarkAs/int64 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "8597403 times\n4 procs"
          },
          {
            "name": "BenchmarkAs/string (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 120.1,
            "unit": "ns/op\t      80 B/op\t       2 allocs/op",
            "extra": "9901095 times\n4 procs"
          },
          {
            "name": "BenchmarkAs/string (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 120.1,
            "unit": "ns/op",
            "extra": "9901095 times\n4 procs"
          },
          {
            "name": "BenchmarkAs/string (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 80,
            "unit": "B/op",
            "extra": "9901095 times\n4 procs"
          },
          {
            "name": "BenchmarkAs/string (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "9901095 times\n4 procs"
          },
          {
            "name": "BenchmarkAs/struct (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 808.8,
            "unit": "ns/op\t     320 B/op\t      14 allocs/op",
            "extra": "1422778 times\n4 procs"
          },
          {
            "name": "BenchmarkAs/struct (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 808.8,
            "unit": "ns/op",
            "extra": "1422778 times\n4 procs"
          },
          {
            "name": "BenchmarkAs/struct (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 320,
            "unit": "B/op",
            "extra": "1422778 times\n4 procs"
          },
          {
            "name": "BenchmarkAs/struct (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 14,
            "unit": "allocs/op",
            "extra": "1422778 times\n4 procs"
          },
          {
            "name": "BenchmarkAs/slice_string/len=100 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 17388,
            "unit": "ns/op\t   20432 B/op\t     304 allocs/op",
            "extra": "66087 times\n4 procs"
          },
          {
            "name": "BenchmarkAs/slice_string/len=100 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 17388,
            "unit": "ns/op",
            "extra": "66087 times\n4 procs"
          },
          {
            "name": "BenchmarkAs/slice_string/len=100 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 20432,
            "unit": "B/op",
            "extra": "66087 times\n4 procs"
          },
          {
            "name": "BenchmarkAs/slice_string/len=100 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 304,
            "unit": "allocs/op",
            "extra": "66087 times\n4 procs"
          },
          {
            "name": "BenchmarkMerge/changes=10 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 36402,
            "unit": "ns/op\t   33808 B/op\t     442 allocs/op",
            "extra": "30902 times\n4 procs"
          },
          {
            "name": "BenchmarkMerge/changes=10 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 36402,
            "unit": "ns/op",
            "extra": "30902 times\n4 procs"
          },
          {
            "name": "BenchmarkMerge/changes=10 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 33808,
            "unit": "B/op",
            "extra": "30902 times\n4 procs"
          },
          {
            "name": "BenchmarkMerge/changes=10 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 442,
            "unit": "allocs/op",
            "extra": "30902 times\n4 procs"
          },
          {
            "name": "BenchmarkMerge/changes=100 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 292993,
            "unit": "ns/op\t  320624 B/op\t    4312 allocs/op",
            "extra": "4119 times\n4 procs"
          },
          {
            "name": "BenchmarkMerge/changes=100 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 292993,
            "unit": "ns/op",
            "extra": "4119 times\n4 procs"
          },
          {
            "name": "BenchmarkMerge/changes=100 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 320624,
            "unit": "B/op",
            "extra": "4119 times\n4 procs"
          },
          {
            "name": "BenchmarkMerge/changes=100 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 4312,
            "unit": "allocs/op",
            "extra": "4119 times\n4 procs"
          },
          {
            "name": "BenchmarkMerge/changes=1000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 3083020,
            "unit": "ns/op\t 3533105 B/op\t   43780 allocs/op",
            "extra": "388 times\n4 procs"
          },
          {
            "name": "BenchmarkMerge/changes=1000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 3083020,
            "unit": "ns/op",
            "extra": "388 times\n4 procs"
          },
          {
            "name": "BenchmarkMerge/changes=1000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 3533105,
            "unit": "B/op",
            "extra": "388 times\n4 procs"
          },
          {
            "name": "BenchmarkMerge/changes=1000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 43780,
            "unit": "allocs/op",
            "extra": "388 times\n4 procs"
          },
          {
            "name": "BenchmarkList/build/len=100 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 66687,
            "unit": "ns/op\t   78912 B/op\t     488 allocs/op",
            "extra": "17748 times\n4 procs"
          },
          {
            "name": "BenchmarkList/build/len=100 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 66687,
            "unit": "ns/op",
            "extra": "17748 times\n4 procs"
          },
          {
            "name": "BenchmarkList/build/len=100 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 78912,
            "unit": "B/op",
            "extra": "17748 times\n4 procs"
          },
          {
            "name": "BenchmarkList/build/len=100 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 488,
            "unit": "allocs/op",
            "extra": "17748 times\n4 procs"
          },
          {
            "name": "BenchmarkList/build/len=1000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 744817,
            "unit": "ns/op\t  719307 B/op\t    4007 allocs/op",
            "extra": "1689 times\n4 procs"
          },
          {
            "name": "BenchmarkList/build/len=1000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 744817,
            "unit": "ns/op",
            "extra": "1689 times\n4 procs"
          },
          {
            "name": "BenchmarkList/build/len=1000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 719307,
            "unit": "B/op",
            "extra": "1689 times\n4 procs"
          },
          {
            "name": "BenchmarkList/build/len=1000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 4007,
            "unit": "allocs/op",
            "extra": "1689 times\n4 procs"
          },
          {
            "name": "BenchmarkList/build/len=10000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 7606669,
            "unit": "ns/op\t10017223 B/op\t   40146 allocs/op",
            "extra": "158 times\n4 procs"
          },
          {
            "name": "BenchmarkList/build/len=10000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 7606669,
            "unit": "ns/op",
            "extra": "158 times\n4 procs"
          },
          {
            "name": "BenchmarkList/build/len=10000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 10017223,
            "unit": "B/op",
            "extra": "158 times\n4 procs"
          },
          {
            "name": "BenchmarkList/build/len=10000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 40146,
            "unit": "allocs/op",
            "extra": "158 times\n4 procs"
          },
          {
            "name": "BenchmarkList/save/len=100 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 19673,
            "unit": "ns/op\t  14.28 MB/s\t    6976 B/op\t     179 allocs/op",
            "extra": "56398 times\n4 procs"
          },
          {
            "name": "BenchmarkList/save/len=100 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 19673,
            "unit": "ns/op",
            "extra": "56398 times\n4 procs"
          },
          {
            "name": "BenchmarkList/save/len=100 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 14.28,
            "unit": "MB/s",
            "extra": "56398 times\n4 procs"
          },
          {
            "name": "BenchmarkList/save/len=100 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 6976,
            "unit": "B/op",
            "extra": "56398 times\n4 procs"
          },
          {
            "name": "BenchmarkList/save/len=100 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 179,
            "unit": "allocs/op",
            "extra": "56398 times\n4 procs"
          },
          {
            "name": "BenchmarkList/save/len=1000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 179826,
            "unit": "ns/op\t   9.30 MB/s\t   28909 B/op\t    1087 allocs/op",
            "extra": "6934 times\n4 procs"
          },
          {
            "name": "BenchmarkList/save/len=1000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 179826,
            "unit": "ns/op",
            "extra": "6934 times\n4 procs"
          },
          {
            "name": "BenchmarkList/save/len=1000 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 9.3,
            "unit": "MB/s",
            "extra": "6934 times\n4 procs"
          },
          {
            "name": "BenchmarkList/save/len=1000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 28909,
            "unit": "B/op",
            "extra": "6934 times\n4 procs"
          },
          {
            "name": "BenchmarkList/save/len=1000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 1087,
            "unit": "allocs/op",
            "extra": "6934 times\n4 procs"
          },
          {
            "name": "BenchmarkList/save/len=10000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 1295687,
            "unit": "ns/op\t  14.65 MB/s\t  303006 B/op\t   10096 allocs/op",
            "extra": "897 times\n4 procs"
          },
          {
            "name": "BenchmarkList/save/len=10000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 1295687,
            "unit": "ns/op",
            "extra": "897 times\n4 procs"
          },
          {
            "name": "BenchmarkList/save/len=10000 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 14.65,
            "unit": "MB/s",
            "extra": "897 times\n4 procs"
          },
          {
            "name": "BenchmarkList/save/len=10000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 303006,
            "unit": "B/op",
            "extra": "897 times\n4 procs"
          },
          {
            "name": "BenchmarkList/save/len=10000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 10096,
            "unit": "allocs/op",
            "extra": "897 times\n4 procs"
          },
          {
            "name": "BenchmarkList/load/len=100 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 27087,
            "unit": "ns/op\t  10.37 MB/s\t   12336 B/op\t     221 allocs/op",
            "extra": "43388 times\n4 procs"
          },
          {
            "name": "BenchmarkList/load/len=100 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 27087,
            "unit": "ns/op",
            "extra": "43388 times\n4 procs"
          },
          {
            "name": "BenchmarkList/load/len=100 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 10.37,
            "unit": "MB/s",
            "extra": "43388 times\n4 procs"
          },
          {
            "name": "BenchmarkList/load/len=100 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 12336,
            "unit": "B/op",
            "extra": "43388 times\n4 procs"
          },
          {
            "name": "BenchmarkList/load/len=100 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 221,
            "unit": "allocs/op",
            "extra": "43388 times\n4 procs"
          },
          {
            "name": "BenchmarkList/load/len=1000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 268454,
            "unit": "ns/op\t   6.23 MB/s\t  149232 B/op\t    2116 allocs/op",
            "extra": "4598 times\n4 procs"
          },
          {
            "name": "BenchmarkList/load/len=1000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 268454,
            "unit": "ns/op",
            "extra": "4598 times\n4 procs"
          },
          {
            "name": "BenchmarkList/load/len=1000 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 6.23,
            "unit": "MB/s",
            "extra": "4598 times\n4 procs"
          },
          {
            "name": "BenchmarkList/load/len=1000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 149232,
            "unit": "B/op",
            "extra": "4598 times\n4 procs"
          },
          {
            "name": "BenchmarkList/load/len=1000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 2116,
            "unit": "allocs/op",
            "extra": "4598 times\n4 procs"
          },
          {
            "name": "BenchmarkList/load/len=10000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 2498507,
            "unit": "ns/op\t   7.60 MB/s\t 1016936 B/op\t   22165 allocs/op",
            "extra": "460 times\n4 procs"
          },
          {
            "name": "BenchmarkList/load/len=10000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 2498507,
            "unit": "ns/op",
            "extra": "460 times\n4 procs"
          },
          {
            "name": "BenchmarkList/load/len=10000 (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 7.6,
            "unit": "MB/s",
            "extra": "460 times\n4 procs"
          },
          {
            "name": "BenchmarkList/load/len=10000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 1016936,
            "unit": "B/op",
            "extra": "460 times\n4 procs"
          },
          {
            "name": "BenchmarkList/load/len=10000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 22165,
            "unit": "allocs/op",
            "extra": "460 times\n4 procs"
          },
          {
            "name": "BenchmarkList/iterate/len=100 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 7722,
            "unit": "ns/op\t   12992 B/op\t     101 allocs/op",
            "extra": "142172 times\n4 procs"
          },
          {
            "name": "BenchmarkList/iterate/len=100 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 7722,
            "unit": "ns/op",
            "extra": "142172 times\n4 procs"
          },
          {
            "name": "BenchmarkList/iterate/len=100 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 12992,
            "unit": "B/op",
            "extra": "142172 times\n4 procs"
          },
          {
            "name": "BenchmarkList/iterate/len=100 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 101,
            "unit": "allocs/op",
            "extra": "142172 times\n4 procs"
          },
          {
            "name": "BenchmarkList/iterate/len=1000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 70069,
            "unit": "ns/op\t  121794 B/op\t    1001 allocs/op",
            "extra": "16934 times\n4 procs"
          },
          {
            "name": "BenchmarkList/iterate/len=1000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 70069,
            "unit": "ns/op",
            "extra": "16934 times\n4 procs"
          },
          {
            "name": "BenchmarkList/iterate/len=1000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 121794,
            "unit": "B/op",
            "extra": "16934 times\n4 procs"
          },
          {
            "name": "BenchmarkList/iterate/len=1000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 1001,
            "unit": "allocs/op",
            "extra": "16934 times\n4 procs"
          },
          {
            "name": "BenchmarkList/iterate/len=10000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 727516,
            "unit": "ns/op\t 1206617 B/op\t   10019 allocs/op",
            "extra": "1638 times\n4 procs"
          },
          {
            "name": "BenchmarkList/iterate/len=10000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 727516,
            "unit": "ns/op",
            "extra": "1638 times\n4 procs"
          },
          {
            "name": "BenchmarkList/iterate/len=10000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 1206617,
            "unit": "B/op",
            "extra": "1638 times\n4 procs"
          },
          {
            "name": "BenchmarkList/iterate/len=10000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 10019,
            "unit": "allocs/op",
            "extra": "1638 times\n4 procs"
          },
          {
            "name": "BenchmarkList/as_slice/len=1000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 184154,
            "unit": "ns/op\t  170147 B/op\t    4004 allocs/op",
            "extra": "6325 times\n4 procs"
          },
          {
            "name": "BenchmarkList/as_slice/len=1000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 184154,
            "unit": "ns/op",
            "extra": "6325 times\n4 procs"
          },
          {
            "name": "BenchmarkList/as_slice/len=1000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 170147,
            "unit": "B/op",
            "extra": "6325 times\n4 procs"
          },
          {
            "name": "BenchmarkList/as_slice/len=1000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 4004,
            "unit": "allocs/op",
            "extra": "6325 times\n4 procs"
          },
          {
            "name": "BenchmarkTextRead/poorly_simulated_typing/n=1000 (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 4351,
            "unit": "ns/op\t    1064 B/op\t       3 allocs/op",
            "extra": "279714 times\n4 procs"
          },
          {
            "name": "BenchmarkTextRead/poorly_simulated_typing/n=1000 (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 4351,
            "unit": "ns/op",
            "extra": "279714 times\n4 procs"
          },
          {
            "name": "BenchmarkTextRead/poorly_simulated_typing/n=1000 (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 1064,
            "unit": "B/op",
            "extra": "279714 times\n4 procs"
          },
          {
            "name": "BenchmarkTextRead/poorly_simulated_typing/n=1000 (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "279714 times\n4 procs"
          },
          {
            "name": "BenchmarkTextRead/edit_trace (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 512922,
            "unit": "ns/op\t  106520 B/op\t       2 allocs/op",
            "extra": "2320 times\n4 procs"
          },
          {
            "name": "BenchmarkTextRead/edit_trace (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 512922,
            "unit": "ns/op",
            "extra": "2320 times\n4 procs"
          },
          {
            "name": "BenchmarkTextRead/edit_trace (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 106520,
            "unit": "B/op",
            "extra": "2320 times\n4 procs"
          },
          {
            "name": "BenchmarkTextRead/edit_trace (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "2320 times\n4 procs"
          },
          {
            "name": "BenchmarkSaveIncremental/incremental (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 2114,
            "unit": "ns/op\t 457.99 MB/s\t    2512 B/op\t       6 allocs/op",
            "extra": "572478 times\n4 procs"
          },
          {
            "name": "BenchmarkSaveIncremental/incremental (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 2114,
            "unit": "ns/op",
            "extra": "572478 times\n4 procs"
          },
          {
            "name": "BenchmarkSaveIncremental/incremental (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 457.99,
            "unit": "MB/s",
            "extra": "572478 times\n4 procs"
          },
          {
            "name": "BenchmarkSaveIncremental/incremental (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 2512,
            "unit": "B/op",
            "extra": "572478 times\n4 procs"
          },
          {
            "name": "BenchmarkSaveIncremental/incremental (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 6,
            "unit": "allocs/op",
            "extra": "572478 times\n4 procs"
          },
          {
            "name": "BenchmarkSaveIncremental/full (github.com/MichaelMure/gotomerge/docproxy)",
            "value": 12009,
            "unit": "ns/op\t  13.91 MB/s\t    6480 B/op\t      94 allocs/op",
            "extra": "97510 times\n4 procs"
          },
          {
            "name": "BenchmarkSaveIncremental/full (github.com/MichaelMure/gotomerge/docproxy) - ns/op",
            "value": 12009,
            "unit": "ns/op",
            "extra": "97510 times\n4 procs"
          },
          {
            "name": "BenchmarkSaveIncremental/full (github.com/MichaelMure/gotomerge/docproxy) - MB/s",
            "value": 13.91,
            "unit": "MB/s",
            "extra": "97510 times\n4 procs"
          },
          {
            "name": "BenchmarkSaveIncremental/full (github.com/MichaelMure/gotomerge/docproxy) - B/op",
            "value": 6480,
            "unit": "B/op",
            "extra": "97510 times\n4 procs"
          },
          {
            "name": "BenchmarkSaveIncremental/full (github.com/MichaelMure/gotomerge/docproxy) - allocs/op",
            "value": 94,
            "unit": "allocs/op",
            "extra": "97510 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteChange/ops=100 (github.com/MichaelMure/gotomerge/format)",
            "value": 61913,
            "unit": "ns/op\t   5.94 MB/s\t   15442 B/op\t     181 allocs/op",
            "extra": "19581 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteChange/ops=100 (github.com/MichaelMure/gotomerge/format) - ns/op",
            "value": 61913,
            "unit": "ns/op",
            "extra": "19581 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteChange/ops=100 (github.com/MichaelMure/gotomerge/format) - MB/s",
            "value": 5.94,
            "unit": "MB/s",
            "extra": "19581 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteChange/ops=100 (github.com/MichaelMure/gotomerge/format) - B/op",
            "value": 15442,
            "unit": "B/op",
            "extra": "19581 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteChange/ops=100 (github.com/MichaelMure/gotomerge/format) - allocs/op",
            "value": 181,
            "unit": "allocs/op",
            "extra": "19581 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteChange/ops=1000 (github.com/MichaelMure/gotomerge/format)",
            "value": 403343,
            "unit": "ns/op\t   9.59 MB/s\t  111117 B/op\t    1094 allocs/op",
            "extra": "2636 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteChange/ops=1000 (github.com/MichaelMure/gotomerge/format) - ns/op",
            "value": 403343,
            "unit": "ns/op",
            "extra": "2636 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteChange/ops=1000 (github.com/MichaelMure/gotomerge/format) - MB/s",
            "value": 9.59,
            "unit": "MB/s",
            "extra": "2636 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteChange/ops=1000 (github.com/MichaelMure/gotomerge/format) - B/op",
            "value": 111117,
            "unit": "B/op",
            "extra": "2636 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteChange/ops=1000 (github.com/MichaelMure/gotomerge/format) - allocs/op",
            "value": 1094,
            "unit": "allocs/op",
            "extra": "2636 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteChange/ops=10000 (github.com/MichaelMure/gotomerge/format)",
            "value": 4204713,
            "unit": "ns/op\t   9.84 MB/s\t 1914402 B/op\t   10226 allocs/op",
            "extra": "284 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteChange/ops=10000 (github.com/MichaelMure/gotomerge/format) - ns/op",
            "value": 4204713,
            "unit": "ns/op",
            "extra": "284 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteChange/ops=10000 (github.com/MichaelMure/gotomerge/format) - MB/s",
            "value": 9.84,
            "unit": "MB/s",
            "extra": "284 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteChange/ops=10000 (github.com/MichaelMure/gotomerge/format) - B/op",
            "value": 1914402,
            "unit": "B/op",
            "extra": "284 times\n4 procs"
          },
          {
            "name": "BenchmarkWriteChange/ops=10000 (github.com/MichaelMure/gotomerge/format) - allocs/op",
            "value": 10226,
            "unit": "allocs/op",
            "extra": "284 times\n4 procs"
          },
          {
            "name": "BenchmarkReadExamplar (github.com/MichaelMure/gotomerge/format)",
            "value": 11230,
            "unit": "ns/op\t  36.15 MB/s\t         1.000 chunks\t    2888 B/op\t      46 allocs/op",
            "extra": "105147 times\n4 procs"
          },
          {
            "name": "BenchmarkReadExamplar (github.com/MichaelMure/gotomerge/format) - ns/op",
            "value": 11230,
            "unit": "ns/op",
            "extra": "105147 times\n4 procs"
          },
          {
            "name": "BenchmarkReadExamplar (github.com/MichaelMure/gotomerge/format) - MB/s",
            "value": 36.15,
            "unit": "MB/s",
            "extra": "105147 times\n4 procs"
          },
          {
            "name": "BenchmarkReadExamplar (github.com/MichaelMure/gotomerge/format) - chunks",
            "value": 1,
            "unit": "chunks",
            "extra": "105147 times\n4 procs"
          },
          {
            "name": "BenchmarkReadExamplar (github.com/MichaelMure/gotomerge/format) - B/op",
            "value": 2888,
            "unit": "B/op",
            "extra": "105147 times\n4 procs"
          },
          {
            "name": "BenchmarkReadExamplar (github.com/MichaelMure/gotomerge/format) - allocs/op",
            "value": 46,
            "unit": "allocs/op",
            "extra": "105147 times\n4 procs"
          },
          {
            "name": "BenchmarkReadLarge (github.com/MichaelMure/gotomerge/format)",
            "value": 373150899,
            "unit": "ns/op\t  78.39 MB/s\t    259779 chunks\t462359632 B/op\t 5792673 allocs/op",
            "extra": "3 times\n4 procs"
          },
          {
            "name": "BenchmarkReadLarge (github.com/MichaelMure/gotomerge/format) - ns/op",
            "value": 373150899,
            "unit": "ns/op",
            "extra": "3 times\n4 procs"
          },
          {
            "name": "BenchmarkReadLarge (github.com/MichaelMure/gotomerge/format) - MB/s",
            "value": 78.39,
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
            "value": 462359632,
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
            "value": 26841,
            "unit": "ns/op\t   7.15 MB/s\t         2.000 chunks\t   85768 B/op\t      66 allocs/op",
            "extra": "44797 times\n4 procs"
          },
          {
            "name": "BenchmarkCompressed (github.com/MichaelMure/gotomerge/format) - ns/op",
            "value": 26841,
            "unit": "ns/op",
            "extra": "44797 times\n4 procs"
          },
          {
            "name": "BenchmarkCompressed (github.com/MichaelMure/gotomerge/format) - MB/s",
            "value": 7.15,
            "unit": "MB/s",
            "extra": "44797 times\n4 procs"
          },
          {
            "name": "BenchmarkCompressed (github.com/MichaelMure/gotomerge/format) - chunks",
            "value": 2,
            "unit": "chunks",
            "extra": "44797 times\n4 procs"
          },
          {
            "name": "BenchmarkCompressed (github.com/MichaelMure/gotomerge/format) - B/op",
            "value": 85768,
            "unit": "B/op",
            "extra": "44797 times\n4 procs"
          },
          {
            "name": "BenchmarkCompressed (github.com/MichaelMure/gotomerge/format) - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "44797 times\n4 procs"
          },
          {
            "name": "BenchmarkTextEdits (github.com/MichaelMure/gotomerge/opset)",
            "value": 733098,
            "unit": "ns/op\t    104852 chars\t  194233 B/op\t     367 allocs/op",
            "extra": "1506 times\n4 procs"
          },
          {
            "name": "BenchmarkTextEdits (github.com/MichaelMure/gotomerge/opset) - ns/op",
            "value": 733098,
            "unit": "ns/op",
            "extra": "1506 times\n4 procs"
          },
          {
            "name": "BenchmarkTextEdits (github.com/MichaelMure/gotomerge/opset) - chars",
            "value": 104852,
            "unit": "chars",
            "extra": "1506 times\n4 procs"
          },
          {
            "name": "BenchmarkTextEdits (github.com/MichaelMure/gotomerge/opset) - B/op",
            "value": 194233,
            "unit": "B/op",
            "extra": "1506 times\n4 procs"
          },
          {
            "name": "BenchmarkTextEdits (github.com/MichaelMure/gotomerge/opset) - allocs/op",
            "value": 367,
            "unit": "allocs/op",
            "extra": "1506 times\n4 procs"
          }
        ]
      }
    ]
  }
}