#!/usr/bin/env python
import os,sys,hashlib,re
import subprocess
from apkverify import is_apk_file

global_mismatch_result = []
onedir = '/mnt/oneimage/'
twodir = '/mnt/twoimage/'

def sumfile(fobj):
    m = hashlib.md5()
    while True:
	d = fobj.read(8096)
	if not d:
		break
	m.update(d)
    return m.hexdigest()

def md5sum(fname):
    if fname == '-':
	ret = sumfile(sys.stdin)
    else:
	try:
	    f = file(fname, 'rb')
	except:
  	    return 'Failed to open file'
  	ret = sumfile(f)
	f.close()
    return ret

def walkdirs(dir1, dir2):
    dirsa = []
    dirsb = []

    for directory in os.walk(dir1):
	relativedirname = pathsep + re.sub(dir1, "", directory[0])
	if relativedirname != pathsep:
		dirsa.append(relativedirname)
    for directory in os.walk(dir2):
	relativedirname = pathsep + re.sub(dir2, "", directory[0])
	if relativedirname != pathsep:
		dirsb.append(relativedirname)
	# Return matches for further inspection.
    return set(dirsa).intersection(set(dirsb))

def complists(dir1, dir2):
    print("#######################################")
    print("Comparing \n%s \n%s" % (dir1, dir2))
    print("#######################################")
    dira = []
    dirb = []
    for file in os.listdir(dir1):
	dira.append(file)
    for file in os.listdir(dir2):
	dirb.append(file)
    indir1 = set(dira).difference(set(dirb))
    indir2 = set(dirb).difference(set(dira))
    return (indir1, indir2, set(dira).intersection(set(dirb)))

def comfiles(files, fdir, sdir):
    firstdir = {}
    secdir = {}
    # Create an absolute path to them from the relative filename and get the md5.
    for f in files:
        of = fdir + pathsep + f
        sf = sdir + pathsep + f
        if os.path.isdir(of) and os.path.isdir(sf):
            continue
        if os.path.islink(of) and os.path.islink(sf):
            continue
        lite_result_check, firstdir[f], lite_errors = is_apk_file(of, validate=True)
        if lite_result_check:
            full_result_check, secdir[f], full_errors = is_apk_file(sf, validate=True)
        else:
            firstdir[f] = md5sum(of)
            secdir[f] = md5sum(sf)
    for x in firstdir:
        print("firstdir %s :::: %s:" % (x, firstdir[x]))
        print("secdir %s :::: %s:" % (x, secdir[x]))
        if firstdir[x] != secdir[x]:
            global_mismatch_result.append(x)
            print("File %s in both targets but does not match!" % x)
    print("\n")

def outp(datum, fdir, sdir):
    if len(datum[0]) > 0:
	print("Items in \"%s\" and not in \"%s\":" % (fdir, sdir))
	for z in datum[0]:
		print("   ", z)
	print("\n------------------------------")
    else:
	print("No unique items in \"%s\"" % (fdir))
    if len(datum[1]) > 0:
	print("Items in \"%s\" and not in \"%s\":" %(sdir, fdir))
	for z in datum[1]:
		print("   ", z)
	print("\n------------------------------")
    else:
	print("No unique items in \"%s\"" % (sdir))
    if len(datum[2]) > 0:
	comfiles(datum[2], fdir, sdir)

def execute(cmd):
    print(cmd)
    p = subprocess.Popen(cmd, shell=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
    out, err = map(lambda b: b.decode('utf-8'), p.communicate())
    print(p.returncode, out, err)
    return p.returncode == 0, out, err

def make_mount_cmd(onedir, twodir, oneimage, twoimage):
    return ['mkdir -p ' + onedir + ' ' + twodir,
            'mount -o rw ' + oneimage + ' ' + onedir,
            'mount -o rw ' + twoimage + ' ' + twodir]

def make_umount_cmd(onedir, twodir):
    return ['umount ' + onedir,
            'umount ' + twodir]

if __name__ == '__main__':
    if len(sys.argv) != 3:
        print("Usage: %s img1 img2" % sys.argv[0])
        sys.exit(1)
    if os.name == "posix":
        pathsep = "/"
    elif os.name == "nt":
        pathsep = "\\"
    else:
        sys.exit(1)

    (oneimage, twoimage) = sys.argv[1], sys.argv[2]
    if oneimage.endswith('.img') and twoimage.endswith('.img'):
        '''default mount point for image '''
        for cmd in make_mount_cmd(onedir, twodir, oneimage, twoimage):
            print(cmd)
            success, out, error_msg = execute(cmd)
            if success:
                continue
            else:
                print(error_msg)
                sys.exit(1)
    elif os.path.isfile(oneimage) and os.path.isfile(twoimage): 
        one_result = ''
        two_result = ''
        lite_result_check, one_result, lite_errors = is_apk_file(onedir, validate=True)
        if lite_result_check:
            full_result_check, two_result, full_errors = is_apk_file(twodir, validate=True)
        else:
            one_result = md5sum(onedir)
            two_result = md5sum(twodir)
        print("onefile %s:::: %s:" % (onedir,one_result))
        print("twofile %s:::: %s:" % (twodir, two_result))
        if  one_result != two_result:
            print("File %s in both targets but does not match!\n" % onedir)
        sys.exit(0)
    else:
        print("please check the input param!")
        sys.exit(1)

    subdirs = walkdirs(onedir, twodir)
    setinfo = complists(onedir, twodir)
    outp(setinfo, onedir, twodir)
    for subdir in subdirs:
        setinfo = complists(onedir + subdir, twodir + subdir)
        outp(setinfo, onedir + subdir, twodir + subdir)

    print("mismatch files:")
    for dir_item in global_mismatch_result:
        print("###########################")
        print("%s " % dir_item)
        print("###########################")

    for cmd in make_umount_cmd(onedir, twodir):
        success, out, error_msg = execute(cmd)
        if success:
            continue
        else:
            print(error_msg)
            sys.exit(1)
