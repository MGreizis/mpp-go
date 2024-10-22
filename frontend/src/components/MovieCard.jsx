import PropTypes from 'prop-types';

const MovieCard = ({ movie, onShowDetails }) => {
  return (
    <div
      className="bg-green-500 p-4 rounded-lg shadow-md hover:shadow-xl hover:bg-green-600 transition-shadow text-center cursor-pointer"
      onClick={() => onShowDetails(movie.imdb_id)}
    >
      <h2 className="text-xl font-bold mb-2">{movie.title}</h2>
      <p><strong>Year:</strong> {movie.year}</p>
      <p><strong>Rating:</strong> {movie.rating}</p>
      {movie.poster ? (
        <img src={movie.poster} alt={movie.title} className="w-full h-48 object-cover rounded mt-2" />
      ) : (
        <p>No poster available</p>
      )}
    </div>
  );
}

MovieCard.propTypes = {
  movie: PropTypes.shape({
    imdb_id: PropTypes.string.isRequired,
    title: PropTypes.string.isRequired,
    year: PropTypes.number.isRequired,
    rating: PropTypes.number.isRequired,
    poster: PropTypes.string,
  }).isRequired,
  onShowDetails: PropTypes.func.isRequired,
};

export default MovieCard;