#! /bin/bash
workers=120
qps=15
minutes=30m
hey -c $workers -q $qps -z $minutes https://nextdemoblog-ean27jt5ha-uc.a.run.app \
& hey -c $workers -q $qps -z $minutes https://justinsblog.web.app
