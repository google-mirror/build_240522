#!/usr/bin/python
"""Call -m warn.warn to process warning messages.

This script is used by Android continuous build bots for all branches.
Old frozen branches will continue to use the old warn.py, and active
branches will use this new version to call -m warn.warn.
"""

import os
import subprocess
import sys


def main():
  os.environ['PYTHONPATH'] = os.path.dirname(os.path.abspath(__file__))
  subprocess.check_call(['/usr/bin/python', '-m', 'warn.warn'] + sys.argv[1:])


if __name__ == '__main__':
  main()
