# folders
SERVSRC=server/

# base filename for go
SERVER=websock
# build file for go
GOCC=go build
# opens a new terminal, -e executes command in quotes that follow
TERMINAL=gnome-terminal -e 

# firefox/chrome new windows
FIREFOX=firefox -new-window
CHROME=google-chrome -new-window

# localhost http address
LHOST=http://localhost:8787

# chat is the default target
CHAT=chat.html
TARGET=$(CHAT)

# echo client
ECLIENT=echo.html

# test html client
TESTHTML=test_client.html

# test websock client
TCLIENT="test_client.html"

# open the server first, then the chat
run: serve bots chat

# open chat
chat: open-page

# opens a new webpage to target
open-page: 
	sleep 1
	# copy the html files to the same location
	# cp $(SRC)*.html $(BIN)
	$(FIREFOX) "$(LHOST)/$(TARGET)"
	#$(CHROME) "$(LHOST)/$(TARGET)"


# spawn a bunch of bots to populate the room
bots:
	@echo "Not implemented yet. Bots to come"

# opens an echo client
echo: TARGET=$(ECLIENT)
echo: run

# open the test client
test: TARGET=$(TCLIENT)
test: run

# server scripts
# opens server in a new terminal
serve: build-server
	$(TERMINAL) "bash -c ./$(SERVER).exe;bash"

# build the files and move them to bin
build-server: clean-server
	$(GOCC) -o $(SERVER).exe $(SERVSRC)$(SERVER).go

# make a bin folder (ignore errors) and remove run file (ignore if not there)
# clear twice to clean screen
clean-server:
	rm -f $(SERVER).exe
	clear

# installs golang and murciral to user the websocket
try: install serve test

install: go-hg resolve-path

go-hg:
	# install golang and mercurial
	sudo apt-get install golang mercurial
	mkdir -p $$HOME/go

resolve-path:
	# append path updates to your bashrc
	sudo echo "export GOPATH=\$$HOME/go" >> ~/.bashrc
	sudo echo "export PATH=\$$PATH:\$$GOPATH/bin" >> ~/.bashrc
	# demand notification
	@echo "\033[7m\033[1myou must reopen terminal then run 'make get-websock'\033[0m"
	@echo "hit enter..."
	@read foo;

get-websock:
	# get websock package
	go get code.google.com/p/go.net/websocket

# prints a bit of help
help:
	clear
	@echo "try 'make install' to install golang and then follow instructions to install websocket so you can compile the websocket"
	@echo "'make'          -- starts everything (executes first target)"
	@echo "'make echo'     -- opens up an echo page (stolen from internet)"
	@echo "'make server'   -- opens server in a new terminal"
