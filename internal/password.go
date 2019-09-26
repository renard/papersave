package papersave
  
import (  
	"math/rand"
	"time"
)

func GenPassword(size int) string {
	LowerLetters := "abcdefghijklmnopqrstuvwxyz"
	UpperLetters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Digits := "0123456789"
	//	Symbols := "~!@#$%^&*()_+`-={}|[]\\:\"<>?,./"
	Symbols := "~!@#$%&()+-={}[]\\:<>?,./"

	source := LowerLetters + UpperLetters + Digits + Symbols
	lsource := len(source)
	ret := ""

	src := rand.NewSource(time.Now().UnixNano())
    r := rand.New(src)
	
	for i := 0; i < size ; i++ {
		v := r.Intn(lsource-1)
		c := source[v]
		ret += string(c)
	}
	return ret
	
}
