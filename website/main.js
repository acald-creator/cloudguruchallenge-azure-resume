const localFuncApi = "http://localhost:7071/api/GetVisitorCounter"

function visitorCount() {
  let counter = document.getElementById('counter').innerText;
  let count = 0;
  count = parseInt(counter.innerText);
  counter.innerHTML = count;
}

window.onload = visitorCount();

async function main() {
  fetch('https://resume.acaldwell.dev/resume.json', {
    method: 'GET',
    headers: {
      Accept: 'application/json',
    },
  })
    .then((response) => response.json())
    .then((response) => console.log(JSON.stringify(response)));
}

main();
