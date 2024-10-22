import PropTypes from 'prop-types';
import Modal from 'react-modal';

Modal.setAppElement('#root'); // Set the app element for accessibility

const customStyles = {
  content: {
    top: '50%',
    left: '50%',
    right: 'auto',
    bottom: 'auto',
    marginRight: '-50%',
    transform: 'translate(-50%, -50%)',
    padding: '20px',
    borderRadius: '10px',
    width: '50%',
    maxHeight: '80%',
    overflowY: 'auto',
  },
  overlay: {
    backgroundColor: 'rgba(0, 0, 0, 0.75)',
    zIndex: 1000,
  },
};

const MovieModal = ({ isOpen, onClose, movie }) => {
  if (!isOpen || !movie) return null;

  return (
    <Modal
      isOpen={isOpen}
      onRequestClose={onClose}
      contentLabel="Movie Details"
      style={customStyles}
    >
      <button onClick={onClose} className="absolute top-2 right-2 text-gray-500 hover:text-gray-700">
        &times;
      </button>
      <h2 className="text-xl font-semibold mb-2">Movie Details</h2>
      <p><strong>IMDb ID:</strong> {movie.imdb_id}</p>
      <p><strong>Title:</strong> {movie.title}</p>
      <p><strong>Year:</strong> {movie.year}</p>
      <p><strong>Rating:</strong> {movie.rating}</p>
      <p><strong>Poster:</strong> {movie.poster ? <img src={movie.poster} alt={movie.title} className="w-full h-48 object-cover rounded mt-2" /> : 'No poster available'}</p>
    </Modal>
  );
}

MovieModal.propTypes = {
  isOpen: PropTypes.bool.isRequired,
  onClose: PropTypes.func.isRequired,
  movie: PropTypes.shape({
    imdb_id: PropTypes.string.isRequired,
    title: PropTypes.string.isRequired,
    year: PropTypes.number.isRequired,
    rating: PropTypes.number.isRequired,
    poster: PropTypes.string,
  }),
};

export default MovieModal;