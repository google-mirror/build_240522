#!/bin/bash
# This script will measure the memory usage of a process in an easy to read
# format that can be diffed before and after to see the impact of a change
# It also gathers the size of the files generated and installed on the device.

usage() {
  echo "Usage: ./memory_check.sh <before|after|summary> [#: iterations] [processes]"
  echo "Example:       ./memory_check.sh before 4 audioserver"
  echo "               # Rebuild with change, flash and boot device."
  echo "               ./memory_check.sh after 4 audioserver"
  echo "               ./memory_check.sh summary"
  exit 1
}

if [ $# != 3 ]; then
  if [ $# != 1 ] || [ $1 != "summary" ]; then
    usage
  fi
fi
if [ $1 == "summary" ]; then
  echo trying to get the summary!
  exit 0
fi

# is this before or after the change?
current=$1
# number of reboots/samples
iterations=$2
# the process to measure
process=$3

mkdir -p $current
for (( i = 1; i <= $iterations; i++ ))
do
echo $(date) "Rebooting device"
adb shell reboot
sleep 30
echo $(date) "Waiting for device"
adb wait-for-device
echo $(date) "Adb root"
adb root
adb wait-for-device
sleep 10
echo $(date) "Gathering data"
adb shell "showmap \$(pidof $process)" | grep -v "\-\-\-"  > "$current/showmap$i"
cat "$current/showmap$i" | grep -v "\[" | grep -v TOTAL | awk '{rss += $2; pss += $3; pc += $6; pd += $7;}END{printf "Sum of rss %d, pss %d, pc %d, pd %d\n", rss, pss, pc, pd}' >> "$current"/libonly_results.txt
done

avg_pss=$(awk '{pss += $6; cnt +=1}END{printf "%d", pss/cnt}' "$current"/libonly_results.txt)
echo Average PSS for libs only $avg_pss
cnt=$(awk '{cnt +=1}END{printf "%d", cnt}' "$current"/libonly_results.txt)
std_dev=$(awk "function abs(v) {return v < 0 ? -v : v} {pss += (abs(\$6 - $avg_pss))^2}END{printf \"Standard deviation of lib-only results %d\", sqrt(pss/$cnt)}" "$current"/libonly_results.txt)
echo $std_dev

for (( i = 1; i <= $iterations; i++ ))
do
cat "$current/showmap$i" | awk '{rss = $2; pss = $3; sc = $4; sd = $5; pc = $6; pd = $7;}END{printf "Sum of rss %d, pss %d, sc %d, sd %d, pc %d, pd %d. Sum of sc/sd/pc/pd %d, sc/pc %d\n", rss, pss, sc, sd, pc, pd, sc + sd + pc + pd, sc + pc}' >> "$current"/full_results.txt
done

avg_pss=$(awk '{pss += $6; cnt +=1}END{printf "%d", pss/cnt}' "$current"/full_results.txt)
echo Average PSS including anon/scudo $avg_pss
cnt=$(awk '{cnt +=1}END{printf "%d", cnt}' "$current"/full_results.txt)
std_dev=$(awk "function abs(v) {return v < 0 ? -v : v} {pss += (abs(\$6 - $avg_pss))^2}END{printf \"%d\", sqrt(pss/$cnt)}" "$current"/full_results.txt)
echo Standard deviation of full results $std_dev

find $ANDROID_PRODUCT_OUT -name "installed-files*.txt" | xargs cat > "$current"/all_installed_files.txt
cat "$current"/all_installed_files.txt | awk '{size += $1}END{printf "Total size of installed files: %d\n", size}' > "$current"/total_size_installed_files.txt

