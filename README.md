To run the project, in the root folder run:
```
npm install
```
Afterwards, run:
```
npm run dev
```

This will start up the backend server, as well as the frontend.

To test the CLI functionality, navigate to the backend folder with
```
cd backend
```
Compile the project:
```
go build -o main
```

Then you can run various commands to interact with the database, for example:
```
// List all movie titles in the database
./main list

// Details of a movie given a valid IMDb ID
./main details -imdbid {imdb_id}

// Delete a movie given a valid IMDb ID
./main delete -imdbid {imdb_id}

// And more....
```
