# Pyrrhus of Epirus
# Copyright 2012 Strick

import ast
import re
import sys

def NodeStringer(n):
  return n.__class__.__name__ + '~' + str(n.nom)

# def Dir(n):
#   return ','.join([x for x in dir(n) if not re.match('^__', x)])

class T(object):
  def __init__(self, f):
    self.f = f
    print '////@ f=', f
    name = self.f.__name__[1:]
    print '////@ name=', name
    cls = vars(ast)[name]
    print '////@ cls=', cls
    cls.Trans = self.f
  def __call__(self): pass

class V(object):
  def __init__(self, f):
    self.f = f
    print '////@ f=', f
    name = self.f.__name__[1:]
    print '////@ name=', name
    cls = vars(ast)[name]
    print '////@ cls=', cls
    cls.Value = self.f
  def __call__(self): pass

@T
def TModule(p):
  for x in p.body:
    x.Trans()
# ast.Module.Trans = TModule

@T
def TFunctionDef(p):
  args_str = ','.join([x.id for x in p.args.args])
  print 'func %s(%s) Any {' % (p.name, args_str)
  for x in p.body:
    x.Trans()
  print '}  // func %s' % (p.name, )

@T
def TAssign(p):
  print 'var %s = %s;' % (p.targets[0].id, p.value.Value())

def VNum(p):
  return str(p.n)
ast.Num.Value = VNum

@T
def TPrint(p):
  for x in p.values:
    print 'func init() { print (%s); }' % x.Value()
  if p.nl:
    print 'func init() { println(); }'
ast.Print.Trans = TPrint

@V
def VBinOp(p):
  return "((%s) %s (%s))" % (p.left.Value(), '+', p.right.Value())

@V
def VName(p):
  return p.id
  
def Translate(filename):
  a = ast.parse(open(filename).read())
  i = 1

  # Monkeypatch the str() and repr() of all used AST node classes.
  for n in ast.walk(a):
    n.nom = i
    n.__class__.__str__ = NodeStringer
    n.__class__.__repr__ = NodeStringer
    i += 1

  # Walk the AST and print the nodes, for reverse-engineering.
  for n in ast.walk(a):
    print '//', n, '//', vars(n), '//'
    #n.more = 'foo'
    print

  print 'package main'
  a.Trans()
  print 'func main() { println("OK"); }'

Translate(sys.argv[1])
