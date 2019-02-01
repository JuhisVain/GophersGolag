package main

import
(
	"github.com/rthornton128/goncurses"
	"log"
)

const WALL int8 = 61   // =
const EMPTY int8 = 45  // -
const BLOCK int8 = 35  // #
const LIVE_WEASEL int8 = 87 // W
const SLEEP_WEASEL int8 = 119 // w
const GOPHER int8 = 64      // @
const WORM int8 = 38 // &

type Playarea struct {
	width, height int
	area []int8
	gopher Gopher
	weasel_list *Weasel
}

type Gopher struct {
	x, y, score int
}

type Weasel struct {
	x, y int
	alive bool
	next *Weasel
}

func main() {
	stdscr, err := goncurses.Init()

	if err != nil {
		log.Fatal("init", err)
	}
	defer goncurses.End()

	goncurses.Raw(true)
	goncurses.Echo(false)
	goncurses.Cursor(0)
	stdscr.Keypad(true)

	area_w := 23
	area_h := 23
	field := init_field(area_w, area_h)
	form_level(field, 1)

	for {

		stdscr.MoveAddChar(0, 25, goncurses.Char(48+field.gopher.x))
		stdscr.MoveAddChar(1, 25, goncurses.Char(48+field.gopher.y))
		
		draw_field(stdscr, field)
		goncurses.Update()
		input := stdscr.GetChar()
		handle_input(stdscr, field, input)
		
		
	}

	//stdscr.GetChar() // wait for input	
}

func tile_at(field *Playarea, x, y int) *int8{
	return &field.area[x+field.width*y]
}

func handle_input(w *goncurses.Window, field *Playarea, input goncurses.Key) {
	switch goncurses.KeyString(input) {
	case "left":
		move_gopher(field, -1, 0)
	case "up":
		move_gopher(field, 0, -1)
	case "right":
		move_gopher(field, 1, 0)
	case "down":
		move_gopher(field, 0, 1)
	}
	w.MovePrint(2, 25, input)
}

func move_gopher(field *Playarea, dx, dy int) {
	switch *tile_at(field, field.gopher.x+dx, field.gopher.y+dy) {
	case EMPTY:
		field.gopher.x += dx
		field.gopher.y += dy
		return
	case WALL:
		return
	case WORM:
		field.gopher.x += dx
		field.gopher.y += dy
		*tile_at(field, field.gopher.x+dx, field.gopher.y+dy) = EMPTY
		field.gopher.score += 100
		return
	case LIVE_WEASEL: // suicide
		//todo: DEATH
		field.gopher.score += -100
		return
	default:
		if push_block( // what the f are these indent rules...
			field,field.gopher.x+dx,field.gopher.y+dy,dx, dy,
			*tile_at(field,field.gopher.x,field.gopher.y)) {
				//beautiful
				move_gopher(field, dx, dy)
				*tile_at(field, field.gopher.x+dx, field.gopher.y+dy) = EMPTY
			}
		
	}
	
}

func push_block(f *Playarea, x, y, dx, dy int, source_block_type int8) bool {
	switch *tile_at(f, x+dx, y+dy) {
	case EMPTY, WORM:
		*tile_at(f, x+dx, y+dy) = source_block_type
		return true
	case BLOCK:
		return push_block(f, x+dx, y+dy, dx, dy, BLOCK)
	case WALL:
		return false
	}
	return false // temp
	
}

func init_field(width, height int) *Playarea{
	field := &Playarea{
		width,
		height,
		make([]int8, width*height),
		Gopher{width/2, width/2, 0},
		&Weasel{2, 2, true, nil}}

	return field

}

func draw_field(w *goncurses.Window, field *Playarea) {
	// draw field:
	for i := 0; i < field.width; i++ {
		for j := 0; j < field.height; j++ {
			w.MoveAddChar(j,i,goncurses.Char(field.area[i + field.width * j]))
		}
	}
	// draw gopher:
	w.MoveAddChar(field.gopher.y, field.gopher.x, goncurses.Char(GOPHER))
	// draw  weasels:
	for wea := field.weasel_list; wea != nil; wea = wea.next {
		if wea.alive {
			w.MoveAddChar(wea.y, wea.x, goncurses.Char(LIVE_WEASEL))
		} else {
			w.MoveAddChar(wea.y, wea.x, goncurses.Char(SLEEP_WEASEL))
		}
	}
	
	w.Refresh()
}

func form_level(field *Playarea, level_index int) {

	// Fill with empties
	for y := 0; y < field.height; y++ {
		for x := 0; x < field.width; x++ {
			*tile_at(field,x,y) = EMPTY
		}
	}
	
	// Border walls
	for x := 0; x < field.width; x++ {*tile_at(field,x,0) = WALL}
	for x := 0; x < field.width; x++ {*tile_at(field,x,field.height-1) = WALL}
	for y := 0; y < field.height; y++ {*tile_at(field,0,y) = WALL}
	for y := 0; y < field.height; y++ {*tile_at(field,field.width-1,y) = WALL}

	if (level_index != 1) {return} // escape for now

	// Fill insides with a square of blocks
	for y := 3; y < field.height-3; y++ {
		for x := 3; x < field.width-3; x++ {
			*tile_at(field,x,y) = BLOCK
		}
	}
	*tile_at(field, field.width/2, field.height/2) = EMPTY
}
