novncgo: noVNC
	go build

clean:
	go clean

noVNC:
	# git clone https://github.com/kanaka/noVNC.git
	git clone https://github.com/novnc/noVNC.git

update: noVNC
	git pull
	cd noVNC && git pull

