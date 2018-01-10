// Copyright 2017-2018 Dale Lakes <spitfire@spitfy.re>. All rights reserved.
// Use of this source code is governed by the MIT license located in the LICENSE file.

package goctftime

import "net/http"

// DefaultHandler handles any request not defined by an existing handler. This function simply drops spurious requests and
// should not be modified.
func DefaultHandler(_ http.ResponseWriter, _ *http.Request) {}
