module github.com/evolbioinf/biobox/plotTree

go 1.17

replace github.com/evolbioinf/biobox/newick => ../newick

require (
	github.com/evolbioinf/biobox/newick v0.0.0-00010101000000-000000000000
	github.com/evolbioinf/biobox/util v0.0.0-20211111160454-71e4fae3cc66
	github.com/evolbioinf/clio v0.0.0-20210309091639-82cb91a31b0c
)

require github.com/evolbioinf/fasta v0.0.0-20210309091740-48b65faf5d3e // indirect
