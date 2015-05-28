#!/usr/bin/python

from __future__ import print_function
from os import makedirs, path
from re import match
from sys import argv, exit, stderr

def parse_fstab(fstab):
    paths = set()
    with open(fstab, 'r') as f:
        for line in f:
            m = match("^[^#](?:\S+\s+)(\/\S+)\s+(\S+)\s+(?:\S+\s+){2}", line)
            if m is None:
                continue
            if (m.group(2) != "swap" and m.group(2) != "emmc" and m.group(2) != "mtd"):
                paths.add(m.group(1))

    return paths


def create_dirs(parent, dirs):
    for filename in dirs:
        fullpath = "%s%s" % (parent, filename)
        if path.exists(fullpath):
            if path.isdir(fullpath) and not path.islink(fullpath):
                continue
            else:
                print("%s already exists and is not a directory" %fullpath, file=stderr)
                return -1

        makedirs(fullpath)

    return 0


def usage():
    print("Usage: %s <fstab> <out_dir>" % argv[0], file=stderr)
    return -1


def main(args=None):
    if args is None or len(args) != 2:
        return usage()

    dirs = parse_fstab(args[0])
    if not dirs:
        return usage()

    return create_dirs(args[1], dirs)


if __name__ == "__main__":
    exit(main(argv[1:]))
