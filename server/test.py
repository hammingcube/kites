import requests
s = requests.session()
data = {u'Files': {u'aa': {u'Content': u'this',
   u'Language': u'j',
   u'Name': u'k',
   u'Truncated': False}},
 u'Id': u'a',
 u'UserId': u'default_user'}

resp = s.get('http://localhost:8080/gists/a')
print(resp.content)
resp = s.post('http://localhost:8080/gists', json=data)
print(resp.content)
resp = s.get('http://localhost:8080/gists/a')
print(resp.content)