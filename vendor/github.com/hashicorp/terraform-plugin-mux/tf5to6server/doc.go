// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// Package tf5to6server translates a provider that implements protocol version 5, into one that implements protocol version 6.
//
// Supported protocol version 5 provider servers include any which implement
// the tfprotov5.ProviderServer (https://pkg.go.dev/github.com/hashicorp/terraform-plugin-go/tfprotov5#ProviderServer)
// interface, such as:
//
//   - https://pkg.go.dev/github.com/hashicorp/terraform-plugin-go/tfprotov5/tf5server
//   - https://pkg.go.dev/github.com/hashicorp/terraform-plugin-mux/tf5muxserver
//   - https://pkg.go.dev/github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema
//
// Refer to the UpgradeServer() function for wrapping a server.
package tf5to6server
