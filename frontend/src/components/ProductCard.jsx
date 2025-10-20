import { useCart } from "../context/Context";

const ProductCard = ({ product }) => {
  const { dispatch } = useCart();
  const handleAddToCart = () => {
    dispatch({ type: "ADD_ITEM", payload: product });
  };

  return (
    <div className="bg-white rounded-xl shadow-lg overflow-hidden transform hover:scale-105 hover:shadow-xl transition-all duration-300 group">
      <img
        src={product.image}
        alt={product.name}
        className="w-full h-56 object-cover"
      />
      <div className="p-6">
        <h3 className="text-xl font-serif font-semibold text-brown-900">
          {product.name}
        </h3>
        <p className="text-brown-700 mt-2 text-sm">{product.scent}</p>
        <div className="mt-4 flex justify-between items-center">
          <span className="text-2xl font-bold font-sans text-brown-800">
            ${product.price.toFixed(2)}
          </span>
          <button
            onClick={handleAddToCart}
            className="px-4 py-2 bg-brown-700 text-cream-100 text-sm font-semibold rounded-full hover:bg-brown-800 transition-colors transform group-hover:scale-110"
          >
            {product.isSubscription ? "Subscribe" : "Add to Cart"}
          </button>
        </div>
      </div>
    </div>
  );
};

export default ProductCard;
