#!/usr/bin/env python3

"""
The complete list of the remaining Make files in each partition for all lunch targets

How to run?
python3 $(path-to-file)/mk2bp_partition.py
"""

from pathlib import Path

import csv
import datetime
import os
import shutil
import subprocess
import time

class UserError(Exception):
  pass

def get_top():
  path = '.'
  while not os.path.isfile(os.path.join(path, 'build/soong/soong_ui.bash')):
    if os.path.abspath(path) == '/':
      raise UserError('Could not find android source tree root.')
    path = os.path.join(path, '..')
  return os.path.abspath(path)

# get the values of a build variable
def get_build_var(variable, product, build_variant):
  """Returns the result of the shell command get_build_var."""
  env = {
      **os.environ,
      'TARGET_PRODUCT': product if product else '',
      'TARGET_BUILD_VARIANT': build_variant if build_variant else '',
  }
  return subprocess.run([
      'build/soong/soong_ui.bash',
      '--dumpvar-mode',
      variable
  ], check=True, capture_output=True, env=env).stdout.decode('UTF-8').strip()

# get a full list of lunch targets
def get_all_lunch_targets():
    return get_build_var("all_named_products", "", "")

def get_count_lunch_targets():
    count = 0
    lunch_targets = get_all_lunch_targets().split()
    for lunch_target in lunch_targets:
        count += 1
    return count

def get_installed_product_out(product, build_variant):
    return get_build_var("PRODUCT_OUT", product, build_variant)

def get_make_file_partitions():
    total_lunch_targets = get_count_lunch_targets()
    lunch_targets = get_all_lunch_targets().split()
    makefile_by_partition = dict()
    partitions = set()
    current_count = 0
    start_time = time.time()
    # cannot run command `m lunch_target`
    obsolete_targets = {"mainline_sdk", "ndk"}
    for lunch_target in lunch_targets:
        current_count += 1
        current_time = time.time()
        print (current_count, "/", total_lunch_targets, lunch_target, datetime.timedelta(seconds=current_time - start_time))
        if lunch_target in obsolete_targets:
            continue
        installed_product_out = get_installed_product_out(lunch_target, "")
        filename = installed_product_out + "/mk2bp_remaining.csv"
        copy_filename = installed_product_out + "/" + lunch_target + "_mk2bp_remaining.csv"
        # only generate if not exists
        if not os.path.exists(copy_filename):
            bash_cmd = "bash build/soong/soong_ui.bash --make-mode TARGET_PRODUCT=" + lunch_target
            bash_cmd += " TARGET_BUILD_VARIANT= " + filename
            subprocess.run(bash_cmd, shell=True, text=True, stdout=subprocess.DEVNULL)
            # generate a copied .csv file, to avoid possible overwrittings
            with open(copy_filename, "w") as file:
                shutil.copyfile(filename, copy_filename)

        # open mk2bp_remaining.csv file
        with open(copy_filename, "r", errors="ignore") as csvfile:
            reader = csv.reader(csvfile, delimiter=",", quotechar='"')
            next(reader, None)
            data_read = [row for row in reader]
            for data_row in data_read:
                # read partition information
                partition = data_row[2]
                if partition not in makefile_by_partition:
                    makefile_by_partition.setdefault(partition, set())
                makefile_by_partition[partition].add(data_row[0])
                partitions.add(partition)

    # write merged make file list for each partition into a csv file
    installed_path = Path(installed_product_out).parents[0].as_posix()
    csv_path = installed_path + "/mk2bp_partition.csv"
    with open(csv_path, "wt") as csvfile:
        writer = csv.writer(csvfile, delimiter=",")
        count_makefile = 0
        for partition in sorted(partitions):
            number_file = len(makefile_by_partition[partition])
            count_makefile += number_file
            row = [partition, number_file]
            writer.writerow(row)
            for makefile in sorted(makefile_by_partition[partition]):
                row = [makefile]
                writer.writerow(row)
        row = ["The total count of make files is ", count_makefile]
        writer.writerow(row)

def main():
    os.chdir(get_top())
    get_make_file_partitions()

if __name__ == "__main__":
    main()
