package terraform_test

import (
	"testing"

	internal "github.com/growit-io/terrahoot/internal/terraform"
	"github.com/stretchr/testify/assert"
)

func TestModuleParser(t *testing.T) {
	parser := internal.ModuleParser{
		Dir: "../examples/terraform/hello",
	}

	info, err := parser.ParseModule()
	if err != nil {
		t.Fatalf(`ParseModule() failed for %s`, parser.Dir)
	}

	assert.ElementsMatchf(t, info.TFFiles, []string{
		parser.Dir + "/main.tf",
		parser.Dir + "/outputs.tf",
	}, "unexpected list of source files for Terraform module in %s", parser.Dir)
}
