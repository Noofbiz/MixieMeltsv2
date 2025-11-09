import React, { useEffect, useState } from "react";

/**
 * ProductDetailPage
 *
 * Props:
 * - productId (number|string) : id of the product to fetch
 * - onBack (function) : optional callback to navigate back (e.g. setPage('products'))
 *
 * This component fetches a single product from /products/{id} and displays:
 * - image, name, scent/description, price
 * - the recipe (list of ingredients with amount/unit/notes)
 *
 * Note: the products API returns `recipe` (array) for the product.
 *      Each ingredient in the recipe has: id, name, amount, unit, notes
 */
const ProductDetailPage = ({ productId, onBack }) => {
  const [product, setProduct] = useState(null);
  const [loading, setLoading] = useState(Boolean(productId));
  const [error, setError] = useState(null);

  useEffect(() => {
    if (!productId) {
      setProduct(null);
      setLoading(false);
      return;
    }

    let mounted = true;
    const fetchProduct = async () => {
      setLoading(true);
      setError(null);
      try {
        const resp = await fetch(`/products/${productId}`);
        if (!resp.ok) {
          // try to read an error message if present
          let msg = `Failed to fetch product (status ${resp.status})`;
          try {
            const body = await resp.json();
            if (body && body.message) msg = body.message;
          } catch {
            // ignore JSON parse errors
          }
          throw new Error(msg);
        }
        const data = await resp.json();
        if (mounted) {
          setProduct(data);
        }
      } catch (err) {
        if (mounted) setError(err.message || "Unknown error");
      } finally {
        if (mounted) setLoading(false);
      }
    };

    fetchProduct();
    return () => {
      mounted = false;
    };
  }, [productId]);

  const handleBack = () => {
    if (typeof onBack === "function") {
      onBack();
      return;
    }
    // fallback: try browser history
    window.history.back();
  };

  return (
    <div className="container mx-auto px-6 py-12">
      <button
        onClick={handleBack}
        className="mb-6 px-4 py-2 bg-cream-200 rounded-full border border-brown-300"
      >
        ← Back
      </button>

      {loading ? (
        <div className="text-center text-brown-700">Loading...</div>
      ) : error ? (
        <div className="text-center text-red-600">Error: {error}</div>
      ) : !product ? (
        <div className="text-center text-brown-700">Product not found</div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
          <div className="md:col-span-1">
            <img
              src={product.image}
              alt={product.name}
              className="w-full h-80 object-cover rounded-lg shadow"
            />
          </div>

          <div className="md:col-span-2">
            <h1 className="text-3xl font-serif font-bold">{product.name}</h1>
            <p className="mt-2 text-brown-700">
              {product.scent || product.description}
            </p>
            <p className="mt-4 text-brown-800 text-xl font-semibold">
              ${Number(product.price || 0).toFixed(2)}
            </p>

            <h2 className="mt-8 text-2xl font-semibold">Recipe</h2>
            {product.recipe && product.recipe.length > 0 ? (
              <ul className="mt-4 list-disc list-inside text-brown-700">
                {product.recipe.map((ing) => (
                  <li key={ing.id || `${ing.name}-${ing.unit}`}>
                    <span className="font-medium">{ing.name}</span>
                    {" — "}
                    <span>
                      {ing.amount} {ing.unit}
                    </span>
                    {ing.notes ? (
                      <span className="text-sm text-brown-600">
                        {" "}
                        {"(" + ing.notes + ")"}
                      </span>
                    ) : null}
                  </li>
                ))}
              </ul>
            ) : (
              <p className="mt-4 text-brown-700">
                No recipe available for this product.
              </p>
            )}
          </div>
        </div>
      )}
    </div>
  );
};

export default ProductDetailPage;
