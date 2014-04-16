#!/usr/bin/env python3.3
#-*- coding:utf-8 -*-

import subprocess, datetime

t = datetime.datetime.utcnow() - datetime.timedelta(minutes=1)

with open('hosts.txt') as hosts:
    for h in hosts:
        h = h.strip()
        if not h: continue
        r = subprocess.check_output(["ssh", "ubuntu@"+h,
            t.strftime("grep \"%Y/%m/%d %H:%M\""),
            "heartbleed.log"], stderr=subprocess.DEVNULL)
        r = r.decode().split('\n')
        tot = len(r)
        vuln = sum(1 for i in r if 'VULNERABLE' in i)
        mism = sum(1 for i in r if 'MISMATCH' in i)
        vulnf = sum(1 for i in r if 'VULNERABLE [skip: false]' in i)
        print('{:<45}     tot: {:>3}     mism: {:>2}     vuln: {:>2} ({} no-skip)  {:.2%}'.format(
            h+':', tot, mism, vuln, vulnf, float(vuln)/tot))
        for line in (i for i in r if 'VULNERABLE [skip: false]' in i):
            remote = line.split(' ')[2]
            res = subprocess.call(['./tool', remote], stderr=subprocess.DEVNULL, stdout=subprocess.DEVNULL)
            if res != 1:
                print('>>>', remote)
