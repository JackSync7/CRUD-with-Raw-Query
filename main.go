package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"personal-web/connection"
	"strconv"
	"text/template"
	"time"

	"github.com/labstack/echo/v4"
)

type Template struct {
	templates *template.Template
}

type Blog struct {
	ID                                      int
	Title, Content, Author, PostDate, Image string
}

type Project struct {
	Title, Content, React, Python, Node, Golang, Duration, Waktu string
	StartDate, EndDate                                           time.Time
	Id                                                           int
	Tech                                                         []string
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {

	connection.DatabaseConnect()
	e := echo.New()

	// route statis untuk mengakses folder public
	e.Static("/public", "public") // /public

	// renderer
	t := &Template{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}

	e.Renderer = t

	// Routing
	e.GET("/", home)           //localhost:5000
	e.GET("/contact", contact) //localhost:5000/contact
	e.GET("/blog", blog)       //localhost:5000/blog
	e.GET("/myProject", myProject)
	e.POST("/addProject", addProject)
	e.GET("/deleteProject/:id", deleteProject)
	e.GET("/detailProject/:id", detailProject)
	e.GET("/editProject/:id", editProject)
	e.POST("/updateProject/:id", updateProject)

	fmt.Println("Server berjalan di port 5000")
	e.Logger.Fatal(e.Start("localhost:5000"))
}

func home(c echo.Context) error {
	data, _ := connection.Conn.Query(context.Background(), "SELECT id, title, content, tech, start_date, end_date, duration FROM tb_project")

	var result []Project
	for data.Next() {
		var each = Project{}

		err := data.Scan(&each.Id, &each.Title, &each.Content, &each.Tech, &each.StartDate, &each.EndDate, &each.Duration)

		if err != nil {
			fmt.Println(err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]string{"Message ": err.Error()})
		}
		result = append(result, each)
	}
	projects := map[string]interface{}{
		"Project": result,
	}
	return c.Render(http.StatusOK, "index.html", projects)
}

func contact(c echo.Context) error {
	return c.Render(http.StatusOK, "contact.html", nil)
}

func blog(c echo.Context) error {
	data, _ := connection.Conn.Query(context.Background(), "SELECT id, title, content, image, post_date FROM tb_blog")

	var result []Blog
	for data.Next() {
		var each = Blog{}

		err := data.Scan(&each.ID, &each.Title, &each.Content, &each.Image, &each.PostDate)
		if err != nil {
			fmt.Println(err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]string{"Message ": err.Error()})
		}

		result = append(result, each)

	}

	blogs := map[string]interface{}{
		"Blogs": result,
	}
	return c.Render(http.StatusOK, "blog.html", blogs)
}

func myProject(c echo.Context) error {

	return c.Render(http.StatusOK, "myProject.html", nil)
}

func addProject(c echo.Context) error {
	title := c.FormValue("name")
	startDate := c.FormValue("startDate")
	endDate := c.FormValue("endDate")
	content := c.FormValue("textArea")
	react := c.FormValue("react")
	python := c.FormValue("python")
	node := c.FormValue("node")
	golang := c.FormValue("golang")

	var techs [4]string
	techs[0] = react
	techs[1] = python
	techs[2] = node
	techs[3] = golang

	fmt.Print(techs)

	layout := "2006-01-02"
	t1, _ := time.Parse(layout, endDate)
	t2, _ := time.Parse(layout, startDate)

	diff := t1.Sub(t2)

	days := int(diff.Hours() / 24)
	months := int(diff.Hours() / 24 / 30)
	weeks := int(diff.Hours() / 24 / 7)
	years := int(diff.Hours() / 24 / 365)

	var Distance string
	if years > 0 {

		Distance = strconv.Itoa(years) + " Years Ago"
		fmt.Printf("ini tahun : %s --", Distance)
	} else if months > 0 {

		Distance = strconv.Itoa(months) + " Month Ago"
		fmt.Printf("ini bulan : %s --", Distance)
	} else if weeks > 0 {

		Distance = strconv.Itoa(weeks) + " Weeks Ago"
		fmt.Printf("ini minggu : %s --", Distance)
	} else if days > 0 {

		Distance = strconv.Itoa(days) + " Days Ago"
		fmt.Printf("ini hari : %s --", Distance)
	}

	_, err := connection.Conn.Exec(context.Background(), "INSERT INTO tb_project (title, content, start_date, end_date, tech, duration) VALUES ($1, $2, $3, $4, $5, $6)", title, content, t2, t1, techs, Distance)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"Message ": err.Error()})
	}

	return c.Redirect(http.StatusMovedPermanently, "/")
}

func deleteProject(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	_, err := connection.Conn.Exec(context.Background(), "DELETE FROM tb_project WHERE id = $1", id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"Message ": err.Error()})
	}

	return c.Redirect(http.StatusMovedPermanently, "/")
}

func detailProject(c echo.Context) error {

	id, _ := strconv.Atoi(c.Param("id"))

	var ProjectDetail = Project{}
	err := connection.Conn.QueryRow(context.Background(), "SELECT id, title, content, start_date, tech FROM tb_project where id = $1", id).Scan(&ProjectDetail.Id, &ProjectDetail.Title, &ProjectDetail.Content, &ProjectDetail.StartDate, &ProjectDetail.Tech)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"Message ": err.Error()})
	}

	ProjectDetail.Waktu = ProjectDetail.StartDate.Format("2006-01-02")

	detailProject := map[string]interface{}{
		"Project": ProjectDetail,
	}
	return c.Render(http.StatusOK, "detailProject.html", detailProject)
}

func editProject(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	var ProjectDetail = Project{}
	err := connection.Conn.QueryRow(context.Background(), "SELECT id, title, content, start_date, end_date from tb_project WHERE id = $1", id).Scan(&ProjectDetail.Id, &ProjectDetail.Title, &ProjectDetail.Content, &ProjectDetail.StartDate, &ProjectDetail.EndDate)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"Message ": err.Error()})
	}

	detailProject := map[string]interface{}{
		"Project": ProjectDetail,
	}
	return c.Render(http.StatusOK, "editProject.html", detailProject)
}
func updateProject(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	title := c.FormValue("name")
	startDate := c.FormValue("startDate")
	endDate := c.FormValue("endDate")
	content := c.FormValue("textArea")
	react := c.FormValue("react")
	python := c.FormValue("python")
	node := c.FormValue("node")
	golang := c.FormValue("golang")

	var techs [4]string
	techs[0] = react
	techs[1] = python
	techs[2] = node
	techs[3] = golang

	fmt.Print(techs)

	layout := "2006-01-02"
	t1, _ := time.Parse(layout, endDate)
	t2, _ := time.Parse(layout, startDate)

	diff := t1.Sub(t2)

	days := int(diff.Hours() / 24)
	months := int(diff.Hours() / 24 / 30)
	weeks := int(diff.Hours() / 24 / 7)
	years := int(diff.Hours() / 24 / 365)

	var Distance string
	if years > 0 {

		Distance = strconv.Itoa(years) + " Years Ago"
		fmt.Printf("ini tahun : %s --", Distance)
	} else if months > 0 {

		Distance = strconv.Itoa(months) + " Month Ago"
		fmt.Printf("ini bulan : %s --", Distance)
	} else if weeks > 0 {

		Distance = strconv.Itoa(weeks) + " Weeks Ago"
		fmt.Printf("ini minggu : %s --", Distance)
	} else if days > 0 {

		Distance = strconv.Itoa(days) + " Days Ago"
		fmt.Printf("ini hari : %s --", Distance)
	}

	_, err := connection.Conn.Exec(context.Background(), "UPDATE tb_project SET title = $1, content = $2, start_date = $3, end_date = $4, tech = $5, duration = $6 WHERE id = $7 ", title, content, t2, t1, techs, Distance, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"Message ": err.Error()})
	}

	return c.Redirect(http.StatusMovedPermanently, "/")
}
