package main
import("ccproj/udp_handler"
"net"
"fmt"
"os");

func main() {
	ch := make(chan []string)
	/* go func(){
		for{
			s := <-ch
			s = s
		}
		
	}() */
	udp_address,error := net.ResolveUDPAddr("udp","10.0.1.20:9090")
		if error != nil {
			fmt.Println(error)
			os.Exit(1)
		}
	
	connection_, error := net.ListenUDP("udp", udp_address)
	go udp_handler.ListenUdp("","",connection_ ,ch)
	go udp_handler.ListenServer(ch,connection_)
	print("started server.")
	select {}
}
