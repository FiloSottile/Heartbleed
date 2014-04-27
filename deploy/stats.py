#!/usr/bin/env python3
#-*- coding:utf-8 -*-

import subprocess
import datetime
import re

from collections import Counter

last_minute = datetime.datetime.utcnow() - datetime.timedelta(minutes=1)

errors = Counter()

with open('hosts.txt') as hosts:
    for h in hosts:
        h = h.strip()
        if not h: continue

        r = subprocess.check_output(["ssh", "ubuntu@"+h,
            last_minute.strftime("grep \"%Y/%m/%d %H:%M\""),
            "heartbleed.log"], stderr=subprocess.DEVNULL)
        r = r.decode().split('\n')

        tot = len(r)
        vuln = sum(1 for i in r if 'VULNERABLE' in i)
        mism = sum(1 for i in r if 'MISMATCH' in i)
        vulnf = sum(1 for i in r if 'VULNERABLE [skip: false]' in i)
        print('{:<45}     tot: {:>3}     mism: {:>2}     vuln: {:>2} ({} no-skip)  {:.2%}'.format(
            h+':', tot, mism, vuln, vulnf, float(vuln)/tot))

        # for line in (i for i in r if 'VULNERABLE [skip: false]' in i):
        #     remote = line.split(' ')[2]
        #     res = subprocess.call(['../cli_tool', remote], stderr=subprocess.DEVNULL, stdout=subprocess.DEVNULL)
        #     if res != 1:
        #         print('>>>', remote)

        for line in (i for i in r if 'ERROR' in i):
            msg = line.split('ERROR')[1].strip(' []')
            msg = re.sub(r'\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}:\d{1,5}', '*.*.*.*:*', msg)
            msg = re.sub(r'lookup .*?: ', 'lookup *: ', msg)
            msg = re.sub(r'port tcp/.*', 'port tcp/*', msg)
            msg = re.sub(r'certificate is valid for .*, not .*', 'certificate is valid for *, not *', msg)
            errors[msg] += 1

print()
print('Error summary:')
for error, count in errors.most_common():
    # if count < 10: continue
    print('{:>5}  {}'.format(count, error))
