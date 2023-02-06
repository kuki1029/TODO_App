/* This function will obtain all the parameters for the form and then send
   it to the appropriate controller function to be processed. It will also 
   update the page accordingly for the user. */
function signupButton() {
  // Obtain the values from the form and then get their value
  var name = document.getElementById("name").value
  var email = document.getElementById("email").value
  var password = document.getElementById("password").value
  var confirm_pass = document.getElementById("confirm_pass").value
  
  // Now we need to check if the entered information is valid or not
  // We know the email will be valid due to our validators in HTML
  // So, we need to make sure the passwords match
  if (password != confirm_pass) {
    // This will let the user know of the error
    window.alert("Passwords do not match. Please try again.")
  }
  // This else activated if the passwords match.
  else {
    let loginData = {
      Name: name,
      Email: email,
      Password: password
    }
    // First we create a variable to hold the data that will need to be fetched
    let fetchData = {
      method: 'POST',
      // The stringify converts a JS value to JSON string
      body: JSON.stringify(loginData),
      headers: new Headers({
        'Content-Type': 'application/json; charset=UTF-8'
      })
    }
    // Now we can fetch the data using the above variable.
    // Normally, fetch defaults to GET but we redefined it above
    fetch('/signup', fetchData)
      // This allows us to work with the data received from the fetch API call
      // We simple convert it back to JSON
      .then(resposne => {
        return resposne.json();
      })
      // Using the converted value, we can check if the controller function
      // was successful or not.
      .then(result => {
        if (result.success) {
          // If everything was okay, we redirect the user to the login page
          // so that they can signin with their new account and view tasks
          window.location.href = "/login.html";
        }
        else {
          window.alert(result.message)
        }
      })
  }

}


/* This function will allow the user to login. It will send the parameters to the
   appropriate controller functions through the fetch api. */
function loginButton() {
  // Obtain the values from the form and then get their value
  var email = document.getElementById("email").value
  var password = document.getElementById("password").value
  
  // Now we need to check if the entered information is valid or not
  // We know the email will be valid due to our validators in HTML
  // We need to verify that the information is in the database and the password is correct
  // First we create a variable to hold the data that will need to be fetched
  let loginData = {
    Email: email,
    Password: password
  }
  let fetchData = {
    method: 'POST',
    // The stringify converts a JS value to JSON string
    body: JSON.stringify(loginData),
    headers: new Headers({
      'Content-Type': 'application/json; charset=UTF-8'
    })
  }
  
  // Now we can fetch the data using the above variable.
  // Normally, fetch defaults to GET but we redefined it above
  fetch('/login', fetchData)
    // This allows us to work with the data received from the fetch API call
    // We simple convert it back to JSON
    .then(resposne => {
      return resposne.json();
    })
    // Using the converted value, we can check if the controller function
    // was successful or not.
    .then(result => {
      if (result.success) {
        
        // We also need to call the tasks method
        // This defaults to GET
        let fetchData2 = {
          body: JSON.stringify(loginData),
          headers: new Headers({
            'Content-Type': 'application/json; charset=UTF-8'
          })
        }
        fetch('/tasks', fetchData2)
        // If everything was okay, we redirect the user to the tasks page so
        // that they can view their tasks
        window.location.href = "/tasks";
      }
      else {
        window.alert(result.message)
      }
    })
  
}

