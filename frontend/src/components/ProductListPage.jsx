import ProductCard from "./ProductCard";

const mockProducts = [
  {
    id: 1,
    name: "Lavender Dreams",
    scent: "Calming lavender and chamomile",
    price: 5.99,
    image: "https://placehold.co/400x400/e9d5ff/581c87?text=Lavender+Dreams",
  },
  {
    id: 2,
    name: "Spiced Apple Cider",
    scent: "Warm apple, cinnamon, and clove",
    price: 6.5,
    image: "https://placehold.co/400x400/fed7aa/9a3412?text=Spiced+Apple",
  },
  {
    id: 3,
    name: "Ocean Breeze",
    scent: "Fresh sea salt and orchid",
    price: 5.99,
    image: "https://placehold.co/400x400/a5f3fc/155e75?text=Ocean+Breeze",
  },
  {
    id: 4,
    name: "Vanilla Bean Noel",
    scent: "Sweet vanilla bean and caramel",
    price: 7.0,
    image: "https://placehold.co/400x400/fef3c7/b45309?text=Vanilla+Bean",
  },
  {
    id: 5,
    name: "Cozy Cashmere",
    scent: "Soft cashmere, white amber, and musk",
    price: 7.5,
    image: "https://placehold.co/400x400/e5e7eb/4b5563?text=Cozy+Cashmere",
  },
  {
    id: 6,
    name: "Monthly Subscription Box",
    scent: "A curated selection of seasonal melts",
    price: 24.99,
    image: "https://placehold.co/400x400/dcfce7/166534?text=Subscription+Box",
    isSubscription: true,
  },
];

const ProductListPage = () => {
  return (
    <div className="container mx-auto px-6 py-12">
      <h1 className="text-4xl font-bold font-serif text-center text-brown-900 mb-12">
        Our Full Collection
      </h1>
      <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-8">
        {mockProducts.map((product) => (
          <ProductCard key={product.id} product={product} />
        ))}
      </div>
    </div>
  );
};

export default ProductListPage;
