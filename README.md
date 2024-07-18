1. Have an eleven labs account
2. Setup a .env file with the following content
```
XI_API_KEY=<replace_me>
OUTPUT=/user/downloads/example # leave blank to use same directory as executable
```
3. Build the project
```
make chatter
```

4. Run with text and voice values
```
./bin/chatter -t "Hello, World!" -v "your_voice_id" 
```
or point it to a website
```
./bin/chatter -s "https://www.example.com" -v "your_voice_id"
```