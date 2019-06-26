#!/usr/bin/env python

import socket
import json
import sys
import time

# if not len(sys.argv) > 1: sys.exit(1)
# strategy=sys.argv[1]
s = socket.socket()
s.connect(('192.168.20.179', 68008))
FILE = "/tmp/status1.json"
while True:
    try:
        b = b''
        tmp = ''
        while 1:
            tmp = s.recv(1024)
            if not tmp:
                break
            if tmp.find("FIN") >= 0:
                b += tmp.strip("FIN")
                break
            b += tmp
        load = json.loads(b)
        open(FILE, 'w').close()
        with open(FILE, 'w') as outfile:
            json.dump(load, outfile, indent=4)
        time.sleep(1)
    except ValueError:
        print "ValueError"
        open('/tmp/ValueError','w').write(b)
    except KeyboardInterrupt:
        print "Socket", s
        if s:
            print "Interrupt recv."
            s.shutdown(socket.SHUT_RDWR)
            time.sleep(1)
            s.close()
            sys.exit(0)
