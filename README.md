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

## Developing Bindl

### Dependencies

To develop on Bindl, you may need to setup several programs:
- GNU Make
- [direnv](https://direnv.net/)

If your project relies on Bindl, the programs above are optional, but recommended as there are workflows which integrates well with them.

### Git hooks

If your project utilizes Bindl and Makefile, it is recommended to use [Git hooks](https://git-scm.com/book/en/v2/Customizing-Git-Git-Hooks), specifically `post-checkout`:


```sh
touch .bindl-lock.yaml
```

By running the line above every successful `git checkout`, your Make rules will forcibly rebuild binary dependencies the next time its invoked. Bindl will then validate if existing program is consistent with lockfile (i.e. `.bindl-lock.yaml`) and let user proceed if it is. Otherwise, Bindl will attempt to lookup locally and fallback to downloading the program.

This is particularly important when working with branches which has different versions of programs declared in `.bindl-lock.yaml`, because Bindl can enforce consistency with the current branch.
