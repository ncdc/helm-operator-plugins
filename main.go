/*
Copyright 2020 The Operator-SDK Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"
	"log"
	"runtime"

	kustomizev2 "sigs.k8s.io/kubebuilder/v3/pkg/plugins/common/kustomize/v2"

	"github.com/spf13/cobra"
	"sigs.k8s.io/kubebuilder/v3/pkg/cli"
	config "sigs.k8s.io/kubebuilder/v3/pkg/config/v3"
	"sigs.k8s.io/kubebuilder/v3/pkg/model/stage"
	"sigs.k8s.io/kubebuilder/v3/pkg/plugin"

	golangv4 "sigs.k8s.io/kubebuilder/v3/pkg/plugins/golang/v4"

	"github.com/operator-framework/helm-operator-plugins/internal/cmd/hybrid-operator/run"
	"github.com/operator-framework/helm-operator-plugins/internal/version"
	pluginv1alpha "github.com/operator-framework/helm-operator-plugins/pkg/plugins/hybrid/v1alpha"
)

func main() {
	commands := []*cobra.Command{
		run.NewCmd(),
	}
	c, err := cli.New(
		cli.WithCommandName("helm-operator"),
		cli.WithVersion(getVersion()),
		cli.WithPlugins(
			getHybridPlugin(),
			golangv4.Plugin{},
		),
		cli.WithDefaultProjectVersion(config.Version),
		cli.WithDefaultPlugins(config.Version, getHybridPlugin()),
		cli.WithExtraCommands(commands...),
	)
	if err != nil {
		log.Fatal(err)
	}

	if err := c.Run(); err != nil {
		log.Fatal(err)
	}
}

func getVersion() string {
	return fmt.Sprintf("helm-operator version: %q, commit: %q, go version: %q, GOOS: %q, GOARCH: %q\n",
		version.GitVersion, version.GitCommit, runtime.Version(), runtime.GOOS, runtime.GOARCH)
}

func getHybridPlugin() plugin.Bundle {
	hybridBundle, _ := plugin.NewBundleWithOptions(plugin.WithName("hybrid"),
		plugin.WithVersion(plugin.Version{Number: 1, Stage: stage.Alpha}),
		plugin.WithPlugins(kustomizev2.Plugin{}, pluginv1alpha.Plugin{}))

	return hybridBundle
}
