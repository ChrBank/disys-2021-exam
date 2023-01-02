# disys-2021-exam
run server/server.go
Then follow by a 1, 2, or 3.

then run client/client.go
then give a name

## What system model are you assuming in your implementation. Write a full desciption.
    I am assuming a system model with only crash failures.
    I won't assume that a server or client can recover after crashes, a new one have to be started instead.

## What is the minmal number of nodes in your system to fullfill the requirements? Why?
    In my system you need three servers because a client only knows three servers that it can connects to.
## Explain how your system recovers from crash failure.
    If a server or client crashes it does not recover after it has crashed. If the server crashes the client gets notified that the server has crashed. 
    If the client crashes noting happends.
## Explain how you achieve the Monotonicity requirement.
    I use locks to enter the critical section, when a client enters the critical section, only then they can increment their amount. Each time it gets printed in the terminal, the amount will increment after, then I use a time.sleep to make sure it does not enter the critical section right after. 
## Explain how you achieve the Liveliness requirement.
    Each time a client tries to increment their value they, need to get the critical section, then the critical section ends when the client has incremented their value since I use the "defer" keyword in the SendBid function in server/server.go, so the critical section ends when someting is returned. 