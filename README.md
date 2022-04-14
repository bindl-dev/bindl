<img style="width: 20em;" src="/docs/static/lazy-gopher.png">
<i>Image is "Lazy Gopher" from the collection <a href="https://github.com/ashleymcnamara/gophers/blob/master/LazyGopher.png">"Gophers" by Ashley Willis</a>.</i>

# Bindl

Bindl is a downloader for programs used in a project, often not necessary at runtime, but essential for development or infrastructure.

Bindl is an distro-agnostic, offering ease of consistency in managing binaries across operating systems and distributions.

## Why?

At the core of it, Bindl is standardizing and securing the work of `curl && chmod`. Through Bindl, projects can rest assured that dependencies and programs are always verified through checksum (and signature if provided).

The ergonomics of adopting Bindl is about making sure that for a given commit in a project, it will have consistent dependency version and installation mechanism regardless of which machine is running.

To learn more about why Bindl exists and how it works, [take a look at the guides](https://bindl.dev/notes).

## Usage / Installation

Available under [examples/](examples/) directory. In short:

```bash
# Whichever directory you'd like bindl to exist
export OUTDIR=/usr/local/bin

# While it's convenient, please inspect bootstrap.sh before running :)
curl --location https://bindl.dev/bootstrap.sh | bash
```

You may try to install with `go get`, though versioning information may be incomplete as they are stamped in build.

And of course, [assets in releases](https://github.com/bindl-dev/bindl/releases) are downloadable for manual binary installation.

## Contributing to Bindl

Our guide on contributing to Bindl is specified in [`CONTRIBUTING.md`](CONTRIBUTING.md)
