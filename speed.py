import requests
import time
url1 = "http://127.0.0.1/1.html"
url2 = "http://127.0.0.1/2.html"
url3 = "http://127.0.0.1/3.html"

t1 = time.time()
for i in range(10):
    r = requests.get(url1, proxies={'http':'http://127.0.0.1:1234'})
t2 = time.time()
print("test1:",t2-t1)

t1 = time.time()
for i in range(10):
    r = requests.get(url1, proxies={'http':'http://127.0.0.1:1234'})
t2 = time.time()
print("test2:",t2-t1)

t1 = time.time()
for i in range(10):
    r = requests.get(url1, proxies={'http':'http://127.0.0.1:1234'})
t2 = time.time()
print("test3:",t2-t1)