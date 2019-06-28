#!/bin/bash

echo  "IP={}"  > Rack_Gateway_info.py
source rakesh_script/checkall/variables.sh
for i in $stadd  $all_145 $cash $future 
do
	rack=`grep -w ${i} ip.txt |awk '{print $1}'`
	gateway=$(sshpass -p 'password' ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no USER@192.168.$i "ls -t /home/USER/USER-Application/NSE_OMS_FILE_FO*|head -1| xargs awk 'BEGIN{ORS=\"-\"}/IP 172.19./{split(\$3,a,\".\"); print a[4]}' |sed 's/-$//'")
	echo "IP[\"$i\"]=[\"${rack:-0}\",\"${gateway:-0}\"]" >> Rack_Gateway_info.py
done 
