#! /usr/bin/env python
# encoding: utf-8
import os
import string

from waflib import Utils
from waflib.TaskGen import feature, after_method

APPNAME = 'kodo-slide-go'
VERSION = '0.0.0'


def build(bld):
    if not bld.has_tool_option('install_static_libs'):
        return

    # We assume that install_path is available if install_static_libs is set.
    install_path = bld.get_tool_option('install_path')
    install_path = os.path.abspath(os.path.expanduser(install_path))

    # Substitute environment variables
    install_path = string.Template(install_path).substitute(os.environ)
    kodo_slide_c = bld.dependency_path('kodo-slide-c')
    kodo_slide_c_h = bld.root.find_node(os.path.realpath(os.path.join(
        kodo_slide_c, 'src', 'kodo_slide_c', 'kodo_slide_c.h')))

    bld.install_files(install_path, [kodo_slide_c_h])
