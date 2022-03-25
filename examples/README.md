# Usage

To use Bindl in your project, you can start by bootstrapping Bindl:

```bash
export OUTDIR=$PWD/bin # Desired installation directory
curl --location https://bindl.dev/bootstrap.sh | bash
```

Afterwards, create `bindl.yaml` file to declare the programs your project need, as well as the platforms (using `GOOS` and `GOARCH` conventions) that your project would like to support.

```yaml
platforms:
  linux:
    - amd64
  darwin:
    - arm64

# Optionally, define popular platform naming convention for overlay
_uname: &uname
  OS:
    linux: Linux
    darwin: Darwin
  Arch:
    amd64: x86_64

programs:
  - name: terraform
    version: 1.1.8
    provider: url
    paths:
      base: 'https://releases.hashicorp.com/terraform/{{ .Version }}/'
      target: '{{ .Name }}_{{ .Version }}_{{ .OS }}_{{ .Arch }}.zip' #=> terraform_1.1.8_linux_amd64.zip
      checksums:
        artifact: '{{ .Name }}_{{ .Version }}_SHA256SUMS'
  - name: ko
    version: 0.11.2
    provider: github  # shortcut for GitHub releases
    overlay: *uname
    paths:
      base: google/ko  #=> https://github.com/google/ko/releases/download/v0.11.2/
      target: '{{ .Name }}_{{ .Version }}_{{ .OS }}_{{ .Arch }}.tar.gz'
      checksums:
        artifact: checksums.txt
  - name: cosign
    version: 1.7.1
    provider: github
    paths:
      base: sigstore/cosign
      target: '{{ .Name }}-{{ .OS }}-{{ .Arch }}'
      checksums: # will validate signatures through 'cosign verify-blob' if provided
        artifact: '{{ .Name }}_checksums.txt'
        certificate: '{{ .Name }}_checksums.txt-keyless.pem'
        signature: '{{ .Name }}_checksums.txt-keyless.sig'
```

Once done, run `bindl sync` to generate `.bindl-lock.yaml` and try running `bindl get` to install! By default, installations go to `bin/`.
