package utils

import "fmt"

func PrintBannerOps() {
	fmt.Println(`
 ______     ______   ______                   
/\  __ \   /\  == \ /\  ___\                  
\ \ \/\ \  \ \  _-/ \ \___  \                 
 \ \_____\  \ \_\    \/\_____\                
  \/_____/   \/_/     \/_____/                
   
`)
}

func PrintBannerPipeline() {
	fmt.Println(`
 ______   __     ______   ______     __         __     __   __     ______    
/\  == \ /\ \   /\  == \ /\  ___\   /\ \       /\ \   /\ "-.\ \   /\  ___\   
\ \  _-/ \ \ \  \ \  _-/ \ \  __\   \ \ \____  \ \ \  \ \ \-.  \  \ \  __\   
 \ \_\    \ \_\  \ \_\    \ \_____\  \ \_____\  \ \_\  \ \_\\"\_\  \ \_____\ 
  \/_/     \/_/   \/_/     \/_____/   \/_____/   \/_/   \/_/ \/_/   \/_____/ 
   
`)
}

func PrintUsage() {
	fmt.Println(
		`Usage: fx-ops 
   [ -l | --local ]
   [ -i | --install ]
   <phase: version|sops|test|build|install|deploy> [project] [args...]

Examples:
   fx-ops -l build fx-pos
   fx-ops deploy fx-api
`)
}
