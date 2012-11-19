package main
  type Any interface {}

  func CompareLtE(a Any, b Any) Any {
    return a.(int) <= b.(int)
  }

  func BinOpAdd(a Any, b Any) Any {
    return a.(int) + b.(int)
  }

  func BinOpSub(a Any, b Any) Any {
    return a.(int) - b.(int)
  }
