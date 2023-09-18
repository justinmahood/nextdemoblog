#! /bin/bash
workers=100
qps=15
minutes=5m
hey -c $workers -q $qps -z $minutes https:// \ # Uncached service
	&
hey -c $workers -q $qps -z $minutes https:// # Cached service
