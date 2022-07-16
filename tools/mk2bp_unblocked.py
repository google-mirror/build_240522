
import mk2bp_catalog

import csv
import os
import subprocess
import sys


def build(modules):
  "Returns true if the build passed."
  print("building: %s" % (" ".join(modules)))
  result = subprocess.run("build/soong/soong_ui.bash --make-mode %s" % (" ".join(modules),),
      shell=True, check=False)
  return result.returncode == 0

def bpfile_name(makefile):
  return makefile[0:-len("Android.mk")] + "Android.bp"

def check(makefile):
  "Return whether this makefile can be autoconverted"
  return subprocess.run("androidmk %s" % makefile,
                        shell=True, check=False, capture_output=True).returncode == 0

def convert(makefile):
  bpfile = bpfile_name(makefile)
  subprocess.run("androidmk %s > %s" % (makefile, bpfile), shell=True, capture_output=True)


def git_repository_for(filename):
  directory, leaf = os.path.split(filename)
  while directory != "/" and not os.path.exists(os.path.join(directory, ".git")):
    directory, left = os.path.split(directory)
  if directory == "/":
    return None
  return directory


def main(args):
  TARGET_PRODUCT = os.environ["TARGET_PRODUCT"]
  TARGET_DEVICE = os.environ["TARGET_DEVICE"]
  OUT_DIR = os.environ.get("OUT_DIR", "out")

  print("TARGET_DEVICE=%s" % TARGET_DEVICE)
  print("OUT_DIR=%s" % OUT_DIR)

  global HOST_OUT_ROOT
  HOST_OUT_ROOT = OUT_DIR + "/host"
  global PRODUCT_OUT
  PRODUCT_OUT = OUT_DIR + "/target/product/%s" % TARGET_DEVICE

  # Read target information
  # TODO: Pull from configurable location. This is also slightly different because it's
  # only a single build, where as the tree scanning we do below is all Android.mk files.
  with open("%s/obj/PACKAGING/soong_conversion_intermediates/soong_conv_data"
      % PRODUCT_OUT, "r", errors="ignore") as csvfile:
    soong = mk2bp_catalog.SoongData(csv.reader(csvfile))

  # Find the unblocked modules
  unblocked_modules = [m for m in soong.modules if len(soong.deps[m]) == 0]

  # Find the files that contain only unblocked modules
  files_to_check = set()
  blocked_makefiles = set()

  for module in unblocked_modules:
    makefiles = soong.makefiles[module]
    if len(makefiles) > 1:
      print("Will NOT convert (can't locate): %s -- %s" % (module, makefiles))
      continue
    makefile = makefiles[0]

    if not soong.contains_blocked_modules(makefile):
      print("Will convert:                    %s -- %s" % (module, makefiles))
      files_to_check.add(makefile)
    else:
      print("Will NOT convert (blocked):      %s -- %s" % (module, makefiles))
      blocked_makefiles.add(makefile)

  # Preflight the conversion. If it can't be done we won't touch the source tree.
  print("--- Work plan ------------------------------")
  files_to_convert = []
  for makefile in sorted(files_to_check):
    if check(makefile):
      files_to_convert.append(makefile)
      print("Will try to auto-convert: %s" % makefile)
    else:
      print("Can't auto-convert:       %s" % makefile)
  print("--------------------------------------------")

  # Group by git project
  git_projects = {}
  for makefile in files_to_convert:
    git_projects.setdefault(git_repository_for(makefile), []).append(makefile)

  print()
  changed_projects = list()

  # Do conversions one project at a time and make git commits
  for git_project, makefiles in git_projects.items():
    print("Git project: ", git_project)

    # The files in this project that were successfully converted
    project_converted_files = set()
    project_converted_modules = set()

    # Convert file by file, and revert it if the build fails
    for makefile in makefiles:
      # Skip weirdly named files
      if not makefile.endswith("/Android.mk"):
        continue

      # Skip files where there is already an Android.bp (those need to be appended
      # and we can't handle that here yet)
      if os.path.exists(bpfile_name(makefile)):
        continue

      # Try converting it
      print("Converting: %s -- %s" % (git_project, makefile))
      convert(makefile)
      subprocess.run("git rm %s" % makefile[len(git_project)+1:],
          cwd=git_project, shell=True, capture_output=True)

      if build(soong.reverse_makefiles[makefile]):
        print("PASSED %s" % makefile)
        # Save that it worked
        project_converted_files.add(makefile)
        for m in soong.reverse_makefiles[makefile]:
          project_converted_modules.add(m)
        # It worked, git add it
        subprocess.run("git add %s" % bpfile_name(makefile[len(git_project)+1:]),
            cwd=git_project, shell=True, capture_output=True)
      else:
        print("FAILED %s" % makefile)
        # It didn't work, revert it
        subprocess.run("git restore --staged %s" % makefile[len(git_project)+1:],
            cwd=git_project, shell=True, capture_output=True)
        subprocess.run("git restore %s" % makefile[len(git_project)+1:],
            cwd=git_project, shell=True, capture_output=True)
        try:
          print("removing %s" % bpfile_name(makefile))
          os.remove(bpfile_name(makefile))
        except OSError as error:
          print("Can't remove %s because of %s" % (bpfile_name(makefile), error))

    # If there were any successful files, git commit them
    if project_converted_files:
      subprocess.run("git commit -m 'Android.mk to Android.bp for %s/\n\nModules:\n%s\n\nTest: treehugger'" %
            (git_project, "\n".join(["  " + m for m in project_converted_modules])),
          cwd=git_project, shell=True, capture_output=True)
      changed_projects.append(git_project)
      print("Committed changes to %s" % git_project)

    print("Changed git projects so far:")
    for git_project in changed_projects:
      print("  %s" % git_project)


  print("Changed git projects:")
  for git_project in changed_projects:
    print("  %s" % git_project)


if __name__ == "__main__":
  main(sys.argv)

