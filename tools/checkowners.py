#!/usr/bin/python

"""Parse and check syntax errors of a given OWNERS file."""

import argparse
import re
import sys
import urllib.request, urllib.parse, urllib.error
import urllib.request, urllib.error, urllib.parse

parser = argparse.ArgumentParser(description='Check OWNERS file syntax')
parser.add_argument('-v', '--verbose', dest='verbose',
                    action='store_true', default=False,
                    help='Verbose output to debug')
parser.add_argument('-c', '--check_address', dest='check_address',
                    action='store_true', default=False,
                    help='Check email addresses')
parser.add_argument(dest='owners', metavar='OWNERS', nargs='+',
                    help='Path to OWNERS file')
args = parser.parse_args()

gerrit_server = 'https://android-review.googlesource.com'
checked_addresses = {}


def echo(msg):
  if args.verbose:
    print(msg)


def find_address(address):
  if address not in checked_addresses:
<<<<<<< HEAD   (c2b35d Merge "Merge empty history for sparse-8348651-L2230000095368)
    request = (gerrit_server + '/accounts/?n=1&o=ALL_EMAILS&q=email:'
               + urllib.quote(address))
=======
    request = (gerrit_server + '/accounts/?n=1&q=email:'
               + urllib.parse.quote(address))
>>>>>>> BRANCH (697279 Merge "Version bump to TKB1.220411.001.A1 [core/build_id.mk])
    echo('Checking email address: ' + address)
<<<<<<< HEAD   (c2b35d Merge "Merge empty history for sparse-8348651-L2230000095368)
    result = urllib2.urlopen(request).read()
    checked_addresses[address] = (
        result.find('"email":') >= 0 and result.find('"_account_id":') >= 0)
=======
    result = urllib.request.urlopen(request).read()
    checked_addresses[address] = result.find('"_account_id":') >= 0
    if checked_addresses[address]:
      echo('Found email address: ' + address)
>>>>>>> BRANCH (697279 Merge "Version bump to TKB1.220411.001.A1 [core/build_id.mk])
  return checked_addresses[address]


<<<<<<< HEAD   (c2b35d Merge "Merge empty history for sparse-8348651-L2230000095368)
=======
def check_address(fname, num, address):
  if find_address(address):
    return 0
  print('%s:%d: ERROR: unknown email address: %s' % (fname, num, address))
  return 1


>>>>>>> BRANCH (697279 Merge "Version bump to TKB1.220411.001.A1 [core/build_id.mk])
def main():
  # One regular expression to check all valid lines.
  noparent = 'set +noparent'
  email = '([^@ ]+@[^ @]+|\\*)'
  directive = '(%s|%s)' % (email, noparent)
  glob = '[a-zA-Z0-9_\\.\\-\\*\\?]+'
  perfile = 'per-file +' + glob + ' *= *' + directive
  pats = '(|%s|%s|%s)$' % (noparent, email, perfile)
  patterns = re.compile(pats)

  # One pattern to capture email address.
  email_address = '.*(@| |=|^)([^@ =]+@[^ @]+)'
  address_pattern = re.compile(email_address)

  error = 0
  for fname in args.owners:
    echo('Checking file: ' + fname)
    num = 0
    for line in open(fname, 'r'):
      num += 1
      stripped_line = re.sub('#.*$', '', line).strip()
      if not patterns.match(stripped_line):
<<<<<<< HEAD   (c2b35d Merge "Merge empty history for sparse-8348651-L2230000095368)
        error = 1
        print('%s:%d: ERROR: unknown line [%s]'
              % (fname, num, line.strip()))
      elif args.check_address and address_pattern.match(stripped_line):
        address = address_pattern.match(stripped_line).group(2)
        if find_address(address):
          echo('Found email address: ' + address)
        else:
          error = 1
          print('%s:%d: ERROR: unknown email address: %s'
                % (fname, num, address))
=======
        error += 1
        print('%s:%d: ERROR: unknown line [%s]' % (fname, num, line.strip()))
      elif args.check_address:
        if perfile_pattern.match(stripped_line):
          for addr in perfile_pattern.match(stripped_line).group(1).split(','):
            a = addr.strip()
            if a and a != '*':
              error += check_address(fname, num, addr.strip())
        elif address_pattern.match(stripped_line):
          error += check_address(fname, num, stripped_line)
>>>>>>> BRANCH (697279 Merge "Version bump to TKB1.220411.001.A1 [core/build_id.mk])
  sys.exit(error)

if __name__ == '__main__':
  main()
