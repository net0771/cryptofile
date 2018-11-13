from hashlib import sha256
from uuid import uuid4

chain = {}
chain["1"] = {"name":'test1'}
chain["2"] = {"name":'test2'}

print(len(chain))
print(chain.get('1'))

