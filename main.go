package main
import(
	"os"
	"github.com/russross/blackfriday"
	"github.com/fatih/color"
	"net/http"
	"net"
	"fmt"
	"strings"
	"time"
)
func getDir()string{
	dir, err := os.Getwd()
	if err != nil {panic(err)}
	return dir
}
func getHomeDir()string{
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return home
}
func GetfileTime(path string)time.Time{
	file ,err:= os.Stat(path)
	if err != nil {
		panic(err)
	}
	time := file.ModTime()
	return time
}
func ReadFile(path string)string{
	bytes, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			color.Red("× File not found.Please check file name.")
		} else {
			fmt.Println("Error:", err)
		}
		os.Exit(0)
	}
	text := string(bytes)
	return text;
}
func WriteFile(path string ,text string){
	bytes := []byte(text)
	if err := os.WriteFile(path, bytes, 0666); err != nil {
		panic(err)
	}
}
func getHTML(md []byte)string{
	html := string(blackfriday.MarkdownCommon(md))
	return html
}
func getMain(path string){
	home := getHomeDir()
	mdtext := ReadFile(path)
	mdhtml := getHTML([]byte(mdtext))
	template := ReadFile(home + "/qread/template.html")
    html := strings.Replace(template, "{{ . }}", mdhtml, 1)
    WriteFile(home+"/qread/index.html",html)
}
func Update(path string){
	getMain(path)
}
func OpenServer(path string){
	var status string = "ok"
	go func(){
	home := getHomeDir()
	server := home + "/qread"
	fs := http.FileServer(http.Dir(server))
	http.Handle("/", fs)
	http.HandleFunc("/api/status", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, status)
	})
    err := http.ListenAndServe(":8000", nil)
	if err != nil {
		if err == http.ErrServerClosed {
			color.Red("× Server closed:", err)
		} else if opErr, ok := err.(*net.OpError); ok && opErr.Err.Error() == "bind: address already in use" {
			color.Red("× Port already in use.Please kill the program.")
		} else {
			color.Red("× Error:", err)
		}
		os.Exit(0)
	}
	}()
	go func() {
		link := color.GreenString("http://localhost:8000/")
		fmt.Println("You can see your markdown on "+ link)
		var lasttime time.Time = GetfileTime(path)
			for {
				time.Sleep(1 * time.Second)
				if lasttime == GetfileTime(path){
					fmt.Printf("\r We are waiting for you visit our html...")
				}else{
					fmt.Println("")
					color.Yellow("! File had updated.")
					Update(path)
					status = "update"
					lasttime = GetfileTime(path)
					color.Green("○ HTML had updated")
					time.Sleep(1 * time.Second)
					status = "ok"
				}
			}
		}()

		select {}
}
func Show(){
	dir := getDir()
	path := dir +"/" +os.Args[1]
	getMain(path)
	OpenServer(path)
}
func AAprint(){
	aa := `  ___          _        _      ____                    _ 
 / _ \  _   _ (_)  ___ | | __ |  _ \   ___   __ _   __| |
| | | || | | || | / __|| |/ / | |_) | / _ \ / _' | / _' |
| |_| || |_| || || (__ |   <  |  _ < |  __/| (_| || (_| |
 \__\_\ \__,_||_| \___||_|\_\ |_| \_\ \___| \__,_| \__,_|`
	fmt.Println(aa)
}
func Start(){
	home := getHomeDir()
	if _, err := os.Stat(home + "/qread"); err != nil {
		if err := os.Mkdir(home+"/qread", 0777); err != nil {
			fmt.Println(err)
		}	
	}
}
func main(){
	args := os.Args
	time.Sleep(1 * time.Second)
	AAprint()
	if len(args) == 1 {
		rogo := color.GreenString("QucikRead")
		help := color.New(color.FgBlack, color.BgWhite)
		fmt.Printf("Welcome to "+ rogo +"! \n If you want more about "+ rogo +",please type " )
		help.Print("qread -h")
		link := color.HiBlueString("https://github.com/roistaff/qread")
		fmt.Println(" or visit " + link)
		os.Exit(0)
	}else if args[1] == "-h"{
		
}else{
	Show()
}
}