#!/bin/bash

# number of reboots/samples
ITERATIONS=3

for (( i = 1; i <= $ITERATIONS; i++ ))
do
echo $(date) "Rebooting device"
adb shell reboot
sleep 60
echo $(date) "Waiting for device"
adb wait-for-device
sleep 10
echo $(date) "Adb root"
adb root
sleep 10
echo $(date) "Gathering data"
adb shell 'showmap $(pidof audioserver)' | grep -v "\-\-\-"  > "showmap$i"
cat "showmap$i" | grep -v "\[" | grep -v TOTAL | awk '{rss += $2; pss += $3; pc += $6; pd += $7;}END{printf "Sum of rss %d, pss %d, pc %d, pd %d\n", rss, pss, pc, pd}' >> libonly_results.txt
done

avg_pss=$(awk '{pss += $6; cnt +=1}END{printf "%d", pss/cnt}' libonly_results.txt)
echo Average PSS for libs only $avg_pss
cnt=$(awk '{cnt +=1}END{printf "%d", cnt}' libonly_results.txt)
std_dev=$(awk "function abs(v) {return v < 0 ? -v : v} {pss += (abs(\$6 - $avg_pss))^2}END{printf \"Standard deviation of lib-only results %d\", sqrt(pss/$cnt)}" libonly_results.txt)
echo $std_dev

for (( i = 1; i <= $ITERATIONS; i++ ))
do
cat "showmap$i" | awk '{rss = $2; pss = $3; sc = $4; sd = $5; pc = $6; pd = $7;}END{printf "Sum of rss %d, pss %d, sc %d, sd %d, pc %d, pd %d. Sum of sc/sd/pc/pd %d, sc/pc %d\n", rss, pss, sc, sd, pc, pd, sc + sd + pc + pd, sc + pc}' >> full_results.txt
done

avg_pss=$(awk '{pss += $6; cnt +=1}END{printf "%d", pss/cnt}' full_results.txt)
echo Average PSS including anon/scudo $avg_pss
cnt=$(awk '{cnt +=1}END{printf "%d", cnt}' full_results.txt)
std_dev=$(awk "function abs(v) {return v < 0 ? -v : v} {pss += (abs(\$6 - $avg_pss))^2}END{printf \"%d\", sqrt(pss/$cnt)}" full_results.txt)
echo Standard deviation of full results $std_dev

find $ANDROID_PRODUCT_OUT -name "installed-files*.txt" | xargs cat > all_installed_files.txt
cat all_installed_files.txt | awk '{size += $1}END{printf "Total size of installed files: %d\n", size}' > total_size_installed_files.txt

