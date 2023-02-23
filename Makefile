mkfile_path := $(abspath $(lastword $(MAKEFILE_LIST)))
cur_makefile_path := $(patsubst %/Makefile, %, $(mkfile_path))
#window "/" transe to "\"
cur_path := $(patsubst %/, %\,$(cur_makefile_path))
export GOPATH:= ${cur_path}

default: BUILDTAGS=release
default: noVNC-go

noVNC-go:
	cd $(cur_makefile_path)/src && go install -tags $(BUILDTAGS) .
