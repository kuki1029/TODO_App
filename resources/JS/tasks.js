

// Allows user to add tasks. This function will update the UI and also backend with the new task
function newElement() {
    li = document.createElement("li")
    task = document.getElementById("myInput").value
    // Creates a node so we can add the text to the list element
    text = document.createTextNode(task)
    li.appendChild(text)
    // If inputted text was empty, we let the user know
    if (task == '') {
        window.alert("You must write something!")
    }
    // Add the task to database through the fetch api
    let fetchData = {
        method: 'POST',
        // The stringify converts a JS value to JSON string
        body: JSON.stringify({TaskName: task}),
        headers: new Headers({
          'Content-Type': 'application/json; charset=UTF-8'
        })
      }
    // Now we can fetch the data using the above variable.
    // Normally, fetch defaults to GET but we redefined it above
    fetch('/tasks', fetchData)
    // This allows us to work with the data received from the fetch API call
    // We simple convert it back to JSON
    .then(resposne => {
        return resposne.json();
      })
      // Using the converted value, we can check if the controller function
      // was successful or not.
      .then(result => {
        if (result.success) {
            // If everything was okay, we can update the UI to show the new added task
            // Add the list element to the page
            document.getElementById("taskList").appendChild(li)
            // Set the input field to empty again
            document.getElementById("myInput").value = "";
            span = document.createElement("SPAN");
            i1 = document.createElement("i");
            i1.className = "bi bi-trash"
            span.className = "close";
            span.appendChild(i1);
            // Adds the ability to delete the task
            span.onclick = (function() {
              delElement(this)
            })
            span2 = document.createElement("SPAN");
            i2 = document.createElement("i");
            i2.className = "bi bi-pencil-square"
            span2.className = "edit"
            span2.appendChild(i2)
            span2.onclick = (function() {
              editButton(this)
            })
            li.appendChild(span);
            li.appendChild(span2)
            li.id = result.ID;
            

        }
        else {
          window.alert("There was an error with adding the task. Please try again.")
        }
      })

}

// This deletes the task in the backend by calling the appropriate function. Once deleted, the user cannot view the task.
// It also removes it from the frontend
function delElement(elem) {
  var div = elem.parentElement;
  div.style.display = "none";
  // We store the 
  var id = elem.parentElement.id
  // Add the task to database through the fetch api
  let fetchData = {
    method: 'DELETE',
    headers: new Headers({
      'Content-Type': 'application/json; charset=UTF-8'
    })
  }
  // Now we can fetch the data using the above variable.
  // Normally, fetch defaults to GET but we redefined it above
  fetch('/tasks/' + id, fetchData)
  // We simple convert it back to JSON
  .then(resposne => {
    return resposne.json();
  })
  // Using the converted value, we can check if the controller function
  // was successful or not.
  .then(result => {
    if (result.success) {
      // If deletion in database was successful, we can remove element from frontend
      var LI = document.getElementById(id)
      LI.parentNode.removeChild(LI);
    }
    else {
      window.alert("There was an error with deleting the task. Please try again.")
    }
  })
}

// This function will add the checked class to the list element when clicked.
// This will make it appear crossed out. This will also update the backend to 
// mark the task done.
function markElemDone(elem) {

  var id = elem.id
  // Add the task to database through the fetch api
  let fetchData = {
    method: 'POST',
    headers: new Headers({
      'Content-Type': 'application/json; charset=UTF-8'
    })
  }
  // Now we can fetch the data using the above variable.
  // Normally, fetch defaults to GET but we redefined it above
  fetch('/tasksDone/' + id, fetchData)
  // We simple convert it back to JSON
  .then(resposne => {
    return resposne.json();
  })
  // Using the converted value, we can check if the controller function
  // was successful or not.
  .then(result => {
    if (result.success) {
      // If deletion in database was successful, we can mark it done on the frontend
      elem.classList.toggle('checked');
    }
    else {
      window.alert("There was an error with marking this task as done. Please try again.")
    }
  })
}

// This function logs out the user
function logout() {
  let fetchData = {
    method: 'POST',
    headers: new Headers({
      'Content-Type': 'application/json; charset=UTF-8'
    })
  }
  // Now we can fetch the data using the above variable.
  // Normally, fetch defaults to GET but we redefined it above
  fetch('/logout', fetchData)
  .then(resposne => {
    return resposne.json();
  })
  // Using the converted value, we can check if the controller function
  // was successful or not.
  .then(result => {
    if (result.success) {
      window.location.href = "/";
    }
    else {
      window.alert("Could not logout. Please try again.")
    }
  })
}

// This function lets you edit the task
function editButton(elem) {
  let newTaskName = prompt("Please enter the new task name.", "Task Name");
  var id = elem.parentElement.id
  // Send the new name to the controller through the fetch api
  let fetchData = {
    method: 'POST',
    // The stringify converts a JS value to JSON string
    body: JSON.stringify({TaskName: newTaskName}),
    headers: new Headers({
      'Content-Type': 'application/json; charset=UTF-8'
    })
  }
  fetch('/tasksEdit/' + id, fetchData)
  // We simple convert it back to JSON
  .then(resposne => {
    return resposne.json();
  })
  // Using the converted value, we can check if the controller function
  // was successful or not.
  .then(result => {
    if (result.success) {
      // If deletion in database was successful, we can edit the task name on the frontend
      elem.parentElement.innerHTML = newTaskName
    }
    else {
      window.alert("There was an error with editing this task as done. Please try again.")
    }
  })

}

// This code allows the user to click on the x button to delete a task
var close = document.getElementsByClassName("close");
var i;
for (i = 0; i < close.length; i++) {
  close[i].onclick = (function() {
    delElement(this)
  })
}

// This code allows the user to click on the x button to delete a task
var edit = document.getElementsByClassName("edit");
var j;
for (j = 0; j < edit.length; j++) {
  edit[j].onclick = (function() {
    editButton(this)
  })
}

// This allows the user to mark a task as done. It crosses it out by adding the checked class to the list element
var list = document.querySelector('ul');
list.addEventListener('click', function(ev) {
  markElemDone(ev.target)
}, false);

