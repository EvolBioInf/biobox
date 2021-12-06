module github.com/evolbioinf/biobox/genTree

go 1.17

require (
	github.com/evolbioinf/biobox/util v0.0.0-20211119161501-576d171cf5f4
	github.com/evolbioinf/clio v0.0.0-20210309091639-82cb91a31b0c
	github.com/evolbioinf/nwk v0.0.0-20211206105227-7d0e890b4e0a
)

require github.com/evolbioinf/fasta v0.0.0-20210309091740-48b65faf5d3e // indirect

replace github.com/evolbioinf/biobox/newick => ../newick
