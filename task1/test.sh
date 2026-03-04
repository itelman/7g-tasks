# Test for LOW/NO traffic diversity:

# terminal 1
sudo ./task1/task1 lo 127.0.0.1 1000 3
# terminal 2
ping 127.0.0.1

# Test for HIGHER traffic diversity:

# terminal 1
sudo ./task1/task1 wlp0s20f3 10.124.42.58 1000 3
# terminal 2...n
for site in google.com github.com wikipedia.org stackoverflow.com example.com; do
    curl -s https://$site >/dev/null &
done

# Ping multiple servers
for i in {1..10}; do ping -c 3 8.8.8.$i & done

ping -c 5 8.8.8.8 &
ping -c 5 1.1.1.1 &
ping -c 5 8.8.4.4 &
ping -c 5 208.67.222.222 &
ping -c 5 9.9.9.9 &