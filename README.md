# lampstand
A service for automatically formatting Bible verses into HTML5 slides using revealJS. It uses SQLite for the database.

## REST API
The lampstand project uses a single path to service out verses at:  
**/api/:version/verses?passage=Luke 11:33**  
where :version is the version (NIV, HSCB, ESV), and the passage is the passage you wnat to look up.  
The passage json object contains the reference, version, and a JSON array of verses:

```
{"reference":"Luke 11:33","version":"NIV","verses":[{"book":"Luke","chapter":11,"verseNo":33,"text":"\"No-one lights a lamp and puts it in a place where it will be hidden, or under a bowl. Instead he puts it on its stand, so that those who come in may see the light."}]}
```

##Usage

1. Get the Go dependencies by running bootstrap.sh
2. Build the go project with go build
3. Run setup_table.sh to create the bible.db sqlite database, and move it to the location of the binary
4. Run the binary, it should be served at localhost:8080
