// Copyright 2022 Bindl Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package program

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/bindl-dev/bindl/internal"
)

var (
	bootstrapCosignOnce sync.Once
	bootstrapCosignArgs = []string{"get", "--bootstrap", "cosign", "--silent"}
)

func cosignPath(ctx context.Context) (string, error) {
	p, err := exec.LookPath("cosign")
	if err == nil {
		return p, nil
	}

	var stdout, stderr bytes.Buffer
	var bootstrapErr error
	bootstrapCosignOnce.Do(func() {
		bindlExecPath, err := os.Executable()
		if err != nil {
			bootstrapErr = err
			return
		}
		cmd := exec.CommandContext(ctx, bindlExecPath, bootstrapCosignArgs...)
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		err = cmd.Run()
		internal.Log().Debug().Strs("cmd", append([]string{bindlExecPath}, bootstrapCosignArgs...)).Err(err).Msg("bootstrapping cosign")
		if err != nil {
			bootstrapErr = fmt.Errorf("bootstrapping through bindl binary: %s", stderr.String())
			return
		}
	})
	if bootstrapErr != nil {
		return "", fmt.Errorf("bootstrapping cosign: %w", bootstrapErr)
	}
	cosignPath := strings.TrimSpace(stdout.String())
	internal.Log().Debug().Str("cosign", cosignPath).Msg("found cosign")
	return cosignPath, nil
}

type CosignBundle struct {
	Artifact    string `json:"artifact"`
	Certificate string `json:"certificate"`
	Signature   string `json:"signature"`
}

func (c *CosignBundle) Signed() bool {
	return c.Certificate != "" || c.Signature != ""
}

func (c *CosignBundle) VerifySignature(ctx context.Context) error {
	p, err := cosignPath(ctx)
	if err != nil {
		return err
	}

	dir, err := os.MkdirTemp(os.TempDir(), "bindl-cosign-bundle-*")
	if err != nil {
		return fmt.Errorf("creating cosign workspace: %w", err)
	}
	artifactPath := filepath.Join(dir, "artifact")
	if err := os.WriteFile(artifactPath, []byte(c.Artifact), 0666); err != nil {
		return fmt.Errorf("creating artifact file: %w", err)
	}
	certPath := filepath.Join(dir, "cert")
	if err := os.WriteFile(certPath, []byte(c.Certificate), 0666); err != nil {
		return fmt.Errorf("creating certificate file: %w", err)
	}
	sigPath := filepath.Join(dir, "sig")
	if err := os.WriteFile(sigPath, []byte(c.Signature), 0666); err != nil {
		return fmt.Errorf("creating signature file: %w", err)
	}

	cosignArgs := []string{
		"verify-blob",
		artifactPath,
		"--cert",
		certPath,
		"--signature",
		sigPath,
	}

	var stderr bytes.Buffer
	cmd := exec.CommandContext(ctx, p, cosignArgs...)
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err == nil {
		internal.Log().Debug().Str("cosign", stderr.String()).Send()
		os.RemoveAll(dir)
		return nil
	}

	cmdStr := append([]string{p}, cosignArgs...)

	internal.Log().Debug().Strs("cmd", cmdStr).Err(err).Str("stderr", stderr.String()).Send()

	return fmt.Errorf("failed to verify signature: %s", stderr.String())
}
