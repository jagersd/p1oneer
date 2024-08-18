# P1ONEER

This is yet another personal project repository with learning purposes as primary aim and an effective product secondary. 

**What is p1oneer might become:**
A JSON configurable multi processor for container deployed applications (think nginx + php-fpm) which runs at pid 1 and it forks out to the other defined processes. When a long running process exits, pid1 does so as well. It  should also support processes which are supposed to exit like startup scripts or things of that nature. 

Secondarily Id also like to create cli tool which can provide monitoring and dynamic configuration support on the host for containers which run p1oneer without the need for something like a docker exec ... 

