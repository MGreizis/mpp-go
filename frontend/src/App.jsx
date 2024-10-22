import { useState, useEffect } from 'react';
import MovieCard from './components/MovieCard';
import MovieModal from './components/MovieModal';

function App() {
  const [movies, setMovies] = useState([]);
  const [newMovie, setNewMovie] = useState({
    imdbId: '',
    title: '',
    year: '',
    rating: '',
    poster: ''
  });
  const [sortBy, setSortBy] = useState('year');
  const [order, setOrder] = useState('asc');
  const [selectedMovie, setSelectedMovie] = useState(null);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [currentPage, setCurrentPage] = useState(1);
  const [moviesPerPage] = useState(15);

  useEffect(() => {
    const fetchMovies = async () => {
      try {
        const response = await fetch(`http://localhost:8090/movies?sort=${sortBy}&order=${order}&page=${currentPage}&limit=${moviesPerPage}`);
        const data = await response.json();
        setMovies(data);
      } catch (error) {
        console.error('Error fetching movies: ', error);
      }
    };
    fetchMovies();
  }, [sortBy, order, currentPage, moviesPerPage]);

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
        setNewMovie({
          imdbId: '',
          title: '',
          year: '',
          rating: '',
          poster: ''
        });
      } else {
        console.error('Error adding movie: ', response.statusText);
      }
    } catch (error) {
      console.error('Error adding movie: ', error);
    }
  };

  const handleShowDetails = (imdbID) => {
    const movie = movies.find((movie) => movie.imdb_id === imdbID);
    setSelectedMovie(movie);
    setIsModalOpen(true);
  };

  const handleCloseModal = () => {
    setIsModalOpen(false);
    setSelectedMovie(null);
  };

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

  // Calculate the index of the first and last movie on the current page
  const indexOfLastMovie = currentPage * moviesPerPage;
  const indexOfFirstMovie = indexOfLastMovie - moviesPerPage;
  const currentMovies = movies.slice(indexOfFirstMovie, indexOfLastMovie);

  // Calculate total pages
  const totalPages = Math.ceil(movies.length / moviesPerPage);

  // Change page
  const paginate = (pageNumber) => setCurrentPage(pageNumber);

  return (
    <main className="p-6 bg-gray-800 min-h-screen">
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

      <div className="flex space-x-4 mb-4">
        <div>
          <label className="text-white">Sort by: </label>
          <select value={sortBy} onChange={(e) => setSortBy(e.target.value)} className="ml-2 px-2 py-1 rounded">
            <option value="year">Year</option>
            <option value="rating">Rating</option>
          </select>
        </div>

        <div>
          <label className="text-white">Order: </label>
          <select value={order} onChange={(e) => setOrder(e.target.value)} className="ml-2 px-2 py-1 rounded">
            <option value="asc">Ascending</option>
            <option value="desc">Descending</option>
          </select>
        </div>
      </div>

      <h1 className="text-3xl font-bold text-center mb-6 text-white">Movie List</h1>

      <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-4">
        {currentMovies.map((movie) => (
          <MovieCard key={movie.imdb_id} movie={movie} onShowDetails={handleShowDetails} />
        ))}
      </div>

      <div className="flex justify-center mt-4">
        <button
          onClick={() => paginate(currentPage - 1)}
          disabled={currentPage === 1}
          className="px-4 py-2 mx-1 bg-gray-300 rounded hover:bg-gray-400 disabled:opacity-50"
        >
          Previous
        </button>
        <span className="px-4 py-2 mx-1 text-white">
          Page {currentPage} of {totalPages}
        </span>
        <button
          onClick={() => paginate(currentPage + 1)}
          disabled={currentPage === totalPages}
          className="px-4 py-2 mx-1 bg-gray-300 rounded hover:bg-gray-400 disabled:opacity-50"
        >
          Next
        </button>
      </div>

      <MovieModal isOpen={isModalOpen} onClose={handleCloseModal} movie={selectedMovie} />
    </main>
  );
}

export default App;