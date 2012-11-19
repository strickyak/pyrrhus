package main
  import "fmt"

  type Any interface {
    func PrintableValue(a Any) string {
      switch t := a.(type)
      default:
        panic("Type not handled in PrintableValue")
      case int:
        return fmt.SPrintf("%d", t)
      case string:
        return t
    }
  }

  func CompareLtE(a Any, b Any) Any {
    return a.(int) <= b.(int)
  }

  func BinOpAdd(a Any, b Any) Any {
    return a.(int) + b.(int)
  }

  func BinOpSub(a Any, b Any) Any {
    return a.(int) - b.(int)
  }
