#!/usr/bin/python
"""Call -m warn.warn to process warning messages."""

import os
import subprocess
import sys


def main():
  os.environ['PYTHONPATH'] = os.path.dirname(os.path.abspath(__file__))
  subprocess.check_call(['/usr/bin/python', '-m', 'warn.warn'] + sys.argv[1:])


if __name__ == '__main__':
  main()
