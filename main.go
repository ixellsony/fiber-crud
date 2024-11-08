// main.go
package main

import (
	"database/sql"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	_ "modernc.org/sqlite"
)

type Item struct {
	ID    int
	Name  string
	Price float64
}

func main() {
	// Initialisation de la base de données
	db, err := sql.Open("sqlite", "items.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Création de la table si elle n'existe pas
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS items (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			price REAL NOT NULL
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	// Configuration de Fiber avec le moteur de template
	app := fiber.New(fiber.Config{
		Views: html.New("./views", ".html"),
	})

	// Route pour afficher tous les items
	app.Get("/", func(c *fiber.Ctx) error {
		var items []Item
		rows, err := db.Query("SELECT id, name, price FROM items")
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var item Item
			if err := rows.Scan(&item.ID, &item.Name, &item.Price); err != nil {
				return err
			}
			items = append(items, item)
		}

		return c.Render("index", fiber.Map{
			"Title": "CRUD Example",
			"Items": items,
		})
	})

	// Route pour ajouter un item
	app.Post("/items", func(c *fiber.Ctx) error {
		name := c.FormValue("name")
		price, err := strconv.ParseFloat(c.FormValue("price"), 64)
		if err != nil {
			return err
		}

		_, err = db.Exec("INSERT INTO items (name, price) VALUES (?, ?)", name, price)
		if err != nil {
			return err
		}

		return c.Redirect("/")
	})

	// Route pour mettre à jour un item
	app.Post("/items/:id", func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return err
		}

		name := c.FormValue("name")
		price, err := strconv.ParseFloat(c.FormValue("price"), 64)
		if err != nil {
			return err
		}

		_, err = db.Exec("UPDATE items SET name = ?, price = ? WHERE id = ?", name, price, id)
		if err != nil {
			return err
		}

		return c.Redirect("/")
	})

	// Route pour supprimer un item
	app.Delete("/items/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		_, err := db.Exec("DELETE FROM items WHERE id = ?", id)
		if err != nil {
			return err
		}

		return c.SendStatus(200)
	})

	log.Fatal(app.Listen(":3000"))
}
