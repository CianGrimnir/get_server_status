# GO-Monitoring tools to get server details.

## get_server_status_GO

### HOW-TO:

<pre>
* Generate Rack 	  		        bash Rack_Gateway_Creation.sh
* Generate variable 			        bash Variable_File.sh

1. Start script from 192.168.*.*		python Web_Monitor_DataGenerator.py

2. Start client from 192.168.*.*		python WebMonitor/Web_Monitor_JSON_Client.py

3. Start WebServer GoScript		        go run WebMonitor/WebMonitor_Server.go


* To View Monitoring:
* http://192.168.CLIENT.IP:8080/CR		FOR CR
* http://192.168.CLIENT.IP:8080/STADD		FOR Straddler

</pre>
* Use Button (TOP) to SWAP between both

* To View Server LOGS:    &nbsp;&nbsp;&nbsp; `/tmp/Monitor_Web.log`
