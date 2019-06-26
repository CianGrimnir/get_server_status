# GO-Monitoring tools to get server details.

## get_server_status_GO

### HOW-TO:

Generate Rack 				Rack_Gateway_Creation.sh
Generate variable 			Variable_File.sh

1.start script from 192.168.*.*		Web_Monitor_DataGenerator.py
2.start client from 192.168.*.*		WebMonitor/Web_Monitor_JSON_Client.py
3.start WebServer GoScript		WebMonitor/WebMonitor_Server.go

To View Monitoring:
http://192.168.CLIENT.IP:8080/CR		FOR CR
http://192.168.CLIENT.IP:8080/STADD		FOR Straddler

Use Button (TOP) to SWAP between both

To View Server LOGS:
/tmp/Monitor_Web.log