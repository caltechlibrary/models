Installation
============

__model__ is a Go package but also includes a demo application
called __modelgen__. You can install __modelgen__ via instructions
below. NOTE: This is an experimental proof of concept.

Quick install with curl or irm
------------------------------

There is an experimental installer.sh script that can be run with the
following command to install latest table release. This may work for
macOS, Linux and if you're using Windows with the Unix subsystem. This
would be run from your shell (e.g. Terminal on macOS).

~~~
curl https://caltechlibrary.github.io/models/installer.sh | sh
~~~

This will install modelgen in your `$HOME/bin` directory.

If you are running Windows 10 or 11 use the Powershell command
below.

~~~
irm https://caltechlibrary.github.io/models/installer.ps1 | iex
~~~

If your want to install a specific verions set the `PKG_VERSION` environment
variable then download. E.g. version 2.1.5 in tihs example.

For Linux and macOS

~~~
export PKG_VERSION=0.0.1
curl https://caltechlibrary.github.io/models/installer.sh | sh
~~~

For Windows

~~~
$env:PKG_VERSION = '0.0.1'
irm https://caltechlibrary.github.io/models/installer.ps1 | iex
~~~

## Compiling from source

You need to have git, Pandoc, Go compiler and Make (GNU Make) available for 
this recipe to work.  Clone the repository and then compile in the typical
POSIX style. NOTE by default the binaries are installed in `$HOME/bin` and
that is assumed to be in your path.

```shell
    cd
    git clone https://github.com/caltechlibrary/models
    cd models
    make
    # Add any missing dependencies you might need in your Go environment
    make test
    make install
```

On Windows you would perform the following in Powershell.

```shell
    cd
    git clone https://github.com/caltechlibrary/models
    cd models
    .\make.bat
    # Follow the prompts and instruction in the bat file.
```


### Requirements

- Go version 1.23.1 or better
- Pandoc version 3.1 or better
- GNU Make
- Common POSIX/Unix utilities, e.g. cat, sed, grep
- jq >= 1.7.1

### Windows compilation

The tool chain to compile on Windows make several assumptions.

1. You're using Anaconda shell and have the C tool chain installed for
   cgo to work
2. GNU Make, cat, grep and sed
3. You have the latest Go installed

Since I don't assume a POSIX shell environment on windows I have made
batch files to perform some of what Make under Linux and macOS would do.

- make.bat builds our application and depends on go and jq commands

Compilation assumes [Go](https://github.com/golang/go) v1.23.1 or better.

