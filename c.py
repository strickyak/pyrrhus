# Pyrrhus of Epirus
# Copyright 2012 Strick

import ast
import re

def NodeStringer(n):
  return n.__class__.__name__ + '~' + str(n.nom)
  return n.__class__.__name__ + '~' + str(hash(n) % 900 + 100)

def Dir(n):
  return ','.join([x for x in dir(n) if not re.match('^__', x)])

a = ast.parse(open("t1.py").read())
i = 1
for n in ast.walk(a):
  n.nom = i
  n.__class__.__str__ = NodeStringer
  n.__class__.__repr__ = NodeStringer
  i += 1

print ast.dump(a)
print

for n in ast.walk(a):
  print n, '::', Dir(n), '::', str(n), '::', vars(n)
  #n.more = 'foo'
  print
