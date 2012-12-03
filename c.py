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

def Twosies(aa):
  """Iterate pairs of elements from list aa"""
  i = 0
  comma = ""
  while aa:
    x = aa[0]
    y = aa[1]
    yield i, comma, x, y
    aa = aa[2:]
    i += 1
    comma = ","


OpNames = {
  # Map ast node class names to our methods on Pobj.
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
  if p.slice.__class__ is ast.Slice:
    lo = p.slice.lower.Value() if p.slice.lower else "nil"
    hi = p.slice.upper.Value() if p.slice.upper else "nil"
    return "(%s).Xslice(%s, %s)" % (p.value.Value(), lo, hi)

  if p.slice.__class__ is ast.Index:
    ix = p.slice.value.Value()
    # TODO: ix could be a slice instance at runtime.
    return "(%s).Xindex(%s)" % (p.value.Value(), ix)

  raise Exception("VSubscript BAD: %s" % ast.dump(p))

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
    print '  return Pstr("")  // extra return should be Pnone'
    print '}  // end func %s' % p.name
    print ''
  finally:
    Frames = Frames[:-1]

@T
def TIf(p):
  print 'if (%s).Bool() {' % p.test.Value()
  DoBody(p.body)
  if p.orelse:
    raise "TODO" # TODO
  print '}'

@T
def TAssign(p):
  print 'var G_%s = %s' % (p.targets[0].id, p.value.Value())

@T
def TReturn(p):
  print 'return %s' % p.value.Value()

PRIM = dict(
    int='Pint', int32='Pint', int64='Pint',
    uint='Pint', uint32='Pint', uint64='Pint', uintptr='Pint',
    byte='Pint', rune='Pint',
    float='Pint', double='Pint',
    string='Pstr',
    bool='Pbool',
    )
DEPRIM = dict(
    int='Int64', int32='Int64', int64='Int64',
    uint='Int64', uint32='Int64', uint64='Int64', uintptr='Int64',
    byte='Int64', rune='Int64',
    float='Int64', double='Int64',
    string='String',
    bool='Bool',
    )
def IsPrimType(a):
  if a[0] == "'":
    return type(a) is str and PRIM.get(a[1:])
  else:
    return type(a) is str and PRIM.get(a)

def GoType(a):
  # Horrible Kludge: Hardwire these.
  # We must understand these names in the scope they are seen,
  # i.e. in the net/http package.
  if a == "'ResponseWriter":
    return "p_http.ResponseWriter"
  if a == "'Request":
    return "p_http.Request"
  if a == "'Handler":
    return "p_http.Handler"

  if type(a) is str:
    return a[1:] if a[0] == "'" else a
  elif type(a) is not list:
    raise "Bad GoType"
  # list:
  h = a[0]
  if h == "STAR":
    return "*%s" % GoType(a[1])
  if h == "SEL":
    return "p_%s.%s" % (GoType(a[1]), GoType(a[2]))
  if h == "FN":
    z = "func("  
    args = a[1]
    rets = a[2]
    for i, comma, name, typ in Twosies(args):
      z += "%s %s %s" % (comma, name, GoType(typ))
    z += ")("
    for i, comma, name, typ in Twosies(rets):
      z += "%s %s %s" % (comma, name, GoType(typ))
    z += ") "
    return z
  raise Exception("no case in GoType")

NextTmp = 100
def NewTmp(prefix):
  global NextTmp
  NextTmp += 1
  return "%s_%d" % (prefix, NextTmp)

LaterStuff = "// Later\n"
def EmitLater(s):
  global LaterStuff
  LaterStuff += "\n" + s + "\n"

class QuickValue(object):
  def __init__(self, x):
    self.x = x
  def Value(self):
    return self.x

def AdaptArgToFn(a, d):
  # Assuming it's a local function, taking Pobjs, returning Pobjs.
  # i.e. Goify a Python Fn Value.
  print "//## AdaptArgToFn %s %s" % (a, d)
  args = d[1]
  rets = d[2]
  if rets:
    raise Exception("rets not supported")

  tmp = NewTmp("fn")
  fn = "func %s(" % tmp

  for i,comma,_,t in Twosies(args):
    fn += "  %s arg%d %s  " % (comma, i, GoType(t))

  fn += ") {\n"
    
  print "//##dump %s" % ast.dump(a)
  fn += "  _ = %s (" % a.Value()

  for i,comma,_,t in Twosies(args):
    fn += "%s NewPgo(arg%d)" % (comma, i)

  fn += ")\n"
  fn += "}\n"

  EmitLater(fn)

  return tmp
  # return "NewPgo(%s)" % tmp

def AdaptArgToDecl(a, d):
  print "//# AdaptArgToDecl( %s , %s )" % (a, d)

  if a is None:
    return "nil"  # Kludge?
  if a.Value() == "G_None":
    return "nil"  # Kludge?

  if d == 'Pobj':
    return a.Value()
  if IsPrimType(d):
    d = d[1:] if d[0]=="'" else d
    return "%s((%s).%s())" % (d, a.Value(), DEPRIM[d])
  t = GoType(d)
  if type(d) is list and d[0] == "FN":
    return AdaptArgToFn(a, d)
  return "(/*A*/(%s).(%s)/*Z*/)" % (a.Value(), t)

def AdaptArgsToDecl(aa, fdecl):
  print "//# AdaptArgsToDecl( %s , %s )" % (aa, fdecl)
  dcls = fdecl[2] if fdecl[0] == "FUNC" else fdecl[3]
  z = []
  ellipsis = None
  d = "?AdaptArgsToDecl?"
  while aa:
    aname = dcls[0]
    d = dcls[1]
    dcls = dcls[2:]
    if type(d) is list and d[0] == "ELLIPSIS":
      ellipsis = d[1]
      break
    a = aa[0]
    aa = aa[1:]
    z.append(AdaptArgToDecl(a, d))

  if ellipsis:
    for a in aa:
      z.append(AdaptArgToDecl(a, ellipsis))
    
  print "//# AdaptArgsToDecl -> %s" % (z,)
  return z

@V
def VCall(p):
  # aa = ','.join([x.Value() for x in p.args]) if p.args else ''
  aa = p.args

  # Try to turn fval into package & func name.
  fval = p.func.Value()
  fvec = fval.split('.')
  if len(fvec) > 1:
    fpkg = '.'.join(fvec[:-1])
    # Omit the p_.  THis is ugly.
    if fpkg[:2] == 'p_':
      fpkg = fpkg[2:]
    fname = fvec[-1]
    print "//# fval %s fvec %s fpkg %s fname %s" % (fval, fvec, fpkg, fname)

    imp = Imports.get(fpkg)
    grok = Grok.get(imp.replace('.', '/'))
    if imp and grok:
      fdecl = grok.get(fname)
      if fdecl:
        print "//# fdecl %s" % (fdecl, )
    else:
      print "//# imp %s grok %s" % (imp, grok)
      print "//# " + repr(Grok.keys())
      print "//# " + repr(Imports)
  else:
    # Dumb assumption: local function has as many args
    # as given, and returns 1 result.
    print "//# fval %s fvec (len=1)" % (fval,)
    fdecl = ['FUNC', '?FUNC?', (2*len(aa)) * ['Pobj'], ['Pobj']]
    
  aa = AdaptArgsToDecl(aa, fdecl)

  aaa = ','.join(aa) if aa else ''
  return '( %s ( %s ))' % (fval, aaa)

@V
def VNum(p):
  return 'Pobj(Pint(%d))' % p.n  # TODO: Pfloat, Plong.

@V
def VStr(p):
  return "Pobj(Pstr(`%s`))" % (p.s)  # TODO: handle ` inside literal string.

@T
def TPrint(p):
  for x in p.values:
    print 'func init() {'
    v = x.Value()
    print '  print((%s).String())' % v
    print '}'
  if p.nl:
    print 'func init() { println() }'

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
  raise Exception("For now, only imports from go or go.X allowed")
  for x in p.names:
    targ = x.name
    alias = x.asname if x.asname else x.name
    Imports[alias] = targ
    print 'import p_%s "%s"' % (alias, '/'.join(targ.split('.')))


@T
def TImportFrom(p):
  for x in p.names:
    targ = '%s.%s' % (p.module, x.name)
    # Only support go.* for now, and strip the go.
    if targ[:3] == "go.":
      targ = targ[3:]
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
    print "<<<", line.strip()
    w = line.split()
    atat = w[0]
    path = w[1]
    data = ListParser(' '.join(w[2:])).ParseThing()
    print ">>>", atat, path, repr(data)
    print "==============================================="

Grok = {}
def LoadGrok(fname):
  f = open(fname)
  for line in f:
    w = line.split()
    atat = w[0]
    path = w[1]
    data = ListParser(' '.join(w[2:])).ParseThing()
    name = data[1]

    d = Grok.get(path)
    if not d:
      Grok[path] = d = {}
    d[name] = data

if __name__ == '__main__':
  LoadGrok("_grok.txt")
  Translate(sys.argv[1])
  print LaterStuff
