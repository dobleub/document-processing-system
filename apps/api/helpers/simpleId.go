package helpers

import "math/rand"

type SimpleID string

func (s SimpleID) String() string {
	return string(s)
}

func (s SimpleID) IsValid() bool {
	return true
}

func (s SimpleID) FromString(token string) SimpleID {
	return SimpleID(token)
}

/*
 * Generate is a functoin which creates SimpleID
 * @return SimpleID
 * Example:
  * 	simpleID := new(SimpleID).Generate()
  * 	fmt.Println(simpleID)
  * 	// Output: XHOC0-45685-3JH9A-4KJH6
*/
func (s SimpleID) Generate() SimpleID {
	chunks := 4
	lengthPerChunck := 5
	abc := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	simpleID := ""

	for i := 0; i < chunks; i++ {
		for j := 0; j < lengthPerChunck; j++ {
			simpleID += string(abc[rand.Intn(len(abc))])
		}
		if i < chunks-1 {
			simpleID += "-"
		}
	}

	return s.FromString(simpleID)
}
