// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package main

import (
	"github.com/cilium/cilium/operator/cmd"
	"github.com/cilium/cilium/operator/xds"
	"github.com/cilium/cilium/pkg/hive"
)

func main() {
	operatorHive := hive.New(cmd.Operator)
	go xds.RunServer()
	cmd.Execute(cmd.NewOperatorCmd(operatorHive))
	
}
