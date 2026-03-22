package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"

	"github.com/PrimeraAizen/e-comm/config"
)

func main() {
	ctx := context.Background()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	pool, err := pgxpool.New(ctx, cfg.PG.URL)
	if err != nil {
		log.Fatal("Failed to connect to PostgreSQL:", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatal("Failed to ping PostgreSQL:", err)
	}

	fmt.Println("Connected to PostgreSQL.")

	hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	passwordHash := string(hash)

	// Clear existing data (order matters due to foreign keys)
	tables := []string{
		"cart_items", "product_reviews", "order_items", "orders",
		"user_product_likes", "user_product_views",
		"profiles", "user_roles", "users", "categories", "products", "roles",
	}
	fmt.Println("Clearing existing data...")
	for _, t := range tables {
		if _, err := pool.Exec(ctx, "TRUNCATE TABLE "+t+" CASCADE"); err != nil {
			log.Printf("Warning: could not truncate %s: %v", t, err)
		}
	}

	fmt.Println("Seeding roles...")
	roles := []struct{ name, description string }{
		{"admin", "System administrator"},
		{"user", "Regular user"},
		{"moderator", "Content moderator"},
		{"student", "Student user"},
		{"teacher", "Teacher user"},
	}
	roleIDs := make(map[string]string)
	for _, r := range roles {
		var id string
		err := pool.QueryRow(ctx,
			"INSERT INTO roles (name, description) VALUES ($1, $2) RETURNING id",
			r.name, r.description,
		).Scan(&id)
		if err != nil {
			log.Fatalf("Insert role %s: %v", r.name, err)
		}
		roleIDs[r.name] = id
	}

	fmt.Println("Seeding users...")
	users := []string{
		"admin@example.com",
		"moderator@example.com",
		"user1@example.com",
		"user2@example.com",
		"student@example.com",
		"teacher@example.com",
	}
	userIDs := make(map[string]string)
	for _, email := range users {
		var id string
		err := pool.QueryRow(ctx,
			"INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id",
			email, passwordHash,
		).Scan(&id)
		if err != nil {
			log.Fatalf("Insert user %s: %v", email, err)
		}
		userIDs[email] = id
	}

	fmt.Println("Assigning roles...")
	assignments := []struct{ email, role string }{
		{"admin@example.com", "admin"},
		{"admin@example.com", "user"},
		{"moderator@example.com", "moderator"},
		{"user1@example.com", "user"},
		{"user2@example.com", "user"},
		{"student@example.com", "student"},
		{"teacher@example.com", "teacher"},
	}
	for _, a := range assignments {
		if _, err := pool.Exec(ctx,
			"INSERT INTO user_roles (user_id, role_id) VALUES ($1, $2)",
			userIDs[a.email], roleIDs[a.role],
		); err != nil {
			log.Fatalf("Assign role %s to %s: %v", a.role, a.email, err)
		}
	}

	fmt.Println("Seeding categories...")
	var electronicsID string
	if err := pool.QueryRow(ctx,
		"INSERT INTO categories (name, description) VALUES ('Electronics', 'Electronic devices and accessories') RETURNING id",
	).Scan(&electronicsID); err != nil {
		log.Fatal("Insert Electronics:", err)
	}

	subCategories := []struct{ name, desc string }{
		{"Smartphones", "Mobile phones"},
		{"Tablets", "Tablet devices"},
		{"Laptops", "Notebook computers"},
		{"Accessories", "Tech accessories"},
	}
	catIDs := make(map[string]string)
	catIDs["Electronics"] = electronicsID
	for _, c := range subCategories {
		var id string
		if err := pool.QueryRow(ctx,
			"INSERT INTO categories (name, description, parent_id) VALUES ($1, $2, $3) RETURNING id",
			c.name, c.desc, electronicsID,
		).Scan(&id); err != nil {
			log.Fatalf("Insert category %s: %v", c.name, err)
		}
		catIDs[c.name] = id
	}

	fmt.Println("Seeding products...")
	products := []struct {
		name, desc, category, imageURL string
		price                          float64
		stock                          int
	}{
		{"iPhone 15 Pro", "Latest Apple flagship", "Smartphones", "https://via.placeholder.com/300x300?text=iPhone+15+Pro", 999.99, 100},
		{"Samsung Galaxy S24", "Samsung flagship phone", "Smartphones", "https://via.placeholder.com/300x300?text=Galaxy+S24", 899.99, 80},
		{"Google Pixel 8", "Google's latest smartphone", "Smartphones", "https://via.placeholder.com/300x300?text=Pixel+8", 699.99, 60},
		{"iPad Pro 12.9", "Apple's premium tablet", "Tablets", "https://via.placeholder.com/300x300?text=iPad+Pro", 1099.99, 50},
		{"Samsung Galaxy Tab S9", "Samsung premium tablet", "Tablets", "https://via.placeholder.com/300x300?text=Galaxy+Tab", 849.99, 45},
		{"MacBook Air M3", "Apple M3, 8GB RAM, 256GB SSD", "Laptops", "https://via.placeholder.com/300x300?text=MacBook+Air", 1199.99, 30},
		{"MacBook Pro 16", "Apple M3 Pro, 18GB RAM, 512GB SSD", "Laptops", "https://via.placeholder.com/300x300?text=MacBook+Pro", 2499.99, 40},
		{"Dell XPS 15", "Intel i7, 16GB RAM, 512GB SSD", "Laptops", "https://via.placeholder.com/300x300?text=Dell+XPS+15", 1799.99, 60},
		{"AirPods Pro", "Apple wireless earbuds with ANC", "Accessories", "https://via.placeholder.com/300x300?text=AirPods", 249.99, 150},
		{"USB-C Hub", "7-in-1 USB-C adapter", "Accessories", "https://via.placeholder.com/300x300?text=USB-C+Hub", 49.99, 200},
	}
	for _, p := range products {
		if _, err := pool.Exec(ctx,
			"INSERT INTO products (name, description, category_id, price, stock, image_url, is_active) VALUES ($1, $2, $3, $4, $5, $6, true)",
			p.name, p.desc, catIDs[p.category], p.price, p.stock, p.imageURL,
		); err != nil {
			log.Fatalf("Insert product %s: %v", p.name, err)
		}
	}

	fmt.Println("✅ Database seeded successfully!")
	fmt.Println("\nDefault credentials:")
	fmt.Println("  Admin:     admin@example.com / password123")
	fmt.Println("  Moderator: moderator@example.com / password123")
	fmt.Println("  User:      user1@example.com / password123")
}
