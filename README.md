This program extracts data from the Covid tracking project and plots the data
on the screen using a browser

It uses the Lexer2 (github.com/Saied74/Lexer2) for parsing JSON from the covid
tracking project.  

The program is pretty simple and well commented (I think), so I won't belabor
the points anymore.

The environment variable is GRAPHPATH and it must be set to the head of the
project which in my case is $GOPATH/src/graphs.  In the Dockerfile it is set
to /go/src/graphs.

The two flags are covidProjectURL and ipAddress.  covidProjectURL overrides
the URL for getting the covid porject data.  The ipAddress overrides the server
IP address that the web server listens to.
