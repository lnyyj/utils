#!/bin/sh
	for ((i=0; i<10000; i++))
	do
		nohup ./cli &
	done
