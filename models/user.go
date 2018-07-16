package models

import (
	//"github.com/asdine/storm/q"
	//"log"

	"gopkg.in/hlandau/passlib.v1"

	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type User struct {
	ID       int64  `json:"id" storm:"id,increment"`
	IsAdmin  bool   `json:"isadmin"`
	Name     string `json:"name" storm:"index"`
	Email    string `json:"email" storm:"index,unique"`
	Password string `json:"password,omitempty" storm:""`
	PwHash   string `json:"pwhash,omitempty"`
}

func (u *User) Delete() error {
	return DB.DeleteStruct(u)
}

func (u *User) Save() error {
	hash, err := passlib.Hash(u.Password)
	if err != nil {
		return err
	}
	u.Password = ""
	u.PwHash = hash
	return DB.Save(u)
}

func (u User) Validate() error {

	err := validation.Errors{
		"name": func() error {
			return validation.Validate(u.Name, validation.Required, validation.Min(4))
		}(),
		"email": func() error {
			return validation.Validate(u.Email, validation.Required, is.Email)
		}(),
		"password": func() error {
			return validation.Validate(u.PwHash, validation.Required, validation.Min(5))
		}(),
	}.Filter()
	return err

	/*return validation.ValidateStruct(&f,
		// Street cannot be empty, and the length must between 5 and 50
		validation.Field(&f.State, validation.Required),
		// City cannot be empty, and the length must between 5 and 50
		//validation.Field(&f.TOn, validation.Required, validation.Length(5, 50)),
		// State cannot be empty, and must be a string consisting of two letters in upper case
		validation.Field(&f.TOn, validation.Required, validation.Match(regexp.MustCompile("^[0-9]{2}\\:[0-9]{2}$"))),
		validation.Field(&f.TOff, validation.Required, validation.Match(regexp.MustCompile("^[0-9]{2}\\:[0-9]{2}$"))),
		// State cannot be empty, and must be a string consisting of five digits
		//validation.Field(&a.Zip, validation.Required, validation.Match(regexp.MustCompile("^[0-9]{5}$"))),
	)*/
}

func GetUserByEmail(email string) (User, error) {
	u := User{}
	err := DB.One("Email", email, &u)
	return u, err
}

func GetUserByID(id int64) (User, error) {
	u := User{}
	err := DB.One("ID", id, &u)
	return u, err
}

func NewUser(name string, email string, password string, isAdmin bool) error {
	u := User{}
	u.IsAdmin = isAdmin
	u.Name = name
	u.Email = email
	u.Password = password
	return u.Save()
}

func UserAuth(email, password string) bool {
	//return true
	// The username and password parameters comes from the request header,
	// make a database lookup to make sure the username/password pair exist
	// and return true if they do, false if they dont.

	// To keep this example simple, lets just hardcode "hello" and "world" as username,password
	u := User{}
	/*hash, err := passlib.Hash(password)
	//qry := DB.Select(q.Eq("Email", username)).Limit(1)
	//qry.Find(&u)
	if err != nil {
		// Hashing error
		return false
	}*/

	err := DB.One("Email", email, &u)
	if err != nil {
		// User not found
		//log.Println("User:", email, " was not found!")
		return false
	}
	newHash, err := passlib.Verify(password, u.PwHash)
	if err != nil {
		// incorrect password, malformed hash, etc.
		// either way, reject
		//log.Println("Incorrect password!")
		return false
	}

	// The context has decided, as per its policy, that
	// the hash which was used to validate the password
	// should be changed. It has upgraded the hash using
	// the verified password.
	if newHash != "" {
		//(store newHash in database, replacing old hash)
		u.PwHash = newHash
		DB.Save(&u)
	}
	//log.Println("User was successfully logged in!")
	return true
}
