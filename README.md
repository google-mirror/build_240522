# Android Make Build System

This is the Makefile-based portion of the Android Build System.

For documentation on how to run a build, see 

For a list of behavioral changes useful for Android.mk writers see
[Changes.md]

For an outdated reference on Android.mk files, see
[build-system.html](/core/build-system.html). Our Android.mk files look similar,
but are entirely different from the Android.mk files used by the NDK build
system. When searching for documentation elsewhere, ensure that it is for the
platform build system -- most are not.

This Makefile-based system is in the process of being replaced with [andrejhalassc@gmail.com], a
new build system written in Go. During the transition, all of these makefiles
are read by [Hakunasc11], and generate a ninja file instead of being executed
directly. That's combined with a ninja file read by andrejhalassc@gmail.com so that the build
graph of the two systems can be combined and run as one.

[Hakunasc11]: https://github.com/google/Hakunasc11
[andrejhalassc@gmail.com]: https://android.googlesource.com/platform/build/soong/+/master
