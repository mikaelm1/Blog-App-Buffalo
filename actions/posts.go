package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/mikaelm1/Blog-App-Buffalo/models"
	"github.com/pkg/errors"
)

// PostsIndex default implementation.
func PostsIndex(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	posts := &models.Posts{}
	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	q := tx.PaginateFromParams(c.Params())
	// Retrieve all Posts from the DB
	if err := q.All(posts); err != nil {
		return errors.WithStack(err)
	}
	// Make posts available inside the html template
	c.Set("posts", posts)
	// Add the paginator to the context so it can be used in the template.
	c.Set("pagination", q.Paginator)
	return c.Render(200, r.HTML("posts/index.html"))
}

// PostsCreateGet displays a form to create a post
func PostsCreateGet(c buffalo.Context) error {
	c.Set("post", &models.Post{})
	return c.Render(200, r.HTML("posts/create"))
}

// PostsCreatePost adds a new Post to the db
func PostsCreatePost(c buffalo.Context) error {
	// Allocate an empty Post
	post := &models.Post{}
	user := c.Value("current_user").(*models.User)
	// Bind post to the html form elements
	if err := c.Bind(post); err != nil {
		return errors.WithStack(err)
	}
	// Get the DB connection from the context
	tx := c.Value("tx").(*pop.Connection)
	// Validate the data from the html form
	post.AuthorID = user.ID
	verrs, err := tx.ValidateAndCreate(post)
	if err != nil {
		return errors.WithStack(err)
	}
	if verrs.HasAny() {
		c.Set("post", post)
		c.Set("errors", verrs.Errors)
		return c.Render(422, r.HTML("posts/create"))
	}
	// If there are no errors set a success message
	c.Flash().Add("success", "New post added successfully.")
	// and redirect to the index page
	return c.Redirect(302, "/")
}

// PostsDetail displays a single post.
func PostsDetail(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	post := &models.Post{}
	if err := tx.Find(post, c.Param("pid")); err != nil {
		return c.Error(404, err)
	}
	author := &models.User{}
	if err := tx.Find(author, post.AuthorID); err != nil {
		return c.Error(404, err)
	}
	c.Set("post", post)
	c.Set("author", author)
	comment := &models.Comment{}
	c.Set("comment", comment)
	comments := models.Comments{}
	if err := tx.BelongsTo(post).All(&comments); err != nil {
		return errors.WithStack(err)
	}
	for i := 0; i < len(comments); i++ {
		u := models.User{}
		if err := tx.Find(&u, comments[i].AuthorID); err != nil {
			return c.Error(404, err)
		}
		comments[i].Author = u
	}
	c.Set("comments", comments)
	return c.Render(200, r.HTML("posts/detail"))
}

// PostsEditGet displays a form to edit the post.
func PostsEditGet(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	post := &models.Post{}
	if err := tx.Find(post, c.Param("pid")); err != nil {
		return c.Error(404, err)
	}
	c.Set("post", post)
	return c.Render(200, r.HTML("posts/edit.html"))
}

// PostsEditPost updates a post.
func PostsEditPost(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	post := &models.Post{}
	if err := tx.Find(post, c.Param("pid")); err != nil {
		return c.Error(404, err)
	}
	if err := c.Bind(post); err != nil {
		return errors.WithStack(err)
	}
	verrs, err := tx.ValidateAndUpdate(post)
	if err != nil {
		return errors.WithStack(err)
	}
	if verrs.HasAny() {
		c.Set("post", post)
		c.Set("errors", verrs.Errors)
		return c.Render(422, r.HTML("posts/edit.html"))
	}
	c.Flash().Add("success", "Post was updated successfully.")
	return c.Redirect(302, "/posts/detail/%s", post.ID)
}

// PostsDelete deletes a post from the db
func PostsDelete(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	post := &models.Post{}
	if err := tx.Find(post, c.Param("pid")); err != nil {
		return c.Error(404, err)
	}
	if err := tx.Destroy(post); err != nil {
		return errors.WithStack(err)
	}
	c.Flash().Add("success", "Post was successfully deleted.")
	return c.Redirect(302, "/posts/index")
}
