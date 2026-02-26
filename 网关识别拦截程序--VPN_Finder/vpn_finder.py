import json
import socket
import pandas as pd
import multiprocessing  
import warnings
import os
import joblib
import time


warnings.filterwarnings("ignore")

ip_que = multiprocessing.Queue()
rep_que = multiprocessing.Queue()
ana_que = multiprocessing.Queue()  
with open('config.json', 'r') as f:
    config = json.loads(f.read())


feature_list = [
		"srcIP", "srcPort", "dstIP", "dstPort", "l4Proto",
		"connSipDip", "connSipDprt",
		"aveIat", "minIat", "maxIat", "iqrIat", "q1Iat", "q2Iat", "q3Iat", "stdIat", "varIat",
		"avePktSz", "minPktSz", "maxPktSz", "iqrPktSz", "q1PktSz", "q2PktSz", "q3PktSz", "stdPktSz", "varPktSz", "modePktSz",
		"bytps", "pktps",
		"bytAsm", "pktAsm",
		"ipToS",
		"ipMinTTL", "ipMaxTTL", "ipTTLChg ",
		"ipAbnormalLenth",
		"l4AbnormalLenth",
		"tcpPSeqCnt", "tcpPAckCnt",
		"tcpAveWinSz", "tcpMinWinSz", "tcpMaxWinSz", "tcpWinSzDwnCnt", "tcpWinSzUpCnt", "tcpWinSzChgDirCnt",
		"tcpFlag",
]
label_list = ['srcIP','srcPort','dstIP','dstPort','l4Proto']

udp_featrue_list = [
    "connSipDip", "connSipDprt",
    "aveIat", "minIat", "maxIat", "q2Iat", "stdIat", "varIat",
    "avePktSz", "minPktSz", "maxPktSz", "q2PktSz", "stdPktSz", "varPktSz",
    "bytps", "pktps",
    "bytAsm", "pktAsm",
    "ipToS",
    "ipMinTTL", "ipMaxTTL", "ipTTLChg ",
    "ipAbnormalLenth",
]


udp_clf = joblib.load('model/udp.pkl')

class ip_Consumer(multiprocessing.Process):  

    def __init__(self, ip_que):  
        super(ip_Consumer, self).__init__()  
        self.ip_que = ip_que
        self.ip_list = []
        self.start() 

    def run(self):
        while True:
            ip = self.ip_que.get()
            if ip not in  self.ip_list:
                os.system(("ssh root@" +config['router_ip'] + " iptables -A FORWARD -d " + ip + ' -j DROP'))
                os.system(("ssh root@" +config['router_ip'] + " iptables -A FORWARD -s " + ip + ' -j DROP'))
                self.ip_list.append(ip)
                msg = '<b class="error">[{date}][ATTENTION] BAN IP: {IP} </b>\n'.format( 
                    date=time.strftime('%m-%d %H:%M:%S', time.localtime(time.time())),IP=ip)
                rep_que.put(msg)

class rep_Consumer(multiprocessing.Process):  

    def __init__(self, rep_que):  
        super(rep_Consumer, self).__init__()  
        self.rep_que = rep_que
        self.start() 

    def run(self): 
        while True:
            rep = self.rep_que.get()
            with open('./log','a') as f:
                f.write(rep)


class ana_Consumer(multiprocessing.Process):  

    def __init__(self, ana_que):  
        super(ana_Consumer, self).__init__()  
        self.ana_que = ana_que  
        self.start() 

    def run(self): 
        while True:
            try:
                flow = self.ana_que.get()
                df = pd.DataFrame([dict(zip(feature_list,flow.split(" ")))])
                label = df[label_list]
                feature = df[udp_featrue_list]
                if label['l4Proto'][0] == 'UDP': 
                    result = udp_clf.predict(feature)
                    if result[0] == 'v2r':
                        msg = '<b>[{date}][ATTENTION] Find vpn flow {proto} {src}:{sport} -> {dst}:{dport}\n</b>'.format( 
                        date=time.strftime('%m-%d %H:%M:%S', time.localtime(time.time())),src=label['srcIP'][0],sport=label['srcPort'][0],dst=label['dstIP'][0],dport=label['dstPort'][0],proto=label['l4Proto'][0])
                        rep_que.put(msg)
                        ip_que.put(label['dstIP'][0])
                else:
                    pass
            except:
                    pass


def recvived(address, port):
    # 文件缓冲区
    Buffersize = 4096*10

    while True:
        udp_socket = socket.socket(socket.AF_INET,socket.SOCK_DGRAM)
        udp_socket.bind((address, port))
        recv_data = udp_socket.recvfrom(Buffersize)
        flowList =  recv_data[0].decode('UTF-8').split("\n")
        msg = "[{date}][INFO] Find {cnt} flows\n".format(date=time.strftime('%m-%d %H:%M:%S', time.localtime(time.time())), cnt=len(flowList))
        rep_que.put(msg)
        for flow in flowList:
            ana_que.put(flow)
        udp_socket.close()

def run(ip,port):
    ip_Consumer(ip_que)
    rep_Consumer(rep_que)
    ana_Consumer(ana_que)
    recvived(ip,port)