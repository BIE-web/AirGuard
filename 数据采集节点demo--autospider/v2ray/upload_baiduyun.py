import os
from datetime import datetime
date_ = datetime.now().strftime("%Y%m%d")
upload_files = []

for pcap_file in os.listdir():
    if date_ in pcap_file:
            upload_files.append(pcap_file)
            
for upload_file in upload_files:
    print(datetime.now().strftime("%Y-%D-%H:%M:%S") +'[start upload]' + upload_file)
    if 'bro' in upload_file:
        os.system('bypy upload '+ upload_file +' /bro_new/`date +%Y%m%d`/')
    elif 'vpn' in upload_file:
        os.system('bypy upload '+ upload_file +' /v2ray_new/`date +%Y%m%d`/')
    else:
        print(datetime.now().strftime("%Y-%D-%H:%M:%S") +'[ERROR]' + upload_file + ' named unexpect')
    print(datetime.now().strftime("%Y-%D-%H:%M:%S") +'[end upload]' + upload_file)