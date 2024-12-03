	package main
	import("ccproj/udp_handler");

	func main() {
		ch := make(chan []string)
		go func(){
			for{
				s := <-ch
				s = s
			}
			
		}()
		udp_handler.ListenUdp("server","127.0.0.1:8008",nil ,ch )
	}
