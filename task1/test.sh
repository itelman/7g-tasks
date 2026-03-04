# Test for LOW/NO traffic diversity:

# 1)
# sudo ./task1/task1 lo 127.0.0.1 100 3

# 2)
# ping 127.0.0.1

# Test for HIGHER traffic diversity:

# 1)
sudo ./task1/task1 wlp0s20f3 10.124.42.58 100 3

# 2)
# ping 8.8.8.8
# curl google.com
# curl github.com
# curl wikipedia.org