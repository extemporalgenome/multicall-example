# Overview

Minimal hybrid single-call/multi-call project example in Go.

For background on what a multi-call binary is, see:
http://www.redbooks.ibm.com/abstracts/tips0092.html

The multi-call technique can be used to bundle the behavior of more than one
utility within a single binary.

BusyBox and [ToyBox](http://www.landley.net/toybox/about.html) use this
approach (the only portable approach) to bundle a minimal *nix userland within
a single, compact executable; these are critical for embedded and recovery
systems.  Git uses this to allow invocations like `git-commit` to behave the
same as `git commit`.  Useful examples are plentiful, even though this
technique is often unknown or overlooked.

Single-call binaries, in contrast, only fulfill the behavior of a single
utility. These are the norm. A `command subcommand` approach that is not aware
of `os.Args[0]` is still a single-call binary). Projects that can alternatively
compile both multicall binaries as well as the subcommands as self-contained
binaries often involve complicated makefiles and pre-processor magic to achieve
that result. This provides the flexibility to deploy one or a few of the
subcommands when desired (without the redundancy of the complete set of
sub-commands).

Go has historically made this either-or compilation approach straightforward
and trivial: no makefiles, no magic. This is accomplished by having one
top-level main package import each subcommand's main package. Each subcommand
is self-contained, and when compiled directly, will produce a minimal binary
with only the code needed for that subcommand alone. When the top-level binary
is compiled, it will pull in the required functionality of all the subcommands
transparently, due to the Go toolchain's sensible build semantics and
conventions.

An indirect fix for [https://github.com/golang/go/issues/4210], which prevents
importing any main package, has not removed the ability for Go to utilize this
either-or compilation flexibility, but does force uses of this technique to add
needless boilerplate code in the form of an extra directory level worth of
throwaway subpackages (one per subcommand, in addition to the boilerplate
package main per subcommand to allow single-call compilation) in order to
satisfy the new high-level toolchain restrictions.

The below will work in 1.4, but not likely work in 1.5 (as of commit
[679fd5b](https://github.com/golang/go/commit/679fd5b4479e0b9936344a33e07a0d1f904c362b)).

# Usage

	$ cd ./fakegit
	$ go build

	# Get a list of available subcommands
	$ ./fakegit
	log
	commit

	# Invoke the log subcommand
	$ ./fakegit log
	log invoked with args: []
	(outputting log...)

	# Create a symlink for the commit subcommand. this could also be a hardlink
	$ ln -s fakegit fakegit-commit 

	# Invocations through this symlink will have os.Args[0] == ".../fakegit-commit"
	$ ./fakegit-commit -m "message" file1 file2
	commit invoked with args: [-m message file1 file2]

	# Alternatively call the binary with a customized zeroth argument (bash/syscall).
	$ bash -c 'exec -a fakegit-commit ./fakegit'
	commit invoked with args: []

	# Verify that we're not cheating.
	$ bash -c 'exec -a fakegit-not-a-command ./fakegit'
	unrecognized command "not-a-command"

	# Compile the fakegit-commit subcommand as a top-level, self-contained
	# single-call binary. This binary only contains commit-related code.
	$ cd ../fakegit-commit
	$ go build

	# This fakegit-commit isn't a symlink.
	$ ./fakegit-commit
	commit invoked with args: []

	# Single-call binaries don't care about their invoked name
	$ bash -c 'exec -a fakegit-not-a-command ./fakegit-commit'
	commit invoked with args: []
