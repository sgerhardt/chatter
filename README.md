1. Have an elven labs account
2. Setup a .env file with the following content
```
XI_API_KEY=<replace_me>
OUTPUT=output.mp3 # leave blank to use same directory as executable
```
3. Build the project
```
go build -o chatter
```

4. Run with text and voice values
```
chatter -t "Hello, World!" -v "your_voice_id"
```