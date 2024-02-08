import argparse
import os
import subprocess

from typing import Dict, Text

def build_android_metadata(argv):
  argparser = argparse.ArgumentParser()
  argparser.add_argument('--dist_dir')
  args = argparser.parse_args()

  build_command = []
  build_command.append('time')
  build_command.append('./build/soong/soong_ui.bash')
  build_command.append('--make-mode')
  build_command.append('all_teams')
  build_command.append('DIST_DIR=' + args.dist_dir)
  build_command.append('TARGET_PRODUCT=aosp_x86_64')
  build_command.append('TARGET_RELEASE=trunk_staging')

  run_command(build_command, print_output=True)

def run_command(
    args: list[str],
    env: Dict[Text, Text] = os.environ,
    print_output: bool = False,
):
  result = subprocess.run(
      args=args,
      text=True,
      capture_output=True,
      check=False,
      env=env,
  )
  # If the process failed, print its stdout and propagate the exception.
  if not result.returncode == 0:
    print('Build command failed! output:')
    print('stdout: ' + result.stdout)
    print('stderr: ' + result.stderr)

  result.check_returncode()

  if print_output:
    print(result.stdout)

def main(argv):
  build_android_metadata(argv)
