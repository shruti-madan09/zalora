## Setup Instructions

### Using docker
1. Navigate to project folder zalora
2. Build docker image using command:
  * ****docker-compose build****
3. Run db container using command:
  * ****docker-compose up db****
4. Open another tab in terminal to run zalora container using command:
  * ****docker-compose up zalora****
5. Upload DB schema to container
  * run command: ****docker ps****
  * container ids of both db and zalora will be displayed in terminal
  * access mysql container using command: ****docker exec -it mysql_container_id bash****
  * upload schema using command: ****mysql -uroot -ppassword bennjerry < docker-entrypoint-initdb.d/bennjerry.sql****
6. Upload data from icecream.json
  * run command: ****docker ps****
  * container ids of both db and zalora will be displayed in terminal
  * access zalora container using command: ****docker exec -it zalora_container_id bash****
  * navigate to uploader package using command: cd /workspace/zalora/src/uploader/
  * run file using command: ****go run upload.go****
7. Read logs
  * access the zalora container as explained in step 6
  * view logs using command: tail -f /workspace/zalora/logs/zalora.log
8. Run test cases
  * access the zalora container as explained in step 6
  * navigate to test package using command: cd /workspace/zalora/src/bennjerry/test
  * run create api test cases using command: ****go test -v create_test.go****
  * run read api test cases using command: ****go test -v create_read.go****
  * run update api test cases using command: ****go test -v create_update.go****
  * run delete api test cases using command: ****go test -v create_delete.go****
  * logs created while running test cases will be in /workspace/zalora/src/bennjerry/test/logs/zalora.log

Points to note for docker setup
  * Steps 5 and 6 need to only be executed the first time.
  * Steps 1 to 4 need to be run every time to run the server


## Without using docker
1. Install go, 1.13.1
2. Set GOPATH
  * edit bash_profile using command: ****vim ~/.bash_profile****
  * add lines:<br/>
    export GOPATH=$HOME/workspace/zalora/vendor:$HOME/workspace/zalora<br/>
    export PATH=$PATH:$GOPATH/bin
  * save file and run command: ****source ~/.bash_profile****
3. Install mysql, 8.0.17
4. Create database using command: ****create database bennjerry****
5. Upload DB schema
  * navigate to zalora folder
  * run command: ****mysql -uroot -ppassword bennjerry < bennjerry.sql****
6. Upload data from icecream.json
  * navigate to uploader package using command: cd zalora/src/uploader
  * run script using command: ****go run upload.go****
7. Run Server
  * navigate to zalora folder
  * run command: ****make****
  * run command: ****./bin/zalora****
8. Read logs
  * navigate to zalora folder
  * run command: tail -f logs/zalora.log
9. Run test cases
  * navigate to test package using command: cd zalora/src/bennjerry/test
  * run create api test cases using command: ****go test -v create_test.go****
  * run read api test cases using command: ****go test -v create_read.go****
  * run update api test cases using command: ****go test -v create_update.go****
  * run delete api test cases using command: ****go test -v create_delete.go****
  * logs created while running test cases will be in /workspace/zalora/src/bennjerry/test/logs/zalora.log
 
Points to note for manual setup
  * Steps 1 to 6 need to only be executed the first time.
  * Step 7 needs to be run every time to run the server.
  
  
  
