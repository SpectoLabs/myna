# myna

*"Process virtualization for the people by the people with people." - People*

Myna is a testing tool that captures and replays the output of command line programs.
It's heavily inspired by [SpectoLab's Hoverfly](https://github.com/SpectoLabs/hoverfly), 
which does capture and playback for http(s) web services.

Myna is written in Go and makes use of Boltdb for storage.  It was
initially built to facilitate testing of [KubeFuse](https://github.com/bspaans/kubefuse). 

It does not work with interactive programs at this point in time, but it's
pretty baller with the rest :sparkles:

## Installation

Myna uses the Glide build tool so you can just do:

```
glide install 
go build
```

And you're laughing.


## Usage

Capture the output of `ls -al /`

```sh
myna --capture ls -al /
```

Or by setting the CAPTURE environment variable:

```sh
CAPTURE=1 myna ls -al /
```

Play back the output of `ls -al /`

```sh
myna ls -al /
```

By default myna will save all its state in a file called `processes.db`, 
but we can override that by setting the `DATABASE_LOCATION` environment 
variable:

```sh
DATABASE_LOCATION=/tmp/my.db myna --capture ls -al /
DATABASE_LOCATION=/tmp/my.db myna ls -al /
```

To use the simulated process from our tests we can update our code to use myna
in place of the usual binary, but by writing a wrapper `ls` file and putting
this on the PATH we can seamlessly use myna from our tests without having to
change a single thing:

```sh
cat > ls <<EOF 
#!/bin/bash 
myna "\$0" "\$@"
EOF

chmod +x ls  
export PATH="`pwd`:$PATH"
```

## Usage from the Shell

We can also specify a cheeky little alias to work with myna from the
command line, which is less useful for automated testing, but still manages to
provide fun for the whole family :family:

```sh
$ alias ls="myna ls"
$ ls /
./myna does not know this process. Run the command in capture mode first.
$ export CAPTURE=1
$ ls /
drwxr-xr-x  24 root root      4096 May 14 12:29 .
drwxr-xr-x  24 root root      4096 May 14 12:29 ..
drwxr-xr-x   2 root root     12288 Mar 20 12:37 bin
drwxr-xr-x   4 root root      4096 May 14 12:30 boot
drwxrwxr-x   2 root root      4096 Aug 10  2015 cdrom
drwxr-xr-x  19 root root      4880 May 14 10:17 dev
drwxr-xr-x 159 root root     12288 May 14 12:29 etc
drwxr-xr-x   4 root root      4096 Aug 10  2015 home
lrwxrwxrwx   1 root root        32 May 14 12:29 initrd.img -> boot/initrd.img-4.2.0-36-generic
lrwxrwxrwx   1 root root        32 Apr  6 18:58 initrd.img.old -> boot/initrd.img-4.2.0-35-generic
drwxr-xr-x  29 root root      4096 Feb 17 21:13 lib
drwxr-xr-x   2 root root      4096 Feb 17 21:13 lib32
drwxr-xr-x   2 root root      4096 Feb 17 21:13 lib64
drwx------   2 root root     16384 Aug 10  2015 lost+found
drwxr-xr-x   3 root root      4096 Aug 10  2015 media
drwxr-xr-x   2 root root      4096 Apr 17  2015 mnt
drwxr-xr-x   5 root root      4096 Apr  5 00:07 opt
dr-xr-xr-x 297 root root         0 May 13 22:34 proc
drwx------   8 root root      4096 Nov 12  2015 root
drwxr-xr-x  30 root root       960 May 15 03:05 run
drwxr-xr-x   2 root root     12288 Mar 20 12:37 sbin
drwxr-xr-x   2 root root      4096 Apr 22  2015 srv
dr-xr-xr-x  13 root root         0 May 15 03:31 sys
drwxrwxrwt  12 root root     12288 May 15 04:49 tmp
drwxr-xr-x  12 root root      4096 Jan 30 15:09 usr
drwxr-xr-x  13 root root      4096 Apr 22  2015 var
lrwxrwxrwx   1 root root        29 May 14 12:29 vmlinuz -> boot/vmlinuz-4.2.0-36-generic
lrwxrwxrwx   1 root root        29 Apr  6 18:58 vmlinuz.old -> boot/vmlinuz-4.2.0-35-generic
$ export CAPTURE=0
$ ls /
drwxr-xr-x  24 root root      4096 May 14 12:29 .
drwxr-xr-x  24 root root      4096 May 14 12:29 ..
drwxr-xr-x   2 root root     12288 Mar 20 12:37 bin
drwxr-xr-x   4 root root      4096 May 14 12:30 boot
drwxrwxr-x   2 root root      4096 Aug 10  2015 cdrom
drwxr-xr-x  19 root root      4880 May 14 10:17 dev
drwxr-xr-x 159 root root     12288 May 14 12:29 etc
drwxr-xr-x   4 root root      4096 Aug 10  2015 home
lrwxrwxrwx   1 root root        32 May 14 12:29 initrd.img -> boot/initrd.img-4.2.0-36-generic
lrwxrwxrwx   1 root root        32 Apr  6 18:58 initrd.img.old -> boot/initrd.img-4.2.0-35-generic
drwxr-xr-x  29 root root      4096 Feb 17 21:13 lib
drwxr-xr-x   2 root root      4096 Feb 17 21:13 lib32
drwxr-xr-x   2 root root      4096 Feb 17 21:13 lib64
drwx------   2 root root     16384 Aug 10  2015 lost+found
drwxr-xr-x   3 root root      4096 Aug 10  2015 media
drwxr-xr-x   2 root root      4096 Apr 17  2015 mnt
drwxr-xr-x   5 root root      4096 Apr  5 00:07 opt
dr-xr-xr-x 297 root root         0 May 13 22:34 proc
drwx------   8 root root      4096 Nov 12  2015 root
drwxr-xr-x  30 root root       960 May 15 03:05 run
drwxr-xr-x   2 root root     12288 Mar 20 12:37 sbin
drwxr-xr-x   2 root root      4096 Apr 22  2015 srv
dr-xr-xr-x  13 root root         0 May 15 03:31 sys
drwxrwxrwt  12 root root     12288 May 15 04:49 tmp
drwxr-xr-x  12 root root      4096 Jan 30 15:09 usr
drwxr-xr-x  13 root root      4096 Apr 22  2015 var
lrwxrwxrwx   1 root root        29 May 14 12:29 vmlinuz -> boot/vmlinuz-4.2.0-36-generic
lrwxrwxrwx   1 root root        29 Apr  6 18:58 vmlinuz.old -> boot/vmlinuz-4.2.0-35-generic
```


## A Continuous Integration Workflow


For [KubeFuse](https://github.com/bspaans/kubefuse) I have to generate more
than one test case.  To keep this process repeatable I added a `capture` step
to my build system that I would run whenever the command under test (`kubectl`
in this case) was updated.

The build step looks something like this:

```sh
#!/bin/bash 

# Put myna into capture mode 
export CAPTURE=1

# Remove the old myna database if it exists
rm -f processes.db

# Run the commands I need for testing:
myna kubectl get namespaces
myna kubectl get pods --all-namespaces
myna kubectl get svc --all-namespaces
myna kubectl get rc --all-namespaces
....

# Export the results
myna --export | python -m json.tool > kubectl.json
```

Before I start running my tests I need to create the kubectl shim and import
the definition:

```sh
cat > bin/kubectl <<EOF 
#!/bin/bash 

unset CAPTURE
myna kubectl "\$@"
EOF

chmod +x bin/kubectl
rm processes.db
myna --import kubectl.json
```

And then I can start my tests with a modified PATH so that the tests pick 
up the shim first:

```sh
PATH="bin/:$PATH" nosetests
```

And that's it. Your uncle is Bob.


## Using Myna in your Python Unit Tests

I actually kind of lied in the previous section, because since the release 
of `python-myna` (see the `contrib/python-myna/` directory) I don't have to 
write any shims or modify my path, because I let this library take care of it 
for me.

For this example I'm using `nose` and `python-myna`, which you can both install 
using pip.

The KubeFuse repository is roughly structured like this:

```
kubefuse/kubefuse.py
kubefuse/all_the_rest.py
tests/test_kubefuse.py
tests/test_all_the_rest.py
tests/__init__.py
```

The kubefuse directory contains all the source code and the tests directory the 
tests. I run the suite by starting `nosetests` in the root.

To make my tests work with Myna I use package level `setUp` and `tearDown` 
functions that take care of handling the shim creation and PATH jumbling 
for me. Luckily I can just import these functions straight from the `python-myna`
library, so my `__init__.py` looks like this:

```python
from myna import *
```

The only thing that is missing now is to create and then later import the 
simulations. This is similar to the steps describe above, and something 
for which I use Make, because I love Make and Make is awesome. 

```make
myna:
	./kubectl.myna.sh

test:
	DATABASE_LOCATION="`pwd`/kubectl.db" myna --import tests/kubectl.json
	DATABASE_LOCATION="`pwd`/kubectl.db" nosetests
```


## License

Apache License version 2.0 [See LICENSE for details](./blob/master/LICENSE).

(c) [SpectoLabs](https://specto.io) 2016.

