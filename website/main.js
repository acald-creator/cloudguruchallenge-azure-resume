const localFuncApi = "http://localhost:7071/api/GetVisitorCounter"

function visitorCount() {
    let counter = document.getElementById('counter').innerText;
    let count = 0;
    count = parseInt(counter.innerText);
    counter.innerHTML = count;
}

window.onload = visitorCount();

const params = {
    method: 'GET',
    headers: {
        'accept': 'application/json'
    },
    mode: "no-cors",
};

let URL = "https://resume.acaldwell.dev/resume.json";

fetch(URL, params)
    .then(function (response) {
        return response.text()
    })
    .then(function (data) {
        Promise.resolve(data ? JSON.parse(data) : {})
    })
    .catch(function (error) {
        return console.log(error)
    })