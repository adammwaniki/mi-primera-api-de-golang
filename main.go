package main

func main(){
	// Creating our server and passing in the address
	server := NewAPIServer(":8080")
	// Now we need to run the server
	server.Run()
}