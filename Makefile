GOCMD=go
BUILDCMD=$(GOCMD) build

all: build_mmap build_mmerge build_mgen

build_mmap:
		$(BUILDCMD) ./cmd/mmap/mmap.go

build_mmerge:
		$(BUILDCMD) ./cmd/mmerge/mmerge.go

build_mgen:
		$(BUILDCMD) ./cmd/mgen/mgen.go