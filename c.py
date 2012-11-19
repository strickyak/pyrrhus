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
  "Decorate a Patched Trans Method"
  def __init__(self, f):
    name = f.__name__[1:]
    cls = vars(ast)[name]
    cls.Trans = f
  def __call__(self): pass

class V(object):
  "Decorate a Patched Value Method"
  def __init__(self, f):
    name = f.__name__[1:]
    cls = vars(ast)[name]
    cls.Value = f
  def __call__(self): pass

@V # REPLACE RYAN
def VCompare(p):
  return "0 /*ryan*/"
@V # REPLACE RYAN
def VBinOp(p):
  return "((%s) %s (%s) /*ryan*/)" % (p.left.Value(), '+', p.right.Value())


def DoBody(body):
  for x in body:
    x.Trans()

@T
def TModule(p):
  DoBody(p.body)

@T
def TFunctionDef(p):
  args_str = ','.join([x.id for x in p.args.args])
  print 'func %s(%s) Any {' % (p.name, args_str)
  DoBody(p.body)
  print '}  // func %s' % (p.name, )

@T
def TIf(p):
  print 'if %s {' % p.test.Value()
  DoBody(p.body)
  if p.orelse:
    pass # TODO
  print '}'

@T
def TAssign(p):
  print 'var %s = %s' % (p.targets[0].id, p.value.Value())

@T
def TReturn(p):
  print 'return %s' % p.value.Value()

@V
def VCall(p):
  aa = ','.join([x.Value() for x in p.args]) if p.args else ''
  return '( %s ( %s ))' % (p.func, aa)

@V
def VNum(p):
  return str(p.n)

@T
def TPrint(p):
  for x in p.values:
    print 'func init() { print (%s); }' % x.Value()
  if p.nl:
    print 'func init() { println(); }'
ast.Print.Trans = TPrint

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
