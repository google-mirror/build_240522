#!/usr/bin/python
"""Simple wrapper to run warn_common with Python standard Pool."""

import multiprocessing

# pylint:disable=relative-beyond-top-level
from .warn_common import common_main


# This parallel_process could be changed depending on platform
# and availability of multi-process library functions.
def parallel_process(num_cpu, classify_warnings, groups):
  pool = multiprocessing.Pool(num_cpu)
  return pool.map(classify_warnings, groups)


def main():
  common_main(parallel_process)


if __name__ == '__main__':
  main()
