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

package program_test

import (
	"context"
	"os/exec"
	"testing"
	"time"

	"github.com/bindl-dev/bindl/internal/log"
	"github.com/bindl-dev/bindl/program"
)

var (
	// Taken from bindl v0.0.1 release
	bundleOne = &program.CosignBundle{
		Artifact: `11d819d3641304e856bbb3fd7dce4523598bb2d1ce33229dfd864a469b37df13  bindl_0.0.1_Darwin_x86_64.tar.gz
3c5bb66b8dca06645ff1c6f592a7011d81721f283137960593b801a55e127bdd  bindl_0.0.1_Linux_i386.tar.gz.sbom
91ac00e534c41d779f00f13a3e6e56115af243a250a94b19495144798cd81ba9  bindl_0.0.1_Darwin_arm64.tar.gz
aca5be6f43e4b0b0903d8de3abe27a7ff74dd182920bed396b78415268017515  bindl_0.0.1_Linux_i386.tar.gz
b78ac3815211c9c02a0a91bee11866ec272f8af9479d92d7d757a2d8d2df1054  bindl_0.0.1_Linux_x86_64.tar.gz.sbom
bc6b420b0eadaa5a444806380df60bec66a61112eebc3b7e112617bec99ae034  bindl_0.0.1_Darwin_x86_64.tar.gz.sbom
c06517714b19c28f7ae4979e9665b343ca1aeed4b7c2e1a4461d62eb88039cd6  bindl_0.0.1_Linux_arm64.tar.gz.sbom
d643de7afa6e42fa2f7bb46daa8aefe0674a21ec27efc3f98fecc73bd5a56553  bindl_0.0.1_Linux_arm64.tar.gz
e8af53b1d50b7a1260131956817c5fc85f4b19874f0349e8a9e1204f94ee6a28  bindl_0.0.1_Darwin_arm64.tar.gz.sbom
ff4baf5bc73c1d973cf52a6eb5e84cf5d4808fe16cefbbc6b04e39f2d2096d9a  bindl_0.0.1_Linux_x86_64.tar.gz
`,
		Certificate: `LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUNFakNDQVplZ0F3SUJBZ0lVQU9SOEYyMUsva2thZzlKRDJzZHhMYXdpWnRjd0NnWUlLb1pJemowRUF3TXcKS2pFVk1CTUdBMVVFQ2hNTWMybG5jM1J2Y21VdVpHVjJNUkV3RHdZRFZRUURFd2h6YVdkemRHOXlaVEFlRncweQpNakF6TWpVd01qTTJOVEZhRncweU1qQXpNalV3TWpRMk5UQmFNQUF3V1RBVEJnY3Foa2pPUFFJQkJnZ3Foa2pPClBRTUJCd05DQUFSQ09zRHNCMUlsZGNzVlhmbG5wdnduSjZodktWSENZcU1wZ2g2NGZJYkY2SzkwWCtUSmppcUcKWFVMY2FrZHNJTlBlbWV6b2JDSDZWazFDcU8vbXJNVzVvNEhFTUlIQk1BNEdBMVVkRHdFQi93UUVBd0lIZ0RBVApCZ05WSFNVRUREQUtCZ2dyQmdFRkJRY0RBekFNQmdOVkhSTUJBZjhFQWpBQU1CMEdBMVVkRGdRV0JCUTVCeWNvCmh4a25udEFsbEQyVDBMWXZxOG5mbURBZkJnTlZIU01FR0RBV2dCUll3QjVma1VXbFpxbDZ6SkNoa3lMUUtzWEYKK2pBZUJnTlZIUkVCQWY4RUZEQVNnUkIzYVd4emIyNUFhSFZ6YVc0dVpHVjJNQ3dHQ2lzR0FRUUJnNzh3QVFFRQpIbWgwZEhCek9pOHZaMmwwYUhWaUxtTnZiUzlzYjJkcGJpOXZZWFYwYURBS0JnZ3Foa2pPUFFRREF3TnBBREJtCkFqRUFpdlNvdW9COWZwV1JyMkVHeVN5QXZaV000c3l5RlVzV0ZOQTNxVU1ITVgyUHp4RWFQU0JFd3NoZGhLd2MKZzNCZkFqRUF6ZmlJMm1ITWNHM0VCMmlJZlZmUDYwdFJkOWVCTjhsd2tPWmxtaHJyYTZyMHZvRFBpNGplVUxTSwpvYUM5VHRXWAotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg==
`,
		Signature: `MEQCIGmtkMAYFAEAX7djgTEvRQD3QIU+msFmGzI8hp5XMiJxAiA/4g262hTZwOZD5ZPjtQb9a4VGEFTppPORZaA5uwinhQ==`,
	}

	// Taken from bindl v0.0.2 release
	bundleTwo = &program.CosignBundle{
		Artifact: `802af63886b5a0a9a0789aeaf27fd16f5e1c0ddd4a9739e5b6bb860737baf053  bindl_0.0.2_Darwin_x86_64.tar.gz
94e01c1549350dc690e9ffd5a374e73e27959d7a5d3648e732cb6041cf30e385  bindl_0.0.2_Darwin_x86_64.tar.gz.sbom
a30560ee4253936604861c023898093593b8f8d4912ec9c1f6c77f8d031d59dd  bindl_0.0.2_Linux_i386.tar.gz.sbom
baaf812d5a1bf426bdf872d14ebd4f6697570d28fb1347032497c54b19b0fc68  bindl_0.0.2_Linux_arm64.tar.gz.sbom
be9d8dc65f0aa781e2e3986f10d132cf87285e3d7a728428fd933f0d7388b903  bindl_0.0.2_Linux_x86_64.tar.gz
d14a7592aedf6115f4bb76de37c2b17e8730ff73291873750babf704f124e0bf  bindl_0.0.2_Linux_i386.tar.gz
d62fcbea495678380cb4d57c9279e8613426f3fd642e02964b8f9e7380da9ea8  bindl_0.0.2_Linux_arm64.tar.gz
d782be6f477be9c61f50d2721a741f1b95adfc8a2c99c0b85d2dc45f1aee9b49  bindl_0.0.2_Darwin_arm64.tar.gz
fc1256a887589aac88fad6abe3f5faaa0b6c434479724840598d3cebfc89effe  bindl_0.0.2_Linux_x86_64.tar.gz.sbom
ffcd8a83e51163b8e13631689ac151ac45e1aa6c954991cd066ebed0816f4f5b  bindl_0.0.2_Darwin_arm64.tar.gz.sbom
`,
		Certificate: `LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUNFakNDQVplZ0F3SUJBZ0lVQU5Ba0xsWnk1cUlCOUFJRVZkT0FpZmxLU1Bjd0NnWUlLb1pJemowRUF3TXcKS2pFVk1CTUdBMVVFQ2hNTWMybG5jM1J2Y21VdVpHVjJNUkV3RHdZRFZRUURFd2h6YVdkemRHOXlaVEFlRncweQpNakF6TWpneE56UTRNVEJhRncweU1qQXpNamd4TnpVNE1EbGFNQUF3V1RBVEJnY3Foa2pPUFFJQkJnZ3Foa2pPClBRTUJCd05DQUFSWlFtR2Z1T2dtUyt1eFNLbWp2a3J6b2Y5QUdnbHNrU3pJY1ZMdVhnbTNwVVo2U0ZreFpodWoKb01oMkJsQjdDQzY4L2ROcFBETE8zcVI4VHFlcENUemNvNEhFTUlIQk1BNEdBMVVkRHdFQi93UUVBd0lIZ0RBVApCZ05WSFNVRUREQUtCZ2dyQmdFRkJRY0RBekFNQmdOVkhSTUJBZjhFQWpBQU1CMEdBMVVkRGdRV0JCUktyTjdlCjl3VzZUZjVQdGloa0k5MTIrZkk4T0RBZkJnTlZIU01FR0RBV2dCUll3QjVma1VXbFpxbDZ6SkNoa3lMUUtzWEYKK2pBZUJnTlZIUkVCQWY4RUZEQVNnUkIzYVd4emIyNUFhSFZ6YVc0dVpHVjJNQ3dHQ2lzR0FRUUJnNzh3QVFFRQpIbWgwZEhCek9pOHZaMmwwYUhWaUxtTnZiUzlzYjJkcGJpOXZZWFYwYURBS0JnZ3Foa2pPUFFRREF3TnBBREJtCkFqRUEwc0N1dHYrZDVaVnZhSFZQalJseVk2dVd0RlA0U1g0ZWtjcEg3ZGNnemZxOUEyeURpalViSm1zVmRUKzcKNUhZMkFqRUF5ZlpSM3lwNXNkcEFta0FZQkVrMVlkaHRxTHoxdlpaQVZoNGUxT2xSZzRlMGhldCtVdnJDeFY0UgpXUGt3SExwRgotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg==
`,
		Signature: `MEUCIGPQhWBML3ckqaoCeDBOLY6w/+IqXTdq/64e9E0mZA+LAiEA9b9k7s5Ftnir7fl7X6/ZYubGaP5k4bf7LKeCBWkR7Jc=`,
	}
)

func TestCosignBundleVerify(t *testing.T) {
	_, err := exec.LookPath("cosign")
	if err != nil {
		t.Fatalf("cosign is required for this test, but not found in $PATH: %v", err)
	}
	testCases := []struct {
		bundle     *program.CosignBundle
		expectPass bool
	}{
		{
			bundle:     bundleOne,
			expectPass: true,
		},
		{
			bundle:     bundleTwo,
			expectPass: true,
		},
		{
			bundle: &program.CosignBundle{
				Artifact:    bundleOne.Artifact,
				Certificate: bundleTwo.Certificate,
				Signature:   bundleTwo.Signature,
			},
			expectPass: false,
		},
		{
			bundle: &program.CosignBundle{
				Artifact:    bundleOne.Artifact,
				Certificate: bundleTwo.Certificate,
				Signature:   bundleOne.Signature,
			},
			expectPass: false,
		},
	}

	log.SetLevel("debug")
	for i, tc := range testCases {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		err := tc.bundle.VerifySignature(ctx)
		if tc.expectPass {
			if err != nil {
				t.Fatalf("test case %d: expecting pass, but received error: %v", i, err)
			}
		} else {
			if err == nil {
				t.Fatalf("test case %d: expecting fail, but passed", i)
			}
		}
		cancel()
	}
}
