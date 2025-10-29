import React, { useState, useEffect } from "react";
import ProductCard from "./ProductCard";

const ProductListPage = ({ goTo }) => {
  const [products, setProducts] = useState([]);
  const [subscriptionBoxes, setSubscriptionBoxes] = useState([]);

  useEffect(() => {
    const fetchProducts = async () => {
      try {
        const response = await fetch("/products");
        if (response.ok) {
          const data = await response.json();
          setProducts(data);
        }
      } catch (error) {
        console.error("Failed to fetch products", error);
      }
    };

    const fetchSubscriptionBoxes = async () => {
      try {
        const response = await fetch("/products/subscription-boxes");
        if (response.ok) {
          const data = await response.json();
          setSubscriptionBoxes(data);
        }
      } catch (error) {
        console.error("Failed to fetch subscription boxes", error);
      }
    };

    fetchProducts();
    fetchSubscriptionBoxes();
  }, []);

  return (
    <div className="container mx-auto px-6 py-12">
      <h1 className="text-4xl font-bold font-serif text-center text-brown-900 mb-12">
        Subscription Boxes
      </h1>
      <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-8">
        {subscriptionBoxes.map((box) => (
          <ProductCard
            key={box.id}
            product={{
              ...box,
              scent: box.description,
              image: box.image,
              isSubscription: true,
            }}
            onOpen={(id) => goTo("product", id)}
          />
        ))}
      </div>

      <h2 className="text-3xl font-bold font-serif text-center text-brown-900 mt-20 mb-12">
        Our Collection
      </h2>
      <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-8">
        {products.map((product) => (
          <ProductCard
            key={product.id}
            product={product}
            onOpen={(id) => goTo("product", id)}
          />
        ))}
      </div>
    </div>
  );
};

export default ProductListPage;
