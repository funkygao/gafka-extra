#!/bin/sh

# master
for i in {10300..10800}; do v=$(dig +short w$i.wdds.redis.com); if [ $v ] ;then echo $v：$i; fi; done

# slave
# for i in {10300..10800}; do v=$(dig +short r$i.wdds.redis.com); if [ $v ] ;then echo $v：$i; fi; done

# distribution of redis instances on hosts
for i in {10300..10800}; do v=$(dig +short w$i.wdds.redis.com); if [ $v ] ;then echo $v; fi; done | sort | uniq -c | sort -n
