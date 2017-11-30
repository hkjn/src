# Sample of converting private key to addr
from pycoin import key

keys= ["KyQmCtGbpUsApxvxRRxTQjwzYpoqANTkbRFPjAY57Bbv1YmwBU45",]

for akey in keys:
    key_obj = key.Key.from_text(akey) 
    print key_obj.address()
