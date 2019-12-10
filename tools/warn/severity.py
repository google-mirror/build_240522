"""Severity levels and attributes."""


class Severity(object):
  """Severity levels and attributes."""
  # numbered by dump order
  FIXMENOW = 0
  HIGH = 1
  MEDIUM = 2
  LOW = 3
  ANALYZER = 4
  TIDY = 5
  HARMLESS = 6
  UNKNOWN = 7
  SKIP = 8
  range = range(SKIP + 1)
  attributes = [
      # pylint:disable=bad-whitespace
      ['fuchsia',   'FixNow',    'Critical warnings, fix me now'],
      ['red',       'High',      'High severity warnings'],
      ['orange',    'Medium',    'Medium severity warnings'],
      ['yellow',    'Low',       'Low severity warnings'],
      ['hotpink',   'Analyzer',  'Clang-Analyzer warnings'],
      ['peachpuff', 'Tidy',      'Clang-Tidy warnings'],
      ['limegreen', 'Harmless',  'Harmless warnings'],
      ['lightblue', 'Unknown',   'Unknown warnings'],
      ['grey',      'Unhandled', 'Unhandled warnings']
  ]
  colors = [a[0] for a in attributes]
  column_headers = [a[1] for a in attributes]
  headers = [a[2] for a in attributes]

