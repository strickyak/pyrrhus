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

print '//', ast.dump(a)
print

for n in ast.walk(a):
  print '//', n, '::', Dir(n), '::', str(n), '::', vars(n)
  #n.more = 'foo'
  print

def TModule(p):
  for x in p.body:
    x.Trans()
ast.Module.Trans = TModule

def TAssign(p):
  print 'var %s = %s;' % (p.targets[0].id, p.value.Value())
ast.Assign.Trans = TAssign

def VNum(p):
  return str(p.n)
ast.Num.Value = VNum

def TPrint(p):
  for x in p.values:
    print 'func init() { print (%s); }' % x.Value()
  if p.nl:
    print 'func init() { println(); }'
ast.Print.Trans = TPrint


def VBinOp(p):
  return "((%s) %s (%s))" % (p.left.Value(), '+', p.right.Value())
ast.BinOp.Value = VBinOp

def VName(p):
  return p.id
ast.Name.Value = VName
  
print 'package main'
a.Trans()
print 'func main() { println("OK"); }'
