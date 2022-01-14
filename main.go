package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/m-mine/proto-to-postman/postman"
	"golang.org/x/xerrors"
)

func main() {
	opt, paths, err := parseOption()
	if err != nil || len(paths) == 0 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	w := os.Stdout
	if err := run(paths, opt.importPaths, opt.configName, opt.baseURL, opt.headers, w); err != nil {
		fmt.Fprintf(os.Stderr, "%+v", err)
		os.Exit(1)
	}
}

func run(files []string, importPaths []string, configName, baseURL string, headers []*postman.HeaderParam, w io.Writer) error {
	p := protoparse.Parser{
		ImportPaths: importPaths,
	}

	fds, err := p.ParseFiles(files...)
	if err != nil {
		return xerrors.Errorf("Unable to parse pb file: %v \n", err)
	}

	apiParamBuilder := NewAPIParamsBuilder(baseURL, headers)
	apiParams, err := apiParamBuilder.Build(fds)
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}

	postman := postman.Build(configName, apiParams)
	body, err := json.Marshal(postman)
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}

	_, err = fmt.Fprintf(w, "%s\n", body)
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}

	return nil
}
