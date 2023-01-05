# tickets_12306

TODO：
功能：查找可买的票，（上车补票；多买几站）
from_station, to_station, expected_trainos
1. 如果 from_station, to_station 和期望的trainnos . 能够找到可购买的，就返回
2. 如果没有：trainos 遍历to_station的上一站。看看能不能买，如果没有买的，再往上。直到找到可买的
3. 如果上一站没有，然后遍历下一站。直到遍历可买的。


TODO：

中转：
