#!/usr/bin/python

"""Parse and check syntax errors of a given OWNERS file."""

import argparse
import re
import urllib2

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
    print msg


def check_address(fname, num, address):
  if address not in checked_addresses:
    request = gerrit_server + '/accounts/?suggest&q=' + address
    echo('Checking email address: ' + address)
    result = urllib2.urlopen(request).read()
    expected = '"email": "' + address + '"'
    checked_addresses[address] = (result.find(expected) >= 0)
  if checked_addresses[address]:
    echo('Found email address: ' + address)
  else:
    print '%s:%d: ERROR: unknown email address: %s' % (fname, num, address)


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
  address = re.compile(email_address)

  for fname in args.owners:
    echo('Checking file: ' + fname)
    num = 0
    for line in open(fname, 'r'):
      num += 1
      stripped_line = re.sub('#.*$', '', line).strip()
      if not patterns.match(stripped_line):
        print('%s:%d: ERROR: unknown line [%s]' %
              (args.owners, num, line.strip()))
      elif args.check_address and address.match(stripped_line):
        check_address(fname, num, address.match(stripped_line).group(2))

if __name__ == '__main__':
  main()
