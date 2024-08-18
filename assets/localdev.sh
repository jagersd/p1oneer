#!/usr/bin/env bash

containerName="p1oneer-dev"

check_image () {
    dockerImage=$(docker images | grep -r "\b $containerName \b")
    if [ -n "$dockerImage" ]; then
        return 1
    else
        return 0
    fi
}

check_image
if [ $? -eq 0 ]; then
    cd "$HOME"/projects/p1oneer && docker build -f assets/test-dockerfile -t "$containerName" .
fi

digIn () {
    docker exec -it "$containerName" /bin/bash
    exit 0
}

start_container () {
    activeContainer=$(docker ps | grep -r "\b $containerName \b")
    if [ -n "$activeContainer" ]; then
        echo "$activeContainer"
        echo "Getting into active container"
        sleep 1
        digIn
    else
        inactiveContainer=$(docker ps -a | grep "$containerName")
        if [ -n "$inactiveContainer" ]; then
            echo "Activating existing container...."
            docker start "$containerName"
            sleep 5
            digIn
        else
            echo "Spinning up a fresh container"
            docker run --rm --name "$containerName" -it --mount src="$HOME/projects/p1oneer",target=/app,type=bind "$containerName" /bin/bash && digIn
        fi
    fi
}

start_container
