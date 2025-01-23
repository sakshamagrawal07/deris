package server

import (
	"fmt"
	"log"
	"net"

	"github.com/sakshamagrawal07/deris/commands"
	"github.com/sakshamagrawal07/deris/config"
	"github.com/sakshamagrawal07/deris/data"
	"golang.org/x/sys/unix"
)

var con_clients int = 0

func respond(response string, c net.Conn) error {
	if _, err := c.Write([]byte(response)); err != nil {
		return err
	}
	return nil
}

func closeConnection(c net.Conn) {
	c.Close()
	con_clients -= 1
	log.Println("Client disconnected with address: ", c.RemoteAddr(), "concurrent clients: ", con_clients)
}

// func RunSyncTCPServer() {
// 	log.Println("Starting synchronous TCP server on ", config.Host, ":", config.Port)

// 	lsnr, err := net.Listen("tcp", config.Host+":"+strconv.Itoa(config.Port))

// 	if err != nil {
// 		panic(err)
// 	}

// 	data.InitData()

// 	for {
// 		c, err := lsnr.Accept()
// 		if err != nil {
// 			panic(err)
// 		}

// 		con_clients += 1
// 		log.Println("Client connected with address: ", c.RemoteAddr(), "concurrent clients: ", con_clients)

// 		for {
// 			c.Write([]byte(":> "))
// 			cmd, err := commands.ReadCommand(c)

// 			if err != nil {
// 				closeConnection(c)
// 				if err == io.EOF {
// 					break
// 				}
// 				log.Println("err", err)
// 			}
// 			log.Print("command received: ", cmd)
// 			response := commands.ExecuteCommand(cmd)
// 			if err = respond(response, c); err != nil {
// 				log.Println("err write:", err)
// 			}
// 		}
// 	}
// }

func StartServer(address string) {
	// Create a socket
	fd, err := unix.Socket(unix.AF_INET, unix.SOCK_STREAM, unix.IPPROTO_TCP)
	if err != nil {
		log.Fatalf("failed to create socket: %v", err)
	}

	// Bind the socket to the address
	sa := &unix.SockaddrInet4{Port: config.Port}
	copy(sa.Addr[:], net.ParseIP(config.Host).To4())

	err = unix.Bind(fd, sa)
	if err != nil {
		log.Fatalf("failed to bind socket: %v", err)
	}

	// Start listening for connections
	err = unix.Listen(fd, unix.SOMAXCONN)
	if err != nil {
		log.Fatalf("failed to listen on socket: %v", err)
	}

	fmt.Printf("Listening on %s...\n", address)

	// Set the socket to non-blocking mode
	err = unix.SetNonblock(fd, true)
	if err != nil {
		log.Fatalf("failed to set non-blocking mode: %v", err)
	}

	data.InitData()

	// Event loop
	clients := make(map[int]struct{})
	for {
		// Prepare the poll structure
		var fds []unix.PollFd
		fds = append(fds, unix.PollFd{Fd: int32(fd), Events: unix.POLLIN})

		for clientFd := range clients {
			unix.Write(clientFd, []byte(":> "))
			fds = append(fds, unix.PollFd{Fd: int32(clientFd), Events: unix.POLLIN})
		}

		// Poll for events
		n, err := unix.Poll(fds, -1)
		if err != nil {
			if err == unix.EINTR {
				// Interrupted by a signal, retry the poll
				continue
			}
			log.Fatalf("poll error: %v", err)
		}

		if n > 0 {
			// Handle new connections
			if fds[0].Revents&unix.POLLIN != 0 {
				clientFd, _, err := unix.Accept(fd)
				if err != nil {
					log.Printf("failed to accept connection: %v", err)
				} else {
					unix.SetNonblock(clientFd, true)
					clients[clientFd] = struct{}{}
					fmt.Printf("New connection: %d\n", clientFd)
				}
			}

			// Handle existing connections
			for i := 1; i < len(fds); i++ {
				if fds[i].Revents&unix.POLLIN != 0 {
					clientFd := int(fds[i].Fd)
					buf := make([]byte, 1024)
					n, err := unix.Read(clientFd, buf)
					if err != nil {
						log.Printf("read error: %v", err)
						unix.Close(clientFd)
						delete(clients, clientFd)
						continue
					}
					cmd := string(buf[:n])
					log.Printf("Received command from %d: %s", clientFd, cmd)
					response, err := commands.ExecuteCommand(cmd)
					if(err != nil){
						log.Println("Error executing command: ", err)
					}
					unix.Write(clientFd, []byte(response))
					if n <= 0 || response == "Bye\n" {
						// Client closed connection
						unix.Close(clientFd)
						delete(clients, clientFd)
						fmt.Printf("Connection closed: %d\n", clientFd)
					}

				}
			}
		}
	}
}