// Copyright © 2020 Banzai Cloud
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"path/filepath"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	"emperror.dev/errors"
	"github.com/MakeNowJust/heredoc"
	"github.com/adobe/operator-tools/pkg/docgen"
)

var logger = logf.Log

func main() {
	crds()
}

func crds() {
	lister := docgen.NewSourceLister(
		map[string]docgen.SourceDir{
			"secret":    {Path: "pkg/secret", DestPath: "docs/types"},
			"volume":    {Path: "pkg/volume", DestPath: "docs/types"},
			"base":      {Path: "pkg/types", DestPath: "docs/types"},
			"overrides": {Path: "pkg/typeoverride", DestPath: "docs/overrides"},
		},
		logger.WithName("lister"))

	lister.IncludeSources = []string{
		".*types$",
		".*override$",
	}
	lister.IgnoredSources = []string{
		".*",
	}

	lister.Index = docgen.NewDoc(docgen.DocItem{
		Name:     "Readme",
		DestPath: "docs/types",
	}, logger.WithName("typedoc"))

	lister.Header = heredoc.Doc(`
		# Available Types
		
		For more information please click on the name
		<center>

		| Name | Description |
		|---|---|`,
	)

	lister.Footer = heredoc.Doc(`
		</center>
	`)

	lister.DocGeneratedHook = func(document *docgen.Doc) error {
		relPath, err := filepath.Rel(lister.Index.Item.DestPath, document.Item.DestPath)
		if err != nil {
			return errors.WrapIff(err, "failed to determine relpath for %s", document.Item.DestPath)
		}
		lister.Index.Append(fmt.Sprintf("| **[%s](%s)** | %s |",
			document.DisplayName,
			filepath.Join(relPath, document.Item.Name+".md"),
			document.Desc))
		return nil
	}

	if err := lister.Generate(); err != nil {
		panic(err)
	}
}
