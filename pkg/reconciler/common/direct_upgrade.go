/*
Copyright 2022 The Knative Authors

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

package common

import (
	"fmt"
	"os"
	"strings"
)

const (
	// DirectUpgradeVersionsEnvKey is the key of the environment variable to specify the direct upgrade versions
	DirectUpgradeVersionsEnvKey = "DIRECT_UPGRADE_VERSIONS"
	// DefaultDirectUpgradeVersions is the default direct upgrade versions
	// Use `-` to separate the versions that can be upgraded directly
	// Use `|` to separate between pairs of versions
	DefaultDirectUpgradeVersions = "v0.24-v1.2|v1.2-v1.10"
)

var (
	directUpgradeVersionsString string = DefaultDirectUpgradeVersions
	directUpgradeVersions       [][2]string
)

func init() {
	envInit()
}

// envInit get the direct upgrade versions from environment variable and parse it
func envInit() {
	directUpgradeVersions = [][2]string{}

	// get the direct upgrade versions from environment variable
	if v := os.Getenv(DirectUpgradeVersionsEnvKey); v != "" {
		directUpgradeVersionsString = v
	}

	// parse the direct upgrade versions
	versions := strings.Split(directUpgradeVersionsString, "|")
	for _, v := range versions {
		pair := strings.Split(v, "-")
		if len(pair) != 2 {
			fmt.Printf("Wrong direct upgrade versions configuration %q\n", v)
			continue
		}
		directUpgradeVersions = append(directUpgradeVersions, [2]string{pair[0], pair[1]})
	}
	fmt.Printf("Direct upgrade versions: %+v\n", directUpgradeVersions)
}

// canUpgradeDirectly returns true if the version is in the direct upgrade list
func canUpgradeDirectly(from, to string) bool {
	for _, v := range directUpgradeVersions {
		// upgrade
		if v[0] == from && v[1] == to {
			return true
		}
		// downgrade
		if v[1] == from && v[0] == to {
			return true
		}
	}
	return false
}
