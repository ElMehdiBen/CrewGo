# CrewGo

This repositery is an answer to this exercice : https://github.com/crewdotwork/backend-challenge

Testing the script is easy :

## Option 1 : using Dockefile

We can start by running the following command (same folder as the Dockefile): 

    docker build . -t crewgo

    docker run -p 1112:1112 -d crewgo

## Option 2 : using Docker Compose File

We can run the following command (same folder as the Dockefile): 

    docker-compose up -d