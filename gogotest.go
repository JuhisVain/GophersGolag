package main

import
(
	"fmt"
)

func Check_weasel_list() {
	weasel_list := &Weasel{1,1,true,&Weasel{2,2,true,nil}}

	if check_weasels(weasel_list,
		func(weasel *Weasel)bool{
			fmt.Printf("x %d, y %d", weasel.x, weasel.y)
			return false
		}) {
			fmt.Println("Finished!")
		}
}
