import { useState } from "react";

const AddProductForm = () => {
  const [name, setName] = useState("");
  const [scent, setScent] = useState("");
  const [price, setPrice] = useState("");
  const [image, setImage] = useState("");
  const [description, setDescription] = useState("");
  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError("");
    setSuccess("");

    try {
      const response = await fetch("/products", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          name,
          scent,
          price: parseFloat(price),
          image,
          description,
        }),
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.message || "Failed to add product");
      }

      setSuccess("Product added successfully!");
      // Clear form
      setName("");
      setScent("");
      setPrice("");
      setImage("");
      setDescription("");
    } catch (err) {
      setError(err.message);
    }
  };

  return (
    <div className="bg-white p-6 rounded-xl shadow-lg">
      <h2 className="text-2xl font-serif font-bold text-brown-900 mb-4">
        Add New Product
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
          <label htmlFor="scent" className="text-sm font-medium text-brown-800">
            Scent
          </label>
          <input
            id="scent"
            type="text"
            value={scent}
            onChange={(e) => setScent(e.target.value)}
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
          <label htmlFor="image" className="text-sm font-medium text-brown-800">
            Image URL
          </label>
          <input
            id="image"
            type="text"
            value={image}
            onChange={(e) => setImage(e.target.value)}
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
        <button
          type="submit"
          className="w-full py-3 px-4 text-lg font-semibold text-cream-100 bg-terracotta-500 rounded-lg hover:bg-terracotta-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-terracotta-500"
        >
          Add Product
        </button>
      </form>
    </div>
  );
};

export default AddProductForm;
