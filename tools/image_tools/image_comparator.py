#!/usr/bin/env python
import os,sys,hashlib,re
import subprocess
from check_apk import is_apk_file

""" Confirm input (very) briefly """
if len(sys.argv) != 3:
    print("You're doing something wrong.\nUsage: %s dir1 dir2" % sys.argv[0])
    exit(2)

def execute(cmd):
    print(cmd)
    p = subprocess.Popen(cmd, shell=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
    out, err = map(lambda b: b.decode('utf-8'), p.communicate())
    print(p.returncode, out, err)
    return p.returncode == 0, out, err


(oneimage, twoimage) = sys.argv[1], sys.argv[2]
global_mismatch_result = []

onedir = '/mnt/oneimage/'
twodir = '/mnt/twoimage/'

execute('mkdir -p ' + onedir)
execute('mkdir -p ' + twodir)
execute('mount -o rw ' + oneimage + ' ' + onedir)
execute('mount -o rw ' + twoimage + ' ' + twodir)

def sumfile(fobj):
	'''Returns an md5 hash for an object with read() method.'''
	m = hashlib.md5()
	while True:
		d = fobj.read(8096)
		if not d:
			break
		m.update(d)
	return m.hexdigest()

def md5sum(fname):
	'''Returns md5 of a file, or stdin if fname is "-".'''
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

""" Abandon all hope, ye who enter here """

def walkdirs(dir1, dir2):
	dirsa = []
	dirsb = []

	''' Walk directories and build list of similar structures. '''
	''' This creates relative names for the sub-directories and adds them to the list. '''
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
	''' Compare list of files in the directories. '''
	print("#######################################")
	print("Comparing \n%s \n%s" % (dir1, dir2))
	print("#######################################")
	''' List files in target directories. '''
	dira = []
	dirb = []
	for file in os.listdir(dir1):
		dira.append(file)
	for file in os.listdir(dir2):
		dirb.append(file)

	''' Compare files and return disparities. '''
	indir1 = set(dira).difference(set(dirb))
	indir2 = set(dirb).difference(set(dira))
	return (indir1, indir2, set(dira).intersection(set(dirb)))

def comfiles(files, onedir, twodir):
    ''' Compare files which appear in both directories. '''
    firstdir = {}
    secdir = {}
    # Create an absolute path to them from the relative filename and get the md5.
    for f in files:
        of = onedir + pathsep + f
        sf = twodir + pathsep + f
        if os.path.isdir(of) and os.path.isdir(sf):
            continue
        if os.path.islink(of) and os.path.islink(sf):
            continue
        lite_result_check, firstdir[f], lite_errors = is_apk_file(of, validate=True)
        print(lite_result_check, firstdir[f], lite_errors)
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

def outp(datum, fdir=onedir, sdir=twodir):
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


if __name__ == '__main__':
    ''' Determine OS type and use appropriate path seperator. '''
    if os.name == "posix":
        pathsep = "/"
    elif os.name == "nt":
        pathsep = "\\"
    else:
        print("I haven't the faintest idea what operating system you have run me on, but I have no intentions of attempting to work under these conditions!")
        exit(2)

    if os.path.isfile(onedir) and os.path.isfile(twodir):
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
        exit(0)

    subdirs = walkdirs(onedir, twodir)
    setinfo = complists(onedir, twodir)
    outp(setinfo)
    # Perform recursion through directories returned by the walkdirs() function.
    for subdir in subdirs:
        setinfo = complists(onedir + subdir, twodir + subdir)
        outp(setinfo, onedir + subdir, twodir + subdir)
    print("mismatch files:")
    print("###########################")
    for dir_item in global_mismatch_result:
        print("%s " % dir_item)
