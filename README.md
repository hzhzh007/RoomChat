# RoomChat

high concurrency and scalability room chat backend based on websocket and write in golang

# fitures
1. use websocket as the main communication method (high concurrency connections and low delay)
2. etcd as the service discovery (on plan)

## architecture

pass


## direction

 * /common 
 * /connector
     the front connector server
 * /gate
     return the connector's ip before the client connect to the chat room
 * /router
     like the controller in mvc, just accept the msg ,deal the logic and decide whether broad it or not
 * /request
    receive send msg request
 * /tools 



