#
# Base config.fs file that specifies AID values to nice names
# as well as interpreter and homedir (see section DEFAULT).
#

[root]
value: 0

[system]
value: 1000

[radio]
value: 1001

[bluetooth]
value: 1002

[graphics]
value: 1003

[input]
value: 1004

[audio]
value: 1005

[camera]
value: 1006

[log]
value: 1007

[compass]
value: 1008

[mount]
value: 1009

[wifi]
value: 1010

[adb]
value: 1011

[install]
value: 1012

[media]
value: 1013

[dhcp]
value: 1014

[sdcard_rw]
value: 1015

[vpn]
value: 1016

[keystore]
value: 1017

[usb]
value: 1018

[drm]
value: 1019

[mdnsr]
value: 1020

[gps]
value: 1021

[media_rw]
value: 1023

[mtp]
value: 1024

[drmrpc]
value: 1026

[nfc]
value: 1027

[sdcard_r]
value: 1028

[clat]
value: 1029

[loop_radio]
value: 1030

[mediadrm]
value: 1031

[package_info]
value: 1032

[sdcard_pics]
value: 1033

[sdcard_av]
value: 1034

[sdcard_all]
value: 1035

[logd]
value: 1036

[shared_relro]
value: 1037

[dbus]
value: 1038

[tlsdate]
value: 1039

[mediaex]
value: 1040

[audioserver]
value: 1041

[metrics_coll]
value: 1042

[metricsd]
value: 1043

[webserv]
value: 1044

[debuggerd]
value: 1045

[mediacodec]
value: 1046

[cameraserver]
value: 1047

[firewall]
value: 1048

[trunks]
value: 1049

[nvram]
value: 1050

[shell]
value: 2000

[cache]
value: 2001

[diag]
value: 2002

[net_bt_admin]
value: 3001

[net_bt]
value: 3002

[inet]
value: 3003

[net_raw]
value: 3004

[net_admin]
value: 3005

[net_bw_stats]
value: 3006

[net_bw_acct]
value: 3007

[net_bt_stack]
value: 3008

[readproc]
value: 3009

[wakelock]
value: 3010

[everybody]
value: 9997

[misc]
value: 9998

[nobody]
value: 9999

# These apply accross all config files and can only be set in this
# base config file. DEFAULT set in other config files are scoped
# to that file only.
[DEFAULT]
homedir: /
interpreter: /system/bin/sh
