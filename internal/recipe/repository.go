package recipe

import (
	"context"
	"database/sql"
	"github.com/mamontmodest/go-rest-api/internal/entity"
	"github.com/mamontmodest/go-rest-api/pkg/db"
	errs "github.com/mamontmodest/go-rest-api/pkg/errors"
	"github.com/mamontmodest/go-rest-api/pkg/log"
	"hash/fnv"
	"time"
)

// Repository encapsulates the logic to access recipes from the data source.
type Repository interface {
	// GetRecipes returns the Recipes .
	GetRecipes(ctx context.Context) (entity.RecipesList, error)
	// Get returns the Recipe
	Get(ctx context.Context, recipeId int) (entity.Recipe, error)
	// Create saves a new Recipe in the storage.
	Create(ctx context.Context, recipe entity.Recipe) (entity.Recipe, error)
	// CreateWithId create the Recipe with given ID in the storage.
	CreateWithId(ctx context.Context, recipe entity.Recipe) error
	// Delete removes the Recipe with given ID from the storage.
	Delete(ctx context.Context, recipeId int) error
	// GetRecipesFilter returns the Recipes with filter
	GetRecipesFilter(ctx context.Context, filter entity.Filter) (entity.RecipesList, error)
	// К сожалению тут логин (
	Login(ctx context.Context, identity entity.Identity) error
}

// repository persists Recipe in database
type repository struct {
	db     *db.SDatabase
	logger *log.Logger
}

// NewRepository creates a new recipe repository
func NewRepository(db *db.SDatabase, logger *log.Logger) Repository {
	return repository{db, logger}
}

// GetRecipes returns the list of recipes recordset in the database.
func (r repository) GetRecipes(ctx context.Context) (entity.RecipesList, error) {
	recipes := new(entity.RecipesList)
	conn, err := r.db.ConnWith(ctx)
	if err != nil {
		return *recipes, err
	}
	defer conn.Close()
	ct, cancel := context.WithTimeout(ctx, time.Second*1)
	defer cancel()
	resp := make(chan ChanRecipesList)
	go func() {
		//select all tables and creates maps for data structs
		IngredientMap := make(map[int][]*entity.Ingredient)
		StepsMap := make(map[int][]*entity.Step)
		tx, err := conn.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
		defer tx.Rollback()
		query := "Select recipeId, ingredient  from ingredient where True"
		rows, err := tx.Query(query)
		if err != nil {
			resp <- ChanRecipesList{*recipes, err}
			return
		}
		for rows.Next() {
			ingredient := new(entity.Ingredient)
			recipeId := 0
			err := rows.Scan(&recipeId, &ingredient.Name)
			if err != nil {
				resp <- ChanRecipesList{*recipes, err}
				return
			}
			IngredientMap[recipeId] = append(IngredientMap[recipeId], ingredient)
		}
		query = "Select recipeid, stepnumber, description, timeduration  from step where True"
		rows, err = tx.Query(query)
		if err != nil {
			resp <- ChanRecipesList{*recipes, err}
			return
		}
		for rows.Next() {
			step := new(entity.Step)
			recipeId := 0
			err := rows.Scan(&recipeId, &step.StepNumber, &step.Description, &step.Time)
			if err != nil {
				resp <- ChanRecipesList{*recipes, err}
				return
			}
			StepsMap[recipeId] = append(StepsMap[recipeId], step)
		}
		query = "Select recipeid, name, title  from recipe where True"
		rows, err = tx.Query(query)
		if err != nil {
			resp <- ChanRecipesList{*recipes, err}
			return
		}
		for rows.Next() {
			recipe := new(entity.Recipe)
			err := rows.Scan(&recipe.RecipeId, &recipe.Name, &recipe.Title)
			if err != nil {
				resp <- ChanRecipesList{*recipes, err}
				return
			}
			recipe.Ingredients = IngredientMap[recipe.RecipeId]
			recipe.Steps = StepsMap[recipe.RecipeId]
			recipes.Recipes = append(recipes.Recipes, *recipe)
		}
		resp <- ChanRecipesList{*recipes, err}
		return
	}()
	for {
		select {
		case <-ct.Done():
			return entity.RecipesList{}, errs.CtxError{}
		case RecipeMessage := <-resp:
			return RecipeMessage.RecipesList, RecipeMessage.err
		}
	}
}

// GetRecipesFilter returns the list of recipes with filter recordset in the database.
func (r repository) GetRecipesFilter(ctx context.Context, filter entity.Filter) (entity.RecipesList, error) {
	recipes := new(entity.RecipesList)
	conn, err := r.db.ConnWith(ctx)
	if err != nil {
		return *recipes, err
	}
	defer conn.Close()
	ct, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	resp := make(chan ChanRecipesList)
	go func() {
		//select all tables and creates maps for data structs
		IngredientMap := make(map[int][]*entity.Ingredient)
		StepsMap := make(map[int][]*entity.Step)
		tx, err := conn.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
		defer tx.Rollback()
		query := "Select recipeId, ingredient  from ingredient where True"
		rows, err := tx.Query(query)
		if err != nil {
			resp <- ChanRecipesList{*recipes, err}
			return
		}
		for rows.Next() {
			ingredient := new(entity.Ingredient)
			recipeId := 0
			err := rows.Scan(&recipeId, &ingredient.Name)
			if err != nil {
				resp <- ChanRecipesList{*recipes, err}
				return
			}
			IngredientMap[recipeId] = append(IngredientMap[recipeId], ingredient)
		}
		query = "Select recipeid, stepnumber, description, timeduration  from step where True"
		rows, err = tx.Query(query)
		if err != nil {
			resp <- ChanRecipesList{*recipes, err}
			return
		}
		for rows.Next() {
			step := new(entity.Step)
			recipeId := 0
			err := rows.Scan(&recipeId, &step.StepNumber, &step.Description, &step.Time)
			if err != nil {
				resp <- ChanRecipesList{*recipes, err}
				return
			}
			StepsMap[recipeId] = append(StepsMap[recipeId], step)
		}
		query = "Select recipeid, name, title  from recipe where True"
		rows, err = tx.Query(query)
		if err != nil {
			resp <- ChanRecipesList{*recipes, err}
			return
		}
		for rows.Next() {
			recipe := new(entity.Recipe)
			err := rows.Scan(&recipe.RecipeId, &recipe.Name, &recipe.Title)
			if err != nil {
				resp <- ChanRecipesList{*recipes, err}
				return
			}
			recipe.Ingredients = IngredientMap[recipe.RecipeId]
			recipe.Steps = StepsMap[recipe.RecipeId]
			recipes.Recipes = append(recipes.Recipes, *recipe)
		}
		resp <- ChanRecipesList{*recipes, err}
		return
	}()
	for {
		select {
		case <-ct.Done():
			return entity.RecipesList{}, errs.CtxError{}
		case RecipeMessage := <-resp:
			return RecipeMessage.RecipesList, RecipeMessage.err
		}
	}
}

// Get reads the recipe with the specified ID from the database.
func (r repository) Get(ctx context.Context, id int) (entity.Recipe, error) {
	conn, err := r.db.ConnWith(ctx)
	if err != nil {
		return entity.Recipe{}, err
	}
	defer conn.Close()
	ct, cancel := context.WithTimeout(ctx, time.Second*1)
	defer cancel()
	resp := make(chan ChanRecipe)
	go func() {
		recipe := new(entity.Recipe)
		tx, err := conn.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
		defer tx.Rollback()
		query := "Select recipeId, name, title  from recipe where recipeId = $1"
		err = tx.QueryRow(query, id).Scan(&recipe.RecipeId, &recipe.Name, &recipe.Title)
		if err != nil {
			resp <- ChanRecipe{*recipe, err}
			return
		}
		query = "Select ingredient  from ingredient where recipeId = $1"
		rows, err := tx.Query(query, id)
		if err != nil {
			resp <- ChanRecipe{*recipe, err}
			return
		}
		for rows.Next() {
			ingredient := new(entity.Ingredient)
			rows.Scan(&ingredient.Name)
			recipe.Ingredients = append(recipe.Ingredients, ingredient)
		}
		query = "Select stepNumber, description, timeduration  from step where recipeId = $1"
		rows, err = tx.Query(query, id)
		if err != nil {
			resp <- ChanRecipe{*recipe, err}
			return
		}
		for rows.Next() {
			step := new(entity.Step)
			rows.Scan(&step.StepNumber, &step.Description, &step.Time)
			recipe.Steps = append(recipe.Steps, step)
		}
		resp <- ChanRecipe{*recipe, err}
		return
	}()
	for {
		select {
		case <-ct.Done():
			return entity.Recipe{}, errs.CtxError{}
		case RecipeMessage := <-resp:
			return RecipeMessage.Recipe, RecipeMessage.err
		}
	}
}

// Create saves a new recipe record in the database.
// It returns the ID of the newly inserted recipe record.
func (r repository) Create(ctx context.Context, recipe entity.Recipe) (entity.Recipe, error) {
	conn, err := r.db.ConnWith(ctx)
	if err != nil {
		return entity.Recipe{}, err
	}
	defer conn.Close()
	ct, cancel := context.WithTimeout(ctx, time.Second*1)
	defer cancel()
	resp := make(chan ChanRecipe)
	go func() {
		tx, err := conn.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
		defer tx.Rollback()
		query := "Insert into recipe (name, title) values($1, $2) returning recipeid"
		err = tx.QueryRow(query, recipe.Name, recipe.Title).Scan(&recipe.RecipeId)
		if err != nil {
			resp <- ChanRecipe{recipe, err}
			return
		}
		for _, v := range recipe.Ingredients {
			query := "Insert into ingredient (recipeid, ingredient) values ($1, $2)"
			rows, err := tx.Query(query, recipe.RecipeId, v.Name)
			if err != nil {
				resp <- ChanRecipe{recipe, err}
				return
			}
			err = rows.Close()
			if err != nil {
				resp <- ChanRecipe{recipe, err}
				return
			}
		}
		for _, v := range recipe.Steps {
			query := "INSERT INTO step (recipeid, stepnumber, description, timeduration) values ($1, $2, $3, $4)"
			rows, err := tx.Query(query, recipe.RecipeId, v.StepNumber, v.Description, v.Time)
			if err != nil {
				resp <- ChanRecipe{recipe, err}
				return
			}
			err = rows.Close()
			if err != nil {
				resp <- ChanRecipe{recipe, err}
				return
			}
		}
		err = tx.Commit()
		if err != nil {
			resp <- ChanRecipe{recipe, err}
			return
		}
		resp <- ChanRecipe{recipe, nil}
		return
	}()
	for {
		select {
		case <-ct.Done():
			return entity.Recipe{}, errs.CtxError{}
		case RecipeMessage := <-resp:
			return RecipeMessage.Recipe, RecipeMessage.err
		}
	}
}

// CreateWithId saves a new recipe record with id in the database.
func (r repository) CreateWithId(ctx context.Context, recipe entity.Recipe) error {
	conn, err := r.db.ConnWith(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()
	ct, cancel := context.WithTimeout(ctx, time.Second*1)
	defer cancel()
	resp := make(chan error)
	go func() {
		tx, err := conn.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
		defer tx.Rollback()
		query := "Insert into recipe (recipeid, name, title) values($1, $2, $3)"
		rows, err := tx.Query(query, recipe.RecipeId, recipe.Name, recipe.Title)
		if err != nil {
			resp <- err
			return
		}
		err = rows.Close()
		if err != nil {
			resp <- err
			return
		}
		for _, v := range recipe.Ingredients {
			query := "Insert into ingredient (recipeid, ingredient) values ($1, $2)"
			rows, err := tx.Query(query, recipe.RecipeId, v.Name)
			if err != nil {
				resp <- err
				return
			}
			err = rows.Close()
			if err != nil {
				resp <- err
				return
			}
		}
		for _, v := range recipe.Steps {
			query := "INSERT INTO step (recipeid, stepnumber, description, timeduration) values ($1, $2, $3, $4)"
			rows, err := tx.Query(query, recipe.RecipeId, v.StepNumber, v.Description, v.Time)
			if err != nil {
				resp <- err
				return
			}
			err = rows.Close()
			if err != nil {
				resp <- err
				return
			}
		}
		err = tx.Commit()
		if err != nil {
			resp <- err
			return
		}
		resp <- nil
		return
	}()
	for {
		select {
		case <-ct.Done():
			return errs.CtxError{}
		case err := <-resp:
			return err
		}
	}
}

// Delete deletes a recipe with the specified ID from the database.
func (r repository) Delete(ctx context.Context, id int) error {
	conn, err := r.db.ConnWith(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()
	ct, cancel := context.WithTimeout(ctx, time.Second*1)
	defer cancel()
	resp := make(chan error)
	go func() {
		tx, err := conn.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
		defer tx.Rollback()
		query := "delete from recipe where recipeid = $1"
		_, err = tx.Query(query, id)
		if err != nil {
			resp <- err
			return
		}
		tx.Commit()
		resp <- nil
	}()
	for {
		select {
		case <-ct.Done():
			return errs.CtxError{}
		case err := <-resp:
			return err
		}
	}
}

func (r repository) Login(ctx context.Context, identity entity.Identity) error {
	conn, err := r.db.ConnWith(ctx)
	person := new(entity.Identity)
	person.Login = identity.Login
	person.Password = identity.Password
	if err != nil {
		return err
	}
	defer conn.Close()
	ct, cancel := context.WithTimeout(ctx, time.Second*1)
	defer cancel()
	resp := make(chan error)
	go func() {
		h := hash(identity.Password)
		tx, err := conn.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
		defer tx.Rollback()
		query := "select login from  users where login = $1 and  password = $2"
		err = tx.QueryRow(query, identity.Login, h).Scan(&identity.Login)
		if err != nil {
			resp <- err
		}
		resp <- nil

	}()
	for {
		select {
		case <-ct.Done():
			return errs.CtxError{}
		case err := <-resp:
			return err
		}
	}
}
func hash(s string) uint32 {
	h := fnv.New32()
	h.Write([]byte(s))
	return h.Sum32()
}

type ChanRecipe struct {
	Recipe entity.Recipe
	err    error
}
type ChanRecipesList struct {
	RecipesList entity.RecipesList
	err         error
}
