#!/usr/bin/env python

"""
Check boot jars.

Usage: check_boot_jars.py <package_allow_list_file> <jar1> <jar2> ...
"""
import logging
import os.path
import re
import subprocess
import sys


# The compiled allow list RE.
allow_list_re = None


def LoadAllowList(filename):
  """ Load and compile allow list regular expressions from filename.
  """
  lines = []
  with open(filename, 'r') as f:
    for line in f:
      line = line.strip()
      if not line or line.startswith('#'):
        continue
      lines.append(line)
  combined_re = r'^(%s)$' % '|'.join(lines)
  global allow_list_re
  try:
    allow_list_re = re.compile(combined_re)
  except re.error:
    logging.exception(
        'Cannot compile package allow list regular expression: %r',
        combined_re)
    allow_list_re = None
    return False
  return True


<<<<<<< HEAD   (5c8d84 Merge "Merge empty history for sparse-6676661-L8360000065797)
def CheckJar(jar):
=======
def CheckJar(allow_list_path, jar):
>>>>>>> BRANCH (a10c18 Merge "Version bump to RT11.201014.001.A1 [core/build_id.mk])
  """Check a jar file.
  """
  # Get the list of files inside the jar file.
  p = subprocess.Popen(args='jar tf %s' % jar,
      stdout=subprocess.PIPE, shell=True)
  stdout, _ = p.communicate()
  if p.returncode != 0:
    return False
  items = stdout.split()
  classes = 0
  for f in items:
    if f.endswith('.class'):
      classes += 1
      package_name = os.path.dirname(f)
      package_name = package_name.replace('/', '.')
<<<<<<< HEAD   (5c8d84 Merge "Merge empty history for sparse-6676661-L8360000065797)
      # Skip class without a package name
      if package_name and not whitelist_re.match(package_name):
        print >> sys.stderr, ('Error: %s contains class file %s, which is not in the whitelist'
                              % (jar, f))
=======
      if not package_name or not allow_list_re.match(package_name):
        print >> sys.stderr, ('Error: %s contains class file %s, whose package name %s is empty or'
                              ' not in the allow list %s of packages allowed on the bootclasspath.'
                              % (jar, f, package_name, allow_list_path))
>>>>>>> BRANCH (a10c18 Merge "Version bump to RT11.201014.001.A1 [core/build_id.mk])
        return False
  if classes == 0:
    print >> sys.stderr, ('Error: %s does not contain any class files.' % jar)
    return False
  return True


def main(argv):
  if len(argv) < 2:
    print __doc__
    return 1
<<<<<<< HEAD   (5c8d84 Merge "Merge empty history for sparse-6676661-L8360000065797)
=======
  allow_list_path = argv[0]
>>>>>>> BRANCH (a10c18 Merge "Version bump to RT11.201014.001.A1 [core/build_id.mk])

<<<<<<< HEAD   (5c8d84 Merge "Merge empty history for sparse-6676661-L8360000065797)
  if not LoadWhitelist(argv[0]):
=======
  if not LoadAllowList(allow_list_path):
>>>>>>> BRANCH (a10c18 Merge "Version bump to RT11.201014.001.A1 [core/build_id.mk])
    return 1

  passed = True
  for jar in argv[1:]:
<<<<<<< HEAD   (5c8d84 Merge "Merge empty history for sparse-6676661-L8360000065797)
    if not CheckJar(jar):
=======
    if not CheckJar(allow_list_path, jar):
>>>>>>> BRANCH (a10c18 Merge "Version bump to RT11.201014.001.A1 [core/build_id.mk])
      passed = False
  if not passed:
    return 1

  return 0


if __name__ == '__main__':
  sys.exit(main(sys.argv[1:]))
