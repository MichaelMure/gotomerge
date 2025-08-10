module gotomerge

go 1.23.1

// https://github.com/jcalabro/leb128/pull/1
//replace github.com/jcalabro/leb128 => github.com/MichaelMure/leb128 v0.0.0-20250809082817-dd087cd8bcba

replace github.com/jcalabro/leb128 => ../leb128

require (
	github.com/jcalabro/leb128 v1.0.2
	github.com/stretchr/testify v1.10.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
