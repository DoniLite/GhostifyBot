package types

type EventHandler = func() 


func restE(fn EventHandler, actions ...string) {

}

func main() {
	restE(func () {})
}