# Copyright (C) 2009 The Android Open Source Project
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

"""A module for reading and parsing event-log-tags files."""
from __future__ import print_function

import re
import sys
from typing import Optional, TextIO, Union


class Tag(object):
  __slots__ = ["tagnum", "tagname", "description", "filename", "linenum"]

  def __init__(self, tagnum, tagname, description, filename, linenum):
    self.tagnum = tagnum
    self.tagname = tagname
    self.description = description
    self.filename = filename
    self.linenum = linenum


class TagFile(object):
  """Read an input event-log-tags file."""
  def AddError(self, msg, linenum=None):
    if linenum is None:
      linenum = self.linenum
    self.errors.append((self.filename, linenum, msg))

  def AddWarning(self, msg, linenum=None):
    if linenum is None:
      linenum = self.linenum
    self.warnings.append((self.filename, linenum, msg))

  def __init__(self, filename, file_object=None):
    """'filename' is the name of the file (included in any error
    messages).  If 'file_object' is None, 'filename' will be opened
    for reading."""
    self.errors = []
    self.warnings = []
    self.tags = []
    self.options = {}

    self.filename = filename
    self.linenum = 0

    try:
      if file_object:
        self.read_tags_from_file(file_object)
      else:
        with open(filename, "r") as tag_file:
          self.read_tags_from_file(tag_file)
    except (IOError, OSError) as ex:
      self.AddError(str(ex))

  def read_tags_from_file(self, file_object):
    # type: (TextIO) -> None
    for self.linenum, line in enumerate(file_object):
      self.linenum += 1

      line = line.strip()
      if not line or line[0] == '#':
        continue
      parts = re.split(r"\s+", line, 2)

      if len(parts) < 2:
        self.AddError("failed to parse \"%s\"" % (line,))
        continue

      if parts[0] == "option":
        self.options[parts[1]] = parts[2:]
        continue

      if parts[0] == "?":
        tag = None
      else:
        try:
          tag = int(parts[0])
        except ValueError:
          self.AddError("\"%s\" isn't an integer tag or '?'" % (parts[0],))
          continue

      tagname = parts[1]
      if len(parts) == 3:
        description = parts[2]
      else:
        description = None

      if description:
        # EventLog.java checks that the description field is
        # surrounded by parens, so we should too (to avoid a runtime
        # crash from badly-formatted descriptions).
        if not re.match(r"\(.*\)\s*$", description):
          self.AddError("tag \"%s\" has unparseable description" % (tagname,))
          continue

      self.tags.append(Tag(tag, tagname, description,
                           self.filename, self.linenum))


def BooleanFromString(s):
  """Interpret 's' as a boolean and return its value.  Raise
  ValueError if it's not something we can interpret as true or
  false."""
  s = s.lower()
  if s in ("true", "t", "1", "on", "yes", "y"):
    return True
  if s in ("false", "f", "0", "off", "no", "n"):
    return False
  raise ValueError("'%s' not a valid boolean" % (s,))


def WriteOutput(output_file, data):
  # type: (Optional[TextIO], Union[str, TextIO]) -> None
  """Write 'data' to the given output filename (which may be None to
  indicate stdout).  Emit an error message and die on any failure.
  'data' may be a string or a StringIO object."""
  if not isinstance(data, str):
    data = data.getvalue()
  try:
    if output_file is None:
      out = sys.stdout
      output_file = "<stdout>"
    else:
      out = open(output_file, "w")
    out.write(data)
  except (IOError, OSError) as e:
    print("failed to write %s: %s" % (output_file, e), file=sys.stderr)
    sys.exit(1)
  finally:
    if out != sys.stdout:
      out.close()
