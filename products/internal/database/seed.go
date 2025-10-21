package database

import (
	"context"
	"log"

	"com.MixieMelts.products/internal/models"
)

func (db *DB) Seed(ctx context.Context) {
	// Check if there are any products in the database
	count, err := db.getProductsCount(ctx)
	if err != nil {
		log.Printf("failed to get products count: %v", err)
		return
	}

	if count > 0 {
		log.Println("Database already seeded.")
		return
	}

	log.Println("Seeding database...")

	// Create some sample products
	products := []models.Product{
		{Category: "Year-Round", Name: "Serene Sanctuary", Scent: "Lavender, Chamomile, Cedarwood, Ylang Ylang", Description: "A calming, spa-like scent perfect for relaxation and de-stressing.", Price: 5.49, Image: "https://placehold.co/400x400/e0e7ff/4c1d95?text=Serene+Sanctuary"},
		{Category: "Year-Round", Name: "Citrus Sunshine", Scent: "Lemon, Bergamot, Sweet Orange, Spearmint", Description: "A bright, energizing, and clean aroma that uplifts the mood.", Price: 3.99, Image: "https://placehold.co/400x400/fef9c3/b45309?text=Citrus+Sunshine"},
		{Category: "Year-Round", Name: "Cozy Cashmere", Scent: "Vanilla Absolute, Sandalwood, Amyris, Peru Balsam", Description: "A warm, soft, and comforting scent like being wrapped in a favorite blanket.", Price: 6.49, Image: "https://placehold.co/400x400/f5f5f4/78350f?text=Cozy+Cashmere"},
		{Category: "Year-Round", Name: "Woodland Walk", Scent: "Eucalyptus, Rosemary, Cypress, Peppermint", Description: "A fresh, green, and earthy scent that brings the outdoors in.", Price: 4.99, Image: "https://placehold.co/400x400/dcfce7/14532d?text=Woodland+Walk"},
		{Category: "Year-Round", Name: "Rose Garden", Scent: "Rose Absolute, Palmarosa, Petitgrain", Description: "A classic, romantic, and elegant true floral scent.", Price: 5.99, Image: "https://placehold.co/400x400/fce7f3/9d174d?text=Rose+Garden"},
		{Category: "Spring", Name: "April Showers", Scent: "Oakmoss, Vetiver, Petitgrain, Clary Sage", Description: "The fresh, earthy scent of rain on soil and budding greenery.", Price: 4.99, Image: "https://placehold.co/400x400/dbeafe/1e3a8a?text=April+Showers"},
		{Category: "Spring", Name: "Wildflower Meadow", Scent: "Geranium, Lavender, Chamomile, Lemongrass", Description: "A sweet, light floral reminiscent of a field of blooming wildflowers.", Price: 5.49, Image: "https://placehold.co/400x400/f5d0fe/701a75?text=Wildflower+Meadow"},
		{Category: "Spring", Name: "Lilac Bloom", Scent: "Lilac Natural Fragrance, Ylang Ylang", Description: "The sweet, heady, and iconic fragrance of a blooming lilac bush.", Price: 6.99, Image: "https://placehold.co/400x400/ede9fe/5b21b6?text=Lilac+Bloom"},
		{Category: "Summer", Name: "Coastal Breeze", Scent: "Lime, Spearmint, Amyris, Coconut", Description: "A refreshing, vibrant scent like a mojito on the beach.", Price: 4.49, Image: "https://placehold.co/400x400/a5f3fc/155e75?text=Coastal+Breeze"},
		{Category: "Summer", Name: "Sun-Kissed Peach", Scent: "Peach Natural Fragrance, Sweet Orange, Vanilla", Description: "Sweet, juicy, and warm, like a ripe peach picked from the tree.", Price: 5.99, Image: "https://placehold.co/400x400/ffedd5/f97316?text=Sun-Kissed+Peach"},
		{Category: "Summer", Name: "Tropical Getaway", Scent: "Pineapple, Coconut, Lime, Vanilla", Description: "An exotic and sweet blend that transports you to a tropical island.", Price: 5.49, Image: "https://placehold.co/400x400/fef08a/eab308?text=Tropical+Getaway"},
		{Category: "Autumn", Name: "Autumn Harvest", Scent: "Apple Natural Fragrance, Cinnamon, Clove", Description: "The quintessential scent of fall—warm, spicy, and fruity.", Price: 4.99, Image: "https://placehold.co/400x400/fee2e2/991b1b?text=Autumn+Harvest"},
		{Category: "Autumn", Name: "Bonfire Flannel", Scent: "Cedarwood, Frankincense, Vetiver, Birch Tar", Description: "A smoky, woody, and cozy scent that evokes a crackling bonfire.", Price: 6.49, Image: "https://placehold.co/400x400/737373/171717?text=Bonfire+Flannel"},
		{Category: "Autumn", Name: "Pumpkin Spice", Scent: "Cinnamon, Clove, Ginger, Cardamom, Nutmeg", Description: "A comforting and classic blend of pumpkin and warm baking spices.", Price: 3.99, Image: "https://placehold.co/400x400/fed7aa/c2410c?text=Pumpkin+Spice"},
		{Category: "Winter", Name: "Winter Woods", Scent: "Pine Needle, Fir Balsam, Cypress, Cedarwood", Description: "A crisp, clean scent of a snow-covered evergreen forest.", Price: 5.49, Image: "https://placehold.co/400x400/ecfdf5/065f46?text=Winter+Woods"},
		{Category: "Winter", Name: "Spiced Cranberry", Scent: "Cranberry Natural Fragrance, Orange, Cinnamon", Description: "A festive and bright blend of tart fruit and warm spices.", Price: 4.99, Image: "https://placehold.co/400x400/fee2e2/dc2626?text=Spiced+Cranberry"},
		{Category: "Winter", Name: "Peppermint Cocoa", Scent: "Peppermint, Cocoa Absolute, Vanilla Absolute", Description: "A delicious and nostalgic mix of rich chocolate and cool, sweet mint.", Price: 5.99, Image: "https://placehold.co/400x400/d1fae5/78350f?text=Peppermint+Cocoa"},
		{Category: "Holiday", Name: "Witches' Brew", Scent: "Patchouli, Frankincense, Cinnamon, Clove", Description: "An earthy, spicy, and mysterious scent for a spooky atmosphere.", Price: 5.99, Image: "https://placehold.co/400x400/a78bfa/3b0764?text=Witches'+Brew"},
		{Category: "Holiday", Name: "Haunted Hayride", Scent: "Hay Absolute, Vetiver, Amyris, Cedarwood", Description: "The earthy smell of dry hay, damp fallen leaves, and distant woods.", Price: 6.99, Image: "https://placehold.co/400x400/fde68a/713f12?text=Haunted+Hayride"},
		{Category: "Holiday", Name: "Christmas Tree", Scent: "Pine Needle, Fir Balsam, Sweet Orange", Description: "The fresh, nostalgic, and beloved scent of a freshly cut Christmas tree.", Price: 4.99, Image: "https://placehold.co/400x400/bbf7d0/166534?text=Christmas+Tree"},
		{Category: "Holiday", Name: "Gingerbread House", Scent: "Ginger, Cinnamon, Nutmeg, Clove, Vanilla", Description: "Warm, spicy, and sweet, just like a freshly decorated gingerbread house.", Price: 5.49, Image: "https://placehold.co/400x400/f3e8ff/92400e?text=Gingerbread+House"},
	}

	for _, product := range products {
		query := `
		INSERT INTO products (name, category, scent, price, subscription, image, description)
		VALUES ($1, $2, $3, $4, $5, $6, $7);
		`
		_, err := db.Exec(ctx, query, product.Name, product.Category, product.Scent, product.Price, product.Subscription, product.Image, product.Description)
		if err != nil {
			log.Printf("failed to create product %s: %v", product.Name, err)
		}
	}

	log.Println("Database seeded successfully.")
}

func (db *DB) getProductsCount(ctx context.Context) (int, error) {
	var count int
	query := "SELECT COUNT(*) FROM products"
	err := db.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
