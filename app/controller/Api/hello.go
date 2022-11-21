package Api

import (
	"aiyun_local_srv/library/response"
	"fmt"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gookit/color"
)

var Hello = hello{}

type hello struct{}

// Index is a demonstration route handler for output "Hello World!".
func (*hello) Index(r *ghttp.Request) {

	// quick use package func
	color.Redp("Simple to use color")
	color.Redln("Simple to use color")
	color.Greenp("Simple to use color\n")
	color.Cyanln("Simple to use color")
	color.Yellowln("Simple to use color")

	// quick use like fmt.Print*

	color.Red.Println("Simple to use color")
	color.Green.Print("Simple to use color\n")
	color.Cyan.Printf("Simple to use %s\n", "color")
	color.Yellow.Printf("Simple to use %s\n", "color")

	// use like func
	red := color.FgRed.Render
	green := color.FgGreen.Render
	fmt.Printf("%s line %s library\n", red("Command"), green("color"))

	// custom color
	color.New(color.FgWhite, color.BgBlack).Println("custom color style")

	// can also:
	color.Style{color.FgCyan, color.OpBold}.Println("custom color style")

	// internal theme/style:
	color.Info.Tips("message")
	color.Info.Prompt("message")
	color.Info.Println("message")
	color.Warn.Println("message")
	color.Error.Println("message")

	// use style tag
	color.Print("<suc>he</><comment>llo</>, <cyan>wel</><red>come</>\n")
	// Custom label attr: Supports the use of 16 color names, 256 color values, rgb color values and hex color values
	color.Println("<fg=11aa23>he</><bg=120,35,156>llo</>, <fg=167;bg=232>wel</><fg=red>come</>")

	color.Println("<fg=6600CC>11111111111111111</>")
	color.Println("<fg=FF3366>2222222222222222222</>")
	color.Println("<fg=FF0000>33333333333333333333333</>")
	color.Println("<fg=FF33FF>4444444444444444444444</>")
	color.Println("<fg=FFFF00>55555555555555555</>")
	color.Println("<fg=FF0066>66666666666666666</>")
	color.Println("<fg=CC66FF>77777777777777777777</>")
	color.Println("<fg=CCFF33>8888888888888888888</>")
	color.Println("<fg=9900CC>999999999999</>")
	color.Println("<fg=66FF99>12312321</>")

	// apply a style tag
	color.Tag("info").Println("info style text")

	// prompt message
	color.Info.Prompt("prompt style message")
	color.Warn.Prompt("prompt style message")

	// tips message
	color.Info.Tips("tips style message")
	color.Warn.Tips("tips style message")

	response.Success(r, "1111111111")
}
