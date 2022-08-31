package main

//importing diffrent libraries for web requesting and building an API
import (
	"fmt"
	"log"
	"net/http"

	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

//url of connecting to a database
const url = "postgres://myuser2:123@localhost:5432/mydb"

type DBConection struct {
	Conn *sqlx.DB
}

func ConnectDB() *DBConection {
	db, err := sqlx.Connect("postgres", url) //connecting to database
	if err != nil {
		log.Println("here is the err")
		panic(err)
	}
	return &DBConection{
		Conn: db,
	}
}

////defining structure
type Book struct {
	ID     int    `json:"id" db:"id"`
	Title  string `json:"title" db:"title"`
	Author string `json:"author" db:"author"`
}

//retriving the book by using id/
func GetBook(c echo.Context) error {
	id := c.Param("id")
	fmt.Println("id:", id)
	b := new(Book)

	if err := c.Bind(b); err != nil {
		return err

	}
	//
	//else return an eroor string
	return c.JSON(http.StatusOK, b)
}

//creating book
func CreateBook(c echo.Context) error {
	fmt.Println("here")

	b := new(Book)

	if err := c.Bind(b); err != nil {
		return err

	}

	dbconn := ConnectDB()
	tx := dbconn.Conn.MustBegin()

	_, err := tx.NamedQuery(`INSERT INTO "myschema".book (id,title ,author ) VALUES(:id, :title, :author)`, b)
	if err != nil {
		log.Fatal(err)
	}
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

	return c.JSON(http.StatusOK, "success")
}

//deleting the book by specific id

func DeleteBook(c echo.Context) error {

	dbconn := ConnectDB()

	requestId := c.Param("id")

	intId, err := strconv.Atoi(requestId)
	if err != nil {
		panic(err)
	}
	_, err = dbconn.Conn.Exec(`DELETE FROM  "myschema".book WHERE id=$1`, intId)

	if err != nil {
		fmt.Println(err)
		log.Fatal("exited")
	}

	return c.JSON(http.StatusOK, "deleted")
}

//update the book by index of map
func UpdateBook(c echo.Context) error {
	dbconn := ConnectDB()
	id := c.Param("id")

	b := new(Book)

	if err := c.Bind(b); err != nil {
		return err

	}

	intId, err1 := strconv.Atoi(id)
	if err1 != nil {
		panic(err1)
	}

	updateStmt := `update "myschema".book set "id"=$1 , "title"=$2, "author"=$3 where "id"=$4 `
	_, err := dbconn.Conn.Exec(updateStmt, b.ID, b.Title, b.Author, intId)
	if err != nil {
		fmt.Println("im herr ")
		log.Fatal("Error is ", err)
		// log.Fatal("exited")
	}

	return c.JSON(http.StatusOK, ":: Updated:: ")
}

//main function
func main() {

	e := echo.New()
	ConnectDB() //creating an echo variable for server communication
	//e.GET("/", GetBooks)                 //getting all the books
	e.GET("/books/:id", GetBook)         //grtting the specific books by the index
	e.POST("/books", CreateBook)         //create the book
	e.DELETE("/delete/:id", DeleteBook)  //delete the book by index
	e.PUT("/updatebook/:id", UpdateBook) //update the book by index by 4
	e.Logger.Fatal(e.Start(":1323"))     //request to server
}
