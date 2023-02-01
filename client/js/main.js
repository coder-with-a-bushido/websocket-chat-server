const chatForm = document.getElementById('chat-form');
const chatMessages = document.querySelector('.chat-messages');
const roomName = document.getElementById('room-name');
const userList = document.getElementById('users');

// Get username from URL
const { username } = Qs.parse(location.search, {
  ignoreQueryPrefix: true,
});

let socket = new WebSocket(window.location.host == "localhost" ? "ws" : "wss" + "://" + window.location.host + "/ws")

socket.onopen = (event) => {
  sendJoinEvent();
}

socket.onmessage = (event) => {
  outputMessage(JSON.parse(event.data))
}

// Message submit
chatForm.addEventListener('submit', (e) => {
  e.preventDefault();

  // Get message text
  let msg = e.target.elements.msg.value;

  msg = msg.trim();

  if (!msg) {
    return false;
  }

  // Emit message to server
  sendMessageEvent(msg);

  // Clear input
  e.target.elements.msg.value = '';
  e.target.elements.msg.focus();

  outputMessage({
    name: "You",
    value: msg
  })

});

// Output message to DOM
function outputMessage(message) {
  const div = document.createElement('div');
  div.classList.add('message');
  const p = document.createElement('p');
  p.classList.add('meta');
  p.innerText = message.name;
  //p.innerHTML += `<span>${message.time}</span>`;
  div.appendChild(p);
  const para = document.createElement('p');
  para.classList.add('text');
  para.innerText = message.value;
  div.appendChild(para);
  document.querySelector('.chat-messages').appendChild(div);

  // Scroll Down
  chatMessages.scrollTop = chatMessages.scrollHeight;
}

// // Add room name to DOM
// function outputRoomName(room) {
//   roomName.innerText = room;
// }

// // Add users to DOM
// function outputUsers(users) {
//   console.log({ users })
//   userList.innerHTML = '';
//   users.forEach((user) => {
//     const li = document.createElement('li');
//     li.innerText = user.username;
//     userList.appendChild(li);
//   });
// }

//Prompt the user before leave chat room
document.getElementById('leave-btn').addEventListener('click', () => {
  const leaveRoom = confirm('Are you sure you want to leave the chatroom?');
  if (leaveRoom) {
    sendLeaveEvent()
    window.location = '../index.html';
  } else {
  }
});

function sendJoinEvent() {
  let joinEvent = {
    kind: 0, //PEERJOIN
    data: {
      username: username
    }
  }

  socket.send(JSON.stringify(joinEvent))
}

function sendMessageEvent(msg) {
  let messageEvent = {
    kind: 1, //PEERMESSAGE
    data: {
      content: msg
    }
  }

  socket.send(JSON.stringify(messageEvent))
}

function sendLeaveEvent() {
  let leaveEvent = {
    kind: 2, //PEERLEAVE
    data: null
  }

  socket.send(JSON.stringify(leaveEvent))
}