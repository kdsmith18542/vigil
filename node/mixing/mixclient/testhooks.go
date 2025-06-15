// Copyright (c) 2024 The Vigil developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package mixclient

type hook string

type hookFunc func(*Client, *pairedSessions, *sessionRun, *peer)

const (
	hookBeforeRun           hook = "before run"
	hookBeforePeerCTPublish hook = "before CT publish"
	hookBeforePeerSRPublish hook = "before SR publish"
	hookBeforePeerVGLPublish hook = "before DC publish"
)
