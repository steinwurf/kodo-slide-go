kodo-slide-go
============

Go bindings for `kodo-slide-c`_.


Getting Started
---------------

These instructions will get you ready to start using this project
as a dependency for your go projects.

Prerequisites
~~~~~~~~~~~~~

This project depends on `kodo-slide-c`_, a C library which is not available as
a package. For this reason you will need to build and install `kodo-slide-c`_
before a successful executing of ``go get github.com/steinwurf/kodo-slide-go``
can be performed.

First checkout this git project.

::

    git clone https://github.com/steinwurf/kodo-slide-go


Use Waf to configure and build. This will ensure the correct version of
`kodo-slide-c`_ is used.

::

    cd kodo-slide-go
    python waf configure
    python waf build

After a successful configuration and compilation the products of the build needs
to be made available. This is accomplished with the following Waf install
command. Make sure you have set your $GOPATH environment variable.

::

    python waf install --install_static_libs --install_path $GOPATH/src/github.com/steinwurf/kodo-slide-c

``$GOPATH/src/github.com/steinwurf/kodo-slide-c`` is the path were kodo-slide-go
expects the needed static library and header is located.

Installing
~~~~~~~~~~

After completing the steps specified in `Prerequisites`_, installing
kodo-slide-go is as simple using the following ``go get`` command:

::

    go get github.com/steinwurf/kodo-slide-go

And similarly it can be used as a dependency like so:

::

    import (
        ...
        "github.com/steinwurf/kodo-slide-go"
    )

When using kodo-slide-go as a dependency in your project, the directions in
`Prerequisites`_ needs to be followed before your project can be built.

Running the tests
-----------------

To check if your installation was success you can try to run the tests like so:

::

    Give an example

License
-------
You will need a valid license to build `kodo-slide-c`_.

To obtain a valid Kodo license **you must fill out the license request** form_.

Kodo is available under a research- and education-friendly license, see the
details in the LICENSE.rst file.

.. _form: http://steinwurf.com/license/
.. _kodo-slide-c: https://github.com/steinwurf/kodo-slide-c
