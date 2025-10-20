import React from 'react';
import ProductCard from './ProductCard';

const mockProducts = [
  { id: 1, name: 'Lavender Dreams', scent: 'Calming lavender and chamomile', price: 5.99, image: 'https://placehold.co/400x400/e9d5ff/581c87?text=Lavender+Dreams' },
  { id: 2, name: 'Spiced Apple Cider', scent: 'Warm apple, cinnamon, and clove', price: 6.50, image: 'https://placehold.co/400x400/fed7aa/9a3412?text=Spiced+Apple' },
  { id: 3, name: 'Ocean Breeze', scent: 'Fresh sea salt and orchid', price: 5.99, image: 'https://placehold.co/400x400/a5f3fc/155e75?text=Ocean+Breeze' },
];

const HomePage = ({ setPage }) => {
    return (
        <div className="container mx-auto px-6 py-12">
            <div className="bg-cream-200 rounded-xl p-12 text-center shadow-lg">
                <h1 className="text-5xl font-extrabold font-serif text-brown-900">Warmth & Fragrance, Delivered.</h1>
                <p className="mt-4 text-lg text-brown-800">Hand-poured, high-quality wax melts to bring a cozy atmosphere to your home.</p>
                <button onClick={() => setPage('products')} className="mt-8 px-8 py-3 bg-brown-800 text-cream-100 text-lg font-semibold rounded-full hover:bg-brown-900 transition-transform transform hover:scale-105 shadow-md">Shop The Collection</button>
            </div>
            <div className="mt-20">
                <h2 className="text-4xl font-bold font-serif text-center text-brown-900">Our Favorites</h2>
                <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-10 mt-10">
                    {mockProducts.slice(0, 3).map(p => <ProductCard key={p.id} product={p} />)}
                </div>
            </div>
        </div>
    );
};

export default HomePage;
