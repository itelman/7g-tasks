# TEST 1:

# terminal 1
./task2/analyzer 127.0.0.1:9000
# terminal 2
sudo ./task2/sniffer lo 127.0.0.1 127.0.0.1:9000
# terminal 3
ping 127.0.0.1

# TEST 2:

# terminal 1
./task2/analyzer 127.0.0.1:9000
# terminal 2
sudo ./task2/sniffer wlp0s20f3 10.124.42.58 127.0.0.1:9000
# terminal 3
ping 8.8.8.8