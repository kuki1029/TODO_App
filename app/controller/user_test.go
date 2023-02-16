package controller

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// )

// func TestLogin(t *testing.T) {
// 	email := "test@gmail.com"
// 	password := "password"

// 	values := map[string]string{"email": email, "password": password}
// 	jsonValue, _ := json.Marshal(values)

// 	req, err := http.Post("/login", bytes.NewBuffer(jsonValue))

// 	assert.Equal(t, err, nil)

// 	//var response map[string]string
// 	//json.Unmarshal(writer.Body.Bytes(), &response)
// 	//_, exists := response["jwt"]
// 	fmt.Println(req)
// 	//.Equal(t, true, exists)
// }
