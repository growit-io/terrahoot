package terraform

// Based on: https://github.com/mixpanel/terraform-deps
// FIXME: attribution for terraform-deps

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"log"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
)

type ModuleParser struct {
	Dir string
}

type ModuleInfo struct {
	Dir     string
	deps    []string
	TFFiles []string
}

func (p ModuleParser) ParseModule() (*ModuleInfo, error) {
	parser := hclparse.NewParser()

	schema := &hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{Type: "module", LabelNames: []string{"name"}},
		},
	}

	moduleSchema := &hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{Name: "source", Required: true},
		},
	}

	modInfo := &ModuleInfo{
		Dir:     p.Dir,
		TFFiles: []string{},
		deps:    []string{},
	}

	dirEntries, err := os.ReadDir(p.Dir)
	if err != nil {
		log.Fatal(err)
	}

	/*
	 * Parse all .tf files in the module directory to collect all "module"
	 * statements pointing to a local directory source, and parse those modules
	 * recursively.
	 */
	for _, entry := range dirEntries {
		if entry.IsDir() {
			continue
		}

		if !strings.HasSuffix(entry.Name(), ".tf") {
			continue
		}

		path := path.Join(p.Dir, entry.Name())
		modInfo.TFFiles = append(modInfo.TFFiles, path)

		f, diags := parser.ParseHCLFile(path)
		if diags.HasErrors() {
			return nil, diags
		}

		content, _, diags := f.Body.PartialContent(schema)
		if diags.HasErrors() {
			fmt.Println(diags)
		}

		for _, block := range content.Blocks {
			content, _, _ := block.Body.PartialContent(moduleSchema)
			source := content.Attributes["source"]
			if source == nil {
				continue
			}

			val, err := source.Expr.Value(nil)
			if err.HasErrors() {
				return nil, err
			}

			src := val.AsString()
			if !strings.HasPrefix(src, "./") && !strings.HasPrefix(src, "../") {
				// Ignore non-local module sources silently.
				continue
			}

			src = filepath.Clean(filepath.Join(p.Dir, src))
			modInfo.deps = append(modInfo.deps, src)
		}
	}

	return modInfo, nil
}

func (m ModuleInfo) Deps() ([]*ModuleInfo, error) {
	depInfos := []*ModuleInfo{}

	for _, d := range m.deps {
		p := ModuleParser{Dir: d}
		depInfo, err := p.ParseModule()
		if err != nil {
			return nil, err
		}

		depInfos = append(depInfos, depInfo)
	}

	return depInfos, nil
}
