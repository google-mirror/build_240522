# Android Make Build System

This is the Makefile-based portion of the Android Build System.

For documentation on how to run a build, see [Usage.txt](Usage.txt)

For a list of behavioral changes useful for Android.mk writers see
[Changes.md](Changes.md)

This Makefile-based system is in the process of moving to [Soong], which is
written in Go. All of the makefiles are intepreted by [Kati], and generate ninja
files instead of being executed directly. That ninja file is combined with one
from Soong, so that part of the build graph can be produced by Soong and part by
Kati, while all being built as a single graph.

[Kati]: https://github.com/google/kati
[Soong]: https://android.googlesource.com/platform/build/soong/+/master
