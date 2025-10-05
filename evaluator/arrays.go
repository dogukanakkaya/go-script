package evaluator

func GetArrayProperty(arr Array, property string) Value {
	switch property {
	case "length":
		return float64(len(arr))
	default:
		return nil
	}
}
