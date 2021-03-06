// Copyright 2018 Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package envoy

import (
	"fmt"
	"reflect"
	"testing"

	"istio.io/istio/pkg/config/mesh"
)

func TestEnvoyArgs(t *testing.T) {
	proxyConfig := mesh.DefaultProxyConfig()
	proxyConfig.ServiceCluster = "my-cluster"
	proxyConfig.Concurrency = 8

	cfg := ProxyConfig{
		Config:            proxyConfig,
		Node:              "my-node",
		LogLevel:          "trace",
		ComponentLogLevel: "misc:error",
		NodeIPs:           []string{"10.75.2.9", "192.168.11.18"},
		PodName:           "",
		PodNamespace:      "",
		PodIP:             nil,
		SDSUDSPath:        "udspath",
		SDSTokenPath:      "tokenpath",
	}

	test := &envoy{
		ProxyConfig: cfg,
		extraArgs:   []string{"-l", "trace", "--component-log-level", "misc:error"},
	}

	testProxy := NewProxy(cfg)
	if !reflect.DeepEqual(testProxy, test) {
		t.Errorf("unexpected struct got\n%v\nwant\n%v", testProxy, test)
	}

	got := test.args("test.json", 5, "testdata/bootstrap.json")
	want := []string{
		"-c", "test.json",
		"--restart-epoch", "5",
		"--drain-time-s", "45",
		"--parent-shutdown-time-s", "60",
		"--service-cluster", "my-cluster",
		"--service-node", "my-node",
		"--max-obj-name-len", fmt.Sprint(proxyConfig.StatNameLength),
		"--local-address-ip-version", "v4",
		"--log-format", "[Envoy (Epoch 5)] [%Y-%m-%d %T.%e][%t][%l][%n] %v",
		"-l", "trace",
		"--component-log-level", "misc:error",
		"--config-yaml", `{"key": "value"}`,
		"--concurrency", "8",
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("envoyArgs() => got:\n%v,\nwant:\n%v", got, want)
	}
}

// TestEnvoyRun is no longer used - we are now using v2 bootstrap API.
