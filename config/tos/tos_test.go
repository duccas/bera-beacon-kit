// SPDX-License-Identifier: MIT
//
// Copyright (c) 2023 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package tos_test

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	file "github.com/itsdevbear/bolaris/config"
	beaconflags "github.com/itsdevbear/bolaris/config/flags"
	"github.com/itsdevbear/bolaris/config/tos"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/itsdevbear/bolaris/examples/beacond/app"
	"github.com/itsdevbear/bolaris/examples/beacond/cmd/root"
)

const (
	acceptTosFilename   = "tosaccepted"
	declinedErrorString = "you have to accept Terms and Conditions in order to continue"
)

func expectTosAcceptSuccess(t *testing.T, homeDir string) {
	if ok := file.Exists(filepath.Join(homeDir, acceptTosFilename)); !ok {
		t.Errorf("Expected tosaccepted file to exist in %s", homeDir)
	}
}

func makeTempDir(t *testing.T) string {
	homeDir, err := os.MkdirTemp("", "beacond-test-*")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	t.Log("homeDir: ", homeDir)
	return homeDir
}

func TestAcceptTosFlag(t *testing.T) {
	stdout := os.Stdout
	defer func() { os.Stdout = stdout }()
	os.Stdout = os.NewFile(0, os.DevNull)
	homeDir := makeTempDir(t)
	defer os.RemoveAll(homeDir)

	rootCmd := root.NewRootCmd()
	rootCmd.SetOut(os.NewFile(0, os.DevNull))
	rootCmd.SetArgs([]string{
		"query",
		"--" + flags.FlagHome,
		homeDir,
		"--" + beaconflags.BeaconKitAcceptTos,
	})

	err := svrcmd.Execute(rootCmd, "", app.DefaultNodeHome)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	expectTosAcceptSuccess(t, homeDir)
}

func TestAcceptWithCLI(t *testing.T) {
	stdout := os.Stdout
	defer func() { os.Stdout = stdout }()
	// os.Stdout = os.NewFile(0, os.DevNull)
	homeDir := makeTempDir(t)
	t.Log("homeDir: ", homeDir)

	inputBuffer := bytes.NewReader([]byte("accept\n"))
	rootCmd := root.NewRootCmd()
	// rootCmd.SetOut(os.NewFile(0, os.DevNull))
	rootCmd.SetArgs([]string{
		"query",
		"--" + flags.FlagHome,
		homeDir,
	})
	rootCmd.SetIn(inputBuffer)

	if err := svrcmd.Execute(rootCmd, "", app.DefaultNodeHome); err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	expectTosAcceptSuccess(t, homeDir)
}

func TestDeclineWithCLI(t *testing.T) {
	stdout := os.Stdout
	defer func() { os.Stdout = stdout }()
	homeDir := makeTempDir(t)
	t.Log("homeDir: ", homeDir)

	inputBuffer := bytes.NewReader([]byte("decline\n"))
	rootCmd := root.NewRootCmd()
	rootCmd.SetArgs([]string{
		"query",
		"--" + flags.FlagHome,
		homeDir,
	})
	rootCmd.SetIn(inputBuffer)

	err := svrcmd.Execute(rootCmd, "", app.DefaultNodeHome)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
	if err.Error() != declinedErrorString {
		t.Errorf("Expected %v, got %v", declinedErrorString, err)
	}
}

type ErrReader struct{}

func (e *ErrReader) Read(p []byte) (int, error) {
	return 0, errors.New("forced error in scanner")
}
func TestDeclineWithNonInteractiveCLI(t *testing.T) {
	stdout := os.Stdout
	defer func() { os.Stdout = stdout }()
	homeDir := makeTempDir(t)
	t.Log("homeDir: ", homeDir)

	rootCmd := root.NewRootCmd()
	rootCmd.SetArgs([]string{
		"query",
		"--" + flags.FlagHome,
		homeDir,
	})
	rootCmd.SetIn(&ErrReader{})

	err := svrcmd.Execute(rootCmd, "", app.DefaultNodeHome)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}

	if !strings.Contains(err.Error(), tos.BuildErrorPromptText("")) {
		t.Errorf("Expected %v, got %v", tos.BuildErrorPromptText(""), err)
	}
}