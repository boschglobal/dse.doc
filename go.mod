module github.com/boschglobal/dse.doc

go 1.19

replace (
	github.com/google/docsy => github.com/google/docsy v0.11.0
	github.com/google/docsy/dependencies => github.com/google/docsy/dependencies v0.7.2
	github.com/twbs/bootstrap => github.com/twbs/bootstrap v5.3.3+incompatible
	github.com/FortAwesome/Font-Awesome => github.com/FortAwesome/Font-Awesome v0.0.0-20240716171331-37eff7fa00de
)

require (
	github.com/boschglobal/dse.clib main // indirect
	github.com/boschglobal/dse.fmi main // indirect
	github.com/boschglobal/dse.modelc main // indirect
	github.com/boschglobal/dse.network main // indirect
	github.com/boschglobal/dse.schemas main // indirect
	github.com/boschglobal/dse.sdp main  // indirect
	github.com/boschglobal/dse.standards main // indirect
	github.com/FortAwesome/Font-Awesome v0.0.0-20240716171331-37eff7fa00de // indirect
	github.com/google/docsy v0.11.0 // indirect
	github.com/google/docsy/dependencies v0.7.2 // indirect
	github.com/twbs/bootstrap v5.3.3+incompatible // indirect
)
