package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/PrimeraAizen/e-comm/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	ctx := context.Background()

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Build MongoDB URI
	mongoURI := cfg.Mongo.URI
	if mongoURI == "" {
		if cfg.Mongo.Username != "" && cfg.Mongo.Password != "" {
			mongoURI = fmt.Sprintf("mongodb://%s:%s@%s:%s",
				cfg.Mongo.Username, cfg.Mongo.Password, cfg.Mongo.Host, cfg.Mongo.Port)
		} else {
			mongoURI = fmt.Sprintf("mongodb://%s:%s", cfg.Mongo.Host, cfg.Mongo.Port)
		}
	}

	fmt.Println("Connecting to MongoDB...")
	fmt.Println("URI:", mongoURI)
	fmt.Println("Database:", cfg.Mongo.Database)

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal("Failed to connect:", err)
	}
	defer client.Disconnect(ctx)

	// Ping to verify connection
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("Failed to ping MongoDB:", err)
	}

	db := client.Database(cfg.Mongo.Database)

	// Clear existing data
	fmt.Println("Clearing existing data...")
	collections := []string{"users", "roles", "user_roles", "categories", "products",
		"orders", "order_items", "user_product_views", "user_product_likes", "profiles"}
	for _, coll := range collections {
		db.Collection(coll).Drop(ctx)
	}

	fmt.Println("Seeding data...")

	// Seed Roles
	fmt.Println("Creating roles...")
	rolesCollection := db.Collection("roles")
	roles := []interface{}{
		bson.M{"_id": 1, "name": "admin", "description": "System administrator", "created_at": time.Now(), "updated_at": time.Now()},
		bson.M{"_id": 2, "name": "user", "description": "Regular user", "created_at": time.Now(), "updated_at": time.Now()},
		bson.M{"_id": 3, "name": "moderator", "description": "Content moderator", "created_at": time.Now(), "updated_at": time.Now()},
		bson.M{"_id": 4, "name": "student", "description": "Student user", "created_at": time.Now(), "updated_at": time.Now()},
		bson.M{"_id": 5, "name": "teacher", "description": "Teacher user", "created_at": time.Now(), "updated_at": time.Now()},
	}
	_, err = rolesCollection.InsertMany(ctx, roles)
	if err != nil {
		log.Fatal("Failed to insert roles:", err)
	}

	// Seed Users
	fmt.Println("Creating users...")
	usersCollection := db.Collection("users")

	// Generate password hash for "password123"
	hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	passwordHash := string(hash)

	users := []interface{}{
		bson.M{
			"_id":           1,
			"email":         "admin@example.com",
			"password_hash": passwordHash,
			"status":        "active",
			"created_at":    time.Now(),
			"updated_at":    time.Now(),
		},
		bson.M{
			"_id":           2,
			"email":         "moderator@example.com",
			"password_hash": passwordHash,
			"status":        "active",
			"created_at":    time.Now(),
			"updated_at":    time.Now(),
		},
		bson.M{
			"_id":           3,
			"email":         "user1@example.com",
			"password_hash": passwordHash,
			"status":        "active",
			"created_at":    time.Now(),
			"updated_at":    time.Now(),
		},
		bson.M{
			"_id":           4,
			"email":         "user2@example.com",
			"password_hash": passwordHash,
			"status":        "active",
			"created_at":    time.Now(),
			"updated_at":    time.Now(),
		},
		bson.M{
			"_id":           5,
			"email":         "student@example.com",
			"password_hash": passwordHash,
			"status":        "active",
			"created_at":    time.Now(),
			"updated_at":    time.Now(),
		},
		bson.M{
			"_id":           6,
			"email":         "teacher@example.com",
			"password_hash": passwordHash,
			"status":        "active",
			"created_at":    time.Now(),
			"updated_at":    time.Now(),
		},
	}
	_, err = usersCollection.InsertMany(ctx, users)
	if err != nil {
		log.Fatal("Failed to insert users:", err)
	}

	// Seed User Roles
	fmt.Println("Assigning roles to users...")
	userRolesCollection := db.Collection("user_roles")
	userRoles := []interface{}{
		bson.M{"user_id": 1, "role_id": 1}, // admin
		bson.M{"user_id": 2, "role_id": 3}, // moderator
		bson.M{"user_id": 3, "role_id": 2}, // user
		bson.M{"user_id": 4, "role_id": 2}, // user
		bson.M{"user_id": 5, "role_id": 4}, // student
		bson.M{"user_id": 6, "role_id": 5}, // teacher
	}
	_, err = userRolesCollection.InsertMany(ctx, userRoles)
	if err != nil {
		log.Fatal("Failed to insert user roles:", err)
	}

	// Seed Categories
	fmt.Println("Creating categories...")
	categoriesCollection := db.Collection("categories")
	categories := []interface{}{
		bson.M{"_id": 1, "name": "Electronics", "description": "Electronic devices and accessories", "parent_id": nil, "created_at": time.Now(), "updated_at": time.Now()},
		bson.M{"_id": 2, "name": "Smartphones", "description": "Mobile phones", "parent_id": 1, "created_at": time.Now(), "updated_at": time.Now()},
		bson.M{"_id": 3, "name": "Tablets", "description": "Tablet devices", "parent_id": 1, "created_at": time.Now(), "updated_at": time.Now()},
		bson.M{"_id": 4, "name": "Laptops", "description": "Notebook computers", "parent_id": 1, "created_at": time.Now(), "updated_at": time.Now()},
		bson.M{"_id": 5, "name": "Accessories", "description": "Tech accessories", "parent_id": 1, "created_at": time.Now(), "updated_at": time.Now()},
	}
	_, err = categoriesCollection.InsertMany(ctx, categories)
	if err != nil {
		log.Fatal("Failed to insert categories:", err)
	}

	// Seed Products
	fmt.Println("Creating products...")
	productsCollection := db.Collection("products")
	categorySmartphones := 2
	categoryTablets := 3
	categoryLaptops := 4
	categoryAccessories := 5

	products := []interface{}{
		// Smartphones
		bson.M{"_id": 1, "name": "iPhone 15 Pro", "description": "Latest Apple flagship", "category_id": categorySmartphones, "price": 999.99, "stock": 100, "image_url": "https://store.storeimages.cdn-apple.com/4982/as-images.apple.com/is/iphone-15-pro-finish-select-202309-6-1inch-naturaltitanium?wid=5120&hei=2880&fmt=p-jpg&qlt=80&.v=1692895702708", "is_active": true, "created_at": time.Now(), "updated_at": time.Now()},
		bson.M{"_id": 2, "name": "Samsung Galaxy S24", "description": "Samsung flagship phone", "category_id": categorySmartphones, "price": 899.99, "stock": 80, "image_url": "https://images.unsplash.com/photo-1610945415295-d9bbf067e59c?w=800&auto=format&fit=crop", "is_active": true, "created_at": time.Now(), "updated_at": time.Now()},
		bson.M{"_id": 3, "name": "Google Pixel 8", "description": "Google's latest smartphone", "category_id": categorySmartphones, "price": 699.99, "stock": 60, "image_url": "https://images.unsplash.com/photo-1598327105666-5b89351aff97?w=800&auto=format&fit=crop", "is_active": true, "created_at": time.Now(), "updated_at": time.Now()},

		// Tablets
		bson.M{"_id": 4, "name": "iPad Pro 12.9", "description": "Apple's premium tablet", "category_id": categoryTablets, "price": 1099.99, "stock": 50, "image_url": "https://store.storeimages.cdn-apple.com/4982/as-images.apple.com/is/ipad-pro-13-select-wifi-spacegray-202405?wid=470&hei=556&fmt=png-alpha&.v=1713308272816", "is_active": true, "created_at": time.Now(), "updated_at": time.Now()},
		bson.M{"_id": 5, "name": "Samsung Galaxy Tab S9", "description": "Samsung premium tablet", "category_id": categoryTablets, "price": 849.99, "stock": 45, "image_url": "https://images.unsplash.com/photo-1544244015-0df4b3ffc6b0?w=800&auto=format&fit=crop", "is_active": true, "created_at": time.Now(), "updated_at": time.Now()},

		// Laptops
		bson.M{"_id": 6, "name": "MacBook Air M3", "description": "Apple M3, 8GB RAM, 256GB SSD", "category_id": categoryLaptops, "price": 1199.99, "stock": 30, "image_url": "https://store.storeimages.cdn-apple.com/4982/as-images.apple.com/is/macbook-air-midnight-select-20240708?wid=904&hei=840&fmt=jpeg&qlt=90&.v=1720728160304", "is_active": true, "created_at": time.Now(), "updated_at": time.Now()},
		bson.M{"_id": 7, "name": "MacBook Pro 16", "description": "Apple M3 Pro, 18GB RAM, 512GB SSD", "category_id": categoryLaptops, "price": 2499.99, "stock": 40, "image_url": "https://store.storeimages.cdn-apple.com/4982/as-images.apple.com/is/mbp16-spacegray-select-202310?wid=904&hei=840&fmt=jpeg&qlt=90&.v=1697311054290", "is_active": true, "created_at": time.Now(), "updated_at": time.Now()},
		bson.M{"_id": 8, "name": "Dell XPS 15", "description": "Intel i7, 16GB RAM, 512GB SSD", "category_id": categoryLaptops, "price": 1799.99, "stock": 60, "image_url": "https://images.unsplash.com/photo-1496181133206-80ce9b88a853?w=800&auto=format&fit=crop", "is_active": true, "created_at": time.Now(), "updated_at": time.Now()},

		// Accessories
		bson.M{"_id": 9, "name": "AirPods Pro", "description": "Apple wireless earbuds with ANC", "category_id": categoryAccessories, "price": 249.99, "stock": 150, "image_url": "https://store.storeimages.cdn-apple.com/4982/as-images.apple.com/is/MQD83?wid=1144&hei=1144&fmt=jpeg&qlt=90&.v=1660803972361", "is_active": true, "created_at": time.Now(), "updated_at": time.Now()},
		bson.M{"_id": 10, "name": "USB-C Hub", "description": "7-in-1 USB-C adapter", "category_id": categoryAccessories, "price": 49.99, "stock": 200, "image_url": "https://images.unsplash.com/photo-1625948515291-69613efd103f?w=800&auto=format&fit=crop", "is_active": true, "created_at": time.Now(), "updated_at": time.Now()},
	}
	_, err = productsCollection.InsertMany(ctx, products)
	if err != nil {
		log.Fatal("Failed to insert products:", err)
	}

	fmt.Println("✅ Database seeded successfully!")
	fmt.Println("\nDefault credentials:")
	fmt.Println("  Admin:     admin@example.com / password123")
	fmt.Println("  Moderator: moderator@example.com / password123")
	fmt.Println("  User:      user1@example.com / password123")
}
