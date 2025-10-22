import { useState, useEffect } from "react";

const AddSubscriptionBoxForm = () => {
  const [name, setName] = useState("");
  const [price, setPrice] = useState("");
  const [description, setDescription] = useState("");
  const [products, setProducts] = useState([]);
  const [selectedProducts, setSelectedProducts] = useState([]);
  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");

  useEffect(() => {
    const fetchProducts = async () => {
      try {
        const response = await fetch("/products");
        if (!response.ok) {
          throw new Error("Failed to fetch products");
        }
        const data = await response.json();
        setProducts(data);
      } catch (err) {
        setError(err.message);
      }
    };
    fetchProducts();
  }, []);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError("");
    setSuccess("");

    try {
      const response = await fetch("/products/subscription-boxes", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          name,
          price: parseFloat(price),
          description,
          product_ids: selectedProducts,
        }),
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.message || "Failed to add subscription box");
      }

      setSuccess("Subscription box added successfully!");
      // Clear form
      setName("");
      setPrice("");
      setDescription("");
      setSelectedProducts([]);
    } catch (err) {
      setError(err.message);
    }
  };

  const handleProductSelection = (productId) => {
    setSelectedProducts((prevSelected) =>
      prevSelected.includes(productId)
        ? prevSelected.filter((id) => id !== productId)
        : [...prevSelected, productId]
    );
  };

  return (
    <div className="bg-white p-6 rounded-xl shadow-lg">
      <h2 className="text-2xl font-serif font-bold text-brown-900 mb-4">
        Add New Subscription Box
      </h2>
      {error && <p className="text-red-500 text-center">{error}</p>}
      {success && <p className="text-green-500 text-center">{success}</p>}
      <form className="space-y-4" onSubmit={handleSubmit}>
        <div>
          <label htmlFor="name" className="text-sm font-medium text-brown-800">
            Name
          </label>
          <input
            id="name"
            type="text"
            value={name}
            onChange={(e) => setName(e.target.value)}
            className="mt-1 block w-full px-3 py-2 bg-white border border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-terracotta-500 focus:border-terracotta-500"
            required
          />
        </div>
        <div>
          <label htmlFor="price" className="text-sm font-medium text-brown-800">
            Price
          </label>
          <input
            id="price"
            type="number"
            step="0.01"
            value={price}
            onChange={(e) => setPrice(e.target.value)}
            className="mt-1 block w-full px-3 py-2 bg-white border border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-terracotta-500 focus:border-terracotta-500"
            required
          />
        </div>
        <div>
          <label
            htmlFor="description"
            className="text-sm font-medium text-brown-800"
          >
            Description
          </label>
          <textarea
            id="description"
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            className="mt-1 block w-full px-3 py-2 bg-white border border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-terracotta-500 focus:border-terracotta-500"
            required
          />
        </div>
        <div>
          <label className="text-sm font-medium text-brown-800">
            Select Products
          </label>
          <div className="grid grid-cols-2 gap-4 mt-2">
            {products.map((product) => (
              <div
                key={product.id}
                className={`p-4 border rounded-md cursor-pointer ${
                  selectedProducts.includes(product.id)
                    ? "border-terracotta-500 bg-terracotta-50"
                    : "border-gray-300"
                }`}
                onClick={() => handleProductSelection(product.id)}
              >
                <p className="font-semibold">{product.name}</p>
                <p className="text-sm text-gray-500">{product.scent}</p>
              </div>
            ))}
          </div>
        </div>
        <button
          type="submit"
          className="w-full py-3 px-4 text-lg font-semibold text-cream-100 bg-terracotta-500 rounded-lg hover:bg-terracotta-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-terracotta-500"
        >
          Add Subscription Box
        </button>
      </form>
    </div>
  );
};

export default AddSubscriptionBoxForm;
