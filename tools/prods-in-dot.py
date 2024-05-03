import argparse
import os.path
import queue
import re
import subprocess
import sys


def get_args():
  parser = argparse.ArgumentParser()
  parser.add_argument('-v', '--verbose', action='store_true', default=False)
  parser.add_argument('-p', '--product', required=True, help='Word to search for in value of PRODUCT_NAME.')
  parser.add_argument('-d', '--device', required=True, help='The product config .mk file node is highlighted if its PRODUCT_DEVICE value equals to this.')
  parser.add_argument('-o', '--output', default='prods-config.dot')
  parser.add_argument('--inherit', action='store_true', default=True)
  parser.add_argument('--include', action='store_true', default=False)
  parser.add_argument('search_dirs', nargs='*', help='Directories to search for .mk files, e.g. "build/make/target device/google"')

  return parser.parse_args()


def must_run_in_top():
  if not os.path.exists('build/make/core/Makefile'):
    print('Run it from the top of the Android source tree.', file=sys.stderr)
    sys.exit(1)


def main():
  must_run_in_top()
  args = get_args()
  find_command = ['find']
  for dir in args.search_dirs:
    find_command.append(dir.strip())
  find_command += ['-name', '*.mk', '-type', 'f', '-print']
  grep_command = ['xargs', 'grep', '-l', 'PRODUCT_NAME := .*' + args.product + '.*', ]

  if args.verbose:
    print(find_command)
    print(grep_command)

  find_result = subprocess.Popen(find_command, encoding='utf-8', stdout=subprocess.PIPE)
  grep_result = subprocess.Popen(grep_command, encoding='utf-8', stdin=find_result.stdout, stdout=subprocess.PIPE)
  outs, err = grep_result.communicate()
  if err:
    print(err)
    sys.exit(1)

  mk_files = queue.Queue()
  roots = outs.strip().split('\n')
  for f in roots:
    mk_files.put(f)

  p0 = '\s*#.*'
  p1 = '\s*\$\(call inherit-product.*,(.*)\)$'
  p2 = '\s*-?include (.*)'
  p3 = '\s*PRODUCT_DEVICE := (.*)'
  mk_files_map = {}
  while not mk_files.empty():
    f = mk_files.get()
    if f in mk_files_map.keys():
      continue
    mk_files_map[f] = {'name': f}
    data = mk_files_map[f]
    with open(f) as lines:
      for line in lines:
        match = re.match(p0, line)
        if match:
          continue

        # inherit
        if args.inherit:
          match = re.match(p1, line)
          if match:
            if 'inherit' not in data:
              data['inherit'] = []
            filepath = replace_var(match.group(1))
            data['inherit'].append(filepath)
            if os.path.isfile(filepath) and filepath not in mk_files_map.keys():
              mk_files.put(filepath)

        # include
        if args.include:
          match = re.match(p2, line)
          if match:
            if 'include' not in data:
              data['include'] = []
            filepath = replace_var(match.group(1))
            data['include'].append(filepath)
            if os.path.isfile(filepath) and filepath not in mk_files_map.keys():
              mk_files.put(filepath)

        # PRODUCT_DEVICE
        match = re.match(p3, line)
        if match:
          data['device'] = match.group(1)

  # Print
  if args.verbose:
    for f in mk_files_map.values():
      # print out
      print(f['name'])
      for (k, v) in f.items():
        if k == 'device':
          print(f'\t{k}: {v}')
        if k not in ('name', 'device'):
          print(f'\t{k}: [')
          for e in v:
            print(f'\t\t{e.strip()}')
          print('\t]')

  # generate dot
  content = [f'digraph "Products config"' + '{']
  content.append('node [shape=rectangle]')
  content.append('ranksep=1')
  for f in mk_files_map.values():
    for (k, v) in f.items():
      if k == 'name':
        group = ""
        if f["name"].startswith('vendor/google/products/'):
          group = 'vendor/google/products'
        elif f["name"].startswith('device/google/'):
          pos_of_3rd_slash = f["name"].index('/', len('device/google/'))
          group = f["name"][:pos_of_3rd_slash]
        elif f["name"].startswith('build/make/target/product'):
          group = 'build/make/target/product'
        elif f["name"].startswith('vendor/google_devices/'):
          pos_of_3rd_slash = f["name"].index('/', len('vendor/google_devices/'))
          group = f["name"][:pos_of_3rd_slash]
        elif f["name"].startswith('vendor/google/'):
          pos_of_3rd_slash = f["name"].index('/', len('vendor/google/'))
          group = f["name"][:pos_of_3rd_slash]
        elif f["name"].startswith('external/'):
          pos_of_2dn_slash = f["name"].index('/', len('external/'))
          group = f["name"][:pos_of_2dn_slash]

        color = ''
        shape = ''
        style = ''
        fontcolor = ''
        if v in roots:
          color = 'color=green'
          shape = 'shape=oval'

        if f.get('device') == args.device:
          color = 'color=blue'
          style = 'style=filled'
          fontcolor = 'fontcolor=white'

        content.append(f'"{f["name"]}" [group="{group}" {color} {shape} {style} {fontcolor}]')
      if args.inherit and k == 'inherit':
        inherits = '"' + '" "'.join(v) + '"'
        content.append(f'"{f["name"]}" -> {{ {inherits} }} [arrowhead = empty]')
      if args.include and k == 'include':
        includes = '"' + '" "'.join(v) + '"'
        content.append(f'"{f["name"]}" -> {{ {includes} }}  [arrowhead = open, style=dashed]')
  content.append('}')

  with open(args.output, 'w', encoding="utf-8") as dot:
    dot.write('\n'.join(content))


def replace_var(str):
  return str.strip().replace('$(SRC_TARGET_DIR)', 'build/make/target')


if __name__ == '__main__':
  main()
