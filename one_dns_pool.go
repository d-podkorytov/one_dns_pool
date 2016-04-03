// Speed 47459 qps
//       57040 qps on single CPU core by taskset -c 3 ./cdns_ch 
package main

import (
	"flag"
	"fmt"
	"net"
	"log"
)

var port = flag.Int("p", 53, "The listening port ")
var udpPackageBufferSize = flag.Int("l", 1024, "The size of the udp package buffer")



type ReadResult_t struct {
buf  []byte
addr *net.UDPAddr
size int
}

func main() {
	flag.Parse()
	// open UDP socket
	log.Println("try to open UDP socket, port ", port)
	socket, err := net.ListenUDP("udp4", &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: *port,
	})

	if err != nil { fmt.Println("listen failed!", err)
		        return }

        log.Println("listen ", socket,port,err)

	defer socket.Close()

        var ch [256]chan ReadResult_t
        for i := range ch { ch[i] = make(chan ReadResult_t,1024) }

        fmt.Println("created chanalls:",len(ch))
         
        for i := 0; i < 256; i++ { go worker_loop( socket , ch ,i) }

        var id0 byte

// endless loop wait clients asks and pass it to channal array
	for {
		ask_data := make([]byte, *udpPackageBufferSize)
		readn, remoteAddr, _ := socket.ReadFromUDP(ask_data)
		
                //readn, remoteAddr, err := socket.ReadFromUDP(ask_data)

                //for fast working dont control returned error
		//if err != nil { fmt.Println("recvfrom error!", err)
		//	        continue }

                // get Message ID low byte and use it as address for worker chanall in pool
                id0=ask_data[0]
		
		//go  process(socket, ask_data[:readn], id, remoteAddr)
                ret:= ReadResult_t{ ask_data, remoteAddr ,readn }
                
                // use Msg_ID[0] as address for channal id0 and pass data
                ch[id0] <- ret
	}
}

func worker_loop( conn *net.UDPConn, ch [256]chan ReadResult_t ,id int) {
 for 
 { Ask:= <- ch[id]
  // fmt.Println("got from ch[id] Ask=",Ask)
   process(conn, Ask.buf , Ask.addr )    
 }
}


func process(conn *net.UDPConn, ask_data []byte, remoteAddr *net.UDPAddr) {
// here do reply on ask data

//    1000.dip A DNS record for return from reply
rr_record := []byte{0,0,1,0,0,1,0,1,0,0,0,0,4,49,48,48,48,3,100,105,112,0,0,1,0,1,192,12,0,1,0,1,0,0,8,212,0,4,87,118,90,81}

                id:=make([]byte,2)

                // get Message ID from ask
                id[0]=ask_data[0]
                id[1]=ask_data[1]

                // change message ID according to ask
                rr_record[0]=id[0]
                rr_record[1]=id[1]

        //dont handle errors for speed working  
	conn.WriteToUDP(rr_record, remoteAddr)

//if err != nil { fmt.Println("send data error", err) }
}
