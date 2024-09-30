import {useState, useEffect} from 'react'
// import "./output.css"

function App() {
  const [movies, setMovies] = useState([]);
  const [newMovie, setNewMovie] = useState({
    imdbId: '',
    title: '',
    year: '',
    rating: '',
    poster: ''
  })
  const [movieDetails, setMovieDetails] = useState(null);

  useEffect(() => {
    const fetchMovies = async () => {
      try {
        const response = await fetch('http://localhost:8090/movies');
        const data = await response.json();
        setMovies(data);
      } catch (error) {
        console.error('Error fetching movies: ', error);
      }
    };
    fetchMovies();
  }, []);

  /**
   * Handles submitting the form for adding a movie.
   * @param {React.FormEvent<HTMLFormElement>} e - The event object.
   */
  const handleAddMovie = async (e) => {
    e.preventDefault();
    try {
      const response = await fetch('http://localhost:8090/movies', {
        method: 'POST',
        mode: 'no-cors',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(newMovie)
      });

      if (response.ok) {
        const movie = await response.json();
        setMovies([...movies, movie]);
        // Reset the form
        setNewMovie({
          imdbId: '',
          title: '',
          year: 0,
          rating: 0.0,
          poster: ''
        })
      } else {
        console.error('Error adding movie: ', response.statusText);
      }
    } catch (error) {
      console.error('Error adding movie: ', error);
    }
  };

  /**
   * Handles fetching the movie details by IMDb ID. If the movie can be found,
   * it updates the movieDetails state with the movie object. If the movie
   * cannot be found, it logs an error message to the console.
   * @param {string} imdbID - The IMDb ID of the movie to show the details for.
   */
  const handleShowDetails = async (imdbID) => {
    try {
      const response = await fetch(`http://localhost:8090/movies/${imdbID}`);
      
      if (response.ok) {
        const movie = await response.json();
        setMovieDetails(movie);
      } else {
        console.error('Error fetching movie details: ', response.statusText);
      }
    } catch (error) {
      console.error('Error fetching movie details: ', error);
    }
  }

  /**
   * Handles deleting a movie by its IMDb ID. If the movie can be deleted,
   * it updates the movies state to exclude the movie. If the movie cannot be
   * deleted, it logs an error message to the console.
   * @param {string} imdbID - The IMDb ID of the movie to delete.
   */
  const handleDeleteMovie = async (imdbID) => {
    try {
      const response = await fetch(`http://localhost:8090/movies/${imdbID}`, {
        method: 'DELETE'
      });

      if (response.ok) {
        setMovies(movies.filter((movie) => movie.imdb_id !== imdbID));
      } else {
        console.error('Error deleting movie: ', response.statusText);
      }
    } catch (error) {
      console.error('Error deleting movie: ', error);
    }
  };

  return (
    <main className="p-6 bg-gray-600 min-h-screen">
      <h2 className="text-2xl text-white font-semibold mb-4">Add New Movie</h2>
      
      <form onSubmit={handleAddMovie} className="bg-white p-6 rounded-lg shadow-md space-y-4 mb-4">
        <input
          type='text'
          placeholder='IMDb ID'
          value={newMovie.imdbId}
          onChange={(e) => setNewMovie({ ...newMovie, imdbId: e.target.value })}
          className="w-full px-4 py-2 border border-gray-300 rounded-lg"
        />
        <input
          type='text'
          placeholder='Title'
          value={newMovie.title}
          onChange={(e) => setNewMovie({ ...newMovie, title: e.target.value })}
          className="w-full px-4 py-2 border border-gray-300 rounded-lg"
        />
        <input
          type='number'
          placeholder='Year'
          value={newMovie.year}
          onChange={(e) => setNewMovie({ ...newMovie, year: e.target.value })}
          className="w-full px-4 py-2 border border-gray-300 rounded-lg"
        />
        <input
          type='number'
          step="0.1"
          placeholder='Rating'
          value={newMovie.rating}
          onChange={(e) => setNewMovie({ ...newMovie, rating: e.target.value })}
          className="w-full px-4 py-2 border border-gray-300 rounded-lg"
        />
        <input
          type='text'
          placeholder='Poster'
          value={newMovie.poster}
          onChange={(e) => setNewMovie({ ...newMovie, poster: e.target.value })}
          className="w-full px-4 py-2 border border-gray-300 rounded-lg"
        />
        <button 
          type='submit'
          className="bg-green-500 hover:bg-green-600 text-white px-6 py-2 rounded-lg w-full"
        >
          Add Movie
        </button>
      </form>

      <h1 className="text-3xl font-bold text-center mb-6 text-white">Movie List</h1>
      <p className="text-center mb-6 text-white">FYI, the details of the movie show up at the bottom of the page :/ (will be fixed soon)</p>

      <ul className="space-y-4 mb-6">
        {movies.map((movie) => (
          <li key={movie.imdb_id} className="flex justify-between items-center bg-white p-4 rounded-lg shadow-md">
            <span>{movie.title}</span>
            <div className="space-x-4">
              <button 
                className="bg-blue-500 hover:bg-blue-600 text-white px-4 py-2 rounded"
                onClick={() => handleShowDetails(movie.imdb_id)}
              >
                Details
              </button>
              <button 
                className="bg-red-500 hover:bg-red-600 text-white px-4 py-2 rounded"
                onClick={() => handleDeleteMovie(movie.imdb_id)}
              >
                Delete
              </button>
            </div>
          </li>
        ))}
      </ul>

      {movieDetails && (
        <div className="mt-6 bg-white p-6 rounded-lg shadow-md">
          <h2 className="text-xl font-semibold mb-2">Movie Details</h2>
          <p><strong>IMDb ID:</strong> {movieDetails.imdb_id}</p>
          <p><strong>Title:</strong> {movieDetails.title}</p>
          <p><strong>Year:</strong> {movieDetails.year}</p>
          <p><strong>Rating:</strong> {movieDetails.rating}</p>
          <p><strong>Poster:</strong> {movieDetails.poster?.Valid ? movieDetails.poster.String : 'No poster available'}</p>
        </div>
      )}
    </main>
  )
}

export default App