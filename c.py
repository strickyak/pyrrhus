# Pyrrhus of Epirus
# Copyright 2012 Strick

import ast
import re
import sys

Imports = {}
Frames = []

def NodeStringer(n):
  return n.__class__.__name__ + '~' + str(n.nom)

class T(object):
  "Decorate a Patched Trans Method"
  def __init__(self, f):
    name = f.__name__[1:]
    cls = vars(ast)[name]
    cls.Trans = f
  def __call__(self):
    pass

class V(object):
  "Decorate a Patched Value Method"
  def __init__(self, f):
    name = f.__name__[1:]
    cls = vars(ast)[name]
    cls.Value = f
  def __call__(self):
    pass

OpNames = {

  'Add': 'Xadd',
  'Sub': 'Xsub',
  'Mult': 'Xmul',
  'Div': 'Xdiv',
  'Mod': 'Xmod',
  
  'Gt': 'Xgt',
  'GtE': 'Xge',
  'Lt': 'Xlt',
  'LtE': 'Xle',
  'Eq': 'Xeq',
  'NotEq': 'Xne',
}

@V
def VBinOp(p):
  return "(%s).%s(%s)" % (p.left.Value(), OpNames[p.op.__class__.__name__], p.right.Value())

@V
def VCompare(p):
  # TODO:  double compare ops.
  if len(p.ops) != 1: raise Exception('only simple compare supported: ' + str(p))
  return "Pobj((%s).%s(%s))" % (p.left.Value(), OpNames[p.ops[0].__class__.__name__], p.comparators[0].Value())

@V
def VSubscript(p):
  return "%s[%s]" % (p.value.Value(), p.slice.Value())

@V
def VSlice(p):
  lower = p.lower.Value() if p.lower is not None else ""
  upper = p.upper.Value() if p.upper is not None else ""
  # TODO: negative slices.
  return "%s:%s" % (lower, upper)

def DoBody(body):
  print "//--- body len is %d, body=%s" % (len(body), body)
  i = 0
  for x in body:
    print "//-------- DOING ", i, ';', ast.dump(x)
    x.Trans()
    i += 1

@T
def TModule(p):
  DoBody(p.body)

@T
def TFunctionDef(p):
  global Frames
  Frames.append([x.id for x in p.args.args])
  try:
    args_str = ','.join(['v_%s Pobj' % x.id for x in p.args.args])
    print ''
    print 'func G_%s(%s) Pobj {' % (p.name, args_str)
    DoBody(p.body)
    print '}  // end func %s' % p.name
    print ''
  finally:
    Frames = Frames[:-1]

@T
def TIf(p):
  print 'if (%s).Bool() {' % p.test.Value()
  DoBody(p.body)
  if p.orelse:
    pass # TODO
  print '}'

@T
def TAssign(p):
  print 'var G_%s = %s' % (p.targets[0].id, p.value.Value())

@T
def TReturn(p):
  print 'return %s' % p.value.Value()

@V
def VCall(p):
  aa = ','.join([x.Value() for x in p.args]) if p.args else ''
  return '( %s ( %s ))' % (p.func.Value(), aa)

@V
def VNum(p):
  return 'Pint(%d)' % p.n  # TODO: Pfloat, Plong.

@V
def VStr(p):
  return "Pstr(`%s`)" % (p.s)  # TODO: handle ` inside literal string.

@T
def TPrint(p):
  for x in p.values:
    print 'func init() { print((%s).String()); }' % x.Value()
  if p.nl:
    print 'func init() { println(); }'

@T
def TExpr(p):
  print p.value.Value()

@V
def VAttribute(p):
  if p.value.__class__ is ast.Name:
    # Cannot have parens around import symbol.
    return '%s.%s' % (p.value.Value(), p.attr)
  else:
    return '(%s).%s' % (p.value.Value(), p.attr)

@V
def VIndex(p):
  return p.value.Value()

@V
def VName(p):
  id = p.id
  for f in Frames:
    for v in f:
      if id == v:
        return 'v_' + id
  if Imports.get(id):
    return 'p_' + id
  return 'G_' + id

@T
def TImport(p):
  for x in p.names:
    targ = x.name
    alias = x.asname if x.asname else x.name
    Imports[alias] = targ
    print 'import p_%s "%s"' % (alias, '/'.join(targ.split('.')))


@T
def TImportFrom(p):
  for x in p.names:
    targ = '%s.%s' % (p.module, x.name)
    alias = x.asname if x.asname else x.name
    Imports[alias] = targ
    print 'import p_%s "%s"' % (alias, '/'.join(targ.split('.')))


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

class ListParser(object):
  def __init__(self, s):
    self.ww = s.split()
    self.p = 0

  def ParseThing(self):
    if self.ww[self.p] == "{":
      self.p += 1
      v = []
      while self.ww[self.p] != "}":
          v.append(self.ParseThing())
      self.p += 1
      # Commas, colons, and any word ending in colon are OMITTED.
      return [x for x in v if (x[-1] != ':' and x != "," if type(x) is str else True)]
    else:
      z = self.ww[self.p]
      self.p += 1
      return z

def Demo(fname):
  f = open(fname)
  for line in f:
    line = line.strip()[3:]  # Drop initial "@@ "
    print "<<<", line
    print ">>>", repr(ListParser(line).ParseThing())
    print "==============================================="

if __name__ == '__main__':
  # Translate(sys.argv[1])
  Demo(sys.argv[1])
