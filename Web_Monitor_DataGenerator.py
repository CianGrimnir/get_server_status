#!/usr/bin/env python

import subprocess
import os
import sys
import string
import threading
import errno
import logging
import signal
import time
import multiprocessing
import json
import socket
import psycopg2
from multiprocessing import Pool

blaclist = []
a = []
global Ferror
Ferror = {}
result = {}
serverdict = {}
jsondata = {}
server1 = []
USER = 'user'
DB = 'dbname'
PASSWD = 'password'
PORT = '5432'
pun = string.punctuation

# To make python unbuffered, other option -> run script via 'python -u'
os.environ["PYTHONUNBUFFERED"] = "1"


def timeout(p):												# subprocess don't have builtin timeout function to terminate session, custom timeout to handle it
    if p.poll() is None:
        global Ferror
        global exi
        try:
            p.kill()
            Ferror[ip.strip()] = [ip.strip('\n'), 'unreachable']
            exi = 0
        except OSError as e:
            if e.errno != errno.ESRCH:
                raise


def INT_HANDLER():											# SIGNAL INTERRUPT HANDLER
    signal.signal(signal.SIGINT, signal.SIG_IGN)


def subprocess_function(j):
    global ip
    global Ferror
    Ferror = {}
    a = []
    ip = j
    global exi
    exi = 1
    timer = 2.0
    # if '.'.join(j.strip().split('.')[-2:]) in blaclist:
    # 	timer=4.0
    try:
        conn = psycopg2.connect(database=DB, user=USER, password=PASSWD, host=j.strip(
        ), port=PORT, connect_timeout=2)
    except:
        Ferror[j.strip()] = [j.strip('\n'), 'unreachable', 0, 0]
        return Ferror
    cur = conn.cursor()
    query = "select (select count(distinct order_number ) from ms_oe_request_fo)  , (select count(distinct response_order_number) from ms_trade_confirm_fo) ;"
    cur.execute(query)
    data = cur.fetchmany()
    conn.close()
    order = int(data[0][0])
    trades = int(data[0][1])
    host = "USER@" + j.strip()
    command1 = '''(echo "import os";echo "name='USER_Application'";echo "r=os.popen('ps -ef').read().strip().split(\n)";echo "print [r[i] for i in range(len(r)) if name in r[i]]")| python'''
    ssh = subprocess.Popen(["ssh", "%s" % host, command1],
                           shell=False, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
    t = threading.Timer(timer, timeout, [ssh])								# creating thread timer
    # starting thread, timeout function will be called after time defined in timer variable
    t.start()
    result1 = ssh.stdout.readlines()
    error = ssh.stderr.readlines()
# 	print(ip,result,error,ssh.poll(),ssh.returncode)
    result2 = result1
    ll = str(''.join([o for o in result2 if o not in pun]).replace('\n', ''))
    for word in ll.split():
        a.append(word)
    for i in range(0, 1):
        try:
            if a[i].find('name'):
                a.pop(i)
        except IndexError:
            pass
# 	print(j,a,exi,error,t.is_alive())
    if len(a) > 0:
        Ferror[j.strip()] = [j.strip('\n'), 'running', order, trades]
    elif not a and exi != 0 and not error and ssh.poll() != 0:
        Ferror[j.strip()] = [j.strip('\n'), 'unreachable', 0, 0]
        return Ferror
    elif not a and exi != 0 and not error and t.is_alive():
        Ferror[j.strip()] = [j.strip('\n'), 'not running', order, trades]
    try:
        if error[0].startswith('ssh:') and exi != 0:
            Ferror[j.strip()] = [j.strip('\n'), 'unreachable', 0, 0]
    except IndexError:
        pass
    if len(Ferror.items()) > 0:									# return empty dict if Ferror is empty, else return actual value
        return Ferror
    else:
        return {}


def main():
    strategy = ["STADD", "CR"]
    i = 0
    s = socket.socket()
    s.bind(('1.1.20.179', 12345))
    s.listen(8)
    while True:
        strategy = ["STADD", "CR"]
        result = {}
        jsondata = {}
        Final_jsondata = {}
        Final_jsondata['monitor'] = []
        server1 = []
        c, addr = s.accept()
        print "connection from: " + str(addr)
        # write std.error to file mentioned
        sys.stderr = open(
            '/home/mint/DAILY_CHECK/user_script/checkall/test/PoolHandler.log', 'a')
        # write strings to stderr
        sys.stderr.write("\n" + sys.argv[0] + "\t" + time.ctime() + "\n\n")
        # LOGGING multiprocess threads
        multiprocessing.log_to_stderr(logging.DEBUG)
        # create pool of n process, having INT_HANDLER as initializer function
        pool = Pool(12, INT_HANDLER)
        while True:
            result = {}
            jsondata = {}
            server1 = []
            if strategy[i] == "CR":
                with open("/home/mint/DAILY_CHECK/user_script/checkall/test/hostfile", "r") as f:
                    server = f.readlines()
            elif strategy[i] == "STADD":
                with open("/home/mint/DAILY_CHECK/user_script/checkall/test/hostfile1", "r") as f:
                    server = f.readlines()

            try:
                # fetch return dict from subprocess_function
                for Ferror in pool.imap(subprocess_function, server):
                    try:
                        # if returned dict not empty store it in result dict
                        if len(Ferror.items()) > 0:
                            for k, v in Ferror.items():
                                result[k] = v
                    except AttributeError:
                        pass
                pass
            except KeyboardInterrupt:									# handle KeyboardInterrupt, it works because INT_HANDLER is initialized with each process
                print('\n\033[0m\033[1;31m' +
                      'SIGINT signal recieved\033[1;m\n')
                sys.exit(0)

            jsondata["name"] = strategy[i]
            for k, v in result.items():									# print error along with it's IP at the end
                m = list(result[k])
                if len(m) == 4:
                    m = result[k]
                    server1.append(
                        {"ip": m[0], "status": m[1], "order": m[2], "trades": m[3]})
                    continue
            jsondata["server"] = server1
            Final_jsondata["monitor"].append(jsondata)
            if i == 1:
                i = 0
                senddata = json.dumps(Final_jsondata).encode('utf-8')
                try:
                    c.settimeout(3.0)
                    c.send(senddata)
# 				print("SEND: ",senddata)
                    c.send("FIN")
                except socket.error, e:
                    print e
                    if isinstance(e.args, tuple):
                        print "errno is %d" % e[0]
                        if e[0] == errno.EPIPE:
                            print "Remote disconnection detected"
                            c.close()
                            break
                Final_jsondata['monitor'] = []
                continue
            i += 1
        c.shutdown(socket.SHUT_RDWR)
        c.close()


if __name__ == "__main__":
    main()
