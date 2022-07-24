package datastore

const (
	GetAuthor     = "select id,first_name,last_name,dob,pen_name from Authors;"
	GetByIDAuthor = "select id,first_name,last_name,dob,pen_name from Authors where id=?"
	InsertAuthor  = "INSERT INTO Authors (first_name, last_name, dob, pen_name) VALUES (?,?,?,?);"
	UpdateAuthor  = "UPDATE Authors SET first_name = ? ,last_name = ? ,dob = ? ,pen_name = ?  WHERE id =?"
	DeleteAuthor  = "delete from Authors where id=?;"

	GetBook     = "select id,title,publication,publication_date,author_id from Books;"
	GetByIDBook = "select id,title,publication,publication_date,author_id from Books where id=?"
	InsertBook  = "INSERT INTO Books (title, publication, publication_date, author_id) VALUES (?,?,?,?);"
	UpdateBook  = "UPDATE Books SET title = ? ,publication = ? ,publication_date = ?,author_id=?  WHERE id =?"
	DeleteBook  = "delete from Books where id=?;"
)
