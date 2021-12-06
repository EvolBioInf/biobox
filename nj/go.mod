module github.com/evolbioinf/biobox/nj

go 1.17

require (
	github.com/evolbioinf/biobox/util v0.0.0-20211124162559-cfcac00de39f
	github.com/evolbioinf/clio v0.0.0-20210309091639-82cb91a31b0c
	github.com/evolbioinf/dist v0.0.0-20211201102954-64b418a79d17
	github.com/evolbioinf/nwk v0.0.0-20211206105227-7d0e890b4e0a
)

require github.com/evolbioinf/fasta v0.0.0-20210309091740-48b65faf5d3e // indirect

replace github.com/evolbioinf/biobox/newick => ../newick
